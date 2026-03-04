# Jet Sharing — Database Schema Reference

> Order service database schema (PostgreSQL). Used by the backend for ride lifecycle, payments, and scooter management.

## Tables

### customers
| Column | Type | Description |
|--------|------|-------------|
| `customer_id` | UUID (PK) | Unique customer identifier |
| `name` | VARCHAR(120) | Full name |
| `phone` | VARCHAR(20) | Phone number (Azerbaijan format: +994-XX-XXX-XX-XX) |
| `email` | VARCHAR(120) | Email address |
| `created_at` | TIMESTAMPTZ | Account creation timestamp |
| `status` | VARCHAR(20) | `active`, `suspended`, or `deleted` |

### scooters
| Column | Type | Description |
|--------|------|-------------|
| `scooter_id` | VARCHAR(40) (PK) | Scooter identifier (e.g., `scooter-7b3f2e91` or `S.006068`) |
| `model` | VARCHAR(60) | Model name: `Jet Pro X1`, `Jet City S2`, `Jet Max Z3` |
| `battery_level` | INT (0-100) | Current battery percentage |
| `lat` | DOUBLE PRECISION | Current latitude |
| `lng` | DOUBLE PRECISION | Current longitude |
| `status` | VARCHAR(20) | `available`, `in_use`, `maintenance`, `offline` |
| `geo_cluster` | VARCHAR(40) | Geographic cluster for spatial indexing. **Must not be NULL** — scooters with NULL geo_cluster are invisible to the routing service |
| `last_seen` | TIMESTAMPTZ | Last telemetry update |

### rides
| Column | Type | Description |
|--------|------|-------------|
| `ride_id` | UUID (PK) | Unique ride identifier |
| `customer_id` | UUID (FK → customers) | Who took the ride |
| `scooter_id` | VARCHAR(40) (FK → scooters) | Which scooter was used |
| `start_time` | TIMESTAMPTZ | Ride start |
| `end_time` | TIMESTAMPTZ | Ride end (NULL if active) |
| `start_lat` / `start_lng` | DOUBLE PRECISION | Starting coordinates |
| `end_lat` / `end_lng` | DOUBLE PRECISION | Ending coordinates (NULL if active) |
| `distance_km` | NUMERIC(6,2) | Total distance (NULL if active) |
| `status` | VARCHAR(20) | `active`, `completed`, `cancelled` |

### orders
| Column | Type | Description |
|--------|------|-------------|
| `order_id` | UUID (PK) | Unique order identifier |
| `ride_id` | UUID (FK → rides) | Associated ride |
| `customer_id` | UUID (FK → customers) | Customer who placed the order |
| `amount_cents` | INT | Total amount in cents (e.g., 450 = $4.50) |
| `currency` | VARCHAR(3) | Currency code (default: `USD`) |
| `status` | INT | Current status code (see Order Status Codes below) |
| `created_at` | TIMESTAMPTZ | Order creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last status update |

### order_status_history
| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL (PK) | Auto-incrementing ID |
| `order_id` | UUID (FK → orders) | Which order |
| `old_status` | INT | Previous status (NULL for first entry) |
| `new_status` | INT | New status |
| `changed_at` | TIMESTAMPTZ | When the transition occurred |
| `reason` | TEXT | Human-readable reason |

### payments
| Column | Type | Description |
|--------|------|-------------|
| `payment_id` | UUID (PK) | Unique payment identifier |
| `order_id` | UUID (FK → orders) | Which order |
| `amount_cents` | INT | Payment amount in cents |
| `currency` | VARCHAR(3) | Currency code |
| `method` | VARCHAR(30) | `card`, `apple_pay`, `google_pay` |
| `status` | VARCHAR(20) | `pending`, `processing`, `completed`, `failed`, `refunded` |
| `provider_ref` | VARCHAR(120) | Stripe payment intent reference |
| `created_at` | TIMESTAMPTZ | When payment was initiated |
| `completed_at` | TIMESTAMPTZ | When finalized (NULL if pending) |

### payment_webhooks
| Column | Type | Description |
|--------|------|-------------|
| `webhook_id` | UUID (PK) | Unique webhook event ID |
| `payment_id` | UUID (FK → payments) | Associated payment |
| `event_type` | VARCHAR(60) | `payment_initiated`, `payment_processing`, `payment_completed` |
| `payload` | JSONB | Raw webhook payload from Stripe |
| `received_at` | TIMESTAMPTZ | When received |
| `processed_at` | TIMESTAMPTZ | When processed |

## Order Status Codes

| Code | Name | Description |
|------|------|-------------|
| 1 | `pending` | Order created, awaiting processing |
| 2 | `created` | Order confirmed in system |
| 3 | `scooter_assigned` | A scooter has been reserved |
| 4 | `route_calculated` | Navigation route is ready |
| 5 | `ride_active` | Customer is currently riding |
| 6 | `ride_paused` | Ride temporarily paused |
| 7 | `ride_ending` | Ride is wrapping up |
| 10 | `payment_pending` | Awaiting payment initiation |
| 11 | `payment_processing` | Payment in flight |
| 12 | `payment_completed` | Payment successful |
| 13 | `payment_failed` | Payment declined |
| 14 | `completed` | Order fully completed |
| 15 | `cancelled` | Order was cancelled |
| 16 | `refunded` | Payment was refunded |

> Codes 8-9 are reserved for future ride states. Codes jump from 7 to 10 for the payment phase.

## Relationships

```
customers ──┐
             ├──> rides ──> orders ──> order_status_history
scooters ───┘                  │
                               └──> payments ──> payment_webhooks
```

## Common Investigation Queries

```sql
-- Trace an order end-to-end
SELECT o.order_id, c.name, s.scooter_id, r.distance_km,
       o.amount_cents, o.status, p.status as pay_status
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
JOIN rides r ON o.ride_id = r.ride_id
JOIN scooters s ON r.scooter_id = s.scooter_id
LEFT JOIN payments p ON p.order_id = o.order_id
WHERE o.order_id::text LIKE '%<partial-id>%';

-- Order status timeline
SELECT old_status, new_status, reason, changed_at
FROM order_status_history
WHERE order_id = '<full-uuid>'
ORDER BY changed_at;

-- Find scooters with missing geo cluster (invisible to routing)
SELECT scooter_id, model, status, geo_cluster
FROM scooters
WHERE geo_cluster IS NULL;
```

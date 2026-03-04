# Jet Sharing - Orders Database Schema

> Fallback reference for participants who cannot connect via MCP.
> This document mirrors the live Supabase database used in Workshop 3.

## Table Overview

```
customers ──┐
             ├──> rides ──> orders ──> order_status_history
scooters ───┘                  │
                               ├──> payments ──> payment_webhooks
                               │
```

## Table Definitions

### customers

| Column | Type | Description |
|--------|------|-------------|
| `customer_id` | UUID (PK) | Unique customer identifier |
| `name` | VARCHAR(120) | Full name |
| `phone` | VARCHAR(20) | Phone number (Azerbaijan format) |
| `email` | VARCHAR(120) | Email address |
| `created_at` | TIMESTAMPTZ | Account creation timestamp |
| `status` | VARCHAR(20) | `active`, `suspended`, or `deleted` |

### scooters

| Column | Type | Description |
|--------|------|-------------|
| `scooter_id` | VARCHAR(40) (PK) | Scooter identifier (e.g., `scooter-7b3f2e91`) |
| `model` | VARCHAR(60) | Model name (Jet Pro X1, Jet City S2, etc.) |
| `battery_level` | INT (0-100) | Current battery percentage |
| `lat` | DOUBLE PRECISION | Current latitude |
| `lng` | DOUBLE PRECISION | Current longitude |
| `status` | VARCHAR(20) | `available`, `in_use`, `maintenance`, `offline` |
| `geo_cluster` | VARCHAR(40) | Geographic cluster assignment (can be NULL!) |
| `last_seen` | TIMESTAMPTZ | Last telemetry update |

> **Important**: `geo_cluster = NULL` means the scooter is not indexed in the geospatial lookup. This is the root cause of "transport not found" errors - the routing service skips scooters without a geo_cluster.

### rides

| Column | Type | Description |
|--------|------|-------------|
| `ride_id` | UUID (PK) | Unique ride identifier |
| `customer_id` | UUID (FK -> customers) | Who took the ride |
| `scooter_id` | VARCHAR(40) (FK -> scooters) | Which scooter was used |
| `start_time` | TIMESTAMPTZ | Ride start timestamp |
| `end_time` | TIMESTAMPTZ | Ride end (NULL if still active) |
| `start_lat` | DOUBLE PRECISION | Starting latitude |
| `start_lng` | DOUBLE PRECISION | Starting longitude |
| `end_lat` | DOUBLE PRECISION | Ending latitude (NULL if active) |
| `end_lng` | DOUBLE PRECISION | Ending longitude (NULL if active) |
| `distance_km` | NUMERIC(6,2) | Total distance (NULL if active) |
| `status` | VARCHAR(20) | `active`, `completed`, `cancelled` |

### orders

| Column | Type | Description |
|--------|------|-------------|
| `order_id` | UUID (PK) | Unique order identifier |
| `ride_id` | UUID (FK -> rides) | Associated ride |
| `customer_id` | UUID (FK -> customers) | Customer who placed the order |
| `amount_cents` | INT | Total amount in cents (e.g., 450 = $4.50) |
| `currency` | VARCHAR(3) | Currency code (default: `USD`) |
| `status` | INT | Current status code (see table below) |
| `created_at` | TIMESTAMPTZ | Order creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last status update |

### order_status_history

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL (PK) | Auto-incrementing ID |
| `order_id` | UUID (FK -> orders) | Which order this event belongs to |
| `old_status` | INT | Previous status (NULL for first entry) |
| `new_status` | INT | New status |
| `changed_at` | TIMESTAMPTZ | When the transition occurred |
| `reason` | TEXT | Human-readable reason for the change |

### payments

| Column | Type | Description |
|--------|------|-------------|
| `payment_id` | UUID (PK) | Unique payment identifier |
| `order_id` | UUID (FK -> orders) | Which order this payment is for |
| `amount_cents` | INT | Payment amount in cents |
| `currency` | VARCHAR(3) | Currency code |
| `method` | VARCHAR(30) | `card`, `apple_pay`, `google_pay` |
| `status` | VARCHAR(20) | `pending`, `processing`, `completed`, `failed`, `refunded` |
| `provider_ref` | VARCHAR(120) | External payment provider reference (e.g., Stripe PI) |
| `created_at` | TIMESTAMPTZ | When payment was initiated |
| `completed_at` | TIMESTAMPTZ | When payment was finalized (NULL if pending) |

### payment_webhooks

| Column | Type | Description |
|--------|------|-------------|
| `webhook_id` | UUID (PK) | Unique webhook event ID |
| `payment_id` | UUID (FK -> payments) | Associated payment |
| `event_type` | VARCHAR(60) | Event type: `payment_initiated`, `payment_processing`, `payment_completed` |
| `payload` | JSONB | Raw webhook payload from payment provider |
| `received_at` | TIMESTAMPTZ | When the webhook was received |
| `processed_at` | TIMESTAMPTZ | When our system processed it |

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
| 11 | `payment_processing` | Payment is being processed |
| 12 | `payment_completed` | Payment successful |
| 13 | `payment_failed` | Payment was declined |
| 14 | `completed` | Order fully completed |
| 15 | `cancelled` | Order was cancelled |
| 16 | `refunded` | Payment was refunded |

> **Note the gap**: Codes 8-9 are reserved for future ride states. Codes jump from 7 to 10 for the payment phase.

## Relationships Diagram

```
                    ┌──────────────┐
                    │  customers   │
                    │  (customer_id)│
                    └──────┬───────┘
                           │
              ┌────────────┴────────────┐
              │                         │
       ┌──────▼───────┐         ┌──────▼───────┐
       │    rides      │         │   orders     │
       │  (ride_id)    │◄────────│  (order_id)  │
       └──────┬────────┘         └──────┬───────┘
              │                         │
       ┌──────▼───────┐    ┌───────────┼───────────┐
       │   scooters   │    │           │           │
       │ (scooter_id) │    ▼           ▼           ▼
       └──────────────┘  order_     payments   (customer_id
                         status_    (payment_    FK back to
                         history      _id)      customers)
                         (id)          │
                                       ▼
                                  payment_
                                  webhooks
                                  (webhook_id)
```

## Example: Tracing a Single Order

For order `a82997e7`, here is how the tables connect:

1. **Customer** `ef15c6de...` (Alex Petrov) requests a ride
2. **Scooter** `scooter-7b3f2e91` (Jet Pro X1) is assigned
3. **Ride** `b93aa8f8...` is created linking customer + scooter
4. **Order** `a82997e7...` tracks the business lifecycle (status transitions)
5. **Status history** records each state change with timestamps
6. **Payment** `c04bb9f9...` is created when ride ends
7. **Webhooks** record raw events from Stripe as the payment processes

```sql
-- Trace an order end-to-end
SELECT o.order_id, c.name, s.scooter_id, r.distance_km,
       o.amount_cents, o.status, p.status as pay_status
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
JOIN rides r ON o.ride_id = r.ride_id
JOIN scooters s ON r.scooter_id = s.scooter_id
LEFT JOIN payments p ON p.order_id = o.order_id
WHERE o.order_id = 'a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d';
```

# Jet Sharing - Sample Query Results

> Pre-captured query results for Workshop 3 demos.
> Use this if MCP/Supabase connection is unavailable during the workshop.

## 1. Incident Order (a82997e7)

### orders table

```sql
SELECT order_id, ride_id, customer_id, amount_cents, currency, status, created_at, updated_at
FROM orders WHERE order_id = 'a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d';
```

| order_id | ride_id | customer_id | amount_cents | currency | status | created_at | updated_at |
|----------|---------|-------------|-------------|----------|--------|------------|------------|
| a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d | b93aa8f8-2c4e-5a9b-0d3f-7e6a5b4c3d2e | ef15c6de-22d7-4183-88e4-33682cad8380 | 450 | USD | 14 | 2026-02-27 15:25:03+00 | 2026-02-27 15:47:44+00 |

### order_status_history for a82997e7

```sql
SELECT id, old_status, new_status, changed_at, reason
FROM order_status_history
WHERE order_id = 'a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d'
ORDER BY changed_at;
```

| id | old_status | new_status | changed_at | reason |
|----|-----------|------------|------------|--------|
| 1 | NULL | 2 | 2026-02-27 15:25:03+00 | Order created |
| 2 | 2 | 3 | 2026-02-27 15:25:45+00 | Scooter scooter-7b3f2e91 assigned |
| 3 | 3 | 4 | 2026-02-27 15:26:12+00 | Route calculated: 2.3 km |
| 4 | 4 | 5 | 2026-02-27 15:42:18+00 | Ride started |
| 5 | 5 | 14 | 2026-02-27 15:47:44+00 | Ride completed, payment processed |

> **Key anomaly**: Look at the gap between status 4 (route_calculated at 15:26:12) and status 5 (ride_active at 15:42:18) - that is a **16-minute delay**. The route was calculated but the ride did not start for 16 minutes. Then status 5 to 14 was only 5.5 minutes (the actual ride).

### Payment for a82997e7

```sql
SELECT payment_id, amount_cents, currency, method, status, provider_ref, created_at, completed_at
FROM payments
WHERE order_id = 'a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d';
```

| payment_id | amount_cents | currency | method | status | provider_ref | created_at | completed_at |
|------------|-------------|----------|--------|--------|--------------|------------|--------------|
| c04bb9f9-3d5f-6a0c-1e4a-8f7a6c5d4e3f | 450 | USD | card | completed | stripe_pi_3Nk2Abc123def456 | 2026-02-27 15:47:30+00 | 2026-02-27 15:47:42+00 |

### Payment Webhooks for a82997e7

```sql
SELECT webhook_id, event_type, payload, received_at, processed_at
FROM payment_webhooks
WHERE payment_id = 'c04bb9f9-3d5f-6a0c-1e4a-8f7a6c5d4e3f'
ORDER BY received_at;
```

| event_type | payload (summary) | received_at | processed_at |
|------------|-------------------|-------------|--------------|
| payment_initiated | `{"event":"payment_initiated","amount":450,"currency":"usd","payment_method":"card_visa_4242","idempotency_key":"idk_a82997e7_001"}` | 2026-02-27 15:47:30+00 | 2026-02-27 15:47:30+00 |
| payment_processing | `{"event":"payment_processing","provider":"stripe","provider_ref":"stripe_pi_3Nk2Abc123def456","risk_score":12}` | 2026-02-27 15:47:35+00 | 2026-02-27 15:47:35+00 |
| payment_completed | `{"event":"payment_completed","provider_ref":"stripe_pi_3Nk2Abc123def456","net_amount":421,"fee":29,"currency":"usd"}` | 2026-02-27 15:47:42+00 | 2026-02-27 15:47:42+00 |

> Payment processed normally: 12 seconds from initiation to completion. Stripe fee was $0.29, net $4.21.

## 2. The "Transport Not Found" Scooter

```sql
SELECT scooter_id, model, battery_level, lat, lng, status, geo_cluster, last_seen
FROM scooters WHERE scooter_id = 'S.006068';
```

| scooter_id | model | battery_level | lat | lng | status | geo_cluster | last_seen |
|------------|-------|--------------|-----|-----|--------|-------------|-----------|
| S.006068 | Jet City S2 | 85 | 40.4021 | 49.8592 | available | **NULL** | 2026-02-27 14:30:00+00 |

> **Root cause**: `geo_cluster` is NULL. The routing service uses `WHERE geo_cluster IS NOT NULL` in its spatial queries, so this scooter is invisible to the assignment algorithm despite being available and having a valid location.

### Compare with a normal scooter

```sql
SELECT scooter_id, model, battery_level, status, geo_cluster
FROM scooters WHERE scooter_id = 'scooter-7b3f2e91';
```

| scooter_id | model | battery_level | status | geo_cluster |
|------------|-------|--------------|--------|-------------|
| scooter-7b3f2e91 | Jet Pro X1 | 72 | available | baku-central-03 |

> This scooter has a proper `geo_cluster` value and is correctly indexed.

## 3. Normal Orders for Comparison

```sql
SELECT o.order_id, c.name, o.amount_cents, o.status,
       o.created_at, o.updated_at,
       EXTRACT(EPOCH FROM (o.updated_at - o.created_at))/60 AS duration_min
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
WHERE o.status = 14
ORDER BY o.created_at DESC
LIMIT 5;
```

| order_id (short) | customer | amount_cents | status | created_at | updated_at | duration_min |
|-----------------|----------|-------------|--------|------------|------------|-------------|
| a82997e7... | Alex Petrov | 450 | 14 | 2026-02-27 15:25:03 | 2026-02-27 15:47:44 | 22.7 |
| 3f8a1b2c... | Leyla Mammadov | 890 | 14 | 2026-02-27 14:10:22 | 2026-02-27 14:42:15 | 31.9 |
| 7d4e5f6a... | Kamran Hasanov | 350 | 14 | 2026-02-27 12:05:11 | 2026-02-27 12:18:33 | 13.4 |
| 9c2b3a4d... | Nigar Guliyev | 1250 | 14 | 2026-02-27 09:30:45 | 2026-02-27 10:12:08 | 41.4 |
| 1e6f7g8h... | Farid Ismailov | 520 | 14 | 2026-02-26 22:15:33 | 2026-02-26 22:38:19 | 22.8 |

> Typical completed orders show total duration (created to completed) proportional to ride length.

## 4. Cancelled and Failed Orders

```sql
SELECT o.order_id, o.status, osh.reason, osh.changed_at
FROM orders o
JOIN order_status_history osh ON o.order_id = osh.order_id
WHERE o.status IN (13, 15)
  AND osh.new_status = o.status
ORDER BY osh.changed_at DESC
LIMIT 5;
```

| order_id (short) | status | reason | changed_at |
|-----------------|--------|--------|------------|
| 5a3b2c1d... | 15 | User cancelled | 2026-02-27 16:20:11 |
| 8d7e6f5g... | 15 | Timeout | 2026-02-27 11:45:33 |
| 2b1c0d9e... | 15 | Scooter unavailable | 2026-02-26 19:22:07 |
| 4f3e2d1c... | 13 | Card declined - insufficient funds | 2026-02-26 17:55:41 |
| 6h5g4f3e... | 15 | User cancelled | 2026-02-26 14:10:28 |

## 5. Analytical Query: Orders with Long Status 5 Duration

This is the kind of investigation query participants will build during the workshop:

```sql
-- Find orders where status 5 (ride_active) lasted more than 3 minutes
-- by comparing when status 5 was entered vs when it transitioned out

SELECT
    o.order_id,
    c.name AS customer,
    entry.changed_at AS ride_started,
    exit_h.changed_at AS ride_ended,
    EXTRACT(EPOCH FROM (exit_h.changed_at - entry.changed_at))/60 AS ride_active_minutes,
    o.amount_cents
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
JOIN order_status_history entry ON entry.order_id = o.order_id AND entry.new_status = 5
JOIN order_status_history exit_h ON exit_h.order_id = o.order_id AND exit_h.old_status = 5
WHERE EXTRACT(EPOCH FROM (exit_h.changed_at - entry.changed_at)) > 180  -- > 3 minutes
ORDER BY ride_active_minutes DESC;
```

| order_id (short) | customer | ride_started | ride_ended | ride_active_minutes | amount_cents |
|-----------------|----------|-------------|-----------|-------------------|-------------|
| a82997e7... | Alex Petrov | 2026-02-27 15:42:18 | 2026-02-27 15:47:44 | 5.4 | 450 |
| 7d4e5f6a... | Kamran Hasanov | 2026-02-27 12:06:52 | 2026-02-27 12:18:33 | 11.7 | 350 |
| 9c2b3a4d... | Nigar Guliyev | 2026-02-27 09:33:15 | 2026-02-27 10:12:08 | 38.9 | 1250 |

> **Insight**: Most orders spend < 1 minute in status 5 before transitioning. The incident order (a82997e7) shows 5.4 minutes, but the real anomaly is the 16-minute gap *before* status 5 (between route_calculated and ride_active). The other long-duration entries are legitimate long rides.

## 6. Scooters Missing Geo Cluster

```sql
SELECT scooter_id, model, battery_level, lat, lng, status, geo_cluster
FROM scooters
WHERE geo_cluster IS NULL;
```

| scooter_id | model | battery_level | lat | lng | status | geo_cluster |
|------------|-------|--------------|-----|-----|--------|-------------|
| S.006068 | Jet City S2 | 85 | 40.4021 | 49.8592 | available | NULL |

> Only one scooter in the demo data has a NULL geo_cluster. In production, this could affect many scooters after a failed batch update or migration.

## 7. Full Order Lifecycle Join

```sql
-- Complete picture of order a82997e7 across all tables
SELECT
    c.name AS customer,
    c.email,
    s.scooter_id,
    s.model AS scooter_model,
    r.start_time,
    r.end_time,
    r.distance_km,
    o.amount_cents,
    o.status AS order_status,
    p.method AS payment_method,
    p.status AS payment_status,
    p.provider_ref
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
JOIN rides r ON o.ride_id = r.ride_id
JOIN scooters s ON r.scooter_id = s.scooter_id
LEFT JOIN payments p ON p.order_id = o.order_id
WHERE o.order_id = 'a82997e7-1b3d-4f8a-9c2e-6d5f4a3b2c1d';
```

| customer | email | scooter_id | scooter_model | start_time | end_time | distance_km | amount_cents | order_status | payment_method | payment_status | provider_ref |
|----------|-------|-----------|---------------|------------|----------|-------------|-------------|-------------|---------------|---------------|-------------|
| Alex Petrov | alex.petrov@gmail.com | scooter-7b3f2e91 | Jet Pro X1 | 2026-02-27 15:25:03 | 2026-02-27 15:47:44 | 2.30 | 450 | 14 | card | completed | stripe_pi_3Nk2Abc123def456 |

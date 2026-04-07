# Metrics

The LDAP Password Change service exposes a Prometheus-compatible `/metrics` endpoint to help you monitor the health,
performance, and usage of the application.

This endpoint combines standard HTTP/application metrics with custom metrics specifically designed to track LDAP
operation performance.

## Endpoint Details

* **Path:** `/metrics`
* **Method:** `GET`
* **Format:** Prometheus Text-based format (`text/plain`)

---

## Custom LDAP Metrics

To help monitor the performance and reliability of your connected LDAP server, the service exposes the following custom metric:

### `ldap_operation_duration_seconds` (Histogram)

Tracks the duration and count of specific LDAP operations executed by the service. Because it is a Histogram,
it automatically provides `_count` and `_sum` metrics, allowing you to calculate the average operation duration
and monitor total operation volume.

**Labels:**

| Label | Description | Example Values                              |
| :--- | :--- |:--------------------------------------------|
| `operation` | The specific LDAP action being performed. | `ping`, `bind`, `search`, `change_password` |
| `status` | The result of the operation. | `success`, `error`                          |

**Buckets:**
This metric uses the default Prometheus time buckets (`.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10` seconds)
to help you accurately track percentiles (e.g., P95, P99) for your LDAP query times.

#### Example Queries (PromQL)

**Average LDAP Search Duration (Last 5 minutes):**

```promql
rate(ldap_operation_duration_seconds_sum{operation="search", status="success"}[5m])
/
rate(ldap_operation_duration_seconds_count{operation="search", status="success"}[5m])

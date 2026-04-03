# Health Endpoints

The service exposes two dedicated health endpoints for use with container orchestrators and load balancers.

## Liveness — `GET /health/live`

Signals that the **process is running**. It has no external dependencies and always returns `200 OK`.

```shell
GET /health/live
```

```json
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

## Readiness — `GET /health/ready`

Signals that the service is **ready to serve traffic**. It verifies LDAP connectivity by performing a
service-account bind. Returns `503 Service Unavailable` if the LDAP server is unreachable.

```shell
GET /health/ready
```

**LDAP reachable:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

**LDAP unreachable:**

```json
HTTP/1.1 503 Service Unavailable
Content-Type: application/json

{"status":"unavailable","message":"ldap unreachable"}
```

## Log Behaviour

Both `/health/live` and `/health/ready` are **excluded from application logs** to prevent polling noise
in production environments.

## Kubernetes Example

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 3000
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 3000
  initialDelaySeconds: 3
  periodSeconds: 5
```

## Docker Compose Example

```yaml
services:
  ldap-password-change:
    image: ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/health/live"]
      interval: 30s
      timeout: 5s
      retries: 3
```

> **Tip:** Use `health/live` for the Docker Compose `healthcheck` (lightweight, no LDAP call).
> Use `health/ready` as the Kubernetes readiness probe so traffic is not routed until LDAP is reachable.

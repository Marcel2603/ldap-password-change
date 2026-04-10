# Version

The `/version` endpoint returns the current build version and commit hash of the running service.

## Endpoint Details

* **Path:** `/version`
* **Method:** `GET`
* **Format:** `application/json`

---

## Response

| Field | Type | Description |
| :--- | :--- | :--- |
| `version` | string | The application version (e.g. `v1.2.3`) |
| `commit` | string | The Git commit hash of the build |

### Example Response

```json
{
  "version": "v1.2.3",
  "commit": "a1b2c3d"
}
```

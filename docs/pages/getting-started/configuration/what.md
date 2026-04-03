# What to configure

All options are nested under their respective top-level YAML keys. Below is a full annotated reference.

## Server

Controls the HTTP listener.

```yaml
server:
  host: localhost   # Hostname used for CORS allowed origins
  port: 3000        # Port the service listens on
```

| Key            | Type   | Default     | Description                        |
|----------------|--------|-------------|------------------------------------|
| `server.host`  | string | `localhost` | Allowed CORS origin hostname       |
| `server.port`  | string | `3000`      | HTTP port to listen on             |

---

## LDAP

!!! warning
    Always set `ignoreTLS: false` in production and supply a valid `tlsCert` or use a trusted CA.

Configures the connection to your directory server.

```yaml
ldap:
  host: localhost:1389
  userDn: cn=admin,dc=example,dc=org
  password: password
  baseDn: ou=users,dc=example,dc=org
  searchFilter: "(objectClass=*)"
  ignoreTLS: true
  tlsCert: ""
```

| Key                   | Type    | Default                        | Description                                              |
|-----------------------|---------|--------------------------------|----------------------------------------------------------|
| `ldap.host`           | string  | `localhost:1389`               | LDAP host and port (`host:port`)                         |
| `ldap.userDn`         | string  | `cn=admin,dc=example,dc=org`   | Bind DN used to search for the user                      |
| `ldap.password`       | string  | `password`                     | Bind password for the search account                     |
| `ldap.baseDn`         | string  | `ou=users,dc=example,dc=org`   | Search base for user lookups                             |
| `ldap.searchFilter`   | string  | `(objectClass=*)`              | Additional LDAP filter applied during user search        |
| `ldap.ignoreTLS`      | bool    | `true`                         | Disable TLS (only for local dev — never in production!)  |
| `ldap.tlsCert`        | string  | `""`                           | Path to a custom CA certificate for TLS verification     |

---

## Validation

Controls the client- and server-side input validation rules.

```yaml
validation:
  username: ^[a-zA-Z0-9]+$
  password: ^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$
```

| Key                    | Type   | Description                                    |
|------------------------|--------|------------------------------------------------|
| `validation.username`  | string | Regex pattern the username must fully match    |
| `validation.password`  | string | Regex pattern the new password must fully match|

---

## Log

Controls the log level.

```yaml
log:
  level: info
```

| Key         | Type   | Default | Values                            |
|-------------|--------|---------|-----------------------------------|
| `log.level` | string | `info`  | `debug`, `info`, `warn`, `error`  |

All logs are emitted as **structured JSON** to stdout with source file and request ID (`req_id`) fields attached.

---

## UI

Controls visual customisation. All values are optional — defaults fall back to the bundled assets.

```yaml
ui:
  backgroundImage: ""
  customCss: ""
  favicon: ""
  icon: ""
```

| Key                    | Type   | Description                                                  |
|------------------------|--------|--------------------------------------------------------------|
| `ui.backgroundImage`   | string | URL or filename of a background image                        |
| `ui.customCss`         | string | URL or filename of an additional CSS stylesheet              |
| `ui.favicon`           | string | URL or filename to replace the browser favicon               |
| `ui.icon`              | string | URL or filename to replace the logo shown above the form     |

See [Customisation](../customisation.md) for more detail.

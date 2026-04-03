# LDAP Password Change

A lightweight, self-hosted web service that allows users to change their LDAP password
through a secure, modern browser interface — no admin intervention required.

[![Test & Lint](https://github.com/Marcel2603/ldap-password-change/actions/workflows/go-test.yml/badge.svg)](https://github.com/Marcel2603/ldap-password-change/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Marcel2603/ldap-password-change)](https://goreportcard.com/report/github.com/Marcel2603/ldap-password-change)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/Marcel2603/ldap-password-change)](https://github.com/Marcel2603/ldap-password-change/releases/latest)

## Features

- Self-service password change via Material Design V3 web UI
- Dark / Light / System theme switching
- Fully configurable via YAML or environment variables
- Custom branding: background image, favicon, logo, CSS
- Structured JSON logging with request IDs

## Docker

```bash
docker pull ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest

docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
```

## Configuration

Copy and edit `app.default.yml`:

```yaml
ldap:
  host: ldap.mycompany.com:636
  userDn: cn=svc-ldap,dc=mycompany,dc=com
  password: s3cr3t
  baseDn: ou=employees,dc=mycompany,dc=com
  ignoreTLS: false
```

> **Tip:** Avoid storing the bind password in `app.yml`. Use the `LDAP_PASSWORD` environment variable instead:
>
> ```bash
> docker run -e LDAP_PASSWORD=s3cr3t ...
> ```

Full reference → [docs/Configuration](https://marcel2603.github.io/ldap-password-change/getting-started/configuration/)

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) and the [docs](https://marcel2603.github.io/ldap-password-change/development/).

```bash
# Install pre-commit hooks
make init-precommit

# Run tests
make test

# Run linter
make lint
```

## Documentation

Full documentation is available at: <https://marcel2603.github.io/ldap-password-change/>

## Acknowledgements

Project default images are generated using Gemini 3 Pro.

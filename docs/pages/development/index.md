# Contributing

We welcome contributions of all kinds — bug fixes, new features, documentation improvements, and more!

## Before You Start

- Read the [Code of Conduct](https://www.contributor-covenant.org/).
- For **security issues**, follow the [Security Policy](../about/security.md) — do **not** open a public issue.
- For significant changes, open an [Issue](https://github.com/Marcel2603/ldap-password-change/issues) first
  to discuss your approach.

## Development Setup

### Prerequisites

- Go 1.24+
- Docker + Docker Compose (for local LDAP)
- `pre-commit` (optional but recommended)

### Clone & Run

```bash
git clone https://github.com/Marcel2603/ldap-password-change.git
cd ldap-password-change

# Start a local LDAP server
docker-compose -f local-dev/docker-compose.yaml up -d

# Install dependencies and generate templates
go mod tidy
make generate

# Run the service
make run
```

### Install Pre-Commit Hooks

```bash
make init-precommit
```

## Workflow

1. Fork the repository and create a feature branch from `main`.
2. Make your changes, following the coding standards below.
3. Write tests — we use table-driven tests and mock interfaces.
4. Run the test suite: `make test`
5. Run the linter: `make lint`
6. Submit a Pull Request targeting `main`.

## Coding Standards

This project follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

Key rules:

- **Initialisms** capitalised: `LDAP`, `UI`, `URL`, `ID`, `JSON`, `HTML`
- **No panics** — return errors, wrapped with `%w`
- **Pointer receivers** for LDAP clients and handlers
- **Functional Options** for complex constructors
- **Never log** passwords or LDAP bind credentials
- Always run `make generate` after editing `.templ` files

## Semantic Commits

We use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0-beta.4/).
PRs are merged via squash — only the PR title needs to follow the convention.

| Prefix     | Effect                   |
|------------|--------------------------|
| `feat:`    | Triggers a minor release |
| `fix:`     | Triggers a patch release |
| `chore:`   | No release               |
| `docs:`    | No release               |
| `refactor:`| No release               |
| `BREAKING CHANGE:` | Triggers a major release |

# Getting Started

## Installation

### Homebrew

```shell
brew tap Marcel2603/tap
brew install ldap-password-change
```

### From Release

Download the latest binary from the [Releases](https://github.com/Marcel2603/ldap-password-change/releases)
page and place it in your `$PATH`.

### Docker

The recommended way to run the service in production is via Docker:

```bash
docker pull ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
```

Run with a custom config file:

```bash
docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
```

Or with custom branding assets mounted:

```bash
docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  -v $PWD/custom:/app/custom \
  ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
```

### From Source

```bash
go install github.com/Marcel2603/ldap-password-change@latest
```

This installs `ldap-password-change` into your `$GOPATH/bin` or `$GOBIN`.

## First Run

### Prerequisites

You need a reachable LDAP server. For local development, a Docker Compose stack is provided:

```bash
docker-compose -f local-dev/docker-compose.yaml up -d
```

This starts an OpenLDAP instance pre-populated with test users at `localhost:1389`.

### Start the service

```bash
go mod tidy
make generate   # only needed once, or when framework versions change
make run
```

The service is now available at [http://localhost:3000](http://localhost:3000).

> **Note:** The default configuration in `app.default.yml` points to `localhost:1389` with `ignoreTLS: true`,
> which matches the local development LDAP container.

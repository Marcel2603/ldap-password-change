# Getting started

## Installation

### Homebrew

```shell
brew tap Marcel2603/tap
brew install ldap-password-change
```

### From release

Download the latest binary from the [Releases](https://github.com/Marcel2603/ldap-password-change/releases) page and place
it in your `$PATH`.

### From docker

```bash
docker pull ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest
docker run -v $PWD:/app ghcr.io/marcel2603/ldap-password-change/ldap-password-change:latest lint .
```

### From source

```bash
go install github.com/Marcel2603/ldap-password-change@latest
```

This will install `ldap-password-change` into your `$GOPATH/bin` or `$GOBIN`.

## First run

```bash
ldap-password-change
```

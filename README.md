# LDAP Password Change

## Execute

Assuming a reachable LDAP server at the URL in `app.default.yml` (see `local-dev` for a local setup):

```shell
go mod tidy

make generate  # you need to do this only once or when the framework versions change
make run

# open localhost:4000
```

## Docker

```shell
docker build . -t ldap-password-change
docker run -p 3333:3333 -e HOST=localhost ldap-password-change
```

## Contribute

Add the pre-commit hook:

```shell
cp .pre-commit-hook .git/hooks/pre-commit
```

See Makefile for development commands.

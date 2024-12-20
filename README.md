# LDAP Password Change

## Execute

```shell
go mod tidy
make run

# open localhost:3333
```

## Docker

```shell
docker build . -t ldap-password-change
docker run -p 3333:3333 -e HOST=localhost ldap-password-change
```
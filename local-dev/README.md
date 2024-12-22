# LocalDev

## Useful commands

**List users**
  
  ```bash
  ldapsearch -H "ldap://localhost:1389" -x -D "cn=admin,dc=example,dc=org" -w password -b "ou=users,dc=example,dc=org
  ```

**Play with ldap**

- start docker-compose `docker-compose -f local-dev/docker-compose.yaml up`
- run following commands
- list users: `go run cmd/ldap/handler.go query`
- change pw: `go run cmd/ldap/handler.go change cn=user02,ou=users,dc=example,dc=org password2 test1234`

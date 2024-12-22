# LocalDev

## Useful commands

**List users**
  
  ```bash
  ldapsearch -H "ldap://localhost:1389" -x -D "cn=admin,dc=example,dc=org" -w password -b "ou=users,dc=example,dc=org
  ```

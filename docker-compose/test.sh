#!/bin/bash

BATS_TEST_NAME_PREFIX='[LDAP] '

function setup_file() {
  # Setup local openldap service:
  docker run --rm -d --name "ldap" \
    --env LDAP_ADMIN_PASSWORD=admin \
    --env LDAP_ROOT='dc=example,dc=test' \
    --env LDAP_PORT_NUMBER=389 \
    --env LDAP_SKIP_DEFAULT_TREE=yes \
    --volume "./config/ldap/ldifs/:/ldifs/:ro" \
    --volume "./config/ldap/schemas/:/schemas/:ro" \
    bitnami/openldap:latest
}

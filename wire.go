//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/service/ldap"
	"ldap-password-change/internal/validation"
)

func InitializeService(ldapConfig config.LdapConfig, validationConfig config.ValidationConfig) (ldap.ServiceImpl, error) {
	wire.Build(ldap.CreateService, ldap.CreateWrapper)
	return ldap.ServiceImpl{}, nil
}

func InitializeValidator(validationConfig config.ValidationConfig) (validation.ValidatorImpl, error) {
	wire.Build(validation.CreateValidator)
	return validation.ValidatorImpl{}, nil
}

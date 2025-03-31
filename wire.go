//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"ldap-password-change/cmd/config"
	change_password "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/service/ldap"
	"ldap-password-change/internal/validation"
)

func InitChangePasswordHandler(ldapConfig config.LdapConfig, validationConfig config.ValidationConfig) (change_password.HandlerImpl, error) {
	wire.Build(change_password.CreateHandler, ldap.CreateService, ldap.CreateWrapper, validation.CreateValidator)
	return change_password.HandlerImpl{}, nil
}

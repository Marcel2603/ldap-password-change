package ldap

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"ldap-password-change/cmd/config"
	"log"
)

type Service interface {
	ChangePassword(username string, currentPassword string, newPassword string) error
}

type serviceImpl struct {
}

var (
	configuration = config.Configuration.Ldap
)

func CreateService() Service {
	testClient := createClient(configuration.UserDn, configuration.Password, configuration.Domain)
	defer testClient.Close()
	return serviceImpl{}
}

func (s serviceImpl) ChangePassword(username string, currentPassword string, newPassword string) error {
	client := createClient(configuration.UserDn, configuration.Password, configuration.Domain)
	defer client.Close()
	usernameDn := fmt.Sprintf("cn=%s,%s", username, configuration.BaseDn)
	passwdModifyRequest := ldap.NewPasswordModifyRequest(usernameDn, currentPassword, newPassword)
	if _, err := client.PasswordModify(passwdModifyRequest); err != nil {
		return err
	}
	fmt.Println("Password changed successfully")
	return nil
}

func createClient(username string, password string, domain string) *ldap.Conn {
	l, err := ldap.DialURL(domain, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}

	err = l.Bind(username, password)
	if err != nil {
		log.Println("Failed to bind ldap user")
		log.Fatal(err)
	}
	return l
}

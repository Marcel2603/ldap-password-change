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
	baseDn   string
	userDn   string
	password string
	domain   string
}

func CreateService(config config.LdapConfig) Service {
	testClient := createClient(config.UserDn, config.Password, config.Domain)
	defer testClient.Close()
	return &serviceImpl{
		baseDn:   config.BaseDn,
		userDn:   config.UserDn,
		password: config.Password,
		domain:   config.Domain,
	}
}

func (s *serviceImpl) ChangePassword(username string, currentPassword string, newPassword string) error {
	client := createClient(s.userDn, s.password, s.domain)
	defer client.Close()
	usernameDn := fmt.Sprintf("cn=%s,%s", username, s.baseDn)
	passwdModifyRequest := ldap.NewPasswordModifyRequest(usernameDn, currentPassword, newPassword)
	if _, err := client.PasswordModify(passwdModifyRequest); err != nil {
		return err
	}
	fmt.Println("Password changed successfully")
	return nil
}

func createClient(username string, password string, domain string) *ldap.Conn {
	conn, err := ldap.DialURL(domain, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Bind(username, password)
	if err != nil {
		log.Println("Failed to bind ldap user")
		log.Fatal(err)
	}
	return conn
}

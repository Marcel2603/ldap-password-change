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

type Conn interface {
	Bind(username, password string) error
	Close() error
	PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
}

type serviceImpl struct {
	baseDn      string
	userDn      string
	password    string
	domain      string
	ldapWrapper Wrapper
}

func CreateService(c config.LdapConfig, wrapper Wrapper) (Service, error) {
	testClient, err := createClient(wrapper, c.UserDn, c.Password, c.Domain)
	if err != nil {
		return nil, err
	}
	defer testClient.Close()
	return &serviceImpl{
		baseDn:      c.BaseDn,
		userDn:      c.UserDn,
		password:    c.Password,
		domain:      c.Domain,
		ldapWrapper: wrapper,
	}, nil
}

func (s *serviceImpl) ChangePassword(username string, currentPassword string, newPassword string) error {
	client, err := createClient(s.ldapWrapper, s.userDn, s.password, s.domain)
	if err != nil {
		return err
	}
	defer client.Close()
	usernameDn := fmt.Sprintf("cn=%s,%s", username, s.baseDn)
	passwdModifyRequest := ldap.NewPasswordModifyRequest(usernameDn, currentPassword, newPassword)
	if _, err := client.PasswordModify(passwdModifyRequest); err != nil {
		return err
	}
	log.Println("Password changed successfully")
	return nil
}

func createClient(wrapper Wrapper, username string, password string, domain string) (Conn, error) {
	conn, err := wrapper.DialURL(domain, wrapper.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}

	err = conn.Bind(username, password)
	if err != nil {
		log.Println("Failed to bind ldap user")
		return nil, err
	}
	return conn, nil
}

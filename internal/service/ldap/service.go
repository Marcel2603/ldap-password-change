package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"ldap-password-change/cmd/config"
	"log/slog"
	"os"

	"github.com/go-ldap/ldap/v3"
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
	host        string
	ignoreTLS   bool
	tlsCert     string
	logger      *slog.Logger
	ldapWrapper Wrapper
}

func CreateService(c config.LdapConfig, wrapper Wrapper, logger *slog.Logger) (Service, error) {
	testClient, err := createClient(wrapper, c.UserDn, c.Password, c.Host, c.IgnoreTLS, c.TlsCert, logger)
	if err != nil {
		return nil, err
	}
	defer testClient.Close()
	return &serviceImpl{
		baseDn:      c.BaseDn,
		userDn:      c.UserDn,
		password:    c.Password,
		host:        c.Host,
		ignoreTLS:   c.IgnoreTLS,
		tlsCert:     c.TlsCert,
		logger:      logger.With(slog.String("class", "service_ldap")),
		ldapWrapper: wrapper,
	}, nil
}

func (s *serviceImpl) ChangePassword(username string, currentPassword string, newPassword string) error {
	client, err := createClient(s.ldapWrapper, s.userDn, s.password, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return err
	}
	defer client.Close()
	usernameDn := fmt.Sprintf("cn=%s,%s", username, s.baseDn)
	passwdModifyRequest := ldap.NewPasswordModifyRequest(usernameDn, currentPassword, newPassword)
	if _, err := client.PasswordModify(passwdModifyRequest); err != nil {
		return err
	}
	s.logger.Info("Password changed successfully", slog.String("username", username))
	return nil
}

func createClient(wrapper Wrapper, username string, password string, host string, ignoreTLS bool, tlsCert string, logger *slog.Logger) (Conn, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: ignoreTLS,
	}

	if tlsCert != "" {
		cert, err := os.ReadFile(tlsCert)
		if err != nil {
			return nil, fmt.Errorf("failed to read tls cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(cert)
		tlsConfig.RootCAs = caCertPool
	}
	var connectionURL string
	if ignoreTLS {
		connectionURL = fmt.Sprintf("ldap://%s", host)
	} else {
		connectionURL = fmt.Sprintf("ldaps://%s", host)
	}

	conn, err := wrapper.DialURL(connectionURL, wrapper.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return nil, err
	}

	err = conn.Bind(username, password)
	if err != nil {
		logger.Error("Failed to bind ldap user", slog.String("error", err.Error()))
		return nil, err
	}
	return conn, nil
}

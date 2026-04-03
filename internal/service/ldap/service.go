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

var ErrUserNotFound = fmt.Errorf("user not found")

type Service interface {
	ChangePassword(username string, currentPassword string, newPassword string) error
	Ping() error
}

type Conn interface {
	Bind(username, password string) error
	Close() error
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
}

type serviceImpl struct {
	baseDn       string
	userDn       string
	password     string
	host         string
	ignoreTLS    bool
	tlsCert      string
	searchFilter string
	logger       *slog.Logger
	ldapWrapper  Wrapper
}

func CreateService(c config.LdapConfig, wrapper Wrapper, logger *slog.Logger) Service {
	return &serviceImpl{
		baseDn:       c.BaseDn,
		userDn:       c.UserDn,
		password:     c.Password,
		host:         c.Host,
		ignoreTLS:    c.IgnoreTLS,
		tlsCert:      c.TlsCert,
		searchFilter: c.SearchFilter,
		logger:       logger.With(slog.String("class", "service_ldap")),
		ldapWrapper:  wrapper,
	}
}

func (s *serviceImpl) Ping() error {
	client, err := createClient(s.ldapWrapper, s.userDn, s.password, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return err
	}
	return client.Close()
}

func (s *serviceImpl) ChangePassword(username string, currentPassword string, newPassword string) error {
	svcClient, err := createClient(s.ldapWrapper, s.userDn, s.password, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return err
	}
	defer svcClient.Close()

	filter := fmt.Sprintf("(&%s(cn=%s))", s.searchFilter, ldap.EscapeFilter(username))
	searchReq := ldap.NewSearchRequest(
		s.baseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)
	result, err := svcClient.Search(searchReq)
	if err != nil {
		return fmt.Errorf("ldap search failed: %w", err)
	}
	if len(result.Entries) == 0 {
		return ErrUserNotFound
	}
	userDN := result.Entries[0].DN

	userClient, err := createClient(s.ldapWrapper, userDN, currentPassword, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	defer userClient.Close()

	passwdModifyRequest := ldap.NewPasswordModifyRequest(userDN, currentPassword, newPassword)
	if _, err := userClient.PasswordModify(passwdModifyRequest); err != nil {
		return fmt.Errorf("password modify failed: %w", err)
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

	if err = conn.Bind(username, password); err != nil {
		logger.Error("Failed to bind ldap user", slog.String("error", err.Error()))
		return nil, err
	}
	return conn, nil
}

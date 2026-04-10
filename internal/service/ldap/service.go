package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Marcel2603/ldap-password-change/cmd/config"
	"github.com/Marcel2603/ldap-password-change/internal/types"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ldapOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ldap_operation_duration_seconds",
			Help:    "Duration of LDAP operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)
	ErrUserNotFound = fmt.Errorf("user not found")
)

type Service interface {
	SearchUser(username string) (*ldap.Entry, error)
	ChangePassword(userDN string, username string, currentPassword string, newPassword string) error
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
		tlsCert:      c.TLSCert,
		searchFilter: c.SearchFilter,
		logger:       logger.With(slog.String("class", "service_ldap")),
		ldapWrapper:  wrapper,
	}
}

func (s *serviceImpl) Ping() error {
	start := time.Now()
	client, err := createClient(s.ldapWrapper, s.userDn, s.password, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		ldapOperationDuration.WithLabelValues("ping", string(types.StatusError)).Observe(time.Since(start).Seconds())
		return err
	}
	ldapOperationDuration.WithLabelValues("ping", string(types.StatusSuccess)).Observe(time.Since(start).Seconds())
	return client.Close()
}

func (s *serviceImpl) SearchUser(username string) (*ldap.Entry, error) {
	start := time.Now()
	svcClient, err := createClient(s.ldapWrapper, s.userDn, s.password, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return nil, err
	}
	defer func(svcClient Conn) {
		err := svcClient.Close()
		if err != nil {
			s.logger.Error("Failed to close ldap service client", slog.String("error", err.Error()))
		}
	}(svcClient)

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
		ldapOperationDuration.WithLabelValues("search", string(types.StatusError)).Observe(time.Since(start).Seconds())
		return nil, fmt.Errorf("ldap search failed: %w", err)
	}
	if len(result.Entries) == 0 {
		ldapOperationDuration.WithLabelValues("search", string(types.StatusError)).Observe(time.Since(start).Seconds())
		return nil, ErrUserNotFound
	}
	ldapOperationDuration.WithLabelValues("search", string(types.StatusSuccess)).Observe(time.Since(start).Seconds())
	return result.Entries[0], nil
}

func (s *serviceImpl) ChangePassword(userDN string, username string, currentPassword string, newPassword string) error {
	start := time.Now()
	userClient, err := createClient(s.ldapWrapper, userDN, currentPassword, s.host, s.ignoreTLS, s.tlsCert, s.logger)
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	defer func(userClient Conn) {
		err := userClient.Close()
		if err != nil {
			s.logger.Error("Failed to close ldap user client", slog.String("error", err.Error()))
		}
	}(userClient)

	passwdModifyRequest := ldap.NewPasswordModifyRequest(userDN, currentPassword, newPassword)
	if _, err := userClient.PasswordModify(passwdModifyRequest); err != nil {
		ldapOperationDuration.WithLabelValues("change_password", string(types.StatusError)).Observe(time.Since(start).Seconds())
		return fmt.Errorf("password modify failed: %w", err)
	}

	ldapOperationDuration.WithLabelValues("change_password", string(types.StatusSuccess)).Observe(time.Since(start).Seconds())
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

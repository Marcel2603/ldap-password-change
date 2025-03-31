package ldap

import "crypto/tls"
import "github.com/go-ldap/ldap/v3"

type LdapWrapper interface {
	DialURL(addr string, opts ...ldap.DialOpt) (Conn, error)
	DialWithTLSConfig(tc *tls.Config) ldap.DialOpt
}

type wrapperImpl struct{}

func CreateWrapper() LdapWrapper {
	return &wrapperImpl{}
}

func (w *wrapperImpl) DialURL(addr string, opts ...ldap.DialOpt) (Conn, error) {
	return ldap.DialURL(addr, opts...)
}

func (w *wrapperImpl) DialWithTLSConfig(tc *tls.Config) ldap.DialOpt {
	return ldap.DialWithTLSConfig(tc)
}

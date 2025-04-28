package ldap

import "crypto/tls"
import "github.com/go-ldap/ldap/v3"

type Wrapper interface {
	DialURL(addr string, opts ...ldap.DialOpt) (Conn, error)
	DialWithTLSConfig(tc *tls.Config) ldap.DialOpt
}

type wrapperImpl struct{}

func CreateWrapper() Wrapper {
	return &wrapperImpl{}
}

func (*wrapperImpl) DialURL(addr string, opts ...ldap.DialOpt) (Conn, error) {
	return ldap.DialURL(addr, opts...)
}

func (*wrapperImpl) DialWithTLSConfig(tc *tls.Config) ldap.DialOpt {
	return ldap.DialWithTLSConfig(tc)
}

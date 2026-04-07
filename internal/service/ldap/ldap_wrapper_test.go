package ldap

import (
	"crypto/tls"
	"testing"
)

func TestLdapWrapper(t *testing.T) {
	wrapper := CreateWrapper()
	if wrapper == nil {
		t.Fatal("Expected wrapper")
	}

	opt := wrapper.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	if opt == nil {
		t.Fatal("Expected DialOpt")
	}

	_, err := wrapper.DialURL("ldap://invalid:123", opt)
	if err == nil {
		t.Fatal("Expected dial error")
	}
}

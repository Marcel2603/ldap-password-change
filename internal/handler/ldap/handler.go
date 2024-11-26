package ldap

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func Handler() {
	username := "admin"
	bindusername := "cn=admin,dc=example,dc=org"
	bindpassword := "admin"
	binddomain := "ldap://localhost:389"
	l := ldapbind(bindusername, bindpassword, binddomain)

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=org",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName={login})(!(userAccountControl:1.2.840.113556.1.4.803:=2)))(uid=%s))", ldap.EscapeFilter(username)), //TODO
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(sr.Entries)
	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}
}

func ldapbind(bindusername string, bindpassword string, binddomain string) *ldap.Conn {

	l, err := ldap.DialURL(binddomain) //TODO Add TLS Function
	if err != nil {
		log.Fatal(err) //TODO ADD Errors from here to display in UI
	}
	defer l.Close()

	// Reconnect with TLS
	// err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}
	return (l)
}

package ldap

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func Handler() {
	bindusername := "cn=admin,dc=example,dc=org"
	bindpassword := "admin"
	username := "admin"
	l, err := ldap.DialURL("ldap://localhost:389") //TODO Add TLS Function and add ldap var
	if err != nil {
		log.Fatal(err)
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

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=or",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}
	fmt.Println(sr)
}

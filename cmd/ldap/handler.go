package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
)

const (
	bindusername = "cn=admin,dc=example,dc=org"
	bindpassword = "password"
	binddomain   = "ldap://localhost:1389"
)

func queryUser() {

	// Search for the given username
	// Filters must start and finish with ()!
	searchRequest := ldap.NewSearchRequest(
		"ou=users,dc=example,dc=org", // The base DN to search
		ldap.ScopeWholeSubtree,       // Search the entire subtree
		ldap.NeverDerefAliases,       // Do not dereference aliases
		0,                            // No size limit
		0,                            // No time limit
		false,                        // Do not return types only
		"(objectClass=*)",            // The search filter
		[]string{"*"},                // The attributes to retrieve
		nil,                          // Controls
	)

	client := createClient(bindusername, bindpassword, binddomain)
	defer client.Close()
	result, err := client.Search(searchRequest)
	if err != nil {
		log.Fatalf("Failed to perform search: %v", err)
	}

	for _, entry := range result.Entries {
		entry.Print()
		fmt.Println("")
	}
}

func createClient(username string, password string, domain string) *ldap.Conn {
	l, err := ldap.DialURL(domain) //TODO Add TLS Function
	//l, err := ldap.DialURL(domain, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}

	err = l.Bind(username, password)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func main() {
	args := os.Args[1:]
	switch args[0] {
	case "query":
		queryUser()
	case "change":
		changePassword(args[1], args[2], args[3])
	}
}

func changePassword(username string, currentPassword string, newPassword string) {
	client := createClient(bindusername, bindpassword, binddomain)
	defer client.Close()
	passwdModifyRequest := ldap.NewPasswordModifyRequest(username, currentPassword, newPassword)
	if _, err := client.PasswordModify(passwdModifyRequest); err != nil {
		log.Fatalf("failed to modify password: %v", err)
	}
	fmt.Println("Password changed successfully")
}

// ldapreporter is a tool to list all of the members of an ldap group
// 1. connect to ldap server using bind account
// 2. query to list all groups
// 3. marshal results from query into json object
package main

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Version from build
var Version string

// logLevel sets the logger level from env LOG_LEVEL
var logLevel string = os.Getenv("LOG_LEVEL")

// setup logrus as log
func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	level := log.WarnLevel
	switch logLevel {
	case "INFO":
		level = log.InfoLevel
	case "DEBUG":
		level = log.DebugLevel
	}
	log.SetLevel(level)
}

// client for mocking ldap requests
type client interface {
	Search(s *ldap.SearchRequest) (*ldap.SearchResult, error)
}

type config struct {
	server   string // ldap server and port to connect
	bindUser string // user to bind with
	bindPass string // password to bind with
	isTLS    bool   // true if this will be a TLS session
}

type group struct {
	name    string
	members []user
}

type user string

// session creates a new ldap connection and returns a session object
func session(c *config) (*ldap.Conn, error) {
	conn, err := ldap.DialURL(c.server)
	if err != nil {
		return nil, err
	}

	err = conn.Bind(c.bindUser, c.bindPass)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// search ldap for an object and return it
// user needs to handle defer to close session
func get(c client, s *ldap.SearchRequest) (*ldap.SearchResult, error) {
	results, err := c.Search(s)
	if err != nil {
		return results, err
	}

	return results, nil
}

func main() {
	config := &config{
		server:   "ldap://localhost:8389",
		bindUser: "cn=admin,dc=planetexpress,dc=com",
		bindPass: "GoodNewsEveryone",
	}

	// create a new ldap session
	session, err := session(config)
	if err != nil {
		log.Error(err)
	}
	defer session.Close()

	log.Info(session)

	// search ldap
	search := ldap.NewSearchRequest(
		"dc=planetexpress,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectclass=Group))"),
		[]string{"cn", "member"},
		nil,
	)

	results, err := get(session, search)
	if err != nil {
		log.Error(err)
	}

	for _, result := range results.Entries {
		log.Info(result.DN)

		for _, member := range result.GetAttributeValues("member") {
			parsed, err := ldap.ParseDN(member)
			if err != nil {
				log.Error(err)
			}
			for _, rdn := range parsed.RDNs {
				for _, attrs := range rdn.Attributes {
					if attrs.Type == "cn" {
						log.Info(attrs.Value)
					}
				}
			}
		}
	}
}

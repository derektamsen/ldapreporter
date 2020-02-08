// ldapreporter is a tool to list all of the members of an ldap group
// 1. connect to ldap server using bind account
// 2. query to list all groups
// 3. marshal results from query into json object
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Version from build
var Version string

// cliFlags holds all of the cli options
type cliFlags struct {
	version                                              bool
	ldapServer, user, password, basedn, filter, loglevel string
}

var flags cliFlags = cliFlags{}

// setup logrus as log
func init() {
	flag.BoolVar(&flags.version, "version", false, "version of ldapreporter")
	flag.StringVar(
		&flags.ldapServer,
		"server",
		"",
		"LDAP Server. Ex: ldap://localhost:389",
	)
	flag.StringVar(
		&flags.user,
		"user",
		"",
		"LDAP Bind User. Ex: cn=admin,dc=planetexpress,dc=com",
	)
	flag.StringVar(
		&flags.password,
		"password",
		"",
		"LDAP Bind Password. Ex: GoodNewsEveryone",
	)
	flag.StringVar(
		&flags.basedn,
		"basedn",
		"dc=planetexpress,dc=com",
		"LDAP base DN. Ex: dc=planetexpress,dc=com",
	)
	flag.StringVar(
		&flags.filter,
		"searchfilter",
		"(&(objectclass=Group))",
		"LDAP search filter",
	)
	flag.StringVar(&flags.loglevel, "loglevel", "WARN", "Logging level")

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

// client for mocking ldap requests
type client interface {
	Search(s *ldap.SearchRequest) (*ldap.SearchResult, error)
}

// config holds the server url, and bind info
type config struct {
	server   string // ldap server and port to connect
	bindUser string // user to bind with
	bindPass string // password to bind with
}

type groupMembers map[string][]string

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

	log.Infof("created session: %v", conn)
	return conn, nil
}

// get searches ldap for an object and returns it. user needs to handle defer
// to close session
func get(c client, s *ldap.SearchRequest) (*ldap.SearchResult, error) {
	log.Infof("searching ldap with query: %v", s)
	results, err := c.Search(s)
	if err != nil {
		return results, err
	}

	log.Infof("found: %v", results)
	return results, nil
}

// getMembers returns all groups with their members
func getMembers(r *ldap.SearchResult) (groupMembers, error) {
	allGroups := groupMembers{}

	// loop through all the entries in a search result to create a map
	// keyed with CN of group
	for _, result := range r.Entries {
		group := result.GetAttributeValue("cn")

		log.Infof("finding members of %v", group)

		// loop through the member of a group to get a each user's CN
		for _, member := range result.GetAttributeValues("member") {
			parsed, err := ldap.ParseDN(member)
			if err != nil {
				return nil, err
			}

			// loop through each user and get their CN
			for _, rdn := range parsed.RDNs {
				for _, attrs := range rdn.Attributes {
					if attrs.Type == "cn" {
						allGroups[group] = append(allGroups[group], attrs.Value)
						log.Infof("found user %v in group %v", attrs.Value, group)
					}
				}
			}
		}
	}

	return allGroups, nil
}

func main() {
	// setup cli flags
	flag.Parse()

	// set logging level based on flags. Default to WARN
	var level log.Level
	switch flags.loglevel {
	case "INFO":
		level = log.InfoLevel
	case "DEBUG":
		level = log.DebugLevel
	case "WARN":
		level = log.WarnLevel
	default:
		level = log.WarnLevel
	}
	log.SetLevel(level)

	// print the version and exit(0)
	if flags.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	// create a new config with connection information
	config := &config{
		server:   flags.ldapServer,
		bindUser: flags.user,
		bindPass: flags.password,
	}

	// create a new ldap session
	session, err := session(config)
	if err != nil {
		log.Panicf("error creating a session: %v", err)
	}
	defer session.Close()

	// create the search object
	search := ldap.NewSearchRequest(
		flags.basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		flags.filter,
		[]string{"cn", "member"},
		nil,
	)

	// query ldap with search object
	results, err := get(session, search)
	if err != nil {
		log.Panicf("error searching ldap: %v", err)
	}

	// get a list of all the members of a group
	allGroups, err := getMembers(results)
	if err != nil {
		log.Panicf("error finding group members from search results: %v", err)
	}

	// convert the list of groups to json for output
	jAllGroups, err := json.Marshal(allGroups)
	if err != nil {
		log.Panicf("error converting go structu to json: %v", err)
	}
	fmt.Println(string(jAllGroups))
}

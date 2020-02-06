package main

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
)

// mockClient mocks the ldap client interface
type mockClient struct{}

// Search is a mock for the real ldap search interface
func (m *mockClient) Search(s *ldap.SearchRequest) (*ldap.SearchResult, error) {
	results := &ldap.SearchResult{}
	entry := &ldap.Entry{
		DN: "cn=ship_crew,ou=people,dc=planetexpress,dc=com",
	}

	results.Entries = append(results.Entries, entry)

	return results, nil
}

func TestGet(t *testing.T) {
	expect := "cn=ship_crew,ou=people,dc=planetexpress,dc=com"
	client := &mockClient{}

	results, err := get(client)
	if err != nil {
		t.Errorf("got get() error: %v", err)
	}

	result := results.Entries[0]
	got := result.DN

	if got != expect {
		t.Errorf("got %v expected %v", got, expect)
	}
}

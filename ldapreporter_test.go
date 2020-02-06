package main

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

// mockClient mocks the ldap client interface
type mockClient struct {
	searchDN string
}

// Search is a mock for the real ldap search interface
func (m *mockClient) Search(s *ldap.SearchRequest) (*ldap.SearchResult, error) {
	results := &ldap.SearchResult{}
	entry := &ldap.Entry{DN: m.searchDN}

	results.Entries = append(results.Entries, entry)

	return results, nil
}

// test get function
func TestGet(t *testing.T) {
	cases := []struct {
		stub, expect string
		err          error
	}{
		{
			stub:   "cn=ship_crew,ou=people,dc=planetexpress,dc=com",
			expect: "cn=ship_crew,ou=people,dc=planetexpress,dc=com",
			err:    nil,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s in %s", tc.stub, tc.expect), func(t *testing.T) {
			client := &mockClient{
				searchDN: tc.stub,
			}

			results, err := get(client)
			if err != tc.err {
				t.Errorf("got get() error: %v", err)
			}

			result := results.Entries[0]
			got := result.DN

			if got != tc.expect {
				t.Errorf("got %v expected %v", got, tc.expect)
			}
		})
	}
}

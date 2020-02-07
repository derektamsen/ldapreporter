package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

// mockClient mocks the ldap client interface
type mockClient struct {
	searchDN string
	err      error
}

var mockErr error = errors.New("Did Everything Just Taste Purple For A Second?")

// Search is a mock for the real ldap search interface
func (m *mockClient) Search(s *ldap.SearchRequest) (*ldap.SearchResult, error) {
	results := &ldap.SearchResult{}
	entry := &ldap.Entry{DN: m.searchDN}

	results.Entries = append(results.Entries, entry)

	if m.err != nil {
		return results, m.err
	}

	return results, nil
}

// test get function
func TestGet(t *testing.T) {
	cases := []struct {
		stub, expect   string
		expectErr, err error
	}{
		{
			stub:      "cn=ship_crew,ou=people,dc=planetexpress,dc=com",
			expect:    "cn=ship_crew,ou=people,dc=planetexpress,dc=com",
			err:       nil,
			expectErr: nil,
		},
		{
			stub:      "err",
			expect:    "err",
			err:       mockErr,
			expectErr: mockErr,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s in %s", tc.stub, tc.expect), func(t *testing.T) {
			client := &mockClient{
				searchDN: tc.stub,
				err:      tc.err,
			}

			// create a new search object
			search := ldap.NewSearchRequest(
				"dc=planetexpress,dc=com",
				ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
				fmt.Sprintf("(&(objectclass=Group))"),
				[]string{"cn", "member"},
				nil,
			)

			// test get()
			results, err := get(client, search)
			result := results.Entries[0]
			got := result.DN

			// check errors and results to see if they match what we expect
			if err != tc.expectErr {
				t.Errorf("got %s expected %s", err, tc.expectErr)
			} else if got != tc.expect {
				t.Errorf("got %s expected %s", got, tc.expect)
			}
		})
	}
}

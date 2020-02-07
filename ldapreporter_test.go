package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

// mockClient mocks the ldap client interface
type mockClient struct {
	searchDN string
	err      error
}

var mockErr error = errors.New("Did Everything Just Taste Purple For A Second?")

var mockUser ldap.AttributeTypeAndValue = ldap.AttributeTypeAndValue{
	Type:  "cn",
	Value: "Hubert J. Farnsworth",
}

var mockUserName string = mockUser.Value

var strMockUser string = fmt.Sprintf("%s=%s", mockUser.Type, mockUser.Value)

var mockGroupOneMember ldap.Entry = ldap.Entry{
	DN: "cn=admin,dc=planetexpress,dc=com",
	Attributes: []*ldap.EntryAttribute{
		{
			Name:   "member",
			Values: []string{strMockUser},
		},
	},
}

var mockGroupName string = mockGroupOneMember.GetAttributeValue("cn")

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

// test getMembers
func TestGetMembers(t *testing.T) {
	cases := []struct {
		name  string
		input ldap.SearchResult
		want  groupMembers
	}{
		{
			name: "group with one member",
			input: ldap.SearchResult{
				Entries: []*ldap.Entry{
					&ldap.Entry{DN: mockGroupOneMember.DN, Attributes: mockGroupOneMember.Attributes},
				},
			},
			want: groupMembers{mockGroupName: {mockUserName}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getMembers(&tc.input)
			if err != nil {
				t.Errorf("got %v; want%v", got, nil)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}

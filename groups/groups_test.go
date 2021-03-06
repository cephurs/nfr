package groups

import (
	"net"
	"testing"
)

func TestGroups(t *testing.T) {
	var testsGroups = []struct {
		name   string
		groups []*Group

		domains  []string
		ips      []net.IP
		expected []bool
	}{
		{
			"allow any ips",
			[]*Group{
				{
					Name:     "allow any",
					Includes: []string{"0.0.0.0/0"},
				},
			},
			[]string{"a"},
			[]net.IP{net.IPv4(10, 0, 0, 0)},
			[]bool{true},
		},
		{
			"allow private networks",
			[]*Group{
				{
					Name:     "private network 1",
					Includes: []string{"10.0.0.0/8"},
				},
				{
					Name:     "private network 2",
					Includes: []string{"192.168.0.0/16"},
				},
			},
			[]string{"a", "a", "a"},
			[]net.IP{net.IPv4(10, 0, 0, 0), net.IPv4(192, 168, 0, 0), net.IPv4(11, 0, 0, 0)},
			[]bool{true, true, false},
		},
		{
			"allow private networks with excludes",
			[]*Group{
				{
					Name:     "private network 1",
					Includes: []string{"10.0.0.0/8"},
					Excludes: []string{"10.1.0.0/16"},
				},
				{
					Name:     "private network 2",
					Includes: []string{"192.168.0.0/16"},
					Excludes: []string{"10.2.0.0/16"},
				},
			},
			[]string{"a", "a", "a", "a"},
			[]net.IP{net.IPv4(10, 0, 0, 0), net.IPv4(192, 168, 0, 0), net.IPv4(10, 1, 0, 0), net.IPv4(10, 2, 0, 0)},
			[]bool{true, true, false, true},
		},
		{
			"include in one group then exclude in next group",
			[]*Group{
				{
					Name:     "private network 1",
					Includes: []string{"10.0.0.0/8"},
				},
				{
					Name:     "private network 2",
					Includes: []string{"10.1.0.0/16"},
					Excludes: []string{"10.1.1.0/24"},
				},
			},
			[]string{"a", "a"},
			[]net.IP{net.IPv4(10, 1, 0, 0), net.IPv4(10, 1, 1, 0)},
			[]bool{true, false},
		},
		{
			"exclude domain in multiple groups",
			[]*Group{
				{
					Name:     "private network 1",
					Includes: []string{"10.0.0.0/16"},
					Domains:  []string{"a"},
				},
				{
					Name:     "private network 2",
					Includes: []string{"10.1.0.0/16"},
					Domains:  []string{"b"},
				},
			},
			[]string{"a", "b", "a", "b"},
			[]net.IP{net.IPv4(10, 0, 0, 0), net.IPv4(10, 0, 0, 0), net.IPv4(10, 1, 0, 0), net.IPv4(10, 1, 0, 0)},
			[]bool{false, true, true, false},
		},
	}

	for _, tt := range testsGroups {
		g := New()
		for _, group := range tt.groups {
			if err := g.Add(group); err != nil {
				t.Fatal(err)
			}
		}
		for i := range tt.domains {
			if _, b := g.IsDNSQueryWhitelisted(tt.domains[i], tt.ips[i]); b != tt.expected[i] {
				t.Fatalf("IsDNSQueryWhitelisted(%s, %s) got %t; expected %t",
					tt.domains[i], tt.ips[i], !tt.expected[i], tt.expected[i])
			}
		}
	}
}

func TestEmptyGroup(t *testing.T) {
	var g *Groups
	if _, b := g.IsDNSQueryWhitelisted("a", net.IPv4(10, 0, 0, 0)); !b {
		t.Fatalf("nil groups must whitelist domain")
	}
	g = New()
	if _, b := g.IsDNSQueryWhitelisted("a", net.IPv4(10, 0, 0, 0)); !b {
		t.Fatalf("no groups must whitelist domain")
	}
}

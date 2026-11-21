package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/wkhere/shargs"
)

type r interface{ testResult() }
type cp func(*config) bool // cp = config pattern/predicate
type ep struct{ string }   // ep = error pattern

func (cp) testResult() {}
func (ep) testResult() {}

func (p ep) match(s string) bool { return strings.Contains(s, p.string) }

func e(s string) ep { return ep{s} }

func empty() cp          { return func(*config) bool { return true } }
func hasHelp(b bool) cp  { return func(c *config) bool { return (c.help != nil) == b } }
func hasSep(s string) cp { return func(c *config) bool { return string(c.proc.OutSep) == s } }
func hasSpec(xx []int) cp {
	return func(c *config) bool { return slices.Equal(c.proc.GroupSpec, xx) }
}

func and(ff ...cp) cp {
	return func(c *config) bool {
		for _, f := range ff {
			if !f(c) {
				return false
			}
		}
		return true
	}
}

func TestParseArgs(t *testing.T) {

	type tc struct {
		input string
		want  r
	}
	type ii = []int

	var tab = []struct {
		group string
		cases []tc
	}{
		{"trivial", []tc{
			{"2", and(hasHelp(false), hasSpec(ii{2}))},
			{"-h", hasHelp(true)},
			{"--help", hasHelp(true)},
			{"foo -h", hasHelp(true)},
			{"foo --help", hasHelp(true)},
		}},

		{"flags-given", []tc{
			{"-s, 2", and(hasSep(","), hasSpec(ii{2}))},
			{"2 -s,", and(hasSep(","), hasSpec(ii{2}))},
			{"-s , 2", hasSep(",")},
			{"-s=, 2", hasSep(",")},
			{"--separator , 2", hasSep(",")},
			{"--separator=, 2", hasSep(",")},
		}},

		{"spec", []tc{
			{"2", hasSpec(ii{2})},
			{"2 4", hasSpec(ii{2, 4})},
			{"2 1 3", hasSpec(ii{2, 1, 3})},
		}},

		{"defaults", []tc{
			{"2", hasSep(" ")},
		}},

		{"invalid", []tc{
			{"foo", e("error reading group spec")},
			{"-x", e("unknown flag")},
			{"-s", e("flag -s needs a value")},
			{"--separator", e("flag --separator needs a value")},
			{"-s foo", e("output separator needs to be 1 character")},
			{"--separator foo", e("output separator needs to be 1 character")},
			{"2 0", e("group spec needs to be positive")},
			{"0 2", e("group spec needs to be positive")},
			{"2 -1", e("unknown flag -1")},
		}},
	}

	for _, g := range tab {
		for i, tc := range g.cases {
			t.Run(fmt.Sprintf("%s-%d", g.group, i), func(t *testing.T) {

				args, err := shargs.Fields(tc.input)
				if err != nil {
					t.Fatalf("shell args: %v", err)
				}
				c, err := parseArgs(args)

				switch errp, wantErr := tc.want.(ep); {
				case err != nil && !wantErr:
					t.Errorf("unexpected error: %s", err)

				case err != nil && wantErr:
					if !errp.match(err.Error()) {
						t.Errorf("expected error with %q; have: %q", errp.string, err)
					}

				case err == nil && wantErr:
					t.Errorf("expected error with %q", errp.string)

				default:
					if f := tc.want.(cp); !f(&c) {
						t.Errorf("predicate failed on %+v", c)
					}
				}
			})
		}
	}
}

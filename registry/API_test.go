package registry

import (
	"reflect"
	"testing"
)

func TestExtractKeyValuePairs(t *testing.T) {
	cases := []struct {
		in  string
		out map[string]string
	}{
		{``, map[string]string{}},
		{`foo="bar"`, map[string]string{`foo`: `bar`}},
		{`foo="bar",foz="baz"`, map[string]string{`foo`: `bar`, `foz`: `baz`}},
		{`foo="bar", foz="baz"`, map[string]string{`foo`: `bar`, `foz`: `baz`}},
	}

	for _, c := range cases {
		out := extractKeyValuePairs(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Fatalf("got %s expected %s", out, c.out)
		}
	}
}

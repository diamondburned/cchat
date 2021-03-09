package repository

import "testing"

func TestTypeQual(t *testing.T) {
	type test struct {
		typePath string
		path     string
		typ      string
	}

	var tests = []test{
		{"string", "", "string"},
		{"context.Context", "context", "Context"},
		{
			"github.com/diamondburned/cchat/text.Rich",
			"github.com/diamondburned/cchat/text", "Rich",
		},
		{
			"(github.com/diamondburned/cchat/text).Rich",
			"github.com/diamondburned/cchat/text", "Rich",
		},
	}

	for _, test := range tests {
		path, typ := TypeQual(test.typePath)
		if path != test.path {
			t.Errorf("Unexpected path %q != %q", path, test.path)
		}
		if typ != test.typ {
			t.Errorf("Unexpected type %q != %q", typ, test.typ)
		}
	}
}

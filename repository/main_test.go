package repository

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/go-test/deep"
)

func TestGob(t *testing.T) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(Main); err != nil {
		t.Fatal("Failed to gob encode:", err)
	}

	t.Log("Marshaled; total bytes:", buf.Len())

	var unmarshaled Packages

	if err := gob.NewDecoder(&buf).Decode(&unmarshaled); err != nil {
		t.Fatal("Failed to gob decode:", err)
	}

	if eq := deep.Equal(Main, unmarshaled); eq != nil {
		t.Fatal("Inequalities after unmarshaling:", eq)
	}
}

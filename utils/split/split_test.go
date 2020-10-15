package split

import "testing"

type testEntry struct {
	input  string
	offset int64
	output []string
	index  int64
}

func (test testEntry) compare(t *testing.T, words []string, index int64) {
	if !strsleq(words, test.output) {
		t.Error("Mismatch output (input/got/expected)", test.input, words, test.output)
	}
	if index != test.index {
		t.Error("Mismatch index (input/got/expected)", test.input, index, test.index)
	}
}

// string slice equal
func strsleq(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

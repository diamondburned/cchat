package split

import "testing"

func TestSpaceIndexed(t *testing.T) {
	var tests = []struct {
		input  string
		offset int64
		output []string
		index  int64
	}{{
		input:  "bruhemus momentus lorem ipsum",
		offset: 13, //       ^
		output: []string{"bruhemus", "momentus", "lorem", "ipsum"},
		index:  1,
	}, {
		input: "Yoohoo! My name's Astolfo! I belong to the Rider-class! And, and... uhm, nice " +
			"to meet you!",
		offset: 37, //                               ^
		output: []string{
			"Yoohoo!", "My", "name's", "Astolfo!", "I", "belong", "to", "the", "Rider-class!",
			"And,", "and...", "uhm,", "nice", "to", "meet", "you!"},
		index: 6,
	}, {
		input:  "sorry, what were you typing?",
		offset: int64(len("sorry, what were you typing?")) - 1,
		output: []string{"sorry,", "what", "were", "you", "typing?"},
		index:  4,
	}, {
		input:  "zeroed out input",
		offset: 0,
		output: []string{"zeroed", "out", "input"},
		index:  0,
	}, {
		input:  "に　ほ　ん　ご",
		offset: 3,
		output: []string{"に", "ほ", "ん", "ご"},
		index:  1,
	}}

	for _, test := range tests {
		a, j := SpaceIndexed(test.input, test.offset)
		if !strsleq(a, test.output) {
			t.Error("Mismatch output (input/got/expected)", test.input, a, test.output)
		}
		if j != test.index {
			t.Error("Mismatch index (input/got/expected)", test.input, j, test.index)
		}
	}
}

const benchstr = "Alright, Master! I'm your blade, your edge and your arrow! You've placed " +
	"so much trust in me, despite how weak I am - I'll do everything in my power to not " +
			"disappoint you!"
const benchcursor = 32 // arbitrary

func BenchmarkSpaceIndexed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SpaceIndexed(benchstr, benchcursor)
	}
}

func BenchmarkSpaceIndexedLong(b *testing.B) {
	const benchstr = benchstr + benchstr + benchstr + benchstr + benchstr + benchstr

	for i := 0; i < b.N; i++ {
		SpaceIndexed(benchstr, benchcursor)
	}
}

// same as benchstr but w/ a horizontal line (outside ascii)
const benchstr8 = "Alright, Master! I'm your blade, your edge and your arrow! You've placed " +
	"so much trust in me, despite how weak I am ― I'll do everything in my power to not " +
	"disappoint you!"

func BenchmarkSpaceIndexedUTF8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SpaceIndexed(benchstr8, benchcursor)
	}
}

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

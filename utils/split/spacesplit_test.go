package split

import "testing"

var spaceSplitTests = []testEntry{
	{
		input:  "bruhemus momentus lorem ipsum",
		offset: 13, //       ^
		output: []string{"bruhemus", "momentus", "lorem", "ipsum"},
		index:  1,
	},
	{
		input: "Yoohoo! My name's Astolfo! I belong to the Rider-class! And, and... uhm, nice " +
			"to meet you!",
		offset: 37, //                               ^
		output: []string{
			"Yoohoo!", "My", "name's", "Astolfo!", "I", "belong", "to", "the", "Rider-class!",
			"And,", "and...", "uhm,", "nice", "to", "meet", "you!"},
		index: 6,
	},
	{
		input:  "sorry, what were you typing?",
		offset: int64(len("sorry, what were you typing?")) - 1,
		output: []string{"sorry,", "what", "were", "you", "typing?"},
		index:  4,
	},
	{
		input:  "zeroed out input",
		offset: 0,
		output: []string{"zeroed", "out", "input"},
		index:  0,
	},
	{
		input:  "に　ほ　ん　ご",
		offset: 3,
		output: []string{"に", "ほ", "ん", "ご"},
		index:  1,
	},
}

func TestSpaceIndexed(t *testing.T) {
	for _, test := range spaceSplitTests {
		a, j := SpaceIndexed(test.input, test.offset)
		test.compare(t, a, j)
	}
}

const benchstr = "Alright, Master! I\\'m your blade, your edge and your arrow! You\\'ve placed " +
	"so much trust in me, despite how weak I am - I\\'ll do everything in my power to not " +
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

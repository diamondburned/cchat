package split

import "testing"

var argsSplitTests = []testEntry{
	{
		input:  "bruhemus 'momentus lorem' \"ipsum\"",
		offset: 13, //       ^
		output: []string{"bruhemus", "momentus lorem", "ipsum"},
		index:  1,
	},
	{
		input: "Yoohoo! My name\\'s Astolfo! I belong to the Rider-class! And, and... uhm, nice " +
			"to meet you!",
		offset: 37, //                                 ^
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
		output: []string{"に　ほ　ん　ご"},
		index:  0,
	},
	{
		input:  `echo "this \"quote\" is a regular test"`,
		offset: 5,
		output: []string{"echo", `this "quote" is a regular test`},
		index:  1,
	},
	{
		input:  `echo "this \"quote\" is a regular test"`,
		offset: 4,
		output: []string{"echo", `this "quote" is a regular test`},
		index:  0,
	},
}

func TestArgsIndexed(t *testing.T) {
	for _, test := range argsSplitTests {
		a, j := ArgsIndexed(test.input, test.offset)
		test.compare(t, a, j)

		if expect, got := approxArgLen(test.input), len(test.output); expect != got {
			t.Error("Approximated arg len is off, (expected/got)", expect, got)
		}
	}
}

const argsbenchstr = "Alright, Master! I\\'m your blade, your edge and your arrow! You\\'ve " +
	"placed so much trust in me, despite how weak I am - I\\'ll do everything in my power to not " +
	"disappoint you!"

func BenchmarkArgsIndexed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ArgsIndexed(benchstr, benchcursor)
	}
}

func BenchmarkArgsIndexedLong(b *testing.B) {
	const benchstr = benchstr + benchstr + benchstr + benchstr + benchstr + benchstr

	for i := 0; i < b.N; i++ {
		ArgsIndexed(benchstr, benchcursor)
	}
}

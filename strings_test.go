package godok_test

import (
	"testing"

	"github.com/ronaudinho/godok"
)

var ntcTests = []struct {
	name string
	in   string
	out  string
}{
	{
		name: "test",
		in:   "test",
		out:  "tests",
	},
	{
		name: "testA",
		in:   "testA",
		out:  "tests a",
	},
	{
		name: "testAA",
		in:   "testAA",
		out:  "tests aa",
	},
	{
		name: "testAa",
		in:   "testAa",
		out:  "tests aa",
	},
	{
		name: "testAAa",
		in:   "testAAa",
		out:  "tests a aa",
	},
}

func TestNameToComment(t *testing.T) {
	for _, tt := range ntcTests {
		t.Run(tt.name, func(t *testing.T) {
			out := godok.ToComment(tt.in)
			if out != tt.out {
				t.Errorf("got %s, wanted %s", out, tt.out)
			}
		})
	}
}

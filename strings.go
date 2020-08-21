package godok

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToComment(src string) string {
	spl := split(src)
	spl[0] = strings.ToLower(spl[0]) + "s"
	if len(spl) > 1 {
		for i := range spl {
			if i == 0 {
				continue
			}
			spl[i] = strings.ToLower(spl[i])
		}
	}
	return strings.Join(spl, " ")
}

// split splits CamelCase into slice of string
// copied with modification from https://github.com/fatih/camelcase
// not sure this is the correct way to attribute the source
// why not import?
// A little copying is better than a little dependency - Go Proverbs
func split(src string) []string {
	var spl []string
	if !utf8.ValidString(src) {
		return []string{src}
	}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	for _, s := range runes {
		if len(s) > 0 {
			spl = append(spl, string(s))
		}
	}
	return spl
}

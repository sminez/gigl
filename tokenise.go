package gigl

import (
	"fmt"
	"regexp"
	"strings"
)

type tag struct {
	name  string
	regex string
}

func buildGroups(tags []tag) []string {
	groups := make([]string, 0)
	for _, t := range tags {
		groups = append(groups, fmt.Sprintf("(?P<%v>%v)", t.name, t.regex))
	}
	return groups
}

var tags = []tag{
	tag{"PAREN", `[(){}\[\]]`},
	tag{"COMPLEX", `-?\d+\.?\d*[+-]\d+\.?\d*j`},
	tag{"COMPLEX_PURE", `-?\d+\.?\d*j`},
	tag{"FLOAT", `-?\d+\.\d+`},
	tag{"INT", `-?\d+`},
	tag{"QUOTE", "'"},
	tag{"QUASI_QUOTE", "`"},
	tag{"UNQUOTE", ","},
	tag{"UNQUOTE_SPLICE", ",@"},
	tag{"NEWLINE", "\n"},
	tag{"WHITESPACE", `\s+`},
	tag{"STRING", `"([^"]*)"`},
	tag{"SYMBOL", "."},
}

var lexTags = strings.Join(buildGroups(tags), "|")

// embed regexp.Regexp in a new type so we can extend it
type tokeniser struct {
	*regexp.Regexp
}

type token struct {
	tag  string
	text string
}

// Split an input string into tokens for parsing
func (t *tokeniser) tokenise(s string) []token {
	tokens := make([]token, 0)

	matches := t.FindAllStringSubmatchIndex(s, -1)
	if matches == nil {
		return tokens
	}
	fmt.Println(matches)
	groupNames := t.SubexpNames()

	for i, indices := range matches {
		// Ignore the whole regexp match
		if i == 0 {
			continue
		}

		tokens = append(tokens, token{groupNames[i], s[indices[0]:indices[1]]})

	}
	return tokens
}

// an example regular expression
// var myExp = tokeniser{regexp.MustCompile(`(?P<PAREN>[(){}\[\]])`)}
// fmt.Printf("%+v", myExp.FindStringSubmatchMap("1234.5678.9"))

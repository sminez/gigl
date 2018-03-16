package gigl

import (
	"fmt"
	"regexp"
	"strconv"
)

// Regex objects for constructing atoms
// (Ensure that keywords and symbols contain valid characters)
var quotes = map[string]SYMBOL{
	"'": SYMBOL("quote"), "`": SYMBOL("quasiquote"),
	"~": SYMBOL("unquote"), "~@": SYMBOL("unquote-splicing"),
}

type tag struct {
	name  string
	regex *regexp.Regexp
}

type token struct {
	Tag  string
	Text string
}

// Tokeniser turns a string into a slice of tokens for parsing
type Tokeniser struct {
	tags   []tag
	input  string
	ix     int
	tokens []token
}

// NewTokeniser constructs a new Tokeniser...!
func NewTokeniser() *Tokeniser {
	return &Tokeniser{
		tags: []tag{
			tag{"LINECOMMENT", regexp.MustCompile("^;;")},
			tag{"LIST_START", regexp.MustCompile(`^\(`)},
			tag{"LIST_END", regexp.MustCompile(`^\)`)},
			tag{"VEC_START", regexp.MustCompile(`^\[`)},
			tag{"VEC_END", regexp.MustCompile(`^\]`)},
			tag{"MAP_START", regexp.MustCompile("^{")},
			tag{"SET_START", regexp.MustCompile("^#{")},
			tag{"MAP_OR_SET_END", regexp.MustCompile("^}")},
			tag{"COMPLEX", regexp.MustCompile(`^-?\d+\.?\d*[+-]\d+\.?\d*j`)},
			tag{"COMPLEX_PURE", regexp.MustCompile(`^-?\d+\.?\d*j`)},
			tag{"FLOAT", regexp.MustCompile(`^-?\d+\.\d+`)},
			tag{"INT", regexp.MustCompile(`^-?\d+`)},
			tag{"BOOL", regexp.MustCompile(`^#[tf]`)},
			tag{"SPLICE", regexp.MustCompile("^~@")},
			tag{"QUOTE", regexp.MustCompile("^['`~]")},
			tag{"NEWLINE", regexp.MustCompile(`^\n`)},
			tag{"WHITESPACE", regexp.MustCompile(`^\s+`)},
			tag{"COMMA", regexp.MustCompile(`^,`)},
			tag{"STRING", regexp.MustCompile(`^"([^"]*)"`)},
			tag{"KEYWORD", regexp.MustCompile("^:[^(){}\\[\\],'`@:; \t\n]*")},
			tag{"SYMBOL", regexp.MustCompile("^[^(){}\\[\\],'`@:; \t\n]*")},
			tag{"ERROR", regexp.MustCompile(".*")},
		},
	}
}

// Tokenise splits an input string into tokens for parsing
func (t *Tokeniser) Tokenise(s string) {
	t.tokens = make([]token, 0)
	t.ix = 0
	t.input = s
	inComment := false

	for len(s) > 0 {
		for _, tag := range t.tags {
			if loc := tag.regex.FindStringIndex(s); loc != nil {
				switch tag.name {
				case "LINECOMMENT":
					inComment = true

				case "WHITESPACE", "COMMA":
					// A la Clojure/edn, commas are also whitespace

				case "NEWLINE":
					// Terminate a line comment if we are in one
					inComment = false

				default:
					if !inComment {
						t.tokens = append(t.tokens, token{tag.name, s[loc[0]:loc[1]]})
					}
				}

				s = s[loc[1]:]
				break
			}
		}
	}
}

// NextToken return the next token in the stream
func (t *Tokeniser) NextToken() (token, error) {
	if t.ix >= len(t.tokens) {
		return token{}, fmt.Errorf("Ran out of tokens")
	}
	tok := t.tokens[t.ix]
	t.ix++
	return tok, nil
}

// Tokenise an input string and then parse the result
func (t *Tokeniser) read(s string) (lispVal, error) {
	t.Tokenise(s)
	return t.parseTokens()
}

// Convert tokens into internal data structures
func (t *Tokeniser) parseTokens() (lispVal, error) {
	// Pull off the first token
	token, err := t.NextToken()
	if err != nil {
		return nil, fmt.Errorf("Syntax error: %v", t.input)
	}

	switch token.Tag {
	case "LIST_START":
		// Start of a list so recuse and build it up
		lst := make([]lispVal, 0)
		parsedToken, err := t.parseTokens()
		for {
			if parsedToken == "LIST_END" {
				if err != nil {
					return nil, fmt.Errorf("Syntax error")
				}
				return List(lst...), nil
			}
			lst = append(lst, parsedToken)
			parsedToken, err = t.parseTokens()
			if err != nil {
				return nil, err
			}
		}

	case "LIST_END":
		return "LIST_END", nil

	case "VEC_START":
		lst := make([]lispVal, 0)
		parsedToken, err := t.parseTokens()
		for {
			if parsedToken == "VEC_END" {
				if err != nil {
					return nil, fmt.Errorf("Syntax error")
				}
				return lst, nil
			}
			lst = append(lst, parsedToken)
			parsedToken, err = t.parseTokens()
			if err != nil {
				return nil, err
			}
		}

	case "VEC_END":
		return "VEC_END", nil

	case "QUOTE", "SPLICE":
		// Something is being quoted or unquoted
		quotedList := make([]lispVal, 0)
		quotedList = append(quotedList, quotes[token.Text])
		parsed, err := t.parseTokens()
		if err != nil {
			return nil, err
		}
		quotedList = append(quotedList, parsed)
		return List(quotedList...), nil

	default:
		// if it"s not a list then it"s an atom
		return makeAtom(token)
	}
}

// makeAtom determines the correct type for an atom
// This will need extending as and when more primative types are added
func makeAtom(t token) (lispVal, error) {
	switch t.Tag {
	case "STRING":
		return string(t.Text[1 : len(t.Text)-1]), nil

	case "INT", "FLOAT":
		f, _ := strconv.ParseFloat(t.Text, 64)
		return float64(f), nil

	case "COMPLEX", "COMPLEX_PURE":
		return nil, fmt.Errorf("Complex numbers not implemented yet!")

	case "BOOL":
		if t.Text == "#t" {
			return true, nil
		}
		return false, nil

	case "KEYWORD":
		return KEYWORD(t.Text[1:]), nil

	case "SYMBOL":
		return SYMBOL(t.Text), nil

	default:
		return nil, fmt.Errorf("Unable to parse input: %v", t.Text)
	}
}

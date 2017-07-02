package gigl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regex objects for constructing atoms
// (Ensure that keywords and symbols contain valid characters)
var (
	reFloat, _    = regexp.Compile(`-?\d+\.\d+`)
	reInt, _      = regexp.Compile(`-?\d+`)
	reComp, _     = regexp.Compile(`-?\d+\.?\d*[+-]\d+\.?\d*j`)
	reCompPure, _ = regexp.Compile(`-?\d+\.?\d*j`)
	// reKeyword, _  = regexp.Compile(`:[^()[\]{}\s\#,\.]+(?=[\)\]}\s])?`)
	// reSymbol, _   = regexp.Compile(`[^()[\]{}\s\#,\.]+(?=[\)\]}\s])?`)
	quotes = map[string]SYMBOL{
		"'": SYMBOL("quote"), "`": SYMBOL("quasiquote"),
		",": SYMBOL("unquote"), ",@": SYMBOL("unquote-splicing"),
	}
)

// read a string and convert it to values we can work with
func read(s string) lispVal {
	tokens := tokenise(s)
	return parse(&tokens)
}

// split an input string into individual tokens, padding around parens & quotes
// This should probably be replaced with proper tokenisation using regex...
func tokenise(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	// s = strings.Replace(s, "[", "[ ", -1)
	// s = strings.Replace(s, "]", " ]", -1)
	// s = strings.Replace(s, "{", "{ ", -1)
	// s = strings.Replace(s, "}", " }", -1)
	s = strings.Replace(s, "'", " ' ", -1)
	s = strings.Replace(s, "`", " ` ", -1)
	s = strings.Replace(s, ",", " , ", -1)
	s = strings.Replace(s, ", @", " ,@ ", -1)

	split := strings.Split(s, " ")
	tokens := make([]string, 0)

	// Remove empty strings
	for _, token := range split {
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// parse the token stream and convert to values
// NOTE :: at present, this will only parse a single, complete s-expression
func parse(tokens *[]string) lispVal {
	// NOTE :: need to dereference tokens so we can slice
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	switch token {
	case "(":
		// Start of a list so recuse and build it up
		lst := make([]lispVal, 0)
		for (*tokens)[0] != ")" {
			nextToken := parse(tokens)
			if nextToken != SYMBOL("") {
				lst = append(lst, nextToken)
			}
		}
		// Slice off that last paren
		*tokens = (*tokens)[1:]
		return List(lst...)

	case "'", "`", ",", ",@":
		// Something is being quoted or unquoted
		quotedList := make([]lispVal, 0)
		quotedList = append(quotedList, quotes[token])
		quotedList = append(quotedList, parse(tokens))
		return List(quotedList...)

	default:
		// if it"s not a list then it"s an atom
		atom, err := makeAtom(token)
		if err != nil {
			fmt.Println("PARSE ERROR => ", err)
			return nil
		}
		return atom
	}
}

// makeAtom determines the correct type for an atom
// This will need extending as and when more primative types are added
func makeAtom(token string) (lispVal, error) {
	switch {
	case token[0] == '"' && token[len(token)-1] == '"':
		return string(token[1 : len(token)-1]), nil

	case reInt.MatchString(token), reFloat.MatchString(token):
		f, _ := strconv.ParseFloat(token, 64)
		return float64(f), nil

	case token == "true", token == "#t":
		return true, nil

	case token == "false", token == "#f":
		return false, nil

	// case reKeyword.MatchString(token):
	// 	log.Println("tis a keyword!")
	// 	return KEYWORD(token[1:]), nil

	// case reSymbol.MatchString(token):
	// 	return SYMBOL(token), nil

	default:
		return SYMBOL(token), nil
		// return token, fmt.Errorf("Unable to parse input: %v", token)
	}
}

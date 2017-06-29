package gigl

import (
	"fmt"
	"strings"
)

/*
	Type constructors and helper functions for the REPL
*/

// This is a catchall interface for functions to use in order to allow
// dynamic typing...I hope!
type lispVal interface{}

// Convert a lispVal to a string
// TODO :: make this work nicely for the list type
func String(val lispVal) string {
	switch val := val.(type) {
	case []lispVal:
		// display the slice as a LISPy list
		lst := make([]string, len(val))
		for i, element := range val {
			lst[i] = String(element)
		}
		return "(" + strings.Join(lst, " ") + ")"

	case STRING:
		return fmt.Sprintf("\"%v\"", val)

	case nil:
		return ""

	default:
		// Try to just print the value
		return fmt.Sprint(val)
	}
}

// A lispFunc takes values and returns a value
type lispFunc func(...lispVal) lispVal

// A lisp comp takes values and returns a bool
type lispComp func(...lispVal) bool

// Only basic data types so far
type SYMBOL string

type STRING string

type KEYWORD string

type NUM float64

// TODO :: It would be nice to have the full numeric tower...
// type INT int64

// type FLOAT float64

// type FRACTION struct {
// 	numerator   int64
// 	denominator int64
// }

// type COMPLEX struct {
// 	real float64
// 	imag float64
// }

// A procedure that stores paramaters, the function body and an
// execution environment to be called in.
type procedure struct {
	params lispVal
	body   lispVal
	env    *environment
}

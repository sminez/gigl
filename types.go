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
func String(val lispVal) string {
	switch val := val.(type) {
	case LispList:
		return val.String()

	case []lispVal:
		// display the slice as a LISPy list
		lst := make([]string, len(val))
		for i, element := range val {
			lst[i] = String(element)
		}
		return "(" + strings.Join(lst, " ") + ")"

	case string:
		return fmt.Sprintf("\"%v\"", val)

	case float64:
		// Try to print ints correctly
		i := int64(val)
		if float64(i) == val {
			return fmt.Sprint(i)
		} else {
			return fmt.Sprint(val)
		}

	case bool:
		if val {
			return "#t"
		}
		return "#f"

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

type KEYWORD string

// TODO :: It would be nice to have the full numeric tower...
// type float64 float64

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

// Make a callable procedure
func makeProc(params, body lispVal, env *environment, e *evaluator) (func(...lispVal) (lispVal, error), error) {
	innerEnv := &environment{
		vals:  make(map[SYMBOL]lispVal),
		outer: env,
	}

	proc := func(args ...lispVal) (lispVal, error) {
		switch params := params.(type) {
		case []lispVal:
			// Bind a list of paramaters into the new environment
			for i, param := range params {
				innerEnv.vals[param.(SYMBOL)] = args[i]
			}

		case LispList:
			// Bind a list of paramaters into the new environment
			for i, param := range params.toSlice() {
				innerEnv.vals[param.(SYMBOL)] = args[i]
			}
		default:
			// Bind as a single argument
			innerEnv.vals[params.(SYMBOL)] = args
		}

		// evaluate the result in the new environment
		result, err := e.eval(body, innerEnv)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return proc, nil
}

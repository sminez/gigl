package gigl

import (
	"log"
	"math"
)

/*
	Builtin functions for the global environment.

	:: NOTE ::
	Most (if not all) functions are variadic and take/return lispVals (interface{})

	It should be noted that this is a LISP and as such, is dynamically typed
	and also throws caution to the wind in terms of passing higher order
	functions around.  Very few things are in place to prevent you shooting
	yourself in the foot...!
*/

/*
	Arithmetic operations
*/

// add together two or more numbers
func add(lst ...lispVal) lispVal {
	total := lst[0].(NUM)
	for _, value := range lst[1:] {
		total += value.(NUM)
	}
	return total
}

// subtract two or more numbers in succession
func sub(lst ...lispVal) lispVal {
	total := lst[0].(NUM)
	for _, value := range lst[1:] {
		total -= value.(NUM)
	}
	return total
}

// multiply two or more numbers in succession
func mul(lst ...lispVal) lispVal {
	total := lst[0].(NUM)
	for _, value := range lst[1:] {
		total *= value.(NUM)
	}
	return total
}

// divide two or more numbers in succession
func div(lst ...lispVal) lispVal {
	total := lst[0].(NUM)
	for _, value := range lst[1:] {
		total /= value.(NUM)
	}
	return total
}

/*
	LISP list manipulations
*/

// construct a new list by prepending the new element
func cons(lst ...lispVal) lispVal {
	return append([]lispVal{lst[0]}, lst[1:]...)
}

// return the first element of a list
func car(lst ...lispVal) lispVal {
	return lst[0].([]lispVal)[0]
}

// everything but the first element of a list
func cdr(lst ...lispVal) lispVal {
	return lst[0].([]lispVal)[1:]
}

/*
	Higher order functions
*/

// map a function over a collection
// func mapfunc(f lispFunc, col []lispVal) lispVal {
func mapfunc(args ...lispVal) lispVal {
	f := args[0].(func(...lispVal) lispVal)
	col := args[1].([]lispVal)

	result := make([]lispVal, len(col))
	for i, element := range col {
		result[i] = f(element)
	}
	return result
}

// use a boolean function to filter a collection
// func filter(f lispComp, col []lispVal) lispVal {
func filter(args ...lispVal) lispVal {
	f := args[0].(func(...lispVal) bool)
	col := args[1].([]lispVal)

	result := make([]lispVal, 0)
	for _, element := range col {
		if f(element) {
			result = append(result, element)
		}
	}
	return result
}

// collapse a collection from the left using a binary
// operator and an optional accumulator value
// func foldl(f lispFunc, col []lispVal, acc ...lispVal) lispVal {
func foldl(args ...lispVal) lispVal {
	var result lispVal

	f := args[0].(func(...lispVal) lispVal)
	col := args[1].([]lispVal)

	switch l := len(args); l {
	case 2:
		result = col[0]
		col = col[1:]
	case 3:
		// Use the provided accumulator if there is one
		result = args[2].(lispVal)
	}

	for _, element := range col {
		result = f(result, element)
	}
	return result
}

/*
	Sequence functions
*/

// Python style range
func makeRange(args ...lispVal) lispVal {
	var min, max, step int

	switch l := len(args); l {
	case 1:
		min, max, step = 0, int(args[0].(NUM)), 1
	case 2:
		min, max, step = int(args[0].(NUM)), int(args[1].(NUM)), 1
	case 3:
		min, max, step = int(args[0].(NUM)), int(args[1].(NUM)), int(args[2].(NUM))
	default:
		log.Println("invalid args for range")
	}

	r := make([]lispVal, int(math.Ceil(float64(max-min)/float64(step))))
	for i := range r {
		r[i] = NUM(min + (step * i))
	}
	return r
}

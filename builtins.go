package gigl

import "log"

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
func mapfunc(f lispFunc, col []lispVal) lispVal {
	result := make([]lispVal, len(col))
	for i, element := range col {
		result[i] = f(element)
	}
	return result
}

// use a boolean function to filter a collection
func filter(f lispComp, col []lispVal) lispVal {
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
func foldl(f lispFunc, col []lispVal, acc ...lispVal) lispVal {
	var result lispVal
	if len(acc) == 1 {
		result = acc[0]
	}
	if len(acc) > 1 {
		log.Println("acc must be a value")
	}

	for _, element := range col {
		result = f(result, element)
	}
	return result
}

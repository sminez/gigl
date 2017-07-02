package gigl

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"reflect"
)

/*
	Builtin functions for the global environment.

	:: NOTE ::
	Most (if not all) functions are variadic and take/return lispVals (interface{})

	It should be noted that this is a LISP and as such, is dynamically typed
	and also throws caution to the wind in terms of passing higher order
	functions around.  Very few things are in place to prevent you shooting
	yourself in the foot...!

	TODO :: look at https://golang.org/pkg/math/big/ for arbitrary precision values
*/

// helper to convert to arbitray precision bignums
func getNum(l lispVal) (big.Float, error) {
	var f big.Float

	switch l := l.(type) {
	case uint8:
		f.SetInt64(int64(l))
	case int8:
		f.SetInt64(int64(l))
	case uint32:
		f.SetInt64(int64(l))
	case int32:
		f.SetInt64(int64(l))
	case uint64:
		f.SetUint64(l)
	case int64:
		f.SetInt64(l)
	case int:
		f.SetInt64(int64(l))
	case float32:
		f.SetFloat64(float64(l))
	case float64:
		f.SetFloat64(float64(l))
	default:
		return f, fmt.Errorf("Non-numeric argument: %v", l)
	}
	return f, nil
}

func getFloat(l lispVal) (float64, error) {
	val, ok := l.(float64)
	if !ok {
		return 0, fmt.Errorf("Non-numeric argument: %v", l)
	}
	return val, nil
}

/*
	Arithmetic operations
*/

// add together two or more numbers
func add(lst ...lispVal) (lispVal, error) {
	total, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}

	for _, value := range lst[1:] {
		v, err := getFloat(value)
		if err != nil {
			return nil, err
		}
		total += v
	}
	return total, nil
}

// subtract two or more numbers in succession
func sub(lst ...lispVal) (lispVal, error) {
	total, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}

	for _, value := range lst[1:] {
		v, err := getFloat(value)
		if err != nil {
			return nil, err
		}
		total -= v
	}
	return total, nil
}

// multiply two or more numbers in succession
func mul(lst ...lispVal) (lispVal, error) {
	total, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}

	for _, value := range lst[1:] {
		v, err := getFloat(value)
		if err != nil {
			return nil, err
		}
		total *= v
	}
	return total, nil
}

// divide two or more numbers in succession
func div(lst ...lispVal) (lispVal, error) {
	total, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}

	for _, value := range lst[1:] {
		v, err := getFloat(value)
		if err != nil {
			return nil, err
		}
		total /= v
	}
	return total, nil
}

// compute the remainder on division
func mod(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return float64(int64(a) % int64(b)), nil
}

/*
	Numeric Comparisons
*/
func lessThan(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a < b, nil
}

func lessThanOrEqual(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a <= b, nil
}

func greaterThan(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a > b, nil
}

func greaterThanOrEqual(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a >= b, nil
}

func equal(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a == b, nil
}

func notEqual(lst ...lispVal) (lispVal, error) {
	a, err := getFloat(lst[0])
	if err != nil {
		return nil, err
	}
	b, err := getFloat(lst[1])
	if err != nil {
		return nil, err
	}
	return a != b, nil
}

func isEqual(lst ...lispVal) (lispVal, error) {
	return reflect.DeepEqual(lst[0], lst[1]), nil
}

func null(lst ...lispVal) (lispVal, error) {
	list, ok := lst[0].([]lispVal)
	if !ok {
		return nil, fmt.Errorf("null? can only be called on lists")
	}
	return len(list) == 0, nil
}

/*
	LISP list manipulations
*/

// construct a new list by prepending the new element
func cons(lst ...lispVal) (lispVal, error) {
	switch lst[1].(type) {
	case []lispVal:
		return List(append([]lispVal{lst[0]}, lst[1].([]lispVal)...)...), nil
	default:
		return nil, fmt.Errorf("The second argument to cons must be a list")
	}
}

// return the first element of a list
func car(lst ...lispVal) (lispVal, error) {
	// l, ok := lst[0].([]lispVal)
	// if !ok {
	// 	return nil, fmt.Errorf("car called on an atom")
	// }
	// return l[0], nil
	l, ok := lst[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("car called on an atom")
	}
	return l.Head(), nil
}

// everything but the first element of a list
func cdr(lst ...lispVal) (lispVal, error) {
	// l, ok := lst[0].([]lispVal)
	// if !ok {
	// 	return nil, fmt.Errorf("cdr called on an atom")
	// }
	// return List(l[1:]), nil
	l, ok := lst[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("cdr called on an atom")
	}
	return l.Tail(), nil
}

/*
	Sequence functions
*/

// Python style range
func makeRange(args ...lispVal) (lispVal, error) {
	var (
		min  = int64(0)
		max  = int64(0)
		step = int64(1)
	)

	switch l := len(args); l {
	case 1:
		fmax, err := getFloat(args[0])
		if err != nil {
			return nil, err
		}
		max = int64(fmax)

	case 2:
		fmin, err := getFloat(args[0])
		if err != nil {
			return nil, err
		}
		min = int64(fmin)

		fmax, err := getFloat(args[1])
		if err != nil {
			return nil, err
		}
		max = int64(fmax)

	case 3:
		fmin, err := getFloat(args[0])
		if err != nil {
			return nil, err
		}
		min = int64(fmin)

		fmax, err := getFloat(args[1])
		if err != nil {
			return nil, err
		}
		max = int64(fmax)
		fstep, err := getFloat(args[2])
		if err != nil {
			return nil, err
		}
		step = int64(fstep)

	default:
		log.Println("invalid args for range")
	}

	r := make([]lispVal, int64(math.Ceil(float64(max-min)/float64(step))))
	for i := range r {
		r[i] = float64(min + (step * int64(i)))
	}
	return List(r...), nil
}

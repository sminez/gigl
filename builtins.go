package gigl

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"reflect"
)

/*
 * Internal versions of list operations for easier use from Go
 */

// Construct a new list by prepending a new element
// NOTE :: Does not fit the API for gigl builtins)
func consInternal(v lispVal, lst *LispList) *LispList {
	newList := NewList(v)
	newList.root.next = lst.root
	newList.length = lst.length + 1
	return newList
}

// Append two lists together, creating a new list
// NOTE :: Does not fit the API for gigl builtins)
func lispAppendInternal(l1, l2 *LispList) *LispList {
	s1 := l1.toSlice()
	s2 := l2.toSlice()
	return List(append(s1, s2...)...)
}

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
	Type checks
*/

func isInt(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	f, ok := lst[0].(float64)
	if !ok {
		return false, nil
	}
	return float64(int64(f)) == f, nil
}

func isFloat(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(float64)
	return ok, nil
}

func isString(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(string)
	return ok, nil
}

func isBool(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(bool)
	return ok, nil
}

func isSymbol(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(SYMBOL)
	return ok, nil
}

func isKeyword(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(KEYWORD)
	return ok, nil
}

// something is only a list if it contains items
func isList(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	_, ok := lst[0].(*LispList)
	return ok, nil
}

func isPair(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type check on non-atom: %v", lst)
	}
	l, ok := lst[0].(*LispList)
	if !ok {
		return false, nil
	}
	return l.Len() > 0, nil
}

func isNull(lst ...lispVal) (lispVal, error) {
	list, ok := lst[0].(*LispList)
	if !ok {
		fmt.Println(list)
		return nil, fmt.Errorf("null? can only be called on lists")
	}
	return list.Len() == 0, nil
}

func str(lst ...lispVal) (lispVal, error) {
	if len(lst) != 1 {
		return nil, fmt.Errorf("Type conversion on list: %v", lst)
	}
	return String(lst[0]), nil
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

/*
	LISP list manipulations
*/

// Construct a new list by prepending a new element
func cons(lst ...lispVal) (lispVal, error) {
	if len(lst) != 2 {
		return nil, fmt.Errorf("Cons takes two arguments")
	}

	newList := NewList(lst[0])
	oldList, ok := lst[1].(*LispList)
	if ok {
		newList.root.next = oldList.root
		newList.length = oldList.length + 1
	} else {
		newList = List(lst[0], lst[1])
		// return nil, fmt.Errorf("The second argument to cons must be a list")
	}
	return newList, nil
}

// Append several lists together, creating a new list
func lispAppend(lst ...lispVal) (lispVal, error) {
	slices := []lispVal{}

	// extract all of the other lists
	for _, l := range lst {
		switch l.(type) {
		case *LispList:
			slices = append(slices, l.(*LispList).toSlice()...)
		default:
			return nil, fmt.Errorf("Arguments to append must lists")
		}
	}

	return List(slices...), nil
}

// return the first element of a list
func car(lst ...lispVal) (lispVal, error) {
	l, ok := lst[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("car called on an atom")
	}
	return l.Head(), nil
}

// everything but the first element of a list
func cdr(lst ...lispVal) (lispVal, error) {
	l, ok := lst[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("cdr called on an atom")
	}
	return l.Tail(), nil
}

// length of a list
func lispLength(lst ...lispVal) (lispVal, error) {
	switch lst[0].(type) {
	case string:
		return len(lst[0].(string)), nil
	case *LispList:
		return float64(lst[0].(*LispList).Len()), nil
	default:
		return nil, fmt.Errorf("len called on a non-sequence")
	}
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

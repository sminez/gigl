package gigl

// An environment is a map of symbols to values that we can look up
// bindings in, along with a reference to the enclosing environment
// that we can backtrack to if we can't find something.
type environment struct {
	vals  map[SYMBOL]lispVal
	outer *environment
}

// Find attempts to find the closest environment that contains the
// requested symbol.
func (e *environment) find(sym SYMBOL) *environment {
	_, inEnvironment := e.vals[sym]
	if inEnvironment {
		return e
	}

	if e.outer != nil {
		return e.outer.find(sym)
	}

	return nil
}

// newGlobalEnvironment constructs a new global environment with the
// predefined builtin functions.
// NOTE :: builtins are found in builtin.go
func newGlobalEnvironment(e evaluator) *environment {
	return &environment{
		map[SYMBOL]lispVal{
			// arithmetic operators
			"+": add,
			"-": sub,
			"*": mul,
			"/": div,
			// "abs": abs,
			// LISP list manipulations
			"cons": cons,
			":":    cons,
			"car":  car,
			"head": car,
			"cdr":  cdr,
			"tail": cdr,
			// higher order functions
			"map":    mapfunc,
			"filter": filter,
			"foldl":  foldl,
		},
		nil,
	}
}

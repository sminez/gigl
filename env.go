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
	_, known := e.vals[sym]
	if known {
		return e
	}

	if e.outer != nil {
		return e.outer.find(sym)
	}

	// Return an empty environment to the caller
	return nil
}

// newGlobalEnvironment constructs a new global environment with the
// predefined builtin functions.
// NOTE :: builtins are found in builtin.go
func newGlobalEnvironment(e evaluator) *environment {
	return &environment{
		map[SYMBOL]lispVal{
			// higher order functions
			// NOTE :: using LISP versions in prelude.go
			// "map":    mapfunc,
			// "filter": filter,
			// "foldl":  foldl,
			// sequence functions
			"+":      add,
			"-":      sub,
			"*":      mul,
			"/":      div,
			"%":      mod,
			"<":      lessThan,
			"<=":     lessThanOrEqual,
			">":      greaterThan,
			">=":     greaterThanOrEqual,
			"=":      equal,
			"equal?": isEqual,
			"null?":  null,
			"cons":   cons,
			":":      cons,
			"car":    car,
			"head":   car,
			"cdr":    cdr,
			"tail":   cdr,
			"range":  makeRange,
		},
		nil,
	}
}

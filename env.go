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
func newGlobalEnvironment(e Evaluator) *environment {
	return &environment{
		map[SYMBOL]lispVal{
			"+":        add,
			"-":        sub,
			"*":        mul,
			"/":        div,
			"%":        mod,
			"modulo":   mod,
			"<":        lessThan,
			"<=":       lessThanOrEqual,
			">":        greaterThan,
			">=":       greaterThanOrEqual,
			"=":        equal,
			"!=":       notEqual,
			"eq?":      isEqual,
			"null?":    isNull,
			"int?":     isInt,
			"float?":   isFloat,
			"string?":  isString,
			"symbol?":  isSymbol,
			"keyword?": isKeyword,
			"list?":    isList,
			"pair?":    isPair,
			"car":      car,
			"cdr":      cdr,
			"head":     car,
			"tail":     cdr,
			"len":      lispLength,
			"cons":     cons,
			"append":   lispAppend,
			"range":    makeRange,
			"str":      str,
		},
		nil,
	}
}

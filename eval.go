package gigl

/*
	This is the main eval/apply loop as described in SICP
*/

import "fmt"

type evaluator struct {
	globalEnv *environment
}

func NewEvaluator() *evaluator {
	e := evaluator{}
	e.globalEnv = newGlobalEnvironment(e)
	return &e
}

// eval evaluates an expression in an environment
func (e *evaluator) eval(expression lispVal, env *environment) (lispVal, error) {
	var (
		result lispVal
		err    error
	)

	// Ensure that we always have an execution environment!
	if env == nil {
		env = e.globalEnv
	}

	for {
		switch expr := expression.(type) {
		case float64, string, bool, KEYWORD:
			// Just return the value as is
			return expr, nil

		case SYMBOL:
			// Find what this symbol refers to and return that
			location := env.find(expr)
			if location != nil {
				return location.vals[expr], nil
			}
			return nil, fmt.Errorf("Unknown symbol: %v", expr)

		case *LispList:
			// Pull off the head of the list and see what we need to do
			head, rest := expr.popHead()
			head, ok := head.(SYMBOL)
			if !ok {
				err = fmt.Errorf("Unknown procedure: %v", expr)
				return nil, err
			}

			switch head.(SYMBOL) {
			case "quote":
				// return the second element of the list unevaluated
				return rest.Head(), nil

			case "quasiquote":
				// recursively expand any quasi-quoted expressions
				return e.expandQuasiQuote(rest.Head(), env)
				// if err != nil {
				// 	return nil, err
				// }
				// return unquotedExp, nil
				// return e.eval(unquotedExp, env)

			case "unquote", "unquote-splicing":
				// This is handled inside of expandQuasiQuote
				return nil, fmt.Errorf("Cannot unquote outside of a quasi-quoted expression.")

			case "if":
				// Evaluate the conditional and cast to a bool
				check, rest := rest.popHead()
				check, err := e.eval(check, env)
				if err != nil {
					return nil, err
				}
				trueBranch, rest := rest.popHead()
				if check.(bool) {
					// evaluate the true branch
					return e.eval(trueBranch, env)
				} else {
					// evaluate the false branch or return nil
					falseBranch, _ := rest.popHead()
					if falseBranch != nil {
						return e.eval(falseBranch, env)
					}
					return nil, nil
				}

			case "set!":
				// find this symbol in its environment and update it
				sym, rest := rest.popHead()
				sym, ok := sym.(SYMBOL)
				if !ok {
					err = fmt.Errorf("Attempt to set non-symbol: %v", expr)
					return nil, err
				}
				if env.find(sym.(SYMBOL)) == nil {
					err = fmt.Errorf("Attempt to set! a new symbol: use define instead.")
					return nil, err
				}

				value, rest := rest.popHead()
				result, err = e.eval(value, env)
				if err != nil {
					return nil, err
				}
				env.find(sym.(SYMBOL)).vals[sym.(SYMBOL)] = result
				return nil, nil

			case "define":
				// Bind this symbol in the current environment
				sym, rest := rest.popHead()
				sym, ok := sym.(SYMBOL)
				if !ok {
					err = fmt.Errorf("Attempt to define non-symbol: %v", expr)
					return nil, err
				}
				if env.find(sym.(SYMBOL)) != nil {
					err = fmt.Errorf("Unable to redefine an existing symbol, use set!")
					return nil, err
				} else {
					value, _ := rest.popHead()
					result, err = e.eval(value, env)
					if err != nil {
						return nil, err
					}
					env.vals[sym.(SYMBOL)] = result
					return nil, nil
				}

			case "lambda", "λ":
				// Define a new procedure and return it
				params, rest := rest.popHead()
				body, rest := rest.popHead()
				return makeProc(params, body, env, e)

			case "defn":
				// define a new procedure and bind it to a symbol
				sym, rest := rest.popHead()
				sym, ok := sym.(SYMBOL)
				if !ok {
					err = fmt.Errorf("Attempt to define non-symbol: %v", expr)
					return nil, err
				}
				if env.find(sym.(SYMBOL)) != nil {
					err = fmt.Errorf("Unable to redefine an existing symbol, use set!")
					return nil, err
				} else {
					params, rest := rest.popHead()
					body, rest := rest.popHead()
					proc, err := makeProc(params, body, env, e)
					if err != nil {
						return nil, err
					}
					env.vals[sym.(SYMBOL)] = proc
					return nil, nil
				}

			case "begin":
				// Execute a collection of statements and return the
				// value of the last statement.
				var subExpr lispVal
				allButOne := rest.Len() - 1
				for i := 0; i < allButOne; i++ {
					subExpr, rest = rest.popHead()
					_, err = e.eval(subExpr, env)
					if err != nil {
						return nil, err
					}
				}
				// Loop back to evaluate the last form and return it
				expression = rest.Head()

			default:
				// Assume that the head is a callable and that the remaining
				// elements of the list are paramaters.
				var elem lispVal
				args := make([]lispVal, rest.Len())

				l := rest.Len()
				for i := 0; i < l; i++ {
					elem, rest = rest.popHead()
					result, err := e.eval(elem, env)
					if err != nil {
						return nil, err
					}
					args[i] = result
				}
				proc, err := e.eval(head, env)
				if err != nil {
					return nil, err
				}

				switch p := proc.(type) {
				case procedure:
					// This is a LISP procedure so create a new nested environment
					// to use as the execution environment and then evaluate the body
					innerEnv := &environment{
						vals:  make(map[SYMBOL]lispVal),
						outer: p.env,
					}

					switch params := p.params.(type) {
					case []lispVal:
						// Bind a list of paramaters into the new environment
						for i, param := range params {
							innerEnv.vals[param.(SYMBOL)] = args[i]
						}
					default:
						// Bind as a single argument
						innerEnv.vals[params.(SYMBOL)] = args
					}

					// loop and evaluate the result in the new environment
					expression = p.body
					env = innerEnv

				default:
					// apply a built-in procedure to some arguments directly and return the result
					return e.apply(proc, args)
				}
			}

		default:
			err = fmt.Errorf("Unknown expression in input: %v", expr)
			return nil, err
		}
	}
}

// apply a procedure to a list of arguments and return the result
// NOTE: built-in/primative operations will execute without any outer environment,
//		 procedures will bind their arguments before executing their statements.
func (e *evaluator) apply(proc lispVal, args []lispVal) (lispVal, error) {
	switch p := proc.(type) {
	case func(...lispVal) (lispVal, error):
		return p(args...)

	case func(...lispVal) (bool, error):
		return p(args...)

	case func(lispVal, *LispList) *LispList:
		switch lst := args[1].(type) {
		case LispList:
			return p(args[0], &lst), nil

		case *LispList:
			return p(args[0], lst), nil

		default:
			return nil, fmt.Errorf("Not a list: %v", args[1])
		}

	case func(*LispList, *LispList) *LispList:
		return p(args[0].(*LispList), args[1].(*LispList)), nil

	default:
		return nil, fmt.Errorf("Unknown procedure type: %v\n%v", p, args)
	}
}

// Expand quasi-quotes: expand `x -> 'x   `,x -> x   `(,@x y) -> (append x y)
func (e *evaluator) expandQuasiQuote(expression lispVal, env *environment) (lispVal, error) {
	switch expr := expression.(type) {
	case *LispList:
		// Make sure we aren't splicing a list into the head position of the new list
		if expr.Head() == SYMBOL("unquote-splicing") {
			return nil, fmt.Errorf("Can't splice at the head of a list: %v", expr)
		}

		// Collecting things up in a slice is conceptually easier to think about
		// when compared the repeated appends of lists or cons -> reverse.
		expandedList := make([]lispVal, 0)

		// Iterate through the terms and evaluate anything that has been unquoted
		element, originalList := expr.popHead()
		for {
			if originalList.Len() == 0 && element == nil {
				return List(expandedList...), nil
			}

			switch element.(type) {
			case *LispList:
				if element.(*LispList).Len() < 2 {
					expandedList = append(expandedList, element)
				} else {
					head, tail := element.(*LispList).popHead()
					switch head {
					case SYMBOL("unquote"):
						// Check that we actually have something to unquote
						if tail.Len() == 0 {
							return nil, fmt.Errorf("Unquoting error: %v", expr)
						}
						// Pop off the head and evaulate it
						toUnquote := tail.Head()
						unquotedElement, err := e.eval(toUnquote, env)
						if err != nil {
							return nil, err
						}
						// If everything looks good, add it to the resulting list
						expandedList = append(expandedList, unquotedElement)

					case SYMBOL("unquote-splicing"):
						// Check that we actually have something to unquote
						if tail.Len() == 0 {
							return nil, fmt.Errorf("Unquoting error: %v", expr)
						}
						// Pop off the head and evaulate it
						toUnquote := tail.Head()
						unquotedElements, err := e.eval(toUnquote, env)
						if err != nil {
							return nil, err
						}
						unquotedElements, ok := unquotedElements.(*LispList)
						if !ok {
							return nil, fmt.Errorf("Cannot call unquote-splicing on non-list: %v", unquotedElements)
						}
						// If everything looks good, add it to the resulting list
						expandedList = append(expandedList, unquotedElements.(*LispList).toSlice()...)
					default:
						// append the element unevaluated
						expandedList = append(expandedList, element)
					}
				}
			default:
				// append the element unevaluated
				expandedList = append(expandedList, element)
			}
			element, originalList = originalList.popHead()
		}

	default:
		// Still quote any un-marked forms for quoting
		return List(SYMBOL("quote"), expression), nil
	}
}

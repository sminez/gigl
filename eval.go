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

	switch expr := expression.(type) {
	case float64, string, bool:
		// Just return the value as is
		return expr, nil

	case SYMBOL:
		// Find what this symbol refers to and return that
		location := env.find(expr)
		if location != nil {
			return location.vals[expr], nil
		}
		return nil, fmt.Errorf("Unknown symbol: %v", expr)

	case LispList:
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
			l := rest.Len()
			for i := 0; i < l; i++ {
				subExpr, rest = rest.popHead()
				result, err = e.eval(subExpr, env)
				if err != nil {
					return nil, err
				}
			}
			return result, nil

		default:
			// Assume that the head is a callable and that the remaining
			// elements of the list are paramaters. Eagerly evaluate the
			// paramaters and then send everything off to apply
			var elem lispVal
			params := make([]lispVal, rest.Len())

			l := rest.Len()
			for i := 0; i < l; i++ {
				elem, rest = rest.popHead()
				result, err := e.eval(elem, env)
				if err != nil {
					return nil, err
				}
				params[i] = result
			}
			result, err = e.eval(head, env)
			if err != nil {
				return nil, err
			}
			return e.apply(result, params)
		}

	default:
		err = fmt.Errorf("Unknown expression in input: %v", expr)
		return nil, err
	}
}

// apply a procedure to a list of arguments and return the result
// NOTE: built-in/primative operations will execute without any outer environment,
//		 procedures will bind their arguments before executing their statements.
func (e *evaluator) apply(proc lispVal, args []lispVal) (lispVal, error) {
	switch p := proc.(type) {
	case func(...lispVal) (lispVal, error):
		// This is a built-in/primitive variadic-function
		return p(args...)

	case func(...lispVal) (bool, error):
		// This is a built-in/primitive comparison function
		return p(args...)

	case procedure:
		// This is a LISP procedure so create a new nested environment
		// to use as the execution environment
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

		// evaluate the result in the new environment
		return e.eval(p.body, innerEnv)

	default:
		return nil, fmt.Errorf("Unknown procedure type: %v", p)
	}
}

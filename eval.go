package gigl

/*
	This is the main eval/apply loop as described in SICP
*/

import "log"

type evaluator struct {
	globalEnv *environment
}

func NewEvaluator() *evaluator {
	e := evaluator{}
	e.globalEnv = newGlobalEnvironment(e)
	return &e
}

// eval evaluates an expression in an environment
func (e *evaluator) eval(expression lispVal, env *environment) lispVal {
	var result lispVal

	// Ensure that we always have an execution environment!
	if env == nil {
		env = e.globalEnv
	}

	switch expr := expression.(type) {
	case float64, string:
		// Just return the value as is
		result = expr

	case SYMBOL:
		// Find what this symbol refers to and return that
		result = env.find(expr).vals[expr]

	case []lispVal:
		// Pull off the head of the list and see what we need to do
		head, _ := expr[0].(SYMBOL)

		switch head {
		case "quote":
			// return the second element of the list unevaluated
			result = expr[1]

		case "if":
			// Evaluate the conditional and cast to a bool
			if e.eval(expr[1], env).(bool) {
				// evaluate the true branch
				result = e.eval(expr[2], env)
			} else {
				// evaluate the false branch or return nil
				if len(expr) == 4 {
					result = e.eval(expr[3], env)
				} else {
					result = nil
				}
			}

		case "set!":
			// find this symbol in its environment and update it
			// TODO:: do we need to worry about scope here?
			sym := expr[1].(SYMBOL)
			env.find(sym).vals[sym] = e.eval(expr[2], env)
			result = nil

		case "define":
			// Bind this symbol in the current environment
			sym := expr[1].(SYMBOL)
			if env.find(sym) != nil {
				log.Println("Unable to redefine an existing symbol, use set!")
			}

			env.vals[expr[1].(SYMBOL)] = e.eval(expr[2], env)
			result = nil

		case "lambda", "λ":
			// Define a new procedure and return it
			result = makeProc(expr[1], expr[2], env, e)

		case "defn":
			// define a new procedure and bind it to a symbol
			sym := expr[1].(SYMBOL)
			if env.find(sym) != nil {
				log.Println("unable to redefine an existing symbol, use set!")
			}

			proc := makeProc(expr[2], expr[3], env, e)
			env.vals[expr[1].(SYMBOL)] = proc
			result = nil

		case "begin":
			// Execute a collection of statements and return the
			// value of the last statement.
			for _, subExpr := range expr[1:] {
				result = e.eval(subExpr, env)
			}

		default:
			// Assume that expr[0] is a callable and that the remaining
			// elements of the list are paramaters. Eagerly evaluate the
			// paramaters and then send everything off to apply
			rawElements := expr[1:]
			params := make([]lispVal, len(rawElements))
			for i, elem := range rawElements {
				params[i] = e.eval(elem, env)
			}
			result = e.apply(e.eval(expr[0], env), params)
		}

	default:
		// If we don't know what to do, tell the user!
		log.Println("Unknown expression in input: ", expr)
	}

	return result
}

// apply a procedure to a list of arguments and return the result
// NOTE: built-in/primative operations will execute without any outer environment,
//		 procedures will bind their arguments before executing their statements.
func (e *evaluator) apply(proc lispVal, args []lispVal) lispVal {
	var result lispVal

	switch p := proc.(type) {
	case func(...lispVal) lispVal:
		// This is a built-in/primitive variadic-function
		result = p(args...)

	case func(...lispVal) bool:
		// This is a built-in/primitive comparison function
		result = p(args...)

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
		result = e.eval(p.body, innerEnv)

	default:
		log.Println("Unknown procedure type", p)
	}
	return result
}

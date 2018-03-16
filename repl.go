package gigl

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

var (
	InPrompt  = "λ > "
	OutPrompt = "   "
	input     string
	prevInput string
)

// REPL is the read-eval-print-loop
func REPL() {
	evaluator := NewEvaluator()
	tokeniser := NewTokeniser()

	// Load the prelude
	fmt.Printf("((Welcome to GIGL!)\n  (Loading prelude...)\n")
	for _, proc := range prelude {
		parsed, err := tokeniser.read(proc)
		if err != nil {
			panic(fmt.Sprintf("\n\nError in prelude!\n%v\n\n", err))
		}
		evaluator.eval(parsed, nil)
	}
	fmt.Println("  (...done!))")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 InPrompt,
		HistoryFile:            "/tmp/gigl-repl",
		DisableAutoSaveHistory: true,
	})

	// If we can't create the REPL we're boned...
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	previousInput := ""

	for {
		input, err := rl.Readline()
		if err != nil {
			fmt.Println(err)
			break
		}

		// Prepend the previous user input if there is any
		if previousInput != "" {
			// input = previousInput + " " + strings.TrimSpace(input)
			input = previousInput + "\n" + input
		}

		if len(input) > 0 {
			if !hasMatchingParens(input) {
				rl.SetPrompt(OutPrompt)
				previousInput = input
				continue
			}

			rl.SetPrompt(InPrompt)
			rl.SaveHistory(input)

			parsed, parseErr := tokeniser.read(input)
			if parseErr != nil {
				fmt.Printf("PARSE ERROR:\n%v\n=> %v\n\n", input, parseErr)
				previousInput = ""
				continue
			}
			result, evalErr := evaluator.eval(parsed, nil)
			if evalErr != nil {
				fmt.Printf("ERROR => %v\n\n", evalErr)
				previousInput = ""
			} else {
				res := String(result)
				if res != "" {
					fmt.Println(OutPrompt, res)
				}
				previousInput = ""
			}
		}
	}
}

// Check that we have a complete s-expression
// NOTE :: This _will_ fail if string literals have unmatched parens...
func hasMatchingParens(input string) bool {
	parenOpen := strings.Count(input, "(")
	parenClose := strings.Count(input, ")")
	return parenOpen == parenClose
}

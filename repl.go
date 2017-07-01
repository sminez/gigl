package gigl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	InPrompt  = "λ > "
	OutPrompt = "... "
	input     string
	prevInput string
)

// REPL is the read-eval-print-loop
// TODO: look at https://github.com/chzyer/readline
func REPL() {
	scanner := bufio.NewScanner(os.Stdin)
	evaluator := NewEvaluator()

	fmt.Printf("((Welcome to GIGL!)\n  (Loading prelude...)\n")
	// Load the prelude
	for _, proc := range prelude {
		evaluator.eval(read(proc), nil)
	}
	fmt.Println("  (...done!))")

	for {
		prevInput = ""
		fmt.Print(InPrompt)

		for {
			scanner.Scan()
			input = prevInput + " " + scanner.Text()

			if len(input) > 0 {
				if hasMatchingParens(input) {
					result, err := evaluator.eval(read(input), nil)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(OutPrompt, String(result))
					}
					break
				} else {
					prevInput = input
					fmt.Print(OutPrompt)
				}
			}
		}
	}
}

// Check that we have a complete s-expression
// NOTE :: This _will_ fail if string literals have unmatched parens
func hasMatchingParens(input string) bool {
	parenOpen := strings.Count(input, "(")
	parenClose := strings.Count(input, ")")
	return parenOpen == parenClose
}

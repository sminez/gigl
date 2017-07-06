package main

import (
	"fmt"

	"github.com/sminez/gigl"
)

func main() {
	lst := gigl.List(
		gigl.SYMBOL("let"),
		gigl.List(gigl.List(gigl.SYMBOL("arg"), gigl.SYMBOL("val")), gigl.SYMBOL("...")),
		gigl.SYMBOL("body"), gigl.SYMBOL("..."),
	)

	fmt.Println("Template: ", lst.String())

	pat, err := gigl.NewMatchPattern(lst)
	if err != nil {
		fmt.Println(err)
		return
	}

	target := gigl.List(
		gigl.SYMBOL("let"),
		gigl.List(gigl.List(gigl.SYMBOL("a"), 2), gigl.List(gigl.SYMBOL("b"), 4), gigl.List(gigl.SYMBOL("c"), 6)),
		gigl.List(gigl.SYMBOL("do")), gigl.List(gigl.SYMBOL("some")),
		gigl.List(gigl.SYMBOL("stuff!")),
	)

	fmt.Println("Target: ", target.String())

	matching := pat.Matches(target)
	fmt.Println("Successful match: ", matching)
	pat.PrintMatch()
}

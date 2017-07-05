package main

import (
	"fmt"

	"github.com/sminez/gigl"
)

func main() {
	lst := gigl.List(
		gigl.List(gigl.SYMBOL("a"), gigl.SYMBOL("b")),
		gigl.SYMBOL("..."),
	)

	pat, err := gigl.NewMatchPattern(lst)
	if err != nil {
		fmt.Println(err)
		return
	}

	target := gigl.List(gigl.List(1, 2), gigl.List(3, 4), gigl.List(5, 6))

	matching := pat.Matches(target)
	fmt.Println("Successful match: ", matching)
	pat.PrintMatch()
}

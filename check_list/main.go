package main

import (
	"fmt"

	"github.com/sminez/gigl"
)

func main() {
	l := gigl.NewList(5)
	fmt.Printf("%v\n", l)
	l2 := gigl.Cons(7, l)
	fmt.Printf("%v\n", l2)
	l3 := gigl.List(1, 2, 3, 4, 5)
	fmt.Printf("%v\n", l3)
	head := l3.Head()
	fmt.Printf("head is %v\n", head)
	tail := l3.Tail()
	fmt.Printf("tail is %v\n", tail)
	fmt.Println(l3.Len())
	fmt.Println(l3.Tail().Len())
}

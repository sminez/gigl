package gigl

import "strings"

/*
	To get slicing to work we need this:

	type Sequencer interface {
		Mean() float64
		Slice(start, end int) Sequencer
	}

*/

// A pair in a singly linked list
type Pair struct {
	Value lispVal // the value stored in this pair
	next  *Pair   // the next pair in the list
}

// A VERY simple singly linked list implementation based on the
// stdlib container/list doubly linked list
type LispList struct {
	root   *Pair // the first value in this list
	length int   // a cached, known length for the list
}

// New returns an initialized list.
func NewList(head lispVal) *LispList {
	return &LispList{
		root:   &Pair{Value: head},
		length: 1,
	}
}

// Front returns the first element value of list l or nil.
func (l *LispList) Head() lispVal {
	if l.root != nil {
		return l.root.Value
	}
	return &LispList{}
}

// Tail return a new list comprising of all elements but the first
func (l *LispList) Tail() *LispList {
	if l.root != nil {
		return &LispList{
			root:   l.root.next,
			length: l.length - 1,
		}
	}
	return &LispList{}
}

// Return the head and tail of the list
func (l *LispList) popHead() (lispVal, *LispList) {
	return l.Head(), l.Tail()
}

func (l *LispList) toSlice() []lispVal {
	lst := make([]lispVal, l.length)
	pair := l.root
	for i := 0; i < l.length; i++ {
		lst[i] = pair.Value
		pair = pair.next
	}
	return lst
}

func (l LispList) String() string {
	lst := make([]string, l.length)
	pair := l.root
	for i := 0; i < l.length; i++ {
		lst[i] = String(pair.Value)
		pair = pair.next
	}
	return "(" + strings.Join(lst, " ") + ")"
}

// Len returns the length of a list
func (l *LispList) Len() int {
	return l.length
}

// List is a repeated cons to build up a list
func List(vals ...lispVal) *LispList {
	if len(vals) == 0 {
		// Empty list
		return &LispList{}
	}

	lst := NewList(vals[0])
	pair := lst.root
	for i := 1; i < len(vals); i++ {
		newPair := &Pair{Value: vals[i]}
		pair.next = newPair
		pair = newPair
	}
	lst.length = len(vals)
	return lst
}

// check if something is an empty list
func isEmptyList(v lispVal) bool {
	lst, ok := v.(*LispList)
	if !ok {
		return false
	}
	return lst.Len() == 0
}

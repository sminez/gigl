package gigl

import "strings"

// A node in a sngly linked list
type Node struct {
	Value lispVal // the value stored in this node
	next  *Node   // the next node in the list
}

// A VERY simple singly linked list implementation based on the
// stdlib container/list doubly linked list
type LispList struct {
	root   *Node // the first value in this list
	length int   // a cached, known length for the list
}

// New returns an initialized list.
func NewList(head lispVal) *LispList {
	return &LispList{
		root:   &Node{Value: head},
		length: 1,
	}
}

// Front returns the first element value of list l or nil.
func (l *LispList) Head() lispVal {
	return l.root.Value
}

// Tail return a new list comprising of all elements but the first
func (l *LispList) Tail() *LispList {
	return &LispList{
		root:   l.root.next,
		length: l.length - 1,
	}
}

func (l LispList) String() string {
	// display the slice as a LISPy list
	lst := make([]string, l.length)
	node := l.root
	for i := 0; i < l.length; i++ {
		lst[i] = String(node.Value)
		node = node.next
	}
	return "(" + strings.Join(lst, " ") + ")"

}

// Len returns the length of a list
func (l *LispList) Len() int {
	return l.length
}

// Construct a new list by prepending a new element
func Cons(v lispVal, lst *LispList) *LispList {
	newList := NewList(v)
	newList.root.next = lst.root
	newList.length = lst.length + 1
	return newList
}

// List is a repeated cons to build up a list
func List(vals ...lispVal) *LispList {
	lst := NewList(vals[0])
	node := lst.root
	for i := 1; i < len(vals); i++ {
		newNode := &Node{Value: vals[i]}
		node.next = newNode
		node = newNode
	}
	lst.length = len(vals)
	return lst
}

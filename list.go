package gigl

// A node in a sngly linked list
type Node struct {
	Value lispVal // the value stored in this node
	next  *Node   // the next node in the list
}

// Next returns the next list node or nil.
func (n *Node) Next() *Node {
	return n.next
}

// A VERY simple singly linked list implementation based on the
// stdlib container/list doubly linked list
type List struct {
	root   *Node // the first value in this list
	length int   // a cached, known length for the list
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = nil
	l.length = 0
	return l
}

// New returns an initialized list.
func NewList() *List {
	return new(List).Init()
}

// Return the current length of the list O(1)
func (l *List) Len() int {
	return l.length
}

// Front returns the first element of list l or nil.
func (l *List) Head() *Node {
	if l.length == 0 {
		return nil
	}
	return l.root
}

// Front returns the first element of list l or nil.
func (l *List) Tail() *Node {
	if l.length < 2 {
		return nil
	}
	return l.root
}

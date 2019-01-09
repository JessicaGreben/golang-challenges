package main

func main() {}

type node struct {
	value    string
	next     *node
	previous *node
}

type dblLinkedList struct {
	root *node
}

func createDblLinkedList(value string) *dblLinkedList {
	root := node{
		value: value,
	}
	return &dblLinkedList{
		root: &root,
	}
}

func (ll *dblLinkedList) find(value string) *node {
	current := ll.root

	for {
		if current.value == value {
			return current
		}

		// If we are at the end of the linked list then that value does
		// not exist in the linked list.
		if current.next == nil {
			return nil
		}
		current = current.next
	}
}

func (ll *dblLinkedList) findEnd() *node {
	current := ll.root
	for {
		if current.next == nil {
			return current
		}
		current = current.next
	}
}

func (ll *dblLinkedList) insert(newNode *node, location string) {
	switch location {
	case "first":
		first := ll.root
		ll.root = newNode
		newNode.next = first
	case "last":
		endNode := ll.findEnd()
		endNode.next = newNode
	}
}

func (ll *dblLinkedList) length() int {
	count := 0
	current := ll.root
	for {
		count++
		if current.next == nil {
			return count
		}
		current = current.next
	}
}

func (ll *dblLinkedList) delete(node *node) {
	previous := node.previous
	next := node.next

	// If we are deleting the end node,
	// then the previous node will now be the end node.
	if next == nil {
		previous.next = nil
	}

	previous.next = next
	next.previous = previous
}

func (ll *dblLinkedList) reverse() {
	current := ll.root

	for {
		oldNext := current.next
		oldPrevious := current.previous

		switch {

		// We are at the root node of the original list.
		case oldPrevious == nil:
			current.next = nil
			current.previous = oldNext

		// We are at the terminal node of the original list.
		case oldNext == nil:
			current.next = oldPrevious
			current.previous = nil
			return

		// We are at a middle node of the original list.
		default:
			current.next = oldPrevious
			current.previous = oldNext
		}
	}
}

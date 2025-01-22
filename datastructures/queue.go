package datastructures

import "mysticsquare/square"

type Queue []*square.MysticSquare

// remove an item from the queue
func (q *Queue) Pop() (item any) {
	old := *q
	n := len(old)
	item = old[0]
	old[0] = nil
	*q = old[1:n]
	return
}

// add an item to the queue
func (q *Queue) Push(item any) {
	if newItem, ok := item.(square.MysticSquare); ok {
		*q = append(*q, &newItem)
		return
	}
	panic("Item is not of type square.MysticSquare")
}

// len of queue
func (q *Queue) Len() int {
	return len(*q)
}

// check if queue is empty
func (q *Queue) Empty() bool {
	return q.Len() == 0
}

// combines checking if the queue has any items and retrieving the item from the queue in a single call.
func (q *Queue) Process() (current square.MysticSquare, hasItem bool) {
	hasItem = !q.Empty()
	if hasItem {
		if data, ok := q.Pop().(*square.MysticSquare); ok {
			current = *data
		} else {
			panic("queue item failed type conversion to *square.MysticSquare")
		}
	} else {
		current = nil
	}
	return
}

// create a new queue
func NewMysticSquareQueue() (q *Queue) {
	q = &Queue{}
	return
}

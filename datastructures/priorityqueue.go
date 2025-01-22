package datastructures

import (
	"container/heap"
	"mysticsquare/square"
)

// a priority queue item
type MysticSquareItem struct {
	Msquare  square.MysticSquare
	priority int
	index    int
}

// create a new item
func NewMysticSquareItem(Msquare square.MysticSquare, priority int) (item *MysticSquareItem) {
	item = &MysticSquareItem{Msquare: Msquare, priority: priority}
	return
}

// get the priority of the mystic square item
func (item MysticSquareItem) Priority() int {
	return item.priority
}

// priority queue struct
type PriorityQueue []*MysticSquareItem

// check the len of the queue
func (pq PriorityQueue) Len() int {

	return len(pq)
}

// compares two items. Used by container/heap
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority() < pq[j].Priority()
}

// swaps two items. Used by container/heap
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// add an item to the priority queue
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*MysticSquareItem)
	item.index = n
	*pq = append(*pq, item)
}

// remove an item from the priority queue
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// check if the priority queue is empty
func (pq *PriorityQueue) Empty() bool {
	return pq.Len() == 0
}

// check an item exists and return it to the user.
func (pq *PriorityQueue) Process() (current *MysticSquareItem, itemExists bool) {
	itemExists = !pq.Empty()
	if itemExists {
		if data, ok := heap.Pop(pq).(*MysticSquareItem); ok {
			current = data
		} else {
			panic("failed to convert item to *MysticSquareItem")
		}
	} else {
		current = nil
	}

	return
}

// public external function
func (pq *PriorityQueue) Update(item *MysticSquareItem, priority int) {
	if idx := item.index; idx >= 0 {
		pq.update(item, priority)
	}
}

// internal update function
func (pq *PriorityQueue) update(item *MysticSquareItem, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

// create a new priority queue
func NewMysticSquarePriorityQueue() (pq *PriorityQueue) {
	pq = &PriorityQueue{}
	return
}

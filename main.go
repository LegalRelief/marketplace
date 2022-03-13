package LegalRelief

import (
	"container/heap"
	"errors"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	volume   uint64
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

type Exchange struct {
	pq     map[string]PriorityQueue
	prices map[string]float64
}

var errNotFound = errors.New("could not find stock")

func (e *Exchange) peek(stock string) (*Item, error) {
	if val, ok := e.pq[stock]; ok {
		return heap.Pop(&val).(*Item), nil
	}
	return nil, errNotFound
}

type TransactionStatus string

const (
	SUCCESS TransactionStatus = "Success"
	PENDING TransactionStatus = "Pending"
	FAILURE TransactionStatus = "Failure"
)

func (e *Exchange) buy(stock string, volume uint64) TransactionStatus {

}

func NewExchange() *Exchange {
	return &Exchange{
		pq: make(map[string]PriorityQueue, 0),
	}
}

func main() {
	e := NewExchange()
	
}

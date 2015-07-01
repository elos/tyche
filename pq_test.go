package tyche

import (
	"container/heap"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestPriorityQueueBasic(t *testing.T) {
	size := 1000
	maxPriority := int64(100000)

	pq := make(PriorityQueue, size)

	for i := 0; i < size; i++ {
		p, err := rand.Int(rand.Reader, big.NewInt(maxPriority))
		if err != nil {
			t.Fatal("Failed to generate *big.Int: %s", err)
		}

		f := float64(p.Uint64())

		pq[i] = &Item{
			value:    i,
			priority: f,
			index:    i,
		}
	}

	heap.Init(&pq)

	last := float64(maxPriority + 1)

	t.Logf("Priority Queue: %+v", pq)
	t.Logf("Items:")
	for _, v := range pq {
		t.Logf(" * %+v", *v)
	}

	count := size

	for pq.Len() > 0 {
		t.Logf("Last priority: %f", last)
		t.Logf("This priority: %f", pq.max().priority)

		if !(pq.max().priority <= last) {
			t.Fatalf("Expected priority queue to be in descending order of priority")
		}

		v := heap.Pop(&pq)
		t.Logf("Value popped off queue %+v", v)
		item, ok := v.(*Item)
		if !ok {
			t.Fatalf("Expected to be able to type case interface{} to *Item")
		}

		if !(item.priority <= last) {
			t.Fatal("Poppoed item's priority is not less than previos")
		}

		last = item.priority

		if count < 0 {
			t.Fatal("Appears to be more items in pq than there should be")
		}

		count -= 1
	}
}

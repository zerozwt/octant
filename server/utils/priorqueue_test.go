package utils

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	q := PriorityQueue(func(a, b *int) bool { return *a < *b })
	ptr := func(i int) *int { return &i }

	q.Put(ptr(3))
	q.Put(ptr(1))
	q.Put(ptr(2))
	q.Put(ptr(4))
	q.Put(ptr(5))
	q.Put(ptr(4))

	ans := []int{1, 2, 3, 4, 4, 5}

	if len(q.data) != len(ans) {
		t.Error(fmt.Errorf("data len %d ans len %d", len(q.data), len(ans)))
	}

	for idx, value := range ans {
		if *q.data[idx] != value {
			t.Error(fmt.Errorf("idx %d value %d ans %d", idx, *q.data[idx], ans))
		}
	}
}

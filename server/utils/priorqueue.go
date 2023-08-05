package utils

type PriorQueue[T any] struct {
	data []*T
	cmp  func(*T, *T) bool
}

func PriorityQueue[T any](cmp func(*T, *T) bool) *PriorQueue[T] {
	return &PriorQueue[T]{
		data: make([]*T, 0),
		cmp:  cmp,
	}
}

func (q *PriorQueue[T]) Put(value *T) {
	idx := q.binarySearch(value, 0, len(q.data))
	if idx == len(q.data) {
		q.data = append(q.data, value)
		return
	}
	tmp := make([]*T, 0, len(q.data)+1)
	tmp = append(tmp, q.data[:idx]...)
	tmp = append(tmp, value)
	tmp = append(tmp, q.data[idx:]...)
	q.data = tmp
}

func (q *PriorQueue[T]) Head() *T {
	if len(q.data) > 0 {
		return q.data[0]
	}
	return nil
}

func (q *PriorQueue[T]) Pop() *T {
	if len(q.data) == 0 {
		return nil
	}
	ret := q.data[0]
	q.data = q.data[1:]
	return ret
}

func (q *PriorQueue[T]) binarySearch(value *T, l, r int) int {
	if l == r {
		return l
	}
	mid := l + ((r - l) >> 1)
	if q.cmp(value, q.data[mid]) {
		return q.binarySearch(value, l, mid)
	}
	return q.binarySearch(value, mid+1, r)
}

// Package heap holds the runnable templates from heap-study-plan.md.
//
// Inside this package, the identifier heap refers to the standard library's container/heap.
package heap

import (
	"container/heap"
	"sort"
)

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func findKthLargest(nums []int, k int) int {
	h := &MinHeap{}
	heap.Init(h)
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return (*h)[0]
}

type Pair struct {
	val   int
	count int
}

type MinHeapByCount []Pair

func (h MinHeapByCount) Len() int           { return len(h) }
func (h MinHeapByCount) Less(i, j int) bool { return h[i].count < h[j].count }
func (h MinHeapByCount) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeapByCount) Push(x any)        { *h = append(*h, x.(Pair)) }
func (h *MinHeapByCount) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func topKFrequent(nums []int, k int) []int {
	freq := map[int]int{}
	for _, n := range nums {
		freq[n]++
	}

	h := &MinHeapByCount{}
	heap.Init(h)
	for val, count := range freq {
		heap.Push(h, Pair{val, count})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	result := make([]int, h.Len())
	for i := range result {
		result[i] = (*h)[i].val
	}
	return result
}

type MedianFinder struct {
	lo *MaxHeap
	hi *MinHeap
}

func newMedianFinder() MedianFinder {
	lo, hi := &MaxHeap{}, &MinHeap{}
	heap.Init(lo)
	heap.Init(hi)
	return MedianFinder{lo, hi}
}

func (mf *MedianFinder) AddNum(num int) {
	heap.Push(mf.lo, num)
	heap.Push(mf.hi, heap.Pop(mf.lo))

	if mf.hi.Len() > mf.lo.Len() {
		heap.Push(mf.lo, heap.Pop(mf.hi))
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if mf.lo.Len() > mf.hi.Len() {
		return float64((*mf.lo)[0])
	}
	return float64((*mf.lo)[0]+(*mf.hi)[0]) / 2.0
}

type ListItem struct {
	val     int
	listIdx int
	nodeIdx int
}

type MinHeapByVal []ListItem

func (h MinHeapByVal) Len() int           { return len(h) }
func (h MinHeapByVal) Less(i, j int) bool { return h[i].val < h[j].val }
func (h MinHeapByVal) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeapByVal) Push(x any)        { *h = append(*h, x.(ListItem)) }
func (h *MinHeapByVal) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func mergeKSorted(lists [][]int) []int {
	h := &MinHeapByVal{}
	heap.Init(h)
	for i, list := range lists {
		if len(list) > 0 {
			heap.Push(h, ListItem{list[0], i, 0})
		}
	}

	var result []int
	for h.Len() > 0 {
		top := heap.Pop(h).(ListItem)
		result = append(result, top.val)

		if top.nodeIdx+1 < len(lists[top.listIdx]) {
			next := lists[top.listIdx][top.nodeIdx+1]
			heap.Push(h, ListItem{next, top.listIdx, top.nodeIdx + 1})
		}
	}
	return result
}

func minMeetingRooms(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

	endTimes := &MinHeap{}
	heap.Init(endTimes)

	for _, iv := range intervals {
		if endTimes.Len() > 0 && (*endTimes)[0] <= iv[0] {
			heap.Pop(endTimes)
		}
		heap.Push(endTimes, iv[1])
	}
	return endTimes.Len()
}

package heap

import (
	"container/heap"
	"slices"
	"testing"
)

func TestMinMaxHeapOrder(t *testing.T) {
	in := []int{5, 2, 8, 1, 9, 3}

	minH := &MinHeap{}
	maxH := &MaxHeap{}
	for _, n := range in {
		heap.Push(minH, n)
		heap.Push(maxH, n)
	}

	var gotMin, gotMax []int
	for minH.Len() > 0 {
		gotMin = append(gotMin, heap.Pop(minH).(int))
		gotMax = append(gotMax, heap.Pop(maxH).(int))
	}
	if want := []int{1, 2, 3, 5, 8, 9}; !slices.Equal(gotMin, want) {
		t.Errorf("MinHeap pop order = %v, want %v", gotMin, want)
	}
	if want := []int{9, 8, 5, 3, 2, 1}; !slices.Equal(gotMax, want) {
		t.Errorf("MaxHeap pop order = %v, want %v", gotMax, want)
	}
}

func TestHeapInit(t *testing.T) {
	h := &MinHeap{5, 2, 8}
	heap.Init(h)
	if (*h)[0] != 2 {
		t.Errorf("after Init, top = %d, want 2", (*h)[0])
	}
}

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		nums    []int
		k, want int
	}{
		{[]int{3, 2, 1, 5, 6, 4}, 2, 5},
		{[]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
		{[]int{1}, 1, 1},
		{[]int{2, 1}, 2, 1},
	}
	for _, tt := range tests {
		if got := findKthLargest(tt.nums, tt.k); got != tt.want {
			t.Errorf("findKthLargest(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
		}
	}
}

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{1, 1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{[]int{1}, 1, []int{1}},
		{[]int{4, 4, 4, 4, 5, 5, 5, 6, 6, 7}, 3, []int{4, 5, 6}},
	}
	for _, tt := range tests {
		got := topKFrequent(tt.nums, tt.k)
		slices.Sort(got) // heap order is not sorted order
		if !slices.Equal(got, tt.want) {
			t.Errorf("topKFrequent(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
		}
	}
}

func TestMedianFinder(t *testing.T) {
	mf := newMedianFinder()
	steps := []struct {
		add  int
		want float64
	}{
		{5, 5},
		{15, 10},
		{1, 5},
		{3, 4},
		{8, 5},
		{-2, 4},
	}
	for _, s := range steps {
		mf.AddNum(s.add)
		if got := mf.FindMedian(); got != s.want {
			t.Fatalf("after adding %d, median = %v, want %v", s.add, got, s.want)
		}
	}
}

func TestMergeKSorted(t *testing.T) {
	tests := []struct {
		name  string
		lists [][]int
		want  []int
	}{
		{"three lists", [][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{"with empties", [][]int{{}, {2}, {}}, []int{2}},
		{"no lists", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeKSorted(tt.lists); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMeetingRooms(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{"overlap", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
		{"disjoint", [][]int{{7, 10}, {2, 4}}, 1},
		{"back to back reuses a room", [][]int{{1, 5}, {5, 10}}, 1},
		{"all overlap", [][]int{{1, 10}, {2, 9}, {3, 8}}, 3},
		{"none", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minMeetingRooms(tt.intervals); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

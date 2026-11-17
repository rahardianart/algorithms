// Package greedy holds the runnable templates from greedy-study-plan.md.
package greedy

import (
	"container/heap"
	"sort"
)

// maxSubArray is Kadane's algorithm: the largest sum of a non-empty contiguous subarray.
func maxSubArray(nums []int) int {
	cur, best := nums[0], nums[0]
	for _, x := range nums[1:] {
		cur = max(x, cur+x)
		best = max(best, cur)
	}
	return best
}

// canReach reports whether the last index is reachable when nums[i] is the max jump from i.
func canReach(nums []int) bool {
	farthest := 0
	for i, jump := range nums {
		if i > farthest {
			return false
		}
		farthest = max(farthest, i+jump)
	}
	return true
}

// maxSatisfied matches each item to at most one need, where an item satisfies a need if item >= need.
func maxSatisfied(needs, items []int) int {
	sort.Ints(needs)
	sort.Ints(items)

	matched := 0
	for _, item := range items {
		if matched < len(needs) && item >= needs[matched] {
			matched++
		}
	}
	return matched
}

func partitionLabels(s string) []int {
	var last [26]int
	for i := 0; i < len(s); i++ {
		last[s[i]-'a'] = i
	}

	var sizes []int
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		end = max(end, last[s[i]-'a'])
		if i == end {
			sizes = append(sizes, end-start+1)
			start = i + 1
		}
	}
	return sizes
}

type maxHeap []int

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// minRefuelStops drives toward target, remembering every station passed and only
// refueling (from the biggest one) when the tank would otherwise run dry.
func minRefuelStops(target, startFuel int, stations [][]int) int {
	passed := &maxHeap{}
	fuel, stops, i := startFuel, 0, 0

	for fuel < target {
		for i < len(stations) && stations[i][0] <= fuel {
			heap.Push(passed, stations[i][1])
			i++
		}
		if passed.Len() == 0 {
			return -1
		}
		fuel += heap.Pop(passed).(int)
		stops++
	}
	return stops
}

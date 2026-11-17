package slidingwindow

import (
	"slices"
	"testing"
)

func TestMaxSumFixed(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{2, 1, 5, 1, 3, 2}, 3, 9},
		{[]int{2, 3, 4, 1, 5}, 2, 7},
		{[]int{-3, -1, -2}, 2, -3},
		{[]int{5}, 1, 5},
	}
	for _, tt := range tests {
		if got := maxSumFixed(tt.nums, tt.k); got != tt.want {
			t.Errorf("maxSumFixed(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
		}
	}
}

func TestLongestSumAtMost(t *testing.T) {
	tests := []struct {
		nums  []int
		limit int
		want  int
	}{
		{[]int{3, 1, 2, 1, 1, 4, 1}, 5, 4},
		{[]int{1, 1, 1, 1}, 10, 4},
		{[]int{6, 7}, 5, 0},
		{nil, 5, 0},
	}
	for _, tt := range tests {
		if got := longestSumAtMost(tt.nums, tt.limit); got != tt.want {
			t.Errorf("longestSumAtMost(%v, %d) = %d, want %d", tt.nums, tt.limit, got, tt.want)
		}
	}
}

func TestMinSubArrayLen(t *testing.T) {
	tests := []struct {
		target int
		nums   []int
		want   int
	}{
		{7, []int{2, 3, 1, 2, 4, 3}, 2},
		{4, []int{1, 4, 4}, 1},
		{11, []int{1, 1, 1, 1, 1, 1, 1, 1}, 0},
		{15, []int{1, 2, 3, 4, 5}, 5},
	}
	for _, tt := range tests {
		if got := minSubArrayLen(tt.target, tt.nums); got != tt.want {
			t.Errorf("minSubArrayLen(%d, %v) = %d, want %d", tt.target, tt.nums, got, tt.want)
		}
	}
}

func TestCountSum(t *testing.T) {
	brute := func(nums []int, goal int) (atMost, exactly int) {
		for i := range nums {
			s := 0
			for j := i; j < len(nums); j++ {
				s += nums[j]
				if s <= goal {
					atMost++
				}
				if s == goal {
					exactly++
				}
			}
		}
		return
	}
	inputs := []struct {
		nums []int
		goal int
	}{
		{[]int{1, 0, 1, 0, 1}, 2},
		{[]int{0, 0, 0, 0, 0}, 0},
		{[]int{2, 1, 3, 0, 1}, 3},
		{[]int{5}, 1},
	}
	for _, in := range inputs {
		wantAtMost, wantExactly := brute(in.nums, in.goal)
		if got := countSumAtMost(in.nums, in.goal); got != wantAtMost {
			t.Errorf("countSumAtMost(%v, %d) = %d, want %d", in.nums, in.goal, got, wantAtMost)
		}
		if got := countSumExactly(in.nums, in.goal); got != wantExactly {
			t.Errorf("countSumExactly(%v, %d) = %d, want %d", in.nums, in.goal, got, wantExactly)
		}
	}
}

func TestMinSlidingWindow(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{4, 2, 12, 3, 8, 1, 7}, 3, []int{2, 2, 3, 1, 1}},
		{[]int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{-1, -3, -3, -3, 3, 3}},
		{[]int{5, 5, 5}, 2, []int{5, 5}},
		{[]int{9}, 1, []int{9}},
	}
	for _, tt := range tests {
		if got := minSlidingWindow(tt.nums, tt.k); !slices.Equal(got, tt.want) {
			t.Errorf("minSlidingWindow(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
		}
	}
}

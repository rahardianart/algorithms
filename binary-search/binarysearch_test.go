package binarysearch

import (
	"slices"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	nums := []int{-1, 0, 3, 5, 9, 12}
	tests := []struct{ target, want int }{
		{9, 4}, {-1, 0}, {12, 5}, {2, -1}, {13, -1}, {-5, -1},
	}
	for _, tt := range tests {
		if got := binarySearch(nums, tt.target); got != tt.want {
			t.Errorf("binarySearch(%d) = %d, want %d", tt.target, got, tt.want)
		}
	}
	if got := binarySearch(nil, 1); got != -1 {
		t.Errorf("binarySearch(nil) = %d, want -1", got)
	}
}

func TestLeftmostTrue(t *testing.T) {
	nums := []int{1, 3, 5, 6}
	tests := []struct {
		name   string
		target int
		want   int
	}{
		{"present", 5, 2},
		{"between", 2, 1},
		{"past the end", 7, 4},
		{"before the start", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := leftmostTrue(0, len(nums), func(i int) bool { return nums[i] >= tt.target })
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSearchRange(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   []int
	}{
		{[]int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{[]int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{[]int{}, 0, []int{-1, -1}},
		{[]int{2, 2, 2}, 2, []int{0, 2}},
		{[]int{1}, 1, []int{0, 0}},
	}
	for _, tt := range tests {
		if got := searchRange(tt.nums, tt.target); !slices.Equal(got, tt.want) {
			t.Errorf("searchRange(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
		}
	}
}

func TestSearchRotated(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 5, 1},
		{[]int{1}, 0, -1},
		{[]int{3, 1}, 1, 1},
		{[]int{1, 2, 3, 4}, 4, 3}, // not rotated at all
	}
	for _, tt := range tests {
		if got := searchRotated(tt.nums, tt.target); got != tt.want {
			t.Errorf("searchRotated(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
		}
	}
}

func TestMinEatingSpeed(t *testing.T) {
	tests := []struct {
		piles []int
		h     int
		want  int
	}{
		{[]int{3, 6, 7, 11}, 8, 4},
		{[]int{30, 11, 23, 4, 20}, 5, 30},
		{[]int{30, 11, 23, 4, 20}, 6, 23},
		{[]int{1, 1, 1}, 10, 1},
	}
	for _, tt := range tests {
		if got := minEatingSpeed(tt.piles, tt.h); got != tt.want {
			t.Errorf("minEatingSpeed(%v, %d) = %d, want %d", tt.piles, tt.h, got, tt.want)
		}
	}
}

func TestSearchMatrix(t *testing.T) {
	matrix := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	tests := []struct {
		target int
		want   bool
	}{
		{3, true}, {16, true}, {60, true}, {1, true}, {13, false}, {0, false}, {61, false},
	}
	for _, tt := range tests {
		if got := searchMatrix(matrix, tt.target); got != tt.want {
			t.Errorf("searchMatrix(%d) = %v, want %v", tt.target, got, tt.want)
		}
	}
}

func TestSplitArray(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{7, 2, 5, 10, 8}, 2, 18},
		{[]int{1, 2, 3, 4, 5}, 2, 9},
		{[]int{1, 4, 4}, 3, 4},
		{[]int{5}, 1, 5},
	}
	for _, tt := range tests {
		if got := splitArray(tt.nums, tt.k); got != tt.want {
			t.Errorf("splitArray(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
		}
	}
}

package twopointer

import (
	"reflect"
	"slices"
	"testing"
)

func TestTwoSumSorted(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   [2]int
	}{
		{[]int{2, 7, 11, 15}, 9, [2]int{1, 2}},
		{[]int{2, 3, 4}, 6, [2]int{1, 3}},
		{[]int{-1, 0}, -1, [2]int{1, 2}},
		{[]int{1, 2, 3}, 100, [2]int{-1, -1}},
	}
	for _, tt := range tests {
		if got := twoSumSorted(tt.nums, tt.target); got != tt.want {
			t.Errorf("twoSumSorted(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
		}
	}
}

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"classic", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"all zeros deduped", []int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{"none", []int{0, 1, 1}, [][]int{}},
		{"too short", []int{0, 0}, [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := threeSum(tt.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxArea(t *testing.T) {
	tests := []struct {
		height []int
		want   int
	}{
		{[]int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{[]int{1, 1}, 1},
		{[]int{4, 3, 2, 1, 4}, 16},
	}
	for _, tt := range tests {
		if got := maxArea(tt.height); got != tt.want {
			t.Errorf("maxArea(%v) = %d, want %d", tt.height, got, tt.want)
		}
	}
}

func TestTrap(t *testing.T) {
	tests := []struct {
		height []int
		want   int
	}{
		{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{[]int{4, 2, 0, 3, 2, 5}, 9},
		{[]int{1, 2, 3}, 0},
		{nil, 0},
	}
	for _, tt := range tests {
		if got := trap(tt.height); got != tt.want {
			t.Errorf("trap(%v) = %d, want %d", tt.height, got, tt.want)
		}
	}
}

// inPlace checks a "return the new length k" function: k must match and nums[:k] must match.
func inPlace(t *testing.T, name string, fn func([]int) int, in, want []int) {
	t.Helper()
	nums := slices.Clone(in)
	k := fn(nums)
	if k != len(want) || !slices.Equal(nums[:k], want) {
		t.Errorf("%s(%v) = %d, %v; want %d, %v", name, in, k, nums[:k], len(want), want)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	inPlace(t, "removeDuplicates", removeDuplicates, []int{1, 1, 2}, []int{1, 2})
	inPlace(t, "removeDuplicates", removeDuplicates, []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, []int{0, 1, 2, 3, 4})
	inPlace(t, "removeDuplicates", removeDuplicates, []int{}, []int{})
}

func TestRemoveElement(t *testing.T) {
	rm := func(val int) func([]int) int { return func(n []int) int { return removeElement(n, val) } }
	inPlace(t, "removeElement", rm(3), []int{3, 2, 2, 3}, []int{2, 2})
	inPlace(t, "removeElement", rm(2), []int{0, 1, 2, 2, 3, 0, 4, 2}, []int{0, 1, 3, 0, 4})
	inPlace(t, "removeElement", rm(1), []int{1, 1}, []int{})
}

func TestRemoveDuplicatesII(t *testing.T) {
	inPlace(t, "removeDuplicatesII", removeDuplicatesII, []int{1, 1, 1, 2, 2, 3}, []int{1, 1, 2, 2, 3})
	inPlace(t, "removeDuplicatesII", removeDuplicatesII, []int{0, 0, 1, 1, 1, 1, 2, 3, 3}, []int{0, 0, 1, 1, 2, 3, 3})
	inPlace(t, "removeDuplicatesII", removeDuplicatesII, []int{1}, []int{1})
}

func TestMoveZeroes(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{[]int{0}, []int{0}},
		{[]int{1, 2}, []int{1, 2}},
	}
	for _, tt := range tests {
		nums := slices.Clone(tt.in)
		moveZeroes(nums)
		if !slices.Equal(nums, tt.want) {
			t.Errorf("moveZeroes(%v) = %v, want %v", tt.in, nums, tt.want)
		}
	}
}

func TestSortColors(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{2, 0, 2, 1, 1, 0}, []int{0, 0, 1, 1, 2, 2}},
		{[]int{2, 0, 1}, []int{0, 1, 2}},
		{[]int{2, 2, 0}, []int{0, 2, 2}}, // swapped-in value from high must be re-examined
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		nums := slices.Clone(tt.in)
		sortColors(nums)
		if !slices.Equal(nums, tt.want) {
			t.Errorf("sortColors(%v) = %v, want %v", tt.in, nums, tt.want)
		}
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		nums1 []int
		m     int
		nums2 []int
		want  []int
	}{
		{"interleaved", []int{1, 2, 3, 0, 0, 0}, 3, []int{2, 5, 6}, []int{1, 2, 2, 3, 5, 6}},
		{"nums2 empty", []int{1}, 1, []int{}, []int{1}},
		{"nums1 empty", []int{0}, 0, []int{1}, []int{1}},
		{"all of nums2 first", []int{4, 5, 0, 0}, 2, []int{1, 2}, []int{1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums1 := slices.Clone(tt.nums1)
			merge(nums1, tt.m, tt.nums2, len(tt.nums2))
			if !slices.Equal(nums1, tt.want) {
				t.Errorf("got %v, want %v", nums1, tt.want)
			}
		})
	}
}

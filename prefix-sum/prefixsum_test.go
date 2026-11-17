package prefixsum

import (
	"reflect"
	"slices"
	"testing"
)

func TestCountPairsWithSum(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{1, 5, 3, 3, 7, 1}, 6, 3}, // (1,5) (3,3) (5,1)
		{[]int{2, 2, 2}, 4, 3},
		{[]int{1, 2, 3}, 100, 0},
		{nil, 0, 0},
	}
	for _, tt := range tests {
		if got := countPairsWithSum(tt.nums, tt.target); got != tt.want {
			t.Errorf("countPairsWithSum(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
		}
	}
}

func TestGroupByCharSet(t *testing.T) {
	got := groupByCharSet([]string{"abc", "cab", "ab", "aabbc", "ba", "xyz"})
	want := [][]string{{"abc", "cab", "aabbc"}, {"ab", "ba"}, {"xyz"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestRangeSum(t *testing.T) {
	nums := []int{3, -1, 4, 1, 5}
	prefix := prefixSums(nums)
	if want := []int{0, 3, 2, 6, 7, 12}; !slices.Equal(prefix, want) {
		t.Fatalf("prefixSums = %v, want %v", prefix, want)
	}
	for l := range nums {
		for r := l; r < len(nums); r++ {
			want := 0
			for _, x := range nums[l : r+1] {
				want += x
			}
			if got := rangeSum(prefix, l, r); got != want {
				t.Errorf("rangeSum(%d, %d) = %d, want %d", l, r, got, want)
			}
		}
	}
}

func TestLongestSubarrayWithSum(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{1, -1, 5, -2, 3}, 3, 4},
		{[]int{-2, -1, 2, 1}, 1, 2},
		{[]int{1, 2, 3}, 7, 0},
		{[]int{0, 0, 0}, 0, 3},
	}
	for _, tt := range tests {
		if got := longestSubarrayWithSum(tt.nums, tt.k); got != tt.want {
			t.Errorf("longestSubarrayWithSum(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
		}
	}
}

func TestApplyRangeAdds(t *testing.T) {
	got := applyRangeAdds(5, [][3]int{{1, 3, 2}, {2, 4, 3}, {0, 2, -2}})
	want := []int{-2, 0, 3, 5, 3}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPrefix2D(t *testing.T) {
	grid := [][]int{
		{3, 0, 1, 4, 2},
		{5, 6, 3, 2, 1},
		{1, 2, 0, 1, 5},
		{4, 1, 0, 1, 7},
		{1, 0, 3, 0, 5},
	}
	p := newPrefix2D(grid)
	tests := []struct {
		r1, c1, r2, c2, want int
	}{
		{2, 1, 4, 3, 8},
		{1, 1, 2, 2, 11},
		{1, 2, 2, 4, 12},
		{0, 0, 0, 0, 3},
		{0, 0, 4, 4, 58},
	}
	for _, tt := range tests {
		if got := p.sumRegion(tt.r1, tt.c1, tt.r2, tt.c2); got != tt.want {
			t.Errorf("sumRegion(%d,%d,%d,%d) = %d, want %d", tt.r1, tt.c1, tt.r2, tt.c2, got, tt.want)
		}
	}
}

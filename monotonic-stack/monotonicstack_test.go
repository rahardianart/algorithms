package monotonicstack

import (
	"slices"
	"testing"
)

// bruteNext is the O(n²) definition the stack versions must agree with.
func bruteNext(nums []int, better func(a, b int) bool) []int {
	out := make([]int, len(nums))
	for i := range nums {
		out[i] = -1
		for j := i + 1; j < len(nums); j++ {
			if better(nums[j], nums[i]) {
				out[i] = nums[j]
				break
			}
		}
	}
	return out
}

func TestNextGreaterAndSmaller(t *testing.T) {
	inputs := [][]int{
		{2, 1, 2, 4, 3},
		{5, 4, 3, 2, 1},
		{1, 2, 3, 4, 5},
		{3, 3, 3},
		{},
	}
	greater := func(a, b int) bool { return a > b }
	smaller := func(a, b int) bool { return a < b }
	for _, in := range inputs {
		if got, want := nextGreaterElement(in), bruteNext(in, greater); !slices.Equal(got, want) {
			t.Errorf("nextGreaterElement(%v) = %v, want %v", in, got, want)
		}
		if got, want := nextSmallerElement(in), bruteNext(in, smaller); !slices.Equal(got, want) {
			t.Errorf("nextSmallerElement(%v) = %v, want %v", in, got, want)
		}
	}
}

func TestDailyTemperatures(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{[]int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{[]int{30, 60, 90}, []int{1, 1, 0}},
		{[]int{50, 50}, []int{0, 0}}, // equal is not warmer
	}
	for _, tt := range tests {
		if got := dailyTemperatures(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("dailyTemperatures(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestNextGreaterElements(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{1, 2, 1}, []int{2, -1, 2}},
		{[]int{1, 2, 3, 4, 3}, []int{2, 3, 4, -1, 4}},
		{[]int{5}, []int{-1}},
		{[]int{2, 2}, []int{-1, -1}},
	}
	for _, tt := range tests {
		if got := nextGreaterElements(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("nextGreaterElements(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestLargestRectangleArea(t *testing.T) {
	tests := []struct {
		heights []int
		want    int
	}{
		{[]int{2, 1, 5, 6, 2, 3}, 10},
		{[]int{2, 4}, 4},
		{[]int{1, 1, 1, 1}, 4},
		{[]int{5}, 5},
		{[]int{}, 0},
	}
	for _, tt := range tests {
		if got := largestRectangleArea(tt.heights); got != tt.want {
			t.Errorf("largestRectangleArea(%v) = %d, want %d", tt.heights, got, tt.want)
		}
	}
}

func TestRemoveKdigits(t *testing.T) {
	tests := []struct {
		num  string
		k    int
		want string
	}{
		{"1432219", 3, "1219"},
		{"10200", 1, "200"},
		{"10", 2, "0"},
		{"12345", 2, "123"}, // increasing: the loop never pops, the final trim does the work
		{"9", 1, "0"},
	}
	for _, tt := range tests {
		if got := removeKdigits(tt.num, tt.k); got != tt.want {
			t.Errorf("removeKdigits(%q, %d) = %q, want %q", tt.num, tt.k, got, tt.want)
		}
	}
}

func TestSumSubarrayMins(t *testing.T) {
	brute := func(arr []int) int {
		total := 0
		for i := range arr {
			m := arr[i]
			for j := i; j < len(arr); j++ {
				m = min(m, arr[j])
				total += m
			}
		}
		return total % 1_000_000_007
	}
	inputs := [][]int{
		{3, 1, 2, 4},
		{11, 81, 94, 43, 3},
		{2, 2, 2}, // repeated minimum: the > / >= asymmetry prevents double counting
		{1},
	}
	for _, in := range inputs {
		if got, want := sumSubarrayMins(in), brute(in); got != want {
			t.Errorf("sumSubarrayMins(%v) = %d, want %d", in, got, want)
		}
	}
}

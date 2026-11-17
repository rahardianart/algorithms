package dp

import "testing"

func TestFib(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for n, w := range want {
		if got := fibMemo(n, map[int]int{}); got != w {
			t.Errorf("fibMemo(%d) = %d, want %d", n, got, w)
		}
		if got := fibTab(n); got != w {
			t.Errorf("fibTab(%d) = %d, want %d", n, got, w)
		}
		if got := fibRolling(n); got != w {
			t.Errorf("fibRolling(%d) = %d, want %d", n, got, w)
		}
	}
	if got := fibMemo(90, map[int]int{}); got != 2880067194370816120 {
		t.Errorf("fibMemo(90) = %d — memo should make this instant", got)
	}
}

func TestClimbStairs(t *testing.T) {
	tests := []struct{ n, want int }{
		{1, 1}, {2, 2}, {3, 3}, {4, 5}, {5, 8}, {10, 89},
	}
	for _, tt := range tests {
		if got := climbStairs(tt.n); got != tt.want {
			t.Errorf("climbStairs(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestLengthOfLIS(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{[]int{0, 1, 0, 3, 2, 3}, 4},
		{[]int{7, 7, 7, 7}, 1},
		{[]int{5}, 1},
		{nil, 0},
	}
	for _, tt := range tests {
		if got := lengthOfLIS(tt.nums); got != tt.want {
			t.Errorf("lengthOfLIS(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestUniquePaths(t *testing.T) {
	tests := []struct{ m, n, want int }{
		{3, 7, 28}, {3, 2, 3}, {1, 1, 1}, {1, 5, 1},
	}
	for _, tt := range tests {
		if got := uniquePaths(tt.m, tt.n); got != tt.want {
			t.Errorf("uniquePaths(%d, %d) = %d, want %d", tt.m, tt.n, got, tt.want)
		}
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"abcde", "ace", 3},
		{"abc", "abc", 3},
		{"abc", "def", 0},
		{"", "abc", 0},
	}
	for _, tt := range tests {
		if got := longestCommonSubsequence(tt.a, tt.b); got != tt.want {
			t.Errorf("LCS(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestCanPartition(t *testing.T) {
	tests := []struct {
		nums []int
		want bool
	}{
		{[]int{1, 5, 11, 5}, true},
		{[]int{1, 2, 3, 5}, false},
		{[]int{2, 2}, true},
		{[]int{1}, false},
		{[]int{4}, false},
		{[]int{2, 6}, false}, // target 4: forward loop would wrongly use the 2 twice
	}
	for _, tt := range tests {
		if got := canPartition(tt.nums); got != tt.want {
			t.Errorf("canPartition(%v) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

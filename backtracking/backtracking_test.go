package backtracking

import (
	"reflect"
	"testing"
)

func TestCombine(t *testing.T) {
	tests := []struct {
		n, k int
		want [][]int
	}{
		{4, 2, [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}},
		{3, 3, [][]int{{1, 2, 3}}},
		{1, 1, [][]int{{1}}},
	}
	for _, tt := range tests {
		if got := combine(tt.n, tt.k); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("combine(%d, %d) = %v, want %v", tt.n, tt.k, got, tt.want)
		}
	}
	if got := len(combine(10, 4)); got != 210 {
		t.Errorf("len(combine(10, 4)) = %d, want C(10,4) = 210", got)
	}
}

func TestCombinationSum(t *testing.T) {
	tests := []struct {
		candidates []int
		target     int
		want       [][]int
	}{
		{[]int{2, 3, 6, 7}, 7, [][]int{{2, 2, 3}, {7}}},
		{[]int{2, 3, 5}, 8, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}},
		{[]int{2}, 1, nil},
	}
	for _, tt := range tests {
		if got := combinationSum(tt.candidates, tt.target); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("combinationSum(%v, %d) = %v, want %v", tt.candidates, tt.target, got, tt.want)
		}
	}
}

func TestPermuteUnique(t *testing.T) {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 1, 2}, [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}}},
		{[]int{1, 2, 3}, [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}},
		{[]int{2, 2, 2}, [][]int{{2, 2, 2}}},
	}
	for _, tt := range tests {
		if got := permuteUnique(tt.nums); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("permuteUnique(%v) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func TestPartition(t *testing.T) {
	got := partition("aab", isPalindrome)
	want := [][]string{{"a", "a", "b"}, {"aa", "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("palindromes of aab = %v, want %v", got, want)
	}

	short := func(p string) bool { return len(p) <= 2 }
	got = partition("abc", short)
	want = [][]string{{"a", "b", "c"}, {"a", "bc"}, {"ab", "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pieces of abc up to length 2 = %v, want %v", got, want)
	}
}

func TestTotalNQueens(t *testing.T) {
	want := map[int]int{1: 1, 2: 0, 3: 0, 4: 2, 5: 10, 6: 4, 8: 92}
	for n, w := range want {
		if got := totalNQueens(n); got != w {
			t.Errorf("totalNQueens(%d) = %d, want %d", n, got, w)
		}
	}
}

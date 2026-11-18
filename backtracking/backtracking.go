// Package backtracking holds the runnable templates from backtracking-study-plan.md.
package backtracking

import "sort"

// combine returns every k-element combination of 1..n, pruning branches that can't reach k.
func combine(n, k int) [][]int {
	var result [][]int
	var current []int

	var backtrack func(start int)
	backtrack = func(start int) {
		if len(current) == k {
			result = append(result, append([]int(nil), current...))
			return
		}
		need := k - len(current)
		for i := start; i <= n-need+1; i++ {
			current = append(current, i)
			backtrack(i + 1)
			current = current[:len(current)-1]
		}
	}
	backtrack(1)
	return result
}

// combinationSum returns every multiset of candidates (each reusable) summing to target.
func combinationSum(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	var result [][]int
	var current []int

	var backtrack func(start, remaining int)
	backtrack = func(start, remaining int) {
		if remaining == 0 {
			result = append(result, append([]int(nil), current...))
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				break
			}
			current = append(current, candidates[i])
			backtrack(i, remaining-candidates[i])
			current = current[:len(current)-1]
		}
	}
	backtrack(0, target)
	return result
}

// permuteUnique returns every distinct permutation of nums, which may contain duplicates.
func permuteUnique(nums []int) [][]int {
	sort.Ints(nums)
	var result [][]int
	current := make([]int, 0, len(nums))
	used := make([]bool, len(nums))

	var backtrack func()
	backtrack = func() {
		if len(current) == len(nums) {
			result = append(result, append([]int(nil), current...))
			return
		}
		for i := range nums {
			if used[i] {
				continue
			}
			if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			current = append(current, nums[i])
			backtrack()
			current = current[:len(current)-1]
			used[i] = false
		}
	}
	backtrack()
	return result
}

// partition returns every way to cut s into pieces that all satisfy isValid.
func partition(s string, isValid func(string) bool) [][]string {
	var result [][]string
	var pieces []string

	var backtrack func(start int)
	backtrack = func(start int) {
		if start == len(s) {
			result = append(result, append([]string(nil), pieces...))
			return
		}
		for end := start + 1; end <= len(s); end++ {
			piece := s[start:end]
			if !isValid(piece) {
				continue
			}
			pieces = append(pieces, piece)
			backtrack(end)
			pieces = pieces[:len(pieces)-1]
		}
	}
	backtrack(0)
	return result
}

// totalNQueens counts the ways to place n non-attacking queens on an n×n board.
func totalNQueens(n int) int {
	cols := make([]bool, n)
	diag := make([]bool, 2*n-1)     // r + c is constant along a "/" diagonal
	antiDiag := make([]bool, 2*n-1) // r - c + n - 1 is constant along a "\" diagonal
	count := 0

	var place func(r int)
	place = func(r int) {
		if r == n {
			count++
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || diag[r+c] || antiDiag[r-c+n-1] {
				continue
			}
			cols[c], diag[r+c], antiDiag[r-c+n-1] = true, true, true
			place(r + 1)
			cols[c], diag[r+c], antiDiag[r-c+n-1] = false, false, false
		}
	}
	place(0)
	return count
}

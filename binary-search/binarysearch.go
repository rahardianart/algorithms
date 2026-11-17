// Package binarysearch holds the runnable templates from binary-search-study-plan.md.
package binarysearch

import "slices"

func binarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// leftmostTrue returns the first index in [lo, hi) where condition is true, or hi if there is none.
func leftmostTrue(lo, hi int, condition func(int) bool) int {
	for lo < hi {
		mid := lo + (hi-lo)/2
		if condition(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func searchRange(nums []int, target int) []int {
	first := leftmostTrue(0, len(nums), func(i int) bool { return nums[i] >= target })
	if first == len(nums) || nums[first] != target {
		return []int{-1, -1}
	}
	last := leftmostTrue(0, len(nums), func(i int) bool { return nums[i] > target }) - 1
	return []int{first, last}
}

// searchRotated assumes distinct values.
func searchRotated(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		}

		if nums[lo] <= nums[mid] {
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else {
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}

func minEatingSpeed(piles []int, h int) int {
	hoursNeeded := func(speed int) int {
		hours := 0
		for _, p := range piles {
			hours += (p + speed - 1) / speed
		}
		return hours
	}

	lo, hi := 1, slices.Max(piles)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if hoursNeeded(mid) <= h {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func searchMatrix(matrix [][]int, target int) bool {
	rows, cols := len(matrix), len(matrix[0])
	lo, hi := 0, rows*cols-1

	for lo <= hi {
		mid := lo + (hi-lo)/2
		val := matrix[mid/cols][mid%cols]
		switch {
		case val == target:
			return true
		case val < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return false
}

func splitArray(nums []int, k int) int {
	canSplit := func(maxSum int) bool {
		pieces, cur := 1, 0
		for _, n := range nums {
			if cur+n > maxSum {
				pieces++
				cur = 0
			}
			cur += n
		}
		return pieces <= k
	}

	lo, hi := slices.Max(nums), sum(nums)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if canSplit(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

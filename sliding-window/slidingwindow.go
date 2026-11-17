// Package slidingwindow holds the runnable templates from sliding-window-study-plan.md.
package slidingwindow

import "math"

// maxSumFixed is the fixed-size template: the largest sum of any k consecutive elements.
func maxSumFixed(nums []int, k int) int {
	sum, best := 0, math.MinInt
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		if right >= k {
			sum -= nums[right-k]
		}
		if right >= k-1 {
			best = max(best, sum)
		}
	}
	return best
}

// longestSumAtMost is the "longest valid window" template. nums must be non-negative.
func longestSumAtMost(nums []int, limit int) int {
	left, sum, best := 0, 0, 0
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		for sum > limit {
			sum -= nums[left]
			left++
		}
		best = max(best, right-left+1)
	}
	return best
}

// minSubArrayLen is the "shortest valid window" template. nums must be positive.
func minSubArrayLen(target int, nums []int) int {
	left, sum := 0, 0
	best := len(nums) + 1
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		for sum >= target {
			best = min(best, right-left+1)
			sum -= nums[left]
			left++
		}
	}
	if best == len(nums)+1 {
		return 0
	}
	return best
}

// countSumAtMost counts subarrays whose sum is <= goal. nums must be non-negative.
func countSumAtMost(nums []int, goal int) int {
	if goal < 0 {
		return 0
	}
	left, sum, count := 0, 0, 0
	for right := 0; right < len(nums); right++ {
		sum += nums[right]
		for sum > goal {
			sum -= nums[left]
			left++
		}
		count += right - left + 1
	}
	return count
}

// countSumExactly uses exactly(k) = atMost(k) - atMost(k-1).
func countSumExactly(nums []int, goal int) int {
	return countSumAtMost(nums, goal) - countSumAtMost(nums, goal-1)
}

// minSlidingWindow returns the minimum of every window of size k using a monotonic deque.
func minSlidingWindow(nums []int, k int) []int {
	var result []int
	deque := []int{} // indices; their values increase from front to back

	for right := 0; right < len(nums); right++ {
		if len(deque) > 0 && deque[0] <= right-k {
			deque = deque[1:]
		}
		for len(deque) > 0 && nums[deque[len(deque)-1]] >= nums[right] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, right)

		if right >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}
	return result
}

// Package monotonicstack holds the runnable templates from monotonic-stack-study-plan.md.
package monotonicstack

func nextGreaterElement(nums []int) []int {
	result := make([]int, len(nums))
	for i := range result {
		result[i] = -1
	}
	stack := []int{} // indices, values decreasing from bottom to top

	for i, n := range nums {
		for len(stack) > 0 && nums[stack[len(stack)-1]] < n {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top] = n
		}
		stack = append(stack, i)
	}
	return result
}

// nextSmallerElement is nextGreaterElement with the comparison flipped: an increasing stack.
func nextSmallerElement(nums []int) []int {
	result := make([]int, len(nums))
	for i := range result {
		result[i] = -1
	}
	stack := []int{} // indices, values increasing from bottom to top

	for i, n := range nums {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > n {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top] = n
		}
		stack = append(stack, i)
	}
	return result
}

func dailyTemperatures(temperatures []int) []int {
	result := make([]int, len(temperatures))
	stack := []int{}

	for i, t := range temperatures {
		for len(stack) > 0 && temperatures[stack[len(stack)-1]] < t {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top] = i - top
		}
		stack = append(stack, i)
	}
	return result
}

func nextGreaterElements(nums []int) []int {
	n := len(nums)
	result := make([]int, n)
	for i := range result {
		result[i] = -1
	}
	stack := []int{}

	for i := 0; i < 2*n; i++ {
		idx := i % n
		for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[idx] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top] = nums[idx]
		}
		if i < n {
			stack = append(stack, idx)
		}
	}
	return result
}

func largestRectangleArea(heights []int) int {
	stack := []int{}
	maxArea := 0

	for i := 0; i <= len(heights); i++ {
		h := 0
		if i < len(heights) {
			h = heights[i]
		}
		for len(stack) > 0 && heights[stack[len(stack)-1]] >= h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			height := heights[top]

			width := i
			if len(stack) > 0 {
				width = i - stack[len(stack)-1] - 1
			}
			maxArea = max(maxArea, height*width)
		}
		stack = append(stack, i)
	}
	return maxArea
}

func removeKdigits(num string, k int) string {
	stack := []byte{}

	for i := 0; i < len(num); i++ {
		for k > 0 && len(stack) > 0 && stack[len(stack)-1] > num[i] {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, num[i])
	}
	stack = stack[:len(stack)-k]

	start := 0
	for start < len(stack)-1 && stack[start] == '0' {
		start++
	}
	result := string(stack[start:])
	if result == "" {
		return "0"
	}
	return result
}

func sumSubarrayMins(arr []int) int {
	const mod = 1_000_000_007
	n := len(arr)
	prevLess := make([]int, n)
	nextLessEq := make([]int, n)

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && arr[stack[len(stack)-1]] > arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			prevLess[i] = i + 1
		} else {
			prevLess[i] = i - stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && arr[stack[len(stack)-1]] >= arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			nextLessEq[i] = n - i
		} else {
			nextLessEq[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}

	total := 0
	for i := 0; i < n; i++ {
		total = (total + arr[i]*prevLess[i]*nextLessEq[i]) % mod
	}
	return total
}

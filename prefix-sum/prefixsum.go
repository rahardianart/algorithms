// Package prefixsum holds the runnable templates from prefix-sum-study-plan.md.
package prefixsum

// countPairsWithSum counts index pairs i < j with nums[i] + nums[j] == target.
func countPairsWithSum(nums []int, target int) int {
	seen := map[int]int{} // value -> how many times it appeared so far
	count := 0
	for _, x := range nums {
		count += seen[target-x]
		seen[x]++
	}
	return count
}

// groupByCharSet groups words that use exactly the same set of letters, in first-seen order.
func groupByCharSet(words []string) [][]string {
	groupOf := map[uint32]int{} // letter-set bitmask -> index into groups
	var groups [][]string
	for _, w := range words {
		var key uint32
		for i := 0; i < len(w); i++ {
			key |= 1 << (w[i] - 'a')
		}
		idx, ok := groupOf[key]
		if !ok {
			idx = len(groups)
			groupOf[key] = idx
			groups = append(groups, nil)
		}
		groups[idx] = append(groups[idx], w)
	}
	return groups
}

// prefixSums returns p where p[i] is the sum of nums[:i]; len(p) == len(nums)+1.
func prefixSums(nums []int) []int {
	prefix := make([]int, len(nums)+1)
	for i, x := range nums {
		prefix[i+1] = prefix[i] + x
	}
	return prefix
}

// rangeSum returns nums[l] + ... + nums[r] (inclusive) in O(1).
func rangeSum(prefix []int, l, r int) int {
	return prefix[r+1] - prefix[l]
}

// longestSubarrayWithSum returns the length of the longest subarray summing to k.
func longestSubarrayWithSum(nums []int, k int) int {
	firstSeen := map[int]int{0: -1} // prefix sum -> earliest index where it occurred
	sum, best := 0, 0
	for i, x := range nums {
		sum += x
		if j, ok := firstSeen[sum-k]; ok {
			best = max(best, i-j)
		}
		if _, ok := firstSeen[sum]; !ok {
			firstSeen[sum] = i
		}
	}
	return best
}

// applyRangeAdds starts from n zeros and applies each update {l, r, delta} to nums[l..r].
func applyRangeAdds(n int, updates [][3]int) []int {
	diff := make([]int, n+1)
	for _, u := range updates {
		l, r, delta := u[0], u[1], u[2]
		diff[l] += delta
		diff[r+1] -= delta
	}
	result := make([]int, n)
	running := 0
	for i := 0; i < n; i++ {
		running += diff[i]
		result[i] = running
	}
	return result
}

// prefix2D[i][j] is the sum of the rectangle grid[0..i-1][0..j-1].
type prefix2D [][]int

func newPrefix2D(grid [][]int) prefix2D {
	rows, cols := len(grid), len(grid[0])
	p := make(prefix2D, rows+1)
	for i := range p {
		p[i] = make([]int, cols+1)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			p[i+1][j+1] = grid[i][j] + p[i][j+1] + p[i+1][j] - p[i][j]
		}
	}
	return p
}

// sumRegion returns the sum of grid[r1..r2][c1..c2] (inclusive).
func (p prefix2D) sumRegion(r1, c1, r2, c2 int) int {
	return p[r2+1][c2+1] - p[r1][c2+1] - p[r2+1][c1] + p[r1][c1]
}

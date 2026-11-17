// Package dfs holds the runnable templates from dfs-study-plan.md.
package dfs

// dfs is the generic recursive template: base case is "already visited".
func dfs[T comparable](node T, getNeighbors func(T) []T, visited map[T]bool) {
	if visited[node] {
		return
	}
	visited[node] = true

	for _, neighbor := range getNeighbors(node) {
		dfs(neighbor, getNeighbors, visited)
	}
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// preorder and postorder are the tree specialization with "process node" filled in.
func preorder(node *TreeNode, out *[]int) {
	if node == nil {
		return
	}
	*out = append(*out, node.Val) // process before children
	preorder(node.Left, out)
	preorder(node.Right, out)
}

func postorder(node *TreeNode, out *[]int) {
	if node == nil {
		return
	}
	postorder(node.Left, out)
	postorder(node.Right, out)
	*out = append(*out, node.Val) // process after children
}

func maxDepth(node *TreeNode) int {
	if node == nil {
		return 0
	}

	left := maxDepth(node.Left)
	right := maxDepth(node.Right)

	return 1 + max(left, right)
}

func hasPathSum(node *TreeNode, remaining int) bool {
	if node == nil {
		return false
	}
	remaining -= node.Val
	if node.Left == nil && node.Right == nil {
		return remaining == 0
	}
	return hasPathSum(node.Left, remaining) || hasPathSum(node.Right, remaining)
}

// collectPaths is the Day 3 root-to-leaf template.
func collectPaths(node *TreeNode, path []int, result *[][]int) {
	if node == nil {
		return
	}

	path = append(path, node.Val)

	if node.Left == nil && node.Right == nil {
		tmp := make([]int, len(path))
		copy(tmp, path)
		*result = append(*result, tmp)
	}

	collectPaths(node.Left, path, result)
	collectPaths(node.Right, path, result)
}

// dfsGraph is the Day 4 adjacency-list template.
func dfsGraph(node int, graph [][]int, visited []bool) {
	visited[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] {
			dfsGraph(neighbor, graph, visited)
		}
	}
}

func countComponents(n int, graph [][]int) int {
	visited := make([]bool, n)
	count := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfsGraph(i, graph, visited)
			count++
		}
	}
	return count
}

var dirs = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

// fill is the Day 5 grid template. Callers must ensure target != replacement.
func fill(grid [][]int, r, c, target, replacement int) {
	if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
		return
	}
	if grid[r][c] != target {
		return
	}

	grid[r][c] = replacement

	for _, d := range dirs {
		fill(grid, r+d[0], c+d[1], target, replacement)
	}
}

// combine is the Day 6 backtracking template with isDone = "picked k" and
// isValid = "strictly larger than the last pick", which yields k-combinations of sorted candidates.
func combine(candidates []int, k int) [][]int {
	var result [][]int
	backtrack(candidates, k, nil, &result)
	return result
}

func backtrack(candidates []int, k int, current []int, result *[][]int) {
	if len(current) == k {
		tmp := make([]int, len(current))
		copy(tmp, current)
		*result = append(*result, tmp)
		return
	}

	for _, candidate := range candidates {
		if len(current) == 0 || candidate > current[len(current)-1] {
			current = append(current, candidate)
			backtrack(candidates, k, current, result)
			current = current[:len(current)-1]
		}
	}
}

// dfsIterative is the Day 7 explicit-stack template; it returns preorder.
func dfsIterative(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	stack := []*TreeNode{root}
	var order []int

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		order = append(order, node.Val)

		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
	return order
}

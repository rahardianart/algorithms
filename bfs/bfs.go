// Package bfs holds the runnable templates from bfs-study-plan.md.
package bfs

// bfs is the one template: it returns nodes in the order they were dequeued.
func bfs[T comparable](start T, getNeighbors func(T) []T) []T {
	queue := []T{start}
	visited := map[T]bool{start: true}
	var order []T

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, neighbor := range getNeighbors(node) {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// levelCounts is the Day 2 level-order template; it returns how many nodes sit on each level.
func levelCounts(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root}
	var counts []int

	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		counts = append(counts, levelSize)
	}
	return counts
}

// bfsGraph is the Day 3 adjacency-list template; it returns the visit order.
func bfsGraph(graph [][]int, start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	var order []int

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}

// bfsFrom runs BFS from start, marking all reachable nodes in the shared visited slice.
func bfsFrom(start int, graph [][]int, visited []bool) {
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

func countComponents(n int, graph [][]int) int {
	visited := make([]bool, n)
	count := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			bfsFrom(i, graph, visited)
			count++
		}
	}
	return count
}

type State struct {
	node int
	dist int
}

func shortestPath(graph [][]int, start, end int) int {
	visited := map[int]bool{start: true}
	queue := []State{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.node == end {
			return cur.dist
		}

		for _, neighbor := range graph[cur.node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, State{neighbor, cur.dist + 1})
			}
		}
	}
	return -1
}

type Cell struct{ r, c int }

// multiSourceBFS returns how many levels the wave needed to reach every reachable cell.
func multiSourceBFS(grid [][]int, sources []Cell) int {
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	queue := []Cell{}
	for _, s := range sources {
		queue = append(queue, s)
		visited[s.r][s.c] = true
	}

	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	minutes := 0

	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			cur := queue[0]
			queue = queue[1:]
			for _, d := range dirs {
				nr, nc := cur.r+d[0], cur.c+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
					visited[nr][nc] = true
					queue = append(queue, Cell{nr, nc})
				}
			}
		}
		if len(queue) > 0 {
			minutes++
		}
	}
	return minutes
}

// minSteps is the Day 6 state-space template; getNeighbors generates the next states.
func minSteps(start, end string, deadends map[string]bool, getNeighbors func(string) []string) int {
	if deadends[start] {
		return -1
	}
	visited := map[string]bool{start: true}
	queue := []string{start}
	steps := 0

	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			cur := queue[0]
			queue = queue[1:]

			if cur == end {
				return steps
			}

			for _, next := range getNeighbors(cur) {
				if !visited[next] && !deadends[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		}
		steps++
	}
	return -1
}

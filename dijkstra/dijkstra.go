// Package dijkstra holds the runnable templates from dijkstra-study-plan.md.
package dijkstra

import (
	"container/heap"
	"math"
)

const inf = math.MaxInt

type edge struct{ to, w int }

type item struct{ node, dist int }

type minHeap []item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// dijkstra returns the shortest distance from src to every node (inf if unreachable).
// Edge weights must be non-negative.
func dijkstra(n int, adj [][]edge, src int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0
	h := &minHeap{{src, 0}}

	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > dist[cur.node] {
			continue // stale entry: a shorter path was already found
		}
		for _, e := range adj[cur.node] {
			if nd := cur.dist + e.w; nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(h, item{e.to, nd})
			}
		}
	}
	return dist
}

// minPathCostGrid returns the cheapest cost from the top-left to the bottom-right cell,
// moving in 4 directions, where entering a cell (including the first) costs its value.
func minPathCostGrid(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	id := func(r, c int) int { return r*cols + c }

	dist := make([]int, rows*cols)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = grid[0][0]
	h := &minHeap{{0, grid[0][0]}}
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > dist[cur.node] {
			continue
		}
		r, c := cur.node/cols, cur.node%cols
		if r == rows-1 && c == cols-1 {
			return cur.dist
		}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				continue
			}
			if nd := cur.dist + grid[nr][nc]; nd < dist[id(nr, nc)] {
				dist[id(nr, nc)] = nd
				heap.Push(h, item{id(nr, nc), nd})
			}
		}
	}
	return -1
}

// zeroOneBFS is shortest paths when every edge weight is 0 or 1, using a deque instead of a heap.
func zeroOneBFS(n int, adj [][]edge, src int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0
	deque := []int{src}

	for len(deque) > 0 {
		u := deque[0]
		deque = deque[1:]
		for _, e := range adj[u] {
			if nd := dist[u] + e.w; nd < dist[e.to] {
				dist[e.to] = nd
				if e.w == 0 {
					deque = append([]int{e.to}, deque...) // same distance: front
				} else {
					deque = append(deque, e.to) // one more: back
				}
			}
		}
	}
	return dist
}

// cheapestWithOneFree returns the cheapest src→dst cost when at most one edge may be used for free.
// The state is (node, used), where used is 0 or 1.
func cheapestWithOneFree(n int, adj [][]edge, src, dst int) int {
	dist := make([]int, 2*n) // state id = node*2 + used
	for i := range dist {
		dist[i] = inf
	}
	dist[src*2] = 0
	h := &minHeap{{src * 2, 0}}

	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > dist[cur.node] {
			continue
		}
		u, used := cur.node/2, cur.node%2
		if u == dst {
			return cur.dist
		}
		for _, e := range adj[u] {
			if next := e.to*2 + used; cur.dist+e.w < dist[next] {
				dist[next] = cur.dist + e.w
				heap.Push(h, item{next, dist[next]})
			}
			if used == 0 {
				if next := e.to*2 + 1; cur.dist < dist[next] {
					dist[next] = cur.dist
					heap.Push(h, item{next, cur.dist})
				}
			}
		}
	}
	return -1
}

// bellmanFord returns the shortest distances from src using at most maxEdges edges.
// Weights may be negative. Each edge is {u, v, weight}.
func bellmanFord(n int, edges [][3]int, src, maxEdges int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0

	for round := 0; round < maxEdges; round++ {
		prev := make([]int, n)
		copy(prev, dist) // relax from last round's values only
		for _, e := range edges {
			u, v, w := e[0], e[1], e[2]
			if prev[u] != inf && prev[u]+w < dist[v] {
				dist[v] = prev[u] + w
			}
		}
	}
	return dist
}

// floydWarshall returns the shortest distance between every pair of nodes in an undirected graph.
func floydWarshall(n int, edges [][3]int) [][]int {
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			if i != j {
				dist[i][j] = inf
			}
		}
	}
	for _, e := range edges {
		dist[e[0]][e[1]] = min(dist[e[0]][e[1]], e[2])
		dist[e[1]][e[0]] = min(dist[e[1]][e[0]], e[2])
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k] != inf && dist[k][j] != inf && dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	return dist
}

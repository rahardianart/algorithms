// Package toposort holds the runnable templates from topological-sort-study-plan.md.
//
// An edge {u, v} means "u must come before v".
package toposort

// topoOrder is Kahn's algorithm. It returns nil if the graph has a cycle.
// Ready nodes are taken in index order, so the result is deterministic.
func topoOrder(n int, edges [][2]int) []int {
	adj := make([][]int, n)
	indegree := make([]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		indegree[e[1]]++
	}

	var queue []int
	for v := 0; v < n; v++ {
		if indegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	order := make([]int, 0, n)
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			indegree[v]--
			if indegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(order) < n {
		return nil
	}
	return order
}

const (
	white = iota // not visited
	gray         // on the current DFS path
	black        // finished: everything reachable from it is explored
)

// hasCycle detects a cycle in a directed graph with three-color DFS.
func hasCycle(n int, adj [][]int) bool {
	color := make([]int, n)

	var visit func(u int) bool
	visit = func(u int) bool {
		color[u] = gray
		for _, v := range adj[u] {
			if color[v] == gray {
				return true
			}
			if color[v] == white && visit(v) {
				return true
			}
		}
		color[u] = black
		return false
	}

	for u := 0; u < n; u++ {
		if color[u] == white && visit(u) {
			return true
		}
	}
	return false
}

// alienOrder returns a letter order consistent with the sorted word list, or "" if none exists.
// Only letters that appear in words are included; ties are broken alphabetically.
func alienOrder(words []string) string {
	present := [26]bool{}
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			present[w[i]-'a'] = true
		}
	}

	var edges [][2]int
	for i := 0; i+1 < len(words); i++ {
		a, b := words[i], words[i+1]
		j := 0
		for j < len(a) && j < len(b) && a[j] == b[j] {
			j++
		}
		if j == len(a) || j == len(b) {
			if len(a) > len(b) {
				return "" // "abc" before "ab" is impossible
			}
			continue // one is a prefix of the other: no information
		}
		edges = append(edges, [2]int{int(a[j] - 'a'), int(b[j] - 'a')})
	}

	order := topoOrder(26, edges)
	if order == nil {
		return ""
	}
	out := make([]byte, 0, 26)
	for _, c := range order {
		if present[c] {
			out = append(out, byte('a'+c))
		}
	}
	return string(out)
}

// longestPath returns the largest total weight of any path in a DAG. Each edge is {u, v, weight}.
func longestPath(n int, edges [][3]int) int {
	pairs := make([][2]int, len(edges))
	adj := make([][][2]int, n)
	for i, e := range edges {
		pairs[i] = [2]int{e[0], e[1]}
		adj[e[0]] = append(adj[e[0]], [2]int{e[1], e[2]})
	}

	dist := make([]int, n) // best path weight ending at each node
	best := 0
	for _, u := range topoOrder(n, pairs) {
		best = max(best, dist[u])
		for _, next := range adj[u] {
			dist[next[0]] = max(dist[next[0]], dist[u]+next[1])
		}
	}
	return best
}

// peelLeaves removes the leaves of an undirected tree layer by layer and returns the layers.
// The last layer holds the tree's center(s).
func peelLeaves(n int, edges [][2]int) [][]int {
	if n == 1 {
		return [][]int{{0}}
	}
	adj := make([][]int, n)
	degree := make([]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
		degree[e[0]]++
		degree[e[1]]++
	}

	var layer []int
	for v := 0; v < n; v++ {
		if degree[v] == 1 {
			layer = append(layer, v)
		}
	}

	var layers [][]int
	remaining := n
	for len(layer) > 0 {
		layers = append(layers, layer)
		remaining -= len(layer)
		if remaining == 0 {
			break
		}
		var next []int
		for _, leaf := range layer {
			for _, nb := range adj[leaf] {
				degree[nb]--
				if degree[nb] == 1 {
					next = append(next, nb)
				}
			}
		}
		layer = next
	}
	return layers
}

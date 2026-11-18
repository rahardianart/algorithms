// Package unionfind holds the runnable templates from union-find-study-plan.md.
package unionfind

import "sort"

type DSU struct {
	parent []int
	size   []int
	count  int // number of disjoint sets
}

func newDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n), count: n}
	for i := range d.parent {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

// find returns x's root and points every node on the way directly at it.
func (d *DSU) find(x int) int {
	root := x
	for d.parent[root] != root {
		root = d.parent[root]
	}
	for d.parent[x] != root {
		d.parent[x], x = root, d.parent[x]
	}
	return root
}

// union merges the sets of a and b; it returns false if they were already together.
func (d *DSU) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.count--
	return true
}

func countComponents(n int, edges [][2]int) int {
	d := newDSU(n)
	for _, e := range edges {
		d.union(e[0], e[1])
	}
	return d.count
}

// isTree reports whether the undirected graph is connected and has no cycle.
func isTree(n int, edges [][2]int) bool {
	if len(edges) != n-1 {
		return false
	}
	d := newDSU(n)
	for _, e := range edges {
		if !d.union(e[0], e[1]) {
			return false
		}
	}
	return true
}

// groupEquivalent groups strings connected by "same group" pairs. Groups and their members
// are sorted so the output is deterministic.
func groupEquivalent(pairs [][2]string) [][]string {
	parent := map[string]string{}
	var find func(x string) string
	find = func(x string) string {
		if _, ok := parent[x]; !ok {
			parent[x] = x
		}
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	for _, p := range pairs {
		ra, rb := find(p[0]), find(p[1])
		if ra != rb {
			parent[rb] = ra
		}
	}

	byRoot := map[string][]string{}
	for x := range parent {
		r := find(x)
		byRoot[r] = append(byRoot[r], x)
	}
	groups := make([][]string, 0, len(byRoot))
	for _, g := range byRoot {
		sort.Strings(g)
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	return groups
}

// kruskal returns the total weight of a minimum spanning tree, or -1 if the graph is disconnected.
// Each edge is {u, v, weight}.
func kruskal(n int, edges [][3]int) int {
	sorted := make([][3]int, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i][2] < sorted[j][2] })

	d := newDSU(n)
	total := 0
	for _, e := range sorted {
		if d.union(e[0], e[1]) {
			total += e[2]
		}
	}
	if d.count != 1 {
		return -1
	}
	return total
}

// potentialDSU tracks constraints of the form value[a] - value[b] = diff.
type potentialDSU struct {
	parent []int
	offset []int // value[x] - value[parent[x]]
}

func newPotentialDSU(n int) *potentialDSU {
	p := &potentialDSU{parent: make([]int, n), offset: make([]int, n)}
	for i := range p.parent {
		p.parent[i] = i
	}
	return p
}

// find returns x's root and value[x] - value[root], compressing the path.
func (p *potentialDSU) find(x int) (int, int) {
	if p.parent[x] == x {
		return x, 0
	}
	root, parentOffset := p.find(p.parent[x])
	p.offset[x] += parentOffset
	p.parent[x] = root
	return root, p.offset[x]
}

// union records value[a] - value[b] = diff. It returns false if that contradicts earlier facts.
func (p *potentialDSU) union(a, b, diff int) bool {
	ra, oa := p.find(a) // value[a] = value[ra] + oa
	rb, ob := p.find(b) // value[b] = value[rb] + ob
	if ra == rb {
		return oa-ob == diff
	}
	// attach ra under rb: value[ra] - value[rb] = diff + ob - oa
	p.parent[ra] = rb
	p.offset[ra] = diff + ob - oa
	return true
}

// diff returns value[a] - value[b], if a and b are connected.
func (p *potentialDSU) diff(a, b int) (int, bool) {
	ra, oa := p.find(a)
	rb, ob := p.find(b)
	if ra != rb {
		return 0, false
	}
	return oa - ob, true
}

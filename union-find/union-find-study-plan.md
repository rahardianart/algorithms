# Union-Find Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 16 — an alternative to BFS/DFS for connectivity)  
**Goal:** Write a DSU with path compression and union by size from memory, use `union`'s return value to detect cycles, and model problems as "which things belong to the same group?".

---

## The Core Template

```go
type DSU struct {
    parent []int
    size   []int
    count  int // number of disjoint sets
}

func newDSU(n int) *DSU {
    d := &DSU{parent: make([]int, n), size: make([]int, n), count: n}
    for i := range d.parent {
        d.parent[i] = i // everyone starts as their own root
        d.size[i] = 1
    }
    return d
}

// find returns x's root, pointing every node on the way directly at it
func (d *DSU) find(x int) int {
    root := x
    for d.parent[root] != root {
        root = d.parent[root]
    }
    for d.parent[x] != root { // path compression
        d.parent[x], x = root, d.parent[x]
    }
    return root
}

// union merges a's and b's sets; false means they were already together
func (d *DSU) union(a, b int) bool {
    ra, rb := d.find(a), d.find(b)
    if ra == rb {
        return false
    }
    if d.size[ra] < d.size[rb] { // union by size: small tree goes under big tree
        ra, rb = rb, ra
    }
    d.parent[rb] = ra
    d.size[ra] += d.size[rb]
    d.count--
    return true
}
```

**Three rules — never break them:**

1. **Always compare roots, never raw nodes.** "Are `a` and `b` connected?" is `find(a) == find(b)`. Comparing `parent[a] == parent[b]` is wrong as soon as trees are more than one level deep.
2. **Use both optimizations.** Path compression flattens trees during `find`; union by size keeps them shallow during `union`. Together they make every operation effectively O(1) (amortized O(α(n)), where α is the inverse Ackermann function — ≤ 4 for any realistic n). Either one alone still leaves O(log n).
3. **`union` returning `false` is information.** It means "these were already connected", which in an undirected graph means **this edge closes a cycle**. Many problems are answered by that single boolean.

**Brute force → union-find:** "are `a` and `b` connected?" after each new edge could be answered with a fresh BFS/DFS: O(V + E) per query. Union-find handles edges **as they arrive** and answers each query in near-O(1), never revisiting the whole graph. It can't *remove* edges, though — when you need deletions or actual paths, go back to BFS/DFS.

**Complexity:** O(α(n)) amortized per `find`/`union` — treat it as O(1). O(n) space.

**Run it:** every template in this file is a real function in [`unionfind.go`](unionfind.go), tested in [`unionfind_test.go`](unionfind_test.go). From the repo root: `go test ./union-find/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Core template  
**Key insight:** Each set is a tree, and the root is the set's name. `find` asks "what's your set's name?", `union` makes one root point at the other.

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each line.
2. Close this file. Write `newDSU`, `find`, and `union` from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does union by size keep trees shallow?*

   > A node's depth only increases when its tree is attached under another root — and with union by size, that only happens when its tree is the smaller one, so the combined tree is at least twice the size of the node's old tree. Sizes can double at most log₂ n times, so no node is ever deeper than log₂ n. Path compression then flattens those paths further on every `find`.

**Dry run** — 6 nodes, a sequence of unions (`parent` after each step):

| Operation | Bigger root | Smaller root | Action | `parent` after | `count` |
|---|---|---|---|---|---|
| union(0, 1) | 0 | 1 | attach 1 under 0 | `[0, 0, 2, 3, 4, 5]` | 5 |
| union(1, 2) | 0 | 2 | attach 2 under 0 | `[0, 0, 0, 3, 4, 5]` | 4 |
| union(3, 4) | 3 | 4 | attach 4 under 3 | `[0, 0, 0, 3, 3, 5]` | 3 |
| union(2, 0) | 0 | 0 | same root → **false** (cycle) | `[0, 0, 0, 3, 3, 5]` | 3 |
| union(4, 2) | 0 | 3 | attach 3 under 0 | `[0, 0, 0, 0, 3, 5]` | 2 |

In the last step, `find(4)` returned 3 and `find(2)` returned 0. The set `{0,1,2}` (size 3) is bigger than `{3,4}` (size 2), so root 3 went under root 0. Node 4 is now two levels deep — the next `find(4)` compresses it to point at 0 directly.

---

## Day 2 — Connected Components and Redundant Edges

**Pattern:** Union every edge; `count` is the number of components; a failed union is a redundant edge  
**Key insight:** You don't need an adjacency list at all — union-find consumes the edge list directly.

**Template adaptation:**

```go
func countComponents(n int, edges [][2]int) int {
    d := newDSU(n)
    for _, e := range edges {
        d.union(e[0], e[1])
    }
    return d.count
}
```

**Why `count` is right:** it starts at `n` (every node alone) and drops by exactly one per *successful* union — each one merges two different components into one. Failed unions (already connected) don't change it.

**Run it:** `newDSU`, `find`, `union`, and `countComponents` in [`unionfind.go`](unionfind.go) (dry run on Day 1).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Provinces | Medium | https://leetcode.com/problems/number-of-provinces/ |
| 2 | Redundant Connection | Medium | https://leetcode.com/problems/redundant-connection/ |

**What to notice — Problem 1:** You solved this with BFS (week 1) and DFS (week 2). Now with union-find: union `i` and `j` for every `isConnected[i][j] == 1` with `i < j`, return `count`. Compare the three solutions — union-find has no recursion and no queue.

**What to notice — Problem 2:** Union edges in order; the first `union` that returns `false` is the edge that closed the cycle — return it. Nodes are 1-indexed, so size the DSU `n + 1`.

**Complexity:** both O(E · α(n)) ≈ O(E). Provinces reads the full matrix: O(n²).

**Common mistakes:**

- Counting components by `parent[i] == i` *without* the unions being complete — fine at the end, wrong mid-way. Prefer `count`.
- Redundant Connection: 0-indexing a 1-indexed input (index out of range on node `n`).

---

## Day 3 — Is It a Tree? (Cycle Detection)

**Pattern:** A graph is a tree iff it's connected **and** acyclic — check one with a count and the other with `union`'s return value  
**Key insight:** For `n` nodes, a tree has exactly `n − 1` edges. With that many edges, "no cycle" and "connected" imply each other — so checking edge count + no cycle is enough.

**Template adaptation:**

```go
func isTree(n int, edges [][2]int) bool {
    if len(edges) != n-1 {
        return false
    }
    d := newDSU(n)
    for _, e := range edges {
        if !d.union(e[0], e[1]) {
            return false // closes a cycle
        }
    }
    return true // n-1 successful unions → exactly one component
}
```

**Why `n − 1` edges + no cycle means connected:** every successful union reduces the component count by one. Starting from `n` components, `n − 1` successful unions leave exactly 1.

**Dry run** — `n = 5`:

| Edges | Edge count | Unions | Result |
|---|---|---|---|
| `[0,1] [0,2] [0,3] [1,4]` | 4 = n − 1 | all succeed | **tree** |
| `[0,1] [1,2] [2,3] [1,3] [1,4]` | 5 ≠ 4 | — | **not a tree** (early exit) |
| `n = 4`: `[0,1] [1,2] [2,0]` | 3 = n − 1 | `[2,0]` fails | **not a tree** (cycle, and node 3 is alone) |

The third row is the sneaky case: the edge count is right, but a cycle "used up" an edge that should have connected node 3.

**Run it:** `isTree` in [`unionfind.go`](unionfind.go) — the test includes all three cases.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Graph Valid Tree | Medium | https://leetcode.com/problems/graph-valid-tree/ (LeetCode Premium — use lintcode.com/problem/178 if needed) |
| 2 | Number of Operations to Make Network Connected | Medium | https://leetcode.com/problems/number-of-operations-to-make-network-connected/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Two numbers decide it: redundant cables (failed unions) and components (`count`). You need `count − 1` moves, and you need at least that many cables in total: if `len(connections) < n − 1`, return −1.

**Complexity:** both O(E · α(n)).

**Common mistakes:**

- Graph Valid Tree: checking only "no cycle" (misses disconnected forests) or only "connected" (misses cycles).
- Network: counting redundant cables but forgetting the `len(connections) < n − 1` impossibility check.

---

## Day 4 — Grouping Things That Aren't Numbered

**Pattern:** Union-find over strings (or any comparable key) with a map instead of slices  
**Key insight:** The structure doesn't care what the elements are — swap `[]int` for `map[string]string` and create entries lazily on first `find`.

**Template adaptation:**

```go
// group strings connected by "same group" pairs
func groupEquivalent(pairs [][2]string) [][]string {
    parent := map[string]string{}
    var find func(x string) string
    find = func(x string) string {
        if _, ok := parent[x]; !ok {
            parent[x] = x // first time we see x: its own root
        }
        if parent[x] != x {
            parent[x] = find(parent[x]) // recursive path compression
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
    // … then sort members and groups so the output is deterministic
}
```

**Dry run** — pairs `(a,b) (c,d) (b,e) (x,x)`:

| Pair | Roots before | Action |
|---|---|---|
| (a, b) | a, b | b → a |
| (c, d) | c, d | d → c |
| (b, e) | a, e | e → a |
| (x, x) | x, x | same root — just registers `x` |

Groups: `{a, b, e}`, `{c, d}`, `{x}`.

**Run it:** `groupEquivalent` in [`unionfind.go`](unionfind.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Accounts Merge | Medium | https://leetcode.com/problems/accounts-merge/ |
| 2 | Satisfiability of Equality Equations | Medium | https://leetcode.com/problems/satisfiability-of-equality-equations/ |

**What to notice — Problem 1:** Union every email in an account with that account's first email. Then group emails by root, sort each group, and prepend the owner's name (keep an `email → name` map).

**What to notice — Problem 2:** **Two passes.** First union every `a==b`. Only then check every `a!=b`: if `find(a) == find(b)`, it's a contradiction. Checking `!=` during the first pass misses equalities that come later in the input.

**Complexity:** Accounts Merge O(E log E) for the final sorts. Equations O(n) (26 variables).

**Common mistakes:**

- Iterating a Go map to build output and expecting a fixed order — sort.
- Equations: handling `==` and `!=` in a single pass.
- Accounts Merge: unioning accounts by *name* — different people can share a name.

---

## Day 5 — Kruskal's Minimum Spanning Tree

**Pattern:** Sort edges by weight; take an edge only if it connects two different components  
**Key insight:** The cheapest edge that doesn't create a cycle is always safe to add to a minimum spanning tree. Union-find answers "does it create a cycle?" in O(1).

**Template adaptation:**

```go
// total weight of a minimum spanning tree, or -1 if disconnected; edges are {u, v, weight}
func kruskal(n int, edges [][3]int) int {
    sorted := make([][3]int, len(edges))
    copy(sorted, edges)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i][2] < sorted[j][2] })

    d := newDSU(n)
    total := 0
    for _, e := range sorted {
        if d.union(e[0], e[1]) { // different components → this edge is needed
            total += e[2]
        }
    }
    if d.count != 1 {
        return -1
    }
    return total
}
```

**Why the cheapest non-cycle edge is safe (cut property):** take the two components the edge connects. *Any* spanning tree must cross between them with some edge. If an optimal tree used a more expensive crossing edge, swapping it for this cheaper one keeps it a spanning tree and lowers the cost — so some optimal tree uses the cheap one. (The Greedy week's exchange argument, on graphs.)

**Dry run** — 5 nodes, edges sorted by weight:

| Edge | Weight | Decision | Total |
|---|---|---|---|
| 0–2 | 1 | take | 1 |
| 1–2 | 2 | take | 3 |
| 3–4 | 3 | take | 6 |
| 0–1 | 4 | skip (0 and 1 already connected via 2) | 6 |
| 1–3 | 5 | take | 11 |
| 2–3 | 8 | skip (cycle) | 11 |

MST weight **11**, using exactly n − 1 = 4 edges.

**Run it:** `kruskal` in [`unionfind.go`](unionfind.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Min Cost to Connect All Points | Medium | https://leetcode.com/problems/min-cost-to-connect-all-points/ |
| 2 | Swim in Rising Water | Hard | https://leetcode.com/problems/swim-in-rising-water/ |

**What to notice — Problem 1:** Build all n² edges with Manhattan distances, then Kruskal. (Prim's algorithm with a heap — or the O(n²) array version — avoids building every edge; worth trying after.)

**What to notice — Problem 2:** "Kruskal by time": process cells in increasing elevation; when a cell opens, union it with already-open neighbors. The answer is the elevation at which the top-left and bottom-right first share a root. You'll solve this again with Dijkstra in week 18.

**Complexity:** Kruskal O(E log E) for the sort. Problem 1 has E = n², so O(n² log n). Problem 2 O(n² log n).

**Common mistakes:**

- Sorting the caller's slice in place when they didn't expect it (the template copies first).
- Forgetting the "disconnected → no spanning tree" case.
- Swim: unioning a cell with neighbors that aren't open yet.

---

## Day 6 — Union-Find with Weights

**Pattern:** Store, for each node, how it relates to its parent (an offset or ratio); combine those along the path during `find`  
**Key insight:** If you know `a − parent(a)` and `parent(a) − root`, you know `a − root` — so path compression can keep the relation exact while flattening the tree.

**Template adaptation:**

```go
// tracks constraints value[a] - value[b] = diff
type potentialDSU struct {
    parent []int
    offset []int // value[x] - value[parent[x]]
}

// returns x's root and value[x] - value[root]
func (p *potentialDSU) find(x int) (int, int) {
    if p.parent[x] == x {
        return x, 0
    }
    root, parentOffset := p.find(p.parent[x])
    p.offset[x] += parentOffset // x→parent plus parent→root
    p.parent[x] = root
    return root, p.offset[x]
}

// record value[a] - value[b] = diff; false if it contradicts earlier facts
func (p *potentialDSU) union(a, b, diff int) bool {
    ra, oa := p.find(a) // value[a] = value[ra] + oa
    rb, ob := p.find(b) // value[b] = value[rb] + ob
    if ra == rb {
        return oa-ob == diff
    }
    p.parent[ra] = rb
    p.offset[ra] = diff + ob - oa // value[ra] - value[rb]
    return true
}
```

**Where `diff + ob − oa` comes from:** we want `value[a] − value[b] = diff`. Substitute `value[a] = value[ra] + oa` and `value[b] = value[rb] + ob`: `value[ra] + oa − value[rb] − ob = diff`, so `value[ra] − value[rb] = diff + ob − oa`. That's the new offset of `ra` under `rb`.

**Dry run** — variables `a=0, b=1, c=2`:

| Fact / query | Result |
|---|---|
| `a − b = 3` | accepted |
| `b − c = 2` | accepted |
| `diff(a, c)` | **5** (3 + 2, read off the offsets after `find`) |
| `diff(c, a)` | **−5** |
| `c − a = 1` | **rejected** — contradicts −5 |
| `c − a = −5` | accepted — consistent |

**Run it:** `potentialDSU` in [`unionfind.go`](unionfind.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Evaluate Division | Medium | https://leetcode.com/problems/evaluate-division/ |
| 2 | Most Stones Removed with Same Row or Column | Medium | https://leetcode.com/problems/most-stones-removed-with-same-row-or-column/ |

**What to notice — Problem 1:** The same structure with **ratios** instead of differences: `a / b = k`. Replace `+` with `×` and `−` with `÷` everywhere (`offset` becomes a `float64` "weight", identity is 1 instead of 0). Variables are strings, so use maps.

**What to notice — Problem 2:** Modelling, not weights: the answer is `stones − components`, where stones in the same row or column are connected. Trick: union the stone's **row** with its **column** (offset columns by 10001 so they don't collide with rows) — then you never compare stones pairwise.

**Complexity:** both near-linear: O((equations + queries) · α).

**Common mistakes:**

- Updating `offset[x]` without first compressing the parent (the recursion order in `find` matters — parent first, then `x`).
- Evaluate Division: returning 1.0 for `x / x` when `x` never appeared — the answer is −1.0.
- Stones: unioning stone indices pairwise (O(n²)) instead of rows with columns.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Reach for union-find when the problem only ever *adds* connections and asks "same group?" / "how many groups?" / "does this connection create a cycle?". Reach for BFS/DFS when you need paths, distances, or removals.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Pick one problem from this week and solve it with BFS or DFS too. Which version is shorter? Which handles edges arriving one at a time?

**Self-assessment — can you answer these without hesitation?**

- [ ] Write `find` (with path compression) and `union` (by size) from memory
- [ ] What does `union` returning `false` mean in an undirected graph?
- [ ] Why does `n − 1` edges plus no cycle imply connected?
- [ ] Why must Satisfiability of Equality Equations use two passes?
- [ ] State the cut property and how Kruskal uses union-find
- [ ] In a weighted DSU, how is `offset[x]` updated during path compression?

---

## Pattern Cheatsheet

| Pattern | What you union | What you read | Time |
|---------|----------------|---------------|------|
| Components | every edge | `count` | O(E · α) |
| Redundant / cycle edge | edges in order | first `union` returning `false` | O(E · α) |
| Valid tree | every edge | `len(edges) == n-1` and no failed union | O(E · α) |
| Grouping by key | pairs of keys (map-based DSU) | members by root | O(n · α) + sorting |
| Kruskal MST | edges sorted by weight | sum of weights of successful unions | O(E log E) |
| Weighted DSU | constraints with a difference / ratio | offset between two nodes with the same root | O(n · α) |

---

# Topological Sort Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 17 — BFS week's queue + DFS week's recursion, on directed graphs)  
**Goal:** Write Kahn's algorithm and three-color DFS from memory, detect cycles in directed graphs, build the dependency graph from a problem statement, and run DP in topological order.

---

## The Two Templates

### Kahn's Algorithm (BFS on in-degrees)

```go
// edge {u, v} means "u must come before v"
adj := make([][]int, n)
indegree := make([]int, n)
for _, e := range edges {
    adj[e[0]] = append(adj[e[0]], e[1])
    indegree[e[1]]++
}

var queue []int
for v := 0; v < n; v++ {
    if indegree[v] == 0 { // nothing has to come before v
        queue = append(queue, v)
    }
}

order := make([]int, 0, n)
for len(queue) > 0 {
    u := queue[0]
    queue = queue[1:]
    order = append(order, u)
    for _, v := range adj[u] {
        indegree[v]-- // u is done, so one of v's requirements is satisfied
        if indegree[v] == 0 {
            queue = append(queue, v)
        }
    }
}

if len(order) < n {
    // a cycle: some nodes never reached in-degree 0
}
```

### Three-Color DFS (cycle detection)

```go
const (
    white = iota // not visited
    gray         // on the current DFS path
    black        // finished
)

var visit func(u int) bool // returns true if a cycle is found
visit = func(u int) bool {
    color[u] = gray
    for _, v := range adj[u] {
        if color[v] == gray {
            return true // back edge to the current path: cycle
        }
        if color[v] == white && visit(v) {
            return true
        }
    }
    color[u] = black
    return false
}
```

**Three rules — never break them:**

1. **Get the edge direction right before anything else.** An edge means "this before that". LeetCode's Course Schedule gives `[a, b]` meaning *b* before *a* — so the edge is `b → a`. Half of all topological-sort bugs are a reversed edge.
2. **A topological order exists iff the graph has no directed cycle.** Kahn's tells you by `len(order) < n`; DFS tells you by hitting a **gray** node. Always check — don't assume the input is a DAG.
3. **Directed cycles need three colors, not a `visited` bool.** Reaching a node that's already *finished* (black) through a different path is fine (a diamond); only reaching a node that's still *on the current path* (gray) is a cycle. A plain visited set can't tell those apart.

**Brute force → topological sort:** trying every ordering of n tasks is n!. Kahn's observes that some task with no remaining prerequisites can always go next (if none exists, there's a cycle), so it builds a valid order in a single pass over nodes and edges.

**Complexity:** O(V + E) time and space for both templates.

**Run it:** every template in this file is a real function in [`toposort.go`](toposort.go), tested in [`toposort_test.go`](toposort_test.go). From the repo root: `go test ./topological-sort/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** In-degree = "number of things that must happen before me that haven't happened yet". When it hits 0, I'm ready.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both templates from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does `len(order) < n` mean there's a cycle?*

   > Nodes on a cycle each wait for the one before them on the cycle. None of them can reach in-degree 0 until another one is output first — which never happens. So they're never enqueued, never added to `order`, and the order comes up short by at least the length of the cycle (plus anything that depends on it).

**Dry run** — Kahn's on 6 nodes with edges `5→2, 5→0, 4→0, 4→1, 2→3, 3→1` (starting in-degrees `[2, 2, 1, 1, 0, 0]`):

| Dequeue | In-degree decrements | Newly ready | Queue after | Order |
|---|---|---|---|---|
| 4 | 0→1, 1→1 | — | `[5]` | `[4]` |
| 5 | 2→0, 0→0 | 2, 0 | `[2, 0]` | `[4, 5]` |
| 2 | 3→0 | 3 | `[0, 3]` | `[4, 5, 2]` |
| 0 | — | — | `[3]` | `[4, 5, 2, 0]` |
| 3 | 1→0 | 1 | `[1]` | `[4, 5, 2, 0, 3]` |
| 1 | — | — | `[]` | `[4, 5, 2, 0, 3, 1]` |

All 6 nodes are output → no cycle. Node `1` waited longest: it needed both `4` and `3`, and `3` needed `2`, which needed `5`.

---

## Day 2 — Kahn's Algorithm: Can It Be Done, and In What Order?

**Pattern:** Build adjacency + in-degrees; process ready nodes; check the count  
**Key insight:** The same run answers both "is it possible?" (`len(order) == n`) and "give me an order" (`order`).

**Template adaptation:**

```go
func topoOrder(n int, edges [][2]int) []int {
    // … the Kahn's template above …
    if len(order) < n {
        return nil // cycle
    }
    return order
}
```

**Why any order Kahn's produces is valid:** a node is only output when its in-degree is 0, meaning every node with an edge into it was already output. So every edge `u → v` has `u` before `v`.

**Run it:** `topoOrder` in [`toposort.go`](toposort.go) (dry run on Day 1) — the test also checks the order respects every edge.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Course Schedule | Medium | https://leetcode.com/problems/course-schedule/ |
| 2 | Course Schedule II | Medium | https://leetcode.com/problems/course-schedule-ii/ |

**What to notice — Problem 1:** Return `len(order) == numCourses`. Mind rule 1: `prerequisites[i] = [a, b]` means the edge `b → a`.

**What to notice — Problem 2:** Return `order` itself, or an empty slice on a cycle. Any valid order is accepted.

**Complexity:** both O(V + E).

**Common mistakes:**

- Reversing the edge direction (rule 1). Your order comes out backwards, and Course Schedule II fails.
- Forgetting nodes with no edges at all — they start with in-degree 0 and must be in the initial queue.
- Building `adj` with a map and iterating it: Go map order is random, so your output changes between runs. Slices of slices are deterministic.

---

## Day 3 — Three-Color DFS

**Pattern:** DFS where "on the current path" (gray) differs from "finished" (black)  
**Key insight:** Black nodes are fully explored and known to be safe (no cycle reachable through them). Many problems are really asking "which nodes end up black?".

**Template adaptation:**

```go
func hasCycle(n int, adj [][]int) bool {
    color := make([]int, n)
    // … visit from the template above …
    for u := 0; u < n; u++ {
        if color[u] == white && visit(u) {
            return true
        }
    }
    return false
}
```

**Why a plain `visited` set fails — the diamond:** edges `0→1, 0→2, 1→3, 2→3`. DFS from 0 visits 1, then 3 (finishes it), then backs up and visits 2, which points at 3 again. With a visited bool, "3 is already visited" looks like a cycle. With colors, 3 is **black** (finished), so it's correctly ignored. Only a **gray** neighbor — one still on the current path — closes a loop.

**Dry run** — edges `0→1, 1→2, 2→0, 3→1`:

| Step | Visit | Colors after (0, 1, 2, 3) | Note |
|---|---|---|---|
| 1 | 0 | G, W, W, W | |
| 2 | 1 | G, G, W, W | |
| 3 | 2 | G, G, G, W | |
| 4 | 2 → 0 | — | 0 is **gray** → cycle `0 → 1 → 2 → 0` |

**Run it:** `hasCycle` in [`toposort.go`](toposort.go) — the test includes the diamond.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Find Eventual Safe States | Medium | https://leetcode.com/problems/find-eventual-safe-states/ |
| 2 | Course Schedule IV | Medium | https://leetcode.com/problems/course-schedule-iv/ |

**What to notice — Problem 1:** A node is safe iff no path from it reaches a cycle — exactly the nodes that finish **black** in the template. Don't reset colors between starting nodes; a node's color is its final verdict. (Alternative: Kahn's on the *reversed* graph — the nodes it outputs are the safe ones.)

**What to notice — Problem 2:** Reachability queries. Process nodes in topological order and propagate: if `u → v`, then everything that is a prerequisite of `u` is also a prerequisite of `v`. With n ≤ 100, a `[100][100]bool` (or a bitset per node) works.

**Complexity:** Problem 1 O(V + E). Problem 2 O(V · (V + E)) for the propagation.

**Common mistakes:**

- Using a two-state visited set for directed cycle detection (the diamond).
- Marking a node black **before** exploring its neighbors.
- Problem 1: resetting colors for each start node — turns O(V + E) into O(V · (V + E)).

---

## Day 4 — Building the Graph from Constraints

**Pattern:** The input isn't a graph; turn each comparison or dependency into an edge, then topo-sort  
**Key insight:** In a sorted list, each **adjacent** pair gives at most one fact: the first position where they differ tells you one letter comes before another. Non-adjacent pairs add nothing new.

**Template adaptation:**

```go
// letter order consistent with a sorted list of words in an alien language
func alienOrder(words []string) string {
    // … mark which letters appear …
    var edges [][2]int
    for i := 0; i+1 < len(words); i++ {
        a, b := words[i], words[i+1]
        j := 0
        for j < len(a) && j < len(b) && a[j] == b[j] {
            j++
        }
        if j == len(a) || j == len(b) {
            if len(a) > len(b) {
                return "" // "abc" listed before "ab" is impossible
            }
            continue // one is a prefix of the other: no information
        }
        edges = append(edges, [2]int{int(a[j] - 'a'), int(b[j] - 'a')})
    }

    order := topoOrder(26, edges)
    if order == nil {
        return "" // contradictory facts: a cycle
    }
    // … keep only letters that appear …
}
```

**Dry run** — `["wrt", "wrf", "er", "ett", "rftt"]`:

| Adjacent pair | First difference | Edge |
|---|---|---|
| wrt / wrf | t vs f | t → f |
| wrf / er | w vs e | w → e |
| er / ett | r vs t | r → t |
| ett / rftt | e vs r | e → r |

Edges chain into `w → e → r → t → f`, so the order is **`wertf`**.

**Run it:** `alienOrder` in [`toposort.go`](toposort.go) — the test includes the prefix and contradiction cases.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Alien Dictionary | Hard | https://leetcode.com/problems/alien-dictionary/ (LeetCode Premium — use lintcode.com/problem/892 if needed) |
| 2 | Find All Possible Recipes from Given Supplies | Medium | https://leetcode.com/problems/find-all-possible-recipes-from-given-supplies/ |

**What to notice — Problem 1:** This IS the template. LintCode's version wants the smallest order in normal lexicographic order when several are valid — the template's index-ordered queue already does that; a min-heap is the general way.

**What to notice — Problem 2:** Nodes are strings: ingredients and recipes. Edge `ingredient → recipe`; a recipe's in-degree is its ingredient count. Seed the queue with the supplies (not with in-degree-0 nodes) — anything missing from the supplies never gets "cooked", which is exactly right.

**Complexity:** Alien Dictionary O(total characters + 26). Recipes O(total ingredients listed).

**Common mistakes:**

- Comparing every pair of words instead of adjacent ones (O(n²) and no extra information).
- Missing the "longer word is a prefix-extension listed first" invalid case.
- Recipes: starting the queue with in-degree-0 nodes instead of the supplies.

---

## Day 5 — DP in Topological Order

**Pattern:** Process nodes in topological order; push each node's best value forward along its edges  
**Key insight:** In a DAG, topological order guarantees that when you process a node, **every** node that can reach it is already final. That's exactly the condition DP needs.

**Template adaptation:**

```go
// largest total weight of any path in a DAG; edges are {u, v, weight}
func longestPath(n int, edges [][3]int) int {
    // … build adj and a pairs list for topoOrder …
    dist := make([]int, n) // best path weight ending at each node
    best := 0
    for _, u := range topoOrder(n, pairs) {
        best = max(best, dist[u]) // dist[u] is final: all predecessors are done
        for _, next := range adj[u] {
            dist[next[0]] = max(dist[next[0]], dist[u]+next[1])
        }
    }
    return best
}
```

**Why not BFS or Dijkstra?** Longest path is NP-hard in general graphs — but in a DAG there are no cycles to loop around, and topological order lets you relax each edge exactly once. (The same idea gives *shortest* paths in a DAG in O(V + E), even with negative weights.)

**Dry run** — edges `0→1 (3)`, `0→2 (2)`, `1→3 (4)`, `2→3 (1)`, `3→4 (2)`, `2→4 (9)`; topological order `0, 1, 2, 3, 4`:

| Process | `dist[u]` (final) | Relaxations |
|---|---|---|
| 0 | 0 | dist[1] = 3; dist[2] = 2 |
| 1 | 3 | dist[3] = 7 |
| 2 | 2 | dist[3] stays 7; dist[4] = 11 |
| 3 | 7 | dist[4] stays 11 (7 + 2 = 9) |
| 4 | 11 | — |

Longest path = **11** (`0 → 2 → 4`), even though `0 → 1 → 3 → 4` has more edges.

**Run it:** `longestPath` in [`toposort.go`](toposort.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Parallel Courses III | Hard | https://leetcode.com/problems/parallel-courses-iii/ |
| 2 | Longest Increasing Path in a Matrix | Hard | https://leetcode.com/problems/longest-increasing-path-in-a-matrix/ |

**What to notice — Problem 1:** Weights are on **nodes** (course durations). `finish[v] = time[v] + max(finish[u] for every prerequisite u)`. Same template, relax with `finish[u]` and add `time[v]` when `v` becomes ready. The answer is `max(finish)`.

**What to notice — Problem 2:** The graph is implicit: an edge from each cell to every strictly larger neighbor, which can never form a cycle. Either Kahn's on that graph (peel cells with no smaller neighbor, count layers) or DFS with memoization (DFS week + DP week). Both are O(R · C).

**Complexity:** both O(V + E).

**Common mistakes:**

- Running the DP in input order instead of topological order — values aren't final when you read them.
- Parallel Courses III: adding the node's own time at the wrong moment (it belongs to the node, not the edge).
- Longest Increasing Path: plain DFS without memo — exponential.

---

## Day 6 — Peeling Layers

**Pattern:** Kahn's idea on other structures: repeatedly remove everything at "degree 1" (leaves) or "in-degree 0", layer by layer  
**Key insight:** What's left after peeling is the "core": the center of a tree, or the nodes on cycles in a directed graph.

**Template adaptation:**

```go
// remove leaves of an undirected tree layer by layer; the last layer is the center(s)
func peelLeaves(n int, edges [][2]int) [][]int {
    // … build adj and degree for an undirected graph …
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
                if degree[nb] == 1 { // just became a leaf
                    next = append(next, nb)
                }
            }
        }
        layer = next
    }
    return layers
}
```

**Why the last layer is the center:** each round removes the two ends of every longest path, shortening all of them by one at each end. The last nodes standing are the middle of the longest path — the nodes that minimize the distance to the farthest leaf. A tree always has one or two centers.

**Dry run** — edges `3–0, 3–1, 3–2, 3–4, 5–4`:

```
0   1   2
 \  |  /
    3
    |
    4
    |
    5
```

| Layer | Removed (current leaves) | Degrees that drop to 1 |
|---|---|---|
| 1 | 0, 1, 2, 5 | 3 (4 → 1), 4 (2 → 1) |
| 2 | 3, 4 | — (nothing left) |

Centers: **3 and 4** — rooting the tree at either gives height 2.

**Run it:** `peelLeaves` in [`toposort.go`](toposort.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Minimum Height Trees | Medium | https://leetcode.com/problems/minimum-height-trees/ |
| 2 | Longest Cycle in a Graph | Hard | https://leetcode.com/problems/longest-cycle-in-a-graph/ |

**What to notice — Problem 1:** Return the last layer of the template. Handle `n == 1` (answer `[0]`) — that node has degree 0, not 1.

**What to notice — Problem 2:** Every node has at most one outgoing edge. Run Kahn's on in-degrees to peel away everything that is *not* on a cycle; whatever remains with in-degree > 0 lies on cycles. Walk each remaining cycle once to measure it, marking nodes as you go.

**Complexity:** both O(V + E).

**Common mistakes:**

- Minimum Height Trees: running BFS from every node (O(n²)) instead of peeling.
- Minimum Height Trees: stopping when one node remains — two centers is also a valid stopping point.
- Longest Cycle: re-walking the same cycle from each of its nodes.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Three questions decide every problem this week: what are the nodes and which way do edges point; do I need an order, a cycle check, or a DP; and Kahn's or DFS?

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Solve Course Schedule once with Kahn's and once with three-color DFS, back to back.

**Self-assessment — can you answer these without hesitation?**

- [ ] Write Kahn's algorithm from memory, including the cycle check
- [ ] Why does directed cycle detection need three colors? Draw the diamond
- [ ] For Course Schedule's `[a, b]`, which way does the edge point?
- [ ] Which pairs of words in Alien Dictionary give information, and why?
- [ ] Why does DP in topological order see final values for every predecessor?
- [ ] Why does peeling leaves find the center of a tree?

---

## Pattern Cheatsheet

| Pattern | Technique | Key move | Detects cycle by | Time |
|---------|-----------|----------|------------------|------|
| Order / feasibility | Kahn's | dequeue in-degree 0, decrement neighbors | `len(order) < n` | O(V + E) |
| Cycle detection / safe nodes | three-color DFS | gray neighbor = cycle; black = done | gray → gray edge | O(V + E) |
| Constraints → graph | build edges, then Kahn's | adjacent pairs, first difference | Kahn's count | O(input + V + E) |
| DP on a DAG | Kahn's order + relaxation | push final value along each edge | — (input is a DAG) | O(V + E) |
| Peel layers | Kahn's on degree 1 / in-degree 0 | remove a whole layer per round | what's left after peeling | O(V + E) |

---

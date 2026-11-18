# Dijkstra & Weighted Shortest Paths Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 18 — BFS week's shortest path plus Heap week's priority queue)  
**Goal:** Write Dijkstra with lazy deletion from memory, know exactly when it's correct (non-negative weights), and pick the right tool: BFS, 0-1 BFS, Dijkstra, Bellman-Ford, or Floyd–Warshall.

---

## What Is Dijkstra's Algorithm?

**Dijkstra's algorithm** finds the shortest path from one node to every other node in a graph where edges have **non-negative weights** (costs, times, distances). It's like BFS, but ordered by total distance instead of number of edges: a **min-heap** always hands you the closest node not yet finalized, and once a node is taken out of the heap, its distance is final.

This week also covers the relatives you pick by looking at the weights: **0-1 BFS** (weights are only 0 or 1), **Bellman-Ford** (negative weights, or at most k edges), and **Floyd–Warshall** (shortest paths between all pairs).

**Typical problems it solves:**

- Network delay time / shortest travel time from a source
- Cheapest or most reliable path (maximum probability)
- Grid paths where moves have costs, or "minimize the worst step" (minimum effort)
- Shortest paths with an extra resource (a coupon, a time limit)
- Cheapest route with at most k stops (Bellman-Ford)
- Distances between every pair of cities (Floyd–Warshall)

---

## The Core Template

```go
type edge struct{ to, w int }
type item struct{ node, dist int } // heap entries, ordered by dist (a min-heap like the Heap week's)

func dijkstra(n int, adj [][]edge, src int) []int {
    dist := make([]int, n)
    for i := range dist {
        dist[i] = math.MaxInt // "not reached yet"
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
```

**Three rules — never break them:**

1. **Non-negative weights only.** Dijkstra's guarantee — "the first time a node is popped, its distance is final" — breaks as soon as an edge can make a path *shorter* later. Negative edges → Bellman-Ford (Day 6).
2. **Skip stale heap entries.** Go's heap has no "decrease key", so when a node's distance improves you just push it again. Older, larger entries for the same node are still in the heap; the `cur.dist > dist[cur.node]` check throws them away. Without it the code is still correct but can do far more work.
3. **Relax with `<`, and update `dist` when you push, not when you pop.** Pushing only on a strict improvement keeps the heap small, and an up-to-date `dist` is what makes the stale check work.

**Brute force → Dijkstra:** BFS finds shortest paths only when every edge costs the same — it explores by number of edges, not by total weight. Trying every path is exponential. Dijkstra is "BFS ordered by distance": the heap always hands you the closest unfinished node, so nodes are finalized in order of distance, exactly like BFS's levels.

**Complexity:** O((V + E) log V) with a binary heap (each successful relaxation pushes once). O(V + E) space.

**Run it:** every template in this file is a real function in [`dijkstra.go`](dijkstra.go), tested in [`dijkstra_test.go`](dijkstra_test.go). From the repo root: `go test ./dijkstra/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Core template  
**Key insight:** When the closest unfinished node is popped, no shorter route to it can still exist: any other route would have to pass through some unfinished node that's already at least as far away, and weights can't be negative.

**No problems today.** Do this instead:

1. Re-read BFS week, Day 4 ("why the first visit is the shortest").
2. Read the template above slowly. Understand each line, especially the stale check.
3. Close this file. Write `dijkstra` (and the min-heap from the Heap week) from memory in Go.
4. Open this file. Compare. Fix any differences. Repeat twice more.
5. Answer this in your head: *Why does one negative edge break Dijkstra?*

   > Take edges `0→1 (4)`, `0→2 (5)`, `2→1 (−3)`. Dijkstra pops `1` at distance 4 and declares it final. Later it pops `2` at distance 5 and finds `5 + (−3) = 2 < 4` — but `1` was already "finished", and anything computed from it used the wrong value. The proof assumed a path can only get longer as you add edges; a negative edge makes it shorter.

**Dry run** — edges `0→1 (4)`, `0→2 (1)`, `2→1 (2)`, `1→3 (1)`, `2→3 (5)`, `3→4 (3)`, source `0`:

| Pop `(dist, node)` | Stale? | Relaxations | `dist` after | Heap after |
|---|---|---|---|---|
| (0, 0) | no | 1: 4, 2: 1 | `[0, 4, 1, ∞, ∞]` | `[(1,2), (4,1)]` |
| (1, 2) | no | 1: 3, 3: 6 | `[0, 3, 1, 6, ∞]` | `[(3,1), (4,1), (6,3)]` |
| (3, 1) | no | 3: 4 | `[0, 3, 1, 4, ∞]` | `[(4,1), (4,3), (6,3)]` |
| (4, 1) | **yes** — dist[1] is already 3 | — | unchanged | `[(4,3), (6,3)]` |
| (4, 3) | no | 4: 7 | `[0, 3, 1, 4, 7]` | `[(6,3), (7,4)]` |
| (6, 3) | **yes** — dist[3] is already 4 | — | unchanged | `[(7,4)]` |
| (7, 4) | no | — | `[0, 3, 1, 4, 7]` | `[]` |

Node `1` was first reached at 4 (directly), then improved to 3 via `2`. The old `(4, 1)` entry stayed in the heap and was skipped as stale — that's rule 2 in action.

---

## Day 2 — Single-Source Shortest Paths

**Pattern:** The template on an explicit weighted graph  
**Key insight:** Build the adjacency list from the edge list once; everything else is the template.

**Template adaptation:** the core template, plus building `adj`:

```go
adj := make([][]edge, n)
for _, e := range edges { // e = {u, v, w}
    adj[e[0]] = append(adj[e[0]], edge{e[1], e[2]})
    // undirected? also: adj[e[1]] = append(adj[e[1]], edge{e[0], e[2]})
}
```

**Run it:** `dijkstra` in [`dijkstra.go`](dijkstra.go) (dry run on Day 1) — the test includes an unreachable node.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Network Delay Time | Medium | https://leetcode.com/problems/network-delay-time/ |
| 2 | Path with Maximum Probability | Medium | https://leetcode.com/problems/path-with-maximum-probability/ |

**What to notice — Problem 1:** Run the template from `k`; the answer is `max(dist)`, or −1 if any node is still unreachable. Nodes are 1-indexed.

**What to notice — Problem 2:** "Shortest" becomes "most probable": a path's value is the **product** of probabilities, and you want the **largest**. Use a max-heap and relax with `prob[u] * p > prob[v]`. It works because multiplying by a number in [0, 1] can never increase a path's value — the same "can't improve later" property that non-negative weights give.

**Complexity:** both O((V + E) log V).

**Common mistakes:**

- Network Delay: forgetting the 1-indexing, or returning `max(dist)` without checking for unreachable nodes.
- Max Probability: keeping a min-heap (pops the *worst* path first).
- Using `math.MaxInt` and then adding to it — check `dist[u] != inf` before adding when weights come from outside the heap.

---

## Day 3 — Dijkstra on a Grid

**Pattern:** Cells are nodes, moves are edges, and the "weight" comes from the cells  
**Key insight:** The DP week's grid problems only moved right/down, so a table filled in order worked. Once you can move in **all four** directions, paths can loop back and there's no fill order — you need Dijkstra.

**Template adaptation:**

```go
// cheapest top-left → bottom-right, 4 directions, entering a cell costs its value
func minPathCostGrid(grid [][]int) int {
    rows, cols := len(grid), len(grid[0])
    id := func(r, c int) int { return r*cols + c } // flatten (r, c) to one node id

    dist := make([]int, rows*cols)
    for i := range dist {
        dist[i] = math.MaxInt
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
            return cur.dist // popped = final, so we can stop early
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
```

**Dry run** — why right/down DP isn't enough:

```
1  1  1
9  9  1
1  1  1
1  9  9
1  1  1
```

The cheap route snakes: right along row 0, down the right column, **left** along row 2, down the left column, then right along row 4 — cost **11**. A right/down-only DP can't take the leftward step, so it has to cross a 9 (its best is 15). Dijkstra finds 11 because the left move is just another edge.

**Run it:** `minPathCostGrid` in [`dijkstra.go`](dijkstra.go) — the test uses this grid.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Path With Minimum Effort | Medium | https://leetcode.com/problems/path-with-minimum-effort/ |
| 2 | Swim in Rising Water | Hard | https://leetcode.com/problems/swim-in-rising-water/ |

**What to notice — Problem 1:** A path's cost is its **largest** single step (`|height difference|`), not the sum. Relax with `max(cur.dist, step)` instead of `cur.dist + step`. Dijkstra still works because taking another step can never make the maximum smaller.

**What to notice — Problem 2:** Same "bottleneck" shape: a path's cost is the highest cell on it, so relax with `max(cur.dist, grid[nr][nc])`. You solved this with Kruskal-style union-find in week 16 — compare. (Binary search on the answer + BFS also works — Binary Search week, Day 4.)

**Complexity:** both O(R · C · log(R · C)).

**Common mistakes:**

- Using a `visited` array and marking on push — a cell can be reached first by a worse path. Mark (or check staleness) on pop.
- Effort: summing the steps instead of taking the max.
- Forgetting to count the starting cell (Swim starts at `grid[0][0]`, not 0).

---

## Day 4 — 0-1 BFS

**Pattern:** When every weight is 0 or 1, replace the heap with a deque: 0-edges go to the front, 1-edges to the back  
**Key insight:** The deque stays sorted by distance — the front holds distance `d`, the back holds `d + 1` — so it works like Dijkstra's heap without the log factor.

**Template adaptation:**

```go
func zeroOneBFS(n int, adj [][]edge, src int) []int {
    dist := make([]int, n)
    for i := range dist {
        dist[i] = math.MaxInt
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
```

**Why it works:** BFS's invariant was "the queue holds distance `d` then `d+1`". A 0-edge leads to another node at distance `d`, so it belongs at the front with the other `d`s; a 1-edge leads to `d + 1`, the back. The invariant survives both, so nodes still come out in distance order.

**Dry run** — edges `0→1 (1)`, `0→2 (0)`, `2→3 (1)`, `1→3 (0)`, `3→4 (0)`, `2→1 (0)`:

| Pop | Improvements | Deque after | `dist` after |
|---|---|---|---|
| 0 | 1 = 1 (back), 2 = 0 (front) | `[2, 1]` | `[0, 1, 0, ∞, ∞]` |
| 2 | 3 = 1 (back), 1 = 0 (front) | `[1, 1, 3]` | `[0, 0, 0, 1, ∞]` |
| 1 | 3 = 0 (front) | `[3, 1, 3]` | `[0, 0, 0, 0, ∞]` |
| 3 | 4 = 0 (front) | `[4, 1, 3]` | `[0, 0, 0, 0, 0]` |
| 4, 1, 3 | — | `[]` | final |

Node 1 was first reached at cost 1 (direct), then improved to 0 via the free edge from 2. Duplicate entries (the later `1` and `3`) find nothing to improve — the same "stale entry" idea as Dijkstra.

**Run it:** `zeroOneBFS` in [`dijkstra.go`](dijkstra.go) — the test checks it against `dijkstra` on the same graph.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Minimum Cost to Make at Least One Valid Path in a Grid | Hard | https://leetcode.com/problems/minimum-cost-to-make-at-least-one-valid-path-in-a-grid/ |
| 2 | Minimum Obstacle Removal to Reach Corner | Hard | https://leetcode.com/problems/minimum-obstacle-removal-to-reach-corner/ |

**What to notice — Problem 1:** Moving in the direction a cell's arrow points costs 0; any other direction costs 1 (you changed the sign). Pure 0-1 BFS on the grid.

**What to notice — Problem 2:** Entering an empty cell costs 0, entering an obstacle costs 1 (you remove it). Again pure 0-1 BFS.

**Complexity:** both O(R · C) — no log factor.

**Common mistakes:**

- Treating it as plain BFS (ignoring the 0-weights) — wrong distances.
- `append([]int{x}, deque...)` is O(n) per push-front. Fine for learning; for big inputs use a ring buffer or two slices (a "front stack" and a "back queue").

---

## Day 5 — Dijkstra with Extra State or Bookkeeping

**Pattern:** The node alone isn't enough to describe "where you are" — make the state `(node, something)` and run Dijkstra over states  
**Key insight:** If two arrivals at the same node can have different futures (one still has a coupon, one doesn't), they're different states. Give each its own `dist` entry.

**Template adaptation:**

```go
// cheapest src→dst when at most one edge may be used for free; state = (node, used)
func cheapestWithOneFree(n int, adj [][]edge, src, dst int) int {
    dist := make([]int, 2*n) // state id = node*2 + used
    for i := range dist {
        dist[i] = math.MaxInt
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
            if next := e.to*2 + used; cur.dist+e.w < dist[next] { // pay for the edge
                dist[next] = cur.dist + e.w
                heap.Push(h, item{next, dist[next]})
            }
            if used == 0 { // or use the free pass on this edge
                if next := e.to*2 + 1; cur.dist < dist[next] {
                    dist[next] = cur.dist
                    heap.Push(h, item{next, cur.dist})
                }
            }
        }
    }
    return -1
}
```

**Dry run** — edges `0→1 (10)`, `1→3 (10)`, `0→2 (3)`, `2→3 (15)`, `3→4 (1)`, from 0 to 4. States in the order they're finalized:

| Popped state `(node, used)` | Cost |
|---|---|
| (0, used=0) | 0 |
| (1, used=1) | 0 — took `0→1` for free |
| (2, used=1) | 0 — took `0→2` for free |
| (2, used=0) | 3 |
| (3, used=1) | 3 — from (2, used=0), `2→3` for free |
| (4, used=1) | **4** |

Without the free pass the cheapest path costs 19 (`0→2→3→4`). The best use of the coupon is the expensive `2→3` edge — Dijkstra found that without trying each edge as "the free one" separately.

**Run it:** `cheapestWithOneFree` in [`dijkstra.go`](dijkstra.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Ways to Arrive at Destination | Medium | https://leetcode.com/problems/number-of-ways-to-arrive-at-destination/ |
| 2 | Minimum Cost to Reach Destination in Time | Hard | https://leetcode.com/problems/minimum-cost-to-reach-destination-in-time/ |

**What to notice — Problem 1:** Bookkeeping, not state: keep `ways[v]` next to `dist[v]`. Strictly shorter path → `dist[v] = nd`, `ways[v] = ways[u]`. Equal path → `ways[v] += ways[u]`. Times can exceed 32 bits — use `int` (64-bit in Go) and take the modulo on `ways`.

**What to notice — Problem 2:** Two quantities: minimize **fee** with a limit on **time**. Run Dijkstra keyed by fee, and keep `bestTime[v]`; only push `(v, time)` if `time < bestTime[v]` — an arrival that's both more expensive *and* slower is useless, but a more expensive and faster one may still matter.

**Complexity:** O(S log S) where S = number of states (×2 for one coupon, ×k for k coupons).

**Common mistakes:**

- One `dist` per node instead of per state — the coupon-used arrival overwrites the coupon-unused one.
- Ways: adding to `ways[v]` on a *stale* pop, double-counting paths.
- Stopping at the first time you see `dst` in a relaxation instead of when you **pop** it.

---

## Day 6 — Beyond Dijkstra: Bellman-Ford and Floyd–Warshall

**Pattern:** Bellman-Ford relaxes **every edge** for a fixed number of rounds; Floyd–Warshall tries **every intermediate node** for every pair  
**Key insight:** After round `k` of Bellman-Ford, `dist[v]` is the shortest path using **at most `k` edges**. That's exactly what "at most k stops" problems need — and it handles negative weights.

**Template adaptation:**

```go
// shortest distances from src using at most maxEdges edges; weights may be negative
func bellmanFord(n int, edges [][3]int, src, maxEdges int) []int {
    dist := make([]int, n)
    for i := range dist {
        dist[i] = math.MaxInt
    }
    dist[src] = 0

    for round := 0; round < maxEdges; round++ {
        prev := make([]int, n)
        copy(prev, dist) // relax only from last round's values
        for _, e := range edges {
            u, v, w := e[0], e[1], e[2]
            if prev[u] != math.MaxInt && prev[u]+w < dist[v] {
                dist[v] = prev[u] + w
            }
        }
    }
    return dist
}

// all-pairs shortest paths: allow node k as a stopover, for k = 0..n-1
for k := 0; k < n; k++ {
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if dist[i][k] != inf && dist[k][j] != inf && dist[i][k]+dist[k][j] < dist[i][j] {
                dist[i][j] = dist[i][k] + dist[k][j]
            }
        }
    }
}
```

**Why `prev` matters:** without the copy, an edge relaxed early in a round can feed another edge later in the *same* round, so one round could use several edges — breaking "round k = at most k edges". (Same trap as the DP week's 0/1 knapsack loop direction.)

**Dry run** — LeetCode 787's example: `0→1 (100)`, `1→2 (100)`, `0→2 (500)`, `1→3 (600)`, `2→3 (200)`:

| After round (= max edges) | `dist` | Cheapest to 3 |
|---|---|---|
| 1 | `[0, 100, 500, ∞]` | unreachable |
| 2 | `[0, 100, 200, 700]` | 700 (`0→1→3`) |
| 3 | `[0, 100, 200, 400]` | 400 (`0→1→2→3`) |

Plain Dijkstra returns 400 regardless of the limit. With "at most 1 stop" (= 2 edges) the right answer is 700 — only the round-limited version gets it.

**Run it:** `bellmanFord` and `floydWarshall` in [`dijkstra.go`](dijkstra.go) — the test includes this example and a negative edge.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Cheapest Flights Within K Stops | Medium | https://leetcode.com/problems/cheapest-flights-within-k-stops/ |
| 2 | Find the City With the Smallest Number of Neighbors at a Threshold Distance | Medium | https://leetcode.com/problems/find-the-city-with-the-smallest-number-of-neighbors-at-a-threshold-distance/ |

**What to notice — Problem 1:** `k` stops = `k + 1` edges → `bellmanFord(n, flights, src, k+1)[dst]`.

**What to notice — Problem 2:** n ≤ 100, so all-pairs Floyd–Warshall (O(n³) = 10⁶) is the simplest correct choice. Count neighbors within the threshold for each city; ties go to the larger index.

**Complexity:** Bellman-Ford O(k · E). Floyd–Warshall O(n³) time, O(n²) space.

**Common mistakes:**

- Cheapest Flights: relaxing in place (no `prev` copy) — paths sneak in extra edges.
- Floyd–Warshall: putting the `k` loop **inside** — it must be the outermost loop.
- Adding to `inf` without checking — overflow to a negative number that wins every `min`.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Choosing the algorithm is the hard part. Look at the weights and the constraints, then pick from this ladder: all equal → BFS; 0 or 1 → 0-1 BFS; non-negative → Dijkstra; negative or "at most k edges" → Bellman-Ford; all pairs with small n → Floyd–Warshall.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. For each problem this week, say out loud which rung of the ladder it's on and why the rung below it wouldn't work.

**Self-assessment — can you answer these without hesitation?**

- [ ] Write Dijkstra with lazy deletion from memory
- [ ] Why does a negative edge break Dijkstra? Give a 3-node example
- [ ] What is a "stale" heap entry, and what happens if you don't skip it?
- [ ] Why does 0-1 BFS push 0-edges to the front?
- [ ] When do you need `(node, extra)` states instead of just `node`?
- [ ] Why does Bellman-Ford copy `dist` into `prev` each round?

---

## Pattern Cheatsheet

| Weights / constraint | Algorithm | Data structure | Time |
|----------------------|-----------|----------------|------|
| All equal | BFS (BFS week) | queue | O(V + E) |
| 0 or 1 | 0-1 BFS | deque | O(V + E) |
| Non-negative | Dijkstra | min-heap, lazy deletion | O((V + E) log V) |
| Bottleneck (max edge on path) | Dijkstra with `max` instead of `+` | min-heap | O((V + E) log V) |
| Extra resource (coupons, time) | Dijkstra over `(node, state)` | min-heap | O(S log S) |
| Negative, or at most k edges | Bellman-Ford | edge list + `prev` copy | O(k · E) |
| All pairs, small n | Floyd–Warshall | n × n matrix | O(n³) |

---

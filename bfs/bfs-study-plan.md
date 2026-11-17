# BFS Interview Prep — 7-Day Study Plan

**Level:** Beginner  
**Goal:** Write BFS from memory, recognize patterns, adapt the template to any problem type.

---

## The One Template

Every BFS problem is a variation of this. Learn it until you can write it in 60 seconds.

```go
func bfs(start Node) {
    queue := []Node{start}
    visited := map[Node]bool{start: true}

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]

        for _, neighbor := range getNeighbors(node) {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

**Three rules — never break them:**

1. Mark visited **on enqueue**, not on dequeue. If you mark on dequeue, duplicates pile up in the queue.
2. In this template one outer-loop iteration processes **one node**, not one level. When you need levels, snapshot `levelSize := len(queue)` at the top of the outer loop and drain exactly that many nodes in an inner loop (Day 2) — then one outer iteration = one level.
3. `getNeighbors` is the only thing that changes. Everything else stays identical across all BFS problems.

**Brute force → BFS:** To find the shortest route from A to B you could list every path and keep the shortest — exponential in the number of nodes. BFS visits nodes *in order of distance* (everything 1 step away, then everything 2 steps away, …), so the first time it touches B is already the shortest way there. No path is ever explored twice.

**Complexity:** O(V + E) time — every node is enqueued at most once (rule 1 guarantees it) and every edge is looked at once from each end. O(V) space for the queue and the visited set.

**Run it:** every template in this file is a real function in [`bfs.go`](bfs.go), tested in [`bfs_test.go`](bfs_test.go). From the repo root: `go test ./bfs/ -v`.

---

## Day 1 — Template + Queue Mechanics

**Pattern:** Core template  
**Key insight:** Every BFS problem is the same template; only `getNeighbors` changes.

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each line.
2. Close this file. Write the template from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this question in your head: *Why must visited be marked on enqueue, not dequeue?*

   > Because if you mark on dequeue, every neighbor gets added to the queue multiple times before any of them are processed. In a graph with back edges, this means O(E) duplicates in the queue instead of O(V) nodes. Marking on enqueue prevents a node from ever being enqueued twice.

**Dry run** — the core template on this graph, starting at `0`:

```
0 — 1
|   |
2 — 3 — 4
```

| Step | Dequeue | Neighbors checked | Enqueued (and marked) | Queue after |
|---|---|---|---|---|
| 1 | `0` | 1, 2 | 1, 2 | `[1, 2]` |
| 2 | `1` | 0, 3 | 3 | `[2, 3]` |
| 3 | `2` | 0, 3 | — | `[3]` |
| 4 | `3` | 1, 2, 4 | 4 | `[4]` |
| 5 | `4` | 3 | — | `[]` |

Two things to see in this table:

- **Step 3 is rule 1 in action.** Node `2` sees neighbor `3`, but `3` was already marked when `1` enqueued it, so it's skipped. If you marked on dequeue instead, `3` would be pushed a second time here (6 pushes instead of 5). On a dense graph those duplicates multiply.
- **Dequeue order is sorted by distance from `0`:** `0` (0 steps), `1`/`2` (1 step), `3` (2 steps), `4` (3 steps). Day 4's shortest-path guarantee is exactly this property.

**Go note:** `queue = queue[1:]` is O(1), but the dequeued items stay in the backing array until the next `append` reallocates. That's fine for interviews; just don't be surprised that memory tracks *total* enqueued, not the current queue length.

---

## Day 2 — Tree Level-Order Traversal

**Pattern:** BFS on trees  
**Key insight:** `getNeighbors(node)` returns `[node.Left, node.Right]`, skipping nil. Trees have no cycles, so no visited set is needed.

**Template adaptation:**

```go
func bfs(root *TreeNode) {
    if root == nil {
        return
    }
    queue := []*TreeNode{root}
    // no visited map — trees have no cycles

    for len(queue) > 0 {
        levelSize := len(queue) // snapshot: how many nodes are in this level
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]

            // process node here

            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        // one full level processed — increment level counter here if needed
    }
}
```

The `levelSize` snapshot is the key move: it lets you process exactly one level per outer loop iteration, without the queue growing under you mid-level.

**Why it works:** At the top of each outer iteration the queue holds *exactly* one level — every node of depth `d` and nothing else. The snapshot `levelSize` freezes that count, so the inner loop dequeues precisely the depth-`d` nodes while their children (depth `d+1`) pile up behind them. When the inner loop ends, the queue again holds exactly one level. That invariant is the whole trick.

**Dry run** on this tree:

```
      1
     / \
    2   3
   / \   \
  4   5   6
```

| Level | `levelSize` | Dequeued by inner loop | Queue after the level |
|---|---|---|---|
| 0 | 1 | `1` | `[2, 3]` |
| 1 | 2 | `2, 3` | `[4, 5, 6]` |
| 2 | 3 | `4, 5, 6` | `[]` |

**Run it:** `levelCounts` in [`bfs.go`](bfs.go) — the template returning `[1, 2, 3]` for this tree.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Binary Tree Level Order Traversal | Easy | https://leetcode.com/problems/binary-tree-level-order-traversal/ |
| 2 | Minimum Depth of Binary Tree | Easy | https://leetcode.com/problems/minimum-depth-of-binary-tree/ |

**What to notice — Problem 1:** You need to return a `[][]int` where each inner slice is one level. Collect node values inside the inner loop, append the slice after the inner loop ends.

**What to notice — Problem 2:** You don't need to finish all levels. Return `depth` the moment you dequeue a node with no children (a leaf). BFS guarantees that's the shallowest leaf.

**Complexity:** both problems O(n) time. Space is O(w), where `w` is the widest level — up to about n/2 for the bottom level of a full tree. Problem 2 can stop early at the first leaf, so it's often much faster than O(n) in practice.

**Common mistakes:**

- Writing `for i := 0; i < len(queue); i++` instead of snapshotting `levelSize` first. `len(queue)` grows as you append children, so levels bleed into each other.
- Forgetting the `root == nil` check — `[]*TreeNode{nil}` dereferences nil on the first dequeue.
- Problem 2: a node with one child is **not** a leaf. For the chain `1 → 2`, min depth is 2, not 1.

---

## Day 3 — Graph Traversal / Connected Components

**Pattern:** BFS on explicit graphs  
**Key insight:** `getNeighbors(node)` returns entries from an adjacency list or matrix. The visited set is now mandatory — graphs have cycles and BFS will loop forever without it.

**Template adaptation:**

```go
func bfs(graph [][]int, start int) {
    visited := map[int]bool{start: true}
    queue := []int{start}

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]

        for _, neighbor := range graph[node] { // adjacency list
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
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

// To count connected components: run BFS from every unvisited node, increment count each time.
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
```

**Why it works:** A BFS started at `s` enqueues every node reachable from `s` and nothing else, so one BFS run = one connected component. Because `visited` is shared across runs, a node that was reached by an earlier run is skipped by the outer loop — every component is counted exactly once.

**Dry run** — `countComponents(7, graph)` on the graph from Day 1 plus a separate edge `5 — 6`:

| Outer `i` | `visited[i]` before? | BFS from `i` marks | `count` |
|---|---|---|---|
| 0 | no | 0, 1, 2, 3, 4 | 1 |
| 1–4 | yes | — (skipped) | 1 |
| 5 | no | 5, 6 | 2 |
| 6 | yes | — (skipped) | 2 |

**Run it:** `bfsGraph`, `bfsFrom`, and `countComponents` in [`bfs.go`](bfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Islands | Medium | https://leetcode.com/problems/number-of-islands/ |
| 2 | Number of Provinces | Medium | https://leetcode.com/problems/number-of-provinces/ |

**What to notice — Problem 1:** The grid IS the graph. Each `'1'` cell is a node. Its neighbors are the 4 adjacent cells (up/down/left/right) that are also `'1'`. Mark cells visited by flipping them to `'0'` in-place (or use a separate visited set). Count how many times you start a new BFS.

**What to notice — Problem 2:** The input is an `n×n` adjacency matrix. `isConnected[i][j] == 1` means there's an edge between city `i` and city `j`. Build a neighbor list on the fly inside `getNeighbors`. Count connected components.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Number of Islands | O(R·C) — each cell enqueued at most once | O(R·C) worst case for visited/queue |
| Number of Provinces | O(n²) — scanning a full row of the matrix for every city | O(n) |

**Common mistakes:**

- Not marking the **start** node visited before the loop. Its neighbor will enqueue it again.
- Islands: flipping `'1'` → `'0'` on *dequeue* instead of enqueue. Same cell gets enqueued many times and the solution times out.
- Provinces: treating `isConnected` as an adjacency list. It's an `n×n` matrix — `isConnected[i]` is a row of 0/1 flags, not a list of neighbors.

---

## Day 4 — Shortest Path on Unweighted Graph

**Pattern:** BFS for minimum distance  
**Key insight:** BFS explores level by level. The first time it reaches a node, that's guaranteed to be the shortest path — no revisit can ever be shorter. Track distance by storing it alongside each node in the queue.

**Template adaptation:**

```go
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
    return -1 // no path found
}
```

For grids, the same pattern applies but `getNeighbors` returns adjacent cells.

**Why it works:** The queue only ever holds nodes at distance `d` followed by nodes at distance `d+1` — never three different distances at once. So nodes leave the queue in non-decreasing order of distance. When `v` is first discovered from `u`, every node closer than `u` has already been dequeued and none of them reached `v`, so no shorter path to `v` exists. That's why the visited check can be permanent: a later discovery can only be the same length or longer.

This only holds when every edge costs the same. With weighted edges, a 1-hop path of cost 10 can lose to a 3-hop path of cost 3 — that's Dijkstra, which swaps the queue for a min-heap.

**Dry run** — `shortestPath(graph, 0, 4)` on the Day 1 graph:

| Dequeue `(node, dist)` | Enqueued |
|---|---|
| `(0, 0)` | `(1, 1)`, `(2, 1)` |
| `(1, 1)` | `(3, 2)` |
| `(2, 1)` | — |
| `(3, 2)` | `(4, 3)` |
| `(4, 3)` | — **node == end, return 3** |

**Run it:** `shortestPath` in [`bfs.go`](bfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Shortest Path in Binary Matrix | Medium | https://leetcode.com/problems/shortest-path-in-binary-matrix/ |
| 2 | 01 Matrix | Medium | https://leetcode.com/problems/01-matrix/ |

**What to notice — Problem 1:** 8-directional movement (diagonals included). Neighbors are all 8 adjacent cells where the value is `0`. Return `-1` if the bottom-right cell is never reached. Edge case: if `grid[0][0] == 1`, return `-1` immediately.

**What to notice — Problem 2:** You want the distance from each cell to its nearest `0`. Don't BFS from each cell individually (too slow). Instead, enqueue ALL `0` cells at the start — this is multi-source BFS. When you reach a `1` cell, its distance is already the minimum.

**Complexity:** O(V + E) for a graph. For an R×C grid that's O(R·C) — each cell has at most 4 (or 8) neighbors, a constant.

**Common mistakes:**

- Problem 1 counts **cells** on the path, not edges. A 1×1 grid of `0` returns 1, so start `dist` at 1, not 0.
- Problem 1: forgetting to check `grid[n-1][n-1]` as well as `grid[0][0]`.
- Problem 2: one BFS per `1` cell is O((R·C)²). Start from all `0`s at once (that's tomorrow's pattern).

---

## Day 5 — Multi-Source BFS

**Pattern:** BFS from multiple starting nodes simultaneously  
**Key insight:** Enqueue ALL source nodes before the BFS loop begins. BFS spreads outward from all sources at once, like a wave front. This is always more efficient than running BFS once per source.

**Template adaptation:**

```go
type Cell struct{ r, c int }

// returns how many levels the wave needed to reach every reachable cell
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
            minutes++ // only count a level if the wave actually reached new cells
        }
    }
    return minutes
}
```

**Why it works:** Putting every source in the queue at distance 0 is exactly the same as adding one imaginary "super-source" node with an edge to each real source, then running ordinary BFS from it. So Day 4's guarantee carries over unchanged: the first time the wave reaches a cell, that's its distance to the **nearest** source.

**Dry run** — 3×5 grid, sources at `(0,0)` and `(2,4)`. Each number is the level (minute) at which the wave reached that cell:

```
0 1 2 3 2
1 2 3 2 1
2 3 2 1 0
```

The two waves meet in the middle: every cell marked `3` is exactly 3 steps from *both* sources, and every other cell automatically got the distance to whichever source is closer. `multiSourceBFS` returns `3`, the last level that reached new cells. Running BFS once per source would do the same work twice and then need a `min` over the results.

**Run it:** `multiSourceBFS` in [`bfs.go`](bfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Rotting Oranges | Medium | https://leetcode.com/problems/rotting-oranges/ |
| 2 | Walls and Gates | Medium | https://leetcode.com/problems/walls-and-gates/ (LeetCode Premium — use lintcode.com/problem/663 if needed) |

**What to notice — Problem 1:** All rotten oranges (`2`) are sources. Each BFS level = 1 minute. After BFS, scan the grid for any remaining `1` — if found, return `-1` (some oranges can never rot).

**What to notice — Problem 2:** All gates (`0`) are sources. Fill each empty room (`INF`) with its BFS level (distance to nearest gate). Walls (`-1`) are impassable — skip them in `getNeighbors`.

**Complexity:** O(R·C) time and space **no matter how many sources** there are. Once-per-source BFS would be O(S·R·C).

**Common mistakes:**

- Rotting Oranges: incrementing `minutes` after the final level even though it rotted nothing new — that's why the template guards with `if len(queue) > 0`.
- Rotting Oranges: no fresh oranges at all should return 0, and any fresh orange left after BFS means -1.
- Walls and Gates: forgetting that walls (`-1`) are never enqueued, and gates are the sources — not the empty rooms.

---

## Day 6 — BFS on Implicit Graphs (State-Space)

**Pattern:** BFS where the graph is not given — you generate it on the fly  
**Key insight:** `getNeighbors(state)` generates all valid next states from the current state. The "graph" only exists as transitions between states. This pattern covers most interview "hard" BFS problems.

**How to recognize it:** The problem gives you a start state and an end state, and asks for minimum steps. The graph is not an input — you derive it from rules (what moves are legal?).

**Template adaptation:**

```go
func bfs(start, end string, deadends map[string]bool) int {
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

            for _, next := range getNeighbors(cur) { // generate next states
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
```

The visited set here is a `map[string]bool` because states are strings. The type of visited always matches the type of the state.

**Why it works:** BFS never needed the graph to be stored — it only ever calls `getNeighbors`. Treat each state as a node and each legal move as an edge of cost 1, and Day 4's proof applies unchanged: the level at which you first reach the target is the minimum number of moves.

**Dry run** — states are 3-bit strings, a move flips one bit, `deadends = {100, 010}`, start `000`, target `111`:

| `steps` | Frontier (this level) |
|---|---|
| 0 | `000` |
| 1 | `001` (`100` and `010` are dead ends) |
| 2 | `101`, `011` |
| 3 | `111` → return 3 |

Even with two of the three first moves blocked, the answer is still 3 — BFS found the detour automatically.

**Run it:** `minSteps` in [`bfs.go`](bfs.go); `getNeighbors` is passed in as a function so you can plug in any move rule.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Open the Lock | Medium | https://leetcode.com/problems/open-the-lock/ |
| 2 | Minimum Genetic Mutation | Medium | https://leetcode.com/problems/minimum-genetic-mutation/ |

**What to notice — Problem 1:** State = 4-character digit string. `getNeighbors` returns 8 strings (each of the 4 wheels turned +1 or -1, wrapping 0↔9). `deadends` acts as an extra visited filter. Start = `"0000"`.

**What to notice — Problem 2:** State = 8-character gene string. `getNeighbors` returns all strings that differ from the current state by exactly one character AND appear in the bank. If the end state isn't in the bank, return `-1` immediately.

**Complexity:** O(S × N) where S is the number of reachable states and N is the cost of generating one state's neighbors.

| Problem | States | Neighbors per state |
|---|---|---|
| Open the Lock | 10⁴ wheel combinations | 8 (4 wheels × ±1) |
| Minimum Genetic Mutation | at most `len(bank)` + 1 | 8 positions × 3 other letters |

**Common mistakes:**

- Open the Lock: `"0000"` itself can be a dead end → return -1 before the loop.
- Open the Lock: wraparound both ways — `'9' + 1` is `'0'` and `'0' - 1` is `'9'`.
- Genetic Mutation: generating neighbors that aren't in the bank. Only bank strings are valid states.
- Using a non-comparable type (slice) as the state. Strings and arrays like `[4]byte` work as map keys; slices don't.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Pattern recognition is the skill. You don't need to memorize 50 problems — you need to see the template in each one within 2 minutes.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. After solving (or attempting), write out the `getNeighbors` function for that pattern from memory.

**Self-assessment:**

After Day 7, you should be able to answer these without hesitation:

- [ ] Write the BFS template from memory in under 60 seconds
- [ ] What changes between tree BFS and graph BFS?
- [ ] Why don't tree problems need a visited set?
- [ ] How do you track level count in BFS?
- [ ] What is multi-source BFS and when do you use it?
- [ ] How do you represent state in an implicit graph problem?

---

## Pattern Cheatsheet

| Pattern | Visited needed? | `getNeighbors` returns | Level tracking? | Time |
|---------|----------------|----------------------|----------------|------|
| Tree level-order | No | `[node.Left, node.Right]` (skip nil) | Yes — use `levelSize` snapshot | O(n) |
| Graph traversal | Yes | adjacency list entries | Usually no | O(V + E) |
| Grid shortest path | Yes | 4 or 8 adjacent cells | Via `dist` field in queue state | O(R·C) |
| Multi-source | Yes | same as grid/graph | Via `levelSize` snapshot | O(R·C), independent of source count |
| Implicit / state-space | Yes (as map[State]bool) | generated next states | Via `levelSize` snapshot | O(states × neighbors per state) |

---

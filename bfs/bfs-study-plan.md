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
2. One full pass through the outer loop = one level. If you need level count, add a counter that increments after each full pass.
3. `getNeighbors` is the only thing that changes. Everything else stays identical across all BFS problems.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Binary Tree Level Order Traversal | Easy | https://leetcode.com/problems/binary-tree-level-order-traversal/ |
| 2 | Minimum Depth of Binary Tree | Easy | https://leetcode.com/problems/minimum-depth-of-binary-tree/ |

**What to notice — Problem 1:** You need to return a `[][]int` where each inner slice is one level. Collect node values inside the inner loop, append the slice after the inner loop ends.

**What to notice — Problem 2:** You don't need to finish all levels. Return `depth` the moment you dequeue a node with no children (a leaf). BFS guarantees that's the shallowest leaf.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Islands | Medium | https://leetcode.com/problems/number-of-islands/ |
| 2 | Number of Provinces | Medium | https://leetcode.com/problems/number-of-provinces/ |

**What to notice — Problem 1:** The grid IS the graph. Each `'1'` cell is a node. Its neighbors are the 4 adjacent cells (up/down/left/right) that are also `'1'`. Mark cells visited by flipping them to `'0'` in-place (or use a separate visited set). Count how many times you start a new BFS.

**What to notice — Problem 2:** The input is an `n×n` adjacency matrix. `isConnected[i][j] == 1` means there's an edge between city `i` and city `j`. Build a neighbor list on the fly inside `getNeighbors`. Count connected components.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Shortest Path in Binary Matrix | Medium | https://leetcode.com/problems/shortest-path-in-binary-matrix/ |
| 2 | 01 Matrix | Medium | https://leetcode.com/problems/01-matrix/ |

**What to notice — Problem 1:** 8-directional movement (diagonals included). Neighbors are all 8 adjacent cells where the value is `0`. Return `-1` if the bottom-right cell is never reached. Edge case: if `grid[0][0] == 1`, return `-1` immediately.

**What to notice — Problem 2:** You want the distance from each cell to its nearest `0`. Don't BFS from each cell individually (too slow). Instead, enqueue ALL `0` cells at the start — this is multi-source BFS. When you reach a `1` cell, its distance is already the minimum.

---

## Day 5 — Multi-Source BFS

**Pattern:** BFS from multiple starting nodes simultaneously  
**Key insight:** Enqueue ALL source nodes before the BFS loop begins. BFS spreads outward from all sources at once, like a wave front. This is always more efficient than running BFS once per source.

**Template adaptation:**

```go
type Cell struct{ r, c int }

func multiSourceBFS(grid [][]int, sources []Cell) {
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
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Rotting Oranges | Medium | https://leetcode.com/problems/rotting-oranges/ |
| 2 | Walls and Gates | Medium | https://leetcode.com/problems/walls-and-gates/ (LeetCode Premium — use lintcode.com/problem/663 if needed) |

**What to notice — Problem 1:** All rotten oranges (`2`) are sources. Each BFS level = 1 minute. After BFS, scan the grid for any remaining `1` — if found, return `-1` (some oranges can never rot).

**What to notice — Problem 2:** All gates (`0`) are sources. Fill each empty room (`INF`) with its BFS level (distance to nearest gate). Walls (`-1`) are impassable — skip them in `getNeighbors`.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Open the Lock | Medium | https://leetcode.com/problems/open-the-lock/ |
| 2 | Minimum Genetic Mutation | Medium | https://leetcode.com/problems/minimum-genetic-mutation/ |

**What to notice — Problem 1:** State = 4-character digit string. `getNeighbors` returns 8 strings (each of the 4 wheels turned +1 or -1, wrapping 0↔9). `deadends` acts as an extra visited filter. Start = `"0000"`.

**What to notice — Problem 2:** State = 8-character gene string. `getNeighbors` returns all strings that differ from the current state by exactly one character AND appear in the bank. If the end state isn't in the bank, return `-1` immediately.

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

| Pattern | Visited needed? | `getNeighbors` returns | Level tracking? |
|---------|----------------|----------------------|----------------|
| Tree level-order | No | `[node.Left, node.Right]` (skip nil) | Yes — use `levelSize` snapshot |
| Graph traversal | Yes | adjacency list entries | Usually no |
| Grid shortest path | Yes | 4 or 8 adjacent cells | Via `dist` field in queue state |
| Multi-source | Yes | same as grid/graph | Via `levelSize` snapshot |
| Implicit / state-space | Yes (as map[State]bool) | generated next states | Via `levelSize` snapshot |

---

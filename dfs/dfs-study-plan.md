# DFS Interview Prep — 7-Day Study Plan

**Level:** Beginner (completed BFS week)  
**Goal:** Write recursive and iterative DFS from memory, recognize patterns, implement backtracking.

---

## The Recursive Template

### Generic (graphs)

Works for any graph. The tree version below is a specialization of this.

```go
func dfs(node Node, visited map[Node]bool) {
    if visited[node] {
        return
    }
    visited[node] = true

    for _, neighbor := range getNeighbors(node) {
        dfs(neighbor, visited)
    }
}
```

### Tree specialization (Days 1–3)

Trees have no cycles so no visited map is needed. `getNeighbors` is just `Left` and `Right`.

```go
func dfs(node *TreeNode) {
    if node == nil {
        return
    }

    // process node here (preorder)

    dfs(node.Left)
    dfs(node.Right)

    // or process here (postorder)
}
```

| | Generic | Tree |
|-|---------|------|
| Base case | `visited[node]` | `node == nil` |
| Visited set | required | not needed |
| Neighbors | `getNeighbors(node)` | `Left`, `Right` |

**Three rules — never break them:**

1. **Base case first** — always check nil/invalid before doing anything else.
2. **Where you process determines the order** — before recursive calls = preorder (root first), after = postorder (children first).
3. **The call stack IS the stack** — you don't need an explicit stack. Recursion handles backtracking automatically.

---

## Day 1 — Recursive DFS Template

**Pattern:** Core recursive template  
**Key insight:** The call stack handles backtracking — you don't need an explicit stack.

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each line.
2. Close this file. Write the template from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *What is the difference between preorder and postorder, and when would you use each?*

   > **Preorder** (process before recursing): used when you need to pass information *down* the tree — e.g., passing a running sum to children. **Postorder** (process after recursing): used when you need to combine results *up* from children — e.g., computing max depth by taking `1 + max(left, right)`.

---

## Day 2 — Tree DFS: Depth & Existence

**Pattern:** DFS for a single value or measurement on a tree  
**Key insight:** Return a value up the call stack — each recursive call combines results from left and right subtrees.

**Template adaptation:**

```go
// Postorder: combine left and right results after recursing
func maxDepth(node *TreeNode) int {
    if node == nil {
        return 0
    }

    left := maxDepth(node.Left)
    right := maxDepth(node.Right)

    return 1 + max(left, right) // combine after children
}

// Preorder: pass information down to children
func hasPathSum(node *TreeNode, remaining int) bool {
    if node == nil {
        return false
    }
    remaining -= node.Val
    if node.Left == nil && node.Right == nil { // leaf
        return remaining == 0
    }
    return hasPathSum(node.Left, remaining) || hasPathSum(node.Right, remaining)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Maximum Depth of Binary Tree | Easy | https://leetcode.com/problems/maximum-depth-of-binary-tree/ |
| 2 | Path Sum | Easy | https://leetcode.com/problems/path-sum/ |
| 3 | Lowest Common Ancestor of a Binary Tree | Medium | https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/ |

**What to notice — Problem 1:** This is postorder — you need both children's depths before you can compute the current node's depth. Return `1 + max(left, right)`. Base case: `nil` returns 0.

**What to notice — Problem 2:** This is preorder — subtract `node.Val` from `remaining` as you go down. At a leaf, check if `remaining == 0`. Don't return true at a nil node — only at a leaf (both children nil).

**What to notice — Problem 3 (LCA):** Postorder — recurse left and right first, then decide at the current node. If the current node is `p` or `q`, return it. If left and right both return non-nil, the current node is the LCA. If only one side returns non-nil, bubble that up.

---

## Day 3 — Tree DFS: Root-to-Leaf Paths

**Pattern:** DFS while tracking the path from root to current node  
**Key insight:** Pass the current path as a parameter down the call stack; add on the way down, remove on the way back up (backtracking).

**Template adaptation:**

```go
func dfs(node *TreeNode, path []int, result *[][]int) {
    if node == nil {
        return
    }

    path = append(path, node.Val) // add on the way down

    if node.Left == nil && node.Right == nil { // leaf
        // make a copy — path slice is reused across calls
        tmp := make([]int, len(path))
        copy(tmp, path)
        *result = append(*result, tmp)
    }

    dfs(node.Left, path, result)
    dfs(node.Right, path, result)

    // no explicit removal needed — the slice header (pointer, len, cap) is passed by value
    // the caller's len is unchanged when this call returns, so it never sees the appended element
}
```

**Why copy the path at the leaf:** Go slices share underlying arrays. If you append `path` directly to `result`, all entries in `result` will point to the same array and overwrite each other. Always copy.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Binary Tree Paths | Easy | https://leetcode.com/problems/binary-tree-paths/ |
| 2 | Path Sum II | Medium | https://leetcode.com/problems/path-sum-ii/ |

**What to notice — Problem 1:** Build a string path (e.g., `"1->2->5"`). At a leaf, record the full string. Use `strconv.Itoa` to convert int to string, and `strings.Join` or manual concatenation with `"->"`.

**What to notice — Problem 2:** Same as the template above but also track `remaining` sum. Subtract `node.Val` from `remaining` at each node (preorder). At a leaf, record the path when `remaining == 0`. Use the copy pattern to avoid slice aliasing bugs.

---

## Day 4 — Graph DFS: Connected Components

**Pattern:** DFS on explicit graphs with a visited set  
**Key insight:** Unlike trees, graphs have cycles — the visited set is mandatory or DFS loops forever.

**Template adaptation:**

```go
func dfs(node int, graph [][]int, visited []bool) {
    visited[node] = true // mark before recursing

    for _, neighbor := range graph[node] {
        if !visited[neighbor] {
            dfs(neighbor, graph, visited)
        }
    }
}

// Count connected components: run DFS from every unvisited node.
func countComponents(n int, graph [][]int) int {
    visited := make([]bool, n)
    count := 0
    for i := 0; i < n; i++ {
        if !visited[i] {
            dfs(i, graph, visited)
            count++
        }
    }
    return count
}
```

**Difference from tree DFS:** Trees can't have cycles so no visited set is needed. Graphs can — always mark visited before recursing, not after.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Islands | Medium | https://leetcode.com/problems/number-of-islands/ |
| 2 | Number of Provinces | Medium | https://leetcode.com/problems/number-of-provinces/ |

**What to notice — Problem 1:** You solved this with BFS in Week 1. Now solve it with DFS and compare. The DFS version is often shorter — no queue, just recurse into all 4 neighbors. Mark visited by flipping `'1'` to `'0'` in-place.

**What to notice — Problem 2:** Input is an `n×n` adjacency matrix. `isConnected[i][j] == 1` means edge between `i` and `j`. Build the neighbor list on the fly: iterate `j` from 0 to n, add `j` if `isConnected[i][j] == 1 && i != j`.

---

## Day 5 — Grid DFS: Flood Fill & Area

**Pattern:** DFS on a 2D grid treating cells as graph nodes  
**Key insight:** `getNeighbors` returns the 4 adjacent cells. Mark visited by modifying the grid in-place — no separate visited set needed.

**Template adaptation:**

```go
var dirs = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func dfs(grid [][]int, r, c, target, replacement int) {
    if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
        return // out of bounds
    }
    if grid[r][c] != target {
        return // wrong color or already visited
    }

    grid[r][c] = replacement // mark visited in-place

    for _, d := range dirs {
        dfs(grid, r+d[0], c+d[1], target, replacement)
    }
}
```

**In-place marking:** Setting `grid[r][c] = replacement` before recursing ensures you never visit the same cell twice. This replaces the visited set from Day 4.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Flood Fill | Easy | https://leetcode.com/problems/flood-fill/ |
| 2 | Max Area of Island | Medium | https://leetcode.com/problems/max-area-of-island/ |

**What to notice — Problem 1:** Edge case — if `image[sr][sc]` already equals `color`, return immediately. When `target == replacement`, the in-place mark (`grid[r][c] = replacement`) is a no-op — the cell keeps its original value, the base case (`grid[r][c] != target`) never triggers, and DFS recurses into the same cell forever.

**What to notice — Problem 2:** DFS returns an `int` (the area of the island it just explored). Each cell contributes 1 to the area. Return `1 + dfs(up) + dfs(down) + dfs(left) + dfs(right)`. Track the max across all starting cells.

---

## Day 6 — Backtracking

**Pattern:** DFS with explicit state undo on the way back up  
**Key insight:** Make a choice, recurse, then undo the choice. The undo step is what separates backtracking from plain DFS.

**Template adaptation:**

```go
func backtrack(candidates []int, current []int, result *[][]int) {
    // base case: record result when done
    if isDone(current) {
        tmp := make([]int, len(current))
        copy(tmp, current)
        *result = append(*result, tmp)
        return
    }

    for _, candidate := range candidates {
        if isValid(candidate, current) {
            current = append(current, candidate) // make choice
            backtrack(candidates, current, result)
            current = current[:len(current)-1]   // undo choice
        }
    }
}
```

The undo step (`current[:len(current)-1]`) restores `current` to its state before the choice. Without it, you'd accumulate choices across branches.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Permutations | Medium | https://leetcode.com/problems/permutations/ |
| 2 | Subsets | Medium | https://leetcode.com/problems/subsets/ |

**What to notice — Problem 1:** Track which numbers are used with a `used []bool` slice. `isValid` checks `!used[i]`. Make choice: `used[i] = true`, append to current. Undo: `used[i] = false`, remove last from current. Record when `len(current) == len(nums)`.

**What to notice — Problem 2:** At each index, choose to include or exclude. No undo needed — pass index+1 to recurse, so each branch naturally excludes the current element. Record `current` at every call (not just leaves) since all subsets are valid.

---

## Day 7 — Iterative DFS + Mixed Review

**Pattern:** Iterative DFS using an explicit stack  
**Key insight:** Replace the call stack with an explicit stack slice. Pop from the end (LIFO). Push right before left so left is processed first.

**The Iterative Template:**

```go
func dfs(root *TreeNode) {
    if root == nil {
        return
    }
    stack := []*TreeNode{root}

    for len(stack) > 0 {
        node := stack[len(stack)-1]   // peek top
        stack = stack[:len(stack)-1]  // pop

        // process node here

        if node.Right != nil {
            stack = append(stack, node.Right) // push right first
        }
        if node.Left != nil {
            stack = append(stack, node.Left)  // push left last (processed first)
        }
    }
}
```

**BFS vs iterative DFS — the only difference:**

```go
// BFS: queue — pop from front
node := queue[0]
queue = queue[1:]

// DFS: stack — pop from end
node := stack[len(stack)-1]
stack = stack[:len(stack)-1]
```

**When to use iterative over recursive:**
- Very deep trees/graphs where recursion risks stack overflow (depth > ~10,000)
- When you need to pause/resume traversal mid-way
- Competitive environments with strict stack size limits

**No new problems today.** Do this instead:

1. Write the iterative DFS template from memory.
2. Pick the 1 day that felt hardest (Days 2–6). Redo one problem from that day without looking at your notes.
3. Answer out loud: *When would you use BFS over DFS, and vice versa?*

   > Use **BFS** when: shortest path matters, or you need level-by-level processing. Use **DFS** when: you need to explore all paths, detect cycles, do backtracking, or the problem is naturally recursive (trees, nested structures).

---

## BFS vs DFS Cheatsheet

| | BFS | DFS |
|-|-----|-----|
| Data structure | Queue (FIFO) | Call stack / explicit stack (LIFO) |
| Traversal order | Level by level | Depth first |
| Finds shortest path? | Yes (unweighted) | No |
| Memory usage | O(width of tree/graph) | O(depth of tree/graph) |
| Best for | Shortest path, level problems, multi-source | Path existence, backtracking, connected components, tree problems |

## Pattern Cheatsheet

| Pattern | Visited needed? | Key move | Level tracking? |
|---------|----------------|----------|----------------|
| Tree DFS (depth/existence) | No | Return value up the stack | No |
| Tree DFS (paths) | No | Pass path down, copy at leaf | No |
| Graph DFS | Yes | Mark visited before recursing | No |
| Grid DFS | In-place | Modify grid cell to mark visited | No |
| Backtracking | Depends | Undo choice after recursing | No |
| Iterative DFS | Yes (for graphs) | Pop from end of stack slice | No |

---

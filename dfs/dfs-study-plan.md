# DFS Interview Prep — 7-Day Study Plan

**Level:** Beginner (completed BFS week)  
**Goal:** Write recursive and iterative DFS from memory, recognize patterns, implement backtracking.

---

## What Is DFS?

**Depth-First Search (DFS)** explores by going **as deep as possible** along one path before backing up and trying the next one. It's naturally written with **recursion** (the call stack remembers where to come back to), or with an explicit **stack**.

DFS only keeps the current path in memory, which makes it the natural tool for "explore every path" and "combine answers from children" problems. Unlike BFS, it does **not** find shortest paths.

**Typical problems it solves:**

- Tree problems: depth, path sums, lowest common ancestor, validating structure
- Root-to-leaf paths and collecting every path
- Connected components, flood fill, counting islands, area of a region
- Cycle detection and reachability in graphs
- Backtracking: generating permutations, subsets, and combinations

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

**What DFS buys you:** DFS follows one path all the way down before backing up to try the next. That makes it the natural fit for "explore every path / every choice" problems, and it only needs memory for the *current* path — O(depth) — where BFS needs memory for a whole level.

**Complexity:** O(n) for a tree, O(V + E) for a graph — every node is entered once, every edge followed once. Space is O(h), the deepest recursion: `log n` for a balanced tree, but up to `n` for a skewed tree or a long chain in a graph.

**Run it:** every template in this file is a real function in [`dfs.go`](dfs.go), tested in [`dfs_test.go`](dfs_test.go). From the repo root: `go test ./dfs/ -v`.

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

**Dry run** — the call trace of the tree template on this tree. Indentation = recursion depth (= the call stack).

```
      1
     / \
    2   3
   / \
  4   5
```

```
enter 1
    enter 2
        enter 4
        exit  4
        enter 5
        exit  5
    exit  2
    enter 3
    exit  3
exit  1
```

Read the `enter` lines top to bottom: `1 2 4 5 3` — that's **preorder**. Read the `exit` lines: `4 5 2 3 1` — that's **postorder**. Same walk, same calls; the only difference is which line you put your processing on. At any moment the "stack" is just the chain of `enter`s that haven't exited yet — e.g. while inside `4` it's `1 → 2 → 4`.

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
```

`max` and `min` are Go builtins since Go 1.21 — no helper needed.

**Why it works — "trust the recursion":** For postorder problems, don't trace the whole tree in your head. Assume `maxDepth` already returns the right answer for the left and right subtrees (they're smaller, so by induction it does), and only ask: *given those two answers, what's the answer for this node?* Here: `1 + max(left, right)`. Then check the base case (`nil` → 0). If both are right, the whole function is right.

**Dry run** — `maxDepth` on this tree. Rows appear in the order calls *return* (postorder):

```
        5
       / \
      4   8
     /   / \
    11  13  4
   /  \      \
  7    2      1
```

| Node | `left` | `right` | returns `1 + max` |
|---|---|---|---|
| `7` | 0 | 0 | **1** |
| `2` | 0 | 0 | **1** |
| `11` | 1 | 1 | **2** |
| `4` | 2 | 0 | **3** |
| `13` | 0 | 0 | **1** |
| `1` | 0 | 0 | **1** |
| `4` | 0 | 1 | **2** |
| `8` | 1 | 2 | **3** |
| `5` | 3 | 3 | **4** |

`hasPathSum(root, 22)` goes the other way (preorder, passing state down): `remaining` is `22 → 17 → 13 → 2 → 0` along `5 → 4 → 11 → 2`, and `2` is a leaf with `remaining == 0`, so `true` bubbles back up through the `||`.

**Run it:** `maxDepth` and `hasPathSum` in [`dfs.go`](dfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Maximum Depth of Binary Tree | Easy | https://leetcode.com/problems/maximum-depth-of-binary-tree/ |
| 2 | Path Sum | Easy | https://leetcode.com/problems/path-sum/ |
| 3 | Lowest Common Ancestor of a Binary Tree | Medium | https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/ |

**What to notice — Problem 1:** This is postorder — you need both children's depths before you can compute the current node's depth. Return `1 + max(left, right)`. Base case: `nil` returns 0.

**What to notice — Problem 2:** This is preorder — subtract `node.Val` from `remaining` as you go down. At a leaf, check if `remaining == 0`. Don't return true at a nil node — only at a leaf (both children nil).

**What to notice — Problem 3 (LCA):** Postorder — recurse left and right first, then decide at the current node. If the current node is `p` or `q`, return it. If left and right both return non-nil, the current node is the LCA. If only one side returns non-nil, bubble that up.

**Complexity:** all three problems are O(n) time (each node visited once) and O(h) space for the recursion.

**Common mistakes:**

- Path Sum: returning `remaining == 0` at a `nil` node instead of at a leaf. For the tree `1 → 2` (one child) and target `1`, that wrongly returns `true` — `1` alone isn't a root-to-leaf path.
- LCA: forgetting that a node can be its own ancestor. If `p` is above `q`, the answer is `p` itself — the "if current is `p` or `q`, return it" line handles this, so don't recurse past it.
- `hasPathSum` short-circuits: with `||`, the right subtree is never searched once the left returns `true`. That's intended, but remember it when you add side effects.

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

**Why the copy matters — a real run:** this is the template *without* the copy (`*result = append(*result, path)`), run on the Day 2 tree:

```
[[5 4 11 2] [5 4 11 2] [5 8 13] [5 8 4 1]]
```

The first path should be `[5 4 11 7]`. Here's what happened: after appending `5`, `4`, `11`, Go had grown `path`'s backing array to capacity 4. Appending `7` used the spare slot (no new array), and `result` stored a slice pointing at that array. Back up at `11`, appending `2` wrote into **the same slot**, overwriting the `7` that `result` was still pointing at. The right subtree survived only because its appends happened to trigger fresh allocations.

That's what makes this bug nasty: whether it bites depends on spare capacity, so it passes some tests and fails others. Always copy before storing.

**Run it:** `collectPaths` in [`dfs.go`](dfs.go) — the test checks all four paths.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Binary Tree Paths | Easy | https://leetcode.com/problems/binary-tree-paths/ |
| 2 | Path Sum II | Medium | https://leetcode.com/problems/path-sum-ii/ |

**What to notice — Problem 1:** Build a string path (e.g., `"1->2->5"`). At a leaf, record the full string. Use `strconv.Itoa` to convert int to string, and `strings.Join` or manual concatenation with `"->"`.

**What to notice — Problem 2:** Same as the template above but also track `remaining` sum. Subtract `node.Val` from `remaining` at each node (preorder). At a leaf, record the path when `remaining == 0`. Use the copy pattern to avoid slice aliasing bugs.

**Complexity:** O(n · h) time — every node is visited once, but each of up to n/2 leaves copies a path of length h. Space O(h) for the recursion and `path`, not counting the output.

**Common mistakes:**

- Storing `path` without copying (see the run above).
- Recording at `nil` instead of at a leaf: each path then gets recorded twice (once from each `nil` child).
- Path Sum II: pruning when `remaining < 0`. Node values can be negative, so a negative `remaining` can still come back to 0 further down.

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

**Why it works:** Marking `visited[node] = true` *before* recursing means a node is marked the moment it's entered. When a cycle leads back to a node that's still on the call stack, that node is already marked, so the recursion stops instead of looping forever. Every node gets exactly one call.

**Dry run** — `dfsGraph(0, …)` on the Day 1 BFS graph (`0–1, 0–2, 1–3, 2–3`):

```
dfs(0)  visited=[0]
    dfs(1)  visited=[0, 1]
        skip 0 (visited)
        dfs(3)  visited=[0, 1, 3]
            skip 1 (visited)
            dfs(2)  visited=[0, 1, 2, 3]
                skip 0 (visited)
                skip 3 (visited)
    skip 2 (visited)
```

Compare to BFS's order `0, 1, 2, 3`: DFS reached `2` *last*, via the long way round `0 → 1 → 3 → 2`. Same set of nodes, different order — which is why DFS can't give shortest paths.

**Run it:** `dfsGraph` and `countComponents` in [`dfs.go`](dfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Number of Islands | Medium | https://leetcode.com/problems/number-of-islands/ |
| 2 | Number of Provinces | Medium | https://leetcode.com/problems/number-of-provinces/ |

**What to notice — Problem 1:** You solved this with BFS in Week 1. Now solve it with DFS and compare. The DFS version is often shorter — no queue, just recurse into all 4 neighbors. Mark visited by flipping `'1'` to `'0'` in-place.

**What to notice — Problem 2:** Input is an `n×n` adjacency matrix. `isConnected[i][j] == 1` means edge between `i` and `j`. Build the neighbor list on the fly: iterate `j` from 0 to n, add `j` if `isConnected[i][j] == 1 && i != j`.

**Complexity:** O(V + E) time. Recursion depth is O(V) in the worst case (a long chain). For Number of Islands that's O(R·C) depth on an all-land grid — fine in Go, because goroutine stacks grow on demand (up to 1 GB by default). In Python or Java, the same code can overflow the stack.

**Common mistakes:**

- Marking visited *after* the recursive calls — on a cycle you recurse forever.
- Provinces: using `j` from `0` and forgetting that `isConnected[i][i] == 1` — harmless here because `i` is already visited, but know why it's harmless.

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

**Why in-place marking works:** `target` means "not visited yet" and `replacement` means "visited". The moment a cell is overwritten it stops matching `target`, so the base case rejects it on any later visit. This is the visited set from Day 4, stored inside the grid itself — which is exactly why it breaks when `target == replacement`.

**Dry run** — `fill(grid, 0, 0, 1, 2)`, neighbors tried in the order right, left, down, up:

```
before      after
1 1 0       2 2 0
1 0 1       2 0 1
0 1 1       0 1 1
```

Cells filled, in order: `(0,0) → (0,1) → (1,0)`. The bottom-right group of `1`s is untouched: it only touches the filled region diagonally, and diagonals aren't neighbors.

**Run it:** `fill` in [`dfs.go`](dfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Flood Fill | Easy | https://leetcode.com/problems/flood-fill/ |
| 2 | Max Area of Island | Medium | https://leetcode.com/problems/max-area-of-island/ |

**What to notice — Problem 1:** Edge case — if `image[sr][sc]` already equals `color`, return immediately. When `target == replacement`, the in-place mark (`grid[r][c] = replacement`) is a no-op — the cell keeps its original value, the base case (`grid[r][c] != target`) never triggers, and DFS recurses into the same cell forever.

**What to notice — Problem 2:** DFS returns an `int` (the area of the island it just explored). Each cell contributes 1 to the area. Return `1 + dfs(up) + dfs(down) + dfs(left) + dfs(right)`. Track the max across all starting cells.

**Complexity:** O(R·C) time — each cell is filled at most once and checked at most 4 times. Recursion depth is O(R·C) in the worst case (a snake-shaped region).

**Common mistakes:**

- Flood Fill: not returning early when the start pixel already has the new color (infinite recursion, explained above).
- Max Area: forgetting `max` across all starting cells, or counting the area with a global counter you never reset between islands.
- Mutating the input when the problem says not to. Use a separate `visited [][]bool` in that case.

---

## Day 6 — Backtracking

**Pattern:** DFS with explicit state undo on the way back up  
**Key insight:** Make a choice, recurse, then undo the choice. The undo step is what separates backtracking from plain DFS.

**Core idea:** Backtracking explores a decision tree where each node is a partial solution, built incrementally via three steps at every level: **choose** a candidate, **explore** by recursing with it applied, **un-choose** by undoing it before trying the next candidate. This is what makes it DFS rather than brute-force enumeration — you generate and check one full candidate at a time, on a shared, reused state, instead of materializing every possibility upfront.

The other lever is pruning: an `isValid` check before recursing cuts off entire subtrees that can't lead to a valid solution, so you never even explore them. Weak pruning degenerates into brute force with extra overhead — the value of backtracking comes almost entirely from how early and how tight that check is.

**Generic pseudocode:**

```
backtrack(state):
    if state is a complete solution:
        record/return it
        return
    for each candidate in next_candidates(state):
        if is_valid(candidate, state):       // pruning
            apply(candidate, state)
            backtrack(state)
            undo(candidate, state)           // the "backtrack" step
```

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

**Dry run** — the template with `isDone = "picked k"` and `isValid = "bigger than the last pick"`, choosing 2 from `[1, 2, 3]`:

```
choose 1 -> [1]
    prune  1 (not > 1)
    choose 2 -> [1, 2]
        [1, 2]  record
    undo   2 -> [1]
    choose 3 -> [1, 3]
        [1, 3]  record
    undo   3 -> [1]
undo   1 -> []
choose 2 -> [2]
    prune  1 (not > 2)
    prune  2 (not > 2)
    choose 3 -> [2, 3]
        [2, 3]  record
    undo   3 -> [2]
undo   2 -> []
choose 3 -> [3]
    prune  1 (not > 3)
    prune  2 (not > 3)
    prune  3 (not > 3)
undo   3 -> []
```

Two things to see:

- **Every `undo` puts `current` back exactly as the parent had it**, so the next sibling starts from the same state. That's what lets one shared slice serve the whole tree.
- **The `[3]` branch is wasted work.** It can never reach length 2, but the check only rejects candidates one at a time. A tighter check — "are there enough candidates left to finish?" — would cut the whole branch before entering it. That's what "pruning quality decides the speed" means.

**Run it:** `combine` / `backtrack` in [`dfs.go`](dfs.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Permutations | Medium | https://leetcode.com/problems/permutations/ |
| 2 | Subsets | Medium | https://leetcode.com/problems/subsets/ |

**What to notice — Problem 1:** Track which numbers are used with a `used []bool` slice. `isValid` checks `!used[i]`. Make choice: `used[i] = true`, append to current. Undo: `used[i] = false`, remove last from current. Record when `len(current) == len(nums)`.

**What to notice — Problem 2:** No `used` array. Pass a `start` index: loop `i` from `start` to the end, append `nums[i]`, recurse with `start = i+1`, then undo. Recursing with `i+1` means you only ever pick elements to the right of the last pick, so `[1,2]` and `[2,1]` can't both appear. Record a copy of `current` at every call (not just leaves) — every partial pick is a valid subset. The undo is still there: either truncate explicitly, or pass `current` by value so the caller's length never changes (the same slice-header trick as Day 3).

**Complexity:**

| Problem | Time | Why |
|---|---|---|
| Permutations | O(n · n!) | n! complete permutations, each copied in O(n) |
| Subsets | O(n · 2ⁿ) | 2ⁿ subsets, each copied in O(n) |

Space is O(n) for the recursion and `current`, not counting the output.

**Common mistakes:**

- Forgetting to copy `current` when recording — the Day 3 aliasing bug, but guaranteed to bite here because `current` is reused for every branch.
- Permutations: setting `used[i] = true` but forgetting `used[i] = false` in the undo step.
- Subsets: recording only at leaves (`len(current) == len(nums)`) — you'd return just the full set.

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

**Dry run** on the Day 1 tree (`1` with children `2`, `3`; `2` with children `4`, `5`):

| Pop | Push (right first, then left) | Stack after (top is last) |
|---|---|---|
| `1` | 3, 2 | `[3, 2]` |
| `2` | 5, 4 | `[3, 5, 4]` |
| `4` | — | `[3, 5]` |
| `5` | — | `[3]` |
| `3` | — | `[]` |

Pop order `1 2 4 5 3` — identical to the recursive preorder from Day 1.

**For graphs, mark visited when you pop, not when you push.** Marking on push (like BFS) still visits every node, but no longer in true DFS order. The safe iterative graph version is: pop, `if visited[node] { continue }`, mark, push unvisited neighbors.

**Run it:** `dfsIterative` in [`dfs.go`](dfs.go) — the test checks it matches recursive preorder.

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

| Pattern | Visited needed? | Key move | Time |
|---------|----------------|----------|------|
| Tree DFS (depth/existence) | No | Return value up the stack | O(n) |
| Tree DFS (paths) | No | Pass path down, copy at leaf | O(n · h) |
| Graph DFS | Yes | Mark visited before recursing | O(V + E) |
| Grid DFS | In-place | Modify grid cell to mark visited | O(R·C) |
| Backtracking | Depends | Undo choice after recursing | O(output size × copy cost) |
| Iterative DFS | Yes (for graphs, mark on pop) | Pop from end of stack slice | same as recursive |

---

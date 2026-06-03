# DFS Interview Prep — 7-Day Study Plan

**Level:** Beginner (completed BFS week)  
**Goal:** Write recursive and iterative DFS from memory, recognize patterns, implement backtracking.

---

## The Recursive Template

Every DFS problem (Days 1–6) is a variation of this. Learn it until you can write it in 60 seconds.

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

**What to notice — Problem 1:** This is postorder — you need both children's depths before you can compute the current node's depth. Return `1 + max(left, right)`. Base case: `nil` returns 0.

**What to notice — Problem 2:** This is preorder — subtract `node.Val` from `remaining` as you go down. At a leaf, check if `remaining == 0`. Don't return true at a nil node — only at a leaf (both children nil).

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

    // no explicit removal needed — append returns a new slice header
    // the caller's path is unchanged when this call returns
}
```

**Why copy the path at the leaf:** Go slices share underlying arrays. If you append `path` directly to `result`, all entries in `result` will point to the same array and overwrite each other. Always copy.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Binary Tree Paths | Easy | https://leetcode.com/problems/binary-tree-paths/ |
| 2 | Path Sum II | Medium | https://leetcode.com/problems/path-sum-ii/ |

**What to notice — Problem 1:** Build a string path (e.g., `"1->2->5"`). At a leaf, record the full string. Use `strconv.Itoa` to convert int to string, and `strings.Join` or manual concatenation with `"->"`.

**What to notice — Problem 2:** Same as the template above but also track `remaining` sum. Only record the path at a leaf when `remaining == node.Val`. Use the copy pattern to avoid slice aliasing bugs.

---

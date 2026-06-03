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

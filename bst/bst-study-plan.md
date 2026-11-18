# BST & Inorder Traversal Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 14 — builds on DFS week's tree recursion and iterative stack)  
**Goal:** Use the BST ordering property to skip half the tree at every step, write iterative inorder from memory, and treat "inorder = sorted" as a tool.

---

## What Is a Binary Search Tree?

A **Binary Search Tree (BST)** is a binary tree with an ordering rule: every value in a node's **left** subtree is smaller than the node, and every value in its **right** subtree is larger. That rule lets you decide which way to go with one comparison, so search and insert take O(h) — O(log n) when the tree is balanced.

**Inorder traversal** (left, node, right) visits a BST's values in **sorted order** — which turns many BST problems into "walk the sorted sequence".

**Typical problems it solves:**

- Search, insert, and delete
- Validating that a tree really is a BST
- K-th smallest element, minimum difference between values
- Range sums and trimming to a range
- Lowest common ancestor in a BST
- Building a balanced BST from sorted data, fixing two swapped nodes

---

## The Core Templates

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

### Walk One Path (search / insert / LCA)

```go
node := root
for node != nil && node.Val != val {
    if val < node.Val {
        node = node.Left // everything smaller lives on the left
    } else {
        node = node.Right
    }
}
```

### Iterative Inorder (sorted order, one node at a time)

```go
var stack []*TreeNode
node := root
for node != nil || len(stack) > 0 {
    for node != nil { // go as far left as possible
        stack = append(stack, node)
        node = node.Left
    }
    node = stack[len(stack)-1] // the smallest unvisited value
    stack = stack[:len(stack)-1]

    // process node.Val here — values arrive in ascending order

    node = node.Right // then do the same for its right subtree
}
```

**Three rules — never break them:**

1. **The BST property is about whole subtrees, not just children.** *Every* value in the left subtree is smaller than the node, and *every* value in the right subtree is larger. Checking only `left.Val < node.Val < right.Val` is the classic wrong validation.
2. **One comparison picks one side.** Search, insert, floor/ceiling and LCA all walk a single root-to-leaf path: O(h). If your BST code visits both children unconditionally, you're not using the BST property.
3. **Inorder = sorted.** Whenever a BST problem mentions "k-th smallest", "closest", "adjacent difference", or "two values that are out of order", think "walk the inorder sequence and remember the previous value".

**Brute force → BST property:** without the ordering you'd search the whole tree, O(n). With it, every comparison discards an entire subtree, so search/insert are O(h) — O(log n) when the tree is balanced. Beware: a BST built from sorted input degenerates into a linked list and h becomes n. That's why Day 5 builds balanced trees.

**Complexity:** search/insert/delete O(h); full inorder O(n) time, O(h) stack space.

**Run it:** every template in this file is a real function in [`bst.go`](bst.go), tested in [`bst_test.go`](bst_test.go). From the repo root: `go test ./bst/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** A BST is a sorted array folded into a tree. Walking a path is binary search; inorder traversal unfolds it back into the sorted array.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does the iterative inorder push everything on the way left before popping?*

   > The smallest value in any subtree is its leftmost node, and every node on the way down to it is a larger value that must come *after* its left subtree. The stack holds exactly those "waiting" ancestors, in the order they'll be output. Popping one means its entire left subtree is done; then its right subtree is handled the same way.

**Dry run** — iterative inorder on the tree built by inserting `4, 2, 6, 1, 3, 5, 7`:

```
        4
      /   \
     2     6
    / \   / \
   1   3 5   7
```

| Step | Pushed going left | Pop (output) | Then go to | Stack after (bottom→top) |
|---|---|---|---|---|
| 1 | 4, 2, 1 | 1 | nil (right) | `[4, 2]` |
| 2 | — | 2 | right child 3 | `[4]` |
| 3 | 3 | 3 | nil (right) | `[4]` |
| 4 | — | 4 | right child 6 | `[]` |
| 5 | 6, 5 | 5 | nil (right) | `[6]` |
| 6 | — | 6 | right child 7 | `[]` |
| 7 | 7 | 7 | nil (right) | `[]` |

Output `1 2 3 4 5 6 7`. The stack never holds more than h = 3 nodes.

---

## Day 2 — Search, Insert, Validate

**Pattern:** Walk one path to search or insert; pass bounds down to validate  
**Key insight:** Each node doesn't just have a parent — it has a **valid range** `(lo, hi)` inherited from *all* its ancestors. Going left tightens `hi`; going right tightens `lo`.

**Template adaptation:**

```go
// insert as a new leaf; returns the (possibly new) root
func insert(root *TreeNode, val int) *TreeNode {
    newNode := &TreeNode{Val: val}
    if root == nil {
        return newNode
    }
    node := root
    for {
        if val < node.Val {
            if node.Left == nil {
                node.Left = newNode
                return root
            }
            node = node.Left
        } else {
            if node.Right == nil {
                node.Right = newNode
                return root
            }
            node = node.Right
        }
    }
}

func isValidBST(root *TreeNode) bool {
    var valid func(node *TreeNode, lo, hi int) bool
    valid = func(node *TreeNode, lo, hi int) bool {
        if node == nil {
            return true
        }
        if node.Val <= lo || node.Val >= hi {
            return false
        }
        return valid(node.Left, lo, node.Val) && valid(node.Right, node.Val, hi)
    }
    return valid(root, math.MinInt, math.MaxInt)
}
```

**Why bounds and not parent checks — the trap tree:**

```
    5
   / \
  1   6
     / \
    3   7
```

Every parent/child pair looks fine (`1 < 5 < 6`, `3 < 6 < 7`), but `3` is in 5's **right** subtree, so it must be `> 5`. With bounds:

| Node | Range it must be in | OK? |
|---|---|---|
| 5 | (−∞, +∞) | yes |
| 1 | (−∞, 5) | yes |
| 6 | (5, +∞) | yes |
| 3 | (5, 6) | **no** — 3 ≤ 5 |

**Run it:** `search`, `insert`, and `isValidBST` in [`bst.go`](bst.go) — the test includes this exact trap tree.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Validate Binary Search Tree | Medium | https://leetcode.com/problems/validate-binary-search-tree/ |
| 2 | Delete Node in a BST | Medium | https://leetcode.com/problems/delete-node-in-a-bst/ |

**What to notice — Problem 1:** The template. An equally good alternative: iterative inorder and check every value is strictly greater than the previous one.

**What to notice — Problem 2:** Three cases once you find the node: no children → remove it; one child → replace it with that child; two children → copy in the **inorder successor's** value (leftmost node of the right subtree), then delete that successor from the right subtree (it has at most one child, so it's an easy case). Recursive `node.Left = deleteNode(node.Left, key)` keeps the rewiring simple.

**Complexity:** both O(h) for delete; validation O(n) (it must check every node).

**Common mistakes:**

- Validating with parent/child comparisons only (the trap tree).
- Initializing bounds with `math.MinInt32` / `math.MaxInt32` — a node holding exactly that value is wrongly rejected. Use `math.MinInt`/`math.MaxInt` (or `nil` pointers).
- Delete: forgetting to return the new subtree root to the caller.

---

## Day 3 — Inorder as a Sorted Stream

**Pattern:** Iterative inorder, stop early or compare with the previous value  
**Key insight:** You rarely need the whole sorted array. The iterative version produces it one value at a time, so you can stop at the k-th value or track adjacent pairs in O(1) extra space.

**Template adaptation:**

```go
func inorderIterative(root *TreeNode) []int {
    var out []int
    var stack []*TreeNode
    node := root
    for node != nil || len(stack) > 0 {
        for node != nil {
            stack = append(stack, node)
            node = node.Left
        }
        node = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        out = append(out, node.Val)
        node = node.Right
    }
    return out
}
```

**Why stopping early is safe:** values are produced strictly in ascending order (Day 1's dry run), so the k-th value produced is the k-th smallest. Nothing produced later can be smaller.

**Run it:** `inorderIterative` in [`bst.go`](bst.go) (dry run on Day 1).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Kth Smallest Element in a BST | Medium | https://leetcode.com/problems/kth-smallest-element-in-a-bst/ |
| 2 | Minimum Absolute Difference in BST | Easy | https://leetcode.com/problems/minimum-absolute-difference-in-bst/ |

**What to notice — Problem 1:** Replace `out = append(...)` with `k--; if k == 0 { return node.Val }`. O(h + k) instead of O(n).

**What to notice — Problem 2:** In a sorted sequence, the minimum difference is always between **neighbors**. Keep `prev` and compute `node.Val - prev` at each step — no need to compare all pairs.

**Complexity:** Problem 1 O(h + k) time, O(h) space. Problem 2 O(n) time, O(h) space.

**Common mistakes:**

- Collecting the full inorder slice and then indexing it — correct, but O(n) space and time when O(h + k) is possible.
- Problem 2: initializing `prev` to 0 — the first difference is then meaningless. Use a "have I seen a previous value" flag, or a pointer.

---

## Day 4 — Pruning with the Ordering

**Pattern:** Like DFS, but skip a whole subtree when the ordering proves it's irrelevant  
**Key insight:** If `node.Val < lo`, the node's entire left subtree is even smaller — skip both and go right. That single check is what makes BST range queries fast.

**Template adaptation:**

```go
// sum of all values in [lo, hi]
func rangeSum(node *TreeNode, lo, hi int) int {
    if node == nil {
        return 0
    }
    if node.Val < lo {
        return rangeSum(node.Right, lo, hi) // node and its left subtree are too small
    }
    if node.Val > hi {
        return rangeSum(node.Left, lo, hi) // node and its right subtree are too big
    }
    return node.Val + rangeSum(node.Left, lo, hi) + rangeSum(node.Right, lo, hi)
}
```

**Dry run** — tree from inserting `10, 5, 15, 3, 7, 18`; `rangeSum(root, 7, 15)`:

```
        10
       /  \
      5    15
     / \     \
    3   7     18
```

```
10 in range: add, search both sides
    5 < 7: skip it and its LEFT subtree, go right
        7 in range: add, search both sides
    15 in range: add, search both sides
        18 > 15: skip it and its RIGHT subtree, go left
```

Total `10 + 7 + 15 = 32`. Node `3` was never visited — the check at `5` ruled it out.

**Run it:** `rangeSum` in [`bst.go`](bst.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Lowest Common Ancestor of a Binary Search Tree | Medium | https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/ |
| 2 | Trim a Binary Search Tree | Medium | https://leetcode.com/problems/trim-a-binary-search-tree/ |

**What to notice — Problem 1:** Much simpler than the general-tree LCA from DFS week. Walk one path: if both `p` and `q` are smaller, go left; both larger, go right; otherwise the current node is where they split — that's the LCA. O(h), no recursion needed.

**What to notice — Problem 2:** Same three cases as `rangeSum`, but returning nodes: too small → return `trim(node.Right)` (the whole left side is gone); too big → return `trim(node.Left)`; in range → trim both children and keep the node.

**Complexity:** LCA O(h). Trim O(n) in the worst case (every node may be in range).

**Common mistakes:**

- LCA: using the general binary-tree LCA (visits everything) — correct but ignores the BST property.
- Trim: returning `nil` for an out-of-range node instead of its trimmed subtree — you lose valid values below it.

---

## Day 5 — Building Balanced Trees and Iterators

**Pattern:** Pick the middle as the root, recurse on each half; or turn iterative inorder into an iterator  
**Key insight:** The middle of a sorted slice splits it into two halves of (almost) equal size — exactly what a balanced BST root needs.

**Template adaptation:**

```go
func sortedArrayToBST(nums []int) *TreeNode {
    if len(nums) == 0 {
        return nil
    }
    mid := len(nums) / 2
    return &TreeNode{
        Val:   nums[mid],
        Left:  sortedArrayToBST(nums[:mid]),
        Right: sortedArrayToBST(nums[mid+1:]),
    }
}
```

**Why it's balanced and valid:** everything in `nums[:mid]` is smaller than `nums[mid]`, everything in `nums[mid+1:]` is larger — valid BST. The two halves differ in size by at most 1, at every level — so the height is ⌈log₂(n+1)⌉.

**Dry run** — `[-10, -3, 0, 5, 9, 12, 20]`:

| Call on | `mid` value | Left half | Right half |
|---|---|---|---|
| `[-10, -3, 0, 5, 9, 12, 20]` | 5 | `[-10, -3, 0]` | `[9, 12, 20]` |
| `[-10, -3, 0]` | -3 | `[-10]` | `[0]` |
| `[9, 12, 20]` | 12 | `[9]` | `[20]` |

```
         5
       /   \
     -3     12
     / \   /  \
   -10  0 9    20
```

7 nodes, height 3 — perfectly balanced. Inserting the same values in sorted order with `insert` would give height 7.

**Run it:** `sortedArrayToBST` in [`bst.go`](bst.go) — the test checks validity, inorder, and height.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Convert Sorted Array to Binary Search Tree | Easy | https://leetcode.com/problems/convert-sorted-array-to-binary-search-tree/ |
| 2 | Binary Search Tree Iterator | Medium | https://leetcode.com/problems/binary-search-tree-iterator/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Split the iterative inorder into pieces: the constructor does the first "push everything going left"; `Next()` pops one node, then pushes the left spine of its right child; `HasNext()` is `len(stack) > 0`. Each node is pushed and popped once, so `Next()` is O(1) amortized with O(h) memory.

**Complexity:** Problem 1 O(n). Problem 2 O(1) amortized per call, O(h) space.

**Common mistakes:**

- Problem 1: passing index bounds and getting the mid formula off by one — slicing (`nums[:mid]`, `nums[mid+1:]`) avoids it.
- Problem 2: flattening the whole tree into a slice in the constructor — works, but O(n) memory when the problem asks for O(h).

---

## Day 6 — Finding What's Out of Order

**Pattern:** Inorder with a `prev` pointer; every `prev.Val > node.Val` is an "inversion"  
**Key insight:** Swapping two values in a sorted sequence creates either **one** inversion (if they were neighbors) or **two** (if they weren't). The first inversion's `prev` and the last inversion's `node` are the swapped pair.

**Template adaptation:**

```go
// the two values that were swapped in an otherwise valid BST
func findSwapped(root *TreeNode) (first, second int) {
    var prev *TreeNode
    found := false

    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil {
            return
        }
        inorder(node.Left)
        if prev != nil && prev.Val > node.Val { // an inversion
            if !found {
                first = prev.Val // the bigger value that moved left
                found = true
            }
            second = node.Val // the smaller value that moved right (latest one wins)
        }
        prev = node
        inorder(node.Right)
    }
    inorder(root)
    return first, second
}
```

**Dry run** — the Day 1 tree with two values swapped; inversions in the inorder sequence:

| Swapped | Inorder sequence | Inversions | `first`, `second` |
|---|---|---|---|
| 3 and 4 (neighbors) | `1 2 4 3 5 6 7` | (4, 3) | 4, 3 |
| 1 and 7 (far apart) | `7 2 3 4 5 6 1` | (7, 2), (6, 1) | 7, 1 |

In the second row, `first` comes from the *first* inversion's left side and `second` from the *last* inversion's right side.

**Run it:** `findSwapped` in [`bst.go`](bst.go) — the test swaps values in a real tree and checks all three cases.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Recover Binary Search Tree | Medium | https://leetcode.com/problems/recover-binary-search-tree/ |
| 2 | Two Sum IV - Input is a BST | Easy | https://leetcode.com/problems/two-sum-iv-input-is-a-bst/ |

**What to notice — Problem 1:** The template, but keep the *nodes* (not values) and swap their `Val`s at the end.

**What to notice — Problem 2:** Inorder gives a sorted sequence, so the Two Pointer week's opposite-ends Two Sum works. For O(h) space, run two iterators at once: a normal inorder (ascending) and a reverse inorder (descending, push right spines) as `left` and `right`.

**Complexity:** Problem 1 O(n) time, O(h) space. Problem 2 O(n) time; O(n) space with a slice or a hash set, O(h) with two iterators.

**Common mistakes:**

- Recover: overwriting `first` on the second inversion — only `second` should update.
- Recover: assuming the swapped values are always neighbors (only one inversion).
- Two Sum IV: forgetting the two pointers must be different nodes (`left != right`).

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Every BST problem uses one of two superpowers: "one comparison picks a side" (O(h) paths, pruning) or "inorder is sorted" (stream with a `prev`).

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing code, say out loud: *Am I walking a path, pruning a DFS, or streaming the inorder sequence?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write iterative inorder from memory
- [ ] Why does validating only parent/child pairs fail? Draw the counterexample
- [ ] What are the three cases for deleting a node, and what replaces a node with two children?
- [ ] Why is BST LCA O(h) while general-tree LCA is O(n)?
- [ ] Why does picking the middle element build a balanced tree?
- [ ] After swapping two values, why are there one or two inversions — and which values are the swapped pair?

---

## Pattern Cheatsheet

| Pattern | Technique | Key move | Time | Space |
|---------|-----------|----------|------|-------|
| Search / insert / LCA | walk one path | `val < node.Val` → left, else right | O(h) | O(1) iterative |
| Validate | DFS with bounds | left gets `(lo, node)`, right gets `(node, hi)` | O(n) | O(h) |
| K-th smallest / neighbors | iterative inorder | stop at k, or compare with `prev` | O(h + k) / O(n) | O(h) |
| Range / trim | pruned DFS | `< lo` → only right; `> hi` → only left | O(h + output) | O(h) |
| Build balanced | divide at middle | root = `nums[len/2]` | O(n) | O(log n) |
| Fix swapped values | inorder with `prev` | first inversion's left, last inversion's right | O(n) | O(h) |

---

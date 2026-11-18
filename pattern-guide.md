# Which Pattern Fits? — A Cross-Topic Guide

Use this when you've read a problem and don't know where to start. Work through the three sections in order: the input size tells you the complexity you need, the input shape narrows the candidates, and the wording usually names the pattern outright.

---

## 1. Let the Constraints Pick the Complexity

The input size limit is a strong hint at the intended solution. A rough budget is about 10⁸ simple operations.

| Largest `n` | Target complexity | Usually means |
|-------------|-------------------|---------------|
| ≤ 10–12 | O(n!) | permutations ([Backtracking](backtracking/backtracking-study-plan.md), Day 4) |
| ≤ 20–25 | O(2ⁿ) | subsets / backtracking with pruning, bitmask DP |
| ≤ 100–500 | O(n³) | Floyd–Warshall ([Dijkstra](dijkstra/dijkstra-study-plan.md), Day 6), 3-index DP |
| ≤ 5,000 | O(n²) | 2D DP ([DP](dp/dp-study-plan.md)), all pairs with pruning |
| ≤ 10⁵–10⁶ | O(n log n) | sort + scan, heap, binary search, Dijkstra |
| ≥ 10⁶ | O(n) | two pointers, sliding window, prefix sums, monotonic stack, hashing |
| huge values (≤ 10⁹), small `n` | O(n log V) | binary search on the answer ([Binary Search](binary-search/binary-search-study-plan.md), Day 4) |

---

## 2. Let the Input Shape Narrow It Down

| The input is… | First patterns to consider |
|---------------|----------------------------|
| A **sorted** array | [Binary Search](binary-search/binary-search-study-plan.md), [Two Pointer](two-pointer/two-pointer-study-plan.md) (opposite ends) |
| An unsorted array, asking about **pairs** | [Hashing](prefix-sum/prefix-sum-study-plan.md) (complement lookup), or sort + [Two Pointer](two-pointer/two-pointer-study-plan.md) |
| An array, asking about **contiguous subarrays** | [Sliding Window](sliding-window/sliding-window-study-plan.md) (non-negative / monotonic), [Prefix Sum + hash map](prefix-sum/prefix-sum-study-plan.md) (negatives, "== k"), [Greedy](greedy/greedy-study-plan.md) (Kadane) |
| An array, asking about the **nearest greater/smaller** element | [Monotonic Stack](monotonic-stack/monotonic-stack-study-plan.md) |
| A list of **[start, end]** pairs | [Intervals](intervals/intervals-study-plan.md) (sort by start or end), sweep line |
| A **linked list** | [Linked List](linked-list/linked-list-study-plan.md) (dummy head, fast/slow, reversal) |
| A **binary tree** | [DFS](dfs/dfs-study-plan.md) (recursion), [BFS](bfs/bfs-study-plan.md) (levels) |
| A **binary search tree** | [BST](bst/bst-study-plan.md) (walk one path, inorder = sorted) |
| A **grid** | [BFS](bfs/bfs-study-plan.md) (shortest steps), [DFS](dfs/dfs-study-plan.md) (regions), [DP](dp/dp-study-plan.md) (right/down only), [Dijkstra](dijkstra/dijkstra-study-plan.md) (weighted, any direction) |
| A **graph** with unweighted edges | [BFS](bfs/bfs-study-plan.md), [DFS](dfs/dfs-study-plan.md), [Union-Find](union-find/union-find-study-plan.md) |
| A **graph** with weighted edges | [Dijkstra](dijkstra/dijkstra-study-plan.md) (and its ladder), [Union-Find](union-find/union-find-study-plan.md) (Kruskal MST) |
| **Dependencies** / prerequisites / "must come before" | [Topological Sort](topological-sort/topological-sort-study-plan.md) |
| Many **words** / prefixes | [Trie](trie/trie-study-plan.md), [Hashing](prefix-sum/prefix-sum-study-plan.md) (canonical key) |
| A **stream** of numbers | [Heap](heap/heap-study-plan.md) (top k, running median) |

---

## 3. Let the Wording Name It

| The problem says… | Pattern | Where |
|-------------------|---------|-------|
| "shortest path", "minimum steps", "fewest moves" (unweighted) | BFS | [BFS](bfs/bfs-study-plan.md), Days 4–6 |
| "minimum cost path" (weighted) | Dijkstra / 0-1 BFS | [Dijkstra](dijkstra/dijkstra-study-plan.md) |
| "at most k stops / edges" | Bellman-Ford | [Dijkstra](dijkstra/dijkstra-study-plan.md), Day 6 |
| "all paths", "every combination / permutation / subset" | Backtracking | [Backtracking](backtracking/backtracking-study-plan.md) |
| "number of ways", "minimum/maximum … with choices that overlap" | DP | [DP](dp/dp-study-plan.md) |
| "maximum number of non-overlapping …", "minimum to remove" | Greedy, sort by end | [Intervals](intervals/intervals-study-plan.md), Day 4 |
| "can you reach", "minimum jumps" | Greedy farthest reach | [Greedy](greedy/greedy-study-plan.md), Day 3 |
| "k-th largest / smallest", "top k", "k closest" | Heap capped at k | [Heap](heap/heap-study-plan.md), Days 2–3 |
| "median of a stream" | Two heaps | [Heap](heap/heap-study-plan.md), Day 4 |
| "merge k sorted …" | K-way merge with a heap | [Heap](heap/heap-study-plan.md), Day 5 |
| "minimum capacity / speed / time such that …" | Binary search on the answer | [Binary Search](binary-search/binary-search-study-plan.md), Day 4 |
| "longest / shortest substring or subarray with …" | Sliding window | [Sliding Window](sliding-window/sliding-window-study-plan.md), Days 3–4 |
| "number of subarrays with sum k" | Prefix sum + hash map | [Prefix Sum](prefix-sum/prefix-sum-study-plan.md), Day 5 |
| "maximum in every window of size k" | Monotonic deque | [Sliding Window](sliding-window/sliding-window-study-plan.md), Day 6 |
| "next greater", "days until warmer", "largest rectangle" | Monotonic stack | [Monotonic Stack](monotonic-stack/monotonic-stack-study-plan.md) |
| "connected components", "are these connected", "redundant edge" | Union-Find (or BFS/DFS) | [Union-Find](union-find/union-find-study-plan.md) |
| "minimum cost to connect all" | Kruskal MST | [Union-Find](union-find/union-find-study-plan.md), Day 5 |
| "course order", "build order", "is there a valid ordering" | Topological sort | [Topological Sort](topological-sort/topological-sort-study-plan.md) |
| "starts with", "prefix", "autocomplete", "word search with many words" | Trie | [Trie](trie/trie-study-plan.md) |
| "maximum XOR" | Binary trie | [Trie](trie/trie-study-plan.md), Day 6 |
| "in-place", "O(1) extra space" on an array | Two pointers (fast/slow) | [Two Pointer](two-pointer/two-pointer-study-plan.md), Days 4–6 |
| "cycle" in a linked list | Fast/slow (Floyd) | [Linked List](linked-list/linked-list-study-plan.md), Days 3–4 |
| "cycle" in a directed graph | Three-color DFS / Kahn's | [Topological Sort](topological-sort/topological-sort-study-plan.md), Day 3 |
| "cycle" in an undirected graph | Union-Find | [Union-Find](union-find/union-find-study-plan.md), Day 3 |

---

## Commonly Confused Pairs

| If you're torn between… | Ask yourself | Then use |
|-------------------------|--------------|----------|
| Sliding window vs prefix sum + hash map | Can numbers be negative, or is the condition "== k"? | Yes → prefix sum. No → sliding window |
| Greedy vs DP | Can I prove the local choice is always safe (exchange argument)? Can I find a 3-element counterexample? | Proof → greedy. Counterexample → DP |
| BFS vs Dijkstra | Are all edges the same cost? | Yes → BFS. Only 0/1 → 0-1 BFS. Otherwise → Dijkstra |
| BFS/DFS vs Union-Find | Do I need paths/distances, or only "same group?" while edges keep arriving? | Paths → BFS/DFS. Groups → Union-Find |
| Heap vs sort | Do I need everything in order once, or the min/max repeatedly as data changes? | Once → sort. Repeatedly → heap |
| Backtracking vs DP | Do I need every solution, or just the count/best? | Every solution → backtracking. Count/best → DP (often backtracking + memo) |
| Monotonic stack vs heap | Is the question about the *nearest* greater/smaller, or the overall max/min? | Nearest → stack. Overall → heap |
| Binary search on the array vs on the answer | Is the thing I'm searching for an index in the input? | Yes → classic/boundary. No → search on the answer |

# Algorithms — Interview Study Plans in Go

Eighteen one-week study plans, one per pattern. Each plan teaches a template you can write from memory, explains **why** it works, walks through it on a small input, and gives two practice problems per day. Every template is real Go code with tests.

## Roadmap

| Week | Topic | Study plan | Code | Builds on |
|------|-------|------------|------|-----------|
| 1 | BFS | [bfs-study-plan.md](bfs/bfs-study-plan.md) | [`bfs/`](bfs/) | — |
| 2 | DFS & intro to backtracking | [dfs-study-plan.md](dfs/dfs-study-plan.md) | [`dfs/`](dfs/) | 1 |
| 3 | Two Pointer | [two-pointer-study-plan.md](two-pointer/two-pointer-study-plan.md) | [`two-pointer/`](two-pointer/) | — |
| 4 | Dynamic Programming | [dp-study-plan.md](dp/dp-study-plan.md) | [`dp/`](dp/) | 2 |
| 5 | Heap / Priority Queue | [heap-study-plan.md](heap/heap-study-plan.md) | [`heap/`](heap/) | — |
| 6 | Binary Search | [binary-search-study-plan.md](binary-search/binary-search-study-plan.md) | [`binary-search/`](binary-search/) | — |
| 7 | Linked List | [linked-list-study-plan.md](linked-list/linked-list-study-plan.md) | [`linked-list/`](linked-list/) | 3 |
| 8 | Monotonic Stack | [monotonic-stack-study-plan.md](monotonic-stack/monotonic-stack-study-plan.md) | [`monotonic-stack/`](monotonic-stack/) | — |
| 9 | Sliding Window | [sliding-window-study-plan.md](sliding-window/sliding-window-study-plan.md) | [`sliding-window/`](sliding-window/) | 3, 8 |
| 10 | Prefix Sum & Hashing | [prefix-sum-study-plan.md](prefix-sum/prefix-sum-study-plan.md) | [`prefix-sum/`](prefix-sum/) | 9 |
| 11 | Intervals | [intervals-study-plan.md](intervals/intervals-study-plan.md) | [`intervals/`](intervals/) | 3, 5 |
| 12 | Greedy | [greedy-study-plan.md](greedy/greedy-study-plan.md) | [`greedy/`](greedy/) | 4, 5, 11 |
| 13 | Backtracking | [backtracking-study-plan.md](backtracking/backtracking-study-plan.md) | [`backtracking/`](backtracking/) | 2 |
| 14 | BST & Inorder | [bst-study-plan.md](bst/bst-study-plan.md) | [`bst/`](bst/) | 2 |
| 15 | Trie | [trie-study-plan.md](trie/trie-study-plan.md) | [`trie/`](trie/) | 2, 13 |
| 16 | Union-Find | [union-find-study-plan.md](union-find/union-find-study-plan.md) | [`union-find/`](union-find/) | 1, 2 |
| 17 | Topological Sort | [topological-sort-study-plan.md](topological-sort/topological-sort-study-plan.md) | [`topological-sort/`](topological-sort/) | 1, 2, 4 |
| 18 | Dijkstra & Weighted Shortest Paths | [dijkstra-study-plan.md](dijkstra/dijkstra-study-plan.md) | [`dijkstra/`](dijkstra/) | 1, 5 |

**Not sure which pattern a problem needs?** See the [pattern guide](pattern-guide.md).

Plans refer back to earlier weeks by number ("you solved this with BFS in week 1"), so the numbered order has no forward references. If you're starting over and want the gentlest difficulty curve instead, a common order is: Two Pointer → Sliding Window → Prefix Sum → Binary Search → Linked List → Monotonic Stack → Heap → Intervals → Greedy → BFS → DFS → Backtracking → BST → Trie → Union-Find → Topological Sort → Dijkstra → DP.

## How Each Week Works

| Day | What you do |
|-----|-------------|
| 1 | Learn the template. Write it from memory three times, answer one "why" question, read the dry run. No problems. |
| 2–6 | One variation per day: template adaptation → **why it works** → **dry run** → two problems → complexity and common mistakes. |
| 7 | Review: redo the two hardest days without notes, then check yourself against the self-assessment list. |

Each plan ends with a pattern cheatsheet for quick revision.

Most dry-run tables were produced by actually running the algorithm (with scripts that mirror the Go code), so you can trust the numbers when you trace an example by hand and compare.

## Running the Code

Requires Go 1.22 or newer (the code uses the built-in `min`/`max` and the `slices` package).

```sh
go test ./...                          # every topic
go test ./bfs/ -v                      # one topic, with test names
go test ./dp/ -run TestCanPartition -v # one test
```

Each topic directory has:

- `<topic>-study-plan.md` — the plan
- `<topic>.go` — every template from the plan as a real function (unexported names match the plan where possible; where they differ, the plan's "Run it" line gives the Go name)
- `<topic>_test.go` — table-driven tests, several of which compare against a brute-force version

The code contains the **templates**, not solutions to the practice problems — those are yours to write. When a practice problem *is* the template verbatim, the plan says so.

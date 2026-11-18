# Heap / Priority Queue Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, and DP weeks)  
**Goal:** Implement Go's `container/heap` interface from memory, recognize when a heap beats sorting, adapt min-heap ↔ max-heap ↔ two-heap patterns.

---

## What Is a Heap / Priority Queue?

A **heap** is a tree stored in a slice that always keeps the **smallest** (min-heap) or **largest** (max-heap) element at the top. Adding an element or removing the top costs O(log n); looking at the top costs O(1). The rest of the elements are only partially ordered — which is exactly why it's cheaper than keeping everything sorted.

A **priority queue** is the abstract idea ("always give me the most important item next"); a heap is how it's implemented. In Go you get one by implementing `heap.Interface` and using `container/heap`.

**Typical problems it solves:**

- K-th largest / smallest, top k frequent, k closest points
- Running median of a stream (two heaps)
- Merging k sorted lists or arrays
- Scheduling: meeting rooms, task scheduler, "pick the best available option now"
- Shortest paths with weights (Dijkstra, week 18)

---

## The Core Template

Go has no built-in heap type. You implement the `heap.Interface` (`Len`, `Less`, `Swap`, `Push`, `Pop`) over a slice, then let the `container/heap` package handle sift-up/sift-down.

```go
import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] } // flip to > for a max-heap
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
    *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```

**Max-heap:** copy `MinHeap` and flip `Less`. Nothing else changes.

```go
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // the only difference
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
    *h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```

(Shortcut when you're in a hurry: push `-x` into a `MinHeap` and negate again when you pop.)

**Usage:**

```go
h := &MinHeap{5, 2, 8}
heap.Init(h)             // O(n) — arrange an existing slice into heap order
heap.Push(h, 1)          // O(log n)
top := heap.Pop(h).(int) // O(log n) — always removes h[0], the top
peek := (*h)[0]          // O(1) — look at the top without removing it
```

**Three rules — never break them:**

1. **`Less(i, j)` is the only line that changes between min-heap and max-heap.** `h[i] < h[j]` = min-heap (smallest on top). `h[i] > h[j]` = max-heap (largest on top). Everything else in the interface stays identical.
2. **`Push`/`Pop` operate on the last element of the slice, not the heap top.** The package appends via your `Push`, then sifts it up; it removes the top by swapping it to the end first, then calls your `Pop` to chop it off. You never write sift logic yourself.
3. **Reach for a heap when you only ever need the min/max repeatedly, not a fully sorted order.** Sorting is O(n log n) once; a heap gives O(log n) per insert/extract — better when the data streams in or you only need the top `k`.

**What a heap actually is:** a binary tree stored in a slice, where every parent is ≤ its children (min-heap). For the node at index `i`, its children are at `2i+1` and `2i+2` and its parent is at `(i-1)/2` — no pointers needed. Only the **top** is guaranteed to be the minimum; the rest is only partially ordered, which is exactly why it's cheaper than keeping everything sorted.

```
slice: [1, 2, 8, 5]          tree:      1
index:  0  1  2  3                     / \
                                      2   8
                                     /
                                    5
```

- **Push:** append at the end, then *sift up* — swap with the parent while smaller. At most `log n` swaps (the tree's height).
- **Pop:** swap the top with the last element, shrink the slice, then *sift down* from the top — swap with the smaller child while bigger. At most `log n` swaps.

**Brute force → heap:** "give me the minimum, repeatedly, while elements keep arriving" brute-forces as a linear scan per query (O(n) each) or a re-sort per insert (O(n log n) each). A heap does both operations in O(log n).

**Run it:** every type and function in this file is real code in [`heap.go`](heap.go), tested in [`heap_test.go`](heap_test.go). From the repo root: `go test ./heap/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Core `heap.Interface`  
**Key insight:** `heap.Push`/`heap.Pop` are just `sort.Interface` (`Len`, `Less`, `Swap`) plus `Push`/`Pop` for growing/shrinking the slice — sift-up/down is handled for you.

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each method.
2. Close this file. Write `MinHeap` (all 5 methods) from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does `Pop` remove from the end of the slice instead of the front?*

   > Removing from the front of a slice is O(n) (everything shifts left). The `container/heap` package avoids this by swapping the top element (`h[0]`) with the last element before calling your `Pop`, so your `Pop` only ever does an O(1) slice-shrink at the end. The sift-down that restores heap order after the swap is handled internally by `heap.Pop`, before your `Pop` method is even called.

**Dry run** — pushing `5, 2, 8, 1` into an empty `MinHeap`, then one pop. This mirrors what `container/heap` does internally:

| Operation | Slice after |
|---|---|
| push 5 | `[5]` |
| push 2 | `[2, 5]` — 2 appended at index 1, smaller than parent 5, swapped up |
| push 8 | `[2, 5, 8]` — 8 is bigger than parent 2, stays |
| push 1 | `[1, 2, 8, 5]` — appended at index 3, swapped with parent 5 (index 1), then with parent 2 (index 0) |
| pop → 1 | `[2, 5, 8]` — 1 swapped to the end and cut off by *your* `Pop`; 5 moved to the top and sifted down past 2 |

Notice the slice is never fully sorted (`[1, 2, 8, 5]`) — only `h[0]` is guaranteed. Don't iterate a heap's slice expecting sorted order.

---

## Day 2 — Kth Largest / Smallest

**Pattern:** Bounded heap of size `k`  
**Key insight:** To find the kth largest, keep a **min-heap capped at size k**. Once the heap has k elements, anything smaller than the top can't be in the top-k — pop it immediately. The top of the heap is always the answer.

**Template adaptation:**

```go
func findKthLargest(nums []int, k int) int {
    h := &MinHeap{}
    heap.Init(h)
    for _, n := range nums {
        heap.Push(h, n)
        if h.Len() > k {
            heap.Pop(h) // discard the smallest — it can't be in the top k
        }
    }
    return (*h)[0]
}
```

**Why a min-heap for the *largest*:** the heap is a "club of the top `k` so far". To decide whether a new number gets in, you only need to compare it against the **weakest member** — the smallest one — and that's exactly what a min-heap keeps on top. When the club is over capacity, kick out the weakest. After all numbers, the weakest member of the top-`k` club is, by definition, the `k`th largest.

**Dry run** — `findKthLargest([3, 2, 1, 5, 6, 4], 2)`:

| Push | Over size 2? | Heap after | Top (current answer) |
|---|---|---|---|
| 3 | — | `[3]` | 3 |
| 2 | — | `[2, 3]` | 2 |
| 1 | pop 1 | `[2, 3]` | 2 |
| 5 | pop 2 | `[3, 5]` | 3 |
| 6 | pop 3 | `[5, 6]` | 5 |
| 4 | pop 4 | `[5, 6]` | **5** |

**Run it:** `findKthLargest` in [`heap.go`](heap.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Kth Largest Element in an Array | Medium | https://leetcode.com/problems/kth-largest-element-in-an-array/ |
| 2 | Kth Largest Element in a Stream | Easy | https://leetcode.com/problems/kth-largest-element-in-a-stream/ |

**What to notice — Problem 1:** Counterintuitive at first — finding the *largest* uses a *min*-heap. The heap only ever holds the k largest values seen so far; its smallest member (the top) is exactly the kth largest overall. O(n log k) instead of O(n log n) for a full sort.

**What to notice — Problem 2:** Same heap, now wrapped in a struct with an `Add` method — the heap persists across calls instead of being built once. Initialize by pushing all of `nums`, trimming to size `k` immediately. Each `Add` pushes then trims exactly like Problem 1's loop body.

**Complexity:** O(n log k) time — n pushes/pops on a heap that never exceeds `k + 1` elements. O(k) space. Compare: sorting is O(n log n), and Quickselect is O(n) on average but O(n²) worst case.

**Common mistakes:**

- Using a max-heap and popping `k` times — correct, but O(n + k log n) with O(n) space, and it doesn't work for a stream.
- Calling `h.Push(x)` (your method) instead of `heap.Push(h, x)` (the package function). Your method only appends; it never sifts, so the heap property silently breaks.
- Kth Largest in a Stream: forgetting to trim to `k` in the constructor when the initial `nums` has more than `k` elements.

---

## Day 3 — Top-K with Custom Ordering

**Pattern:** Heap over structs, not primitives  
**Key insight:** `Less` can compare any field — frequency, distance, whatever the problem ranks by. The bounded-heap-of-size-k trick from Day 2 still applies.

**Template adaptation:**

```go
type Pair struct {
    val   int
    count int
}
type MinHeapByCount []Pair

func (h MinHeapByCount) Len() int           { return len(h) }
func (h MinHeapByCount) Less(i, j int) bool { return h[i].count < h[j].count }
func (h MinHeapByCount) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeapByCount) Push(x any)        { *h = append(*h, x.(Pair)) }
func (h *MinHeapByCount) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func topKFrequent(nums []int, k int) []int {
    freq := map[int]int{}
    for _, n := range nums {
        freq[n]++
    }

    h := &MinHeapByCount{}
    heap.Init(h)
    for val, count := range freq {
        heap.Push(h, Pair{val, count})
        if h.Len() > k {
            heap.Pop(h)
        }
    }

    result := make([]int, h.Len())
    for i := range result {
        result[i] = (*h)[i].val
    }
    return result
}
```

**Why it works:** nothing about the heap changes — only what "smaller" means. `container/heap` never looks at your data; it only calls `Less(i, j)`. Rank by count, distance, deadline — anything — and the Day 2 "club of the top `k`" logic carries over unchanged.

**Run it:** `topKFrequent` (with `Pair` and `MinHeapByCount`) in [`heap.go`](heap.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Top K Frequent Elements | Medium | https://leetcode.com/problems/top-k-frequent-elements/ |
| 2 | K Closest Points to Origin | Medium | https://leetcode.com/problems/k-closest-points-to-origin/ |

**What to notice — Problem 1:** Build the frequency map first (a separate O(n) pass), then run the exact bounded-heap loop from Day 2, but ranking by `count` instead of the raw value.

**What to notice — Problem 2:** `Less` compares squared Euclidean distance (`x*x + y*y`) — never take the actual square root, it's slower and unnecessary since you only need relative order. Same bounded-min-heap-of-size-k shape as Problem 1.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Top K Frequent Elements | O(n + u log k), u = distinct values | O(u) for the frequency map |
| K Closest Points | O(n log k) | O(k) |

(Top K Frequent also has an O(n) bucket-sort solution — frequencies are at most `n`, so bucket index = frequency.)

**Common mistakes:**

- Top K Frequent: pushing every element of `nums` instead of each distinct value once.
- K Closest Points: using `math.Sqrt` or floats. Squared integer distances compare the same way and avoid rounding.
- Forgetting that the result order is heap order, not sorted — fine for these problems, but check what's asked.

---

## Day 4 — Two Heaps (Running Median)

**Pattern:** A max-heap for the lower half, a min-heap for the upper half  
**Key insight:** Keep the two halves balanced in size (differing by at most 1). The median is always at the top of one heap, or the average of both tops — no sorting or scanning required per query.

**Template adaptation:**

```go
type MedianFinder struct {
    lo *MaxHeap // lower half, largest on top
    hi *MinHeap // upper half, smallest on top
}

func Constructor() MedianFinder {
    lo, hi := &MaxHeap{}, &MinHeap{}
    heap.Init(lo)
    heap.Init(hi)
    return MedianFinder{lo, hi}
}

func (mf *MedianFinder) AddNum(num int) {
    heap.Push(mf.lo, num)
    heap.Push(mf.hi, heap.Pop(mf.lo)) // always route through lo first, keeps hi's min >= lo's max

    if mf.hi.Len() > mf.lo.Len() {
        heap.Push(mf.lo, heap.Pop(mf.hi))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64((*mf.lo)[0])
    }
    return float64((*mf.lo)[0]+(*mf.hi)[0]) / 2.0
}
```

**Why it works:** two invariants, restored after every insert:

1. **Order:** every number in `lo` ≤ every number in `hi`. Pushing into `lo` then moving `lo`'s max to `hi` guarantees the new number ends up on the correct side.
2. **Size:** `len(lo)` is `len(hi)` or `len(hi) + 1`. The final `if` moves one back if `hi` got bigger.

With both true, the median is the middle of the combined sorted order: `lo`'s top (odd count) or the average of both tops (even count).

**Dry run** — adding `5, 15, 1, 3, 8` (heaps shown sorted for readability):

| Add | Moved back to `lo`? | `lo` (max-heap) | `hi` (min-heap) | Median |
|---|---|---|---|---|
| 5 | yes | `[5]` | `[]` | 5 |
| 15 | no | `[5]` | `[15]` | 10 |
| 1 | yes | `[1, 5]` | `[15]` | 5 |
| 3 | no | `[1, 3]` | `[5, 15]` | 4 |
| 8 | yes | `[1, 3, 5]` | `[8, 15]` | 5 |

Look at the row for `3`: it went into `lo`, but `lo`'s max `5` was bounced to `hi` — that's invariant 1 keeping `5` on the correct side.

**Run it:** `MedianFinder` in [`heap.go`](heap.go) (the constructor is `newMedianFinder` there, since `Constructor` is only a LeetCode naming convention).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Find Median from Data Stream | Hard | https://leetcode.com/problems/find-median-from-data-stream/ |
| 2 | IPO | Hard | https://leetcode.com/problems/ipo/ |

**What to notice — Problem 1:** `lo` (max-heap) never has more than one extra element over `hi` (min-heap). Routing every insert through `lo` first, then rebalancing into `hi`, guarantees `lo`'s max is always ≤ `hi`'s min — that invariant is what makes the median readable in O(1).

**What to notice — Problem 2:** Two heaps, different split: a min-heap of projects by required capital (to find what's affordable), and a max-heap of profits (to greedily pick the best affordable project each round). At each of `k` rounds, pop every project from the capital-heap that's now affordable into the profit-heap, then take the profit-heap's max.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Find Median from Data Stream | O(log n) per `AddNum`, O(1) per `FindMedian` | O(n) |
| IPO | O(n log n) — sort/heapify by capital, then each project enters and leaves the profit heap once | O(n) |

**Common mistakes:**

- Median: pushing into whichever heap "looks right" by comparing to one top. It works but has more edge cases than always routing through `lo`.
- Median: integer division — `(a + b) / 2` with ints truncates. Convert to `float64` first.
- IPO: re-scanning all projects every round (O(n·k)) instead of using a capital-sorted min-heap so each project is moved once.

---

## Day 5 — K-Way Merge

**Pattern:** Heap holding one "current" element per list/source  
**Key insight:** Never merge lists pairwise. Put the head of every list in a heap keyed by value; repeatedly pop the smallest, then push that list's next element. This merges k sources in O(n log k) instead of O(nk).

**Template adaptation:**

```go
type ListItem struct {
    val      int
    listIdx  int
    nodeIdx  int
}
type MinHeapByVal []ListItem

func (h MinHeapByVal) Len() int           { return len(h) }
func (h MinHeapByVal) Less(i, j int) bool { return h[i].val < h[j].val }
func (h MinHeapByVal) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeapByVal) Push(x any)        { *h = append(*h, x.(ListItem)) }
func (h *MinHeapByVal) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func mergeKSorted(lists [][]int) []int {
    h := &MinHeapByVal{}
    heap.Init(h)
    for i, list := range lists {
        if len(list) > 0 {
            heap.Push(h, ListItem{list[0], i, 0})
        }
    }

    var result []int
    for h.Len() > 0 {
        top := heap.Pop(h).(ListItem)
        result = append(result, top.val)

        if top.nodeIdx+1 < len(lists[top.listIdx]) {
            next := lists[top.listIdx][top.nodeIdx+1]
            heap.Push(h, ListItem{next, top.listIdx, top.nodeIdx + 1})
        }
    }
    return result
}
```

**Why it works:** the smallest remaining element overall must be the head of *some* list — every list is sorted, so nothing behind a head can be smaller than it. The heap holds exactly the k heads, so its top is the global minimum. After popping it, only that one list's head changed, so pushing its successor restores "heap = all current heads".

**Dry run** — merging `[1, 4, 5]`, `[1, 3, 4]`, `[2, 6]`:

| Pop | Push (same list's next) | Heap after | Output |
|---|---|---|---|
| 1 (list 0) | 4 (list 0) | `[1, 2, 4]` | `[1]` |
| 1 (list 1) | 3 (list 1) | `[2, 3, 4]` | `[1, 1]` |
| 2 (list 2) | 6 (list 2) | `[3, 4, 6]` | `[1, 1, 2]` |
| 3 (list 1) | 4 (list 1) | `[4, 4, 6]` | `[1, 1, 2, 3]` |
| 4 (list 0) | 5 (list 0) | `[4, 5, 6]` | `[1, 1, 2, 3, 4]` |
| 4 (list 1) | — | `[5, 6]` | `[1, 1, 2, 3, 4, 4]` |
| 5 (list 0) | — | `[6]` | `[1, 1, 2, 3, 4, 4, 5]` |
| 6 (list 2) | — | `[]` | `[1, 1, 2, 3, 4, 4, 5, 6]` |

The heap never holds more than 3 items (k = 3), however long the lists are.

**Run it:** `mergeKSorted` in [`heap.go`](heap.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Merge k Sorted Lists | Hard | https://leetcode.com/problems/merge-k-sorted-lists/ |
| 2 | Kth Smallest Element in a Sorted Matrix | Medium | https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/ |

**What to notice — Problem 1:** Same shape as the template but over `*ListNode` instead of `[]int` — push each list's head, pop the min, push its `.Next` if non-nil. The heap never holds more than k nodes at once.

**What to notice — Problem 2:** Each matrix row is sorted, so treat every row like a "list" in k-way merge. Push `(matrix[i][0], row=i, col=0)` for every row, pop k times; each pop advances that row's column pointer. The kth pop is the answer — no need to drain the whole heap.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Merge k Sorted Lists | O(N log k), N = total nodes | O(k) |
| Kth Smallest in Sorted Matrix | O(k log n) for an n×n matrix | O(n) |

Merging pairwise one list at a time is O(N·k); merging in pairs like merge sort is also O(N log k) — a good alternative if you're asked to avoid a heap.

**Common mistakes:**

- Pushing `nil` list heads (or empty slices) into the heap.
- Storing only the value in the heap — then you don't know which list to advance. Store `(value, list index, position)` or the node pointer itself.

---

## Day 6 — Greedy Scheduling with a Heap

**Pattern:** Heap tracks "what's currently in progress" or "what's available now"  
**Key insight:** Greedy interval/scheduling problems often need to know the min (soonest-ending) or max (highest-priority) of an active set at every step — that's exactly what a heap tracks in O(log n) per update.

**Template adaptation:**

```go
func minMeetingRooms(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

    endTimes := &MinHeap{} // end times of meetings currently in progress
    heap.Init(endTimes)

    for _, iv := range intervals {
        if endTimes.Len() > 0 && (*endTimes)[0] <= iv[0] {
            heap.Pop(endTimes) // earliest-ending meeting is done — reuse its room
        }
        heap.Push(endTimes, iv[1])
    }
    return endTimes.Len() // one room per meeting still in the heap
}
```

**Why it works (Meeting Rooms II):** process meetings in start order. The heap holds the end time of every room currently in use, so its top is the room that frees up **soonest**. If even that room isn't free when the next meeting starts, no room is, so you need a new one. Rooms are only ever reused, never released, so the heap's final size is the peak number of simultaneous meetings.

**Dry run** — `[[0, 30], [5, 10], [15, 20], [25, 35]]`:

| Meeting | Action | Heap (end times) | Rooms |
|---|---|---|---|
| [0, 30] | new room | `[30]` | 1 |
| [5, 10] | new room (30 > 5) | `[10, 30]` | 2 |
| [15, 20] | reuse (room free at 10) | `[20, 30]` | 2 |
| [25, 35] | reuse (room free at 20) | `[30, 35]` | 2 |

**Run it:** `minMeetingRooms` in [`heap.go`](heap.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Meeting Rooms II | Medium | https://leetcode.com/problems/meeting-rooms-ii/ (LeetCode Premium — use lintcode.com/problem/919 if needed) |
| 2 | Task Scheduler | Medium | https://leetcode.com/problems/task-scheduler/ |

**What to notice — Problem 1:** Sort by start time first — greedy correctness depends on processing meetings in the order they begin. The heap holds end times of rooms currently occupied; its size at the end is the max rooms ever needed simultaneously.

**What to notice — Problem 2:** Max-heap of task frequencies. Each cycle, pop up to `n+1` of the most frequent tasks, decrement their counts, and push back any with count remaining. The heap always surfaces "which task is most urgent to schedule next" — count idle slots when the heap empties before a full `n+1` cycle.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Meeting Rooms II | O(n log n) — the sort, plus one push/pop per meeting | O(n) |
| Task Scheduler | O(total tasks × log 26) = O(n) — at most 26 distinct tasks | O(1) |

**Common mistakes:**

- Meeting Rooms: forgetting to sort by start time first.
- Meeting Rooms: `<` instead of `<=` when checking whether a room is free — a meeting ending at 10 frees the room for one starting at 10.
- Task Scheduler: pushing a task back *during* the same cycle it was popped in, letting it run twice within `n` slots. Collect the cycle's popped tasks and push them back after the cycle.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** The heap itself is always the same 5 methods. The skill is recognizing which of the five patterns (bounded-k, custom ordering, two-heap, k-way merge, greedy-active-set) a new problem needs, and what `Less` should compare.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Answer out loud: *For a brand-new problem, how do you decide whether you need a heap at all, and if so, min or max?*

   > Reach for a heap when you repeatedly need the current min/max of a changing set, or only the top `k` of a stream — not when you need the data fully sorted (then just sort). For direction: if you're discarding everything except the largest values seen so far, you paradoxically want a **min**-heap capped at size k (Day 2); if you want direct access to the largest/most-urgent item right now, you want a **max**-heap (Day 6).

**Self-assessment — can you answer these without hesitation?**

- [ ] Write all 5 methods of `heap.Interface` from memory
- [ ] Why does finding the kth *largest* element use a *min*-heap?
- [ ] In the two-heap median pattern, why route every insert through `lo` before rebalancing?
- [ ] Why is k-way merge O(n log k) instead of O(nk)?
- [ ] What invariant does `minMeetingRooms` rely on for its heap to give the right answer?

---

## Pattern Cheatsheet

| Pattern | Heap type | What it tracks | Key move | Time |
|---------|-----------|-----------------|----------|------|
| Kth largest/smallest | min-heap, capped at k | top-k values seen so far | pop when size exceeds k | O(n log k) |
| Top-K custom ordering | min-heap over structs, capped at k | top-k by any field | same as above, `Less` compares the ranking field | O(n log k) |
| Two heaps (median) | max-heap (lower) + min-heap (upper) | balanced split at the median | route inserts through one heap, rebalance | O(log n) per insert |
| K-way merge | min-heap over `(val, source, idx)` | current head of each source | pop min, push that source's next | O(N log k) |
| Greedy active-set | min or max-heap | in-progress/available items | pop when an item's constraint is resolved | O(n log n) |

---

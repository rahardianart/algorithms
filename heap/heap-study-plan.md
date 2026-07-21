# Heap / Priority Queue Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, and DP weeks)  
**Goal:** Implement Go's `container/heap` interface from memory, recognize when a heap beats sorting, adapt min-heap ↔ max-heap ↔ two-heap patterns.

---

## The Core Template

Go has no built-in heap type. You implement the `heap.Interface` (`Len`, `Less`, `Swap`, `Push`, `Pop`) over a slice, then let the `container/heap` package handle sift-up/sift-down.

```go
import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] } // flip to > for a max-heap
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }

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

**Usage:**

```go
h := &MinHeap{5, 2, 8}
heap.Init(h)          // O(n) — arrange an existing slice into heap order
heap.Push(h, 1)        // O(log n)
min := heap.Pop(h).(int) // O(log n) — always removes h[0], the top
```

**Three rules — never break them:**

1. **`Less(i, j)` is the only line that changes between min-heap and max-heap.** `h[i] < h[j]` = min-heap (smallest on top). `h[i] > h[j]` = max-heap (largest on top). Everything else in the interface stays identical.
2. **`Push`/`Pop` operate on the last element of the slice, not the heap top.** The package appends via your `Push`, then sifts it up; it removes the top by swapping it to the end first, then calls your `Pop` to chop it off. You never write sift logic yourself.
3. **Reach for a heap when you only ever need the min/max repeatedly, not a fully sorted order.** Sorting is O(n log n) once; a heap gives O(log n) per insert/extract — better when the data streams in or you only need the top `k`.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Kth Largest Element in an Array | Medium | https://leetcode.com/problems/kth-largest-element-in-an-array/ |
| 2 | Kth Largest Element in a Stream | Easy | https://leetcode.com/problems/kth-largest-element-in-a-stream/ |

**What to notice — Problem 1:** Counterintuitive at first — finding the *largest* uses a *min*-heap. The heap only ever holds the k largest values seen so far; its smallest member (the top) is exactly the kth largest overall. O(n log k) instead of O(n log n) for a full sort.

**What to notice — Problem 2:** Same heap, now wrapped in a struct with an `Add` method — the heap persists across calls instead of being built once. Initialize by pushing all of `nums`, trimming to size `k` immediately. Each `Add` pushes then trims exactly like Problem 1's loop body.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Top K Frequent Elements | Medium | https://leetcode.com/problems/top-k-frequent-elements/ |
| 2 | K Closest Points to Origin | Medium | https://leetcode.com/problems/k-closest-points-to-origin/ |

**What to notice — Problem 1:** Build the frequency map first (a separate O(n) pass), then run the exact bounded-heap loop from Day 2, but ranking by `count` instead of the raw value.

**What to notice — Problem 2:** `Less` compares squared Euclidean distance (`x*x + y*y`) — never take the actual square root, it's slower and unnecessary since you only need relative order. Same bounded-min-heap-of-size-k shape as Problem 1.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Find Median from Data Stream | Hard | https://leetcode.com/problems/find-median-from-data-stream/ |
| 2 | IPO | Hard | https://leetcode.com/problems/ipo/ |

**What to notice — Problem 1:** `lo` (max-heap) never has more than one extra element over `hi` (min-heap). Routing every insert through `lo` first, then rebalancing into `hi`, guarantees `lo`'s max is always ≤ `hi`'s min — that invariant is what makes the median readable in O(1).

**What to notice — Problem 2:** Two heaps, different split: a min-heap of projects by required capital (to find what's affordable), and a max-heap of profits (to greedily pick the best affordable project each round). At each of `k` rounds, pop every project from the capital-heap that's now affordable into the profit-heap, then take the profit-heap's max.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Merge k Sorted Lists | Hard | https://leetcode.com/problems/merge-k-sorted-lists/ |
| 2 | Kth Smallest Element in a Sorted Matrix | Medium | https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/ |

**What to notice — Problem 1:** Same shape as the template but over `*ListNode` instead of `[]int` — push each list's head, pop the min, push its `.Next` if non-nil. The heap never holds more than k nodes at once.

**What to notice — Problem 2:** Each matrix row is sorted, so treat every row like a "list" in k-way merge. Push `(matrix[i][0], row=i, col=0)` for every row, pop k times; each pop advances that row's column pointer. The kth pop is the answer — no need to drain the whole heap.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Meeting Rooms II | Medium | https://leetcode.com/problems/meeting-rooms-ii/ (LeetCode Premium — use lintcode.com/problem/919 if needed) |
| 2 | Task Scheduler | Medium | https://leetcode.com/problems/task-scheduler/ |

**What to notice — Problem 1:** Sort by start time first — greedy correctness depends on processing meetings in the order they begin. The heap holds end times of rooms currently occupied; its size at the end is the max rooms ever needed simultaneously.

**What to notice — Problem 2:** Max-heap of task frequencies. Each cycle, pop up to `n+1` of the most frequent tasks, decrement their counts, and push back any with count remaining. The heap always surfaces "which task is most urgent to schedule next" — count idle slots when the heap empties before a full `n+1` cycle.

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

| Pattern | Heap type | What it tracks | Key move |
|---------|-----------|-----------------|----------|
| Kth largest/smallest | min-heap, capped at k | top-k values seen so far | pop when size exceeds k |
| Top-K custom ordering | min-heap over structs, capped at k | top-k by any field | same as above, `Less` compares the ranking field |
| Two heaps (median) | max-heap (lower) + min-heap (upper) | balanced split at the median | route inserts through one heap, rebalance |
| K-way merge | min-heap over `(val, source, idx)` | current head of each source | pop min, push that source's next |
| Greedy active-set | min or max-heap | in-progress/available items | pop when an item's constraint is resolved |

---

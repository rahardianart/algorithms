# Intervals Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 11 — uses sorting, two pointers, and the Heap week's Meeting Rooms II)  
**Goal:** Sort by the right endpoint for the job, write the merge and overlap checks without off-by-ones, and recognize merge / schedule / cover / sweep problems on sight.

---

## The Core Templates

### Overlap Check (closed intervals `[start, end]`)

```go
// a and b overlap when each one starts before the other ends
overlaps := a[0] <= b[1] && b[0] <= a[1]

// and the overlapping part (if any) is
lo, hi := max(a[0], b[0]), min(a[1], b[1]) // overlap exists iff lo <= hi
```

### Sort + Single Pass

```go
sort.Slice(intervals, func(i, j int) bool {
    return intervals[i][0] < intervals[j][0] // by start (merge, cover) — or [1] by end (scheduling)
})

for _, iv := range intervals {
    // compare iv with the last thing you kept / the running boundary
}
```

**Three rules — never break them:**

1. **Sort first, and know which endpoint you sort by.** Sort by **start** when you're combining or covering (merge, insert, cover a range). Sort by **end** when you're choosing which intervals to keep (max non-overlapping, min arrows). Picking the wrong one gives plausible but wrong answers.
2. **Decide whether touching endpoints overlap, then be consistent.** `[1, 3]` and `[3, 5]`: overlap for merging closed ranges (`<=`), don't conflict for back-to-back meetings (`<`). Every off-by-one in this topic is this decision made inconsistently.
3. **After sorting, you only ever compare against one thing** — the last merged interval, the last kept end, the current coverage boundary. If you're comparing against all previous intervals, you've lost the point of sorting.

**Brute force → sort + one pass:** checking every pair of intervals for overlap is O(n²). After sorting by start, an interval can only overlap the ones *just before it*, so a single pass that tracks "the last thing I kept" is enough: O(n log n) for the sort plus O(n) for the pass.

**Complexity:** O(n log n) time (the sort dominates), O(n) space for the output.

**Run it:** every template in this file is a real function in [`intervals.go`](intervals.go), tested in [`intervals_test.go`](intervals_test.go). From the repo root: `go test ./intervals/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Overlap check + sort/scan  
**Key insight:** Two intervals *don't* overlap only when one ends before the other starts. Negate that and you get the overlap check — no case analysis needed.

**No problems today.** Do this instead:

1. Read the templates above slowly. Draw the six ways two intervals can sit relative to each other (before, after, overlapping left, overlapping right, one inside the other, identical) and check the overlap formula on each.
2. Close this file. Write the overlap check and the "overlap part" formula from memory.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why is `a[0] <= b[1] && b[0] <= a[1]` enough, with no special case for one interval containing the other?*

   > Because it's the negation of the only two ways *not* to overlap: `a` entirely before `b` (`a[1] < b[0]`) or `b` entirely before `a` (`b[1] < a[0]`). Containment, partial overlap, and equality all fall out automatically — they're just "neither of those two".

---

## Day 2 — Merge Overlapping Intervals

**Pattern:** Sort by start; extend the last merged interval or start a new one  
**Key insight:** After sorting by start, the next interval can only overlap the **last** merged interval — everything earlier ends before that one does, or it would have been merged into it.

**Template adaptation:**

```go
func mergeIntervals(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

    var merged [][]int
    for _, iv := range intervals {
        if len(merged) > 0 && iv[0] <= merged[len(merged)-1][1] {
            last := merged[len(merged)-1]
            last[1] = max(last[1], iv[1]) // extend — but never shrink
        } else {
            merged = append(merged, []int{iv[0], iv[1]})
        }
    }
    return merged
}
```

**Why it works:** invariant — `merged` is always the merge of every interval processed so far, sorted and disjoint. A new interval starts at or after every previous start, so it either reaches back into the last merged block (`iv[0] <= last end`) or starts after it and therefore after everything.

**Dry run** — `[[1,3], [8,10], [2,6], [15,18], [17,20]]`, sorted to `[[1,3], [2,6], [8,10], [15,18], [17,20]]`:

| Interval | Overlaps last merged? | `merged` after |
|---|---|---|
| [1, 3] | first interval | `[[1, 3]]` |
| [2, 6] | yes (2 ≤ 3) → extend end to 6 | `[[1, 6]]` |
| [8, 10] | no → start a new one | `[[1, 6], [8, 10]]` |
| [15, 18] | no → start a new one | `[[1, 6], [8, 10], [15, 18]]` |
| [17, 20] | yes (17 ≤ 18) → extend end to 20 | `[[1, 6], [8, 10], [15, 20]]` |

**Run it:** `mergeIntervals` in [`intervals.go`](intervals.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Merge Intervals | Medium | https://leetcode.com/problems/merge-intervals/ |
| 2 | Insert Interval | Medium | https://leetcode.com/problems/insert-interval/ |

**What to notice — Problem 1:** This IS the template. Get it automatic.

**What to notice — Problem 2:** The input is already sorted, so you can do it in O(n) without re-sorting: copy every interval that ends before `newInterval` starts, then merge everything that overlaps into `newInterval` (min of starts, max of ends), then copy the rest.

**Complexity:** Problem 1 O(n log n). Problem 2 O(n) — the input is already sorted.

**Common mistakes:**

- `last[1] = iv[1]` instead of `max(last[1], iv[1])` — breaks when `iv` is entirely inside `last` (`[1, 10]` then `[2, 3]` would shrink to `[1, 3]`).
- Appending `iv` itself to `merged` and then modifying it — that mutates the caller's input. Append a copy.
- Insert Interval: forgetting to add `newInterval` when it goes after every existing interval.

---

## Day 3 — Intersections and Conflict Checks

**Pattern:** Walk two sorted lists with two pointers, or scan one sorted list for any conflict  
**Key insight:** For two lists, always advance the pointer whose interval **ends first** — it can't intersect anything further along the other list.

**Template adaptation:**

```go
// intersect two lists that are each sorted and internally disjoint
func intervalIntersection(a, b [][]int) [][]int {
    var result [][]int
    i, j := 0, 0
    for i < len(a) && j < len(b) {
        lo := max(a[i][0], b[j][0])
        hi := min(a[i][1], b[j][1])
        if lo <= hi {
            result = append(result, []int{lo, hi})
        }
        if a[i][1] < b[j][1] { // a[i] ends first — it's done
            i++
        } else {
            j++
        }
    }
    return result
}
```

**Why advancing the earlier end is safe:** say `a[i]` ends before `b[j]`. Every later interval in `b` starts after `b[j]` ends, which is after `a[i]` ends — so `a[i]` can't touch any of them. Discard it. This is the Two Pointer week's "each step rules out one element for good", applied to intervals.

**Dry run** — `a = [[0,2], [5,10], [13,23], [24,25]]`, `b = [[1,5], [8,12], [15,24], [25,26]]`:

| `a[i]` | `b[j]` | `[max start, min end]` | Overlap? | Advance |
|---|---|---|---|---|
| [0, 2] | [1, 5] | [1, 2] | yes → [1, 2] | `i` (a ends first) |
| [5, 10] | [1, 5] | [5, 5] | yes → [5, 5] | `j` (b ends first or tie) |
| [5, 10] | [8, 12] | [8, 10] | yes → [8, 10] | `i` (a ends first) |
| [13, 23] | [8, 12] | [13, 12] | no | `j` (b ends first) |
| [13, 23] | [15, 24] | [15, 23] | yes → [15, 23] | `i` (a ends first) |
| [24, 25] | [15, 24] | [24, 24] | yes → [24, 24] | `j` (b ends first or tie) |
| [24, 25] | [25, 26] | [25, 25] | yes → [25, 25] | `i` (a ends first) |

**Run it:** `intervalIntersection` in [`intervals.go`](intervals.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Interval List Intersections | Medium | https://leetcode.com/problems/interval-list-intersections/ |
| 2 | Meeting Rooms | Easy | https://leetcode.com/problems/meeting-rooms/ (LeetCode Premium — use lintcode.com/problem/920 if needed) |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Sort by start, then a person can attend everything iff no meeting starts before the previous one ends: `intervals[i][0] < intervals[i-1][1]` means conflict. Meetings are back-to-back-friendly, so `[1,5]` then `[5,8]` is fine — strict `<`.

**Complexity:** Problem 1 O(m + n). Problem 2 O(n log n).

**Common mistakes:**

- Intersections: advancing both pointers after a match — you skip intersections when one interval overlaps several on the other side.
- Meeting Rooms: checking only adjacent pairs *without sorting first*.

---

## Day 4 — Scheduling: Keep the Most Intervals

**Pattern:** Sort by **end**; greedily keep every interval that starts after the last kept one ends  
**Key insight:** The interval that ends first leaves the most room for everything else. Always keep it.

**Template adaptation:**

```go
// max number of intervals with no two overlapping (touching is fine)
func maxNonOverlapping(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })

    kept := 0
    lastEnd := 0
    for i, iv := range intervals {
        if i == 0 || iv[0] >= lastEnd {
            kept++
            lastEnd = iv[1]
        }
    }
    return kept
}
```

**Why sorting by end is correct (exchange argument):** take any optimal selection. Its first interval ends no earlier than the interval with the earliest end, `E`. Swap that first interval for `E`: `E` ends even sooner, so it can't conflict with the rest of the selection. Same count, now starting with `E`. Repeat the argument on what's left — greedy matches optimal at every step.

**Why sorting by start is wrong — a real run:** `[[1,10], [2,3], [4,5], [6,7]]`.

| Sort by | Order | Kept |
|---|---|---|
| start | `[1,10]`, `[2,3]`, `[4,5]`, `[6,7]` | `[1,10]` only → **1** (it blocks everything) |
| end | `[2,3]`, `[4,5]`, `[6,7]`, `[1,10]` | `[2,3]`, `[4,5]`, `[6,7]` → **3** |

**Run it:** `maxNonOverlapping` in [`intervals.go`](intervals.go) — the test includes this exact case.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Non-overlapping Intervals | Medium | https://leetcode.com/problems/non-overlapping-intervals/ |
| 2 | Minimum Number of Arrows to Burst Balloons | Medium | https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/ |

**What to notice — Problem 1:** Asks for the minimum number to **remove** = `n − maxNonOverlapping`. Same template.

**What to notice — Problem 2:** Same greedy, different question: shoot an arrow at the end of the earliest-ending balloon; it pops every balloon that starts at or before that point. Balloons are closed intervals, so touching **does** count as overlapping here (`iv[0] > lastEnd` starts a new arrow). Values go down to −2³¹, so don't initialize `lastEnd` to 0.

**Complexity:** both O(n log n).

**Common mistakes:**

- Sorting by start (see the run above).
- Arrows: using `>=` (meeting semantics) instead of `>` (closed-interval semantics).
- Initializing `lastEnd := 0` when coordinates can be negative — use the first interval, as the template does with `i == 0`.

---

## Day 5 — Covering a Range

**Pattern:** Sort by start; from the current coverage boundary, jump to the farthest reachable end  
**Key insight:** Among all intervals that start inside what you've already covered, only the one reaching **farthest** matters. It's the BFS-level idea from Jump Game II, on intervals.

**Template adaptation:**

```go
// fewest intervals whose union covers [start, end], or -1
func minIntervalsToCover(intervals [][]int, start, end int) int {
    sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

    count, covered, i := 0, start, 0
    for covered < end {
        farthest := covered
        for i < len(intervals) && intervals[i][0] <= covered {
            farthest = max(farthest, intervals[i][1]) // best reach from here
            i++
        }
        if farthest == covered {
            return -1 // nothing extends the coverage: a gap
        }
        count++
        covered = farthest
    }
    return count
}
```

**Why it works:** each round answers "using one more interval, how far can coverage get?" Any interval that starts within `[start, covered]` is a legal next pick, and picking the one that reaches farthest can only leave you better placed for every later round. Each interval is looked at once (the `i` pointer never resets), so the loop is O(n) after the sort.

**Dry run** — cover `[0, 10]` with `[[0,2], [4,6], [8,10], [1,9], [1,5], [5,9]]` (sorted: `[0,2], [1,5], [1,9], [4,6], [5,9], [8,10]`):

| Round | Covered up to | Candidates starting ≤ covered | Farthest reach | Count |
|---|---|---|---|---|
| 1 | 0 | [0, 2] | 2 | 1 |
| 2 | 2 | [1, 5], [1, 9] | 9 | 2 |
| 3 | 9 | [4, 6], [5, 9], [8, 10] | 10 | 3 |

Answer **3**: `[0,2]`, `[1,9]`, `[8,10]`.

**Run it:** `minIntervalsToCover` in [`intervals.go`](intervals.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Video Stitching | Medium | https://leetcode.com/problems/video-stitching/ |
| 2 | Minimum Number of Taps to Open to Water a Garden | Hard | https://leetcode.com/problems/minimum-number-of-taps-to-open-to-water-a-garden/ |

**What to notice — Problem 1:** The template with `start = 0`, `end = time`.

**What to notice — Problem 2:** Convert each tap into an interval first: tap `i` with range `r` covers `[i − r, i + r]`, clipped to `[0, n]`. After that it's the template with `start = 0`, `end = n`.

**Complexity:** both O(n log n) for the sort plus O(n) for the scan. (Both problems also have an O(n) version that skips the sort by recording, for each start position, the farthest end — the Jump Game II trick.)

**Common mistakes:**

- Advancing `covered` to the end of the *first* candidate instead of the *farthest* one.
- Not detecting the gap (`farthest == covered`) — infinite loop.
- Taps: forgetting to clip `i − r` to 0.

---

## Day 6 — Sweep Line

**Pattern:** Turn each interval into two events (+1 at start, −1 at end), sort the events, and keep a running count  
**Key insight:** You don't care which intervals are active — only how many. Events let you track that count by touching each endpoint once.

**Template adaptation:**

```go
// most half-open [start, end) intervals active at the same time
func maxConcurrent(intervals [][]int) int {
    type event struct{ time, delta int }
    events := make([]event, 0, 2*len(intervals))
    for _, iv := range intervals {
        events = append(events, event{iv[0], +1}, event{iv[1], -1})
    }
    sort.Slice(events, func(i, j int) bool {
        if events[i].time != events[j].time {
            return events[i].time < events[j].time
        }
        return events[i].delta < events[j].delta // ends before starts at the same time
    })

    active, best := 0, 0
    for _, e := range events {
        active += e.delta
        best = max(best, active)
    }
    return best
}
```

**Why the tie-break matters:** with half-open intervals, a meeting ending at 4 and one starting at 4 don't overlap. Processing the −1 before the +1 at the same time makes `active` dip before it rises, so they're never counted together. Flip the tie-break and you'd count them as concurrent — that's rule 2 again.

**Dry run** — `[[1,4], [2,6], [4,7], [5,8]]`:

| Time | Event | Active | Best |
|---|---|---|---|
| 1 | start (+1) | 1 | 1 |
| 2 | start (+1) | 2 | 2 |
| 4 | end (−1) | 1 | 2 |
| 4 | start (+1) | 2 | 2 |
| 5 | start (+1) | 3 | 3 |
| 6 | end (−1) | 2 | 3 |
| 7 | end (−1) | 1 | 3 |
| 8 | end (−1) | 0 | 3 |

At time 4, the end of `[1,4)` is processed first, so the count never reaches 3 there. It reaches 3 at time 5 (`[2,6)`, `[4,7)`, `[5,8)`). Compare with Meeting Rooms II from the Heap week: same answer, no heap.

**Run it:** `maxConcurrent` in [`intervals.go`](intervals.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | My Calendar II | Medium | https://leetcode.com/problems/my-calendar-ii/ |
| 2 | Number of Flowers in Full Bloom | Hard | https://leetcode.com/problems/number-of-flowers-in-full-bloom/ |

**What to notice — Problem 1:** On each `book`, tentatively add the two events, sweep, and if `active` ever hits 3, remove them and return false. O(n) per booking with a sorted map; O(n log n) if you re-sort a slice — both pass.

**What to notice — Problem 2:** Flowers in bloom at time `t` = (number of starts ≤ t) − (number of ends < t). Sort the starts and ends separately and answer each person with two binary searches (Binary Search week's `leftmostTrue`). Note: bloom intervals are **closed**, so an end at `t` still counts as blooming.

**Complexity:** Problem 1 O(n²) total over n bookings with the simple sweep. Problem 2 O((n + p) log n) for n flowers and p people.

**Common mistakes:**

- Wrong tie-break order for the interval semantics you need.
- Problem 2: treating bloom intervals as half-open.
- Using a `map[int]int` of deltas and iterating it directly — Go map iteration order is random. Sort the keys first.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Classify first: merging/combining (sort by start), keeping/selecting (sort by end), covering (sort by start, jump farthest), or counting overlap (sweep events). Then decide whether touching endpoints count.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing code, say out loud: *What do I sort by, and do touching endpoints overlap here?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the overlap check and the "overlap part" formula from memory
- [ ] Why does merging only ever compare against the last merged interval?
- [ ] Why does max-non-overlapping sort by end, and what goes wrong sorting by start?
- [ ] In interval intersection, which pointer advances and why is that safe?
- [ ] In a sweep line, why does the tie-break between +1 and −1 at the same time matter?

---

## Pattern Cheatsheet

| Pattern | Sort by | Compare against | Key move | Time |
|---------|---------|-----------------|----------|------|
| Merge | start | last merged interval | extend end with `max`, or start new | O(n log n) |
| Intersection (two lists) | already sorted | the other list's current interval | `[max starts, min ends]`; advance earlier end | O(m + n) |
| Max non-overlapping | **end** | last kept end | keep if `start >= lastEnd` | O(n log n) |
| Cover a range | start | current coverage boundary | jump to farthest reachable end | O(n log n) |
| Sweep line | event time (+ tie-break) | running active count | +1 at start, −1 at end | O(n log n) |

---

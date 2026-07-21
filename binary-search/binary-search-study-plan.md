# Binary Search Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, DP, and Heap weeks)  
**Goal:** Write both binary search templates from memory, recognize monotonic-predicate problems even when there's no sorted array in sight.

---

## The Two Templates

### Classic (find an exact value)

```go
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2 // avoids overflow vs (lo+hi)/2
        switch {
        case nums[mid] == target:
            return mid
        case nums[mid] < target:
            lo = mid + 1
        default:
            hi = mid - 1
        }
    }
    return -1
}
```

### Boundary (find the first index where a condition flips false→true)

```go
func leftmostTrue(lo, hi int, condition func(int) bool) int {
    for lo < hi { // half-open: [lo, hi)
        mid := lo + (hi-lo)/2
        if condition(mid) {
            hi = mid // mid could be the answer — keep it in range
        } else {
            lo = mid + 1 // mid is definitely not the answer — exclude it
        }
    }
    return lo // lo == hi: the boundary
}
```

**Three rules — never break them:**

1. **`lo + (hi-lo)/2`, not `(lo+hi)/2`.** Habit worth keeping even where Go's int range makes overflow unlikely — it's the version that ports cleanly to any language.
2. **Pick one loop invariant and don't mix templates.** Classic search uses inclusive bounds (`lo <= hi`, both `mid+1`/`mid-1`). Boundary search uses half-open bounds (`lo < hi`, `hi = mid` because `mid` might be the answer). Mixing inclusive and half-open logic is the #1 source of infinite loops and off-by-ones.
3. **Binary search needs a monotonic predicate, not a sorted array.** If you can define a `condition(x)` that's `false` for a while then `true` forever (or vice versa) over some ordered search space, you can binary search that space — the space doesn't have to be the input array itself. This is "search on the answer" (Day 4).

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** Binary search halves the search space every iteration because the predicate is monotonic — once you know `condition(mid)` is true/false, you know the answer can't be on one whole side.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line, especially why the boundary template uses `hi = mid` instead of `hi = mid - 1`.
2. Close this file. Write both templates from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does the boundary template use `lo < hi` while the classic template uses `lo <= hi`?*

   > The classic template is looking for a possibly-absent exact value — when `lo == hi`, there's still one unchecked element, so the loop must run once more (`<=`). The boundary template is converging `lo` and `hi` onto the same answer index — once they're equal, that index *is* the answer, so nothing left to check (`<`). Using `<=` in the boundary template with `hi = mid` risks an infinite loop when `mid == lo == hi - 0`.

---

## Day 2 — Boundary Search: First/Last Occurrence

**Pattern:** Find the leftmost or rightmost index satisfying a condition  
**Key insight:** "Find first index where `nums[i] >= target`" and "find last index where `nums[i] <= target`" are the same boundary template with the condition flipped — most boundary problems reduce to one of these two.

**Template adaptation:**

```go
func searchRange(nums []int, target int) []int {
    first := leftmostTrue(0, len(nums), func(i int) bool { return nums[i] >= target })
    if first == len(nums) || nums[first] != target {
        return []int{-1, -1}
    }
    last := leftmostTrue(0, len(nums), func(i int) bool { return nums[i] > target }) - 1
    return []int{first, last}
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Find First and Last Position of Element in Sorted Array | Medium | https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/ |
| 2 | Search Insert Position | Easy | https://leetcode.com/problems/search-insert-position/ |

**What to notice — Problem 1:** Two boundary searches, not one classic search — `nums[i] >= target` gives the first occurrence, `nums[i] > target` minus one gives the last. Note `hi` starts at `len(nums)` (not `len(nums)-1`) because the answer might be "off the end" (target larger than everything).

**What to notice — Problem 2:** This IS `leftmostTrue` directly — "first index where `nums[i] >= target`" is exactly the insert position, whether or not `target` is present.

---

## Day 3 — Rotated Sorted Array

**Pattern:** Binary search where only half the array is sorted at any given `mid`  
**Key insight:** At each step, one half (`lo` to `mid` or `mid` to `hi`) is always properly sorted. Check which half is sorted first, then decide whether `target` falls inside that sorted half.

**Template adaptation:**

```go
func search(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return mid
        }

        if nums[lo] <= nums[mid] { // left half is sorted
            if nums[lo] <= target && target < nums[mid] {
                hi = mid - 1
            } else {
                lo = mid + 1
            }
        } else { // right half is sorted
            if nums[mid] < target && target <= nums[hi] {
                lo = mid + 1
            } else {
                hi = mid - 1
            }
        }
    }
    return -1
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Search in Rotated Sorted Array | Medium | https://leetcode.com/problems/search-in-rotated-sorted-array/ |
| 2 | Find Minimum in Rotated Sorted Array | Medium | https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/ |

**What to notice — Problem 1:** `nums[lo] <= nums[mid]` is the check that identifies the sorted half — if the left half were rotated, `nums[lo]` would be greater than `nums[mid]`. Once you know which half is sorted, a normal range check (`nums[lo] <= target < nums[mid]`) tells you whether to search there.

**What to notice — Problem 2:** Simpler than Problem 1 — you're not searching for a value, you're finding the rotation point itself. Compare `nums[mid]` to `nums[hi]`: if `nums[mid] > nums[hi]`, the minimum is to the right (`lo = mid + 1`); otherwise it's at `mid` or to the left (`hi = mid`). This is a boundary search in disguise.

---

## Day 4 — Search on the Answer

**Pattern:** Binary search over a range of possible answers, not over the input array  
**Key insight:** Define `condition(x)` = "is `x` a feasible answer?" over the range `[minPossible, maxPossible]`. If feasibility is monotonic (once feasible, stays feasible as `x` grows), binary search finds the smallest feasible `x`.

**Template adaptation:**

```go
func minEatingSpeed(piles []int, h int) int {
    hoursNeeded := func(speed int) int {
        hours := 0
        for _, p := range piles {
            hours += (p + speed - 1) / speed // ceil division
        }
        return hours
    }

    lo, hi := 1, slices.Max(piles) // speed 1 is always feasible eventually; max(piles) finishes any single pile in 1 hour
    for lo < hi {
        mid := lo + (hi-lo)/2
        if hoursNeeded(mid) <= h { // feasible — try slower
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Koko Eating Bananas | Medium | https://leetcode.com/problems/koko-eating-bananas/ |
| 2 | Capacity To Ship Packages Within D Days | Medium | https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/ |

**What to notice — Problem 1:** The search space is "possible eating speeds" (`1` to `max(piles)`), not the `piles` array. `condition(speed)` = "can Koko finish within `h` hours at this speed" — monotonic because a faster speed never needs more hours. This is the exact `leftmostTrue` template with a custom feasibility function.

**What to notice — Problem 2:** Same shape as Problem 1: search space is "possible ship capacities" (`max(weights)` to `sum(weights)`), `condition(capacity)` = "can all packages ship within `D` days at this capacity." Once you can name the search space and write `condition`, the binary search code is identical to Problem 1's.

---

## Day 5 — 2D and Peak Search

**Pattern:** Binary search over a matrix, or over an array with no single sorted order  
**Key insight:** A row-and-column sorted matrix can be treated as one flattened sorted array via index math. A "peak" doesn't require full sorting — only a local monotonic slope at `mid`.

**Template adaptation:**

```go
func searchMatrix(matrix [][]int, target int) bool {
    rows, cols := len(matrix), len(matrix[0])
    lo, hi := 0, rows*cols-1

    for lo <= hi {
        mid := lo + (hi-lo)/2
        val := matrix[mid/cols][mid%cols] // flatten: row = mid/cols, col = mid%cols
        switch {
        case val == target:
            return true
        case val < target:
            lo = mid + 1
        default:
            hi = mid - 1
        }
    }
    return false
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Search a 2D Matrix | Medium | https://leetcode.com/problems/search-a-2d-matrix/ |
| 2 | Find Peak Element | Medium | https://leetcode.com/problems/find-peak-element/ |

**What to notice — Problem 1:** Only works because each row's last element is smaller than the next row's first — that's what makes the matrix equivalent to one flat sorted array. `mid/cols, mid%cols` converts a flat index back to 2D coordinates.

**What to notice — Problem 2:** The array isn't sorted, but at any `mid`, comparing `nums[mid]` to `nums[mid+1]` tells you which direction has a peak: if `nums[mid] < nums[mid+1]`, a peak must exist to the right (slope going up), so `lo = mid + 1`; otherwise `hi = mid`. You never need to look at the whole array — just the local slope.

---

## Day 6 — Advanced Search on the Answer

**Pattern:** Search-on-answer where `condition` itself requires a greedy or DP pass  
**Key insight:** Same template as Day 4, but `condition(x)` is now more expensive to evaluate (a full greedy simulation, or a merge across two arrays). The binary search shape doesn't change — only the cost of checking feasibility.

**Template adaptation:**

```go
func splitArray(nums []int, k int) int {
    canSplit := func(maxSum int) bool {
        pieces, sum := 1, 0
        for _, n := range nums {
            if sum+n > maxSum {
                pieces++
                sum = 0
            }
            sum += n
        }
        return pieces <= k
    }

    lo, hi := slices.Max(nums), sum(nums) // lo: 1 piece per element; hi: everything in 1 piece
    for lo < hi {
        mid := lo + (hi-lo)/2
        if canSplit(mid) { // feasible with maxSum = mid — try smaller
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}

func sum(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Split Array Largest Sum | Hard | https://leetcode.com/problems/split-array-largest-sum/ |
| 2 | Median of Two Sorted Arrays | Hard | https://leetcode.com/problems/median-of-two-sorted-arrays/ |

**What to notice — Problem 1:** The search space is "possible values for the largest subarray sum" (`max(nums)` to `sum(nums)`). `canSplit(maxSum)` greedily counts how many pieces are needed if no piece may exceed `maxSum` — a smaller `maxSum` never needs fewer pieces, so feasibility is monotonic.

**What to notice — Problem 2:** The hardest classic binary search problem — instead of searching on the answer's *value*, you binary search on the **partition point** of the smaller array. A correct partition splits both arrays so every left element ≤ every right element; that condition is monotonic in the partition index, which is what makes binary search applicable at all. Budget extra time here — get the O(n) two-pointer merge working first if the binary-search partition logic doesn't click immediately.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Every binary search problem reduces to naming the search space and writing one `condition` function. The loop logic itself never changes.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing code, say out loud: *What is the search space, and what does `condition(mid)` mean here?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the classic binary search template from memory
- [ ] Write the boundary (`leftmostTrue`) template from memory
- [ ] Why does mixing `lo <= hi` with `hi = mid` risk an infinite loop?
- [ ] How do you know which half of a rotated sorted array is properly sorted?
- [ ] What makes a problem "search on the answer" instead of search on the input array?

---

## Pattern Cheatsheet

| Pattern | Search space | `condition` / comparison | Template |
|---------|--------------|---------------------------|----------|
| Classic exact match | the array itself | `nums[mid] == target` | inclusive (`lo <= hi`) |
| First/last occurrence | the array itself | `nums[i] >= target` (or `>`) | half-open boundary |
| Rotated array | the array itself | which half is sorted, then range check | inclusive (`lo <= hi`) |
| Search on answer | range of possible answers | feasibility function, monotonic | half-open boundary |
| 2D matrix | flattened index space | `mid/cols, mid%cols` | inclusive (`lo <= hi`) |
| Peak / local slope | the array itself | compare `nums[mid]` to `nums[mid+1]` | half-open boundary |

---

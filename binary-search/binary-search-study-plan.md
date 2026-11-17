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

**Brute force → binary search:** a linear scan checks every candidate: O(n). Binary search asks one yes/no question at `mid` and throws away **half** of the remaining candidates with each answer, so after `k` questions only `n / 2ᵏ` are left. That reaches 1 after `log₂ n` steps — 20 steps for a million elements, 30 for a billion.

**The mental model that makes every variant easy:** picture the search space as a row of answers to `condition(i)`:

```
index:      0  1  2  3  4  5  6
condition:  F  F  F  T  T  T  T
                     ^ leftmostTrue returns this index
```

Every binary search problem is "find where F turns into T". Once you can write `condition` and say it's monotonic (all F's before all T's), the loop is always the same.

**Run it:** every function in this file is real code in [`binarysearch.go`](binarysearch.go), tested in [`binarysearch_test.go`](binarysearch_test.go). From the repo root: `go test ./binary-search/ -v`.

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

   > The classic template is looking for a possibly-absent exact value — when `lo == hi`, there's still one unchecked element, so the loop must run once more (`<=`). The boundary template is converging `lo` and `hi` onto the same answer index — once they're equal, that index *is* the answer, so nothing left to check (`<`). Using `<=` in the boundary template with `hi = mid` loops forever: once `lo == hi`, `mid` equals both, `condition(mid)` is true, `hi = mid` changes nothing, and `lo <= hi` is still true.

**Dry run — classic**, `[-1, 0, 3, 5, 9, 12]`:

Target `9` (present):

| `lo` | `hi` | `mid` | `nums[mid]` | Action |
|---|---|---|---|---|
| 0 | 5 | 2 | 3 | too small → `lo = 3` |
| 3 | 5 | 4 | 9 | found → return 4 |

Target `2` (absent):

| `lo` | `hi` | `mid` | `nums[mid]` | Action |
|---|---|---|---|---|
| 0 | 5 | 2 | 3 | too big → `hi = 1` |
| 0 | 1 | 0 | -1 | too small → `lo = 1` |
| 1 | 1 | 1 | 0 | too small → `lo = 2` |
| 2 | 1 | — | — | `lo > hi` → return -1 |

Look at the third row: `lo == hi == 1` and that element still had to be checked — that's why the classic loop is `lo <= hi`. When it ends, `lo` (2) is where `2` *would* be inserted.

**Dry run — boundary**, first index where `nums[i] >= 8` in `[5, 7, 7, 8, 8, 10]`:

| `lo` | `hi` | `mid` | `condition(mid)` | Action |
|---|---|---|---|---|
| 0 | 6 | 3 | `8 >= 8` | true → `hi = 3` |
| 0 | 3 | 1 | `7 >= 8` | false → `lo = 2` |
| 2 | 3 | 2 | `7 >= 8` | false → `lo = 3` |
| 3 | 3 | — | — | `lo == hi` → return 3 |

Invariant the whole way: the answer is always inside `[lo, hi]`. A `true` keeps `mid` (it might be the first true); a `false` excludes it.

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

**Why two boundary searches:** in a sorted array, the copies of `target` form one contiguous block. Its **first** index is the first `i` with `nums[i] >= target`. Its **last** index is one before the first `i` with `nums[i] > target`. Both conditions are monotonic, so both are `leftmostTrue`.

**Dry run** — the second search on `[5, 7, 7, 8, 8, 10]`, first index where `nums[i] > 8`:

| `lo` | `hi` | `mid` | `condition(mid)` | Action |
|---|---|---|---|---|
| 0 | 6 | 3 | `8 > 8` | false → `lo = 4` |
| 4 | 6 | 5 | `10 > 8` | true → `hi = 5` |
| 4 | 5 | 4 | `8 > 8` | false → `lo = 5` |
| 5 | 5 | — | — | `lo == hi` → return 5 |

Last occurrence = 5 − 1 = **4**. Combined with the Day 1 boundary run (first = 3), the answer is `[3, 4]`.

**Run it:** `leftmostTrue` and `searchRange` in [`binarysearch.go`](binarysearch.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Find First and Last Position of Element in Sorted Array | Medium | https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/ |
| 2 | Search Insert Position | Easy | https://leetcode.com/problems/search-insert-position/ |

**What to notice — Problem 1:** Two boundary searches, not one classic search — `nums[i] >= target` gives the first occurrence, `nums[i] > target` minus one gives the last. Note `hi` starts at `len(nums)` (not `len(nums)-1`) because the answer might be "off the end" (target larger than everything).

**What to notice — Problem 2:** This IS `leftmostTrue` directly — "first index where `nums[i] >= target`" is exactly the insert position, whether or not `target` is present.

**Complexity:** both problems O(log n) time, O(1) space. Two binary searches are still O(log n).

**Common mistakes:**

- Starting `hi` at `len(nums) - 1` in the boundary template. If the answer is "past the end" (target bigger than everything), `len(nums)` must be a possible return value.
- Checking `nums[first] != target` without first checking `first == len(nums)` — index out of range.
- Finding the first occurrence and then scanning right for the last. That's O(n) on `[8, 8, 8, …, 8]`.

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

**Why one half is always sorted:** a rotated sorted array is two sorted runs glued together, with a single "drop" where the big values end. `mid` splits the range in two, and the drop can only be in one of those halves — so the other half is a normal sorted run. `nums[lo] <= nums[mid]` means no drop between them → the left half is sorted. Then a plain range check on the sorted half tells you whether the target can be there.

**Dry run** — `[4, 5, 6, 7, 0, 1, 2]`, target `0`:

| `lo` | `hi` | `mid` | `nums[mid]` | Sorted half | Action |
|---|---|---|---|---|---|
| 0 | 6 | 3 | 7 | left `[4..7]` | 0 not in `[4, 7)` → `lo = 4` |
| 4 | 6 | 5 | 1 | left `[0..1]` | 0 in `[0, 1)` → `hi = 4` |
| 4 | 4 | 4 | 0 | — | found → return 4 |

**Run it:** `searchRotated` in [`binarysearch.go`](binarysearch.go) (called `search` above, as on LeetCode).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Search in Rotated Sorted Array | Medium | https://leetcode.com/problems/search-in-rotated-sorted-array/ |
| 2 | Find Minimum in Rotated Sorted Array | Medium | https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/ |

**What to notice — Problem 1:** `nums[lo] <= nums[mid]` is the check that identifies the sorted half — if the left half were rotated, `nums[lo]` would be greater than `nums[mid]`. Once you know which half is sorted, a normal range check (`nums[lo] <= target < nums[mid]`) tells you whether to search there.

**What to notice — Problem 2:** Simpler than Problem 1 — you're not searching for a value, you're finding the rotation point itself. Compare `nums[mid]` to `nums[hi]`: if `nums[mid] > nums[hi]`, the minimum is to the right (`lo = mid + 1`); otherwise it's at `mid` or to the left (`hi = mid`). This is a boundary search in disguise.

**Complexity:** both problems O(log n) time, O(1) space — for distinct values. With duplicates (the "II" versions), `nums[lo] == nums[mid] == nums[hi]` gives no information and you have to shrink by one, so the worst case becomes O(n).

**Common mistakes:**

- `nums[lo] < nums[mid]` instead of `<=`. When `lo == mid` (two elements left), the "left half" is just one element and is sorted — strict `<` sends you the wrong way.
- Range checks with the wrong strictness: `nums[lo] <= target < nums[mid]` (mid was already checked, so it's excluded).
- Find Minimum: comparing `nums[mid]` with `nums[lo]` instead of `nums[hi]`. On an un-rotated array that comparison can't tell which side the minimum is on.

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

    lo, hi := 1, slices.Max(piles) // lo: slowest possible speed; hi: always feasible (every pile takes 1 hour, and h >= len(piles))
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

**Why it works:** hours needed only goes **down** (or stays the same) as speed goes up. So over the speeds `1, 2, 3, …`, the question "can she finish in `h` hours?" is a row of F's followed by T's — exactly the shape from the top of this file. You're looking for the first T.

Hours needed for `piles = [3, 6, 7, 11]` at each speed (`h = 8`):

| speed | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| hours | 27 | 15 | 10 | **8** | 8 | 6 | 5 | 5 | 5 | 5 | 4 |
| ≤ 8? | F | F | F | **T** | T | T | T | T | T | T | T |

**Dry run** — binary search finds that first T without computing the whole row:

| `lo` | `hi` | `mid` | `hoursNeeded(mid)` | Feasible? |
|---|---|---|---|---|
| 1 | 11 | 6 | 6 | yes → `hi = 6` |
| 1 | 6 | 3 | 10 | no → `lo = 4` |
| 4 | 6 | 5 | 8 | yes → `hi = 5` |
| 4 | 5 | 4 | 8 | yes → `hi = 4` |

`lo == hi == 4` → answer **4**, after evaluating only 4 of the 11 speeds.

**Run it:** `minEatingSpeed` in [`binarysearch.go`](binarysearch.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Koko Eating Bananas | Medium | https://leetcode.com/problems/koko-eating-bananas/ |
| 2 | Capacity To Ship Packages Within D Days | Medium | https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/ |

**What to notice — Problem 1:** The search space is "possible eating speeds" (`1` to `max(piles)`), not the `piles` array. `condition(speed)` = "can Koko finish within `h` hours at this speed" — monotonic because a faster speed never needs more hours. This is the exact `leftmostTrue` template with a custom feasibility function.

**What to notice — Problem 2:** Same shape as Problem 1: search space is "possible ship capacities" (`max(weights)` to `sum(weights)`), `condition(capacity)` = "can all packages ship within `D` days at this capacity." Once you can name the search space and write `condition`, the binary search code is identical to Problem 1's.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Koko Eating Bananas | O(n log M), M = max pile — log M guesses, O(n) check each | O(1) |
| Capacity to Ship Packages | O(n log S), S = sum of weights | O(1) |

**Common mistakes:**

- Ship Packages: setting `lo = 1` instead of `max(weights)`. A capacity smaller than the heaviest package can never work, and the greedy check may loop or miscount.
- Ceiling division written as `p / speed + 1` — wrong when `p` divides evenly. Use `(p + speed - 1) / speed`.
- Not convincing yourself the condition is monotonic before coding. If it isn't, binary search gives a wrong answer silently.

---

## Day 5 — 2D and Peak Search

**Pattern:** Binary search over a matrix, or over an array with no single sorted order  
**Key insight:** A matrix whose rows are sorted *and* where each row starts after the previous row ends can be treated as one flattened sorted array via index math. (A matrix that is only row-sorted and column-sorted can't be flattened — that's Search a 2D Matrix II, solved by a staircase walk from the top-right corner.) A "peak" doesn't require full sorting — only a local monotonic slope at `mid`.

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

**Why flattening works:** read the matrix row by row and you get one sorted list of `rows × cols` numbers. Flat index `k` lives in row `k / cols`, column `k % cols` (each row holds `cols` numbers). So run the classic template over `0 … rows×cols − 1` and convert `mid` back to 2D only when you need the value.

**Dry run** — `[[1,3,5,7],[10,11,16,20],[23,30,34,60]]`, target `16`, `cols = 4`:

| `lo` | `hi` | `mid` | `(mid/4, mid%4)` | Value | Action |
|---|---|---|---|---|---|
| 0 | 11 | 5 | (1, 1) | 11 | `lo = 6` |
| 6 | 11 | 8 | (2, 0) | 23 | `hi = 7` |
| 6 | 7 | 6 | (1, 2) | 16 | found |

**Why Find Peak works:** if `nums[mid] < nums[mid+1]`, you're on an upward slope. Keep walking right: either it keeps rising until the last element (a peak, since `nums[n] = -∞`) or it turns down somewhere (a peak). Either way, a peak exists to the right, so discarding the left is safe.

**Run it:** `searchMatrix` in [`binarysearch.go`](binarysearch.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Search a 2D Matrix | Medium | https://leetcode.com/problems/search-a-2d-matrix/ |
| 2 | Find Peak Element | Medium | https://leetcode.com/problems/find-peak-element/ |

**What to notice — Problem 1:** Only works because each row's last element is smaller than the next row's first — that's what makes the matrix equivalent to one flat sorted array. `mid/cols, mid%cols` converts a flat index back to 2D coordinates.

**What to notice — Problem 2:** The array isn't sorted, but at any `mid`, comparing `nums[mid]` to `nums[mid+1]` tells you which direction has a peak: if `nums[mid] < nums[mid+1]`, a peak must exist to the right (slope going up), so `lo = mid + 1`; otherwise `hi = mid`. You never need to look at the whole array — just the local slope.

**Complexity:** Search a 2D Matrix O(log(m·n)); Find Peak O(log n). Both O(1) space.

**Common mistakes:**

- Converting with `mid / rows` instead of `mid / cols`.
- Find Peak with `lo <= hi` and `hi = mid` — the Day 1 infinite loop. It's a boundary-style search; use `lo < hi`.
- Find Peak: reading `nums[mid+1]` when `mid` is the last index. With `lo < hi`, `mid < hi` always, so `mid+1` is safe — another reason to stick to the boundary template.

---

## Day 6 — Advanced Search on the Answer

**Pattern:** Search-on-answer where `condition` itself requires a greedy or DP pass  
**Key insight:** Same template as Day 4, but `condition(x)` is now more expensive to evaluate (a full greedy simulation, or a merge across two arrays). The binary search shape doesn't change — only the cost of checking feasibility.

**Template adaptation:**

```go
func splitArray(nums []int, k int) int {
    canSplit := func(maxSum int) bool {
        pieces, cur := 1, 0
        for _, n := range nums {
            if cur+n > maxSum { // n doesn't fit in the current piece — start a new one
                pieces++
                cur = 0
            }
            cur += n
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

**Why it works (Split Array):** if every piece may sum to at most `x`, a larger `x` never needs **more** pieces. So "can I split into ≤ k pieces with max sum `x`?" is F…F T…T over `x` — search on the answer again. The greedy check (start a new piece only when the next number doesn't fit) uses the fewest pieces possible for a given `x`, so it answers the question exactly.

**Dry run** — `splitArray([7, 2, 5, 10, 8], 2)`, search space `[10, 32]`:

| `lo` | `hi` | `mid` | Greedy pieces | Count | ≤ 2? |
|---|---|---|---|---|---|
| 10 | 32 | 21 | [7, 2, 5] [10, 8] | 2 | yes → `hi = 21` |
| 10 | 21 | 15 | [7, 2, 5] [10] [8] | 3 | no → `lo = 16` |
| 16 | 21 | 18 | [7, 2, 5] [10, 8] | 2 | yes → `hi = 18` |
| 16 | 18 | 17 | [7, 2, 5] [10] [8] | 3 | no → `lo = 18` |

Answer **18**: at 17 the pieces `10` and `8` no longer fit together.

**Run it:** `splitArray` in [`binarysearch.go`](binarysearch.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Split Array Largest Sum | Hard | https://leetcode.com/problems/split-array-largest-sum/ |
| 2 | Median of Two Sorted Arrays | Hard | https://leetcode.com/problems/median-of-two-sorted-arrays/ |

**What to notice — Problem 1:** The search space is "possible values for the largest subarray sum" (`max(nums)` to `sum(nums)`). `canSplit(maxSum)` greedily counts how many pieces are needed if no piece may exceed `maxSum` — a smaller `maxSum` never needs fewer pieces, so feasibility is monotonic.

**What to notice — Problem 2:** The hardest classic binary search problem — instead of searching on the answer's *value*, you binary search on the **partition point** of the smaller array. A correct partition splits both arrays so every left element ≤ every right element; that condition is monotonic in the partition index, which is what makes binary search applicable at all. Budget extra time here — get the O(n) two-pointer merge working first if the binary-search partition logic doesn't click immediately.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Split Array Largest Sum | O(n log S), S = sum(nums) | O(1) |
| Median of Two Sorted Arrays | O(log(min(m, n))) — binary search on the shorter array's partition | O(1) |

**Common mistakes:**

- Split Array: `lo = 0` or `lo = 1`. Any piece holds at least one element, so the max sum is at least `max(nums)`.
- Median: binary searching the longer array — the partition index of the other array can go out of range. Always search the shorter one.
- Median: forgetting the `-∞` / `+∞` sentinels when a partition is at an array edge.

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

All patterns are O(log n) iterations; search-on-answer multiplies that by the cost of one `condition` check.

| Pattern | Search space | `condition` / comparison | Template |
|---------|--------------|---------------------------|----------|
| Classic exact match | the array itself | `nums[mid] == target` | inclusive (`lo <= hi`) |
| First/last occurrence | the array itself | `nums[i] >= target` (or `>`) | half-open boundary |
| Rotated array | the array itself | which half is sorted, then range check | inclusive (`lo <= hi`) |
| Search on answer | range of possible answers | feasibility function, monotonic | half-open boundary |
| 2D matrix | flattened index space | `mid/cols, mid%cols` | inclusive (`lo <= hi`) |
| Peak / local slope | the array itself | compare `nums[mid]` to `nums[mid+1]` | half-open boundary |

---

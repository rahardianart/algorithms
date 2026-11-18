# Prefix Sum & Hashing Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 10 — pairs naturally with the Sliding Window week)  
**Goal:** Trade memory for time with hash maps, answer range-sum questions in O(1), and combine the two to count subarrays — including the negative-number cases sliding windows can't handle.

---

## What Are Prefix Sums and Hashing?

A **prefix sum** array stores, at position `i`, the total of everything before `i`. With it, the sum of **any range** `nums[l..r]` is a single subtraction — `prefix[r+1] - prefix[l]` — instead of a loop.

**Hashing** means using a hash map (Go's `map`) to remember what you've seen so you can answer "have I seen X before?" in O(1). Combined, they find subarrays with a given sum in one pass — even when the numbers are negative, where a sliding window fails.

**Typical problems it solves:**

- Two Sum and other "find the partner" problems
- Grouping items that belong together (group anagrams)
- Range sum queries, product of array except self
- Counting subarrays with sum k, or sum divisible by k
- Many range updates at once (difference arrays), 2D rectangle sums

---

## The Two Templates

### Hash Map: "have I seen what I need?"

```go
seen := map[int]int{} // value -> count (or index) of what came before
for _, x := range nums {
    need := target - x
    if seen[need] > 0 {
        // x pairs with an earlier element
    }
    seen[x]++ // record x AFTER checking, so it can't pair with itself
}
```

### Prefix Sum: "sum of any range in O(1)"

```go
prefix := make([]int, len(nums)+1) // prefix[i] = sum of nums[:i]
for i, x := range nums {
    prefix[i+1] = prefix[i] + x
}

// sum of nums[l..r], inclusive
rangeSum := prefix[r+1] - prefix[l]
```

**Three rules — never break them:**

1. **Check the map, then update it.** Looking up before inserting means every match is with an *earlier* element — no self-pairs, no double counting.
2. **`prefix` has `n + 1` entries and `prefix[0] = 0`.** The extra slot represents the empty prefix. It removes every `if l == 0` special case and is what makes "subarray starting at index 0" work in the hash map patterns.
3. **A subarray sum is a difference of two prefix sums: `sum(l..r) = prefix[r+1] - prefix[l]`.** Every pattern this week is that one equation, rearranged.

**Brute force → hashing / prefix sums:** "find a pair with sum X" is O(n²) over all pairs; a hash map answers "has the partner appeared?" in O(1), so one pass is O(n). "Sum of `nums[l..r]`" is O(n) per query by looping; one O(n) prefix pass makes every later query O(1). Combine them and "count subarrays with sum k" goes from O(n²) to O(n) — **even with negative numbers**, which is exactly where the sliding window fails.

**Complexity:** O(n) time, O(n) space for both templates.

**Run it:** every template in this file is a real function in [`prefixsum.go`](prefixsum.go), tested in [`prefixsum_test.go`](prefixsum_test.go). From the repo root: `go test ./prefix-sum/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** A hash map turns "search everything before me" into one O(1) lookup. A prefix array turns "add up a range" into one subtraction.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both templates from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does `prefix` have `n + 1` entries instead of `n`?*

   > So that every range has a "left edge" to subtract. With `prefix[i]` = sum of the first `i` elements, `sum(l..r) = prefix[r+1] - prefix[l]` works for every `l`, including `l = 0`, because `prefix[0] = 0` is the sum of nothing. With only `n` entries you'd need a special case for ranges that start at index 0 — and in the hash-map patterns (Day 5), forgetting that case silently misses every subarray that starts at the beginning.

**Dry run** — prefix sums of `[3, -1, 4, 1, 5]`:

| `i` | 0 | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|---|
| `prefix[i]` | 0 | 3 | 2 | 6 | 7 | 12 |

`sum(1..3) = prefix[4] - prefix[1] = 7 - 3 = 4`, and indeed `-1 + 4 + 1 = 4`. Negative numbers don't matter — it's just subtraction.

---

## Day 2 — Complement Lookup

**Pattern:** For each element, look up its "partner" among the elements seen so far  
**Key insight:** Instead of asking "which later element pairs with me?", ask "which *earlier* element pairs with me?" — and keep the earlier ones in a map.

**Template adaptation:**

```go
// count index pairs i < j with nums[i] + nums[j] == target
func countPairsWithSum(nums []int, target int) int {
    seen := map[int]int{} // value -> how many times it appeared so far
    count := 0
    for _, x := range nums {
        count += seen[target-x] // every earlier partner forms one pair
        seen[x]++
    }
    return count
}
```

**Why it works:** every pair `(i, j)` with `i < j` is counted exactly once — at the moment `j` is processed, when `i` is already in `seen`. It can't be counted at `i`, because `j` isn't in the map yet. That's rule 1.

**Dry run** — `countPairsWithSum([1, 5, 3, 3, 7, 1], 6)`:

| `x` | Looking for `6 - x` | Seen before (count) | Total | `seen` after |
|---|---|---|---|---|
| 1 | 5 | 0 | 0 | `{1: 1}` |
| 5 | 1 | 1 | 1 | `{1: 1, 5: 1}` |
| 3 | 3 | 0 | 1 | `{1: 1, 5: 1, 3: 1}` |
| 3 | 3 | 1 | 2 | `{1: 1, 5: 1, 3: 2}` |
| 7 | -1 | 0 | 2 | `{…, 7: 1}` |
| 1 | 5 | 1 | 3 | `{1: 2, …}` |

At the first `3`, the lookup for `3` finds nothing because that `3` isn't inserted yet — so it can't pair with itself.

**Run it:** `countPairsWithSum` in [`prefixsum.go`](prefixsum.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Two Sum | Easy | https://leetcode.com/problems/two-sum/ |
| 2 | Contains Duplicate II | Easy | https://leetcode.com/problems/contains-duplicate-ii/ |

**What to notice — Problem 1:** Same shape, but the map stores **value → index** (you need to return indices), and you return on the first match instead of counting.

**What to notice — Problem 2:** The "partner" of `x` is `x` itself, and the condition is about distance: map value → *last* index seen. Overwrite the index every time — the most recent occurrence is always the closest one for later elements.

**Complexity:** both O(n) time, O(n) space.

**Common mistakes:**

- Inserting before looking up — `x` pairs with itself when `target == 2x`.
- Two Sum: sorting first to use two pointers — it works for the values, but sorting loses the original indices the problem asks for.
- Contains Duplicate II: storing the *first* index instead of the last.

---

## Day 3 — Grouping by a Canonical Key

**Pattern:** Map every item to a key that's equal exactly when items "belong together", then bucket by key  
**Key insight:** The hard part is designing the key. Equal keys must mean "same group" — no more, no less.

**Template adaptation:**

```go
// group words that use exactly the same SET of letters ("abc", "cab", "aabbc")
func groupByCharSet(words []string) [][]string {
    groupOf := map[uint32]int{} // letter-set bitmask -> index into groups
    var groups [][]string
    for _, w := range words {
        var key uint32
        for i := 0; i < len(w); i++ {
            key |= 1 << (w[i] - 'a') // bit c is set if letter c appears
        }
        idx, ok := groupOf[key]
        if !ok {
            idx = len(groups)
            groupOf[key] = idx
            groups = append(groups, nil)
        }
        groups[idx] = append(groups[idx], w)
    }
    return groups
}
```

**Why the key works:** a 26-bit mask records *which* letters appear, ignoring order and repetition — which is exactly the definition of "same set of letters". Two words get the same mask if and only if they belong in the same group. Change the definition, change the key: anagrams care about *how many* of each letter, so a bitmask is too coarse for them.

**Dry run** — `groupByCharSet(["abc", "cab", "ab", "aabbc", "ba", "xyz"])`:

| Word | Letters | Key (bits) | Group |
|---|---|---|---|
| `abc` | {a, b, c} | `…0111` | new group 0 |
| `cab` | {a, b, c} | `…0111` | 0 |
| `ab` | {a, b} | `…0011` | new group 1 |
| `aabbc` | {a, b, c} | `…0111` | 0 |
| `ba` | {a, b} | `…0011` | 1 |
| `xyz` | {x, y, z} | bits 23–25 | new group 2 |

Result: `[[abc cab aabbc] [ab ba] [xyz]]`.

**Run it:** `groupByCharSet` in [`prefixsum.go`](prefixsum.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Group Anagrams | Medium | https://leetcode.com/problems/group-anagrams/ |
| 2 | Longest Consecutive Sequence | Medium | https://leetcode.com/problems/longest-consecutive-sequence/ |

**What to notice — Problem 1:** Same skeleton, different key: a `[26]int` letter count. Go arrays are comparable, so `map[[26]int][]string` works directly — no string building or sorting needed.

**What to notice — Problem 2:** Put everything in a set. Only start counting from numbers that begin a run (`x-1` is not in the set), then walk `x+1, x+2, …`. Each number is walked over at most once, so it's O(n) even though there's a loop inside a loop.

**Complexity:** Problem 1 O(n · L) for n words of length L (sorting each word instead would be O(n · L log L)). Problem 2 O(n) time, O(n) space.

**Common mistakes:**

- Group Anagrams: using the sorted word as a key is fine but slower; using a bitmask is **wrong** (`"aab"` and `"abb"` would collide).
- Longest Consecutive: starting a walk from every number — O(n²) on `[1, 2, …, n]`.
- Longest Consecutive: sorting first — correct but O(n log n), which the problem explicitly rules out.

---

## Day 4 — Prefix Sums for Range Queries

**Pattern:** Precompute once, answer every range query with a subtraction  
**Key insight:** Anything that "accumulates" left to right — sums, products, counts — can be prefixed. Many array problems are a prefix pass from the left combined with one from the right.

**Template adaptation:**

```go
func prefixSums(nums []int) []int {
    prefix := make([]int, len(nums)+1)
    for i, x := range nums {
        prefix[i+1] = prefix[i] + x
    }
    return prefix
}

func rangeSum(prefix []int, l, r int) int {
    return prefix[r+1] - prefix[l]
}
```

**Why it works:** `prefix[r+1]` is the sum of `nums[0..r]` and `prefix[l]` is the sum of `nums[0..l-1]`. Subtract, and everything before `l` cancels out, leaving exactly `nums[l..r]`.

**Run it:** `prefixSums` and `rangeSum` in [`prefixsum.go`](prefixsum.go) — the test checks every `(l, r)` pair against a direct loop (dry run on Day 1).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Range Sum Query - Immutable | Easy | https://leetcode.com/problems/range-sum-query-immutable/ |
| 2 | Product of Array Except Self | Medium | https://leetcode.com/problems/product-of-array-except-self/ |

**What to notice — Problem 1:** The template wrapped in a struct: build `prefix` in the constructor, `rangeSum` in `SumRange`. Building once and querying many times is the whole point.

**What to notice — Problem 2:** No division allowed, so you can't "subtract" a product. Instead: `answer[i]` = (product of everything left of `i`) × (product of everything right of `i`). Fill `answer` with left products in one pass, then multiply in right products with a running variable in a second pass from the right — O(1) extra space.

**Complexity:** Problem 1 O(n) build, O(1) per query. Problem 2 O(n) time, O(1) extra space (the output doesn't count).

**Common mistakes:**

- Sizing `prefix` as `n` and writing `prefix[r] - prefix[l-1]` — breaks at `l = 0`.
- Problem 2: using division and then failing on zeros.
- Recomputing the prefix inside every query (O(n) per query again).

---

## Day 5 — Prefix Sum + Hash Map

**Pattern:** Count or find subarrays with a target sum in one pass, negatives allowed  
**Key insight:** `sum(l..r) = k` ⟺ `prefix[r+1] − prefix[l] = k` ⟺ `prefix[l] = prefix[r+1] − k`. So at each position, ask the map: *"which earlier prefix equals `sum − k`?"* It's Day 2's complement lookup, run over prefix sums.

**Template adaptation:**

```go
// length of the longest subarray summing to k (negatives allowed)
func longestSubarrayWithSum(nums []int, k int) int {
    firstSeen := map[int]int{0: -1} // prefix sum -> earliest index; 0 before the array starts
    sum, best := 0, 0
    for i, x := range nums {
        sum += x
        if j, ok := firstSeen[sum-k]; ok {
            best = max(best, i-j) // subarray (j, i] sums to k
        }
        if _, ok := firstSeen[sum]; !ok {
            firstSeen[sum] = i // keep the EARLIEST index — it gives the longest subarray
        }
    }
    return best
}
```

**Why `{0: -1}`:** it's `prefix[0] = 0` from rule 2, stored in the map. Without it, a subarray that starts at index 0 (like `[1, -1, 5, -2]` below) has no "earlier prefix" to match and is never found.

**Dry run** — `longestSubarrayWithSum([1, -1, 5, -2, 3], 3)`:

| `i` | `x` | prefix `sum` | Look up `sum - 3` | Found at | Length | best | Store `sum`? |
|---|---|---|---|---|---|---|---|
| 0 | 1 | 1 | -2 | — | — | 0 | yes → 0 |
| 1 | -1 | 0 | -3 | — | — | 0 | no (0 already at -1) |
| 2 | 5 | 5 | 2 | — | — | 0 | yes → 2 |
| 3 | -2 | 3 | 0 | -1 | 4 | 4 | yes → 3 |
| 4 | 3 | 6 | 3 | 3 | 1 | 4 | yes → 4 |

Answer **4** (`[1, -1, 5, -2]`), found through the `0: -1` entry. A sliding window could never find it: the window would have to grow past a sum of 5 and come back down.

**Run it:** `longestSubarrayWithSum` in [`prefixsum.go`](prefixsum.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Subarray Sum Equals K | Medium | https://leetcode.com/problems/subarray-sum-equals-k/ |
| 2 | Continuous Subarray Sum | Medium | https://leetcode.com/problems/continuous-subarray-sum/ |

**What to notice — Problem 1:** Counting instead of longest: the map stores prefix sum → **how many times** it occurred, and you add `count[sum-k]` at every step. Initialize with `{0: 1}`.

**What to notice — Problem 2:** "sum is a multiple of k" ⟺ two prefix sums have the **same remainder mod k**. Store remainder → earliest index (with `{0: -1}`), and check the length is at least 2.

**Complexity:** both O(n) time, O(n) space (up to n distinct prefix sums).

**Common mistakes:**

- Forgetting the `{0: 1}` / `{0: -1}` initial entry — misses every subarray starting at index 0.
- Longest: overwriting the stored index on every occurrence (you want the earliest). Counting: forgetting that repeated prefix sums must *add up*, not overwrite.
- Trying a sliding window because "it's a subarray sum problem" — check for negatives first.

---

## Day 6 — Difference Arrays and 2D Prefix Sums

**Pattern:** Prefix sums run in reverse (many range updates, one read) and in two dimensions (many rectangle queries)  
**Key insight:** A difference array is the "derivative" of an array; its prefix sum is the array itself. A range update touches 2 cells of the derivative instead of `r - l + 1` cells of the array.

**Template adaptation:**

```go
// apply many "add delta to nums[l..r]" updates, then read the final array once
func applyRangeAdds(n int, updates [][3]int) []int {
    diff := make([]int, n+1)
    for _, u := range updates {
        l, r, delta := u[0], u[1], u[2]
        diff[l] += delta   // from l onward, values go up by delta…
        diff[r+1] -= delta // …and from r+1 onward, they go back down
    }
    result := make([]int, n)
    running := 0
    for i := 0; i < n; i++ {
        running += diff[i] // prefix sum of the differences
        result[i] = running
    }
    return result
}

// 2D: p[i][j] = sum of the rectangle grid[0..i-1][0..j-1]
p[i+1][j+1] = grid[i][j] + p[i][j+1] + p[i+1][j] - p[i][j]

// sum of grid[r1..r2][c1..c2]
p[r2+1][c2+1] - p[r1][c2+1] - p[r2+1][c1] + p[r1][c1]
```

**Why the 2D formula works (inclusion–exclusion):** the big rectangle up to `(r2, c2)`, minus the strip above the target, minus the strip to its left — but the top-left corner was subtracted **twice**, so add it back once. Building `p` uses the same idea in reverse.

**Dry run** — `applyRangeAdds(5, [{1,3,+2}, {2,4,+3}, {0,2,−2}])`:

| After update | `diff` |
|---|---|
| `[1..3] += 2` | `[0, 2, 0, 0, -2, 0]` |
| `[2..4] += 3` | `[0, 2, 3, 0, -2, -3]` |
| `[0..2] -= 2` | `[-2, 2, 3, 2, -2, -3]` |
| running sum | `[-2, 0, 3, 5, 3]` |

Three updates, two cell writes each — no matter how wide the ranges are.

**Run it:** `applyRangeAdds`, `newPrefix2D`, and `sumRegion` in [`prefixsum.go`](prefixsum.go) — the 2D test uses LeetCode 304's example grid.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Range Sum Query 2D - Immutable | Medium | https://leetcode.com/problems/range-sum-query-2d-immutable/ |
| 2 | Car Pooling | Medium | https://leetcode.com/problems/car-pooling/ |

**What to notice — Problem 1:** The 2D template in a struct. Draw the four rectangles once on paper — the `+ p[r1][c1]` term is the one everyone forgets.

**What to notice — Problem 2:** Each trip is a range update on "passengers in the car at location x": `+numPassengers` at `from`, `−numPassengers` at `to` (they get out *at* `to`, so it's `to`, not `to + 1`). Then a running sum — if it ever exceeds `capacity`, return false.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Range Sum Query 2D | O(R·C) build, O(1) per query | O(R·C) |
| Car Pooling | O(trips + max location) | O(max location) |

**Common mistakes:**

- 2D: missing the `+ p[r1][c1]` correction (corner subtracted twice).
- Difference array: sizing `diff` as `n` — `diff[r+1]` overflows when `r = n-1`.
- Car Pooling: subtracting at `to + 1` — passengers leave at `to`, so the seat is free for someone boarding at `to`.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Most problems this week are one equation — `sum(l..r) = prefix[r+1] − prefix[l]` — plus a decision about what to store in the map: a count, the earliest index, or the latest index.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Answer out loud: *Given "count subarrays with sum k", how do you decide between a sliding window and prefix sum + hash map?*

   > Check whether validity is monotonic. If every number is non-negative and the condition is "≤" or "≥", a sliding window works in O(1) space. If numbers can be negative, or the condition is "== k", growing the window can make it valid again later — use prefix sum + hash map.

**Self-assessment — can you answer these without hesitation?**

- [ ] Why look up before inserting in the complement pattern?
- [ ] Why does `prefix` have `n + 1` entries?
- [ ] What does the map store for "count subarrays", "longest subarray", and "nearest duplicate"?
- [ ] Why is `{0: -1}` (or `{0: 1}`) in the map before the loop starts?
- [ ] How does a difference array turn an O(n) range update into O(1)?
- [ ] Write the 2D range-sum formula and explain each of the four terms

---

## Pattern Cheatsheet

| Pattern | Data structure | Map stores | Key move | Time |
|---------|----------------|------------|----------|------|
| Complement lookup | `map[value]…` | count / index of earlier values | look up `target - x`, then insert `x` | O(n) |
| Canonical key | `map[key][]item` | the group | design a key equal iff "same group" | O(n · key cost) |
| Range query | `prefix []int` (n+1) | — | `prefix[r+1] - prefix[l]` | O(n) build, O(1) query |
| Prefix + hash map | `map[prefixSum]…` | count, or earliest index | look up `sum - k`; seed with prefix 0 | O(n) |
| Difference array | `diff []int` (n+1) | — | `diff[l] += d; diff[r+1] -= d`, then prefix-sum | O(n + updates) |
| 2D prefix | `p [][]int` (R+1 × C+1) | — | inclusion–exclusion, 4 terms | O(R·C) build, O(1) query |

---

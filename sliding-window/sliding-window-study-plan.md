# Sliding Window Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 9 — builds directly on the Two Pointer week)  
**Goal:** Write the fixed and variable window templates from memory, know which way to shrink for "longest" vs "shortest", and recognize when a window *can't* work.

---

## The Two Templates

### Fixed Size (window of exactly `k`)

```go
sum, best := 0, math.MinInt
for right := 0; right < len(nums); right++ {
    sum += nums[right] // add the element entering on the right
    if right >= k {
        sum -= nums[right-k] // drop the element leaving on the left
    }
    if right >= k-1 { // window [right-k+1, right] is full
        best = max(best, sum)
    }
}
```

### Variable Size (grow right, shrink left)

```go
left := 0
for right := 0; right < len(nums); right++ {
    add(nums[right]) // update window state for the new element

    for windowIsInvalid() {
        remove(nums[left]) // undo left's contribution
        left++
    }

    best = max(best, right-left+1) // window [left, right] is valid here
}
```

**Three rules — never break them:**

1. **Both pointers only move forward.** `right` advances once per iteration; `left` only ever increases, never resets. That's what makes it O(n): each element is added once and removed at most once.
2. **Window state must update in O(1).** Keep a running sum, a count map, or a "number of bad elements" counter — never re-scan the window. If you can't update the state in O(1) on add/remove, it's not a sliding window problem yet.
3. **Longest → shrink while *invalid*, record after the loop. Shortest → shrink while *valid*, record inside the loop.** Mixing these up is the #1 sliding-window bug.

**Brute force → sliding window:** "best subarray/substring with property X" brute-forces by checking every `(left, right)` pair: O(n²) windows, often with O(n) work each. A sliding window notices that neighboring windows share almost everything — moving one end changes one element — so it updates the answer instead of recomputing it.

**When a window can't work:** the variable template needs **monotonic validity** — if a window is valid, every smaller window inside it is valid too (or: if invalid, every bigger window containing it is invalid). "Sum ≤ limit with non-negative numbers" has it; "sum == k with negative numbers" does not, because adding a negative number can make an invalid window valid again. That case belongs to prefix sums + hashing (next week).

**Complexity:** O(n) time for both templates. Space is O(1), or O(alphabet) / O(distinct values) when the state is a count map.

**Run it:** every template in this file is a real function in [`slidingwindow.go`](slidingwindow.go), tested in [`slidingwindow_test.go`](slidingwindow_test.go). From the repo root: `go test ./sliding-window/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** A window is two pointers that bound a *range* instead of a pair — what matters is the state of everything between them.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both templates from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *The variable template has a `for` loop inside a `for` loop. Why is it still O(n)?*

   > Count work per element instead of per iteration. Each element is added exactly once (when `right` reaches it) and removed at most once (when `left` passes it). `left` never moves backward, so the inner loop runs at most n times **in total** across the whole run, not n times per outer iteration. n adds + at most n removes = O(n).

**Dry run** — fixed window, `maxSumFixed([2, 1, 5, 1, 3, 2], 3)`:

| `right` | add | drop | `sum` | Window | best |
|---|---|---|---|---|---|
| 0 | 2 | — | 2 | (filling) | — |
| 1 | 1 | — | 3 | (filling) | — |
| 2 | 5 | — | 8 | `[2, 1, 5]` | 8 |
| 3 | 1 | 2 | 7 | `[1, 5, 1]` | 8 |
| 4 | 3 | 1 | 9 | `[5, 1, 3]` | 9 |
| 5 | 2 | 5 | 6 | `[1, 3, 2]` | 9 |

Each row is one add and one drop — the window's sum is never recomputed from scratch.

---

## Day 2 — Fixed-Size Window

**Pattern:** Window of exactly `k`; slide by one, update by one add and one drop  
**Key insight:** The only state you need is whatever lets you answer the question for the current window — a sum for averages, a character count for anagrams.

**Template adaptation:**

```go
func maxSumFixed(nums []int, k int) int {
    sum, best := 0, math.MinInt
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        if right >= k {
            sum -= nums[right-k]
        }
        if right >= k-1 {
            best = max(best, sum)
        }
    }
    return best
}
```

**Why it works:** the window `[right-k+1, right]` and the next one `[right-k+2, right+1]` share `k-1` elements. Their sums differ by exactly one added and one dropped element, so `sum` stays correct with O(1) work per step. `best` starts at `math.MinInt` (not 0) so all-negative inputs still work.

**Run it:** `maxSumFixed` in [`slidingwindow.go`](slidingwindow.go) (dry run on Day 1).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Maximum Average Subarray I | Easy | https://leetcode.com/problems/maximum-average-subarray-i/ |
| 2 | Find All Anagrams in a String | Medium | https://leetcode.com/problems/find-all-anagrams-in-a-string/ |

**What to notice — Problem 1:** The template verbatim, divided by `k` at the end. Keep the sum as an `int` and divide once — don't average inside the loop.

**What to notice — Problem 2:** The window size is fixed at `len(p)`. State is a `[26]int` count of the window's letters; the window is an anagram when it equals `p`'s count array. Go arrays are comparable with `==`, so `windowCount == pCount` is a single O(26) check.

**Complexity:** Problem 1 O(n) time, O(1) space. Problem 2 O(n · 26) = O(n) time, O(1) space (two fixed-size arrays).

**Common mistakes:**

- Recording the answer before the window is full (`right < k-1`).
- Initializing `best` to 0 — wrong when every window sum is negative.
- Anagrams: comparing maps instead of arrays. Maps need a loop to compare and leftover zero-count keys make equal windows look different.

---

## Day 3 — Variable Window: Longest Valid

**Pattern:** Grow `right` every step; shrink `left` only while the window is invalid  
**Key insight:** After the inner loop, `[left, right]` is the **longest valid window ending at `right`**. The answer is the best of those over all `right`.

**Template adaptation:**

```go
// longest subarray with sum <= limit (nums non-negative)
func longestSumAtMost(nums []int, limit int) int {
    left, sum, best := 0, 0, 0
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        for sum > limit { // invalid — shrink until valid again
            sum -= nums[left]
            left++
        }
        best = max(best, right-left+1)
    }
    return best
}
```

**Why it works:** suppose the best window ending at `right` starts at `left`. For `right+1`, the best start can't be *earlier* than `left` — if `[left-1, right+1]` were valid, then `[left-1, right]` (smaller, inside it) would be valid too, contradicting that `left` was the earliest valid start for `right`. So `left` never needs to move back, and one forward-only pointer finds every "longest window ending here".

**Dry run** — `longestSumAtMost([3, 1, 2, 1, 1, 4, 1], 5)`:

| `right` | add | Shrink (removed) | `left` | `sum` | Window | best |
|---|---|---|---|---|---|---|
| 0 | 3 | — | 0 | 3 | `[3]` | 1 |
| 1 | 1 | — | 0 | 4 | `[3, 1]` | 2 |
| 2 | 2 | 3 | 1 | 3 | `[1, 2]` | 2 |
| 3 | 1 | — | 1 | 4 | `[1, 2, 1]` | 3 |
| 4 | 1 | — | 1 | 5 | `[1, 2, 1, 1]` | 4 |
| 5 | 4 | 1, 2, 1 | 4 | 5 | `[1, 4]` | 4 |
| 6 | 1 | 1 | 5 | 5 | `[4, 1]` | 4 |

At `right = 5` the inner loop removes three elements in one go — that's fine, because those elements are gone for good.

**Run it:** `longestSumAtMost` in [`slidingwindow.go`](slidingwindow.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Longest Substring Without Repeating Characters | Medium | https://leetcode.com/problems/longest-substring-without-repeating-characters/ |
| 2 | Longest Repeating Character Replacement | Medium | https://leetcode.com/problems/longest-repeating-character-replacement/ |

**What to notice — Problem 1:** State = count of each character in the window. Invalid = the character just added has count > 1. Shrink until it's back to 1.

**What to notice — Problem 2:** Invalid = `windowLength - maxCount > k` (more characters to replace than allowed), where `maxCount` is the highest single-letter count in the window. The well-known trick: `maxCount` never needs to *decrease* when you shrink. A stale, too-high `maxCount` only keeps the window from shrinking more than it needs to, and the answer can only improve when a genuinely higher `maxCount` appears.

**Complexity:** both O(n) time. Space O(alphabet) for the count array.

**Common mistakes:**

- Using `if` instead of `for` for the shrink step. One removal may not be enough to make the window valid.
- Recording `best` before shrinking — you record an invalid window.
- Problem 1: storing "last seen index" and jumping `left` there without `max(left, lastSeen+1)`. `left` would jump backward when the repeat is to the left of the current window.

---

## Day 4 — Variable Window: Shortest Valid

**Pattern:** Grow `right` until valid; then shrink `left` as far as it stays valid, recording each time  
**Key insight:** For "shortest", the interesting moments are when the window **is** valid — so the record happens *inside* the shrink loop, before each removal.

**Template adaptation:**

```go
// shortest subarray with sum >= target (nums positive)
func minSubArrayLen(target int, nums []int) int {
    left, sum := 0, 0
    best := len(nums) + 1 // "infinity"
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        for sum >= target { // valid — record, then try a shorter one
            best = min(best, right-left+1)
            sum -= nums[left]
            left++
        }
    }
    if best == len(nums)+1 {
        return 0
    }
    return best
}
```

**Why it works:** mirror of Day 3. Once `[left, right]` is valid, extending `right` further can only make it longer, so the only way to improve is to shrink from the left. The inner loop keeps shrinking until the window breaks; the last valid window it saw is the shortest valid window ending at `right`.

**Dry run** — `minSubArrayLen(7, [2, 3, 1, 2, 4, 3])`:

| `right` | add | `sum` | Shrink steps (record, then remove) | `left` after | best |
|---|---|---|---|---|---|
| 0 | 2 | 2 | — | 0 | — |
| 1 | 3 | 5 | — | 0 | — |
| 2 | 1 | 6 | — | 0 | — |
| 3 | 2 | 8 | `[2, 3, 1, 2]` len 4 | 1 | 4 |
| 4 | 4 | 10 | `[3, 1, 2, 4]` len 4; `[1, 2, 4]` len 3 | 3 | 3 |
| 5 | 3 | 9 | `[2, 4, 3]` len 3; `[4, 3]` len 2 | 5 | 2 |

**Run it:** `minSubArrayLen` in [`slidingwindow.go`](slidingwindow.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Minimum Size Subarray Sum | Medium | https://leetcode.com/problems/minimum-size-subarray-sum/ |
| 2 | Minimum Window Substring | Hard | https://leetcode.com/problems/minimum-window-substring/ |

**What to notice — Problem 1:** This IS the template. Write it from memory until it's automatic — Problem 2 is the same skeleton with a richer "valid" check.

**What to notice — Problem 2:** State = count of each needed character still missing, plus a single integer `missing` (how many required characters aren't covered yet). Valid = `missing == 0`. Adding a character decrements `missing` only if it was still needed; removing one increments `missing` only if it drops below what `t` requires. That integer turns an O(alphabet) validity check into O(1).

**Complexity:** Problem 1 O(n) time, O(1) space. Problem 2 O(|s| + |t|) time, O(alphabet) space.

**Common mistakes:**

- Recording `best` after the shrink loop (Day 3 style) — by then the window is invalid.
- Problem 1: forgetting to return 0 when no window is ever valid.
- Problem 2: checking validity by comparing full count maps on every step — O(alphabet) per step and easy to get wrong. Use the `missing` counter.

---

## Day 5 — Counting Subarrays

**Pattern:** Count every valid window, not just the best one  
**Key insight:** When validity is monotonic, every subarray `[i, right]` with `left ≤ i ≤ right` is also valid — so the number of valid subarrays **ending at** `right` is `right - left + 1`. Sum that over all `right`.

**Template adaptation:**

```go
// subarrays with sum <= goal (nums non-negative)
func countSumAtMost(nums []int, goal int) int {
    if goal < 0 {
        return 0
    }
    left, sum, count := 0, 0, 0
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        for sum > goal {
            sum -= nums[left]
            left++
        }
        count += right - left + 1 // every start in [left, right] works
    }
    return count
}

// "exactly" isn't monotonic, but it's the difference of two "at most"s
func countSumExactly(nums []int, goal int) int {
    return countSumAtMost(nums, goal) - countSumAtMost(nums, goal-1)
}
```

**Why "exactly = atMost(k) − atMost(k−1)" works:** "sum exactly 2" is not monotonic — shrinking a sum-2 window can give sum 1. But "sum ≤ 2" is, and every subarray with sum ≤ 2 has sum exactly 2 **or** sum ≤ 1. Subtract the second group and you're left with exactly 2.

**Dry run** — `countSumExactly([1, 0, 1, 0, 1], 2)`. First `atMost(2)`:

| `right` | `left` | Window | New subarrays ending here | Count |
|---|---|---|---|---|
| 0 | 0 | `[1]` | 1 | 1 |
| 1 | 0 | `[1, 0]` | 2 | 3 |
| 2 | 0 | `[1, 0, 1]` | 3 | 6 |
| 3 | 0 | `[1, 0, 1, 0]` | 4 | 10 |
| 4 | 1 | `[0, 1, 0, 1]` | 4 | 14 |

Then `atMost(1)`:

| `right` | `left` | Window | New subarrays ending here | Count |
|---|---|---|---|---|
| 0 | 0 | `[1]` | 1 | 1 |
| 1 | 0 | `[1, 0]` | 2 | 3 |
| 2 | 1 | `[0, 1]` | 2 | 5 |
| 3 | 1 | `[0, 1, 0]` | 3 | 8 |
| 4 | 3 | `[0, 1]` | 2 | 10 |

Exactly 2 = 14 − 10 = **4**: `[1,0,1]`, `[1,0,1,0]`, `[0,1,0,1]`, `[1,0,1]` (the second one).

**Run it:** `countSumAtMost` and `countSumExactly` in [`slidingwindow.go`](slidingwindow.go) — tested against a brute force.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Subarray Product Less Than K | Medium | https://leetcode.com/problems/subarray-product-less-than-k/ |
| 2 | Subarrays with K Different Integers | Hard | https://leetcode.com/problems/subarrays-with-k-different-integers/ |

**What to notice — Problem 1:** Running product instead of running sum; invalid = `product >= k`. Every element is ≥ 1, so validity is monotonic. Guard `k <= 1` up front (no product of positive integers is `< 1`), or the shrink loop runs `left` past `right`.

**What to notice — Problem 2:** "exactly k distinct" = `atMost(k) − atMost(k−1)`, where `atMost` uses a count map and a `distinct` counter. Write `atMost` once and call it twice.

**Complexity:** both O(n) time. Problem 2 needs O(n) space for the count map.

**Common mistakes:**

- Trying to count "exactly k" directly with one window — it undercounts because shrinking skips valid windows.
- Adding 1 per valid window instead of `right - left + 1`.
- Problem 1: forgetting the `k <= 1` guard.

---

## Day 6 — Window Min/Max with a Monotonic Deque

**Pattern:** Maintain the min (or max) of a sliding window in O(1) amortized per step  
**Key insight:** It's the Monotonic Stack week's idea plus expiry: elements leave the **back** when a better one arrives, and leave the **front** when they slide out of the window.

**Template adaptation:**

```go
// minimum of every window of size k
func minSlidingWindow(nums []int, k int) []int {
    var result []int
    deque := []int{} // indices; their values increase from front to back

    for right := 0; right < len(nums); right++ {
        if len(deque) > 0 && deque[0] <= right-k {
            deque = deque[1:] // front index slid out of the window
        }
        for len(deque) > 0 && nums[deque[len(deque)-1]] >= nums[right] {
            deque = deque[:len(deque)-1] // can never be a window min again
        }
        deque = append(deque, right)

        if right >= k-1 {
            result = append(result, nums[deque[0]]) // front is the min
        }
    }
    return result
}
```

**Why popping from the back is safe:** if an older element is ≥ the new one, the new one is both **smaller** and **leaves the window later**. Every future window containing the old element also contains the new one, so the old element can never be the minimum again — drop it. What remains is increasing from front to back, so the front is always the current minimum.

**Dry run** — `minSlidingWindow([4, 2, 12, 3, 8, 1, 7], 3)`:

| `right` | value | Expired from front | Popped from back | Deque (values) | Window min |
|---|---|---|---|---|---|
| 0 | 4 | — | — | `[4]` | (filling) |
| 1 | 2 | — | 4 | `[2]` | (filling) |
| 2 | 12 | — | — | `[2, 12]` | 2 |
| 3 | 3 | — | 12 | `[2, 3]` | 2 |
| 4 | 8 | 2 | — | `[3, 8]` | 3 |
| 5 | 1 | — | 8, 3 | `[1]` | 1 |
| 6 | 7 | — | — | `[1, 7]` | 1 |

At `right = 4`, the `2` (index 1) expired because the window is now `[2..4]`.

**Run it:** `minSlidingWindow` in [`slidingwindow.go`](slidingwindow.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Sliding Window Maximum | Hard | https://leetcode.com/problems/sliding-window-maximum/ |
| 2 | Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit | Medium | https://leetcode.com/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/ |

**What to notice — Problem 1:** The template with the comparison flipped (`<=` instead of `>=`) so the deque is decreasing and the front is the max.

**What to notice — Problem 2:** A **variable** window (Day 3) whose validity check is `max - min <= limit`. Keep two deques at once — one for max, one for min. When shrinking, pop a deque's front only if that front index is the one leaving (`== left`).

**Complexity:** both O(n) time — each index enters and leaves each deque at most once. O(k) or O(n) space for the deques.

**Common mistakes:**

- Storing values instead of indices — then you can't tell when the front has expired.
- Using `>` instead of `>=` when popping the back. Both are correct for the min, but `>=` keeps the deque smaller; with `>`, equal values pile up.
- Problem 2: popping the front unconditionally when shrinking.
- A heap also works (O(n log k)) but needs lazy deletion of expired items — the deque is simpler and faster.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Every sliding window problem comes down to three decisions: fixed or variable, what state to keep, and whether you're finding the longest, the shortest, or a count.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing code, answer out loud: *Is validity monotonic here? If not, why not, and what would I use instead?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the fixed and variable templates from memory
- [ ] Why is the variable template O(n) despite the nested loop?
- [ ] Longest vs shortest: where does the `best = ...` line go in each, and why?
- [ ] Why does "sum == k" with negative numbers break the sliding window?
- [ ] How do you count "exactly k" when only "at most k" is monotonic?
- [ ] In the min-deque, why can an element be popped from the back forever?

---

## Pattern Cheatsheet

| Pattern | Window | Shrink condition | Record where | Time |
|---------|--------|------------------|--------------|------|
| Fixed size | exactly `k` | always drop `right-k` | once full (`right >= k-1`) | O(n) |
| Longest valid | variable | `for` invalid | after the shrink loop | O(n) |
| Shortest valid | variable | `for` valid | inside the shrink loop, before removing | O(n) |
| Count valid | variable | `for` invalid | `count += right-left+1` | O(n) |
| Exactly k | two variable windows | — | `atMost(k) - atMost(k-1)` | O(n) |
| Window min/max | fixed or variable + deque | expire front, pop worse from back | front of deque | O(n) |

---

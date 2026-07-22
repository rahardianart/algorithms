# Monotonic Stack Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, DP, Heap, Binary Search, and Linked List weeks)  
**Goal:** Write the monotonic stack template from memory, decide increasing vs decreasing on sight, recognize "next greater/smaller" problems even when disguised.

---

## The Core Template

Maintain a stack that stays monotonic (strictly increasing or decreasing, bottom to top) by popping every element the incoming value invalidates before pushing.

```go
func nextGreaterElement(nums []int) []int {
    result := make([]int, len(nums))
    for i := range result {
        result[i] = -1 // default: no greater element exists
    }
    stack := []int{} // holds INDICES, values decreasing from bottom to top

    for i, n := range nums {
        for len(stack) > 0 && nums[stack[len(stack)-1]] < n {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[top] = n // n is the next greater element for index `top`
        }
        stack = append(stack, i)
    }
    return result
}
```

**Three rules — never break them:**

1. **Store indices, not values.** You almost always need the index to write into a result array, compute a distance (`i - poppedIndex`), or look up other properties of that element. Pushing raw values throws that information away.
2. **Decreasing stack finds "next greater"; increasing stack finds "next smaller."** Pop while `top < current` → decreasing stack → next *greater* element. Pop while `top > current` → increasing stack → next *smaller* element. Read the problem, decide which one you need, before writing any code.
3. **Each element is pushed once and popped at most once — that's what makes it O(n).** The inner `while` loop looks like it could be O(n²), but across the *entire* run every element enters and leaves the stack a single time. Don't second-guess the complexity mid-interview.

---

## Day 1 — Template + Concept

**Pattern:** Core template  
**Key insight:** A monotonic stack answers "what's the nearest element to my left/right that's greater/smaller than me" for every element, in one linear pass — the brute-force nested-loop version of this question is O(n²).

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each line, especially the inner `while` loop condition.
2. Close this file. Write the template from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why does storing indices instead of values make the amortized O(n) argument easier to see?*

   > With indices, it's visible that the stack holds at most `n` items ever (one per index, 0 to n-1), and each index is pushed exactly once. The inner loop only removes items already on the stack, so total pops across the whole run can't exceed total pushes — n. Two O(n) operations (n pushes, ≤n pops) is O(n) overall, not O(n²), even though the loop is nested.

---

## Day 2 — Next Greater / Next Smaller

**Pattern:** One pass, decreasing stack, next-greater lookup  
**Key insight:** "Next greater element" and "how many days until warmer" are the same question — the second is just "next greater" phrased as a distance instead of a value.

**Template adaptation:**

```go
func dailyTemperatures(temperatures []int) []int {
    result := make([]int, len(temperatures))
    stack := []int{} // indices, decreasing temps from bottom to top

    for i, t := range temperatures {
        for len(stack) > 0 && temperatures[stack[len(stack)-1]] < t {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[top] = i - top // distance, not the value itself
        }
        stack = append(stack, i)
    }
    return result
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Next Greater Element I | Easy | https://leetcode.com/problems/next-greater-element-i/ |
| 2 | Daily Temperatures | Medium | https://leetcode.com/problems/daily-temperatures/ |

**What to notice — Problem 1:** Two arrays, but only `nums2` needs the monotonic stack — run the core template once on `nums2` into a `map[int]int` of value → next-greater, then look up each element of `nums1` in that map. Don't re-run the stack per query.

**What to notice — Problem 2:** Identical stack logic to the core template, but `result[top]` stores `i - top` (the wait, in days) instead of `t` (the temperature). This substitution — index math instead of the raw popped value — is the most common variation across "next greater" problems.

---

## Day 3 — Circular Arrays & Running Span

**Pattern:** Wrap the scan around, or maintain a running answer across streaming input  
**Key insight:** For a circular array, iterate `2n` times using `i % n` — the stack still only ever holds real indices `0` to `n-1`, so the O(n) bound is unaffected. For streaming span problems, the stack persists across calls instead of being built once.

**Template adaptation:**

```go
func nextGreaterElements(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result {
        result[i] = -1
    }
    stack := []int{}

    for i := 0; i < 2*n; i++ {
        idx := i % n
        for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[idx] {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[top] = nums[idx]
        }
        if i < n { // only push real indices once — second pass just resolves wraparound
            stack = append(stack, idx)
        }
    }
    return result
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Next Greater Element II | Medium | https://leetcode.com/problems/next-greater-element-ii/ |
| 2 | Online Stock Span | Medium | https://leetcode.com/problems/online-stock-span/ |

**What to notice — Problem 1:** The `i < n` guard is the key line — you scan `2n` times so every index gets a second chance to find a greater element that wraps around from the front, but you must not push duplicate indices on the second pass, or the stack size bound breaks.

**What to notice — Problem 2:** A design problem, not a one-shot array scan — the stack lives in the struct and persists across `next(price)` calls. Store `(price, span)` pairs; when popping a smaller-or-equal price, absorb its span into the current call's span before pushing. This is the decreasing-stack template with the popped value's span folded into the new entry instead of discarded.

---

## Day 4 — Histogram Area

**Pattern:** For each bar, find how far left/right you can extend before hitting something shorter  
**Key insight:** The maximal rectangle at bar `i` is bounded by the nearest shorter bar on each side — exactly the "next smaller element" question from Day 2, computed for both directions at once using one increasing stack.

**Template adaptation:**

```go
func largestRectangleArea(heights []int) int {
    stack := []int{} // indices, increasing heights from bottom to top
    maxArea := 0

    for i := 0; i <= len(heights); i++ {
        h := 0
        if i < len(heights) {
            h = heights[i]
        }
        for len(stack) > 0 && heights[stack[len(stack)-1]] >= h {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            height := heights[top]

            width := i // if stack is now empty, this bar extends back to index 0
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            maxArea = max(maxArea, height*width)
        }
        stack = append(stack, i)
    }
    return maxArea
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Largest Rectangle in Histogram | Hard | https://leetcode.com/problems/largest-rectangle-in-histogram/ |
| 2 | Trapping Rain Water | Hard | https://leetcode.com/problems/trapping-rain-water/ |

**What to notice — Problem 1:** The loop runs to `len(heights)` inclusive, using a sentinel height of `0` on the last iteration — this forces every remaining bar on the stack to be popped and evaluated by the end. When a bar is popped, its rectangle's width spans from the new stack top (exclusive) to the current index (exclusive) — both boundaries are "nearest shorter bar."

**What to notice — Problem 2:** Same family as Problem 1 but tracking trapped *water* instead of rectangle area — use a decreasing stack of indices; when you pop a bar because a taller one arrived, the trapped water above it is bounded by `min(current height, new stack top's height) - popped height`, times the width between them. (The two-pointer solution from your Two Pointer week solves this in O(1) space — worth comparing both approaches.)

---

## Day 5 — Greedy Removal for Optimal Subsequence

**Pattern:** Build a result by popping "worse" earlier characters when a better one arrives  
**Key insight:** To build the smallest (or otherwise optimal) subsequence, greedily remove a previous character whenever the current one would make a better result — but only if you're still allowed to remove (a budget, or "each letter used once").

**Template adaptation:**

```go
func removeKdigits(num string, k int) string {
    stack := []byte{} // increasing digits from bottom to top

    for i := 0; i < len(num); i++ {
        for k > 0 && len(stack) > 0 && stack[len(stack)-1] > num[i] {
            stack = stack[:len(stack)-1]
            k--
        }
        stack = append(stack, num[i])
    }
    stack = stack[:len(stack)-k] // if we didn't use up k removals in the loop, trim from the end

    // strip leading zeros
    start := 0
    for start < len(stack)-1 && stack[start] == '0' {
        start++
    }
    result := string(stack[start:])
    if result == "" {
        return "0"
    }
    return result
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Remove K Digits | Medium | https://leetcode.com/problems/remove-k-digits/ |
| 2 | Remove Duplicate Letters | Medium | https://leetcode.com/problems/remove-duplicate-letters/ |

**What to notice — Problem 1:** `k` acts as a removal budget that shrinks in the inner loop — the greedy insight is that removing a larger digit before a smaller one always produces a smaller number, so pop whenever you still can. The final `stack[:len(stack)-k]` handles a strictly increasing input where the inner loop never fires.

**What to notice — Problem 2:** Two extra conditions layered onto the same increasing-stack skeleton: only pop the top if it reappears later in the string (check a precomputed "last occurrence" map), and never push a character already in the stack (track with a `seen` set). The core pop/push shape is unchanged — the guard conditions are what make it "duplicate letters" instead of "remove k digits."

---

## Day 6 — Contribution Technique

**Pattern:** For every element, count how many subarrays/subproblems it "contributes to" using its next/previous smaller boundaries  
**Key insight:** Instead of asking "what's the answer for this window," ask "for how many windows is this specific element the deciding one" — the monotonic stack gives you the left and right boundary where that's true in O(1) amortized per element.

**Template adaptation:**

```go
func sumSubarrayMins(arr []int) int {
    const mod = 1_000_000_007
    n := len(arr)
    prevLess := make([]int, n) // distance to previous strictly-smaller element (or start)
    nextLessEq := make([]int, n) // distance to next smaller-or-equal element (or end)

    stack := []int{}
    for i := 0; i < n; i++ {
        for len(stack) > 0 && arr[stack[len(stack)-1]] > arr[i] {
            stack = stack[:len(stack)-1]
        }
        if len(stack) == 0 {
            prevLess[i] = i + 1
        } else {
            prevLess[i] = i - stack[len(stack)-1]
        }
        stack = append(stack, i)
    }

    stack = stack[:0]
    for i := n - 1; i >= 0; i-- {
        for len(stack) > 0 && arr[stack[len(stack)-1]] >= arr[i] {
            stack = stack[:len(stack)-1]
        }
        if len(stack) == 0 {
            nextLessEq[i] = n - i
        } else {
            nextLessEq[i] = stack[len(stack)-1] - i
        }
        stack = append(stack, i)
    }

    total := 0
    for i := 0; i < n; i++ {
        total = (total + arr[i]*prevLess[i]*nextLessEq[i]) % mod
    }
    return total
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Sum of Subarray Minimums | Medium | https://leetcode.com/problems/sum-of-subarray-minimums/ |
| 2 | 132 Pattern | Medium | https://leetcode.com/problems/132-pattern/ |

**What to notice — Problem 1:** `arr[i]` is the minimum for exactly `prevLess[i] * nextLessEq[i]` subarrays (choices of left boundary times choices of right boundary) — sum `arr[i] * prevLess[i] * nextLessEq[i]` over all `i` instead of enumerating subarrays directly, which would be O(n²). The `>` vs `>=` asymmetry between the two passes avoids double-counting subarrays where the minimum value repeats.

**What to notice — Problem 2:** Different goal (existence, not a sum) but the same "process right to left, maintain a decreasing stack" shape. Track a running `third` value (the best candidate popped so far while scanning right to left) — if a later (leftward) element is smaller than `third`, a valid 132 pattern exists. This is a variation worth recognizing rather than deriving from scratch under time pressure.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Every monotonic stack problem reduces to two decisions: increasing or decreasing, and what you store when you pop (a value, a distance, or a running total). The push/pop mechanics never change.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Answer out loud: *For a new problem, is it asking for "next greater" or "next smaller," and what do you store on each pop?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the core monotonic stack template from memory
- [ ] When do you use an increasing stack vs a decreasing stack?
- [ ] Why is the amortized complexity O(n) despite the nested loop?
- [ ] In Largest Rectangle in Histogram, what do the two neighbors on the stack represent when a bar is popped?
- [ ] What's the difference between the contribution technique and a direct "next greater" lookup?

---

## Pattern Cheatsheet

| Pattern | Stack direction | What you store on pop | Complexity win |
|---------|------------------|--------------------------|------------------|
| Next greater/smaller | decreasing (greater) / increasing (smaller) | the popped index's answer (value or distance) | O(n²) → O(n) |
| Circular / streaming | same as above, scan `2n` or persist across calls | same, with a guard against double-pushing | O(n²) → O(n) |
| Histogram area | increasing | width = distance between new top and current index | O(n²) → O(n) |
| Greedy removal | increasing (usually) | nothing — the pop itself IS the removal | O(n·k) → O(n) |
| Contribution technique | increasing/decreasing (paired passes) | boundary distances, multiplied into a running total | O(n²) → O(n) |

---

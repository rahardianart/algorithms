# Dynamic Programming Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, and Two Pointer weeks)  
**Goal:** Recognize DP problems, define state and recurrence before coding, write both top-down and bottom-up from memory.

---

## The Two Approaches

### Top-Down (Memoization)

Recursion + cache. Start from the full problem, break it into subproblems, cache results so each is computed once.

```go
func solve(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    if v, ok := memo[n]; ok {
        return v
    }
    result := solve(n-1, memo) + solve(n-2, memo)
    memo[n] = result
    return result
}
```

### Bottom-Up (Tabulation)

Iterative. Build up from the base case to the full answer, filling a table (usually a slice).

```go
func solve(n int) int {
    if n <= 1 {
        return n
    }
    dp := make([]int, n+1)
    dp[0], dp[1] = 0, 1
    for i := 2; i <= n; i++ {
        dp[i] = dp[i-1] + dp[i-2]
    }
    return dp[n]
}
```

**Three rules — never break them:**

1. **Define the state first** — what does `dp[i]` (or `dp[i][j]`) actually represent? Say it in one sentence before writing any code. Every DP bug traces back to a fuzzy state definition.
2. **Find the recurrence** — how does `dp[i]` relate to smaller subproblems? This falls directly out of the state definition; it doesn't need to be guessed.
3. **Base cases are non-negotiable** — get them wrong and every value built on top is wrong too. Write them explicitly before the loop or recursion.

**Top-down vs bottom-up:** Same recurrence, different direction. Top-down is easier to derive (write the brute-force recursion, add a cache). Bottom-up is usually faster in practice (no recursion overhead) and easier to optimize for space (rolling array).

---

## Day 1 — Template + Concept

**Pattern:** Both approaches  
**Key insight:** DP is recursion with overlapping subproblems, cached. If subproblems don't overlap, plain recursion is enough — caching buys nothing.

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both templates (Fibonacci) from memory in Go — top-down with a map, bottom-up with a slice.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *What makes a problem a DP problem instead of plain recursion?*

   > Two things must both be true: **optimal substructure** (the optimal solution can be built from optimal solutions to subproblems) and **overlapping subproblems** (the same subproblem gets solved multiple times by the naive recursion). If subproblems don't repeat, you're doing divide-and-conquer or plain recursion — a cache wouldn't help because nothing is ever looked up twice.

---

## Day 2 — 1D DP: Linear Sequences

**Pattern:** `dp[i]` depends on a fixed number of previous states  
**Key insight:** State is just an index into the input. The recurrence looks back 1 or 2 steps — no inner loop needed.

**Template adaptation:**

```go
func climbStairs(n int) int {
    if n <= 2 {
        return n
    }
    dp := make([]int, n+1)
    dp[1], dp[2] = 1, 2
    for i := 3; i <= n; i++ {
        dp[i] = dp[i-1] + dp[i-2]
    }
    return dp[n]
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Climbing Stairs | Easy | https://leetcode.com/problems/climbing-stairs/ |
| 2 | House Robber | Medium | https://leetcode.com/problems/house-robber/ |

**What to notice — Problem 1:** `dp[i]` = number of ways to reach step `i`. Recurrence: `dp[i] = dp[i-1] + dp[i-2]` — you either took a 1-step from `i-1` or a 2-step from `i-2`. Identical shape to Fibonacci.

**What to notice — Problem 2:** `dp[i]` = max money robbable from the first `i` houses. Recurrence: `dp[i] = max(dp[i-1], dp[i-2]+nums[i])` — either skip house `i` (keep `dp[i-1]`) or rob it (add to `dp[i-2]`, since adjacent houses can't both be robbed). This "skip vs take" choice is the most common DP recurrence shape — you'll see it again on Day 6.

---

## Day 3 — 1D DP: Unbounded Choices & Subsequences

**Pattern:** `dp[i]` depends on ALL previous states, not just the last 1–2  
**Key insight:** When the recurrence needs to scan every earlier index (not a fixed lookback), the inner loop becomes part of the DP itself: O(n²) instead of O(n).

**Template adaptation:**

```go
func lengthOfLIS(nums []int) int {
    dp := make([]int, len(nums))
    for i := range dp {
        dp[i] = 1 // every element is a subsequence of length 1 by itself
    }

    best := 1
    for i := 1; i < len(nums); i++ {
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] {
                dp[i] = max(dp[i], dp[j]+1)
            }
        }
        best = max(best, dp[i])
    }
    return best
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Longest Increasing Subsequence | Medium | https://leetcode.com/problems/longest-increasing-subsequence/ |
| 2 | Coin Change | Medium | https://leetcode.com/problems/coin-change/ |

**What to notice — Problem 1:** `dp[i]` = length of the longest increasing subsequence *ending at* index `i` (not "in the first `i` elements" — that distinction is why the inner loop is necessary). The answer is `max(dp)`, not `dp[n-1]`.

**What to notice — Problem 2:** `dp[i]` = fewest coins to make amount `i`. Choices are unbounded (each coin reusable), so the inner loop is over `coins`, not previous amounts: `dp[i] = min(dp[i], dp[i-coin]+1)` for each `coin <= i`. Initialize `dp` with a sentinel (`amount+1`) to represent "unreachable." Return `-1` if `dp[amount]` is still the sentinel.

---

## Day 4 — 2D DP: Grids

**Pattern:** `dp[i][j]` depends on `dp[i-1][j]` and `dp[i][j-1]`  
**Key insight:** Grid DP is 1D DP with a second index. The recurrence is almost always "came from above or came from the left."

**Template adaptation:**

```go
func uniquePaths(m, n int) int {
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][0] = 1 // only one way to reach any cell in column 0: straight down
    }
    for j := 0; j < n; j++ {
        dp[0][j] = 1 // only one way to reach any cell in row 0: straight right
    }

    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            dp[i][j] = dp[i-1][j] + dp[i][j-1]
        }
    }
    return dp[m-1][n-1]
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Unique Paths | Medium | https://leetcode.com/problems/unique-paths/ |
| 2 | Minimum Path Sum | Medium | https://leetcode.com/problems/minimum-path-sum/ |

**What to notice — Problem 1:** `dp[i][j]` = number of ways to reach cell `(i, j)`. The first row and first column are the base cases (only one straight-line path each). Everywhere else: `dp[i][j] = dp[i-1][j] + dp[i][j-1]`.

**What to notice — Problem 2:** Same shape as Problem 1 but `min` instead of `+`, and you accumulate cost instead of counting paths: `dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])`. The first row/column base cases become running sums, not `1`s.

---

## Day 5 — 2D DP: Two Strings

**Pattern:** `dp[i][j]` compares prefixes of two strings  
**Key insight:** `dp[i][j]` = the answer for `s1[:i]` and `s2[:j]`. The recurrence branches on whether `s1[i-1] == s2[j-1]`.

**Template adaptation:**

```go
func longestCommonSubsequence(text1, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1) // dp[0][j] and dp[i][0] default to 0 — empty prefix has no common subsequence
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Longest Common Subsequence | Medium | https://leetcode.com/problems/longest-common-subsequence/ |
| 2 | Edit Distance | Medium | https://leetcode.com/problems/edit-distance/ |

**What to notice — Problem 1:** `dp` is sized `(m+1) x (n+1)`, not `m x n` — row/column 0 represent the "empty prefix," and Go zero-initializes them to 0, which is exactly the right base case. Characters match → extend the diagonal. No match → carry forward the best of dropping one character from either string.

**What to notice — Problem 2:** Same grid shape as Problem 1. No match → `1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])` (delete, insert, replace — try all three, take the cheapest). Base cases `dp[i][0] = i` and `dp[0][j] = j` must be set explicitly this time (converting to/from an empty string costs one op per character) — they aren't both 0 like Problem 1.

---

## Day 6 — Knapsack Pattern

**Pattern:** Choose a subset of items under a capacity constraint  
**Key insight:** `dp[i][capacity]` = the best answer using the first `i` items at a given capacity. The recurrence is always "skip item `i`, or take it and reduce capacity" — the same skip-vs-take shape from Day 2, now with a second dimension for the constraint.

**Template adaptation:**

```go
func canPartition(nums []int) bool {
    sum := 0
    for _, n := range nums {
        sum += n
    }
    if sum%2 != 0 {
        return false
    }
    target := sum / 2

    dp := make([]bool, target+1)
    dp[0] = true // sum 0 is always reachable — take nothing

    for _, num := range nums {
        for cap := target; cap >= num; cap-- { // iterate backward — 0/1 knapsack, each item used once
            dp[cap] = dp[cap] || dp[cap-num]
        }
    }
    return dp[target]
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Partition Equal Subset Sum | Medium | https://leetcode.com/problems/partition-equal-subset-sum/ |
| 2 | Coin Change II | Medium | https://leetcode.com/problems/coin-change-ii/ |

**What to notice — Problem 1:** This is 0/1 knapsack collapsed to 1D: `dp[cap]` = "is sum `cap` reachable using items seen so far." The backward inner loop (`cap := target; cap >= num; cap--`) is what limits each item to at most one use — a forward loop would let an item be reused, since it would read a `dp[cap-num]` already updated in the same pass.

**What to notice — Problem 2:** Compare to Coin Change (Day 3): that problem minimizes coin *count*, this one counts the number of *combinations*. The loop order flips meaning: coins in the outer loop, amounts in the inner loop (forward this time) — this counts each combination once regardless of coin order. Recurrence: `dp[amt] += dp[amt-coin]`.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** The hard part of DP interviews is never the code — it's stating the state and recurrence out loud before writing anything. If you can say "`dp[i][j]` represents X, and it comes from Y" in one sentence, the code is mechanical.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing any code, say the state definition and recurrence out loud first. Only then write the DP table or recursion.

**Self-assessment — can you answer these without hesitation?**

- [ ] What two properties make a problem solvable with DP?
- [ ] What's the difference between top-down and bottom-up, and when would you prefer one over the other?
- [ ] Why does Longest Increasing Subsequence need an O(n²) inner loop but House Robber doesn't?
- [ ] Why does the 0/1 knapsack inner loop go backward, and what breaks if it goes forward?
- [ ] Given a brand-new problem, what's the first sentence you write before any code?

---

## Pattern Cheatsheet

| Pattern | State | Recurrence shape | Dimensions |
|---------|-------|-------------------|------------|
| Linear (fixed lookback) | `dp[i]` | Looks back 1–2 steps | 1D |
| Subsequence / unbounded | `dp[i]` | Scans all `j < i`, or reuses items | 1D + inner loop |
| Grid | `dp[i][j]` | From above or from the left | 2D |
| Two strings | `dp[i][j]` | Branch on character match | 2D |
| Knapsack | `dp[i][capacity]` (or rolled 1D) | Skip vs take, capacity shrinks | 2D or rolled 1D |

---

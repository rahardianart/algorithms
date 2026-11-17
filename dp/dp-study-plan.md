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

**Brute force → DP:** the brute force for almost every DP problem is a recursion that tries every choice. It's exponential because it re-solves the *same* subproblem over and over (see the call counts on Day 1). DP keeps the same recursion but stores each subproblem's answer the first time it's computed, so the total work becomes **number of distinct states × work per state**. That formula is how you compute every DP's complexity.

**How to find the state — a recipe for a new problem:**

1. Write the brute-force recursion first. What parameters does it take? Those parameters **are** the state (`i`, `(i, j)`, `(i, capacity)`, …).
2. Say what the function returns for those parameters in one sentence. That's your `dp[...]` definition.
3. The recursive calls inside it are your recurrence; its base cases are your base cases.
4. Add a memo → top-down done. Want bottom-up? Fill the table in the order the recurrence needs (smaller indices first).

**Run it:** every template in this file is a real function in [`dp.go`](dp.go), tested in [`dp_test.go`](dp_test.go). From the repo root: `go test ./dp/ -v`.

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

**See the overlap yourself** — number of function calls to compute `fib(n)`:

| `n` | Naive calls | Memo calls |
|---|---|---|
| 5 | 15 | 9 |
| 10 | 177 | 19 |
| 20 | 21,891 | 39 |
| 30 | 2,692,537 | 59 |
| 40 | 331,160,281 | 79 |

Naive recursion makes about φⁿ ≈ 1.6ⁿ calls because `fib(n-2)` is computed inside `fib(n-1)` *and* again directly. With the memo, each `fib(k)` is computed once and every later call is a map lookup: about `2n` calls total. There are only `n` distinct states — that's the whole speedup.

**Run it:** `fibMemo`, `fibTab`, and `fibRolling` in [`dp.go`](dp.go). `TestFib` computes `fib(90)` with the memo instantly; the naive version would take longer than the age of the universe.

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

**Why it works:** to stand on step `i`, your *last* move was either a 1-step from `i-1` or a 2-step from `i-2`. Those two groups don't overlap (different last move) and together cover everything, so `dp[i] = dp[i-1] + dp[i-2]`. Splitting on "what was the last decision" is the most reliable way to find any recurrence.

**Dry run** — `climbStairs(6)`:

| `i` | 1 | 2 | 3 | 4 | 5 | 6 |
|---|---|---|---|---|---|---|
| `dp[i]` | 1 | 2 | 3 | 5 | 8 | **13** |

**Space optimization:** `dp[i]` only reads the two previous values, so the whole slice can be replaced by two variables:

```go
prev2, prev1 := 1, 2 // dp[1], dp[2]
for i := 3; i <= n; i++ {
    prev2, prev1 = prev1, prev1+prev2
}
return prev1
```

O(n) space → O(1). Whenever a recurrence only looks back a fixed number of steps, you can do this.

**Run it:** `climbStairs` and `fibRolling` (the two-variable version) in [`dp.go`](dp.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Climbing Stairs | Easy | https://leetcode.com/problems/climbing-stairs/ |
| 2 | House Robber | Medium | https://leetcode.com/problems/house-robber/ |

**What to notice — Problem 1:** `dp[i]` = number of ways to reach step `i`. Recurrence: `dp[i] = dp[i-1] + dp[i-2]` — you either took a 1-step from `i-1` or a 2-step from `i-2`. Identical shape to Fibonacci.

**What to notice — Problem 2:** `dp[i]` = max money robbable from the first `i` houses. Recurrence: `dp[i] = max(dp[i-1], dp[i-2]+nums[i])` — either skip house `i` (keep `dp[i-1]`) or rob it (add to `dp[i-2]`, since adjacent houses can't both be robbed). This "skip vs take" choice is the most common DP recurrence shape — you'll see it again on Day 6.

**Complexity:** both problems have n states and O(1) work per state → O(n) time. O(n) space with the table, O(1) with rolling variables.

**Common mistakes:**

- House Robber: indexing confusion between "first `i` houses" and "house `i`". Pick one definition, write it down, and size the slice to match (`n` vs `n+1`).
- House Robber: missing the base cases for 1 house (`nums[0]`) and 2 houses (`max(nums[0], nums[1])`).
- Climbing Stairs: `dp[0]` — is it 0 or 1? The template avoids the question by starting at `dp[1]`, `dp[2]`.

---

## Day 3 — 1D DP: Unbounded Choices & Subsequences

**Pattern:** `dp[i]` depends on ALL previous states, not just the last 1–2  
**Key insight:** When the recurrence needs to scan every earlier index (not a fixed lookback), the inner loop becomes part of the DP itself: O(n²) instead of O(n).

**Template adaptation:**

```go
func lengthOfLIS(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
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
```

`max` and `min` are Go builtins since Go 1.21 — no helper needed.

**Why "ending at `i`" matters:** "LIS in the first `i` elements" isn't enough information to extend — you'd need to know the *last value* of that subsequence to check whether `nums[i]` can follow it. "LIS *ending at* `i`" bakes the last value into the state (it's `nums[i]`), so extending is a simple `nums[j] < nums[i]` check. When a recurrence needs information the state doesn't carry, change the state definition.

**Dry run** — `lengthOfLIS([10, 9, 2, 5, 3, 7, 101, 18])`:

| `i` | `nums[i]` | `dp[i]` | Best `j` (extends …) |
|---|---|---|---|
| 0 | 10 | 1 | — |
| 1 | 9 | 1 | — |
| 2 | 2 | 1 | — |
| 3 | 5 | 2 | 2 (`2`) |
| 4 | 3 | 2 | 2 (`2`) |
| 5 | 7 | 3 | 3 (`5`) |
| 6 | 101 | 4 | 5 (`7`) |
| 7 | 18 | 4 | 5 (`7`) |

The answer is `max(dp) = 4` (for example `2, 5, 7, 101`), **not** `dp[7]` — which happens to also be 4 here, but wouldn't be for `[1, 2, 3, 0]`. Following the "Best `j`" column backwards from index 6 rebuilds the actual subsequence: 6 → 5 → 3 → 2.

**Run it:** `lengthOfLIS` in [`dp.go`](dp.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Longest Increasing Subsequence | Medium | https://leetcode.com/problems/longest-increasing-subsequence/ |
| 2 | Coin Change | Medium | https://leetcode.com/problems/coin-change/ |

**What to notice — Problem 1:** `dp[i]` = length of the longest increasing subsequence *ending at* index `i` (not "in the first `i` elements" — that distinction is why the inner loop is necessary). The answer is `max(dp)`, not `dp[n-1]`.

**What to notice — Problem 2:** `dp[i]` = fewest coins to make amount `i`. Choices are unbounded (each coin reusable), so the inner loop is over `coins`, not previous amounts: `dp[i] = min(dp[i], dp[i-coin]+1)` for each `coin <= i`. Initialize `dp` with a sentinel (`amount+1`) to represent "unreachable." Return `-1` if `dp[amount]` is still the sentinel.

**Complexity:**

| Problem | States | Work per state | Time | Space |
|---|---|---|---|---|
| LIS | n | O(n) scan of earlier `j` | O(n²) | O(n) |
| Coin Change | `amount + 1` | O(number of coins) | O(amount × coins) | O(amount) |

(LIS also has an O(n log n) solution using binary search on "smallest tail of each length" — a good stretch exercise after your Binary Search week.)

**Common mistakes:**

- LIS: returning `dp[n-1]` instead of `max(dp)`.
- Coin Change: initializing `dp` with `math.MaxInt` and then computing `dp[i-coin] + 1` — it overflows to a negative number and wins every `min`. Use `amount + 1` as "infinity".
- Coin Change: forgetting `dp[0] = 0`.

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

**Why it works:** you can only move right or down, so the last move into `(i, j)` came from above `(i-1, j)` or from the left `(i, j-1)` — the same "split on the last decision" as Day 2, now in two dimensions. The two groups don't overlap, so the counts add.

**Dry run** — `uniquePaths(3, 4)`. Each cell is the sum of the cell above and the cell to the left:

```
 1  1  1  1
 1  2  3  4
 1  3  6 10
```

The answer is the bottom-right cell, `10`. Fill order matters: row by row, left to right, guarantees both cells you read are already filled.

**Space optimization:** each row only reads the row above, so one slice is enough — `row[j] += row[j-1]` (old `row[j]` is "above", new `row[j-1]` is "left"). O(m·n) space → O(n).

**Run it:** `uniquePaths` in [`dp.go`](dp.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Unique Paths | Medium | https://leetcode.com/problems/unique-paths/ |
| 2 | Minimum Path Sum | Medium | https://leetcode.com/problems/minimum-path-sum/ |

**What to notice — Problem 1:** `dp[i][j]` = number of ways to reach cell `(i, j)`. The first row and first column are the base cases (only one straight-line path each). Everywhere else: `dp[i][j] = dp[i-1][j] + dp[i][j-1]`.

**What to notice — Problem 2:** Same shape as Problem 1 but `min` instead of `+`, and you accumulate cost instead of counting paths: `dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])`. The first row/column base cases become running sums, not `1`s.

**Complexity:** both problems O(m·n) time — m·n states, O(1) work each. O(m·n) space with the full grid, O(n) with one rolling row.

**Common mistakes:**

- Minimum Path Sum: initializing the first row/column to `1` (copied from Unique Paths). They're running sums of `grid` values.
- Minimum Path Sum: using `min` on the first row/column where there's only one direction to come from — `dp[0][j-1]` is the only option there.
- Filling column by column but reading `dp[i][j-1]` before it's computed. Pick a fill order that matches what the recurrence reads.

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

**Why it works:** look at the last character of each prefix. If they match, that character is in some longest common subsequence, so take it and solve the rest: `1 + dp[i-1][j-1]`. If they don't match, at least one of them is *not* in the answer — try dropping each and keep the better result: `max(dp[i-1][j], dp[i][j-1])`.

**Dry run** — `longestCommonSubsequence("abcde", "ace")`. Row = prefix of `abcde`, column = prefix of `ace`:

|   | "" | a | c | e |
|---|---|---|---|---|
| **""** | 0 | 0 | 0 | 0 |
| **a** | 0 | 1 | 1 | 1 |
| **b** | 0 | 1 | 1 | 1 |
| **c** | 0 | 1 | 2 | 2 |
| **d** | 0 | 1 | 2 | 2 |
| **e** | 0 | 1 | 2 | 3 |

Every time the row and column letters match (`a/a`, `c/c`, `e/e`), the value is the up-left diagonal + 1. Everywhere else it's the max of up and left. The answer is the bottom-right `3` ("ace").

**Run it:** `longestCommonSubsequence` in [`dp.go`](dp.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Longest Common Subsequence | Medium | https://leetcode.com/problems/longest-common-subsequence/ |
| 2 | Edit Distance | Medium | https://leetcode.com/problems/edit-distance/ |

**What to notice — Problem 1:** `dp` is sized `(m+1) x (n+1)`, not `m x n` — row/column 0 represent the "empty prefix," and Go zero-initializes them to 0, which is exactly the right base case. Characters match → extend the diagonal. No match → carry forward the best of dropping one character from either string.

**What to notice — Problem 2:** Same grid shape as Problem 1. No match → `1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])` (delete, insert, replace — try all three, take the cheapest). Base cases `dp[i][0] = i` and `dp[0][j] = j` must be set explicitly this time (converting to/from an empty string costs one op per character) — they aren't both 0 like Problem 1.

**Complexity:** both problems O(m·n) time and space. Space can drop to O(n) with two rolling rows (Edit Distance needs the diagonal from the previous row, so keep a `prev` variable for it).

**Common mistakes:**

- Off-by-one between `dp` indices and string indices: `dp[i][j]` is about prefixes of length `i` and `j`, so the characters are `text1[i-1]` and `text2[j-1]`.
- Edit Distance: leaving the base cases at 0. `dp[i][0] = i` (delete everything) and `dp[0][j] = j` (insert everything).
- Edit Distance: on a match, still adding 1. A matching character costs nothing: `dp[i][j] = dp[i-1][j-1]`.

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
        for c := target; c >= num; c-- { // iterate backward — 0/1 knapsack, each item used once
            dp[c] = dp[c] || dp[c-num]
        }
    }
    return dp[target]
}
```

**Why the backward loop matters — a real run:** `dp` is the set of reachable sums. Here's `canPartition` on `[1, 5, 11, 5]` (target `11`), showing which sums are reachable after each number:

| After number | Reachable sums |
|---|---|
| 1 | {0, 1} |
| 5 | {0, 1, 5, 6} |
| 11 | {0, 1, 5, 6, 11} |
| 5 | {0, 1, 5, 6, 10, 11} |

`11` is reachable → `true`. Now `[2, 6]` (target `4`, correct answer `false`):

| After number | Backward loop | Forward loop (bug) |
|---|---|---|
| 2 | {0, 2} | {0, 2, **4**} |
| 6 | {0, 2} | {0, 2, 4} |

With the forward loop, `dp[2]` became `true` while processing `2`, and later in the *same* pass `dp[4]` read it — so `2` was used twice. Going backward, `dp[c-num]` is always read before it can be updated in the current pass, so each number is used at most once. (Coin Change II wants unlimited reuse, so it loops **forward** on purpose.)

**Run it:** `canPartition` in [`dp.go`](dp.go) — `TestCanPartition` includes the `[2, 6]` case.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Partition Equal Subset Sum | Medium | https://leetcode.com/problems/partition-equal-subset-sum/ |
| 2 | Coin Change II | Medium | https://leetcode.com/problems/coin-change-ii/ |

**What to notice — Problem 1:** This is 0/1 knapsack collapsed to 1D: `dp[c]` = "is sum `c` reachable using items seen so far." The backward inner loop (`c := target; c >= num; c--`) is what limits each item to at most one use — a forward loop would let an item be reused, since it would read a `dp[c-num]` already updated in the same pass.

**What to notice — Problem 2:** Compare to Coin Change (Day 3): that problem minimizes coin *count*, this one counts the number of *combinations*. The loop order flips meaning: coins in the outer loop, amounts in the inner loop (forward this time) — this counts each combination once regardless of coin order. Recurrence: `dp[amt] += dp[amt-coin]`.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Partition Equal Subset Sum | O(n × sum/2) | O(sum/2) |
| Coin Change II | O(amount × coins) | O(amount) |

This is "pseudo-polynomial": it depends on the *value* of the target, not just the input length. A target of 10⁹ would make this approach too slow.

**Common mistakes:**

- 0/1 knapsack in 1D with a forward inner loop (see the run above).
- Coin Change II with amounts in the outer loop and coins in the inner loop — that counts *permutations* (`1+2` and `2+1` separately), not combinations.
- Forgetting the odd-sum early exit in Partition — an odd total can never split evenly.

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

| Pattern | State | Recurrence shape | Dimensions | Time |
|---------|-------|-------------------|------------|------|
| Linear (fixed lookback) | `dp[i]` | Looks back 1–2 steps | 1D (rolls to O(1)) | O(n) |
| Subsequence / unbounded | `dp[i]` | Scans all `j < i`, or reuses items | 1D + inner loop | O(n²) or O(amount × items) |
| Grid | `dp[i][j]` | From above or from the left | 2D (rolls to 1 row) | O(m·n) |
| Two strings | `dp[i][j]` | Branch on character match | 2D (rolls to 2 rows) | O(m·n) |
| Knapsack | `dp[i][capacity]` (or rolled 1D) | Skip vs take, capacity shrinks | 2D or rolled 1D | O(n × capacity) |

---

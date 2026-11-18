# Backtracking Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 13 — expands DFS week, Day 6 into a full week)  
**Goal:** Write the choose → explore → un-choose template from memory, pick between "start index" and "used array" instantly, de-duplicate without a set, and prune hard enough that exponential search finishes.

---

## What Is Backtracking?

**Backtracking** builds candidate solutions one choice at a time and **undoes** each choice when it's done exploring it: choose → explore (recurse) → un-choose. It's a DFS over a tree of decisions, where each path from the root is a partial answer.

What makes it fast enough is **pruning**: as soon as a partial answer can't possibly succeed, the whole branch below it is skipped. It's still exponential in the worst case — often because the number of answers itself is exponential.

**Typical problems it solves:**

- Generating all subsets, combinations, and permutations (including with duplicates)
- Combination sum: picking numbers that add up to a target
- Building strings position by position (phone letter combinations, valid parentheses)
- Cutting a string into valid pieces (palindrome partitioning, IP addresses)
- Constraint puzzles on a board (N-Queens, Sudoku, word search)

---

## The Core Template

```go
var result [][]int
var current []int

var backtrack func(start int)
backtrack = func(start int) {
    if isComplete(current) {
        result = append(result, append([]int(nil), current...)) // copy!
        return
    }
    for i := start; i < len(candidates); i++ {
        if !isValid(candidates[i], current) {
            continue // prune this branch
        }
        current = append(current, candidates[i]) // choose
        backtrack(i + 1)                         // explore
        current = current[:len(current)-1]      // un-choose
    }
}
backtrack(0)
```

Using a closure keeps `result` and `current` shared without passing pointers around — the same idea as the DFS week's `*[][]int`, with less noise.

**Three rules — never break them:**

1. **Copy when you record.** `current` is one slice reused by the entire search. Store `append([]int(nil), current...)`, never `current` itself (DFS week, Day 3 shows exactly what goes wrong).
2. **Every choose has a matching un-choose, in reverse order.** Whatever you changed before the recursive call — `current`, `used[i]`, a board cell, a remaining sum — restore it right after. The next sibling must see exactly the state the parent saw.
3. **Pick the right "what's still available" mechanism.** Order doesn't matter (subsets, combinations) → pass a `start` index so you only look to the right. Order matters (permutations) → loop from 0 with a `used []bool`. Choosing the wrong one gives duplicates or missing answers.

**Brute force → backtracking:** the brute force generates every full candidate (every subset, every arrangement), then filters. Backtracking builds candidates one choice at a time and **abandons a partial candidate as soon as it can't succeed**, which cuts off whole subtrees at once. It's still exponential in the worst case — the output itself is often exponential — but pruning decides whether it finishes in milliseconds or never.

**Complexity:** roughly (number of nodes in the decision tree) × (work per node). For outputs, it's at least (number of answers) × (cost to copy one) — e.g. O(n · 2ⁿ) for subsets, O(n · n!) for permutations.

**Run it:** every template in this file is a real function in [`backtracking.go`](backtracking.go), tested in [`backtracking_test.go`](backtracking_test.go). From the repo root: `go test ./backtracking/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** The decision tree  
**Key insight:** Every backtracking problem is a tree: each node is a partial answer, each edge is one choice, each leaf is a complete candidate. Backtracking is a DFS over that tree that never builds it explicitly.

**No problems today.** Do this instead:

1. Re-read DFS week, Day 6, including its dry run.
2. Read the template above slowly. Point at the three lines that are choose / explore / un-choose.
3. Close this file. Write the template from memory in Go (closure style).
4. Open this file. Compare. Fix any differences. Repeat twice more.
5. Answer this in your head: *Why does a `start` index prevent `[1, 2]` and `[2, 1]` from both appearing, while a `used` array allows both?*

   > With `start`, after picking index `i` you only consider indices `> i`, so elements are always picked in increasing index order — `[2, 1]` would require going left, which never happens. With `used`, every unused index is available at every level, so both orders are generated. Combinations want the first behavior, permutations want the second.

---

## Day 2 — Combinations and Subsets (with Pruning)

**Pattern:** `start` index; record at the right depth; cut branches that can't finish  
**Key insight:** If you need `k` more elements and fewer than `k` are left, the branch is dead before it starts. One bound on the loop removes it.

**Template adaptation:**

```go
// every k-element combination of 1..n
func combine(n, k int) [][]int {
    var result [][]int
    var current []int

    var backtrack func(start int)
    backtrack = func(start int) {
        if len(current) == k {
            result = append(result, append([]int(nil), current...))
            return
        }
        need := k - len(current)
        for i := start; i <= n-need+1; i++ { // leave room for the `need` picks
            current = append(current, i)
            backtrack(i + 1)
            current = current[:len(current)-1]
        }
    }
    backtrack(1)
    return result
}
```

**Why the bound is `n - need + 1`:** if you pick `i` now, the remaining `need - 1` picks must come from `i+1 … n`, which has `n - i` numbers. That's only possible when `n - i >= need - 1`, i.e. `i <= n - need + 1`. Anything larger is guaranteed to fail.

**How much pruning buys — real call counts:**

| `n`, `k` | Calls without pruning | Calls with pruning | Results |
|---|---|---|---|
| 4, 2 | 11 | 10 | 6 |
| 10, 4 | 386 | 330 | 210 |
| 20, 10 | 616,666 | 352,716 | 184,756 |
| 20, 18 | 1,048,555 | 1,330 | 190 |

When `k` is close to `n`, almost every branch dies — and the unpruned version wastes about a million calls to produce 190 answers.

**Run it:** `combine` in [`backtracking.go`](backtracking.go) (this is LeetCode 77 — do it as a warm-up).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Subsets II | Medium | https://leetcode.com/problems/subsets-ii/ |
| 2 | Combination Sum III | Medium | https://leetcode.com/problems/combination-sum-iii/ |

**What to notice — Problem 1:** Duplicates in the input. Sort first, then at each level skip `nums[i]` if it equals `nums[i-1]` **and** `i > start` — the same value was already tried as the choice for this position. Record at every node (every node is a subset).

**What to notice — Problem 2:** Choose `k` numbers from 1–9 summing to `n`. Two prunes: the count bound from the template, and `if i > remaining { break }` (sorted candidates — everything after is too big too).

**Complexity:** Subsets II O(n · 2ⁿ). Combination Sum III is bounded by C(9, k) — tiny.

**Common mistakes:**

- Subsets II: skipping duplicates with `i > 0` instead of `i > start` — that also skips valid subsets like `[1, 1]`.
- Forgetting to copy `current` when recording.
- Recursing with `start + 1` instead of `i + 1` — produces repeats and misses combinations.

---

## Day 3 — Combination Sum: Reuse and Sorted Pruning

**Pattern:** Track the remaining target; recurse with `i` (reuse allowed) or `i + 1` (each once); sort so you can `break`  
**Key insight:** Sorting turns "this candidate is too big" into "this candidate **and every one after it** is too big" — `continue` becomes `break`.

**Template adaptation:**

```go
// every multiset of candidates (each reusable) summing to target
func combinationSum(candidates []int, target int) [][]int {
    sort.Ints(candidates)
    var result [][]int
    var current []int

    var backtrack func(start, remaining int)
    backtrack = func(start, remaining int) {
        if remaining == 0 {
            result = append(result, append([]int(nil), current...))
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining {
                break // sorted: the rest are even bigger
            }
            current = append(current, candidates[i])
            backtrack(i, remaining-candidates[i]) // i, not i+1: may reuse it
            current = current[:len(current)-1]
        }
    }
    backtrack(0, target)
    return result
}
```

**Dry run** — `combinationSum([2, 3, 6, 7], 7)`:

```
take 2 → [2], remaining 5
    take 2 → [2, 2], remaining 3
        take 2 → [2, 2, 2], remaining 1
            2 > 1 → break (and skip [3, 6, 7])
        take 3 → [2, 2, 3], remaining 0
            [2, 2, 3] ✓ record
        6 > 3 → break (and skip [7])
    take 3 → [2, 3], remaining 2
        3 > 2 → break (and skip [6, 7])
    6 > 5 → break (and skip [7])
take 3 → [3], remaining 4
    take 3 → [3, 3], remaining 1
        3 > 1 → break (and skip [6, 7])
    6 > 4 → break (and skip [7])
take 6 → [6], remaining 1
    6 > 1 → break (and skip [7])
take 7 → [7], remaining 0
    [7] ✓ record
```

Every `break` line cut off the rest of that level. Recursing with `i` (not `i + 1`) is why `[2, 2, 3]` is possible; recursing with `start = i` (not `0`) is why `[3, 2, 2]` never appears.

**Run it:** `combinationSum` in [`backtracking.go`](backtracking.go) (this is LeetCode 39 — do it as a warm-up).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Combination Sum II | Medium | https://leetcode.com/problems/combination-sum-ii/ |
| 2 | Partition to K Equal Sum Subsets | Medium | https://leetcode.com/problems/partition-to-k-equal-sum-subsets/ |

**What to notice — Problem 1:** Each candidate used once (`i + 1`) and the input has duplicates — combine today's template with Day 2's `i > start` duplicate skip.

**What to notice — Problem 2:** Fill `k` buckets of size `sum / k`. Sort descending (big numbers fail fast), and when placing a number, skip any bucket whose current sum equals a bucket you already tried for this number — identical buckets lead to identical subtrees. Without these prunes it times out.

**Complexity:** exponential in the worst case; pruning is what makes the real inputs fast.

**Common mistakes:**

- Using `continue` instead of `break` after sorting — correct but slower; or using `break` **without** sorting — wrong.
- Combination Sum II: `i + 1` without the duplicate skip → repeated combinations.
- Partition to K: forgetting the "skip identical bucket" prune.

---

## Day 4 — Permutations and Position-by-Position Building

**Pattern:** `used []bool`, loop from 0; or "fill position `p` with each allowed option"  
**Key insight:** To skip duplicate permutations without a set, sort and only allow the second copy of a value once the first copy is **already in use** — so equal values always appear in the same relative order.

**Template adaptation:**

```go
// every distinct permutation of nums (which may contain duplicates)
func permuteUnique(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    current := make([]int, 0, len(nums))
    used := make([]bool, len(nums))

    var backtrack func()
    backtrack = func() {
        if len(current) == len(nums) {
            result = append(result, append([]int(nil), current...))
            return
        }
        for i := range nums {
            if used[i] {
                continue
            }
            if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
                continue // same value as an unused earlier copy: would repeat a branch
            }
            used[i] = true
            current = append(current, nums[i])
            backtrack()
            current = current[:len(current)-1]
            used[i] = false
        }
    }
    backtrack()
    return result
}
```

**Why `!used[i-1]` is the right condition:** if the earlier copy is unused, then at this level it was already tried as the choice for this position — and picking the later copy now would rebuild exactly the same subtree. If the earlier copy *is* used (it's above us in the current path), the later copy is the next legitimate pick.

**Dry run** — `permuteUnique([1, 1, 2])`:

```
pick index 0 (1) → [1]
    pick index 1 (1) → [1, 1]
        pick index 2 (2) → [1, 1, 2]
            [1, 1, 2] ✓
    pick index 2 (2) → [1, 2]
        pick index 1 (1) → [1, 2, 1]
            [1, 2, 1] ✓
skip index 1 (second 1 while first 1 is unused)
pick index 2 (2) → [2]
    pick index 0 (1) → [2, 1]
        pick index 1 (1) → [2, 1, 1]
            [2, 1, 1] ✓
    skip index 1 (second 1 while first 1 is unused)
```

Three permutations instead of 3! = 6, and no set needed. Each "skip" cut a branch that would have produced a duplicate.

**Run it:** `permuteUnique` in [`backtracking.go`](backtracking.go) (this is LeetCode 47 — do it as a warm-up).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Letter Combinations of a Phone Number | Medium | https://leetcode.com/problems/letter-combinations-of-a-phone-number/ |
| 2 | Generate Parentheses | Medium | https://leetcode.com/problems/generate-parentheses/ |

**What to notice — Problem 1:** Position-by-position: the tree depth is the number of digits, and at depth `p` the options are the letters of `digits[p]`. No `used`, no `start` — the position *is* the state. Return an empty list (not `[""]`) for empty input.

**What to notice — Problem 2:** Position-by-position with two counters. You may add `(` while `open < n`, and `)` while `close < open`. That second rule is the pruning: it makes invalid strings impossible to build, so nothing needs to be filtered at the end.

**Complexity:** Permutations II O(n · n!). Letter Combinations O(n · 4ⁿ). Generate Parentheses O(4ⁿ / √n) (the Catalan number of results).

**Common mistakes:**

- Permutations II: using `used[i-1]` instead of `!used[i-1]` — still de-duplicates but explores far more of the tree. Know why the `!` version prunes earlier.
- Generate Parentheses: generating all 2²ⁿ strings and filtering valid ones.
- Building strings with `+=` in Go at every level — fine for these sizes, but a `[]byte` buffer with append/truncate is the backtracking-native way.

---

## Day 5 — Partitioning a String

**Pattern:** At each position, try every possible next piece; recurse on the rest  
**Key insight:** "Cut a string into valid pieces" is choosing the **end** of the next piece. The state is just "where the next piece starts".

**Template adaptation:**

```go
// every way to cut s into pieces that all satisfy isValid
func partition(s string, isValid func(string) bool) [][]string {
    var result [][]string
    var pieces []string

    var backtrack func(start int)
    backtrack = func(start int) {
        if start == len(s) {
            result = append(result, append([]string(nil), pieces...))
            return
        }
        for end := start + 1; end <= len(s); end++ {
            piece := s[start:end]
            if !isValid(piece) {
                continue
            }
            pieces = append(pieces, piece)
            backtrack(end)
            pieces = pieces[:len(pieces)-1]
        }
    }
    backtrack(0)
    return result
}
```

**Why it's generic:** the search shape never changes — only `isValid` does. Palindromes, valid IP octets, dictionary words: plug in a different predicate.

**Dry run** — `partition("aab", isPalindrome)`:

```
cut "a" → [a]
    cut "a" → [a a]
        cut "b" → [a a b]
            [a a b] ✓
    "ab" not a palindrome → skip
cut "aa" → [aa]
    cut "b" → [aa b]
        [aa b] ✓
"aab" not a palindrome → skip
```

**Run it:** `partition` in [`backtracking.go`](backtracking.go) — the test runs it with a palindrome predicate and with a "length ≤ 2" predicate.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Palindrome Partitioning | Medium | https://leetcode.com/problems/palindrome-partitioning/ |
| 2 | Restore IP Addresses | Medium | https://leetcode.com/problems/restore-ip-addresses/ |

**What to notice — Problem 1:** The template with a palindrome check. Optional upgrade: precompute `isPal[i][j]` with DP (DP week, Day 5 style) so each check is O(1).

**What to notice — Problem 2:** Same shape with two extra prunes: exactly 4 pieces (stop if `len(pieces) == 4` before the end, or if the remaining characters can't fill the remaining pieces at 1–3 each), and each piece must be `0–255` with no leading zero unless it's `"0"`.

**Complexity:** Palindrome Partitioning O(n · 2ⁿ) — up to 2ⁿ⁻¹ ways to cut. Restore IP is bounded by 3⁴ = 81 — constant.

**Common mistakes:**

- `end < len(s)` instead of `end <= len(s)` — the last piece can never reach the end of the string.
- IP: accepting `"01"` or `"256"`.
- Re-checking palindromes from scratch in a hot loop when the input is large.

---

## Day 6 — Constraint Search on a Board

**Pattern:** Place one thing per step; keep O(1) "is this square attacked / used?" lookups; undo on the way back  
**Key insight:** The pruning lookups are the algorithm. Checking a placement by scanning the board is O(n) per check; keeping columns and diagonals in boolean arrays makes it O(1).

**Template adaptation:**

```go
// count ways to place n non-attacking queens
func totalNQueens(n int) int {
    cols := make([]bool, n)
    diag := make([]bool, 2*n-1)     // r + c is constant along a "/" diagonal
    antiDiag := make([]bool, 2*n-1) // r - c + n - 1 is constant along a "\" diagonal
    count := 0

    var place func(r int)
    place = func(r int) {
        if r == n {
            count++
            return
        }
        for c := 0; c < n; c++ {
            if cols[c] || diag[r+c] || antiDiag[r-c+n-1] {
                continue // attacked
            }
            cols[c], diag[r+c], antiDiag[r-c+n-1] = true, true, true
            place(r + 1)
            cols[c], diag[r+c], antiDiag[r-c+n-1] = false, false, false
        }
    }
    place(0)
    return count
}
```

**Why one queen per row:** two queens in a row always attack each other, so every solution has exactly one queen per row. Making "row" the recursion depth removes that constraint for free — only columns and diagonals need checking.

**Dry run** — `n = 4`. The search makes **17** calls to `place` and finds 2 solutions (queen column for rows 0–3):

```
solution [1, 3, 0, 2]      solution [2, 0, 3, 1]
. Q . .                    . . Q .
. . . Q                    Q . . .
Q . . .                    . . . Q
. . Q .                    . Q . .
```

Brute force with one queen per row would check 4⁴ = 256 complete boards. The diagonal arrays killed most branches after one or two queens.

**Run it:** `totalNQueens` in [`backtracking.go`](backtracking.go) — tested against the known counts (8 queens → 92).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | N-Queens | Hard | https://leetcode.com/problems/n-queens/ |
| 2 | Word Search | Medium | https://leetcode.com/problems/word-search/ |

**What to notice — Problem 1:** The template plus a `queens []int` (column per row) to build the board strings at each solution.

**What to notice — Problem 2:** Grid DFS (DFS week, Day 5) plus backtracking: mark the cell as used (e.g. overwrite it with `'#'`), recurse into 4 neighbors for the next letter, then **restore** the original letter. Unlike flood fill, the mark must be undone — another path may need that cell.

**Complexity:** N-Queens O(n!) roughly. Word Search O(R · C · 3ᴸ) — each step after the first has at most 3 new directions.

**Common mistakes:**

- N-Queens: using `r - c` as an index without the `+ n - 1` offset (negative index panic).
- Word Search: forgetting to restore the cell — later starting points see a corrupted board.
- Word Search: not returning early once found — keeps searching the whole board.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Every backtracking problem is four decisions: what a node's state is, what the choices at a node are, how to avoid duplicates (`start` / `used` / sort-and-skip), and how to prune.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before writing code, draw the first two levels of the decision tree on paper and mark one branch your pruning would cut.

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the choose / explore / un-choose template from memory
- [ ] When do you use a `start` index vs a `used` array?
- [ ] How do you skip duplicates in Subsets II, and why `i > start` and not `i > 0`?
- [ ] Why does sorting let Combination Sum `break` instead of `continue`?
- [ ] In Permutations II, what does `!used[i-1]` guarantee?
- [ ] How does N-Queens check a diagonal in O(1)?

---

## Pattern Cheatsheet

| Pattern | Availability | Recurse with | Duplicate handling | Typical prune | Time |
|---------|--------------|--------------|--------------------|---------------|------|
| Combinations / subsets | `start` index | `i + 1` | sort; skip `i > start && nums[i] == nums[i-1]` | not enough elements left | O(n · 2ⁿ) |
| Combination sum (reuse) | `start` index | `i` | sort | `break` when candidate > remaining | exponential |
| Permutations | `used []bool` | — (loop from 0) | sort; skip when previous equal is unused | — | O(n · n!) |
| Position by position | the position | `p + 1` | — | counters (open/close) | output-sized |
| Partitioning | start of next piece | `end` | — | `isValid(piece)` | O(n · 2ⁿ) |
| Board placement | row = depth | `r + 1` | — | O(1) column/diagonal arrays | ~O(n!) |

---

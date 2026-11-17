# Greedy Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 12 — after Intervals; contrast with the DP week throughout)  
**Goal:** Recognize when a locally best choice is provably safe, justify it out loud with an exchange argument, and know the classic shapes: running best, farthest reach, sort-and-match, last occurrence, and regret with a heap.

---

## The Core Template

Greedy isn't one piece of code — it's a decision rule applied in a single pass, usually after a sort:

```go
sort.Slice(items, byTheRightKey) // often the most important line

state := initial
for _, item := range items {
    if isSafeToTake(item, state) { // the greedy choice: best-looking option right now
        state = take(item, state)
    }
}
return answerFrom(state)
```

**Three rules — never break them:**

1. **A greedy choice needs a reason, not a feeling.** Before coding, say *why* the local choice can't hurt: "any optimal answer can be swapped to include this choice without getting worse" (exchange argument), or "after every step, greedy is at least as far along as any other strategy" (stays ahead). If you can't say it, try a counterexample.
2. **Try to break it with a tiny counterexample first.** Two or three elements are usually enough. If you find one, it's DP (or search), not greedy.
3. **The sort key IS the algorithm.** Most greedy bugs are sorting by the wrong thing (start vs end, ascending vs descending, value vs ratio). Decide the key from the exchange argument, not from habit.

**Brute force → greedy:** the brute force tries every combination of choices (exponential) or every sub-problem (DP, polynomial). Greedy commits to one choice per step and never reconsiders, so it's usually O(n) or O(n log n). That speed is only correct when rule 1 holds — which is the whole difficulty.

**Greedy vs DP — the classic counterexample:** coin change with coins `[1, 3, 4]` and amount `6`. Greedy (take the biggest coin that fits) picks `4 + 1 + 1` = 3 coins. The optimal answer is `3 + 3` = 2 coins. Taking `4` *looked* best but ruled out the better combination — so Coin Change needs DP (DP week, Day 3). Greedy works for real-world coin systems like `[1, 5, 10, 25]` only because of their specific values.

**Complexity:** usually O(n log n) for the sort plus O(n) for the pass; O(n log n) overall with a heap.

**Run it:** every template in this file is a real function in [`greedy.go`](greedy.go), tested in [`greedy_test.go`](greedy_test.go). From the repo root: `go test ./greedy/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Exchange argument  
**Key insight:** Greedy is correct when you can take *any* optimal solution and transform it, step by step, into the greedy solution without ever making it worse.

**No problems today.** Do this instead:

1. Read the template and the coin-change counterexample above.
2. Re-read the Intervals week, Day 4 ("sort by end"). Write its exchange argument from memory in 3 sentences.
3. Find your own counterexample for "always take the largest coin" with coins `[1, 5, 6]` (hint: amount 10).
4. Answer this in your head: *How is greedy different from DP?*

   > DP considers every choice at each step and keeps the best result for each sub-problem — it's safe even when the locally best choice is wrong. Greedy makes one choice per step and never revisits it, so it's faster but only correct when you can prove that the local choice is always part of some optimal answer. When in doubt, DP is the safe default; greedy is the optimization you justify.

---

## Day 2 — Running Best (Kadane)

**Pattern:** Carry one or two running values through a single pass; at each step decide "extend or restart"  
**Key insight:** A subarray ending at `i` either extends the best subarray ending at `i-1` or starts fresh at `i`. If the running sum is negative, it can only hurt what comes next — drop it.

**Template adaptation:**

```go
func maxSubArray(nums []int) int {
    cur, best := nums[0], nums[0] // best subarray ending here, best overall
    for _, x := range nums[1:] {
        cur = max(x, cur+x) // extend, or restart at x
        best = max(best, cur)
    }
    return best
}
```

**Why it works:** `cur` is "the best sum of a subarray **ending at** the current index". Any subarray ending at `i` is either just `[x]` or a subarray ending at `i-1` plus `x` — and the best of the latter is `cur + x`. So `max(x, cur + x)` is exact. (This is really a DP with one state rolled into a variable — Kadane sits right on the greedy/DP border.)

**Dry run** — `maxSubArray([-2, 1, -3, 4, -1, 2, 1, -5, 4])`:

| `i` | `x` | Extend or restart? | `cur` | `best` |
|---|---|---|---|---|
| 0 | -2 | start | -2 | -2 |
| 1 | 1 | restart (1 > -1) | 1 | 1 |
| 2 | -3 | extend (1 + -3 = -2) | -2 | 1 |
| 3 | 4 | restart (4 > 2) | 4 | 4 |
| 4 | -1 | extend (4 + -1 = 3) | 3 | 4 |
| 5 | 2 | extend (3 + 2 = 5) | 5 | 5 |
| 6 | 1 | extend (5 + 1 = 6) | 6 | **6** |
| 7 | -5 | extend (6 + -5 = 1) | 1 | 6 |
| 8 | 4 | extend (1 + 4 = 5) | 5 | 6 |

The best subarray is `[4, -1, 2, 1]`. At `i = 4` it extended through a `-1` because `3` is still positive — a negative *element* is fine; a negative *running sum* is what triggers a restart.

**Run it:** `maxSubArray` in [`greedy.go`](greedy.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Maximum Subarray | Medium | https://leetcode.com/problems/maximum-subarray/ |
| 2 | Best Time to Buy and Sell Stock | Easy | https://leetcode.com/problems/best-time-to-buy-and-sell-stock/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** The running value is the **lowest price so far**. At each day, the best profit selling today is `price - minSoFar`. One pass, two variables. (It's Kadane in disguise: the daily price *differences* form an array whose maximum subarray is the answer.)

**Complexity:** both O(n) time, O(1) space.

**Common mistakes:**

- Initializing `best` (or `cur`) to 0 — wrong for all-negative arrays; the answer must be a non-empty subarray.
- Stock: updating `minSoFar` *before* computing today's profit is fine; computing profit against a future minimum is not — you must buy before you sell.

---

## Day 3 — Farthest Reach

**Pattern:** Track the farthest position reachable so far; fail the moment you stand somewhere unreachable  
**Key insight:** You never need to know *which* jumps get you somewhere — only how far the best combination of jumps could reach.

**Template adaptation:**

```go
func canReach(nums []int) bool {
    farthest := 0
    for i, jump := range nums {
        if i > farthest {
            return false // index i can't be reached by any sequence of jumps
        }
        farthest = max(farthest, i+jump)
    }
    return true
}
```

**Why it works (stays ahead):** invariant — every index `≤ farthest` is reachable. It holds at the start (index 0). If index `i ≤ farthest` is reachable, then so is everything up to `i + nums[i]`, so extending `farthest` keeps the invariant. The first index beyond `farthest` can't be reached by anything before it, and nothing after it is reachable either.

**Dry run** — `[2, 3, 1, 1, 4]` (reachable) and `[3, 2, 1, 0, 4]` (stuck):

| `i` | `nums[i]` | `[2,3,1,1,4]`: `farthest` after | `[3,2,1,0,4]`: `farthest` after |
|---|---|---|---|
| 0 | 2 / 3 | 2 | 3 |
| 1 | 3 / 2 | 4 | 3 |
| 2 | 1 / 1 | 4 | 3 |
| 3 | 1 / 0 | 4 | 3 |
| 4 | 4 / 4 | 8 → **true** | 4 > 3 → **false** |

In the second array every path funnels into index 3, whose jump is 0 — `farthest` gets stuck at 3.

**Run it:** `canReach` in [`greedy.go`](greedy.go). (This is LeetCode 55, Jump Game — do it as a warm-up.)

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Jump Game II | Medium | https://leetcode.com/problems/jump-game-ii/ |
| 2 | Gas Station | Medium | https://leetcode.com/problems/gas-station/ |

**What to notice — Problem 1:** Minimum jumps = BFS levels without a queue. Keep `currentEnd` (edge of what this many jumps can reach) and `farthest`. When `i` hits `currentEnd`, you must jump again: `jumps++`, `currentEnd = farthest`. Same shape as covering a range (Intervals, Day 5).

**What to notice — Problem 2:** Two facts: (1) if total gas < total cost, it's impossible; (2) otherwise, scan once with a running tank, and whenever it goes negative at station `i`, no station from the current start to `i` can be the answer — restart at `i + 1` with an empty tank. Kadane's "restart when negative", in another costume.

**Complexity:** all O(n) time, O(1) space.

**Common mistakes:**

- Jump Game II: looping to the last index inclusive and counting an extra jump. Stop at `n - 2` — standing on the last index needs no jump.
- Gas Station: trying every start (O(n²)). The restart argument is what makes it O(n).

---

## Day 4 — Sort, Then Match

**Pattern:** Sort both sides; match the smallest thing that works with the smallest need it can satisfy  
**Key insight:** Using a bigger item than necessary on a small need wastes it. Give each need the smallest item that satisfies it.

**Template adaptation:**

```go
// match items to needs (item satisfies need if item >= need), each used at most once
func maxSatisfied(needs, items []int) int {
    sort.Ints(needs)
    sort.Ints(items)

    matched := 0
    for _, item := range items {
        if matched < len(needs) && item >= needs[matched] {
            matched++ // item goes to the smallest unmet need
        }
    }
    return matched
}
```

**Why it works (exchange):** suppose an optimal matching gives the smallest need a bigger item than the smallest one that works. Swap the two items between their needs: the smallest need is still satisfied, and the other need gets the bigger item, so it's still satisfied too. Same count. Repeating this turns any optimal matching into the greedy one.

**Dry run** — needs `[5, 1, 4, 2]` → sorted `[1, 2, 4, 5]`; items `[3, 6, 3, 1]` → sorted `[1, 3, 3, 6]`:

| Item | Smallest unmet need | Satisfies? | Matched |
|---|---|---|---|
| 1 | 1 | yes | 1 |
| 3 | 2 | yes | 2 |
| 3 | 4 | no — too small for anyone left | 2 |
| 6 | 4 | yes | 3 |

The second `3` is skipped for good: every unmet need is ≥ 4, so it can never be used.

**Run it:** `maxSatisfied` in [`greedy.go`](greedy.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Assign Cookies | Easy | https://leetcode.com/problems/assign-cookies/ |
| 2 | Boats to Save People | Medium | https://leetcode.com/problems/boats-to-save-people/ |

**What to notice — Problem 1:** This IS the template (greed factors are needs, cookie sizes are items).

**What to notice — Problem 2:** Sort, then opposite-ends two pointers. The heaviest person always takes a boat; they share it with the lightest person if the two fit together. If the lightest can't ride with the heaviest, they can't ride with anyone heavier either — so pairing them with the lightest is never a mistake.

**Complexity:** both O(n log n) for the sorts.

**Common mistakes:**

- Cookies: giving the biggest cookie to the greediest child first and stopping when it fails — you must keep scanning.
- Boats: trying to pair the heaviest with the *second* heaviest that fits — that's a search; the lightest is always the safe partner.

---

## Day 5 — Last Occurrence and Two Passes

**Pattern:** Precompute something about the whole input (last index, or a left-to-right constraint), then make one greedy pass  
**Key insight:** A greedy cut is safe when you already know nothing to the right can force you to undo it — and "last occurrence" tells you exactly that.

**Template adaptation:**

```go
func partitionLabels(s string) []int {
    var last [26]int
    for i := 0; i < len(s); i++ {
        last[s[i]-'a'] = i // where each letter appears for the final time
    }

    var sizes []int
    start, end := 0, 0
    for i := 0; i < len(s); i++ {
        end = max(end, last[s[i]-'a']) // this part must reach at least here
        if i == end { // every letter seen so far is finished — safe to cut
            sizes = append(sizes, end-start+1)
            start = i + 1
        }
    }
    return sizes
}
```

**Why cutting as early as possible is right:** a cut at `i` is legal iff every letter in `s[start..i]` has its last occurrence at or before `i` — which is exactly `i == end`. Cutting at the earliest legal point never prevents a later cut, so it maximizes the number of parts.

**Dry run** — `partitionLabels("abacbdeedf")`, with last indexes `a:2, b:4, c:3, d:8, e:7, f:9`:

| `i` | char | last[char] | `end` | Cut? |
|---|---|---|---|---|
| 0 | a | 2 | 2 | — |
| 1 | b | 4 | 4 | — |
| 2 | a | 2 | 4 | — |
| 3 | c | 3 | 4 | — |
| 4 | b | 4 | 4 | yes → size 5 (`abacb`) |
| 5 | d | 8 | 8 | — |
| 6 | e | 7 | 8 | — |
| 7 | e | 7 | 8 | — |
| 8 | d | 8 | 8 | yes → size 4 (`deed`) |
| 9 | f | 9 | 9 | yes → size 1 (`f`) |

**Run it:** `partitionLabels` in [`greedy.go`](greedy.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Partition Labels | Medium | https://leetcode.com/problems/partition-labels/ |
| 2 | Candy | Hard | https://leetcode.com/problems/candy/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Two constraints (left neighbor, right neighbor) are hard to satisfy in one pass. Satisfy them one at a time: left-to-right pass gives each child `left[i-1] + 1` if they out-rate their left neighbor; right-to-left pass does the same with the right neighbor. Each child gets `max(left, right)` — the minimum that satisfies both.

**Complexity:** Problem 1 O(n) time, O(1) space (26 letters). Problem 2 O(n) time, O(n) space.

**Common mistakes:**

- Partition Labels: cutting when the *current* letter's last occurrence is `i`, instead of when `end` (the max over all letters so far) is `i`.
- Candy: doing only one pass, or adding the two passes instead of taking the max.

---

## Day 6 — Regret Greedy with a Heap

**Pattern:** Pass over options without committing; keep them in a heap; when you're forced to act, take the best option you've already passed  
**Key insight:** You don't have to decide at each station whether to stop. Decide *later* — when you'd otherwise fail — and at that point choose with full hindsight.

**Template adaptation:**

```go
// fewest refueling stops to reach target; stations are [position, fuel], sorted by position
func minRefuelStops(target, startFuel int, stations [][]int) int {
    passed := &maxHeap{} // fuel amounts of stations we've driven past
    fuel, stops, i := startFuel, 0, 0

    for fuel < target { // fuel doubles as "how far we can reach"
        for i < len(stations) && stations[i][0] <= fuel {
            heap.Push(passed, stations[i][1]) // reachable — remember it, don't stop yet
            i++
        }
        if passed.Len() == 0 {
            return -1 // nothing to fall back on
        }
        fuel += heap.Pop(passed).(int) // "in hindsight, I stopped at the biggest one"
        stops++
    }
    return stops
}
```

**Why it works:** to reach further with `k` stops, you want the `k` stations with the most fuel among those you could have reached. Every station you've passed was reachable, so retroactively "having stopped" at the biggest is always allowed — and if you need one more stop, the biggest remaining one is never a worse choice than any other.

**Dry run** — `target = 100`, `startFuel = 10`, stations `[[10,60], [20,30], [30,30], [60,40]]`:

| Fuel (= reach) | Stations newly passed | Heap | Refuel from | Stops |
|---|---|---|---|---|
| 10 | @10: 60 | `[60]` | 60 → reach 70 | 1 |
| 70 | @20: 30, @30: 30, @60: 40 | `[40, 30, 30]` | 40 → reach 110 | 2 |

Reach 110 ≥ 100 → **2** stops (at positions 10 and 60). The two 30-fuel stations were passed and never used.

**Run it:** `minRefuelStops` in [`greedy.go`](greedy.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Minimum Number of Refueling Stops | Hard | https://leetcode.com/problems/minimum-number-of-refueling-stops/ |
| 2 | Course Schedule III | Hard | https://leetcode.com/problems/course-schedule-iii/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Sort courses by deadline. Take every course; keep their durations in a max-heap and a running `time`. If `time` passes the current deadline, drop the **longest** course taken so far (pop the heap) — it frees the most time and the count only drops by one, the same as dropping any other course.

**Complexity:** both O(n log n).

**Common mistakes:**

- Refueling: stopping at every station "just in case", or greedily stopping at the *nearest* station.
- Course Schedule III: sorting by duration instead of deadline, or dropping the course you just added instead of the longest one.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** The skill isn't the code — it's the one-sentence proof. If you can say why the greedy choice is safe, the code is a single pass.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. For each problem you redo, say the exchange argument (or "stays ahead" argument) out loud before coding.
4. Then try to break it: invent a 3-element input where a slightly different greedy rule (wrong sort key, wrong tie-break) fails.

**Self-assessment — can you answer these without hesitation?**

- [ ] What's an exchange argument? Give one for "sort intervals by end"
- [ ] Why does greedy fail for coin change with `[1, 3, 4]`, and what should you use instead?
- [ ] In Kadane's algorithm, when do you restart, and why is that safe?
- [ ] In Jump Game, what invariant does `farthest` maintain?
- [ ] What does "regret greedy" mean, and why does it need a heap?

---

## Pattern Cheatsheet

| Pattern | Sort / precompute | Greedy rule | Proof idea | Time |
|---------|-------------------|-------------|------------|------|
| Running best (Kadane) | — | extend or restart; restart when running value goes negative | best ending here is one of two options | O(n) |
| Farthest reach | — | track max reach; fail when `i > farthest` | everything ≤ farthest is reachable | O(n) |
| Sort and match | sort both sides | smallest item that fits the smallest need | swap a bigger item for a smaller one | O(n log n) |
| Last occurrence | last index of each value | cut as soon as everything seen has ended | earliest legal cut never blocks a later one | O(n) |
| Two-pass constraints | — | satisfy left constraints, then right, take max | each pass is the minimum for its side | O(n) |
| Regret with heap | sort by position / deadline | defer choices; when forced, take the best passed option | hindsight choice is always allowed | O(n log n) |

---

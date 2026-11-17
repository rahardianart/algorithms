# Two Pointer Interview Prep — 7-Day Study Plan

**Level:** Beginner-Intermediate (completed BFS and DFS weeks)  
**Goal:** Write both Two Pointer templates from memory, recognize patterns, implement in-place array modification.

---

## The Two Templates

### Opposite Ends

Start from both sides of the array, shrink toward the center.

```go
left, right := 0, len(arr)-1

for left < right {
    // process arr[left] and arr[right]

    if someCondition {
        left++
    } else {
        right--
    }
}
```

### Same Direction (Fast/Slow)

Both pointers move left to right. `slow` writes, `fast` scans.

```go
slow := 0

for fast := 0; fast < len(arr); fast++ {
    if someCondition {
        arr[slow] = arr[fast]
        slow++
    }
}
```

**Three rules — never break them:**

1. **Opposite ends** — start from both sides, shrink toward center. Works when array is sorted or when looking for a pair/triplet.
2. **Same direction** — `slow` tracks where to write, `fast` scans ahead. Works for in-place removal or overwrite.
3. **The condition drives the pointer** — which pointer moves and when is the only thing that changes between problems. Everything else stays the same.

**Why Two Pointer?** It eliminates nested loops. A brute force O(n²) search becomes O(n) by using two pointers instead of two nested loops.

**Brute force → two pointers:** "Find a pair with property X" brute-forces as two nested loops over all n² pairs. Opposite-ends pointers work because each comparison lets you **throw away a whole row of pairs at once**: one pointer moves, and every pair it could still have formed is proven useless. n² pairs, eliminated n at a time, means n steps. Same-direction pointers replace "build a new filtered array" (O(n) extra space) with an in-place overwrite (O(1) extra space).

**Run it:** every function in this file is real code in [`twopointer.go`](twopointer.go), tested in [`twopointer_test.go`](twopointer_test.go). From the repo root: `go test ./two-pointer/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Both templates  
**Key insight:** Two pointers eliminate the need for nested loops — O(n²) becomes O(n).

**No problems today.** Do this instead:

1. Read both templates above slowly. Understand each line.
2. Close this file. Write both templates from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *When do you use opposite ends vs same direction?*

   > Use **opposite ends** when the array is sorted and you're searching for a pair, or when you're maximizing/minimizing something computed from both sides. Use **same direction** when you're filtering or overwriting elements in-place and relative order must be preserved.

---

## Day 2 — Opposite Ends: Target Sum

**Pattern:** Shrink from both ends toward a target  
**Key insight:** On a sorted array, if the sum is too large move `right` left; if too small move `left` right. You always make progress — one pointer moves per iteration.

**Template adaptation:**

```go
func twoSumSorted(numbers []int, target int) [2]int {
    left, right := 0, len(numbers)-1

    for left < right {
        sum := numbers[left] + numbers[right]
        if sum == target {
            return [2]int{left + 1, right + 1} // 1-indexed
        } else if sum < target {
            left++  // need larger sum
        } else {
            right-- // need smaller sum
        }
    }
    return [2]int{-1, -1} // not found
}
```

**3Sum extension — fix one, two-pointer the rest:**

```go
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    result := [][]int{}

    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] {
            continue // skip duplicate fixed number
        }
        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                for left < right && nums[left] == nums[left+1] {
                    left++ // skip duplicate left
                }
                for left < right && nums[right] == nums[right-1] {
                    right-- // skip duplicate right
                }
                left++
                right--
            } else if sum < 0 {
                left++
            } else {
                right--
            }
        }
    }
    return result
}
```

**Why it works:** Suppose `numbers[left] + numbers[right] < target`. The array is sorted, so `numbers[right]` is the **largest** value `left` could ever be paired with — and even that's too small. Every pair that uses `left` is therefore too small, so `left` can be discarded for good: `left++`. The mirror argument discards `right` when the sum is too big. Each step rules out one element permanently and never rules out the answer.

**Dry run** — `[1, 3, 4, 6, 8, 11]`, target `10`:

| `left` | `right` | `a[left] + a[right]` | Move |
|---|---|---|---|
| 0 | 5 | 1 + 11 = 12 | too big → `right--` |
| 0 | 4 | 1 + 8 = 9 | too small → `left++` |
| 1 | 4 | 3 + 8 = 11 | too big → `right--` |
| 1 | 3 | 3 + 6 = 9 | too small → `left++` |
| 2 | 3 | 4 + 6 = 10 | found → `[3, 4]` |

Five steps instead of checking all 15 pairs.

**Run it:** `twoSumSorted` and `threeSum` in [`twopointer.go`](twopointer.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Two Sum II - Input Array Is Sorted | Easy | https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/ |
| 2 | 3Sum | Medium | https://leetcode.com/problems/3sum/ |

**What to notice — Problem 1:** The array is 1-indexed in the output. The key insight: on a sorted array, moving `left` increases the sum, moving `right` decreases it. You always know which direction to move.

**What to notice — Problem 2:** Sort first — this enables two pointer on the inner loop. The tricky part is skipping duplicates: skip `nums[i]` at the outer loop, and skip `nums[left]`/`nums[right]` after recording a valid triplet. Without duplicate skipping, you'll get repeated results.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Two Sum II | O(n) | O(1) |
| 3Sum | O(n²) — sort is O(n log n), then n fixed values × O(n) inner scan | O(1) extra (ignoring sort and output) |

**Common mistakes:**

- 3Sum: skipping duplicates with `nums[i] == nums[i+1]` instead of `nums[i] == nums[i-1]`. The first version skips the *first* copy, so `[-1, -1, 2]` is lost.
- 3Sum: skipping duplicates of `left`/`right` before recording the triplet instead of after.
- Two Sum II: returning 0-based indices — the problem wants 1-based.

---

## Day 3 — Opposite Ends: Min/Max Area

**Pattern:** Maximize or minimize a value computed from both ends  
**Key insight:** Always move the pointer with the smaller value — the area/water is limited by the shorter side, so moving the taller side inward can only decrease or maintain the limiting factor.

**Template adaptation:**

```go
func maxArea(height []int) int {
    left, right := 0, len(height)-1
    maxWater := 0

    for left < right {
        water := min(height[left], height[right]) * (right - left)
        if water > maxWater {
            maxWater = water
        }
        if height[left] < height[right] {
            left++  // move the shorter side — only way to potentially improve
        } else {
            right--
        }
    }
    return maxWater
}
```

`min` and `max` are Go builtins since Go 1.21 — no helper needed.

**Trapping Rain Water — track running max from both ends:**

```go
func trap(height []int) int {
    left, right := 0, len(height)-1
    maxLeft, maxRight := 0, 0
    water := 0

    for left < right {
        if height[left] < height[right] {
            if height[left] >= maxLeft {
                maxLeft = height[left]
            } else {
                water += maxLeft - height[left]
            }
            left++
        } else {
            if height[right] >= maxRight {
                maxRight = height[right]
            } else {
                water += maxRight - height[right]
            }
            right--
        }
    }
    return water
}
```

**Why it works (Container With Most Water):** Say `height[left] < height[right]`. Any container that keeps `left` and moves `right` inward is **narrower** and still capped at `height[left]` (the shorter wall can't get taller by moving the other side). So none of them can beat what you just measured — `left` is done, discard it. The pointer on the taller side is the only one with any upside.

**Dry run** — `[1, 8, 6, 2, 5, 4, 8, 3, 7]`:

| `left` | `right` | heights | area | best | Move |
|---|---|---|---|---|---|
| 0 | 8 | 1, 7 | min(1,7)×8 = 8 | 8 | `left++` |
| 1 | 8 | 8, 7 | min(8,7)×7 = 49 | 49 | `right--` |
| 1 | 7 | 8, 3 | min(8,3)×6 = 18 | 49 | `right--` |
| 1 | 6 | 8, 8 | min(8,8)×5 = 40 | 49 | `right--` |
| 1 | 5 | 8, 4 | min(8,4)×4 = 16 | 49 | `right--` |
| 1 | 4 | 8, 5 | min(8,5)×3 = 15 | 49 | `right--` |
| 1 | 3 | 8, 2 | min(8,2)×2 = 4 | 49 | `right--` |
| 1 | 2 | 8, 6 | min(8,6)×1 = 6 | 49 | `right--` |

**Why it works (Trapping Rain Water):** Water above bar `i` is `min(tallest bar to its left, tallest bar to its right) - height[i]`. Invariant: the taller of `height[left]` and `height[right]` is at least as tall as every bar already passed. So when `height[left] < height[right]`, the right side is guaranteed to have a wall ≥ `maxLeft`, the `min` is `maxLeft`, and the water at `left` is final — you never need the exact right max.

**Dry run** — `[4, 2, 0, 3, 2, 5]`:

| `left` | `right` | Side processed | `maxLeft` | `maxRight` | Water added | Total |
|---|---|---|---|---|---|---|
| 0 | 5 | left (h=4) | 4 | 0 | 0 | 0 |
| 1 | 5 | left (h=2) | 4 | 0 | 2 | 2 |
| 2 | 5 | left (h=0) | 4 | 0 | 4 | 6 |
| 3 | 5 | left (h=3) | 4 | 0 | 1 | 7 |
| 4 | 5 | left (h=2) | 4 | 0 | 2 | 9 |

Notice `maxRight` stays 0 the whole time — the right side is never processed. The `5` at index 5 is all the proof needed that every bar on the left is bounded by `maxLeft`.

**Run it:** `maxArea` and `trap` in [`twopointer.go`](twopointer.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Container With Most Water | Medium | https://leetcode.com/problems/container-with-most-water/ |
| 2 | Trapping Rain Water | Hard | https://leetcode.com/problems/trapping-rain-water/ |

**What to notice — Problem 1:** Area = `min(height[left], height[right]) * (right - left)`. Moving the taller pointer inward always reduces width without guarantee of a taller short side — so always move the shorter pointer.

**What to notice — Problem 2:** Water at position `i` = `min(maxLeft, maxRight) - height[i]`. You don't need to precompute the full left/right max arrays — track them with running variables as you move the pointers. Move the side with the shorter *bar*: if `height[left] < height[right]`, there is a wall at least as tall as `maxLeft` somewhere on the right, so the water at `left` is exactly `maxLeft - height[left]` — final, no need to know the true right max.

**Complexity:** both problems O(n) time, O(1) space. (Trapping Rain Water's prefix-max/suffix-max array version is also O(n) time but O(n) space — the two-pointer version is the upgrade.)

**Common mistakes:**

- Container: moving the **taller** side. That only ever shrinks the width without raising the cap.
- Trapping Rain Water: adding water at a bar that is a new max (it holds 0 water — it *becomes* the wall).
- Trapping Rain Water: computing `min(maxLeft, maxRight)` with `maxRight` from the running variable — it may be stale (still 0 above). The invariant is what makes the `min` unnecessary.

---

## Day 4 — Same Direction: Remove/Overwrite

**Pattern:** Fast pointer scans, slow pointer marks where to write  
**Key insight:** `slow` always points to the next valid write position. Only advance `slow` when you write. The elements before `slow` are always the valid result so far.

**Template adaptation:**

```go
// Remove duplicates: write when current != previous valid
func removeDuplicates(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    slow := 0

    for fast := 1; fast < len(nums); fast++ {
        if nums[fast] != nums[slow] { // found a new unique element
            slow++
            nums[slow] = nums[fast] // write it
        }
    }
    return slow + 1 // length of valid portion
}

// Move zeroes: write non-zero, then fill rest with zeros
func moveZeroes(nums []int) {
    slow := 0

    for fast := 0; fast < len(nums); fast++ {
        if nums[fast] != 0 {
            nums[slow] = nums[fast]
            slow++
        }
    }
    // fill remaining positions with zeros
    for slow < len(nums) {
        nums[slow] = 0
        slow++
    }
}
```

**Why it works — the invariant:** at every step, `nums[0..slow]` (inclusive) is exactly the de-duplicated version of `nums[0..fast-1]`. A new value at `fast` either repeats `nums[slow]` (sorted input, so duplicates are adjacent — skip it) or is new (extend the kept part by one). `fast` only reads positions ≥ `slow`, so writes never clobber data that hasn't been read yet.

**Dry run** — `removeDuplicates([0, 0, 1, 1, 1, 2, 3, 3])`:

| `fast` | `a[fast]` | `a[slow]` | Write? | `slow` after | Array (kept part in **bold**) |
|---|---|---|---|---|---|
| 1 | 0 | 0 | no | 0 | [**0**, 0, 1, 1, 1, 2, 3, 3] |
| 2 | 1 | 0 | yes | 1 | [**0**, **1**, 1, 1, 1, 2, 3, 3] |
| 3 | 1 | 1 | no | 1 | [**0**, **1**, 1, 1, 1, 2, 3, 3] |
| 4 | 1 | 1 | no | 1 | [**0**, **1**, 1, 1, 1, 2, 3, 3] |
| 5 | 2 | 1 | yes | 2 | [**0**, **1**, **2**, 1, 1, 2, 3, 3] |
| 6 | 3 | 2 | yes | 3 | [**0**, **1**, **2**, **3**, 1, 2, 3, 3] |
| 7 | 3 | 3 | no | 3 | [**0**, **1**, **2**, **3**, 1, 2, 3, 3] |

Returns `slow + 1 = 4`. Everything after the bold part is leftover junk — the problem only checks the first `k` elements.

**Run it:** `removeDuplicates` and `moveZeroes` in [`twopointer.go`](twopointer.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Remove Duplicates from Sorted Array | Easy | https://leetcode.com/problems/remove-duplicates-from-sorted-array/ |
| 2 | Move Zeroes | Easy | https://leetcode.com/problems/move-zeroes/ |

**What to notice — Problem 1:** `slow` starts at 0 (first element is always valid). `fast` starts at 1. Write when `nums[fast] != nums[slow]` — this preserves one copy of each duplicate. Return `slow + 1` (not `slow`) because `slow` is an index, not a count.

**What to notice — Problem 2:** Two-pass approach is clearest — first pass writes all non-zeros to the front, second pass fills the rest with zeros. One-pass swap also works: swap `nums[fast]` and `nums[slow]` instead of just writing, which naturally moves zeros to the end.

**Complexity:** both O(n) time, O(1) space.

**Common mistakes:**

- Returning `slow` instead of `slow + 1` in Remove Duplicates (here `slow` is the last kept *index*, not the count).
- Move Zeroes: swapping zeros with elements from the *right end* (opposite-ends style). It's also O(n), but it scrambles the relative order of the non-zeros, which the problem requires you to keep.
- Forgetting the empty-array guard in Remove Duplicates (`fast` starts at 1, but `slow + 1` would still return 1).

---

## Day 5 — Same Direction: Conditional Overwrite

**Pattern:** Same direction with a more complex write condition  
**Key insight:** The fast/slow template is identical to Day 4 — only the write condition changes. Recognizing this lets you adapt the template to any filtering problem.

**Template adaptation:**

```go
// Remove element: write when current != target value
func removeElement(nums []int, val int) int {
    slow := 0

    for fast := 0; fast < len(nums); fast++ {
        if nums[fast] != val { // keep everything that isn't val
            nums[slow] = nums[fast]
            slow++
        }
    }
    return slow
}

// Remove duplicates II: allow at most 2 of each — write when current != element 2 positions back
func removeDuplicatesII(nums []int) int {
    slow := 0

    for fast := 0; fast < len(nums); fast++ {
        if slow < 2 || nums[fast] != nums[slow-2] {
            // slow < 2: first two elements always valid
            // nums[fast] != nums[slow-2]: current differs from 2 positions back in result
            nums[slow] = nums[fast]
            slow++
        }
    }
    return slow
}
```

**Why `slow-2` works:** at most two copies are allowed, so a value is a forbidden third copy exactly when it equals the element **two places back in the kept part**. Because the input is sorted, if `nums[slow-2] == nums[fast]`, then `nums[slow-1]` is the same value too — two copies already kept. Comparing against `nums[slow-1]` would only tell you there's one copy, which is still allowed. Generalize: "at most `k` copies" → compare against `nums[slow-k]`.

**Dry run** — `removeDuplicatesII([1, 1, 1, 2, 2, 3])`:

| `fast` | `a[fast]` | `a[slow-2]` | Write? | `slow` after | Array (kept part in **bold**) |
|---|---|---|---|---|---|
| 0 | 1 | — (`slow < 2`) | yes | 1 | [**1**, 1, 1, 2, 2, 3] |
| 1 | 1 | — (`slow < 2`) | yes | 2 | [**1**, **1**, 1, 2, 2, 3] |
| 2 | 1 | 1 | no | 2 | [**1**, **1**, 1, 2, 2, 3] |
| 3 | 2 | 1 | yes | 3 | [**1**, **1**, **2**, 2, 2, 3] |
| 4 | 2 | 1 | yes | 4 | [**1**, **1**, **2**, **2**, 2, 3] |
| 5 | 3 | 2 | yes | 5 | [**1**, **1**, **2**, **2**, **3**, 3] |

Returns `5`.

**Run it:** `removeElement` and `removeDuplicatesII` in [`twopointer.go`](twopointer.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Remove Element | Easy | https://leetcode.com/problems/remove-element/ |
| 2 | Remove Duplicates from Sorted Array II | Medium | https://leetcode.com/problems/remove-duplicates-from-sorted-array-ii/ |

**What to notice — Problem 1:** Identical to Day 4's template. Only the condition changes from `nums[fast] != nums[slow]` to `nums[fast] != val`. Return `slow` (not `slow + 1`) because `slow` is already the count here — it starts at 0 and increments after each write.

**What to notice — Problem 2:** The condition `nums[fast] != nums[slow-2]` is the key insight. You're comparing against what's 2 positions back in the *result* (not in the original array). The `slow < 2` guard handles the first two elements which are always valid regardless of value.

**Complexity:** both O(n) time, O(1) space.

**Common mistakes:**

- Comparing against `nums[fast-2]` (the original array) instead of `nums[slow-2]` (the kept part). The original array has been partly overwritten, and even where it hasn't, it's the wrong thing to compare against.
- Remove Element: returning `slow + 1`. Here `slow` starts at 0 and is incremented *after* each write, so it's already the count. Compare with Day 4, where it's incremented *before* the write.

---

## Day 6 — Partition

**Pattern:** Divide array into sections using pointers as boundaries  
**Key insight:** Pointers mark the boundary between sections. Swap elements across the boundary to place them in the right section. With three sections, you need three pointers.

**Template adaptation:**

```go
// Dutch National Flag — three sections: low(0), mid(1), high(2)
func sortColors(nums []int) {
    low, mid, high := 0, 0, len(nums)-1

    for mid <= high {
        switch nums[mid] {
        case 0:
            nums[low], nums[mid] = nums[mid], nums[low]
            low++
            mid++
        case 1:
            mid++ // already in correct section
        case 2:
            nums[mid], nums[high] = nums[high], nums[mid]
            high-- // don't increment mid — swapped element unexamined
        }
    }
}

// Merge sorted arrays — merge from the END to avoid overwriting nums1
func merge(nums1 []int, m int, nums2 []int, n int) {
    i := m - 1     // last valid element in nums1
    j := n - 1     // last element in nums2
    k := m + n - 1 // last position in nums1

    for j >= 0 {
        if i >= 0 && nums1[i] > nums2[j] {
            nums1[k] = nums1[i]
            i--
        } else {
            nums1[k] = nums2[j]
            j--
        }
        k--
    }
}
```

**Why it works (Sort Colors):** the three pointers split the array into four zones that stay true after every step:

```
[ 0 0 0 | 1 1 1 | ? ? ? ? | 2 2 2 ]
  0..low-1  low..mid-1  mid..high  high+1..end
```

Every iteration takes one `?` from position `mid` and puts it into its zone, so the unknown region shrinks by one each time and the loop ends when it's empty. Swapping with `low` brings back a value from the `1` zone (already examined), so `mid` can advance. Swapping with `high` brings back a `?` — unexamined — so `mid` must stay put.

**Dry run** — `[2, 0, 2, 1, 1, 0]`:

| `low` | `mid` | `high` | `nums[mid]` | Action | Array after |
|---|---|---|---|---|---|
| 0 | 0 | 5 | 2 | swap with `high`, `high--` (mid stays) | `[0, 0, 2, 1, 1, 2]` |
| 0 | 0 | 4 | 0 | swap with `low`, `low++`, `mid++` | `[0, 0, 2, 1, 1, 2]` |
| 1 | 1 | 4 | 0 | swap with `low`, `low++`, `mid++` | `[0, 0, 2, 1, 1, 2]` |
| 2 | 2 | 4 | 2 | swap with `high`, `high--` (mid stays) | `[0, 0, 1, 1, 2, 2]` |
| 2 | 2 | 3 | 1 | `mid++` | `[0, 0, 1, 1, 2, 2]` |
| 2 | 3 | 3 | 1 | `mid++` | `[0, 0, 1, 1, 2, 2]` |

Row 1 is the important one: a `0` came back from `high`, and because `mid` didn't move, row 2 caught it.

**Run it:** `sortColors` and `merge` in [`twopointer.go`](twopointer.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Sort Colors | Medium | https://leetcode.com/problems/sort-colors/ |
| 2 | Merge Sorted Array | Easy | https://leetcode.com/problems/merge-sorted-array/ |

**What to notice — Problem 1:** When you swap `nums[mid]` with `nums[high]`, you don't increment `mid` — the swapped element from `high` is unexamined and might be 0 or 1. When you swap `nums[mid]` with `nums[low]`, you increment both — the swapped element from `low` is always 1 (everything before `low` is already 0, and `mid` only moves forward).

**What to notice — Problem 2:** Merging from the front causes overwriting — you'd clobber elements in `nums1` before reading them. Merging from the end is safe because the extra space in `nums1` is at the end. When `nums2` is exhausted (`j < 0`), the remaining `nums1` elements are already in place — no extra work needed.

**Complexity:**

| Problem | Time | Space |
|---|---|---|
| Sort Colors | O(n), single pass | O(1) |
| Merge Sorted Array | O(m + n) | O(1) — the spare room at the end of `nums1` is the buffer |

**Common mistakes:**

- Sort Colors: loop condition `mid < high` instead of `mid <= high` — the element *at* `high` is still unexamined.
- Merge: looping `for i >= 0 && j >= 0` and then forgetting to copy what's left of `nums2`. The template loops on `j >= 0` only, because leftover `nums1` elements are already in place.

---
## Day 7 — Mixed Review + Cheatsheet

**Pattern:** Consolidation under pressure  
**Key insight:** Pattern recognition is the skill. Identify which pointer setup a problem needs within 2 minutes of reading it.

**No new problems today.** Do this instead:

1. Write both templates (opposite ends and fast/slow) from memory.
2. Pick the 1 day that felt hardest (Days 2–6). Redo one problem without looking at your notes.
3. Answer out loud: *How do you decide between opposite ends and same direction?*

   > Use **opposite ends** when: the array is sorted and you're searching for a pair/triplet, or you're maximizing/minimizing something computed from both ends. Use **same direction** when: you're filtering or overwriting elements in-place, and the relative order of kept elements matters.

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the opposite ends template from memory
- [ ] Write the fast/slow template from memory
- [ ] Why do you sort first in 3Sum?
- [ ] In Remove Duplicates II, why compare against `nums[slow-2]` instead of `nums[slow-1]`?
- [ ] In Sort Colors, why don't you increment `mid` after swapping with `high`?
- [ ] In Merge Sorted Array, why merge from the end?

---

## Pattern Cheatsheet

| Pattern | Pointer setup | Array sorted? | Key move | Time / Space |
|---------|--------------|---------------|----------|--------------|
| Target sum | opposite ends | yes | `left++` if sum too small, `right--` if too large | O(n) / O(1) |
| Min/max area | opposite ends | no | move the pointer with the smaller value | O(n) / O(1) |
| Remove/overwrite | fast/slow same dir | yes (usually) | `slow` writes, `fast` scans; advance `slow` only on write | O(n) / O(1) |
| Conditional overwrite | fast/slow same dir | no | only the write condition changes | O(n) / O(1) |
| Partition (3 sections) | low/mid/high | no | swap to correct section; don't advance `mid` after swap with `high` | O(n) / O(1) |
| Merge sorted | three pointers from end | yes | merge backwards to avoid overwriting | O(m + n) / O(1) |

---

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Two Sum II - Input Array Is Sorted | Easy | https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/ |
| 2 | 3Sum | Medium | https://leetcode.com/problems/3sum/ |

**What to notice — Problem 1:** The array is 1-indexed in the output. The key insight: on a sorted array, moving `left` increases the sum, moving `right` decreases it. You always know which direction to move.

**What to notice — Problem 2:** Sort first — this enables two pointer on the inner loop. The tricky part is skipping duplicates: skip `nums[i]` at the outer loop, and skip `nums[left]`/`nums[right]` after recording a valid triplet. Without duplicate skipping, you'll get repeated results.

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

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Container With Most Water | Medium | https://leetcode.com/problems/container-with-most-water/ |
| 2 | Trapping Rain Water | Hard | https://leetcode.com/problems/trapping-rain-water/ |

**What to notice — Problem 1:** Area = `min(height[left], height[right]) * (right - left)`. Moving the taller pointer inward always reduces width without guarantee of a taller short side — so always move the shorter pointer.

**What to notice — Problem 2:** Water at position `i` = `min(maxLeft, maxRight) - height[i]`. You don't need to precompute the full left/right max arrays — track them with running variables as you move the pointers. Move whichever side has the smaller max (that side's water calculation is already finalized).

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Remove Duplicates from Sorted Array | Easy | https://leetcode.com/problems/remove-duplicates-from-sorted-array/ |
| 2 | Move Zeroes | Easy | https://leetcode.com/problems/move-zeroes/ |

**What to notice — Problem 1:** `slow` starts at 0 (first element is always valid). `fast` starts at 1. Write when `nums[fast] != nums[slow]` — this preserves one copy of each duplicate. Return `slow + 1` (not `slow`) because `slow` is an index, not a count.

**What to notice — Problem 2:** Two-pass approach is clearest — first pass writes all non-zeros to the front, second pass fills the rest with zeros. One-pass swap also works: swap `nums[fast]` and `nums[slow]` instead of just writing, which naturally moves zeros to the end.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Remove Element | Easy | https://leetcode.com/problems/remove-element/ |
| 2 | Remove Duplicates from Sorted Array II | Medium | https://leetcode.com/problems/remove-duplicates-from-sorted-array-ii/ |

**What to notice — Problem 1:** Identical to Day 4's template. Only the condition changes from `nums[fast] != nums[slow]` to `nums[fast] != val`. Return `slow` (not `slow + 1`) because `slow` is already the count here — it starts at 0 and increments after each write.

**What to notice — Problem 2:** The condition `nums[fast] != nums[slow-2]` is the key insight. You're comparing against what's 2 positions back in the *result* (not in the original array). The `slow < 2` guard handles the first two elements which are always valid regardless of value.

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

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Sort Colors | Medium | https://leetcode.com/problems/sort-colors/ |
| 2 | Merge Sorted Array | Easy | https://leetcode.com/problems/merge-sorted-array/ |

**What to notice — Problem 1:** When you swap `nums[mid]` with `nums[high]`, you don't increment `mid` — the swapped element from `high` is unexamined and might be 0 or 1. When you swap `nums[mid]` with `nums[low]`, you increment both — the swapped element from `low` is always 1 (everything before `low` is already 0, and `mid` only moves forward).

**What to notice — Problem 2:** Merging from the front causes overwriting — you'd clobber elements in `nums1` before reading them. Merging from the end is safe because the extra space in `nums1` is at the end. When `nums2` is exhausted (`j < 0`), the remaining `nums1` elements are already in place — no extra work needed.

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

| Pattern | Pointer setup | Array sorted? | Key move |
|---------|--------------|---------------|----------|
| Target sum | opposite ends | yes | `left++` if sum too small, `right--` if too large |
| Min/max area | opposite ends | no | move the pointer with the smaller value |
| Remove/overwrite | fast/slow same dir | yes (usually) | `slow` writes, `fast` scans; advance `slow` only on write |
| Conditional overwrite | fast/slow same dir | no | only the write condition changes |
| Partition (3 sections) | low/mid/high | no | swap to correct section; don't advance `mid` after swap with `high` |
| Merge sorted | three pointers from end | yes | merge backwards to avoid overwriting |

---

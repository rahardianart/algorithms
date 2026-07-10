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

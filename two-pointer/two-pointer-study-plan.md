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

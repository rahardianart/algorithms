# Linked List Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, DP, Heap, and Binary Search weeks)  
**Goal:** Write the reversal and fast/slow templates from memory, use a dummy head by default, never lose a node mid-rewire.

---

## The Core Templates

```go
type ListNode struct {
    Val  int
    Next *ListNode
}
```

### Dummy Head (any operation that might change the head)

```go
dummy := &ListNode{Next: head}
prev := dummy

for cur := head; cur != nil; cur = cur.Next {
    // decide whether to keep or skip cur, rewiring prev.Next as needed
    prev = cur
}

return dummy.Next
```

### Fast/Slow (Floyd's — middle-finding and cycle detection)

```go
slow, fast := head, head
for fast != nil && fast.Next != nil {
    slow = slow.Next
    fast = fast.Next.Next
}
// slow is now the middle node (or the meeting point, if a cycle exists)
```

### Reversal (in-place, iterative)

```go
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    cur := head

    for cur != nil {
        next := cur.Next // save before overwriting
        cur.Next = prev
        prev = cur
        cur = next
    }
    return prev
}
```

**Three rules — never break them:**

1. **Use a dummy head whenever the real head might change or be deleted.** It removes the need to special-case "is this the first node" — `prev` always starts pointing at something valid.
2. **Save `cur.Next` before you overwrite it.** The single most common linked-list bug is rewiring `cur.Next` and losing the reference to the rest of the list. Always capture `next := cur.Next` first.
3. **Fast/slow moves fast at 2x, slow at 1x.** If there's a cycle, they will meet — that's the whole proof, no need to re-derive it under pressure. If there's no cycle, `fast` (or `fast.Next`) hits `nil` first.

---

## Day 1 — Template + Concept

**Pattern:** All three templates  
**Key insight:** Every linked-list problem is dummy-head, fast/slow, or reversal — usually two of the three combined.

**No problems today.** Do this instead:

1. Read all three templates above slowly. Understand each line.
2. Close this file. Write all three from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why must `next := cur.Next` happen before `cur.Next = prev` in the reversal template?*

   > Once `cur.Next = prev` runs, the original link to the rest of the list is gone — if you hadn't already saved it in `next`, there'd be no way to advance `cur` forward; the list beyond the current node would be unreachable garbage.

---

## Day 2 — Reversal

**Pattern:** Reverse all or part of a list in place  
**Key insight:** Reversing a sublist is the same three-line loop as reversing the whole list — the only difference is where you start and where you reconnect the reversed piece back into the original list.

**Template adaptation:**

```go
func reverseBetween(head *ListNode, left, right int) *ListNode {
    dummy := &ListNode{Next: head}
    prev := dummy
    for i := 0; i < left-1; i++ {
        prev = prev.Next // walk to just before the sublist
    }

    cur := prev.Next
    for i := 0; i < right-left; i++ {
        next := cur.Next
        cur.Next = next.Next
        next.Next = prev.Next
        prev.Next = next // reinsert next right after prev, one node at a time
    }
    return dummy.Next
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Reverse Linked List | Easy | https://leetcode.com/problems/reverse-linked-list/ |
| 2 | Reverse Linked List II | Medium | https://leetcode.com/problems/reverse-linked-list-ii/ |

**What to notice — Problem 1:** This IS the core reversal template, verbatim. Get this one automatic before Day 2's second problem — everything else this week builds on it.

**What to notice — Problem 2:** `prev` marks the node just before the sublist (dummy head handles `left == 1`). Instead of the "save next, point back, advance" reversal pattern, this uses "head-insertion": repeatedly pull the node right after `cur` and reinsert it right after `prev`. Trace it on paper with a 5-node list before trusting it from memory.

---

## Day 3 — Fast/Slow: Middle and Cycle Detection

**Pattern:** Two pointers moving at different speeds through one list  
**Key insight:** Fast/slow answers two different questions with the identical loop: "where's the middle?" (stop when `fast` runs out) and "is there a cycle?" (stop when `fast` catches `slow`).

**Template adaptation:**

```go
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Middle of the Linked List | Easy | https://leetcode.com/problems/middle-of-the-linked-list/ |
| 2 | Linked List Cycle | Easy | https://leetcode.com/problems/linked-list-cycle/ |

**What to notice — Problem 1:** Pure fast/slow, no cycle check needed — when `fast` (or `fast.Next`) hits `nil`, `slow` is sitting on the middle. For even-length lists this lands on the *second* middle node — check which one the problem wants.

**What to notice — Problem 2:** Add one line to the fast/slow loop: check `slow == fast` after each step. No cycle → `fast` hits `nil` and the loop ends normally. Don't compare `.Val` — compare pointers directly, since values can repeat.

---

## Day 4 — Fast/Slow: Advanced

**Pattern:** Fast/slow combined with a second technique (math, or reversal)  
**Key insight:** Once you can find the middle and detect a cycle, two harder patterns fall out: finding *where* a cycle starts (a bit of pointer math after the meeting point), and palindrome-checking (find middle, reverse the second half, compare).

**Template adaptation:**

```go
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // meeting point found — reset one pointer to head, advance both 1 step at a time
            ptr := head
            for ptr != slow {
                ptr = ptr.Next
                slow = slow.Next
            }
            return ptr
        }
    }
    return nil
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Linked List Cycle II | Medium | https://leetcode.com/problems/linked-list-cycle-ii/ |
| 2 | Palindrome Linked List | Easy | https://leetcode.com/problems/palindrome-linked-list/ |

**What to notice — Problem 1:** The "reset to head, advance both by 1" step after the meeting point is a known result (Floyd's algorithm) — the distance from `head` to the cycle start equals the distance from the meeting point to the cycle start, going around the loop. Memorize that the second phase always moves both pointers one step at a time, not `slow`/`fast` speeds anymore.

**What to notice — Problem 2:** Three templates in one problem: fast/slow to find the middle, reversal (Day 2) on the second half, then walk both halves comparing values. No need for extra space — this is the in-place O(1) approach interviewers expect once they see you know the pieces.

---

## Day 5 — Dummy Head: Merge and Remove

**Pattern:** Operations where the head itself might be removed or replaced  
**Key insight:** The dummy head means `prev.Next = ...` works identically whether you're rewiring the 1st node or the 50th — no special case for "am I at the head."

**Template adaptation:**

```go
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy

    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            cur.Next = l1
            l1 = l1.Next
        } else {
            cur.Next = l2
            l2 = l2.Next
        }
        cur = cur.Next
    }
    // exactly one of l1, l2 is non-nil here — attach whichever remains
    if l1 != nil {
        cur.Next = l1
    } else {
        cur.Next = l2
    }
    return dummy.Next
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Merge Two Sorted Lists | Easy | https://leetcode.com/problems/merge-two-sorted-lists/ |
| 2 | Remove Nth Node From End of List | Medium | https://leetcode.com/problems/remove-nth-node-from-end-of-list/ |

**What to notice — Problem 1:** No dummy-head *deletion* here, but the dummy still avoids special-casing "which list's head comes first" — `dummy.Next` is correct regardless of which list wins the first comparison.

**What to notice — Problem 2:** Two pointers, `n+1` apart: advance `fast` n+1 steps ahead of `slow` first (starting both from `dummy`), then move both together until `fast` hits `nil` — `slow` now sits just before the node to remove. The dummy head is what makes removing the actual head node (n == length) not a special case.

---

## Day 6 — Multi-Step Combinations

**Pattern:** Problems that chain two or three earlier patterns together  
**Key insight:** By this point every "hard" linked-list problem is a composition of Days 2–5, not a new idea. Recognizing which pieces to chain, in what order, is the actual skill being tested.

**Template adaptation:**

```go
func reorderList(head *ListNode) {
    // 1. find the middle (fast/slow, Day 3)
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    // 2. reverse the second half (Day 2)
    var prev *ListNode
    cur := slow.Next
    slow.Next = nil // cut the list in two
    for cur != nil {
        next := cur.Next
        cur.Next = prev
        prev = cur
        cur = next
    }

    // 3. merge the two halves, alternating (dummy-head merge shape, Day 5)
    first, second := head, prev
    for second != nil {
        first.Next, first = second, first.Next
        second.Next, second = first, second.Next
    }
}
```

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Reorder List | Medium | https://leetcode.com/problems/reorder-list/ |
| 2 | Add Two Numbers | Medium | https://leetcode.com/problems/add-two-numbers/ |

**What to notice — Problem 1:** Three templates chained in sequence: find middle → reverse second half → merge-alternate. `slow.Next = nil` before reversing is easy to forget — without it the first half still points into the second half and you get a cycle.

**What to notice — Problem 2:** Not fast/slow or reversal — this is dummy-head plus manual carry arithmetic, digit by digit, like grade-school addition. Handle both lists being different lengths (treat a missing node as value 0) and a leftover carry after both lists end (append one more node).

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Almost every linked-list "hard" is dummy-head, fast/slow, and reversal composed together. If you can name which of the three you need and in what order, the code is mechanical.

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Answer out loud: *For a new problem, which of the three core templates does it need, and in what order?*

**Self-assessment — can you answer these without hesitation?**

- [ ] Write the reversal template from memory
- [ ] Write the fast/slow template from memory
- [ ] Why use a dummy head instead of special-casing the real head?
- [ ] Why does `slow.Next = nil` matter before reversing a second half?
- [ ] In Linked List Cycle II, why does resetting one pointer to `head` and stepping both by 1 find the cycle start?

---

## Pattern Cheatsheet

| Pattern | Core move | Needs dummy head? | Combines with |
|---------|-----------|---------------------|----------------|
| Reversal | save `next`, flip `cur.Next` to `prev` | only if reversing from the true head | fast/slow (Day 4, 6) |
| Fast/slow (middle) | `fast` 2x speed, stop when `fast`/`fast.Next` is nil | no | reversal, merge |
| Fast/slow (cycle) | stop when `slow == fast` | no | pointer-reset math (Day 4) |
| Dummy-head merge/remove | `prev.Next` rewiring, no head special-case | yes | — |
| Multi-step combination | chain the above in sequence | usually | all of the above |

---

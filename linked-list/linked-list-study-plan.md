# Linked List Interview Prep — 7-Day Study Plan

**Level:** Intermediate (completed BFS, DFS, Two Pointer, DP, Heap, and Binary Search weeks)  
**Goal:** Write the reversal and fast/slow templates from memory, use a dummy head by default, never lose a node mid-rewire.

---

## What Is a Linked List?

A **linked list** is a chain of nodes where each node holds a value and a pointer to the **next** node. Unlike an array, there's no index access — to reach the 5th node you walk from the head. In exchange, inserting or removing a node is O(1) once you're standing next to it: you just rewire pointers.

Linked list problems are about rewiring pointers **in place** without losing track of the rest of the list, usually with three tools: a **dummy head**, **fast/slow pointers**, and **in-place reversal**.

**Typical problems it solves:**

- Reversing all or part of a list
- Finding the middle, detecting a cycle, finding where the cycle starts
- Merging sorted lists
- Removing nodes (n-th from the end, all nodes with a value)
- Reordering and palindrome checks, adding numbers stored as lists

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
    if shouldRemove(cur) {
        prev.Next = cur.Next // unlink cur — prev stays where it is
    } else {
        prev = cur // cur stays in the list, so it becomes the new prev
    }
}

return dummy.Next
```

`prev` only advances onto nodes that stay in the list. If it advanced onto a removed node, the next unlink would rewire a node that's no longer reachable and the real list would keep the node you meant to delete.

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
3. **Fast/slow moves fast at 2x, slow at 1x.** If there's a cycle, they will meet: once both are inside the cycle, `fast` closes the gap by exactly one node per step, so the gap hits 0 — it can't jump over `slow`. If there's no cycle, `fast` (or `fast.Next`) hits `nil` first.

**Why linked lists feel hard:** the algorithms are short, but every line *destroys* information — `cur.Next = prev` erases the only path to the rest of the list. The fix is always the same: before you overwrite a pointer, ask "do I still need what it points to?" and save it first. Draw boxes and arrows on paper for 3–4 nodes and update the drawing line by line; it's the fastest way to debug any of these.

**Brute force → in place:** most of these problems are trivial if you copy the values into a slice (reverse it, index the middle, compare as a palindrome). That costs O(n) extra space. The templates below do the same work by rewiring pointers, in O(1) extra space — that's what interviewers are testing.

**Run it:** every function in this file is real code in [`linkedlist.go`](linkedlist.go), tested in [`linkedlist_test.go`](linkedlist_test.go) (with `fromSlice`/`toSlice` helpers for building lists). From the repo root: `go test ./linked-list/ -v`.

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

**Dry run** — `reverseList(1 → 2 → 3)`. `prev` is the reversed part built so far, `cur` is what's left:

| After iteration | `prev` (reversed part) | `cur` (remaining) |
|---|---|---|
| start | nil | `1 → 2 → 3` |
| 1 | `1` | `2 → 3` |
| 2 | `2 → 1` | `3` |
| 3 | `3 → 2 → 1` | nil |

Each iteration moves exactly one node from the front of "remaining" to the front of "reversed". When nothing remains, `prev` is the new head.

**Run it:** `removeIf` (dummy head), `middle` (fast/slow), and `reverseList` in [`linkedlist.go`](linkedlist.go). The `removeIf` test includes consecutive matches at the head — the case the dummy head exists for.

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

**Why head-insertion works:** `prev` never moves and `cur` never moves — `cur` is the original first node of the sublist, and it drifts toward the end as other nodes are pulled out from behind it. Each iteration takes the node right after `cur` and re-inserts it right after `prev`, i.e. at the front of the sublist. After `right - left` moves, the sublist is reversed and still attached at both ends, so no reconnect step is needed.

**Dry run** — `reverseBetween(1 → 2 → 3 → 4 → 5, 2, 4)`. `prev` = node `1`, `cur` = node `2`:

| Iteration | Node pulled (`next`) | List after |
|---|---|---|
| start | — | `1 → 2 → 3 → 4 → 5` |
| 1 | 3 | `1 → 3 → 2 → 4 → 5` |
| 2 | 4 | `1 → 4 → 3 → 2 → 5` |

Watch node `2` (`cur`): it moves one place right each iteration and ends as the sublist's tail, still pointing at `5`.

**Run it:** `reverseBetween` in [`linkedlist.go`](linkedlist.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Reverse Linked List | Easy | https://leetcode.com/problems/reverse-linked-list/ |
| 2 | Reverse Linked List II | Medium | https://leetcode.com/problems/reverse-linked-list-ii/ |

**What to notice — Problem 1:** This IS the core reversal template, verbatim. Get this one automatic before Day 2's second problem — everything else this week builds on it.

**What to notice — Problem 2:** `prev` marks the node just before the sublist (dummy head handles `left == 1`). Instead of the "save next, point back, advance" reversal pattern, this uses "head-insertion": repeatedly pull the node right after `cur` and reinsert it right after `prev`. Trace it on paper with a 5-node list before trusting it from memory.

**Complexity:** both O(n) time, O(1) space. (Recursive reversal is also O(n) time but O(n) stack space.)

**Common mistakes:**

- Reverse Linked List: returning `cur` (which is `nil` at the end) instead of `prev`.
- Reverse II: walking `left` steps instead of `left - 1` to find `prev` — you end up one node too far.
- Reverse II: not using a dummy head, then special-casing `left == 1`.

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

**Why `slow` lands in the middle:** `fast` moves twice as far as `slow`. When `fast` reaches the end (n steps), `slow` has gone n/2 steps — the middle.

**Why they must meet in a cycle:** once both are inside the cycle, look at the gap from `fast` forward to `slow`. Each step, `slow` moves +1 and `fast` +2, so the gap shrinks by exactly 1. A gap that shrinks by 1 every step can't skip over 0 — it hits 0, and they meet, within one lap.

**Dry run** — where `slow` and `fast` sit after each iteration:

| List | After iteration 1 | 2 | 3 | Loop stops because |
|---|---|---|---|---|
| `1 → 2 → 3 → 4 → 5` | slow 2, fast 3 | slow **3**, fast 5 | — | `fast.Next == nil` |
| `1 → 2 → 3 → 4 → 5 → 6` | slow 2, fast 3 | slow 3, fast 5 | slow **4**, fast nil | `fast == nil` |

Odd length → the exact middle. Even length → the **second** of the two middles.

**Run it:** `middle` and `hasCycle` in [`linkedlist.go`](linkedlist.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Middle of the Linked List | Easy | https://leetcode.com/problems/middle-of-the-linked-list/ |
| 2 | Linked List Cycle | Easy | https://leetcode.com/problems/linked-list-cycle/ |

**What to notice — Problem 1:** Pure fast/slow, no cycle check needed — when `fast` (or `fast.Next`) hits `nil`, `slow` is sitting on the middle. For even-length lists this lands on the *second* middle node — check which one the problem wants.

**What to notice — Problem 2:** Add one line to the fast/slow loop: check `slow == fast` after each step. No cycle → `fast` hits `nil` and the loop ends normally. Don't compare `.Val` — compare pointers directly, since values can repeat.

**Complexity:** both O(n) time, O(1) space. (A `map[*ListNode]bool` of seen nodes also detects cycles in O(n) time, but O(n) space.)

**Common mistakes:**

- Loop condition `fast.Next != nil && fast != nil` — wrong order; `fast.Next` panics when `fast` is nil. Check `fast != nil` first.
- Needing the **first** middle for even lengths (common in "split the list in half" problems): start `fast` at `head.Next` instead.
- Comparing `slow.Val == fast.Val` instead of `slow == fast`.

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

**Why the reset finds the cycle start:** let `a` = steps from head to the cycle start, `b` = steps from the cycle start to where they met, `c` = cycle length. At the meeting point `slow` has walked `a + b` and `fast` has walked `2(a + b)`. The extra distance `fast` covered, `a + b`, must be whole laps of the cycle: `a + b = k·c`. Rearranged: **`a = k·c − b`**. So walking `a` steps from the meeting point = finishing the current lap (`c − b` steps) plus some whole laps → you land exactly on the cycle start. A pointer starting at `head` also needs exactly `a` steps to get there, so they meet at the cycle start.

**Dry run** — `1 → 2 → 3 → 4 → 5 → 6 → 7 → 8`, with `8` pointing back to `3` (so `a = 2`, `c = 6`):

Phase 1 (find the meeting point):

| Step | `slow` | `fast` |
|---|---|---|
| 1 | 2 | 3 |
| 2 | 3 | 5 |
| 3 | 4 | 7 |
| 4 | 5 | 3 |
| 5 | 6 | 5 |
| 6 | **7** | **7** |

They meet at `7`, so `b = 4` (3 → 7). Check: `a + b = 6 = 1 · c` ✓, and `a = c − b = 2`.

Phase 2 (both at speed 1):

| Step | `ptr` (from head) | `slow` (from meeting point) |
|---|---|---|
| 0 | 1 | 7 |
| 1 | 2 | 8 |
| 2 | **3** | **3** ← cycle start |

**Run it:** `detectCycle` in [`linkedlist.go`](linkedlist.go) — `TestCycle` includes this exact list.

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Linked List Cycle II | Medium | https://leetcode.com/problems/linked-list-cycle-ii/ |
| 2 | Palindrome Linked List | Easy | https://leetcode.com/problems/palindrome-linked-list/ |

**What to notice — Problem 1:** The "reset to head, advance both by 1" step after the meeting point is a known result (Floyd's algorithm) — the distance from `head` to the cycle start equals the distance from the meeting point forward to the cycle start, plus some whole number of laps. Why: let `a` = head→cycle start, `b` = cycle start→meeting point, `c` = cycle length. Slow walked `a+b`, fast walked `2(a+b)`; the extra `a+b` is whole laps, `k·c`, so `a = k·c − b` — walking `a` steps from the meeting point finishes the lap and lands on the cycle start. Remember that the second phase always moves both pointers one step at a time, not `slow`/`fast` speeds anymore.

**What to notice — Problem 2:** Three templates in one problem: fast/slow to find the middle, reversal (Day 2) on the second half, then walk both halves comparing values. No need for extra space — this is the in-place O(1) approach interviewers expect once they see you know the pieces.

**Complexity:** both O(n) time, O(1) space.

**Common mistakes:**

- Cycle II: in phase 2, still moving one pointer at double speed.
- Cycle II: starting phase 2 with `ptr := head.Next` (off by one).
- Palindrome: not restoring the second half after comparing. Interviewers sometimes ask you to leave the input unchanged — reverse it back.
- Palindrome on odd lengths: the middle node belongs to neither half. Comparing until the shorter half runs out handles it.

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

**Why the dummy works for merge:** `cur` always points at the tail of the merged list. Without a dummy, the tail doesn't exist until you've picked the first node, so the first iteration needs its own code. The dummy is a fake tail that exists from the start; you throw it away at the end with `dummy.Next`.

**Dry run** — merging `1 → 2 → 4` and `1 → 3 → 4`:

| Compare | Take | Merged so far (after dummy) |
|---|---|---|
| 1 vs 1 | `l1`'s 1 (`<=` prefers `l1`) | `1` |
| 2 vs 1 | `l2`'s 1 | `1 → 1` |
| 2 vs 3 | `l1`'s 2 | `1 → 1 → 2` |
| 4 vs 3 | `l2`'s 3 | `1 → 1 → 2 → 3` |
| 4 vs 4 | `l1`'s 4 | `1 → 1 → 2 → 3 → 4` |
| `l1` empty | attach rest of `l2` | `1 → 1 → 2 → 3 → 4 → 4` |

The last row is one pointer assignment, not a loop — the rest of `l2` is already linked.

**Run it:** `mergeTwoLists` in [`linkedlist.go`](linkedlist.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Merge Two Sorted Lists | Easy | https://leetcode.com/problems/merge-two-sorted-lists/ |
| 2 | Remove Nth Node From End of List | Medium | https://leetcode.com/problems/remove-nth-node-from-end-of-list/ |

**What to notice — Problem 1:** No dummy-head *deletion* here, but the dummy still avoids special-casing "which list's head comes first" — `dummy.Next` is correct regardless of which list wins the first comparison.

**What to notice — Problem 2:** Two pointers, `n+1` apart: advance `fast` n+1 steps ahead of `slow` first (starting both from `dummy`), then move both together until `fast` hits `nil` — `slow` now sits just before the node to remove. The dummy head is what makes removing the actual head node (n == length) not a special case.

**Complexity:** Merge O(m + n) time; Remove Nth O(n) time, single pass. Both O(1) space.

**Common mistakes:**

- Remove Nth: starting `fast` and `slow` at `head` instead of `dummy` — then removing the head itself (n == length) has nowhere for `slow` to stand.
- Remove Nth: advancing `fast` `n` steps instead of `n + 1` — `slow` ends *on* the node to delete instead of just before it.
- Merge: allocating new nodes. Reuse the existing ones; only pointers change.

---

## Day 6 — Multi-Step Combinations

**Pattern:** Problems that chain two or three earlier patterns together  
**Key insight:** By this point every "hard" linked-list problem is a composition of Days 2–5, not a new idea. Recognizing which pieces to chain, in what order, is the actual skill being tested.

**Template adaptation:**

```go
func reorderList(head *ListNode) {
    if head == nil {
        return
    }

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
        firstNext, secondNext := first.Next, second.Next // save both before rewiring
        first.Next = second
        second.Next = firstNext
        first, second = firstNext, secondNext
    }
}
```

**Dry run** — `reorderList(1 → 2 → 3 → 4 → 5)`, stage by stage:

| Stage | Result |
|---|---|
| 1. find middle (fast/slow) | `slow` = node `3` |
| 2. cut after middle, reverse second half | first: `1 → 2 → 3`, second: `5 → 4` |
| 3a. weave: 1, then 5 | `1 → 5 → 2 → 3` |
| 3b. weave: 2, then 4 | `1 → 5 → 2 → 4 → 3` |

The weave stops when `second` runs out. The first half is never shorter than the second, so its leftover node (`3`) is already in place at the end.

**Run it:** `reorderList` in [`linkedlist.go`](linkedlist.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Reorder List | Medium | https://leetcode.com/problems/reorder-list/ |
| 2 | Add Two Numbers | Medium | https://leetcode.com/problems/add-two-numbers/ |

**What to notice — Problem 1:** Three templates chained in sequence: find middle → reverse second half → merge-alternate. `slow.Next = nil` before reversing is easy to forget — without it the first half still points into the second half and you get a cycle.

**What to notice — Problem 2:** Not fast/slow or reversal — this is dummy-head plus manual carry arithmetic, digit by digit, like grade-school addition. Handle both lists being different lengths (treat a missing node as value 0) and a leftover carry after both lists end (append one more node).

**Complexity:** Reorder List O(n) time, O(1) space. Add Two Numbers O(max(m, n)) time, O(1) extra space besides the output list.

**Common mistakes:**

- Reorder: forgetting `slow.Next = nil` (creates a cycle — see above).
- Reorder: weaving without saving both `Next` pointers first. The template saves `firstNext` and `secondNext` before touching anything.
- Add Two Numbers: dropping the final carry. `5 + 5` must produce `0 → 1`, not `0`.

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

Every pattern here is O(n) time and O(1) extra space — that's the point of doing it in place.

| Pattern | Core move | Needs dummy head? | Combines with |
|---------|-----------|---------------------|----------------|
| Reversal | save `next`, flip `cur.Next` to `prev` | only if reversing from the true head | fast/slow (Day 4, 6) |
| Fast/slow (middle) | `fast` 2x speed, stop when `fast`/`fast.Next` is nil | no | reversal, merge |
| Fast/slow (cycle) | stop when `slow == fast` | no | pointer-reset math (Day 4) |
| Dummy-head merge/remove | `prev.Next` rewiring, no head special-case | yes | — |
| Multi-step combination | chain the above in sequence | usually | all of the above |

---

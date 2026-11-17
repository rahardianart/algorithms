// Package linkedlist holds the runnable templates from linked-list-study-plan.md.
package linkedlist

type ListNode struct {
	Val  int
	Next *ListNode
}

// removeIf is the dummy-head template with "shouldRemove" as the keep/skip decision.
func removeIf(head *ListNode, shouldRemove func(*ListNode) bool) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy

	for cur := head; cur != nil; cur = cur.Next {
		if shouldRemove(cur) {
			prev.Next = cur.Next
		} else {
			prev = cur
		}
	}

	return dummy.Next
}

// middle is the fast/slow template; for even lengths it returns the second middle.
func middle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head

	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}

func reverseBetween(head *ListNode, left, right int) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy
	for i := 0; i < left-1; i++ {
		prev = prev.Next
	}

	cur := prev.Next
	for i := 0; i < right-left; i++ {
		next := cur.Next
		cur.Next = next.Next
		next.Next = prev.Next
		prev.Next = next
	}
	return dummy.Next
}

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

func detectCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
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
	if l1 != nil {
		cur.Next = l1
	} else {
		cur.Next = l2
	}
	return dummy.Next
}

func reorderList(head *ListNode) {
	if head == nil {
		return
	}

	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	var prev *ListNode
	cur := slow.Next
	slow.Next = nil
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	first, second := head, prev
	for second != nil {
		firstNext, secondNext := first.Next, second.Next
		first.Next = second
		second.Next = firstNext
		first, second = firstNext, secondNext
	}
}

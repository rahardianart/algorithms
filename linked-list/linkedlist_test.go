package linkedlist

import (
	"slices"
	"testing"
)

func fromSlice(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func toSlice(head *ListNode) []int {
	var out []int
	for ; head != nil; head = head.Next {
		out = append(out, head.Val)
	}
	return out
}

// withCycle builds a list whose tail points back to the node at index pos (pos < 0 means no cycle).
func withCycle(vals []int, pos int) (head, cycleStart *ListNode) {
	head = fromSlice(vals)
	if pos < 0 {
		return head, nil
	}
	var tail *ListNode
	i := 0
	for n := head; n != nil; n = n.Next {
		if i == pos {
			cycleStart = n
		}
		tail = n
		i++
	}
	tail.Next = cycleStart
	return head, cycleStart
}

func TestRemoveIf(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		val  int
		want []int
	}{
		{"middle", []int{1, 2, 6, 3, 4, 5, 6}, 6, []int{1, 2, 3, 4, 5}},
		{"head and consecutive", []int{7, 7, 7, 1, 7, 7}, 7, []int{1}},
		{"everything", []int{7, 7}, 7, nil},
		{"empty", nil, 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeIf(fromSlice(tt.in), func(n *ListNode) bool { return n.Val == tt.val })
			if !slices.Equal(toSlice(got), tt.want) {
				t.Errorf("got %v, want %v", toSlice(got), tt.want)
			}
		})
	}
}

func TestMiddle(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3, 4, 5}, 3},
		{[]int{1, 2, 3, 4, 5, 6}, 4},
		{[]int{1}, 1},
	}
	for _, tt := range tests {
		if got := middle(fromSlice(tt.in)).Val; got != tt.want {
			t.Errorf("middle(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestReverseList(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{nil, nil},
	}
	for _, tt := range tests {
		if got := toSlice(reverseList(fromSlice(tt.in))); !slices.Equal(got, tt.want) {
			t.Errorf("reverseList(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestReverseBetween(t *testing.T) {
	tests := []struct {
		in          []int
		left, right int
		want        []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, 4, []int{1, 4, 3, 2, 5}},
		{[]int{1, 2, 3, 4, 5}, 1, 5, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2, 3}, 1, 2, []int{2, 1, 3}},
		{[]int{5}, 1, 1, []int{5}},
	}
	for _, tt := range tests {
		got := toSlice(reverseBetween(fromSlice(tt.in), tt.left, tt.right))
		if !slices.Equal(got, tt.want) {
			t.Errorf("reverseBetween(%v, %d, %d) = %v, want %v", tt.in, tt.left, tt.right, got, tt.want)
		}
	}
}

func TestCycle(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		pos  int
	}{
		{"tail to index 1", []int{3, 2, 0, -4}, 1},
		{"tail to head", []int{1, 2}, 0},
		{"self loop", []int{1}, 0},
		{"long tail, short cycle", []int{1, 2, 3, 4, 5, 6, 7, 8}, 6},
		{"no cycle", []int{1, 2, 3}, -1},
		{"empty", nil, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, start := withCycle(tt.vals, tt.pos)
			if got := hasCycle(head); got != (start != nil) {
				t.Errorf("hasCycle = %v, want %v", got, start != nil)
			}
			if got := detectCycle(head); got != start {
				t.Errorf("detectCycle returned the wrong node")
			}
		})
	}
}

func TestMergeTwoLists(t *testing.T) {
	tests := []struct{ a, b, want []int }{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{nil, nil, nil},
		{nil, []int{0}, []int{0}},
		{[]int{5, 6}, []int{1}, []int{1, 5, 6}},
	}
	for _, tt := range tests {
		got := toSlice(mergeTwoLists(fromSlice(tt.a), fromSlice(tt.b)))
		if !slices.Equal(got, tt.want) {
			t.Errorf("mergeTwoLists(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestReorderList(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		head := fromSlice(tt.in)
		reorderList(head)
		if got := toSlice(head); !slices.Equal(got, tt.want) {
			t.Errorf("reorderList(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

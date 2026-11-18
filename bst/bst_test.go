package bst

import (
	"slices"
	"testing"
)

const math32Min = -1 << 31

func build(vals ...int) *TreeNode {
	var root *TreeNode
	for _, v := range vals {
		root = insert(root, v)
	}
	return root
}

func height(n *TreeNode) int {
	if n == nil {
		return 0
	}
	return 1 + max(height(n.Left), height(n.Right))
}

func TestInsertSearchInorder(t *testing.T) {
	root := build(8, 3, 10, 1, 6, 14, 4, 7, 13)
	if got, want := inorderIterative(root), []int{1, 3, 4, 6, 7, 8, 10, 13, 14}; !slices.Equal(got, want) {
		t.Errorf("inorder = %v, want %v", got, want)
	}
	if n := search(root, 7); n == nil || n.Val != 7 {
		t.Errorf("search(7) = %v, want node 7", n)
	}
	if n := search(root, 5); n != nil {
		t.Errorf("search(5) = %v, want nil", n)
	}
	if got := inorderIterative(nil); got != nil {
		t.Errorf("inorder(nil) = %v, want nil", got)
	}
}

func TestIsValidBST(t *testing.T) {
	//   5
	//  / \
	// 1   6
	//    / \
	//   3   7     <- 3 < 6 is fine locally, but it sits in 5's right subtree, so it must be > 5
	trap := &TreeNode{Val: 5,
		Left:  &TreeNode{Val: 1},
		Right: &TreeNode{Val: 6, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 7}},
	}
	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{"built by insert", build(8, 3, 10, 1, 6, 14), true},
		{"grandchild breaks the root bound", trap, false},
		{"duplicate not allowed", &TreeNode{Val: 2, Left: &TreeNode{Val: 2}}, false},
		{"empty", nil, true},
		{"single at int32 min", &TreeNode{Val: math32Min}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidBST(tt.root); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeSum(t *testing.T) {
	root := build(10, 5, 15, 3, 7, 18)
	tests := []struct{ lo, hi, want int }{
		{7, 15, 32},
		{6, 10, 17},
		{0, 100, 58},
		{11, 14, 0},
	}
	for _, tt := range tests {
		if got := rangeSum(root, tt.lo, tt.hi); got != tt.want {
			t.Errorf("rangeSum(%d, %d) = %d, want %d", tt.lo, tt.hi, got, tt.want)
		}
	}
}

func TestSortedArrayToBST(t *testing.T) {
	nums := []int{-10, -3, 0, 5, 9, 12, 20}
	root := sortedArrayToBST(nums)
	if !isValidBST(root) {
		t.Fatal("result is not a valid BST")
	}
	if got := inorderIterative(root); !slices.Equal(got, nums) {
		t.Errorf("inorder = %v, want %v", got, nums)
	}
	if h := height(root); h != 3 {
		t.Errorf("height = %d, want 3 for 7 nodes", h)
	}
	if sortedArrayToBST(nil) != nil {
		t.Error("empty input should give nil")
	}
}

func TestFindSwapped(t *testing.T) {
	tests := []struct {
		name          string
		vals          []int // a valid BST is built from vals, then a and b are swapped
		a, b          int
		first, second int
	}{
		{"adjacent in inorder", []int{4, 2, 6, 1, 3, 5, 7}, 3, 4, 4, 3},
		{"far apart", []int{4, 2, 6, 1, 3, 5, 7}, 1, 7, 7, 1},
		{"root and leaf", []int{4, 2, 6, 1, 3, 5, 7}, 4, 5, 5, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := build(tt.vals...)
			search(root, tt.a).Val, search(root, tt.b).Val = tt.b, tt.a
			first, second := findSwapped(root)
			if first != tt.first || second != tt.second {
				t.Errorf("got (%d, %d), want (%d, %d)", first, second, tt.first, tt.second)
			}
		})
	}
}

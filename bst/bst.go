// Package bst holds the runnable templates from bst-study-plan.md.
package bst

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// search walks one root-to-leaf path: O(h).
func search(root *TreeNode, val int) *TreeNode {
	node := root
	for node != nil && node.Val != val {
		if val < node.Val {
			node = node.Left
		} else {
			node = node.Right
		}
	}
	return node
}

// insert adds val as a new leaf (duplicates go right) and returns the root.
func insert(root *TreeNode, val int) *TreeNode {
	newNode := &TreeNode{Val: val}
	if root == nil {
		return newNode
	}
	node := root
	for {
		if val < node.Val {
			if node.Left == nil {
				node.Left = newNode
				return root
			}
			node = node.Left
		} else {
			if node.Right == nil {
				node.Right = newNode
				return root
			}
			node = node.Right
		}
	}
}

// isValidBST requires every value to be strictly between all its ancestors' bounds.
func isValidBST(root *TreeNode) bool {
	var valid func(node *TreeNode, lo, hi int) bool
	valid = func(node *TreeNode, lo, hi int) bool {
		if node == nil {
			return true
		}
		if node.Val <= lo || node.Val >= hi {
			return false
		}
		return valid(node.Left, lo, node.Val) && valid(node.Right, node.Val, hi)
	}
	return valid(root, math.MinInt, math.MaxInt)
}

// inorderIterative returns the values in sorted order using an explicit stack.
func inorderIterative(root *TreeNode) []int {
	var out []int
	var stack []*TreeNode
	node := root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.Left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, node.Val)
		node = node.Right
	}
	return out
}

// rangeSum adds every value in [lo, hi], skipping subtrees the BST order rules out.
func rangeSum(node *TreeNode, lo, hi int) int {
	if node == nil {
		return 0
	}
	if node.Val < lo {
		return rangeSum(node.Right, lo, hi)
	}
	if node.Val > hi {
		return rangeSum(node.Left, lo, hi)
	}
	return node.Val + rangeSum(node.Left, lo, hi) + rangeSum(node.Right, lo, hi)
}

// sortedArrayToBST builds a height-balanced BST from a sorted slice.
func sortedArrayToBST(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	mid := len(nums) / 2
	return &TreeNode{
		Val:   nums[mid],
		Left:  sortedArrayToBST(nums[:mid]),
		Right: sortedArrayToBST(nums[mid+1:]),
	}
}

// findSwapped returns the two values that were swapped in an otherwise valid BST.
func findSwapped(root *TreeNode) (first, second int) {
	var prev *TreeNode
	found := false

	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)
		if prev != nil && prev.Val > node.Val {
			if !found {
				first = prev.Val
				found = true
			}
			second = node.Val
		}
		prev = node
		inorder(node.Right)
	}
	inorder(root)
	return first, second
}

package dfs

import (
	"reflect"
	"slices"
	"testing"
)

//	    5
//	   / \
//	  4   8
//	 /   / \
//	11  13  4
//
// /  \      \
// 7   2      1
func sampleTree() *TreeNode {
	return &TreeNode{Val: 5,
		Left: &TreeNode{Val: 4,
			Left: &TreeNode{Val: 11, Left: &TreeNode{Val: 7}, Right: &TreeNode{Val: 2}},
		},
		Right: &TreeNode{Val: 8,
			Left:  &TreeNode{Val: 13},
			Right: &TreeNode{Val: 4, Right: &TreeNode{Val: 1}},
		},
	}
}

var graph = [][]int{
	0: {1, 2},
	1: {0, 3},
	2: {0, 3},
	3: {1, 2},
	4: {5},
	5: {4},
	6: {},
}

func TestDFSGeneric(t *testing.T) {
	visited := map[int]bool{}
	dfs(0, func(n int) []int { return graph[n] }, visited)
	want := map[int]bool{0: true, 1: true, 2: true, 3: true}
	if !reflect.DeepEqual(visited, want) {
		t.Errorf("visited = %v, want %v", visited, want)
	}
}

func TestTraversalOrder(t *testing.T) {
	var pre, post []int
	preorder(sampleTree(), &pre)
	postorder(sampleTree(), &post)

	if want := []int{5, 4, 11, 7, 2, 8, 13, 4, 1}; !slices.Equal(pre, want) {
		t.Errorf("preorder = %v, want %v", pre, want)
	}
	if want := []int{7, 2, 11, 4, 13, 1, 4, 8, 5}; !slices.Equal(post, want) {
		t.Errorf("postorder = %v, want %v", post, want)
	}
	if got := dfsIterative(sampleTree()); !slices.Equal(got, pre) {
		t.Errorf("dfsIterative = %v, want recursive preorder %v", got, pre)
	}
	if got := dfsIterative(nil); got != nil {
		t.Errorf("dfsIterative(nil) = %v, want nil", got)
	}
}

func TestMaxDepth(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want int
	}{
		{"empty", nil, 0},
		{"single", &TreeNode{Val: 1}, 1},
		{"sample", sampleTree(), 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxDepth(tt.root); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestHasPathSum(t *testing.T) {
	tests := []struct {
		name   string
		root   *TreeNode
		target int
		want   bool
	}{
		{"5-4-11-2", sampleTree(), 22, true},
		{"5-8-4-1", sampleTree(), 18, true},
		{"not a root-to-leaf sum", sampleTree(), 9, false}, // 5+4 stops at a non-leaf
		{"empty tree", nil, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPathSum(tt.root, tt.target); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectPaths(t *testing.T) {
	var got [][]int
	collectPaths(sampleTree(), nil, &got)
	want := [][]int{{5, 4, 11, 7}, {5, 4, 11, 2}, {5, 8, 13}, {5, 8, 4, 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCountComponents(t *testing.T) {
	if got := countComponents(len(graph), graph); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestFill(t *testing.T) {
	grid := [][]int{
		{1, 1, 0},
		{1, 0, 1},
		{0, 1, 1},
	}
	fill(grid, 0, 0, 1, 2)
	want := [][]int{
		{2, 2, 0},
		{2, 0, 1}, // diagonal cells are not neighbors, so this region is untouched
		{0, 1, 1},
	}
	if !reflect.DeepEqual(grid, want) {
		t.Errorf("got %v, want %v", grid, want)
	}
}

func TestCombine(t *testing.T) {
	tests := []struct {
		name       string
		candidates []int
		k          int
		want       [][]int
	}{
		{"choose 2 of 4", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}},
		{"choose all", []int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{"choose none", []int{1, 2}, 0, [][]int{{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combine(tt.candidates, tt.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

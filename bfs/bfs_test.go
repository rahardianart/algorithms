package bfs

import (
	"slices"
	"testing"
)

// 0 — 1
// |   |
// 2 — 3 — 4      5 — 6
var sample = [][]int{
	0: {1, 2},
	1: {0, 3},
	2: {0, 3},
	3: {1, 2, 4},
	4: {3},
	5: {6},
	6: {5},
}

func TestBFS(t *testing.T) {
	got := bfs(0, func(n int) []int { return sample[n] })
	want := []int{0, 1, 2, 3, 4}
	if !slices.Equal(got, want) {
		t.Errorf("bfs(0) = %v, want %v", got, want)
	}
}

func TestLevelCounts(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want []int
	}{
		{"empty", nil, nil},
		{"single", &TreeNode{Val: 1}, []int{1}},
		{
			"three levels",
			&TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
				Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 6}},
			},
			[]int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levelCounts(tt.root); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBFSGraph(t *testing.T) {
	tests := []struct {
		start int
		want  []int
	}{
		{0, []int{0, 1, 2, 3, 4}},
		{4, []int{4, 3, 1, 2, 0}},
		{5, []int{5, 6}},
	}
	for _, tt := range tests {
		if got := bfsGraph(sample, tt.start); !slices.Equal(got, tt.want) {
			t.Errorf("bfsGraph(%d) = %v, want %v", tt.start, got, tt.want)
		}
	}
}

func TestCountComponents(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		graph [][]int
		want  int
	}{
		{"sample", 7, sample, 2},
		{"no edges", 3, [][]int{{}, {}, {}}, 3},
		{"one node", 1, [][]int{{}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countComponents(tt.n, tt.graph); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestShortestPath(t *testing.T) {
	tests := []struct {
		start, end, want int
	}{
		{0, 4, 3},
		{0, 0, 0},
		{2, 1, 2},
		{0, 5, -1},
	}
	for _, tt := range tests {
		if got := shortestPath(sample, tt.start, tt.end); got != tt.want {
			t.Errorf("shortestPath(%d, %d) = %d, want %d", tt.start, tt.end, got, tt.want)
		}
	}
}

func TestMultiSourceBFS(t *testing.T) {
	grid3 := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	tests := []struct {
		name    string
		grid    [][]int
		sources []Cell
		want    int
	}{
		{"one corner", grid3, []Cell{{0, 0}}, 4},
		{"opposite corners", grid3, []Cell{{0, 0}, {2, 2}}, 2},
		{"center", grid3, []Cell{{1, 1}}, 2},
		{"single cell", [][]int{{0}}, []Cell{{0, 0}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multiSourceBFS(tt.grid, tt.sources); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// flipOneBit generates every string that differs from s in exactly one bit.
func flipOneBit(s string) []string {
	var out []string
	for i := range s {
		b := []byte(s)
		if b[i] == '0' {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
		out = append(out, string(b))
	}
	return out
}

func TestMinSteps(t *testing.T) {
	tests := []struct {
		name       string
		start, end string
		deadends   []string
		want       int
	}{
		{"no deadends", "000", "111", nil, 3},
		{"start is end", "000", "000", nil, 0},
		{"detour still 3", "000", "111", []string{"100", "010"}, 3},
		{"walled off", "000", "111", []string{"100", "010", "001"}, -1},
		{"start is a deadend", "000", "111", []string{"000"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dead := map[string]bool{}
			for _, d := range tt.deadends {
				dead[d] = true
			}
			if got := minSteps(tt.start, tt.end, dead, flipOneBit); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

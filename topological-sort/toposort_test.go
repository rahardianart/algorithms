package toposort

import (
	"reflect"
	"slices"
	"testing"
)

// respects reports whether order puts u before v for every edge.
func respects(order []int, edges [][2]int) bool {
	pos := make(map[int]int, len(order))
	for i, v := range order {
		pos[v] = i
	}
	for _, e := range edges {
		if pos[e[0]] > pos[e[1]] {
			return false
		}
	}
	return true
}

func TestTopoOrder(t *testing.T) {
	edges := [][2]int{{5, 2}, {5, 0}, {4, 0}, {4, 1}, {2, 3}, {3, 1}}
	got := topoOrder(6, edges)
	if want := []int{4, 5, 2, 0, 3, 1}; !slices.Equal(got, want) {
		t.Errorf("topoOrder = %v, want %v", got, want)
	}
	if !respects(got, edges) {
		t.Errorf("order %v violates an edge", got)
	}
	if got := topoOrder(3, [][2]int{{0, 1}, {1, 2}, {2, 0}}); got != nil {
		t.Errorf("cycle should give nil, got %v", got)
	}
	if got := topoOrder(3, nil); !slices.Equal(got, []int{0, 1, 2}) {
		t.Errorf("no edges = %v, want [0 1 2]", got)
	}
}

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		want bool
	}{
		{"chain", [][]int{{1}, {2}, {}}, false},
		{"cycle", [][]int{{1}, {2}, {0}, {1}}, true},
		{"self loop", [][]int{{0}}, true},
		{"diamond is not a cycle", [][]int{{1, 2}, {3}, {3}, {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasCycle(len(tt.adj), tt.adj); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlienOrder(t *testing.T) {
	tests := []struct {
		words []string
		want  string
	}{
		{[]string{"wrt", "wrf", "er", "ett", "rftt"}, "wertf"},
		{[]string{"z", "x"}, "zx"},
		{[]string{"z", "x", "z"}, ""},
		{[]string{"abc", "ab"}, ""},
		{[]string{"ab", "abc"}, "abc"},
	}
	for _, tt := range tests {
		if got := alienOrder(tt.words); got != tt.want {
			t.Errorf("alienOrder(%v) = %q, want %q", tt.words, got, tt.want)
		}
	}
}

func TestLongestPath(t *testing.T) {
	edges := [][3]int{{0, 1, 3}, {0, 2, 2}, {1, 3, 4}, {2, 3, 1}, {3, 4, 2}, {2, 4, 9}}
	if got := longestPath(5, edges); got != 11 {
		t.Errorf("longestPath = %d, want 11", got)
	}
	if got := longestPath(3, nil); got != 0 {
		t.Errorf("no edges = %d, want 0", got)
	}
}

func TestPeelLeaves(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][2]int
		want  [][]int
	}{
		{"two centers", 6, [][2]int{{3, 0}, {3, 1}, {3, 2}, {3, 4}, {5, 4}}, [][]int{{0, 1, 2, 5}, {3, 4}}},
		{"one center", 4, [][2]int{{1, 0}, {1, 2}, {1, 3}}, [][]int{{0, 2, 3}, {1}}},
		{"path of 5", 5, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}, [][]int{{0, 4}, {1, 3}, {2}}},
		{"single node", 1, nil, [][]int{{0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := peelLeaves(tt.n, tt.edges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

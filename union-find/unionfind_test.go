package unionfind

import (
	"reflect"
	"testing"
)

func TestDSU(t *testing.T) {
	d := newDSU(6)
	steps := []struct {
		a, b   int
		merged bool
		count  int
	}{
		{0, 1, true, 5},
		{1, 2, true, 4},
		{3, 4, true, 3},
		{2, 0, false, 3}, // already connected: this edge closes a cycle
		{4, 2, true, 2},
	}
	for _, s := range steps {
		if got := d.union(s.a, s.b); got != s.merged {
			t.Fatalf("union(%d, %d) = %v, want %v", s.a, s.b, got, s.merged)
		}
		if d.count != s.count {
			t.Fatalf("after union(%d, %d), count = %d, want %d", s.a, s.b, d.count, s.count)
		}
	}
	if d.find(0) != d.find(3) {
		t.Error("0 and 3 should be connected")
	}
	if d.find(5) == d.find(0) {
		t.Error("5 should be alone")
	}
}

func TestCountComponents(t *testing.T) {
	tests := []struct {
		n     int
		edges [][2]int
		want  int
	}{
		{5, [][2]int{{0, 1}, {1, 2}, {3, 4}}, 2},
		{5, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}, 1},
		{3, nil, 3},
	}
	for _, tt := range tests {
		if got := countComponents(tt.n, tt.edges); got != tt.want {
			t.Errorf("countComponents(%d, %v) = %d, want %d", tt.n, tt.edges, got, tt.want)
		}
	}
}

func TestIsTree(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][2]int
		want  bool
	}{
		{"tree", 5, [][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 4}}, true},
		{"cycle", 5, [][2]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}, false},
		{"n-1 edges but a cycle and a loner", 4, [][2]int{{0, 1}, {1, 2}, {2, 0}}, false},
		{"single node", 1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTree(tt.n, tt.edges); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupEquivalent(t *testing.T) {
	got := groupEquivalent([][2]string{{"a", "b"}, {"c", "d"}, {"b", "e"}, {"x", "x"}})
	want := [][]string{{"a", "b", "e"}, {"c", "d"}, {"x"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestKruskal(t *testing.T) {
	edges := [][3]int{
		{0, 1, 4}, {0, 2, 1}, {1, 2, 2}, {1, 3, 5}, {2, 3, 8}, {3, 4, 3},
	}
	if got := kruskal(5, edges); got != 11 {
		t.Errorf("kruskal = %d, want 11", got)
	}
	if got := kruskal(3, [][3]int{{0, 1, 1}}); got != -1 {
		t.Errorf("disconnected kruskal = %d, want -1", got)
	}
}

func TestPotentialDSU(t *testing.T) {
	p := newPotentialDSU(4) // a=0, b=1, c=2, d=3
	if !p.union(0, 1, 3) || !p.union(1, 2, 2) {
		t.Fatal("consistent facts were rejected")
	}
	if got, ok := p.diff(0, 2); !ok || got != 5 {
		t.Errorf("a - c = %d, %v; want 5, true", got, ok)
	}
	if got, ok := p.diff(2, 0); !ok || got != -5 {
		t.Errorf("c - a = %d, %v; want -5, true", got, ok)
	}
	if _, ok := p.diff(0, 3); ok {
		t.Error("a and d are not connected")
	}
	if p.union(2, 0, 1) {
		t.Error("c - a = 1 contradicts c - a = -5 but was accepted")
	}
	if !p.union(2, 0, -5) {
		t.Error("c - a = -5 is consistent but was rejected")
	}
}

package dijkstra

import (
	"reflect"
	"slices"
	"testing"
)

func adjacency(n int, edges [][3]int) [][]edge {
	adj := make([][]edge, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], edge{e[1], e[2]})
	}
	return adj
}

var sample = [][3]int{{0, 1, 4}, {0, 2, 1}, {2, 1, 2}, {1, 3, 1}, {2, 3, 5}, {3, 4, 3}}

func TestDijkstra(t *testing.T) {
	got := dijkstra(6, adjacency(6, sample), 0)
	want := []int{0, 3, 1, 4, 7, inf}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMinPathCostGrid(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"must go around", [][]int{
			{1, 9, 1, 1, 1},
			{1, 9, 1, 9, 1},
			{1, 1, 1, 9, 1},
		}, 11},
		{"right/down only would be worse", [][]int{
			{1, 1, 1},
			{9, 9, 1},
			{1, 1, 1},
			{1, 9, 9},
			{1, 1, 1},
		}, 11},
		{"single cell", [][]int{{7}}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minPathCostGrid(tt.grid); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestZeroOneBFS(t *testing.T) {
	edges := [][3]int{{0, 1, 1}, {0, 2, 0}, {2, 3, 1}, {1, 3, 0}, {3, 4, 0}, {2, 1, 0}}
	adj := adjacency(5, edges)
	got := zeroOneBFS(5, adj, 0)
	if want := dijkstra(5, adj, 0); !slices.Equal(got, want) {
		t.Errorf("zeroOneBFS = %v, dijkstra = %v", got, want)
	}
	if want := []int{0, 0, 0, 0, 0}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCheapestWithOneFree(t *testing.T) {
	edges := [][3]int{{0, 1, 10}, {1, 3, 10}, {0, 2, 3}, {2, 3, 15}, {3, 4, 1}}
	adj := adjacency(5, edges)
	tests := []struct {
		src, dst, want int
	}{
		{0, 3, 3},  // 0→2 (3) then 2→3 free
		{0, 4, 4},  // 0→2 (3), 2→3 free, 3→4 (1)
		{0, 1, 0},  // the single edge is free
		{3, 0, -1}, // unreachable
		{2, 2, 0},
	}
	for _, tt := range tests {
		if got := cheapestWithOneFree(5, adj, tt.src, tt.dst); got != tt.want {
			t.Errorf("cheapestWithOneFree(%d, %d) = %d, want %d", tt.src, tt.dst, got, tt.want)
		}
	}
}

func TestBellmanFord(t *testing.T) {
	// LeetCode 787 example: 0→1 (100), 1→2 (100), 0→2 (500), 1→3 (600), 2→3 (200)
	edges := [][3]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}, {1, 3, 600}, {2, 3, 200}}
	tests := []struct {
		maxEdges int
		want     int // distance to node 3
	}{
		{1, inf},
		{2, 700},
		{3, 400},
	}
	for _, tt := range tests {
		if got := bellmanFord(4, edges, 0, tt.maxEdges)[3]; got != tt.want {
			t.Errorf("maxEdges %d: dist[3] = %d, want %d", tt.maxEdges, got, tt.want)
		}
	}

	negative := [][3]int{{0, 1, 4}, {0, 2, 5}, {2, 1, -3}}
	if got := bellmanFord(3, negative, 0, 2); !slices.Equal(got, []int{0, 2, 5}) {
		t.Errorf("negative edge: got %v, want [0 2 5]", got)
	}
}

func TestFloydWarshall(t *testing.T) {
	edges := [][3]int{{0, 1, 3}, {1, 2, 1}, {1, 3, 4}, {2, 3, 1}}
	got := floydWarshall(4, edges)
	want := [][]int{
		{0, 3, 4, 5},
		{3, 0, 1, 2},
		{4, 1, 0, 1},
		{5, 2, 1, 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

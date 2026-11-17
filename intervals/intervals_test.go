package intervals

import (
	"reflect"
	"testing"
)

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{"unsorted with overlaps", [][]int{{1, 3}, {8, 10}, {2, 6}, {15, 18}, {17, 20}}, [][]int{{1, 6}, {8, 10}, {15, 20}}},
		{"touching endpoints merge", [][]int{{1, 4}, {4, 5}}, [][]int{{1, 5}}},
		{"contained", [][]int{{1, 10}, {2, 3}, {4, 5}}, [][]int{{1, 10}}},
		{"single", [][]int{{1, 1}}, [][]int{{1, 1}}},
		{"empty", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeIntervals(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntervalIntersection(t *testing.T) {
	a := [][]int{{0, 2}, {5, 10}, {13, 23}, {24, 25}}
	b := [][]int{{1, 5}, {8, 12}, {15, 24}, {25, 26}}
	want := [][]int{{1, 2}, {5, 5}, {8, 10}, {15, 23}, {24, 24}, {25, 25}}
	if got := intervalIntersection(a, b); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got := intervalIntersection([][]int{{1, 3}}, nil); got != nil {
		t.Errorf("with empty list got %v, want nil", got)
	}
}

func TestMaxNonOverlapping(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want int
	}{
		{"long interval first", [][]int{{1, 10}, {2, 3}, {4, 5}, {6, 7}}, 3},
		{"touching is fine", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 3},
		{"all identical", [][]int{{1, 2}, {1, 2}, {1, 2}}, 1},
		{"empty", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxNonOverlapping(tt.in); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMinIntervalsToCover(t *testing.T) {
	tests := []struct {
		name       string
		in         [][]int
		start, end int
		want       int
	}{
		{"video stitching example", [][]int{{0, 2}, {4, 6}, {8, 10}, {1, 9}, {1, 5}, {5, 9}}, 0, 10, 3},
		{"gap", [][]int{{0, 1}, {2, 3}}, 0, 3, -1},
		{"one covers all", [][]int{{0, 1}, {0, 5}, {1, 3}}, 0, 5, 1},
		{"already covered", nil, 4, 4, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minIntervalsToCover(tt.in, tt.start, tt.end); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMaxConcurrent(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want int
	}{
		{"meetings", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
		{"back to back", [][]int{{1, 5}, {5, 10}, {10, 15}}, 1},
		{"nested", [][]int{{1, 10}, {2, 9}, {3, 8}, {4, 7}}, 4},
		{"mixed", [][]int{{1, 4}, {2, 6}, {4, 7}, {5, 8}}, 3},
		{"empty", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxConcurrent(tt.in); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

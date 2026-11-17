package greedy

import (
	"slices"
	"testing"
)

func TestMaxSubArray(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{[]int{1}, 1},
		{[]int{5, 4, -1, 7, 8}, 23},
		{[]int{-3, -1, -2}, -1},
	}
	for _, tt := range tests {
		if got := maxSubArray(tt.nums); got != tt.want {
			t.Errorf("maxSubArray(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestCanReach(t *testing.T) {
	tests := []struct {
		nums []int
		want bool
	}{
		{[]int{2, 3, 1, 1, 4}, true},
		{[]int{3, 2, 1, 0, 4}, false},
		{[]int{0}, true},
		{[]int{0, 1}, false},
		{[]int{1, 0, 1, 0}, false},
	}
	for _, tt := range tests {
		if got := canReach(tt.nums); got != tt.want {
			t.Errorf("canReach(%v) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

func TestMaxSatisfied(t *testing.T) {
	tests := []struct {
		needs, items []int
		want         int
	}{
		{[]int{1, 2, 3}, []int{1, 1}, 1},
		{[]int{1, 2}, []int{1, 2, 3}, 2},
		{[]int{5, 1, 4, 2}, []int{3, 6, 3, 1}, 3},
		{[]int{10}, []int{1, 2, 3}, 0},
	}
	for _, tt := range tests {
		got := maxSatisfied(slices.Clone(tt.needs), slices.Clone(tt.items))
		if got != tt.want {
			t.Errorf("maxSatisfied(%v, %v) = %d, want %d", tt.needs, tt.items, got, tt.want)
		}
	}
}

func TestPartitionLabels(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"ababcbacadefegdehijhklij", []int{9, 7, 8}},
		{"abacbdeedf", []int{5, 4, 1}},
		{"eccbbbbdec", []int{10}},
		{"abc", []int{1, 1, 1}},
	}
	for _, tt := range tests {
		if got := partitionLabels(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("partitionLabels(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestMinRefuelStops(t *testing.T) {
	tests := []struct {
		name              string
		target, startFuel int
		stations          [][]int
		want              int
	}{
		{"no stops needed", 1, 1, nil, 0},
		{"unreachable", 100, 1, [][]int{{10, 100}}, -1},
		{"two stops", 100, 10, [][]int{{10, 60}, {20, 30}, {30, 30}, {60, 40}}, 2},
		{"take the big one late", 100, 50, [][]int{{25, 25}, {50, 50}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minRefuelStops(tt.target, tt.startFuel, tt.stations); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

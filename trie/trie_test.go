package trie

import (
	"reflect"
	"testing"
)

func TestTrieBasics(t *testing.T) {
	tr := newTrie()
	for _, w := range []string{"car", "cat", "cart"} {
		tr.Insert(w)
	}
	tests := []struct {
		query              string
		search, startsWith bool
	}{
		{"car", true, true},
		{"cart", true, true},
		{"ca", false, true},
		{"cab", false, false},
		{"carts", false, false},
		{"", false, true},
	}
	for _, tt := range tests {
		if got := tr.Search(tt.query); got != tt.search {
			t.Errorf("Search(%q) = %v, want %v", tt.query, got, tt.search)
		}
		if got := tr.StartsWith(tt.query); got != tt.startsWith {
			t.Errorf("StartsWith(%q) = %v, want %v", tt.query, got, tt.startsWith)
		}
	}
}

func TestCountWithPrefix(t *testing.T) {
	tr := newTrie()
	for _, w := range []string{"apple", "app", "apply", "bat"} {
		tr.Insert(w)
	}
	tests := []struct {
		prefix string
		want   int
	}{
		{"app", 3},
		{"appl", 2},
		{"apple", 1},
		{"b", 1},
		{"c", 0},
		{"", 4},
	}
	for _, tt := range tests {
		if got := tr.CountWithPrefix(tt.prefix); got != tt.want {
			t.Errorf("CountWithPrefix(%q) = %d, want %d", tt.prefix, got, tt.want)
		}
	}
}

func TestSearchPattern(t *testing.T) {
	tr := newTrie()
	for _, w := range []string{"bad", "dad", "mad", "cat"} {
		tr.Insert(w)
	}
	tests := []struct {
		pattern string
		want    bool
	}{
		{"pad", false},
		{"bad", true},
		{".ad", true},
		{"b..", true},
		{"c.t", true},
		{"...", true},
		{"..", false},
		{"....", false},
	}
	for _, tt := range tests {
		if got := tr.SearchPattern(tt.pattern); got != tt.want {
			t.Errorf("SearchPattern(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}

func TestFindWordsInText(t *testing.T) {
	got := findWordsInText("catsanddog", []string{"cat", "cats", "and", "sand", "dog"})
	want := [][2]int{{0, 2}, {0, 3}, {3, 6}, {4, 6}, {7, 9}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMaxXor(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{3, 10, 5, 25, 2, 8}, 28},
		{[]int{14, 70, 53, 83, 49, 91, 36, 80, 92, 51, 66, 70}, 127},
		{[]int{0}, 0},
	}
	for _, tt := range tests {
		b := newBinaryTrie()
		for _, x := range tt.nums {
			b.insert(x)
		}
		best := 0
		for _, x := range tt.nums {
			best = max(best, b.maxXorWith(x))
		}
		if best != tt.want {
			t.Errorf("max xor of %v = %d, want %d", tt.nums, best, tt.want)
		}
	}
}

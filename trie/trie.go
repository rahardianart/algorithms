// Package trie holds the runnable templates from trie-study-plan.md.
//
// Words are lowercase a–z.
package trie

type trieNode struct {
	children [26]*trieNode
	isEnd    bool
	count    int // how many inserted words pass through this node
}

type Trie struct {
	root *trieNode
}

func newTrie() *Trie {
	return &Trie{root: &trieNode{}}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		if node.children[c] == nil {
			node.children[c] = &trieNode{}
		}
		node = node.children[c]
		node.count++
	}
	node.isEnd = true
}

// walk follows prefix from the root and returns the node it ends on, or nil.
func (t *Trie) walk(prefix string) *trieNode {
	node := t.root
	for i := 0; i < len(prefix); i++ {
		node = node.children[prefix[i]-'a']
		if node == nil {
			return nil
		}
	}
	return node
}

func (t *Trie) Search(word string) bool {
	node := t.walk(word)
	return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

// CountWithPrefix returns how many inserted words start with prefix.
func (t *Trie) CountWithPrefix(prefix string) int {
	if prefix == "" {
		total := 0
		for _, child := range t.root.children {
			if child != nil {
				total += child.count
			}
		}
		return total
	}
	node := t.walk(prefix)
	if node == nil {
		return 0
	}
	return node.count
}

// SearchPattern is Search where '.' matches any single letter.
func (t *Trie) SearchPattern(pattern string) bool {
	var match func(node *trieNode, i int) bool
	match = func(node *trieNode, i int) bool {
		if i == len(pattern) {
			return node.isEnd
		}
		if pattern[i] == '.' {
			for _, child := range node.children {
				if child != nil && match(child, i+1) {
					return true
				}
			}
			return false
		}
		child := node.children[pattern[i]-'a']
		return child != nil && match(child, i+1)
	}
	return match(t.root, 0)
}

// findWordsInText returns every [start, end] (inclusive) where a dictionary word occurs in text.
func findWordsInText(text string, words []string) [][2]int {
	t := newTrie()
	for _, w := range words {
		t.Insert(w)
	}

	var found [][2]int
	for start := 0; start < len(text); start++ {
		node := t.root
		for end := start; end < len(text); end++ {
			node = node.children[text[end]-'a']
			if node == nil {
				break
			}
			if node.isEnd {
				found = append(found, [2]int{start, end})
			}
		}
	}
	return found
}

const maxBit = 30 // enough for values below 2^31

type bitNode struct {
	children [2]*bitNode
}

type binaryTrie struct {
	root *bitNode
}

func newBinaryTrie() *binaryTrie {
	return &binaryTrie{root: &bitNode{}}
}

func (b *binaryTrie) insert(x int) {
	node := b.root
	for bit := maxBit; bit >= 0; bit-- {
		d := (x >> bit) & 1
		if node.children[d] == nil {
			node.children[d] = &bitNode{}
		}
		node = node.children[d]
	}
}

// maxXorWith returns the largest x XOR y over all inserted y. The trie must not be empty.
func (b *binaryTrie) maxXorWith(x int) int {
	node := b.root
	result := 0
	for bit := maxBit; bit >= 0; bit-- {
		d := (x >> bit) & 1
		if want := 1 - d; node.children[want] != nil {
			result |= 1 << bit
			node = node.children[want]
		} else {
			node = node.children[d]
		}
	}
	return result
}

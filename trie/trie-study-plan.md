# Trie Interview Prep — 7-Day Study Plan

**Level:** Intermediate (week 15 — uses DFS week's recursion and backtracking)  
**Goal:** Implement a trie from memory with a fixed-size child array, use it for prefix queries, wildcard search, and multi-word matching, and recognize the binary-trie trick for XOR problems.

---

## The Core Template

```go
type trieNode struct {
    children [26]*trieNode // one slot per letter; nil = no word continues this way
    isEnd    bool          // a whole word ends exactly here
}

type Trie struct {
    root *trieNode
}

func (t *Trie) Insert(word string) {
    node := t.root
    for i := 0; i < len(word); i++ {
        c := word[i] - 'a'
        if node.children[c] == nil {
            node.children[c] = &trieNode{}
        }
        node = node.children[c]
    }
    node.isEnd = true
}

// walk follows prefix from the root; nil means no word starts with it
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

func (t *Trie) Search(word string) bool      { n := t.walk(word); return n != nil && n.isEnd }
func (t *Trie) StartsWith(prefix string) bool { return t.walk(prefix) != nil }
```

**Three rules — never break them:**

1. **A node is a prefix, not a letter.** The path from the root to a node spells a prefix; the node itself just says "this prefix exists" and (via `isEnd`) "it's also a whole word". Thinking of nodes as letters leads to off-by-one bugs at the root.
2. **`isEnd` is what separates "word" from "prefix".** Insert `"cart"` and `"car"` is a prefix but not a word — until you insert `"car"` too. `Search` checks `isEnd`; `StartsWith` doesn't.
3. **Share the walk.** `Search`, `StartsWith`, counting, and autocomplete all start with "follow the prefix" — write `walk` once and build on it.

**Brute force → trie:** checking which of `m` words start with a prefix of length `L` is O(m · L) by comparing each word. A trie answers it in O(L) — independent of how many words there are — because words sharing a prefix share the path. The price is memory: up to 26 pointers per node.

**Complexity:** Insert/Search/StartsWith O(L) for a word of length L. Space O(total characters × 26) pointers in the worst case.

**Run it:** every template in this file is a real function in [`trie.go`](trie.go), tested in [`trie_test.go`](trie_test.go). From the repo root: `go test ./trie/ -v`.

---

## Day 1 — Template + Concept

**Pattern:** Core template  
**Key insight:** A trie is a tree where every edge is a letter and every path from the root is a prefix. Words that share a beginning share nodes.

**No problems today.** Do this instead:

1. Read the template above slowly. Understand each line.
2. Close this file. Write the `trieNode`, `Insert`, `walk`, `Search`, and `StartsWith` from memory in Go.
3. Open this file. Compare. Fix any differences.
4. Repeat steps 2–3 two more times.
5. Answer this in your head: *Why `[26]*trieNode` instead of `map[byte]*trieNode`?*

   > With only lowercase letters, a fixed array gives O(1) child lookup with no hashing, no allocation per lookup, and predictable memory. A map is better when the alphabet is large or sparse (Unicode, arbitrary bytes) — then 26 (or 256) mostly-nil pointers per node would waste a lot of memory.

**Dry run** — insert `"car"`, `"cat"`, `"cart"` (`*` = `isEnd`):

```
root
 └─ c
     └─ a
         ├─ r*
         │   └─ t*
         └─ t*
```

| Query | `walk` ends at | `Search` | `StartsWith` |
|---|---|---|---|
| `"car"` | `r*` | true | true |
| `"ca"` | `a` (not an end) | false | true |
| `"cab"` | nil (no `b` under `a`) | false | false |
| `"carts"` | nil (no `s` under `t`) | false | false |

Seven characters were inserted, but only five nodes exist below the root — `c` and `a` are shared by all three words.

---

## Day 2 — Implement and Use the Basic Trie

**Pattern:** Core template as a data structure  
**Key insight:** Most trie problems are "build the trie once, then answer many prefix questions" — the build cost is paid once.

**Template adaptation:** the core template above, wrapped with a constructor:

```go
func newTrie() *Trie {
    return &Trie{root: &trieNode{}}
}
```

**Run it:** `newTrie`, `Insert`, `Search`, and `StartsWith` in [`trie.go`](trie.go) (dry run on Day 1).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Implement Trie (Prefix Tree) | Medium | https://leetcode.com/problems/implement-trie-prefix-tree/ |
| 2 | Longest Common Prefix | Easy | https://leetcode.com/problems/longest-common-prefix/ |

**What to notice — Problem 1:** This IS the template. LeetCode's constructor is called `Constructor` and returns a value, not a pointer — adjust the signature, not the logic.

**What to notice — Problem 2:** Insert every word, then walk down from the root while the current node has **exactly one** child and is **not** the end of a word. (Vertical scanning without a trie is just as fast here — the trie version is practice for seeing "shared prefix = shared path".)

**Complexity:** Problem 1 O(L) per operation. Problem 2 O(total characters).

**Common mistakes:**

- Forgetting `isEnd` and treating every reachable node as a word.
- Problem 2: not stopping at a node where a word ends — `["ab", "abc"]` has common prefix `"ab"`, and the walk must stop there.
- Allocating a new root per call instead of once in the constructor.

---

## Day 3 — Wildcard Search (DFS through the Trie)

**Pattern:** Follow the exact letter when you can; branch into every child when you see a wildcard  
**Key insight:** A wildcard turns the single path walk into a DFS — but only *at wildcard positions*. Everywhere else it's still one pointer hop.

**Template adaptation:**

```go
// Search where '.' matches any single letter
func (t *Trie) SearchPattern(pattern string) bool {
    var match func(node *trieNode, i int) bool
    match = func(node *trieNode, i int) bool {
        if i == len(pattern) {
            return node.isEnd
        }
        if pattern[i] == '.' {
            for _, child := range node.children {
                if child != nil && match(child, i+1) {
                    return true // stop at the first branch that works
                }
            }
            return false
        }
        child := node.children[pattern[i]-'a']
        return child != nil && match(child, i+1)
    }
    return match(t.root, 0)
}
```

**Why it works:** `match(node, i)` answers "does some word continue from `node` that matches `pattern[i:]`?" A letter has exactly one place to continue; a `.` can continue in any existing child, so try each until one succeeds.

**Dry run** — trie with `bad`, `dad`, `mad`, `cat`; pattern `".ad"`:

| Position | Char | At node | Tries | Result |
|---|---|---|---|---|
| 0 | `.` | root | children `b`, `c`, `d`, `m` in order | — |
| 1 | `a` | `b` | `a` exists | — |
| 2 | `d` | `ba` | `d` exists | — |
| end | — | `bad` | `isEnd` = true | **true** (stops; `c`, `d`, `m` never tried) |

For `"c.t"`, only one branch exists at each fixed letter, and the `.` at position 1 tries the only child of `c`, which is `a`.

**Run it:** `SearchPattern` in [`trie.go`](trie.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Design Add and Search Words Data Structure | Medium | https://leetcode.com/problems/design-add-and-search-words-data-structure/ |
| 2 | Replace Words | Medium | https://leetcode.com/problems/replace-words/ |

**What to notice — Problem 1:** This IS the template.

**What to notice — Problem 2:** Insert every root, then for each word in the sentence walk the trie letter by letter and stop at the **first** node with `isEnd` — that's the shortest root. If the walk falls off the trie first, keep the word unchanged.

**Complexity:** Problem 1: O(L) without wildcards; with them, up to O(26ᵈ · L) for `d` wildcards in the worst case. Problem 2 O(total characters).

**Common mistakes:**

- Returning `true` when the pattern is exhausted instead of `node.isEnd`.
- Collecting results from all branches instead of short-circuiting on the first `true`.
- Replace Words: picking the longest matching root instead of the shortest.

---

## Day 4 — Prefix Counts and Autocomplete

**Pattern:** Store extra information on each node during insert  
**Key insight:** Anything you want to know "for all words with this prefix" can be precomputed into the prefix's node: a count, a sum, the top suggestions.

**Template adaptation:**

```go
type trieNode struct {
    children [26]*trieNode
    isEnd    bool
    count    int // how many inserted words pass through this node
}

// in Insert, after moving to the child:
node = node.children[c]
node.count++

// how many inserted words start with prefix
func (t *Trie) CountWithPrefix(prefix string) int {
    node := t.walk(prefix)
    if node == nil {
        return 0
    }
    return node.count
}
```

**Why `count` is right:** every word that starts with prefix `P` passes through `P`'s node exactly once during its insert, and no other word does. So incrementing on the way down counts exactly the words in that subtree.

**Dry run** — insert `apple`, `app`, `apply`, `bat` (`count` in parentheses):

```
root
 ├─ a (3)
 │   └─ p (3)
 │       └─ p* (3)
 │           └─ l (2)
 │               ├─ e* (1)
 │               └─ y* (1)
 └─ b (1)
     └─ a (1)
         └─ t* (1)
```

`CountWithPrefix("app") = 3`, `("appl") = 2`, `("c") = 0` — each a single walk, no subtree traversal.

**Run it:** `CountWithPrefix` in [`trie.go`](trie.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Search Suggestions System | Medium | https://leetcode.com/problems/search-suggestions-system/ |
| 2 | Map Sum Pairs | Medium | https://leetcode.com/problems/map-sum-pairs/ |

**What to notice — Problem 1:** Sort the products first, then store up to 3 product strings on each node during insert (they arrive in lexicographic order, so the first 3 are the answer). Each keystroke is then one hop plus reading the node's list.

**What to notice — Problem 2:** Like `count`, but store a running **sum** of values. The twist: inserting an existing key *replaces* its value, so add the **difference** (`new - old`, using a side map of key → old value) along the path.

**Complexity:** Problem 1 O(total characters · log) for the sort, O(L) per query. Problem 2 O(L) per operation.

**Common mistakes:**

- Problem 1: running a DFS on every keystroke to collect suggestions — correct, but much slower than storing them.
- Problem 2: adding the full new value on a repeated key — double-counts.

---

## Day 5 — Matching Many Words at Once

**Pattern:** From every starting position in a text, walk the trie as far as the text allows  
**Key insight:** Checking every dictionary word at every position costs O(n · m · L). Walking the trie from each position checks *all* words in one walk, and stops the moment no word can continue.

**Template adaptation:**

```go
// every [start, end] (inclusive) where a dictionary word occurs in text
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
                break // no dictionary word continues this way
            }
            if node.isEnd {
                found = append(found, [2]int{start, end})
            }
        }
    }
    return found
}
```

**Dry run** — text `"catsanddog"` (indices 0–9), words `cat, cats, and, sand, dog`:

| `start` | Walk | Words found |
|---|---|---|
| 0 | c → a → t\* → s\* → (no `a` after `cats`) stop | [0, 2] `cat`, [0, 3] `cats` |
| 1, 2 | `a` → no `t`; `t` → no node | — |
| 3 | s → a → n → d\* → (no `d` after `sand`) stop | [3, 6] `sand` |
| 4 | a → n → d\* → stop | [4, 6] `and` |
| 5, 6 | `n` → no node; d → (no `d` after `d`) stop | — |
| 7 | d → o → g\* | [7, 9] `dog` |

Every walk stops early. The longest walk is bounded by the longest word, not by the text.

**Run it:** `findWordsInText` in [`trie.go`](trie.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Word Break | Medium | https://leetcode.com/problems/word-break/ |
| 2 | Word Search II | Hard | https://leetcode.com/problems/word-search-ii/ |

**What to notice — Problem 1:** DP (DP week) where the transitions come from the template: `dp[start]` = "the text before `start` can be segmented". For each `start` with `dp[start] == true`, walk the trie and set `dp[end+1] = true` at every `isEnd`.

**What to notice — Problem 2:** Backtracking on the grid (Backtracking week, Day 6) while walking the trie in lock-step: move to a neighbor only if the trie has a child for its letter. Store the full word on its end node so you can record it without rebuilding the string, and clear it after recording to avoid duplicates.

**Complexity:** template O(n · L) for text length n and longest word L. Word Break O(n · L). Word Search II O(R · C · 4 · 3ᴸ⁻¹) worst case, much less in practice thanks to the trie pruning.

**Common mistakes:**

- Word Search II: running a separate grid search for each word — O(words × grid) — instead of one search with the trie.
- Word Search II: not pruning trie nodes (or clearing found words) — the same word is found many times.
- Forgetting `break` when the trie walk falls off — keeps walking a nil node.

---

## Day 6 — Binary Trie for XOR

**Pattern:** Insert numbers bit by bit (most significant first); greedily walk toward the opposite bit  
**Key insight:** XOR is maximized by making the **highest** bits differ. Starting from the top bit, pick the child with the opposite bit whenever it exists — a higher bit is worth more than all lower bits combined.

**Template adaptation:**

```go
const maxBit = 30 // enough for values below 2^31

type bitNode struct {
    children [2]*bitNode
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

// largest x XOR y over all inserted y
func (b *binaryTrie) maxXorWith(x int) int {
    node := b.root
    result := 0
    for bit := maxBit; bit >= 0; bit-- {
        d := (x >> bit) & 1
        if want := 1 - d; node.children[want] != nil {
            result |= 1 << bit // this bit of the XOR is 1
            node = node.children[want]
        } else {
            node = node.children[d] // forced: this bit of the XOR is 0
        }
    }
    return result
}
```

**Why greedy on bits is correct:** a 1 at bit `k` is worth `2ᵏ`, and all lower bits together are worth at most `2ᵏ − 1`. So getting the highest possible bit right always beats anything you could gain lower down — the greedy choice is safe (Greedy week, rule 1).

**Dry run** — numbers `[3, 10, 5, 25, 2, 8]`, query `x = 5` (`00101`), showing the low 5 bits:

| Bit | `x`'s bit | Want | Available? | XOR so far | Still possible |
|---|---|---|---|---|---|
| 4 | 0 | 1 | yes (25) | `10000` = 16 | 25 |
| 3 | 0 | 1 | yes | `11000` = 24 | 25 |
| 2 | 1 | 0 | yes | `11100` = 28 | 25 |
| 1 | 0 | 1 | no — follow 0 | 28 | 25 |
| 0 | 1 | 0 | no — follow 1 | 28 | 25 |

`5 XOR 25 = 28`, the maximum over all pairs.

**Run it:** `binaryTrie`, `insert`, and `maxXorWith` in [`trie.go`](trie.go).

**Problems:**

| # | Problem | Difficulty | Link |
|---|---------|------------|------|
| 1 | Maximum XOR of Two Numbers in an Array | Medium | https://leetcode.com/problems/maximum-xor-of-two-numbers-in-an-array/ |
| 2 | Maximum XOR With an Element From Array | Hard | https://leetcode.com/problems/maximum-xor-with-an-element-from-array/ |

**What to notice — Problem 1:** Insert everything, then take the max of `maxXorWith(x)` over every `x`. O(n · 31).

**What to notice — Problem 2:** Each query only allows elements `≤ m`. Answer queries **offline**: sort the numbers and the queries by `m`, and before answering a query insert every number `≤ m`. The trie then only ever contains allowed elements. (Return −1 when nothing is inserted yet.)

**Complexity:** both O((n + q) · 31) after sorting.

**Common mistakes:**

- Walking from the least significant bit — the greedy argument only works from the top.
- Using too few bits (`maxBit` must cover the largest value).
- Problem 2: answering queries in input order — you'd have to delete from the trie. Sort them and remember original indices.

---

## Day 7 — Mixed Review

**Pattern:** Consolidation under pressure  
**Key insight:** Trie problems differ only in what each node stores (`isEnd`, `count`, a word, suggestions) and how you walk (exact path, DFS on wildcards, restart at each position, opposite bit).

**No new problems today.** Do this instead:

1. Look back at Days 2–6. Pick the 2 days that felt hardest.
2. For each of those 2 days, redo one problem without looking at your notes.
3. Before coding, draw the trie for 3–4 example words and mark what each node stores.

**Self-assessment — can you answer these without hesitation?**

- [ ] Write `trieNode`, `Insert`, `Search`, and `StartsWith` from memory
- [ ] What's the difference between `Search` and `StartsWith`, in one line of code?
- [ ] How does a `.` wildcard change the search, and what's its worst-case cost?
- [ ] How do you count words with a given prefix in O(L)?
- [ ] Why does Word Search II use a trie instead of searching for each word separately?
- [ ] Why does the XOR trie go from the highest bit down?

---

## Pattern Cheatsheet

| Pattern | Node stores | Walk | Time per query |
|---------|-------------|------|----------------|
| Word / prefix lookup | `isEnd` | exact path | O(L) |
| Wildcard search | `isEnd` | DFS into all children at `.` | O(L) to O(26ᵈ · L) |
| Prefix counts / sums | `count` / `sum` updated on insert | exact path | O(L) |
| Autocomplete | top-k strings | exact path | O(L) |
| Many words in a text | `isEnd` (or the word) | restart at every position, stop on nil | O(n · L) total |
| Max XOR | 2 children per node | prefer the opposite bit, top bit first | O(bits) |

---

package urltemplate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PathTrieNode is one segment in a path trie.
// A node can be a leaf (an observed endpoint ends here) and still have children
// from longer paths that share this prefix.
type PathTrieNode struct {
	// Segment is the path segment value for this node (no leading slash).
	Segment string

	// Children keyed by segment value.
	Children map[string]*PathTrieNode

	// Leaf is true when at least one full endpoint ends at this node.
	Leaf bool

	// ObservationCount is the Redis observation count for the endpoint that
	// ends at this node when Leaf is true. Zero when Leaf is false.
	ObservationCount int64

	// UniqueLeaves is the number of distinct leaf endpoints under this node
	// (including this node when Leaf is true).
	UniqueLeaves int

	// TotalObservations is the sum of Redis observation counts for all leaf
	// endpoints under this node (including this node when Leaf is true).
	TotalObservations int64
}

// PathTrie indexes URL paths by segment for URL templatization live traffic learning.
type PathTrie struct {
	Root *PathTrieNode
}

func NewPathTrie() *PathTrie {
	return &PathTrie{
		Root: &PathTrieNode{
			Children: make(map[string]*PathTrieNode),
		},
	}
}

// Insert walks path segment-by-segment, creating nodes as needed, and marks
// the final segment as a leaf with the given Redis observation count.
// If the same path is inserted again, observation counts are summed.
func (t *PathTrie) Insert(path string, count int64) {
	segments, _ := SplitPath(path)
	node := t.Root
	for _, segment := range segments {
		child, ok := node.Children[segment]
		if !ok {
			child = &PathTrieNode{
				Segment:  segment,
				Children: make(map[string]*PathTrieNode),
			}
			node.Children[segment] = child
		}
		node = child
	}
	node.Leaf = true
	node.ObservationCount += count
}

// recompute fills UniqueLeaves and TotalObservations bottom-up for every node.
func (t *PathTrie) recompute() {
	t.Root.recompute()
}

func (n *PathTrieNode) recompute() {
	unique := 0
	total := int64(0)
	if n.Leaf {
		unique = 1
		total = n.ObservationCount
	}
	for _, child := range n.Children {
		child.recompute()
		unique += child.UniqueLeaves
		total += child.TotalObservations
	}
	n.UniqueLeaves = unique
	n.TotalObservations = total
}

// BuildPathTrie builds a path trie from a Redis hash of path -> count string.
func BuildPathTrie(pathCounts map[string]string) *PathTrie {
	trie := NewPathTrie()
	for path, countStr := range pathCounts {
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			continue
		}
		trie.Insert(path, count)
	}
	trie.recompute()
	return trie
}

// BuildPathTrieFromCounts builds a path trie from path -> observation count.
func BuildPathTrieFromCounts(pathCounts map[string]int64) *PathTrie {
	trie := NewPathTrie()
	for path, count := range pathCounts {
		trie.Insert(path, count)
	}
	trie.recompute()
	return trie
}

// String returns a multi-line debug rendering of the trie.
func (t *PathTrie) String() string {
	var b strings.Builder
	t.Root.writeString(&b, 0, "/")
	return b.String()
}

func (n *PathTrieNode) writeString(b *strings.Builder, depth int, pathSoFar string) {
	if depth > 0 {
		fmt.Fprintf(b, "%s%s leaf=%v obs=%d uniqueLeaves=%d totalObs=%d\n",
			strings.Repeat("  ", depth-1),
			pathSoFar,
			n.Leaf,
			n.ObservationCount,
			n.UniqueLeaves,
			n.TotalObservations,
		)
	}

	names := make([]string, 0, len(n.Children))
	for name := range n.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		child := n.Children[name]
		childPath := pathSoFar
		if pathSoFar == "/" {
			childPath = "/" + name
		} else {
			childPath = pathSoFar + "/" + name
		}
		child.writeString(b, depth+1, childPath)
	}
}

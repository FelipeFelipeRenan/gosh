package trie

type Node struct {
	Children map[rune]*Node
	IsEnd    bool
}

type Trie struct {
	Root *Node
}

func New() *Trie {
	return &Trie{Root: &Node{Children: make(map[rune]*Node)}}
}

func (t *Trie) Insert(word string) {
	node := t.Root
	for _, char := range word {
		if _, ok := node.Children[char]; !ok {
			node.Children[char] = &Node{Children: make(map[rune]*Node)}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
}

func (t *Trie) SearchPrefix(prefix string) []string {
	node := t.Root
	for _, char := range prefix {
		if _, ok := node.Children[char]; !ok {
			return nil
		}
		node = node.Children[char]
	}
	var results []string
	t.collect(node, prefix, &results)
	return results
}

func (t *Trie) collect(node *Node, current string, results *[]string) {
	if node.IsEnd {
		*results = append(*results, current)
	}
	for char, child := range node.Children {
		t.collect(child, current+string(char), results)
	}
}

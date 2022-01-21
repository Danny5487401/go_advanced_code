package trie

// Trie 树节点
type trieNode struct {
	char     string             // Unicode 字符
	isEnding bool               // 是否是单词结尾
	children map[rune]*trieNode // 该节点的子节点字典
}

// 初始化 Trie 树节点
func NewTrieNode(char string) *trieNode {
	return &trieNode{
		char:     char,
		isEnding: false,
		children: make(map[rune]*trieNode),
	}
}

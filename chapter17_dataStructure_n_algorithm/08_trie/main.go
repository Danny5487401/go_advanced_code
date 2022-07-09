package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/08_trie/trie"
)

func main() {
	trie := trie.NewTrie()
	words := []string{"Golang", "学院君", "Language", "Trie", "Go", "Danny"}
	// 构建 Trie 树
	for _, word := range words {
		trie.Insert(word)
	}
	// 从 Trie 树中查找字符串
	//term := "学院君"
	term := "Dann"
	if trie.Find(term) {
		fmt.Printf("包含单词\"%s\"\n", term)
	} else {
		fmt.Printf("不包含单词\"%s\"\n", term)
	}
}

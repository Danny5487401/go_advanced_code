package main

import (
	"fmt"

	trieV3 "github.com/derekparker/trie/v3"
)

func main() {
	trieInfo := trieV3.New[any]()
	words := []string{"Golang", "学院君", "Language", "Trie", "Go", "Danny", "Danny2"}
	// 构建 Trie 树
	for _, word := range words {
		trieInfo.Add(word, len(word))
	}
	// 从 Trie 树中查找字符串
	//term := "学院君"
	term := "Danny"
	if _, ok := trieInfo.Find(term); ok {
		fmt.Printf("包含单词 %v \n", term)
	} else {
		fmt.Printf("不包含单词\"%s\"\n", term)
	}

	prefixInfo := trieInfo.PrefixSearch("Dan")
	fmt.Println(prefixInfo)
}

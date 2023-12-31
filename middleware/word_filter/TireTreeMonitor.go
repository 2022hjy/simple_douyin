package word_filter

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"os"
)

const replaceString = "&&&"

func NewWordFilterMiddleware(txtFilePath string) (gin.HandlerFunc, error) {
	trie, err := NewTrieFromTxtFile(txtFilePath)
	if err != nil {
		return nil, err
	}
	// 路径位置：middleware/word_filter/sensitive_words.txt
	return func(c *gin.Context) {
		text := c.Query("comment_text")
		if text != "" {
			filteredText := trie.Filter(text)
			c.Set("filteredText", filteredText)
		}

		// 聊天记录模块
		content := c.Query("content")
		if content != "" {
			filteredContent := trie.Filter(content)
			c.Set("filteredContent", filteredContent)
		}

		c.Next()
	}, nil
}

func NewTrieFromTxtFile(path string) (*Trie, error) {
	trie := NewTrie()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		trie.Insert(word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return trie, nil
}

// TrieNode 前缀树的结点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Trie 前缀树
type Trie struct {
	root *TrieNode
}

// NewTrie 创建新的前缀树
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
			isEnd:    false,
		},
	}
}

// Insert 插入一个敏感词
func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isEnd:    false,
			}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Search 查找敏感词
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

// Filter 过滤敏感词
func (t *Trie) Filter(text string) string {
	// 加上一个字符，方便敏感词在最后的处理
	text = text + "!"
	// 设置当前节点为根节点
	node := t.root
	// 记录单词的开始位置
	start := 0
	// 用来存储过滤后的字符串
	var result []rune
	// 遍历字符串
	for i, char := range text {
		// 如果当前字符不在字典树中，说明当前子串不存在于字典树中。
		if _, ok := node.children[char]; !ok {
			// 如果当前单词的起始位置不等于结束位置，说明之前已经找到了一个单词，则需要将该单词替换成 replaceString
			if start != i && node.isEnd {
				result = append(result, []rune(replaceString)...)
			}
			// 直接加入过滤后的字符串
			result = append(result, char)
			// 重置当前节点为根节点，起始位置为下一个字符位置
			node = t.root
			start = i + 1
		} else {
			// 如果当前 Trie 节点存在该字符，说明存在包含该字符的单词
			// 如果该 Trie 节点是一个单词的结尾，则将单词的起始位置设为当前字符位置
			if node.isEnd {
				result = append(result, []rune(replaceString)...)
				start = i + 1
			}
			node = node.children[char]
		}
	}
	res := string(result)
	// 去掉刚才的字符
	return res[:len(res)-1]
}

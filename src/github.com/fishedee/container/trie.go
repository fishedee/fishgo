package container

import (
	"sort"
	"unicode/utf8"
)

type TrieTree struct {
	nodeNum int
	root    *trieTreeNode
}

type TrieTreeWalker func(key string, value interface{}, parentKey string, parentValue interface{})

type trieTreeNode struct {
	value    interface{}
	children map[rune]*trieTreeNode
}

type trieWalk struct {
	key  string
	node *trieTreeNode
}

func NewTrieTree() *TrieTree {
	trieTree := &TrieTree{}
	trieTree.nodeNum = 0
	trieTree.root = &trieTreeNode{
		value:    nil,
		children: map[rune]*trieTreeNode{},
	}
	return trieTree
}

func (this *TrieTree) Walk(walker TrieTreeWalker) {
	queue := NewQueue()
	queue.Push(&trieWalk{
		key:  "",
		node: this.root,
	})
	walker("", this.root.value, "", nil)
	for queue.Len() != 0 {
		top := queue.Pop().(*trieWalk)
		var keySort []int
		for key, _ := range top.node.children {
			keySort = append(keySort, int(key))
		}
		sort.Ints(keySort)
		for _, next := range keySort {
			child := top.node.children[rune(next)]
			childKey := top.key + string(rune(next))
			walker(childKey, child.value, top.key, top.node.value)
			queue.Push(&trieWalk{
				key:  childKey,
				node: child,
			})
		}
	}
}

func (this *TrieTree) Get(key string) interface{} {
	current := this.root
	for _, char := range key {
		child, isExist := current.children[char]
		if isExist == false {
			return nil
		}
		current = child
	}
	return current.value
}

func (this *TrieTree) Set(key string, value interface{}) {
	current := this.root
	for _, char := range key {
		child, isExist := current.children[char]
		if isExist == false {
			child = &trieTreeNode{
				value:    nil,
				children: map[rune]*trieTreeNode{},
			}
			current.children[char] = child
			this.nodeNum++
		}
		current = child
	}
	current.value = value
}

func (this *TrieTree) ToTrieArray() *TrieArray {
	trie := newTrieArray()
	trie.build(this.root)
	return trie
}

type TrieArray struct {
	base  []int
	check []int
	value []interface{}
	dict  *HashListArray
}

type TrieMatch struct {
	key   string
	value interface{}
}

type trieArrayNode struct {
	index  int
	offset int
	node   *trieTreeNode
}

func newTrieArray() *TrieArray {
	trie := &TrieArray{}
	trie.base = []int{}
	trie.check = []int{}
	trie.value = []interface{}{}
	return trie
}

func (this *TrieArray) build(root *trieTreeNode) {
	dict := map[rune]int{}

	queue := NewQueue()
	queue.Push(&trieArrayNode{
		index:  1,
		offset: 0,
		node:   root,
	})
	this.setCheck(1, -1)
	this.setValue(1, root.value)
	for queue.Len() != 0 {
		top := queue.Pop().(*trieArrayNode)
		freeOffset := this.findIndex(top.offset, top.node.children, dict)
		this.setBase(top.index, freeOffset)
		for next, child := range top.node.children {
			childIndex := freeOffset + dict[next]
			this.setCheck(childIndex, top.index)
			this.setValue(childIndex, child.value)
			queue.Push(&trieArrayNode{
				index:  childIndex,
				offset: freeOffset,
				node:   child,
			})
		}
	}

	hashList := NewHashList(len(dict) * 2)
	for key, value := range dict {
		hashList.Set(int(key), value)
	}
	this.dict = hashList.ToHashListArray()
}

func (this *TrieArray) findIndex(offset int, next map[rune]*trieTreeNode, dict map[rune]int) int {

	for char, _ := range next {
		_, isExist := dict[char]
		if isExist == false {
			dict[char] = len(dict) + 1
		}
	}

	for {
		isOk := true
		for char, _ := range next {
			index := offset + dict[char]
			if this.isExist(index) {
				isOk = false
				break
			}
		}
		if isOk {
			break
		}
		offset++
	}
	return offset
}

func (this *TrieArray) isExist(index int) bool {
	if index < 0 || index >= len(this.check) {
		return false
	}
	return this.check[index] != 0
}

func (this *TrieArray) expand(index int) {
	if index < len(this.check) {
		return
	}
	newSize := index - len(this.check) + 1
	newBase := make([]int, newSize)
	newCheck := make([]int, newSize)
	newValue := make([]interface{}, newSize)
	this.base = append(this.base, newBase...)
	this.check = append(this.check, newCheck...)
	this.value = append(this.value, newValue...)
}

func (this *TrieArray) setBase(index int, base int) {
	this.expand(index)
	this.base[index] = base
}

func (this *TrieArray) setCheck(index int, check int) {
	this.expand(index)
	this.check[index] = check
}
func (this *TrieArray) setValue(index int, value interface{}) {
	this.expand(index)
	this.value[index] = value
}

func (this *TrieArray) LongestPrefixMatch(key string) (string, interface{}) {
	length := len(this.check)
	current := 1
	var result interface{}

	if this.value[current] != nil {
		result = this.value[current]
	}

	index := 0
	for _, char := range key {
		charId := this.dict.Get(int(char))
		if charId == nil {
			break
		}
		next := this.base[current] + charId.(int)
		if next >= length {
			break
		}
		if this.check[next] != current {
			break
		}
		index += utf8.RuneLen(char)
		current = next
		if this.value[current] != nil {
			result = this.value[current]
		}
	}
	return key[0:index], result
}

func (this *TrieArray) ExactMatch(key string) interface{} {
	resultKey, resultValue := this.LongestPrefixMatch(key)
	if len(resultKey) != len(key) {
		return nil
	}
	return resultValue
}

func (this *TrieArray) PrefixMatch(key string) []TrieMatch {
	length := len(this.check)
	current := 1
	result := []TrieMatch{}

	if this.value[current] != nil {
		result = append(result, TrieMatch{
			key:   "",
			value: this.value[current],
		})
	}

	index := 0
	for _, char := range key {
		dictId := this.dict.Get(int(char))
		if dictId == nil {
			break
		}
		next := this.base[current] + dictId.(int)
		if next >= length {
			break
		}
		if this.check[next] != current {
			break
		}
		index += utf8.RuneLen(char)
		if this.value[next] != nil {
			result = append(result, TrieMatch{
				key:   key[0:index],
				value: this.value[next],
			})
		}
		current = next
	}
	return result
}

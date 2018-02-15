package container

import (
	. "github.com/fishedee/language"
)

type TrieTree struct {
	nodeNum int
	root    *trieTreeNode
}

type TrieTreeWalker func(key string, value interface{}, parentKey string, parentValue interface{})

type trieTreeNode struct {
	value    interface{}
	children map[byte]*trieTreeNode
}

type trieWalk struct {
	key  []byte
	node *trieTreeNode
}

func NewTrieTree() *TrieTree {
	trieTree := &TrieTree{}
	trieTree.nodeNum = 0
	trieTree.root = &trieTreeNode{
		value:    nil,
		children: map[byte]*trieTreeNode{},
	}
	return trieTree
}

func (this *TrieTree) Walk(walker TrieTreeWalker) {
	queue := NewQueue()
	queue.Push(&trieWalk{
		key:  []byte(""),
		node: this.root,
	})
	walker("", this.root.value, "", nil)
	for queue.Len() != 0 {
		top := queue.Pop().(*trieWalk)
		keySortInterface, _ := ArrayKeyAndValue(top.node.children)
		keySort := keySortInterface.([]byte)
		for _, next := range keySort {
			child := top.node.children[next]
			childKey := append(top.key, next)
			walker(string(childKey), child.value, string(top.key), top.node.value)
			queue.Push(&trieWalk{
				key:  childKey,
				node: child,
			})
		}
	}
}

func (this *TrieTree) Get(key string) interface{} {
	current := this.root
	for i := 0; i != len(key); i++ {
		char := key[i]
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
	for i := 0; i != len(key); i++ {
		char := key[i]
		child, isExist := current.children[char]
		if isExist == false {
			child = &trieTreeNode{
				value:    nil,
				children: map[byte]*trieTreeNode{},
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
		freeOffset := this.findIndex(top.offset, top.node.children)
		this.setBase(top.index, freeOffset)
		for next, child := range top.node.children {
			childIndex := freeOffset + int(next)
			this.setCheck(childIndex, top.index)
			this.setValue(childIndex, child.value)
			queue.Push(&trieArrayNode{
				index:  childIndex,
				offset: freeOffset,
				node:   child,
			})
		}

	}
}

func (this *TrieArray) findIndex(offset int, next map[byte]*trieTreeNode) int {
	for {
		isOk := true
		for char, _ := range next {
			index := offset + int(char)
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

func (this *TrieArray) singleMatch(key string) (interface{}, bool) {
	length := len(this.check)
	current := 1
	var result interface{}

	if this.value[current] != nil {
		result = this.value[current]
	}

	i := 0
	for ; i != len(key); i++ {
		next := this.base[current] + int(key[i])
		if next >= length {
			break
		}
		if this.check[next] != current {
			break
		}
		current = next
		if this.value[current] != nil {
			result = this.value[current]
		}
	}
	return result, i == len(key)
}

func (this *TrieArray) LongestPrefixMatch(key string) interface{} {
	result, _ := this.singleMatch(key)
	return result
}

func (this *TrieArray) ExactMatch(key string) interface{} {
	result, isExact := this.singleMatch(key)
	if isExact == false {
		return nil
	}
	return result
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

	for i := 0; i != len(key); i++ {
		next := this.base[current] + int(key[i])
		if next >= length {
			break
		}
		if this.check[next] != current {
			break
		}
		if this.value[next] != nil {
			result = append(result, TrieMatch{
				key:   key[0 : i+1],
				value: this.value[next],
			})
		}
		current = next
	}
	return result
}

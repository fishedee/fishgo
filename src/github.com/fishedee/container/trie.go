package container

import (
	//"fmt"
	"sort"
)

type TrieTree struct {
	root *trieTreeNode
}

type TrieTreeWalker func(key string, value interface{}, parentKey string, parentValue interface{})

type trieTreeNode struct {
	segment  string
	value    interface{}
	children map[byte]*trieTreeNode
}

type trieWalk struct {
	key  []byte
	node *trieTreeNode
}

func NewTrieTree() *TrieTree {
	trieTree := &TrieTree{}
	trieTree.root = nil
	return trieTree
}

func (this *TrieTree) Walk(walker TrieTreeWalker) {
	if this.root == nil {
		return
	}
	queue := NewQueue()
	queue.Push(&trieWalk{
		key:  []byte(this.root.segment),
		node: this.root,
	})
	walker(this.root.segment, this.root.value, "", nil)
	for queue.Len() != 0 {
		top := queue.Pop().(*trieWalk)
		var keySort []int
		for key, _ := range top.node.children {
			keySort = append(keySort, int(key))
		}
		sort.Ints(keySort)
		for _, next := range keySort {
			nextChar := byte(next)
			child := top.node.children[nextChar]
			childKey := make([]byte, 0, len(top.key)+1+len(child.segment))
			childKey = append(childKey, top.key...)
			childKey = append(childKey, nextChar)
			childKey = append(childKey, []byte(child.segment)...)
			walker(string(childKey), child.value, string(top.key), top.node.value)
			queue.Push(&trieWalk{
				key:  childKey,
				node: child,
			})
		}
	}
}

func (this *TrieTree) Get(key string) interface{} {
	if this.root == nil {
		return nil
	}
	current := this.root
	for {
		segment := current.segment
		if len(key) < len(segment) {
			return nil
		}
		if key[:len(segment)] != segment {
			return nil
		}
		if len(key) == len(segment) {
			return current.value
		}
		char := key[len(segment)]
		next, isExist := current.children[char]
		if isExist == false {
			return nil
		}
		current = next
		key = key[len(segment)+1:]
	}
	return nil
}

func (this *TrieTree) Set(key string, value interface{}) {
	if this.root == nil {
		this.root = &trieTreeNode{
			segment:  key,
			value:    value,
			children: map[byte]*trieTreeNode{},
		}
		return
	}
	current := this.root
	for {
		i := 0
		j := 0
		segment := current.segment
		for i < len(key) && j < len(segment) {
			if key[i] != segment[j] {
				break
			}
			i++
			j++
		}
		if i < len(key) && j < len(segment) {
			//分裂成三个节点
			newChild := &trieTreeNode{
				segment:  segment[j+1:],
				value:    current.value,
				children: current.children,
			}
			insertChild := &trieTreeNode{
				segment:  key[i+1:],
				value:    value,
				children: map[byte]*trieTreeNode{},
			}
			current.segment = segment[:j]
			current.value = nil
			current.children = map[byte]*trieTreeNode{
				segment[j]: newChild,
				key[i]:     insertChild,
			}
			break
		} else if i == len(key) && j == len(segment) {
			//正好是当前节点
			current.value = value
			break
		} else if i < len(key) && j == len(segment) {
			//复用当前节点
			key = key[i:]
		} else {
			//分裂成两个节点
			newChild := &trieTreeNode{
				segment:  segment[j+1:],
				value:    current.value,
				children: current.children,
			}
			current.segment = segment[:j]
			current.value = value
			current.children = map[byte]*trieTreeNode{
				segment[j]: newChild,
			}
			break
		}
		char := key[0]
		key = key[1:]
		child, isExist := current.children[char]
		if isExist == false {
			current.children[char] = &trieTreeNode{
				segment:  key,
				value:    value,
				children: map[byte]*trieTreeNode{},
			}
			break
		}
		current = child
	}
}

func (this *TrieTree) ToTrieArray() *TrieArray {
	trie := newTrieArray()
	trie.build(this.root)
	return trie
}

type TrieArray struct {
	data []trieArrayData
}

type TrieMatch struct {
	Key   string
	Value interface{}
}

type trieArrayData struct {
	base    int
	check   int
	value   interface{}
	segment string
}
type trieArrayNode struct {
	index  int
	offset int
	node   *trieTreeNode
}

func newTrieArray() *TrieArray {
	trie := &TrieArray{}
	trie.data = []trieArrayData{}
	return trie
}

func (this *TrieArray) build(root *trieTreeNode) {
	if root == nil {
		return
	}
	queue := NewQueue()
	queue.Push(&trieArrayNode{
		index:  1,
		offset: 0,
		node:   root,
	})
	this.setCheck(1, -1)
	this.setValue(1, root.value)
	this.setSegment(1, root.segment)
	for queue.Len() != 0 {
		top := queue.Pop().(*trieArrayNode)
		freeOffset := this.findIndex(top.offset, top.node.children)
		this.setBase(top.index, freeOffset)
		for next, child := range top.node.children {
			childIndex := freeOffset + int(next)
			this.setCheck(childIndex, top.index)
			this.setValue(childIndex, child.value)
			this.setSegment(childIndex, child.segment)
			queue.Push(&trieArrayNode{
				index:  childIndex,
				offset: freeOffset,
				node:   child,
			})
		}
	}
	//fmt.Println(len(this.data))
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
	if index < 0 || index >= len(this.data) {
		return false
	}
	return this.data[index].check != 0
}

func (this *TrieArray) expand(index int) {
	if index < len(this.data) {
		return
	}
	newSize := index - len(this.data) + 1
	newData := make([]trieArrayData, newSize)
	this.data = append(this.data, newData...)
}

func (this *TrieArray) setBase(index int, base int) {
	this.expand(index)
	this.data[index].base = base
}

func (this *TrieArray) setCheck(index int, check int) {
	this.expand(index)
	this.data[index].check = check
}
func (this *TrieArray) setValue(index int, value interface{}) {
	this.expand(index)
	this.data[index].value = value
}

func (this *TrieArray) setSegment(index int, segment string) {
	this.expand(index)
	this.data[index].segment = segment
}

func (this *TrieArray) LongestPrefixMatchWithChar(key string, extChar byte) (interface{}, bool) {
	var resultValue interface{}

	length := len(this.data)
	current := 1
	for {
		segment := this.data[current].segment
		if len(key)+1 < len(segment) {
			return resultValue, false
		} else if len(key)+1 == len(segment) {
			if key != segment[0:len(key)] {
				return resultValue, false
			}
			if extChar != segment[len(key)] {
				return resultValue, false
			}
			if this.data[current].value != nil {
				resultValue = this.data[current].value
			}
			return resultValue, true
		}

		if key[:len(segment)] != segment {
			return resultValue, false
		}
		if this.data[current].value != nil {
			resultValue = this.data[current].value
		}
		if len(key) == len(segment) {
			char := extChar
			next := this.data[current].base + int(char)
			if next >= length {
				return resultValue, false
			}
			if this.data[next].check != current {
				return resultValue, false
			}
			current = next
			if this.data[current].segment != "" {
				return resultValue, false
			}
			if this.data[current].value != nil {
				resultValue = this.data[current].value
			}
			return resultValue, true
		} else {
			key = key[len(segment):]
			char := key[0]
			next := this.data[current].base + int(char)
			if next >= length {
				return resultValue, false
			}
			if this.data[next].check != current {
				return resultValue, false
			}
			current = next
			key = key[1:]
		}
	}
}

func (this *TrieArray) LongestPrefixMatch(key string) (string, interface{}) {
	var resultValue interface{}
	var resultKey string

	length := len(this.data)
	current := 1
	origin := key
	for {
		segment := this.data[current].segment
		if len(key) < len(segment) {
			break
		}
		if key[:len(segment)] != segment {
			break
		}
		if len(key) >= len(segment) {
			key = key[len(segment):]
			if this.data[current].value != nil {
				resultKey = origin[0 : len(origin)-len(key)]
				resultValue = this.data[current].value
			}
		}
		if len(key) == 0 {
			break
		}

		char := key[0]
		next := this.data[current].base + int(char)
		if next >= length {
			break
		}
		if this.data[next].check != current {
			break
		}
		current = next
		key = key[1:]
	}

	return resultKey, resultValue
}

func (this *TrieArray) ExactMatch(key string) interface{} {
	resultKey, resultValue := this.LongestPrefixMatch(key)
	if len(resultKey) != len(key) {
		return nil
	}
	return resultValue
}

func (this *TrieArray) PrefixMatch(key string) []TrieMatch {
	result := []TrieMatch{}

	length := len(this.data)
	current := 1
	origin := key
	for {
		segment := this.data[current].segment
		if len(key) < len(segment) {
			break
		}
		if key[:len(segment)] != segment {
			break
		}
		if len(key) >= len(segment) {
			key = key[len(segment):]
			if this.data[current].value != nil {
				result = append(result, TrieMatch{
					Key:   origin[0 : len(origin)-len(key)],
					Value: this.data[current].value,
				})
			}
		}
		if len(key) == 0 {
			break
		}

		char := key[0]
		next := this.data[current].base + int(char)
		if next >= length {
			break
		}
		if this.data[next].check != current {
			break
		}
		current = next
		key = key[1:]
	}

	return result
}

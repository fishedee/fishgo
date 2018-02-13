package container

import (
	"fmt"
	. "github.com/fishedee/language"
)

type RadixFactory struct {
	nodeNum   int
	matchMode []int
	root      *radixTreeNode
}

type radixTreeNode struct {
	match    []interface{}
	children map[byte]*radixTreeNode
}

var RadixMatchMode struct {
	EnumStruct
	ALL    int `enum:"1,精确匹配"`
	PREFIX int `enum:"2,前缀匹配"`
}

func NewRadixFactory(matchMode []int) *RadixFactory {
	radixFactory := &RadixFactory{}
	radixFactory.matchMode = matchMode
	radixFactory.nodeNum = 0
	radixFactory.root = &radixTreeNode{
		match:    make([]interface{}, len(matchMode)),
		children: map[byte]*radixTreeNode{},
	}
	return radixFactory
}

func (this *RadixFactory) Insert(mode int, key string, value interface{}) error {
	if mode < 0 || mode >= len(this.matchMode) {
		return fmt.Errorf("invalid mode:%v,modeSize:%v", mode, len(this.matchMode))
	}
	current := this.root
	for _, char := range []byte(key) {
		child, isExist := current.children[char]
		if isExist == false {
			child = &radixTreeNode{
				match:    make([]interface{}, len(this.matchMode)),
				children: map[byte]*radixTreeNode{},
			}
			current.children[char] = child
			this.nodeNum++
		}
		current = child
	}
	if current.match[mode] != nil {
		return fmt.Errorf("already has exist!mode:%v,key:%v", mode, key)
	}
	current.match[mode] = value
	return nil
}

func (this *RadixFactory) Create() *Radix {
	radix := newRadix(this.matchMode)
	radix.build(this.root)
	return radix
}

type Radix struct {
	matchMode []int
	base      []int
	check     []int
	value     [][]interface{}
}

type radixArrayNode struct {
	index  int
	offset int
	node   *radixTreeNode
}

func newRadix(matchMode []int) *Radix {
	radix := &Radix{}
	radix.matchMode = matchMode
	radix.base = []int{}
	radix.check = []int{}
	radix.value = [][]interface{}{}
	return radix
}

func (this *Radix) build(root *radixTreeNode) {
	queue := NewQueue()
	queue.Push(&radixArrayNode{
		index:  1,
		offset: 0,
		node:   root,
	})
	this.setCheck(1, -1)
	this.setValue(1, root.match)
	for queue.Len() != 0 {
		top := queue.Pop().(*radixArrayNode)
		freeOffset := this.findIndex(top.offset, top.node.children)
		this.setBase(top.index, freeOffset)
		for next, child := range top.node.children {
			childIndex := freeOffset + int(next)
			this.setCheck(childIndex, top.index)
			this.setValue(childIndex, child.match)
			queue.Push(&radixArrayNode{
				index:  childIndex,
				offset: freeOffset,
				node:   child,
			})
		}

	}
}

func (this *Radix) findIndex(offset int, next map[byte]*radixTreeNode) int {
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

func (this *Radix) isExist(index int) bool {
	if index < 0 || index >= len(this.check) {
		return false
	}
	return this.check[index] != 0
}

func (this *Radix) expand(index int) {
	if index < len(this.check) {
		return
	}
	newSize := index - len(this.check) + 1
	newBase := make([]int, newSize)
	newCheck := make([]int, newSize)
	newValue := make([][]interface{}, newSize)
	this.base = append(this.base, newBase...)
	this.check = append(this.check, newCheck...)
	this.value = append(this.value, newValue...)
}

func (this *Radix) setBase(index int, base int) {
	this.expand(index)
	this.base[index] = base
}

func (this *Radix) setCheck(index int, check int) {
	this.expand(index)
	this.check[index] = check
}
func (this *Radix) setValue(index int, value []interface{}) {
	this.expand(index)
	this.value[index] = value
}

func (this *Radix) Find(key string) []interface{} {
	length := len(this.check)
	current := 1
	result := make([]interface{}, len(this.matchMode))
	hasAllMatch := true

	for singleMatchIndex, singleMatchMode := range this.matchMode {
		if singleMatchMode != RadixMatchMode.PREFIX {
			continue
		}
		if this.value[current][singleMatchIndex] != nil {
			result[singleMatchIndex] = this.value[current][singleMatchIndex]
		}
	}
	for _, char := range []byte(key) {
		next := this.base[current] + int(char)
		if next >= length {
			hasAllMatch = false
			break
		}
		if this.check[next] != current {
			hasAllMatch = false
			break
		}
		for singleMatchIndex, singleMatchMode := range this.matchMode {
			if singleMatchMode != RadixMatchMode.PREFIX {
				continue
			}
			if this.value[next][singleMatchIndex] != nil {
				result[singleMatchIndex] = this.value[next][singleMatchIndex]
			}
		}
		current = next
	}
	if hasAllMatch {
		for singleMatchIndex, singleMatchMode := range this.matchMode {
			if singleMatchMode != RadixMatchMode.ALL {
				continue
			}
			if this.value[current][singleMatchIndex] != nil {
				result[singleMatchIndex] = this.value[current][singleMatchIndex]
			}
		}
	}
	return result
}

func init() {
	InitEnumStruct(&RadixMatchMode)
}

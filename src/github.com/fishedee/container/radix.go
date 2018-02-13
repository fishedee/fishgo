package container

type RadixFactory struct {
	nodeNum int
	root    *radixTreeNode
}

type radixTreeNode struct {
	value    interface{}
	children map[byte]*radixTreeNode
}

func NewRadixFactory() *RadixFactory {
	radixFactory := &RadixFactory{}
	radixFactory.nodeNum = 0
	radixFactory.root = &radixTreeNode{
		value:    nil,
		children: map[byte]*radixTreeNode{},
	}
	return radixFactory
}

func (this *RadixFactory) Get(key string) interface{} {
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

func (this *RadixFactory) Set(key string, value interface{}) {
	current := this.root
	for i := 0; i != len(key); i++ {
		char := key[i]
		child, isExist := current.children[char]
		if isExist == false {
			child = &radixTreeNode{
				value:    nil,
				children: map[byte]*radixTreeNode{},
			}
			current.children[char] = child
			this.nodeNum++
		}
		current = child
	}
	current.value = value
}

func (this *RadixFactory) Create() *Radix {
	radix := newRadix()
	radix.build(this.root)
	return radix
}

type Radix struct {
	base  []int
	check []int
	value []interface{}
}

type RadixMatch struct {
	key   string
	value interface{}
}

type radixArrayNode struct {
	index  int
	offset int
	node   *radixTreeNode
}

func newRadix() *Radix {
	radix := &Radix{}
	radix.base = []int{}
	radix.check = []int{}
	radix.value = []interface{}{}
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
	this.setValue(1, root.value)
	for queue.Len() != 0 {
		top := queue.Pop().(*radixArrayNode)
		freeOffset := this.findIndex(top.offset, top.node.children)
		this.setBase(top.index, freeOffset)
		for next, child := range top.node.children {
			childIndex := freeOffset + int(next)
			this.setCheck(childIndex, top.index)
			this.setValue(childIndex, child.value)
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
	newValue := make([]interface{}, newSize)
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
func (this *Radix) setValue(index int, value interface{}) {
	this.expand(index)
	this.value[index] = value
}

func (this *Radix) ExactMatch(key string) interface{} {
	length := len(this.check)
	current := 1

	for i := 0; i != len(key); i++ {
		next := this.base[current] + int(key[i])
		if next >= length {
			return nil
		}
		if this.check[next] != current {
			return nil
		}
		current = next
	}
	return this.value[current]
}

func (this *Radix) PrefixMatch(key string) []RadixMatch {
	length := len(this.check)
	current := 1
	result := []RadixMatch{}

	if this.value[current] != nil {
		result = append(result, RadixMatch{
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
			result = append(result, RadixMatch{
				key:   key[0 : i+1],
				value: this.value[next],
			})
		}
		current = next
	}
	return result
}

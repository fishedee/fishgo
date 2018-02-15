package container

func hashInt(key int, mask int) int {
	return key & mask
}

type hashListNode struct {
	key   int
	value interface{}
	next  *hashListNode
}

type HashList struct {
	slot     []*hashListNode
	slotSize int
	slotMask int
	size     int
}

func NewHashList(size int) *HashList {
	hashList := &HashList{}
	hashList.slotSize = hashList.getSlotSize(size)
	hashList.slotMask = hashList.slotSize - 1
	hashList.slot = make([]*hashListNode, hashList.slotSize, hashList.slotSize)
	headListNode := make([]hashListNode, hashList.slotSize, hashList.slotSize)
	for i := 0; i != hashList.slotSize; i++ {
		hashList.slot[i] = &headListNode[i]
	}
	hashList.size = 0
	return hashList
}

func (this *HashList) getSlotSize(size int) int {
	size = int(float64(size) / 0.75)
	for i := 2; ; i *= 2 {
		if i > size {
			return i
		}
	}
	panic("invalid slot size")
}

func (this *HashList) hash(key int) int {
	return hashInt(key, this.slotMask)
}

func (this *HashList) find(key int) (*hashListNode, *hashListNode, *hashListNode) {
	hash := this.hash(key)

	head := this.slot[hash]
	prev := head
	current := prev.next
	for current != nil {
		if current.key == key {
			return head, prev, current
		}
		prev = current
		current = current.next
	}
	return head, prev, nil
}

func (this *HashList) Set(key int, value interface{}) {
	head, prev, current := this.find(key)
	if current != nil {
		current.value = value
	} else {
		prev.next = &hashListNode{
			key:   key,
			value: value,
			next:  nil,
		}
		head.key++
		this.size++
	}
}

func (this *HashList) Del(key int) {
	head, prev, current := this.find(key)
	if current != nil {
		prev.next = current.next
		this.size--
		head.key--
	}
}

func (this *HashList) Get(key int) interface{} {
	_, _, current := this.find(key)
	if current != nil {
		return current.value
	} else {
		return nil
	}
}

func (this *HashList) Len() int {
	return this.size
}

func (this *HashList) ToHashListArray() *HashListArray {
	hashListArray := newHashListArray()
	hashListArray.build(this.slot)
	return hashListArray
}

type hashListArrayNode struct {
	key   int
	value interface{}
}

type HashListArray struct {
	slot     [][]hashListArrayNode
	slotSize int
	slotMask int
	size     int
}

func newHashListArray() *HashListArray {
	hashListArray := &HashListArray{}
	return hashListArray
}

func (this *HashListArray) build(slot []*hashListNode) {
	this.slotSize = len(slot)
	this.slotMask = this.slotSize - 1
	this.slot = make([][]hashListArrayNode, this.slotSize, this.slotSize)
	allCount := 0
	for index, singleSlot := range slot {
		slotArray := make([]hashListArrayNode, singleSlot.key, singleSlot.key)
		slotArrayIndex := 0
		for head := singleSlot.next; head != nil; head = head.next {
			slotArray[slotArrayIndex] = hashListArrayNode{
				key:   head.key,
				value: head.value,
			}
			slotArrayIndex++
			allCount++
		}
		this.slot[index] = slotArray
	}
	this.size = allCount
}

func (this *HashListArray) Get(key int) interface{} {
	hash := hashInt(key, this.slotMask)
	slot := this.slot[hash]
	for i := 0; i != len(slot); i++ {
		if slot[i].key == key {
			return slot[i].value
		}
	}
	return nil
}

func (this *HashListArray) Len() int {
	return this.size
}

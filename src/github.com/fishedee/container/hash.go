package container

var hashPrimeList = []int{
	53, 97, 193, 389, 769,
	1543, 3079, 6151, 12289, 24593,
	49157, 98317, 196613, 393241,
	786433, 1572869, 3145739, 6291469, 12582917,
	25165843, 50331653, 100663319, 201326611, 402653189,
	805306457, 1610612741, 3221225473, 4294967291,
}

type hashListNode struct {
	key   int
	value interface{}
	next  *hashListNode
}

type HashList struct {
	slot     []*hashListNode
	slotSize int
	size     int
}

func NewHashList(size int) *HashList {
	hashList := &HashList{}
	hashList.slotSize = hashList.getSlotSize(size)
	hashList.slot = make([]*hashListNode, hashList.slotSize, hashList.slotSize)
	headListNode := make([]hashListNode, hashList.slotSize, hashList.slotSize)
	for i := 0; i != hashList.slotSize; i++ {
		hashList.slot[i] = &headListNode[i]
	}
	hashList.size = 0
	return hashList
}

func (this *HashList) getSlotSize(size int) int {
	for i := 0; i != len(hashPrimeList); i++ {
		if hashPrimeList[i] > size {
			return hashPrimeList[i]
		}
	}
	return hashPrimeList[len(hashPrimeList)-1]
}

func (this *HashList) hash(key int) int {
	return key % this.slotSize
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
	size     int
}

func newHashListArray() *HashListArray {
	hashListArray := &HashListArray{}
	return hashListArray
}

func (this *HashListArray) build(slot []*hashListNode) {
	this.slotSize = len(slot)
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
	hash := key % this.slotSize
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

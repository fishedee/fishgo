package container

type Queue struct {
	head *queueNode
	tail *queueNode
	size int
}

type queueNode struct {
	data interface{}
	next *queueNode
}

func NewQueue() *Queue {
	queue := &Queue{}
	queue.head = &queueNode{
		data: nil,
		next: nil,
	}
	queue.tail = queue.head
	queue.size = 0
	return queue
}

func (this *Queue) Push(data interface{}) {
	node := &queueNode{
		data: data,
		next: nil,
	}
	this.tail.next = node
	this.tail = node
	this.size++
}

func (this *Queue) Pop() interface{} {
	if this.size == 0 {
		return nil
	}
	first := this.head.next
	this.head = first
	this.size--
	if this.size == 0 {
		this.tail = this.head
	}
	return first.data
}

func (this *Queue) Top() interface{} {
	if this.size == 0 {
		return nil
	}
	return this.head.next.data
}

func (this *Queue) Len() int {
	return this.size
}

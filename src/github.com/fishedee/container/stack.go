package container

type Stack struct {
	head *stackNode
	tail *stackNode
	size int
}

type stackNode struct {
	data interface{}
	next *stackNode
}

func NewStack() *Stack {
	stack := &Stack{}
	stack.head = &stackNode{
		data: nil,
		next: nil,
	}
	stack.size = 0
	return stack
}

func (this *Stack) Push(data interface{}) {
	node := &stackNode{
		data: data,
		next: nil,
	}
	node.next = this.head.next
	this.head.next = node
	this.size++
}

func (this *Stack) Pop() interface{} {
	if this.size == 0 {
		return nil
	}
	first := this.head.next
	second := first.next
	this.head.next = second
	this.size--
	return first.data
}

func (this *Stack) Top() interface{} {
	if this.size == 0 {
		return nil
	}
	return this.head.next.data
}

func (this *Stack) Len() int {
	return this.size
}

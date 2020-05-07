package web

import (
	"sync"
	"time"
)

type OldestStayContainer struct {
	stayMap sync.Map
}

type OldestStayElem struct {
	Timestamp int64
	Key       int64
	Value     interface{}
}

func NewOldestStayContainer() *OldestStayContainer {
	return &OldestStayContainer{
		stayMap: sync.Map{},
	}
}

func (this *OldestStayContainer) Push(key int64, value interface{}) {
	this.stayMap.Store(key, &OldestStayElem{
		Timestamp: time.Now().UnixNano(),
		Key:       key,
		Value:     value,
	})
}

func (this *OldestStayContainer) Pop(key int64) {
	this.stayMap.Delete(key)
}

func (this *OldestStayContainer) OldestStay(topSize int) []OldestStayElem {
	heap := make([]OldestStayElem, topSize+1, topSize+1)
	heapSize := 0

	this.stayMap.Range(func(key interface{}, value interface{}) bool {
		elem := value.(*OldestStayElem)
		heapSize = this.pushHeap(heap, heapSize, topSize, elem)
		return true
	})

	//堆排序
	var temp OldestStayElem
	oldHeapSize := heapSize
	for heapSize > 1 {
		temp = heap[1]
		heap[1] = heap[heapSize]
		heap[heapSize] = temp
		heapSize--

		this.adjustDownHeap(heap, heapSize)
	}
	return heap[1 : oldHeapSize+1]
}

func (this *OldestStayContainer) adjustUpHeap(heap []OldestStayElem, heapSize int) {
	var temp OldestStayElem
	for i := heapSize; i/2 >= 1; i = i / 2 {
		parent := i / 2
		if heap[i].Timestamp > heap[parent].Timestamp {
			temp = heap[i]
			heap[i] = heap[parent]
			heap[parent] = temp
		}
	}
}

func (this *OldestStayContainer) adjustDownHeap(heap []OldestStayElem, heapSize int) {
	var temp OldestStayElem
	i := 1
	for {
		maxest := i
		if i*2 <= heapSize && heap[i*2].Timestamp > heap[i].Timestamp {
			maxest = i * 2
		}
		if i*2+1 <= heapSize && heap[i*2+1].Timestamp > heap[maxest].Timestamp {
			maxest = i*2 + 1
		}
		if maxest == i {
			break
		}
		temp = heap[i]
		heap[i] = heap[maxest]
		heap[maxest] = temp
		i = maxest
	}
}

func (this *OldestStayContainer) pushHeap(heap []OldestStayElem, heapSize int, topSize int, elem *OldestStayElem) int {
	if heapSize < topSize {
		//堆未满，向上调整
		heap[heapSize+1] = *elem
		heapSize++
		this.adjustUpHeap(heap, heapSize)
	} else {
		//堆已满，向下调整
		if heap[1].Timestamp > elem.Timestamp {
			heap[1] = *elem
			this.adjustDownHeap(heap, heapSize)
		}
	}
	return heapSize
}

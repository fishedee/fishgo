package web

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type OldestStayContainer struct {
	topSize int
	stayMap sync.Map
}

type OldestStayElem struct {
	Timestamp int64
	Key       int64
	Value     interface{}
}

func NewOldestStayContainer(size int) *OldestStayContainer {
	return &OldestStayContainer{
		topSize: size,
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

func (this *OldestStayContainer) OldestStay() []OldestStayElem {
	heap := make([]OldestStayElem, this.topSize+1, this.topSize+1)
	heapSize := 0

	this.stayMap.Range(func(key interface{}, value interface{}) bool {
		elem := value.(*OldestStayElem)
		heapSize = this.pushHeap(heap, heapSize, elem)
		fmt.Println(heap[1 : heapSize+1])
		return true
	})

	result := heap[1 : heapSize+1]
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})
	return result
}

func (this *OldestStayContainer) pushHeap(heap []OldestStayElem, heapSize int, elem *OldestStayElem) int {
	var temp OldestStayElem
	if heapSize < this.topSize {
		//堆未满，向上调整
		heap[heapSize+1] = *elem
		heapSize++
		for i := heapSize; i/2 >= 1; i = i / 2 {
			parent := i / 2
			if heap[i].Timestamp > heap[parent].Timestamp {
				temp = heap[i]
				heap[i] = heap[parent]
				heap[parent] = temp
			}
		}
	} else {
		//堆已满，向下调整
		if heap[1].Timestamp > elem.Timestamp {
			heap[1] = *elem
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
			}
		}
	}
	return heapSize
}

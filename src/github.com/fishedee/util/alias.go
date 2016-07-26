package util

import (
	"fmt"
	"math/rand"
	"time"
)

// 待优化
func Initialization(probabilities []float64) int {
	count := len(probabilities)
	fmt.Printf("%v", count)

	_alias := make([]int, count)
	_probability := make([]float64, count)

	average := float64(1.0) / float64(count)

	small := NewStack()
	large := NewStack()

	for i := 0; i < count; i++ {
		if probabilities[i] >= average {
			large.Push(i)
		} else {
			small.Push(i)
		}

	}

	for {
		if small.Len() <= 0 || large.Len() <= 0 {
			break
		}

		less := small.Pop().(int)
		more := large.Pop().(int)

		_probability[less] = probabilities[less] * float64(count)
		_alias[less] = more

		probabilities[more] = (probabilities[more] + probabilities[less] - average)

		if probabilities[more] >= average {
			large.Push(more)
		} else {
			small.Push(more)
		}
	}

	for {
		if small.Len() <= 0 {
			break
		}

		_probability[small.Pop().(int)] = 1.0
	}

	for {
		if large.Len() <= 0 {
			break
		}

		_probability[large.Pop().(int)] = 1.0
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	num := r.Intn(count)

	result := rand.Float64() < _probability[num]

	if result {
		return num
	} else {
		return _alias[num]
	}
}

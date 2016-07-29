// 别名算法
// 把 N 种可能性拼装成一个方形（整体），分成 N 列，每列高度为 1 且最多 2 种可能性
// 可能性抽象为某种颜色，即每列最多有 2 种颜色，且第 n 列中必有第 n 种可能性，这里将第 n 种可能性称为原色
// 两个数组：
// 		一个记录落在原色的概率是多少，记为 Prob 数组;
// 		另一个记录列上非原色的颜色名称，记为 Alias 数组，若该列只有原色则记为 null
package util

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

type AliasMethod struct {
	probabilities []float64
	prob          []float64
	alias         []int
}

func NewAliasMethod(probabilities []float64) *AliasMethod {
	alias := &AliasMethod{
		probabilities: probabilities,
	}
	alias.initialization()
	return alias
}

func (this *AliasMethod) initialization() {
	// 设置浮点数的精度
	decimal.DivisionPrecision = 4

	// 输入数据校验
	sum := decimal.NewFromFloat(0.0)
	for _, single := range this.probabilities {
		if single < 0 {
			panic(errors.New("传入概率不能小于0,你输入了:" + fmt.Sprintf("%+v", single)))
		}
		singleDecimal := decimal.NewFromFloat(single)
		sum = sum.Add(singleDecimal)
	}
	oneFloat := decimal.NewFromFloat(1.0)
	if !sum.Equals(oneFloat) {
		panic(errors.New("传入概率数组之和不为1～" + fmt.Sprintf("%+v", this.probabilities)))
	}
	count := len(this.probabilities)

	// 原色数组
	this.prob = make([]float64, count)

	// 别名数组
	this.alias = make([]int, count)
	for i := 0; i < count; i++ {
		this.alias[i] = -1
	}

	// 平均概率
	countFloat := decimal.NewFromFloat(float64(count))
	averageFloat := oneFloat.Div(countFloat)
	average, _ := averageFloat.Float64()

	small := NewStack()
	large := NewStack()

	for i := 0; i < count; i++ {
		if this.probabilities[i] >= average {
			// 大于平均概率推入large栈
			large.Push(i)
		} else {
			// 小于平均概率推入small栈
			small.Push(i)
		}
	}

	// 每次处理一个小概率数
	for {

		if small.Len() <= 0 || large.Len() <= 0 {
			break
		}

		// 小概率下标
		less := small.Pop().(int)

		// 大概率下标
		more := large.Pop().(int)

		probabilitiesLessFloat := decimal.NewFromFloat(this.probabilities[less])
		probabilitiesMoreFloat := decimal.NewFromFloat(this.probabilities[more])

		// 概率值翻count倍--原色数组
		lessFloat := probabilitiesLessFloat.Mul(countFloat)
		this.prob[less], _ = lessFloat.Float64()

		// 大概率数移动部分将小概率补为1,纪录小概率数被谁补偿--别名数组
		this.alias[less] = more

		// 补偿后
		leftFloat := probabilitiesMoreFloat.Add(probabilitiesLessFloat).Sub(averageFloat)
		this.probabilities[more], _ = leftFloat.Float64()

		// 判断剩余部分大小

		if this.probabilities[more] >= average {
			large.Push(more)
		} else {
			small.Push(more)
		}
	}

	// 剩下部分
	for {
		if small.Len() <= 0 {
			break
		}

		this.prob[small.Pop().(int)] = 1.0
	}

	for {
		if large.Len() <= 0 {
			break
		}

		this.prob[large.Pop().(int)] = 1.0
	}
}

func (this *AliasMethod) Rand() int {
	count := len(this.prob)
	if count == 0 {
		return -1
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// count列中，随机选取一列
	num := r.Intn(count)

	// 生成概率，并与原色概率比较
	result := rand.Float64() < this.prob[num]

	// 返回概率列表下标
	if result {
		return num
	} else {
		return this.alias[num]
	}
}

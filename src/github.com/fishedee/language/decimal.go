package language

import (
	"math"
	"strconv"
)

/*
* 使用64位整数来实现精度为0.0001的小数运算
* 为什么不直接用decimal包，因为decimal包的转储需要用string，不方便在数据库和前端进行运算，并且绝大多数的业务仅需保留4位小数就已经足够了
* 详情看[这里](https://fishedee.com/2016/03/16/mysql%E7%BB%8F%E9%AA%8C%E6%B1%87%E6%80%BB/)的浮点数章节
* 本目录下的math.go的实现是错误的，将decimal包直接用float64转储会[爆炸](https://github.com/fishedee/Demo/blob/master/go/decimal/main.go)
 */
type Decimal int

func (left Decimal) Add(right Decimal) Decimal {
	a := int(left)
	b := int(right)
	return Decimal(a + b)
}

func (left Decimal) Sub(right Decimal) Decimal {
	a := int(left)
	b := int(right)
	return Decimal(a - b)
}

func (left Decimal) Mul(right Decimal) Decimal {
	a := int(left)
	b := int(right)
	r := a * b
	exp := r % 10000
	main := r / 10000
	if exp >= 5000 {
		return Decimal(main + 1)
	} else {
		return Decimal(main)
	}
}

func (left Decimal) Div(right Decimal) Decimal {
	a := int(left)
	b := int(right)
	a = a * 100000
	r := a / b
	main := r / 10
	exp := r % 10
	if exp >= 5 {
		return Decimal(main + 1)
	} else {
		return Decimal(main)
	}
}

func (left Decimal) Round(precision int) Decimal {
	a := int(left)
	if precision == 4 {
		return Decimal(a)
	}
	precisionMap := []int{10000, 1000, 100, 10}
	p := precisionMap[precision]
	main := a / p * p
	exp := a % p
	if exp >= (p / 2) {
		return Decimal(main + p)
	} else {
		return Decimal(main)
	}
}

func NewDecimalFromString(a string) (Decimal, error) {
	data, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	result := NewDecimal(data)
	return result, nil
}

func NewDecimal(a float64) Decimal {
	return Decimal(int(math.Round(a * 10000)))
}

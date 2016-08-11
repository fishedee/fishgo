package language

import (
	"errors"

	"github.com/shopspring/decimal"
)

//小数点绝对值
func AbsDecimal(x float64) float64 {
	xDecimal := decimal.NewFromFloat(x)

	xAbsDecimal := xDecimal.Abs()

	return toFloat64(xAbsDecimal)
}

//绝对值转成浮点数
func toFloat64(d decimal.Decimal) float64 {
	result, ok := d.Float64()
	if !ok {
		// TODO
	}
	return result
}

//运算方法-加
func AddDecimal(a, b float64) float64 {
	return operate(a, b, "add")
}

//运算方法-减
func SubDecimal(a, b float64) float64 {
	return operate(a, b, "sub")
}

//运算方法-乘
func MulDecimal(a, b float64) float64 {
	return operate(a, b, "mul")
}

//运算方法-除
func DivDecimal(a, b float64) float64 {
	return operate(a, b, "div")
}

//运算方法-取余数
func ModDecimal(a, b float64) float64 {
	return operate(a, b, "mod")
}

//具体计算
func operate(a, b float64, oper string) float64 {
	aDecimal := decimal.NewFromFloat(a)
	bDecimal := decimal.NewFromFloat(b)

	rDecimal := decimal.Decimal{}
	switch oper {
	case "add":
		rDecimal = aDecimal.Add(bDecimal)
	case "sub":
		rDecimal = aDecimal.Sub(bDecimal)
	case "mul":
		rDecimal = aDecimal.Mul(bDecimal)
	case "div":
		rDecimal = aDecimal.Div(bDecimal)
	case "mod":
		rDecimal = aDecimal.Mod(bDecimal)
	default:
		panic(errors.New("非法运算符！"))
	}
	return toFloat64(rDecimal)
}

//比较
func CmpDecimal(a, b float64) int {
	aDecimal := decimal.NewFromFloat(a)
	bDecimal := decimal.NewFromFloat(b)

	return aDecimal.Cmp(bDecimal)
}

//相等比较
func EqualsDecimal(a, b float64) bool {
	return CmpDecimal(a, b) == 0
}

//输出整数绝对值
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

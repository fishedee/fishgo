package language

import (
	"errors"

	"github.com/shopspring/decimal"
)

func AbsDecimal(x float64) float64 {
	xDecimal := decimal.NewFromFloat(x)

	xAbsDecimal := xDecimal.Abs()

	return toFloat64(xAbsDecimal)
}

func toFloat64(d decimal.Decimal) float64 {
	result, ok := d.Float64()
	if !ok {
		// TODO
	}
	return result
}

func AddDecimal(a, b float64) float64 {
	return operate(a, b, "add")
}

func SubDecimal(a, b float64) float64 {
	return operate(a, b, "sub")
}

func MulDecimal(a, b float64) float64 {
	return operate(a, b, "mul")
}

func DivDecimal(a, b float64) float64 {
	return operate(a, b, "div")
}

func ModDecimal(a, b float64) float64 {
	return operate(a, b, "mod")
}

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

func CmpDecimal(a, b float64) int {
	aDecimal := decimal.NewFromFloat(a)
	bDecimal := decimal.NewFromFloat(b)

	return aDecimal.Cmp(bDecimal)
}

func EqualsDecimal(a, b float64) bool {
	return CmpDecimal(a, b) == 0
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

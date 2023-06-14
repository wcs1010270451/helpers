package _type

import (
	"math"
	"math/big"
)

// SafeFloatAdd 安全的浮点数相加
func SafeFloatAdd(a, b float64) float64 {
	aBF := new(big.Float).SetFloat64(a)
	bBF := new(big.Float).SetFloat64(b)
	sumBF := new(big.Float).Add(aBF, bBF)
	c, _ := sumBF.Float64()
	return c
}

// RoundToTwoDecimalPlaces 保留小数点后两位
func RoundToTwoDecimalPlaces(num float64) float64 {
	return math.Round(num*100) / 100
}

// RoundToSevenDecimalPlaces 保留小数点后7位
func RoundToSevenDecimalPlaces(num float64) float64 {
	return math.Round(num*1e7) / 1e7
}

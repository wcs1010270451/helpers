package _type

import (
	"crypto/rand"
	"io"
	"strings"
)

func InStrArray(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// IsContainStr 判断数组中的字符是否在当前字符串中
func IsContainStr(str string, arr []string) bool {
	for _, s := range arr {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

package main

import (
	"fmt"
	_struct "github.com/wcs1010270451/helpers/struct"
)

/**
 * main
 *  @Description: 调试
 */
func main() {
	a := A{"郭德纲", 20}
	b := B{A: A{}, Address: "北京", Info: "旅游"}

	fmt.Println("Before copying:")
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	_struct.Copy(&a, &b)
	fmt.Println("b", b)
}

type A struct {
	Name string
	Age  int
}

type B struct {
	A
	Address string
	Info    string
}

package main

import (
	"fmt"

	"github.com/billylkc/aastock"
)

func main() {
	getCompanyName()
}

func getCompanyName() {
	code := 11
	company, err := aastock.GetCompanyName(code)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(company)
}

// get delayed price
func getPrice() {
	data := makeRange(1, 10)
	result := aastock.GetCurrentPrices(data...)
	fmt.Println(result)
	fmt.Println(len(result))
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

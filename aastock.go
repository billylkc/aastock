package aastock

import (
	"fmt"
)

// GetList gets a list of stocks that we are interested in
func GetList() []string {
	var result []string
	c := []string{"1", "5", "9988", "700"}

	for _, cc := range c {
		result = append(result, fmt.Sprintf("%05s", cc))
	}
	return result
}

// getLastValue gets the last non zero value from a list
func getLastValue(values []float64) float64 {
	var last float64
	for _, v := range values {
		if v > 0 {
			last = v
		}
	}
	return last
}

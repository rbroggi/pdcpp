package cache

import (
	"math"
	"testing"
)

var EPSILON = 0.00000001

func TestComboMultiplePrices(t *testing.T) {
	m := make(map[string]float64)
	m["a"] = 1.0
	m["b"] = 2.0
	m["c"] = -3.0
	f := func(s string) []float64 {
		switch s {
		case "a":
			return []float64{1.0, 2.0, 3.0}
		case "b":
			return []float64{4.0, 5.0, 6.0}
		case "c":
			return []float64{7.0, 8.0, 9.0}
		default:
			return []float64{0.0, 0.0, 0.0}
		}
	}
	expectedResult := []float64{-12.0, -12.0, -12.0}
	result := comboMultiplePrices(m, f)
	for i := range result {
		if !floatEquals(result[i], expectedResult[i]) {
			t.Errorf("Expected result %v, differs from actual result %v, for index %v", expectedResult[i], result[i], i)
		}
	}
}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

func TestSumSingle(t *testing.T) {

	m := make(map[string]float64)
	m["a"] = 1.0
	m["b"] = 2.0
	m["c"] = -2.0
	f := func(s string) float64 {
		switch s {
		case "a":
			return 10.0
		case "b":
			return 20.0
		case "c":
			return 30.0
		default:
			return 0.
		}
	}
	expectedResult := -10.0
	result := sumSingle(m, f)
	if !floatEquals(result, expectedResult) {
		t.Errorf("Expected result %v, differs from actual result %v", expectedResult, result)
	}

}

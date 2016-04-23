package main

import (
	"testing"
)

func TestMean(t *testing.T) {
	floats := []float64{5.0, 10.0, 3.0}
	m := mean(floats)
	if m != 6.0 {
		t.Error("mean is broken (output is ", m, ")")
	}
}

func TestStDev(t *testing.T) {
	floats :[]float64{9, 2, 5, 4, 12, 7, 8, 11, 9, 3, 7, 4, 12, 5, 4, 10, 9, 6, 9, 4}
	stdev := 
}

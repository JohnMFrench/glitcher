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

package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestMean(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	expected := 3.0
	result := mean(data)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPredictCounterfactuals(t *testing.T) {
	coeff := mat.NewVecDense(3, []float64{1, 2, 3})
	covariate := []float64{1, 2, 3}
	newTreatment := 1.0
	expected := []float64{6, 9, 12}
	result := predictCounterfactuals(*coeff, covariate, newTreatment)
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], v)
		}
	}
}

func BenchmarkMean(b *testing.B) {
	data := []float64{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		mean(data)
	}
}

func BenchmarkPredictCounterfactuals(b *testing.B) {
	coeff := mat.NewVecDense(3, []float64{1, 2, 3})
	covariate := []float64{1, 2, 3}
	newTreatment := 1.0
	for i := 0; i < b.N; i++ {
		predictCounterfactuals(*coeff, covariate, newTreatment)
	}
}

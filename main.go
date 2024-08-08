package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	start := time.Now()

	// Set up logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Starting application")

	// Start memory profiling
	f, err := os.Create("memprofile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Log initial memory usage
	logMemoryUsage("Initial memory usage")

	// Open CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Log memory usage after reading CSV
	logMemoryUsage("Memory usage after reading CSV")

	// Initialize slices to hold the data
	n := len(records) - 1 // Assuming the first row is the header
	id := make([]int, n)
	treatment := make([]float64, n)
	outcome := make([]float64, n)
	covariate := make([]float64, n)

	// Parse the CSV data into slices
	for i, record := range records[1:] {
		id[i], err = strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("Failed to parse ID: %v", err)
		}
		treatment[i], err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse treatment: %v", err)
		}
		outcome[i], err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatalf("Failed to parse outcome: %v", err)
		}
		covariate[i], err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatalf("Failed to parse covariate: %v", err)
		}
	}

	// Log memory usage after parsing CSV
	logMemoryUsage("Memory usage after parsing CSV")

	// Fit a linear model
	X := mat.NewDense(n, 3, nil)
	for i := 0; i < n; i++ {
		X.Set(i, 0, 1) // Intercept
		X.Set(i, 1, treatment[i])
		X.Set(i, 2, covariate[i])
	}
	y := mat.NewVecDense(n, outcome)

	var XtX mat.Dense
	XtX.Mul(X.T(), X)

	var XtXInv mat.Dense
	if err := XtXInv.Inverse(&XtX); err != nil {
		log.Fatalf("Error inverting matrix: %v", err)
	}

	var XtY mat.VecDense
	XtY.MulVec(X.T(), y)

	var coeff mat.VecDense
	coeff.MulVec(&XtXInv, &XtY)

	// Log memory usage after fitting model
	logMemoryUsage("Memory usage after fitting model")

	// Calculate residuals
	var yPred mat.VecDense
	yPred.MulVec(X, &coeff)
	var residuals mat.VecDense
	residuals.SubVec(y, &yPred)

	// Calculate residual standard error
	var residualsSquared []float64
	for i := 0; i < residuals.Len(); i++ {
		residualsSquared = append(residualsSquared, residuals.At(i, 0)*residuals.At(i, 0))
	}
	rss := stat.Variance(residualsSquared, nil) * float64(n)
	residualSE := math.Sqrt(rss / float64(n-3))

	// Calculate R-squared and adjusted R-squared
	yMean := mean(y.RawVector().Data)
	tss := 0.0
	for _, v := range y.RawVector().Data {
		tss += (v - yMean) * (v - yMean)
	}
	rSquared := 1 - rss/tss
	adjRSquared := 1 - (1-rSquared)*(float64(n)/float64(n-3))

	// Calculate F-statistic
	fStatistic := (rSquared / 2) / ((1 - rSquared) / float64(2))
	fDist := F{DFn: 2, DFd: float64(n - 3)}
	pValue := 1 - fDist.CDF(fStatistic)

	// Calculate standard errors of coefficients
	coeffSE := make([]float64, 3)
	for i := 0; i < 3; i++ {
		coeffSE[i] = math.Sqrt(XtXInv.At(i, i)) * residualSE
	}

	// Calculate t-values and p-values for coefficients
	tValues := make([]float64, 3)
	pValues := make([]float64, 3)
	studentsT := distuv.StudentsT{Nu: float64(n - 3)}
	for i := 0; i < 3; i++ {
		tValues[i] = coeff.At(i, 0) / coeffSE[i]
		pValues[i] = 2 * (1 - studentsT.CDF(math.Abs(tValues[i])))
	}

	// Write results to a file
	resultsFile, err := os.Create("results.txt")
	if err != nil {
		log.Fatalf("Error creating results file: %v", err)
	}
	defer resultsFile.Close()

	// Print values
	meanObserved := mean(y.RawVector().Data)
	meanCfTreatment1 := mean(predictCounterfactuals(coeff, covariate, 1))
	meanCfTreatment0 := mean(predictCounterfactuals(coeff, covariate, 0))

	fmt.Fprintf(resultsFile, "\nMean Observed Outcome: %.5f\n", meanObserved)
	fmt.Fprintf(resultsFile, "Mean Counterfactual Outcome (Treatment = 1): %.5f\n", meanCfTreatment1)
	fmt.Fprintf(resultsFile, "Mean Counterfactual Outcome (Treatment = 0): %.5f\n", meanCfTreatment0)

	fmt.Fprintf(resultsFile, "Coefficients:\n")
	fmt.Fprintf(resultsFile, "             Estimate Std. Error t value Pr(>|t|)\n")
	fmt.Fprintf(resultsFile, "(Intercept)  %8.4f %10.4f %7.4f %8.4f\n", coeff.At(0, 0), coeffSE[0], tValues[0], pValues[0])
	fmt.Fprintf(resultsFile, "treatment    %8.4f %10.4f %7.4f %8.4f\n", coeff.At(1, 0), coeffSE[1], tValues[1], pValues[1])
	fmt.Fprintf(resultsFile, "covariate    %8.4f %10.4f %7.4f %8.4f\n", coeff.At(2, 0), coeffSE[2], tValues[2], pValues[2])
	fmt.Fprintf(resultsFile, "\nResidual standard error: %.4f on %d degrees of freedom\n", residualSE, n-3)
	fmt.Fprintf(resultsFile, "Multiple R-squared: %.6f, Adjusted R-squared: %.6f\n", rSquared, adjRSquared)
	fmt.Fprintf(resultsFile, "F-statistic: %.4f on 2 and %d DF,  p-value: %.5f\n", fStatistic, n-3, pValue)

	// Finish profiling
	pprof.WriteHeapProfile(f)
	log.Println("Application finished")
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)

	// Log final memory usage
	logMemoryUsage("Final memory usage")
}

// All Helper functions:
func mean(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func predictCounterfactuals(coeff mat.VecDense, covariate []float64, newTreatment float64) []float64 {
	n := len(covariate)
	predictions := make([]float64, n)
	for i := 0; i < n; i++ {
		predictions[i] = coeff.AtVec(0) + coeff.AtVec(1)*newTreatment + coeff.AtVec(2)*covariate[i]
	}
	return predictions
}

func logMemoryUsage(message string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logMessage := fmt.Sprintf("%s: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
		message, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	log.Print(logMessage)
	fmt.Print(logMessage)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Custom F distribution type
type F struct {
	DFn float64 // Degrees of freedom numerator
	DFd float64 // Degrees of freedom denominator
}

// CDF method for the custom F distribution
func (f F) CDF(x float64) float64 {
	return distuv.Beta{Alpha: f.DFn / 2, Beta: f.DFd / 2}.CDF(f.DFn * x / (f.DFn*x + f.DFd))
}

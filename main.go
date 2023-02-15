package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

func main() {
	// Load the dataset from a CSV file
	file, err := os.Open("energy_loss3.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Extract the input and output data from the dataset
	var (
		x             [][]float64
		kwhPenerimaan []float64
	)

	for _, record := range records {
		var input []float64
		for i := 0; i < len(record)-3; i++ {
			val, err := strconv.ParseFloat(record[i], 64)
			if err != nil {
				log.Fatal(err)
			}
			input = append(input, val)
		}
		x = append(x, input)

		kwhPenerimaanOut, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		kwhPenerimaan = append(kwhPenerimaan, kwhPenerimaanOut)

		// y2Out, err := strconv.ParseFloat(record[3], 64)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// y2 = append(y2, y2Out)

		// y3Out, err := strconv.ParseFloat(record[4], 64)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// y3 = append(y3, y3Out)
	}

	// Train the linear regression model
	r := new(regression.Regression)
	r.SetObserved("Energy loss every year")
	r.SetVar(0, "month")
	r.SetVar(1, "year")

	for i, xi := range x {
		r.Train(
			regression.DataPoint(kwhPenerimaan[i], xi),
			// regression.DataPoint(y2[i], xi),
			// regression.DataPoint(y3[i], xi),
		)
	}

	r.Run()

	// fmt.Printf("Regression formula:\n%v\n", r.Formula)
	// fmt.Printf("Regression:\n%s\n", r)

	// Use the trained model to make predictions for the next 5 years
	var predictions []float64
	for year := 2019; year <= 2023; year++ {
		for month := 1; month <= 12; month++ {
			input := []float64{float64(month), float64(year)}
			prediction, err := r.Predict(input)
			if err != nil {
				log.Fatal(err)
			}
			predictions = append(predictions, prediction)
		}
	}

	// Calculate the mean absolute percentage error of the predictions
	var actual []float64
	var absoluteErrors []float64
	actual = append(actual, kwhPenerimaan...)

	for i, value := range predictions {
		absoluteError := 100 * (actual[i] - value) / actual[i]
		absoluteErrors = append(absoluteErrors, absoluteError)
	}

	sum := 0.0
	for _, error := range absoluteErrors {
		sum += error
	}
	meanAbsolutePercentageError := sum / float64(len(absoluteErrors))

	// Print the predictions and the mean absolute percentage error
	fmt.Println("Predictions:")
	for i, prediction := range predictions {
		fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
	}

	fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)
}

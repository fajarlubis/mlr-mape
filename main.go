package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sajari/regression"
	"github.com/xuri/excelize/v2"
)

type Result struct {
	Items []*ResultItem
	Type  string
	MAPE  float64
}

type ResultItem struct {
	Year       int
	Month      int
	Prediction float64
}

func main() {
	// Load the dataset from a CSV file
	file, err := os.Open("training-data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var (
		res   = make([]*Result, 0)
		items = make([]*ResultItem, 0)
	)

	for label := 0; label <= 10; label++ { // there is 11 data label
		// Extract the input and output data from the dataset
		var (
			x       [][]float64
			y       []float64
			kWhType string
		)

		for _, record := range records {
			var input []float64
			for i := 0; i < len(record)-11; i++ {
				val, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					log.Fatal(err)
				}
				input = append(input, val)
			}
			x = append(x, input)

			yOut, err := strconv.ParseFloat(record[label+2], 64)
			if err != nil {
				log.Fatal(err)
			}
			y = append(y, yOut)

			switch label + 2 {
			case 2:
				kWhType = "kWh Penerimaan"
			case 3:
				kWhType = "kWh Penjualan"
			case 4:
				kWhType = "Pemakaian Sendiri"
			case 5:
				kWhType = "Kirim ke Unit Lain"
			case 6:
				kWhType = "Susut Teknis JTM"
			case 7:
				kWhType = "Susut Teknis JTR"
			case 8:
				kWhType = "Susut Teknis Trafo"
			case 9:
				kWhType = "Susut Teknis SR"
			case 10:
				kWhType = "Susut Total"
			case 11:
				kWhType = "Susut Teknis"
			case 12:
				kWhType = "Susut Non-Teknis"
			}
		}

		// Train the linear regression model
		r := new(regression.Regression)
		r.SetObserved("Energy loss every year")
		r.SetVar(0, "month")
		r.SetVar(1, "year")

		for i, xi := range x {
			r.Train(
				regression.DataPoint(y[i], xi),
			)
		}

		r.Run()

		// fmt.Printf("Regression formula:\n%v\n", r.Formula)
		// fmt.Printf("Regression:\n%s\n", r)

		fmt.Printf("R^2:\n%v\n", r.R2)

		// Use the trained model to make predictions for the next 5 years
		var predictions []float64
		for year := 2021; year <= 2025; year++ {
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
		actual = append(actual, y...)

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
		// fmt.Println("Predictions:")
		for i, prediction := range predictions {
			items = append(items, &ResultItem{
				Year:       (i / 12) + 2021,
				Month:      (i % 12) + 1,
				Prediction: prediction,
			})

			// fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
		}

		// fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)

		res = append(res, &Result{
			Items: items,
			Type:  kWhType,
			MAPE:  meanAbsolutePercentageError,
		})
	}

	// for _, v := range res {
	// 	for _, w := range v.Items {
	// 		fmt.Printf("%d-%02d: %.2f\n", w.Year, w.Month, w.Prediction)
	// 	}
	// 	fmt.Printf("Mean absolute percentage error: %.2f%%\n", v.MAPE)
	// }

	if err := export2(res); err != nil {
		log.Fatalln(err)
	}
}

func export2(data []*Result) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet = "Sheet1"
	)

	// create a new sheet
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	var kWhType string
	rowStartIdx := 1
	for _, res := range data {
		kWhType = res.Type

		for _, v := range res.Items {
			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", rowStartIdx), time.Month(v.Month))
			f.SetCellInt(defaultSheet, fmt.Sprintf("B%v", rowStartIdx), v.Year)
			f.SetCellFloat(defaultSheet, fmt.Sprintf("C%v", rowStartIdx), v.Prediction, 2, 64)
			f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx), kWhType)

			rowStartIdx += 1
		}

		f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", rowStartIdx), fmt.Sprintf("MAPE: %v", res.MAPE))
		rowStartIdx += 1
	}

	f.SetActiveSheet(index)
	// save spreadsheet by the given path
	if err := f.SaveAs("result.xlsx"); err != nil {
		return err
	}

	return nil
}

func export(res []*Result) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet = "Sheet1"
	)

	// create a new sheet
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	// START CELL STYLE =========================================================================================== //

	centeredStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	boldYellowBgStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#ffff00"}, Pattern: 1},
	})

	// boldBlueBgStyle, _ := f.NewStyle(&excelize.Style{
	// 	Font: &excelize.Font{
	// 		Bold:  true,
	// 		Color: "#ffffff",
	// 	},
	// 	Fill: excelize.Fill{Type: "pattern", Color: []string{"#1135d4"}, Pattern: 1},
	// })

	boldCenteredStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	customNumFmt := "#,##0.##"
	thousandSeparatorNumberStyle, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &customNumFmt})

	// END CELL STYLE =========================================================================================== //

	// START DEFINING ROW/COL SIZE =========================================================================================== //

	f.SetColWidth(defaultSheet, "A", "A", 5)
	f.SetColWidth(defaultSheet, "B", "B", 20)
	f.SetColWidth(defaultSheet, "D", "O", 20)

	// END DEFINING ROW/COL SIZE =========================================================================================== //

	// START PRINTING PREDICTIONS =========================================================================================== //

	rowStartIdx := 3 // starting point
	for i, v := range res {
		f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("FORJA %v", v.MAPE))
		f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), "Bulan")

		f.MergeCell(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+1+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("A%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+1+i), fmt.Sprintf("B%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("C%v", rowStartIdx+1+i), fmt.Sprintf("C%v", rowStartIdx+2+i))

		f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+2+i), boldCenteredStyle)
		f.SetCellStyle(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("B%v", rowStartIdx+i), boldYellowBgStyle)

		no := 1
		valueStartRowIdx := rowStartIdx + 3
		for j, w := range v.Items {

			monthNameStartColIdx := 4 // D
			monthNameStartRowIdx := rowStartIdx + 2

			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", monthNameStartRowIdx+i), "No.")
			f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", monthNameStartRowIdx+i), "Input")
			f.SetCellValue(defaultSheet, fmt.Sprintf("C%v", monthNameStartRowIdx+i), "Satuan")

			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", valueStartRowIdx+i), no)
			f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", valueStartRowIdx+i), j)
			f.SetCellValue(defaultSheet, fmt.Sprintf("C%v", valueStartRowIdx+i), "kWh")

			f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", valueStartRowIdx+i), fmt.Sprintf("A%v", valueStartRowIdx+i), centeredStyle)
			f.SetCellStyle(defaultSheet, fmt.Sprintf("C%v", valueStartRowIdx+i), fmt.Sprintf("C%v", valueStartRowIdx+i), centeredStyle)

			for m := time.January; m <= time.December; m++ {
				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), monthNameStartRowIdx+i), m.String())
				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), math.Round(w.Prediction))

				f.SetCellStyle(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), thousandSeparatorNumberStyle)

				monthNameStartColIdx += 1
			}

			valueStartRowIdx += 1

			no += 1

		}

		rowStartIdx += 14 // 14 is total length of populated rows for one year
	}

	// END PRINTING PREDICTIONS =========================================================================================== //

	f.SetActiveSheet(index)
	// save spreadsheet by the given path
	if err := f.SaveAs("result.xlsx"); err != nil {
		return err
	}

	return nil
}

func toChar(i int) string {
	return strings.Replace(string(rune('A'-1+i)), "'", "", -1)
}

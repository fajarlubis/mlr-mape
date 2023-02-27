package printout

import (
	"fmt"
	"math"
	"mlr-mape/model"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func Export(dataset, predictions []model.YearlyData) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet = "Sheet1"
		chartsSheet  = "Charts"
	)

	// create a new sheet
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	// create a new sheet for charts
	_, err = f.NewSheet(chartsSheet)
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

	boldBlueBgStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#ffffff",
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#1135d4"}, Pattern: 1},
	})

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

	// START PRINTING DEFAULT DATASET =========================================================================================== //
	rowStartIdx := 3 // starting point

	if dataset != nil {
		sort.Slice(dataset, func(i, j int) bool {
			return dataset[i].Year < dataset[j].Year
		})

		for i, v := range dataset {
			f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("FORJA %v", v.Year))
			f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), "Bulan")

			f.MergeCell(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+1+i))
			f.MergeCell(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("A%v", rowStartIdx+2+i))
			f.MergeCell(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+1+i), fmt.Sprintf("B%v", rowStartIdx+2+i))
			f.MergeCell(defaultSheet, fmt.Sprintf("C%v", rowStartIdx+1+i), fmt.Sprintf("C%v", rowStartIdx+2+i))

			f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+2+i), boldCenteredStyle)
			f.SetCellStyle(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("B%v", rowStartIdx+i), boldYellowBgStyle)

			keys := make([]string, 0, len(v.Data))
			for key := range v.Data {
				keys = append(keys, key)
			}

			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			no := 1
			valueStartRowIdx := rowStartIdx + 3
			for _, j := range keys {
				w := v.Data[j]

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
					value, ok := w[m]
					if !ok {
						continue
					}

					f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), monthNameStartRowIdx+i), m.String())
					f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), math.Round(value))

					f.SetCellStyle(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), thousandSeparatorNumberStyle)

					monthNameStartColIdx += 1
				}

				valueStartRowIdx += 1

				no += 1
			}

			rowStartIdx += 14 // 14 is total length of populated rows for one year
		}
	}

	// END PRINTING DEFAULT DATASET =========================================================================================== //

	// START PRINTING PREDICTIONS =========================================================================================== //

	if dataset != nil {
		rowStartIdx = 78 // starting point
	}

	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Year < predictions[j].Year
	})

	for i, v := range predictions {
		f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("FORJA %v", v.Year))
		f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), "Bulan")

		f.MergeCell(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+1+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("A%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+1+i), fmt.Sprintf("B%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("C%v", rowStartIdx+1+i), fmt.Sprintf("C%v", rowStartIdx+2+i))

		f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+2+i), boldCenteredStyle)
		f.SetCellStyle(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("B%v", rowStartIdx+i), boldBlueBgStyle)

		keys := make([]string, 0, len(v.Data))
		for key := range v.Data {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		no := 1
		valueStartRowIdx := rowStartIdx + 3
		for _, j := range keys {
			w := v.Data[j]

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
				value, ok := w[m]
				if !ok {
					continue
				}

				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), monthNameStartRowIdx+i), m.String())
				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), math.Round(value))

				f.SetCellStyle(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), thousandSeparatorNumberStyle)

				monthNameStartColIdx += 1
			}

			valueStartRowIdx += 1

			no += 1
		}

		rowStartIdx += 14 // 14 is total length of populated rows for one year
	}

	// END PRINTING PREDICTIONS =========================================================================================== //

	// for idx, row := range [][]interface{}{
	// 	{nil, "Apple", "Orange", "Pear"},
	// 	{"Small", 2, 3, 3},
	// 	{"Normal", 5, 2, 4},
	// 	{"Large", 6, 7, 8},
	// } {
	// 	cell, err := excelize.CoordinatesToCellName(1, idx+1)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if err := f.SetSheetRow(chartsSheet, cell, &row); err != nil {
	// 		return err
	// 	}
	// }
	// if err := f.AddChart(chartsSheet, "E1", &excelize.Chart{
	// 	Type: "line",
	// 	Series: []excelize.ChartSeries{
	// 		{
	// 			Name:       fmt.Sprintf("%s!$A$2", chartsSheet),
	// 			Categories: fmt.Sprintf("%s!$B$1:$D$1", chartsSheet),
	// 			Values:     fmt.Sprintf("%s!$B$2:$D$2", chartsSheet),
	// 			Line: excelize.ChartLine{
	// 				Smooth: true,
	// 			},
	// 		},
	// 		{
	// 			Name:       fmt.Sprintf("%s!$A$3", chartsSheet),
	// 			Categories: fmt.Sprintf("%s!$B$1:$D$1", chartsSheet),
	// 			Values:     fmt.Sprintf("%s!$B$3:$D$3", chartsSheet),
	// 			Line: excelize.ChartLine{
	// 				Smooth: true,
	// 			},
	// 		},
	// 		{
	// 			Name:       fmt.Sprintf("%s!$A$4", chartsSheet),
	// 			Categories: fmt.Sprintf("%s!$B$1:$D$1", chartsSheet),
	// 			Values:     fmt.Sprintf("%s!$B$4:$D$4", chartsSheet),
	// 			Line: excelize.ChartLine{
	// 				Smooth: true,
	// 			},
	// 		},
	// 	},
	// 	Format: excelize.GraphicOptions{
	// 		OffsetX: 15,
	// 		OffsetY: 10,
	// 	},
	// 	Legend: excelize.ChartLegend{
	// 		Position: "top",
	// 	},
	// 	Title: excelize.ChartTitle{
	// 		Name: "Fruit Line Chart",
	// 	},
	// 	PlotArea: excelize.ChartPlotArea{
	// 		ShowCatName:     false,
	// 		ShowLeaderLines: false,
	// 		ShowPercent:     true,
	// 		ShowSerName:     true,
	// 		ShowVal:         true,
	// 	},
	// 	ShowBlanksAs: "zero",
	// }); err != nil {
	// 	return err
	// }

	// set active sheet of the workbook
	f.SetActiveSheet(index)
	// save spreadsheet by the given path
	// if err := f.SaveAs(fmt.Sprintf(model.OutputFile, dataset[0].Year, predictions[len(predictions)-1].Year)); err != nil {
	if err := f.SaveAs("result.xlsx"); err != nil {
		return err
	}

	return nil
}

func toChar(i int) string {
	return strings.Replace(string(rune('A'-1+i)), "'", "", -1)
}

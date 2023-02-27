package printout

import (
	"fmt"
	"mlr-mape/model"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportLegacy(data []*model.Result) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet   = "Sheet1"
		formattedSheet = "Formatted"
		chartSheet     = "Charts"
	)

	// sheet default
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	// sheet chart
	_, err = f.NewSheet(chartSheet)
	if err != nil {
		return err
	}

	// sheet formatted
	_, err = f.NewSheet(formattedSheet)
	if err != nil {
		return err
	}

	for _, res := range data {
		resultData := Parse(res.Items, res.Type, res.MAPE)

		rowStartIndex := 1
		for year, _ := range resultData.Items {
			// fmt.Println(year)

			f.SetCellValue(formattedSheet, fmt.Sprintf("A%v", rowStartIndex), year)

			// rrr := 1
			// for month, value := range v {
			// 	f.SetCellValue(formattedSheet, fmt.Sprintf("A%v", rrr), resultData.Type)
			// 	fmt.Println(time.Month(month), value)

			// 	rrr += 1
			// }
			// rrr += 13

			rowStartIndex += 13
		}
	}

	var kWhType string
	rowStartIdx := 1
	for _, res := range data {
		kWhType = res.Type

		// newResult := Parse(res.Items, res.Type, res.MAPE)
		// log.Println(newResult.Type, newResult.MAPE)
		// for year, v := range newResult.Items {
		// 	fmt.Println(year, "--------------------------------------------------->>>>>>")
		// 	for j, w := range v {
		// 		fmt.Println(time.Month(j), w)
		// 	}
		// }
		// log.Println("===============================================================================================")

		// if err := f.AddChart(chartSheet, fmt.Sprintf("A%v", rowStartIdx), printkWhChart(defaultSheet, newResult)); err != nil {
		// 	return err
		// }

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

func printkWhChart(source string, res *Result) *excelize.Chart {
	lineChart := &excelize.Chart{
		Type: "line",
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("%s!$B$1", source),
				Categories: fmt.Sprintf("%s!$A$1:$A$12", source),
				Values:     fmt.Sprintf("%s!$C$1:$C$12", source),
				Line: excelize.ChartLine{
					Smooth: true,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$B$13", source),
				Categories: fmt.Sprintf("%s!$A$1:$A$12", source),
				Values:     fmt.Sprintf("%s!$C$13:$C$23", source),
				Line: excelize.ChartLine{
					Smooth: true,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$B$25", source),
				Categories: fmt.Sprintf("%s!$A$1:$A$12", source),
				Values:     fmt.Sprintf("%s!$C$25:$C$36", source),
				Line: excelize.ChartLine{
					Smooth: true,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$B$37", source),
				Categories: fmt.Sprintf("%s!$A$1:$A$12", source),
				Values:     fmt.Sprintf("%s!$C$37:$C$48", source),
				Line: excelize.ChartLine{
					Smooth: true,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$B$49", source),
				Categories: fmt.Sprintf("%s!$A$1:$A$12", source),
				Values:     fmt.Sprintf("%s!$C$49:$C$60", source),
				Line: excelize.ChartLine{
					Smooth: true,
				},
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
		Title: excelize.ChartTitle{
			Name: res.Type,
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     false,
			ShowSerName:     false,
			ShowVal:         true,
		},
		ShowBlanksAs: "zero",
	}

	return lineChart
}

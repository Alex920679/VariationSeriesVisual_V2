package main

import (
	"fmt"
	v1charts "github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
)

func drawBarChart(uniqueSlice []int, relValFreqSlice []float64) error { // рисовальщик гистограммы
	lenUniqueSlice := len(uniqueSlice)
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Гистограмма данных",
	}))
	bar.SetXAxis(uniqueSlice).
		AddSeries("Category A", generateBarItems(lenUniqueSlice, relValFreqSlice))
	f, err := os.Create("Гистограмма_данных.html")
	if err != nil {
		return fmt.Errorf("creation error")
	}
	bar.Render(f)
	return nil
}

func drawLine(uniqueSlice []int, relValFreqSlice []float64) *v1charts.Line { // рисовальщик полигона распределения
	line := v1charts.NewLine()
	line.SetGlobalOptions(v1charts.TitleOpts{Title: "Полигон распределения"})
	line.AddXAxis(uniqueSlice).AddYAxis("ОЧВ = f(В)", relValFreqSlice)
	return line
}

func renderLine(line *v1charts.Line) error {
	f, err := os.Create("Полигон распределения.html")
	if err != nil {
		return fmt.Errorf("creation error")
	}
	line.Render(f)
	return nil
}

func createExcelTable(uniqueSlice []int, mp map[int]int, relValFreqSlice []float64) error {
	f := excelize.NewFile()
	digitOfCell := 2
	// прописываю стиль заголовков
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: "Times New Roman",
			Size:   12,
			Color:  "#777777",
		},
	})
	if err != nil {
		return fmt.Errorf("filling error")
	} // прописал, далее заполняю ячейки
	f.SetCellValue("Sheet1", "B2", "Вариант")
	f.SetCellValue("Sheet1", "C2", "Частота")
	f.SetCellValue("Sheet1", "D2", "ОЧВ")
	// наделяю заголовки стилем style
	err = f.SetCellStyle("Sheet1", "B2", "B2", style)
	err = f.SetCellStyle("Sheet1", "C2", "C2", style)
	err = f.SetCellStyle("Sheet1", "D2", "D2", style)
	// далее заполняю таблицу
	for i := 0; i < len(uniqueSlice); i++ {
		digitOfCell++
		d := strconv.Itoa(digitOfCell)
		s1 := "B" + d
		f.SetCellValue("Sheet1", s1, uniqueSlice[i])
		s2 := "C" + d
		f.SetCellValue("Sheet1", s2, mp[uniqueSlice[i]])
		s3 := "D" + d
		f.SetCellValue("Sheet1", s3, relValFreqSlice[i])
	}
	if err := f.SaveAs("Таблица данных.xlsx"); err != nil {
		return fmt.Errorf("saving error")
	}
	return nil
}

func generateBarItems(lenUniqueSl int, relValFreqSlice []float64) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < lenUniqueSl; i++ {
		items = append(items, opts.BarData{Value: relValFreqSlice[i]})
	}
	return items
}

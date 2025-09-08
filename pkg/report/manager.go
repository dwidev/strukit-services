package report

import (
	"fmt"
	"strukit-services/internal/models"

	"github.com/xuri/excelize/v2"
)

func Manager() *ReportManager {
	return &ReportManager{}
}

type ReportManager struct {
}

func (r *ReportManager) GenerateExcel(project *models.Project, receipts []*models.Receipt) (*excelize.File, error) {
	f := excelize.NewFile()
	styleBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})

	sheetName := fmt.Sprintf("Laporan keuangan %s", project.Name)
	sheet, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Headers
	headers := [][]string{
		{fmt.Sprintf("LAPORAN KEUANGAN - %s", project.Name)}, // A1
		{""},
		{"Owner", "Fahmi dwi"},
		{"Nama Projek", "Test Project"},
		{"Budget", "Rp. 100.000"},
		{"Total Pengeluaran", "Rp. 100.000"},
		{"Sisa Budget", "Rp. 100.000"},
		{"Persentase Terpakai", "90%"},
		{"Status Projek", "ANjAY"},
		{""},
	}
	headerRow := 1
	for _, header := range headers {
		for col, value := range header {
			cell, _ := excelize.CoordinatesToCellName(col+2, headerRow)
			f.SetCellValue(sheetName, cell, value)

			if col == 0 {
				f.SetCellStyle(sheetName, cell, cell, styleBold)
			}
		}
		headerRow++
	}

	// Contents
	headerContentRow := headerRow + 1
	headerContents := []string{
		"No", "Hari/Tanggal", "Nama Toko", "Nama Item", "Banyaknya", "Total item", "Total bayar",
	}
	for col, contentHeader := range headerContents {
		cell, _ := excelize.CoordinatesToCellName(col+1, headerContentRow)
		f.SetCellValue(sheetName, cell, contentHeader)
	}

	bodyContentRow := headerContentRow + 1
	for col, content := range receipts {
		cell, _ := excelize.CoordinatesToCellName(col+1, bodyContentRow)
		for _, item := range content.Items {
			f.SetCellValue(sheetName, cell, item.ItemName)
		}
		headerContentRow++
	}

	f.SetActiveSheet(sheet)
	return f, nil
}

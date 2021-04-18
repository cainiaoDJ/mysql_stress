package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

const defalut_sheet = "RAW_DATA"

type RowData struct {
	Func    string
	Routine uint
	DBNum   uint
	TbNum   uint
	Cost    float64
	Speed   float64
}

func WriteToExcel(data []RowData, fileName string) error {
	xls := excelize.NewFile()
	sheet := xls.NewSheet(defalut_sheet)
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 1), "Func")
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 2), "Routine")
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 3), "DBNum")
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 4), "TbNum")
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 5), "Cost")
	xls.SetCellStr(defalut_sheet, getExcelKey(1, 6), "Speed")
	xls.SetActiveSheet(sheet)
	for k1, rd := range data {
		xls.SetCellStr(defalut_sheet, getExcelKey(k1+2, 1), rd.Func)
		xls.SetCellInt(defalut_sheet, getExcelKey(k1+2, 2), int(rd.Routine))
		xls.SetCellInt(defalut_sheet, getExcelKey(k1+2, 3), int(rd.DBNum))
		xls.SetCellInt(defalut_sheet, getExcelKey(k1+2, 4), int(rd.TbNum))
		xls.SetCellValue(defalut_sheet, getExcelKey(k1+2, 5), rd.Cost)
		xls.SetCellValue(defalut_sheet, getExcelKey(k1+2, 6), rd.Speed)

	}
	tableNmae := fmt.Sprintf("./%s.xlsx", fileName)
	err := xls.SaveAs(tableNmae)
	if err != nil {
		return err
	}
	return nil
}

func getExcelKey(height int, width int) string {
	code := [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	tmp := make([]byte, 0)

	var index int
	for width > 0 {
		index = width % 26
		tmp = append(tmp, code[index-1])
		width = width / 26
	}
	l := len(tmp)
	for i := 0; i < l/2; i++ {
		tmp[i], tmp[l-1-i] = tmp[l-1-i], tmp[i]
	}

	var buf strings.Builder
	buf.Write(tmp)
	buf.WriteString(strconv.Itoa(height))
	//common.Info.Print("excel pos:",buf.String())
	return buf.String()
}

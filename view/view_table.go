package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/util"
	"image/color"
)

type genericTable struct {
	header    []string
	filter    []string
	colWidth  []float32
	Ctrl      *ctrl.Ctrl
	sortCount int
	MainTable *widget.Table
}

func (obj *genericTable) GetNumCols() int {
	return len(obj.header)
}
func (obj *genericTable) GetColHeader(colIdx int) string {
	var retval string
	if colIdx < len(obj.header) {
		retval = obj.header[colIdx]
	} else {
		obj.Ctrl.Model.LogObj.Printf("GetColHeader invalid colidx %d", colIdx)
	}
	return retval
}
func (obj *genericTable) GetColWidth(colIdx int) float32 {
	var retval float32
	retval = 100
	if colIdx < len(obj.colWidth) {
		retval = obj.colWidth[colIdx]
	} else {
		obj.Ctrl.Model.LogObj.Printf("GetColWidth invalid colidx %d", colIdx)
	}
	return retval
}
func (obj *genericTable) SetFilter(colIdx int, filterStr string) {
	obj.filter[colIdx] = filterStr
}
func (obj *genericTable) GetFilter(colIdx int) (filterStr string) {
	return obj.filter[colIdx]
}

func (obj *genericTable) SetMainTable(mainTable *widget.Table) {
	obj.MainTable = mainTable
}

type guiTableIF2 interface {
	GetColHeader(colIdx int) string
	GetColWidth(colIdx int) float32
	GetNumCols() int
	GetNumRows() int
	GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA)
	GetSumCellStr(colIdx int) (string, color.NRGBA)
	GetCSVCellStr(rowIdx int, colIdx int) string
	UpdateLists()
	SortCol(colIdx int)
	SetFilter(colIdx int, filterStr string)
	GetFilter(colIdx int) (filterStr string)
	SelectedFunc() func(id widget.TableCellID)
	SetMainTable(mainTable *widget.Table)
}

func (obj *OmipGui) createGenTable2(
	table guiTableIF2, monthlyTable bool, sumRowEnable bool, col1FilterName string) fyne.CanvasObject {
	var topRow *widget.Table
	var mainTable *widget.Table
	var bottowmRow *widget.Table
	var sumRow *widget.Table

	topRow = widget.NewTable(
		func() (int, int) { return 1, table.GetNumCols() },
		func() fyne.CanvasObject {
			topButton := widget.NewButton("", func() {
				table.SortCol(0)
				table.UpdateLists()
			})
			return topButton
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			button := cell.(*widget.Button)
			button.SetText(table.GetColHeader(id.Col))
			button.OnTapped = func() {
				table.SortCol(id.Col)
				table.UpdateLists()
				mainTable.Refresh()
			}
		})

	mainTable = widget.NewTable(
		func() (int, int) { return table.GetNumRows(), table.GetNumCols() },
		func() fyne.CanvasObject {
			newText := canvas.NewText("", color.NRGBA{0, 0x80, 0, 0xff})
			newText.Alignment = fyne.TextAlignCenter
			return newText
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			text := cell.(*canvas.Text)
			txt, col := table.GetCellStr(id.Row, id.Col)
			text.Color = col
			text.Text = txt
		})

	bottowmRow = widget.NewTable(
		func() (int, int) {
			cols := table.GetNumCols()
			if monthlyTable {
				cols = 2
			}
			return 1, cols
		},
		func() fyne.CanvasObject {
			newEntry := widget.NewEntry()
			return newEntry
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			entry := cell.(*widget.Entry)
			if monthlyTable && id.Col == 1 {
				entry.SetPlaceHolder(col1FilterName)
			} else {
				entry.SetPlaceHolder(table.GetColHeader(id.Col))
			}

			entry.SetText(table.GetFilter(id.Col))
			entry.OnChanged = func(s string) {
				table.SetFilter(id.Col, s)
				table.UpdateLists()
				mainTable.Refresh()
				sumRow.Refresh()
			}
		})

	sumRow = widget.NewTable(
		func() (int, int) {
			cols := table.GetNumCols()
			return 1, cols
		},
		func() fyne.CanvasObject {
			newText := canvas.NewText("", color.NRGBA{0, 0x80, 0, 0xff})
			newText.Alignment = fyne.TextAlignCenter
			return newText
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			cellCanvasText := cell.(*canvas.Text)
			txt, col := table.GetSumCellStr(id.Col)
			cellCanvasText.Text = txt
			cellCanvasText.Color = col
		})

	for i := 0; i < table.GetNumCols(); i++ {
		topRow.SetColumnWidth(i, table.GetColWidth(i))
		mainTable.SetColumnWidth(i, table.GetColWidth(i))
		bottowmRow.SetColumnWidth(i, table.GetColWidth(i))
		sumRow.SetColumnWidth(i, table.GetColWidth(i))
	}
	mainTable.OnSelected = func(id widget.TableCellID) {
		f := table.SelectedFunc()
		f(id)
		mainTable.Refresh()
	}
	copyCSVbtn := widget.NewButton("COPY CSV", func() {
		var outString string
		for iRow := 0; iRow < table.GetNumRows(); iRow++ {
			for jCol := 0; jCol < table.GetNumCols(); jCol++ {
				outString += fmt.Sprintf("%s\t", table.GetCSVCellStr(iRow, jCol))
			}
			outString += "\n"
		}
		//obj.Ctrl.Model.LogObj.Printf("%s", outString)
		util.ClipboardPaste(outString)
	})
	resetFilterBtn := widget.NewButton("Reset Filter", func() {
		for i := 0; i < table.GetNumCols(); i++ {
			table.SetFilter(i, "")
		}
		table.UpdateLists()
		mainTable.Refresh()
		bottowmRow.Refresh()
		sumRow.Refresh()
	})

	buttonGrid := container.New(layout.NewGridLayout(2), resetFilterBtn, copyCSVbtn)
	botttomGrid := container.New(layout.NewGridLayout(1), sumRow, bottowmRow, buttonGrid)
	if !sumRowEnable {
		sumRow.Hide()
	}
	table.SetMainTable(mainTable)
	return container.NewBorder(topRow, botttomGrid, nil, nil, mainTable)
}

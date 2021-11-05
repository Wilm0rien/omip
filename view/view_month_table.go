package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"regexp"
	"sort"
	"strconv"
)

type monthlyGuiTable struct {
	fulllist     *model.MonthlyTable
	filteredList *model.MonthlyTable
	sortByRow    string
	filterM      bool
	maxColumns   int
	genericTable
}

func NewMonthlyTable(ctrl *ctrl.Ctrl, fullList *model.MonthlyTable, maxMonth int, filterMillions bool) *monthlyGuiTable {
	var obj monthlyGuiTable
	obj.Ctrl = ctrl
	obj.fulllist = fullList
	obj.filterM = filterMillions
	var filteredList model.MonthlyTable
	obj.filteredList = &filteredList
	obj.header = make([]string, 0, 10)
	obj.colWidth = make([]float32, 0, 10)
	obj.colWidth = append(obj.colWidth, 200)
	obj.header = append(obj.header, "CharName")
	obj.maxColumns = maxMonth + 1 // number of month + character column
	for colIdx := maxMonth; colIdx >= 0; colIdx-- {
		_, _, dateStr := ymStrPastMonth(colIdx)
		obj.header = append(obj.header, dateStr)
		obj.colWidth = append(obj.colWidth, 100)
	}
	obj.filter = make([]string, 0, 10)
	for i := 0; i < len(obj.header); i++ {
		obj.filter = append(obj.filter, "")
	}
	obj.sortByRow = "Character"
	return &obj
}

func (obj *monthlyGuiTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	return obj.getCellStrFromList(rowIdx, colIdx, false)
}
func (obj *monthlyGuiTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	var txt string
	if colIdx != 0 {
		dateStr := obj.header[colIdx]

		if value, ok := obj.filteredList.SumInMonth[dateStr]; ok {

			txt = util.HumanizeNumber(value)
			col = *util.GetColor(obj.filteredList.MaxAllTime, value, false)
			if value == 0 {
				txt = ""
			}
		}
	} else {
		txt = "Monthly Sum"
	}

	return txt, col
}
func (obj *monthlyGuiTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	retval, _ := obj.getCellStrFromList(rowIdx, colIdx, true)
	return retval
}
func (obj *monthlyGuiTable) getCellStrFromList(rowIdx int, colIdx int, csv bool) (string, color.NRGBA) {
	var retval string
	var rowCount int
	var colCount int
	color := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	sortedCharList := util.GetSortKeysFromStrMap(obj.filteredList.ValCharPerMon, false)
	if obj.sortByRow != "Character" {
		sort.Slice(sortedCharList, func(i, j int) bool {
			name1 := sortedCharList[i]
			name2 := sortedCharList[j]
			compare1 := obj.fulllist.ValCharPerMon[name1][obj.sortByRow]
			compare2 := obj.fulllist.ValCharPerMon[name2][obj.sortByRow]
			return compare1 >= compare2
		})
	}
	for _, charName := range sortedCharList {
		colCount = 0
		for _, yearMonth := range obj.header {
			if colIdx == colCount && rowIdx == rowCount {
				if colIdx == 0 {
					retval = charName
				} else {
					amount := obj.filteredList.ValCharPerMon[charName][yearMonth]
					if csv {
						retval = fmt.Sprintf("%3.2f", amount)
					} else {
						retval = util.HumanizeNumber(amount)
					}

					color = *util.GetColor(obj.filteredList.MaxInMonth[yearMonth], amount, false)
				}
			}
			colCount++
		}
		rowCount++
	}
	if retval == "0" {
		retval = ""
	}
	return retval, color
}

func (obj *monthlyGuiTable) UpdateLists() {
	obj.filteredList.MaxAllTime = 0
	obj.filteredList.SumInMonth = make(map[string]float64)
	obj.filteredList.MaxInMonth = make(map[string]float64)
	obj.filteredList.ValCharPerMon = make(map[string]map[string]float64)

	sortedCharList := util.GetSortKeysFromStrMap(obj.fulllist.ValCharPerMon, false)

	charFilterStr := obj.filter[0]
	var minimum float64
	if s, err := strconv.ParseFloat(obj.filter[1], 64); err == nil {
		minimum = s
	}

	var filteredMaxAllTime float64
	for _, YearMonth := range obj.header {
		var filteredSum float64
		var filteredMax float64

		for _, charname := range sortedCharList {
			amount := obj.fulllist.ValCharPerMon[charname][YearMonth]
			charNameOK := true
			if charFilterStr != "" {
				charNameOK, _ = regexp.MatchString(fmt.Sprintf("(?i)%s", charFilterStr), charname)
			}
			amountOK := amount >= minimum
			if obj.filterM {
				amountOK = (amount / 1000000) >= minimum
			}

			if charNameOK && amountOK {
				if _, ok := obj.filteredList.ValCharPerMon[charname]; !ok {
					obj.filteredList.ValCharPerMon[charname] = make(map[string]float64)
				}
				obj.filteredList.ValCharPerMon[charname][YearMonth] = amount
				filteredSum += amount
				if obj.filteredList.ValCharPerMon[charname][YearMonth] > filteredMax {
					filteredMax = obj.filteredList.ValCharPerMon[charname][YearMonth]
				}
			}
		}
		obj.filteredList.SumInMonth[YearMonth] = filteredSum
		obj.filteredList.MaxInMonth[YearMonth] = filteredMax
		if filteredSum > filteredMaxAllTime {
			filteredMaxAllTime = filteredSum
		}
	}
	obj.filteredList.MaxAllTime = filteredMaxAllTime
}
func (obj *monthlyGuiTable) SortCol(colIdx int) {
	if colIdx == 0 {
		obj.sortByRow = "Character"
	} else {
		monthIdx := obj.maxColumns - colIdx
		_, _, dateStr := ymStrPastMonth(monthIdx)
		obj.sortByRow = dateStr
	}
}

func (obj *monthlyGuiTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {

	}
}
func (obj *monthlyGuiTable) GetNumRows() int {
	return len(obj.filteredList.ValCharPerMon)
}

func (obj *OmipGui) createBountyTab2(corpID int, maxMonth int) (retTable fyne.CanvasObject, result bool) {
	fullList := obj.Ctrl.Model.GetBountyTable(corpID)
	if fullList.MaxAllTime == 0 {
		result = false
		return retTable, result
	} else {
		table := NewMonthlyTable(obj.Ctrl, fullList, maxMonth, true)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, true, true, "")
		result = true
	}
	return retTable, result
}

package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"math"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type ctrExtInfo struct {
	ItemString    string
	MultipleItems bool
}

type ctrGuiTable struct {
	fulllist      []*model.DBContract
	filteredList  []*model.DBContract
	ctrID2ItemMap map[int]ctrExtInfo
	char          *ctrl.EsiChar
	isCorp        bool
	wList         []*WindowList
	gui           *OmipGui
	genericTable
}

const (
	CTRT0Status = iota
	CTRT1Avail
	CTRT2Items
	CTRT3Price
	CTRT4IssuedTS
	CTRT5CCompTS
	CTRT6Info
	CTRT7Location
)

func NewCTRTable(gui *OmipGui, fullList []*model.DBContract, char *ctrl.EsiChar, isCorp bool) *ctrGuiTable {
	var obj ctrGuiTable
	obj.char = char
	obj.isCorp = isCorp
	obj.Ctrl = gui.Ctrl
	obj.gui = gui
	obj.ctrID2ItemMap = make(map[int]ctrExtInfo)
	obj.fulllist = fullList
	obj.filteredList = make([]*model.DBContract, 0, 10)
	obj.header = make([]string, 0, 10)
	obj.header = append(obj.header, "Status")
	obj.header = append(obj.header, "Availability")
	obj.header = append(obj.header, "Items")
	obj.header = append(obj.header, "Price")
	obj.header = append(obj.header, "DateIssued")
	obj.header = append(obj.header, "DateCompleted")
	obj.header = append(obj.header, "Info")
	obj.header = append(obj.header, "Location")
	obj.colWidth = make([]float32, 0, 10)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 350)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 200)
	obj.colWidth = append(obj.colWidth, 200)
	obj.colWidth = append(obj.colWidth, 300)
	obj.colWidth = append(obj.colWidth, 350)
	obj.filter = make([]string, 0, 10)
	for i := 0; i < len(obj.header); i++ {
		obj.filter = append(obj.filter, "")
	}

	for _, row := range obj.fulllist {
		var extInfo ctrExtInfo
		extInfo.ItemString = "N/A"
		items := obj.Ctrl.Model.GetContrItems(row.Contract_id)
		if len(items) == 1 {
			extInfo.ItemString = obj.Ctrl.Model.GetTypeString(items[0].Type_id)
		} else if len(items) > 1 {
			extInfo.ItemString = gui.GetMaxPricedItem(items)
			extInfo.MultipleItems = true
		}
		obj.ctrID2ItemMap[row.Contract_id] = extInfo
	}
	obj.wList = make([]*WindowList, 0, 10)
	return &obj
}
func (obj *OmipGui) GetMaxPricedItem(items []*model.DBContrItem) (retval string) {
	retval = "multi"
	var maxValue float64
	var sumValue float64
	var foundItemId int
	for _, item := range items {
		if price, ok := obj.Ctrl.Model.ItemAvgPrice[item.Type_id]; ok {
			value := price * float64(item.Quantity)
			sumValue += value
			if value > maxValue {
				maxValue = value
				foundItemId = item.Type_id
			}
		}
	}
	if foundItemId != 0 {
		typeStr := obj.Ctrl.Model.GetTypeString(foundItemId)
		retval = fmt.Sprintf("%s + %d", typeStr, len(items)-1)
	}
	return retval
}
func (obj *ctrGuiTable) GetNumRows() int {
	return len(obj.filteredList)
}
func (obj *ctrGuiTable) getCellStrFromList(rowIdx int, colIdx int, inputList []*model.DBContract) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	var retval string
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {
			case CTRT0Status:
				retval = obj.Ctrl.Model.ContractStatusInt2Str(listElem.Status)
				if listElem.Status == model.Cntr_Stat_outstanding {
					partDuration := listElem.Date_expired - time.Now().Unix()
					if partDuration > 0 {
						duration := listElem.Date_expired - listElem.Date_issued
						ratio := 1 - (float64(partDuration) / float64(duration))
						retval = fmt.Sprintf("%2.0f%%", ratio*100)
					}
				}

			case CTRT1Avail:
				retval = obj.Ctrl.Model.ContractAvailInt2Str(listElem.Availability)
			case CTRT2Items:
				retval = "no items"
				if val, ok := obj.ctrID2ItemMap[listElem.Contract_id]; ok {
					retval = val.ItemString
					if val.MultipleItems {
						col = color.NRGBA{62, 160, 221, 0xff}
					}
				}
			case CTRT3Price:
				retval = util.HumanizeNumber(listElem.Price)
			case CTRT4IssuedTS:
				retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(listElem.Date_issued))
			case CTRT5CCompTS:
				if listElem.Date_completed != 0 {
					retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(listElem.Date_completed))
				} else {
					retval2, future := util.GetTimeDiffStringFromTS(listElem.Date_expired)
					if future {
						retval = retval2
					} else {
						retval = fmt.Sprintf("- %s", retval2)
					}
				}
			case CTRT6Info:
				retval, _ = obj.Ctrl.Model.GetStringEntry(listElem.Title)
			case CTRT7Location:
				retval = obj.Ctrl.GetStructureNameCached(listElem.End_location_id, obj.char)
			}
		}
	}
	return retval, col
}
func (obj *ctrGuiTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	return obj.getCellStrFromList(rowIdx, colIdx, obj.filteredList)
}
func (obj *ctrGuiTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	txt := ""
	switch colIdx {
	case CTRT0Status:
		txt = fmt.Sprintf("Cnt: %d", len(obj.filteredList))
	case CTRT3Price:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += elem.Price
		}
		txt = util.HumanizeNumber(sum)
	case CTRT5CCompTS:
		numFinElem := 0
		var timeSum int64
		for _, elem := range obj.filteredList {
			if elem.Date_completed != 0 {
				if elem.Date_completed < elem.Date_issued {
					obj.Ctrl.Model.LogObj.Printf("warning date completed after date issued %d %d", elem.Date_completed, elem.Date_issued)
				} else {
					timeDiff := elem.Date_completed - elem.Date_issued
					timeSum += timeDiff
					numFinElem++
				}
			}
		}
		if numFinElem > 0 {
			avgDiff := int64(float64(timeSum) / float64(numFinElem))
			txt = util.GetTimeDiffStringFromDiff(avgDiff)
		}

	}

	return txt, col
}

func (obj *ctrGuiTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	var retval string
	if colIdx == CTRT3Price {
		if rowIdx < len(obj.filteredList) {
			retval = fmt.Sprintf("%3.2f", obj.filteredList[rowIdx].Price)
		}
	} else if colIdx == CTRT5CCompTS {
		elem := obj.filteredList[rowIdx]
		if rowIdx < len(obj.filteredList) {
			if elem.Date_completed != 0 {
				retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(elem.Date_completed))
			} else {
				retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(elem.Date_expired))
			}
		}
	} else { // return full float value for csv
		retval, _ = obj.getCellStrFromList(rowIdx, colIdx, obj.filteredList)
	}
	return retval
}
func (obj *ctrGuiTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {
		if id.Col == CTRT2Items {
			if id.Row < len(obj.filteredList) {
				listElem := obj.filteredList[id.Row]
				if val, ok := obj.ctrID2ItemMap[listElem.Contract_id]; ok {
					if val.MultipleItems {

						obj.MainTable.UnselectAll()
						windowTitle := fmt.Sprintf("%d - %s", listElem.Contract_id, val.ItemString)
						windowFound := false
						for _, elem := range obj.wList {
							if elem.Title == windowTitle {
								elem.wRef.Show()
								windowFound = true
								break
							}
						}
						if !windowFound {
							newCtrDT := NewCtrDT(obj.Ctrl)
							items := obj.Ctrl.Model.GetContrItems(listElem.Contract_id)
							for _, item := range items {
								typeStr := obj.Ctrl.Model.GetTypeString(item.Type_id)
								var tableItem ctrDetail
								tableItem.Name = typeStr
								tableItem.Quantity = item.Quantity
								if price, ok2 := obj.Ctrl.Model.ItemAvgPrice[item.Type_id]; ok2 {
									tableItem.AvgValue = price
								}
								newCtrDT.fulllist = append(newCtrDT.fulllist, &tableItem)
							}
							newCtrDT.UpdateLists()
							w := fyne.CurrentApp().NewWindow(windowTitle)
							ctrDetailWidget := obj.gui.createGenTable2(newCtrDT, false, true, "")
							w.SetContent(ctrDetailWidget)
							height := float32(700)
							if len(newCtrDT.fulllist) < 20 {
								height = float32(50*len(newCtrDT.fulllist)) + 150
							}
							w.Resize(fyne.NewSize(1024, height))
							obj.wList = append(obj.wList, &WindowList{windowTitle, w})
							w.Show()
							w.SetOnClosed(func() {
								found := false
								foundIdx := 0
								for idx, elem := range obj.wList {
									if elem.Title == w.Title() {
										found = true
										foundIdx = idx
										break
									}
								}
								if found {
									obj.wList = append(obj.wList[:foundIdx], obj.wList[foundIdx+1:]...)
								}
							})

						}
					}
				}
			}
		}
	}
}
func (obj *ctrGuiTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			filterString := true
			if colIdx == CTRT0Status {
				elem := obj.fulllist[rowIdx]
				if elem.Status == model.Cntr_Stat_outstanding {
					if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
						partDuration := elem.Date_expired - time.Now().Unix()
						duration := elem.Date_expired - elem.Date_issued
						if duration > 0 {
							ratio := (1 - (float64(partDuration) / float64(duration))) * 100
							ratio = math.Round(ratio)
							filterString = false
							if ratio < s {
								filterOK = false
								break
							}
						}
					}
				}
			} else if colIdx == CTRT3Price {
				filterString = false
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (math.Abs(obj.fulllist[rowIdx].Price) / 1000000) < s {
						filterOK = false
						break
					}
				}
			}
			if filterString {
				currentCellStr, _ := obj.getCellStrFromList(rowIdx, colIdx, obj.fulllist)
				fMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", currentFilter), currentCellStr)
				if !fMatch {
					filterOK = false
					break
				}
			}
		}
		if filterOK {
			obj.filteredList = append(obj.filteredList, obj.fulllist[rowIdx])
		}
	}
}
func (obj *ctrGuiTable) SortCol(colIdx int) {
	obj.sortCount++
	sort.Slice(obj.fulllist, func(i, j int) bool {
		a := obj.fulllist[i]
		b := obj.fulllist[j]
		var retval bool
		switch colIdx {
		case CTRT0Status:
			if a.Status == model.Cntr_Stat_outstanding && b.Status == model.Cntr_Stat_outstanding {
				aDuration := a.Date_expired - time.Now().Unix()
				bDuration := b.Date_expired - time.Now().Unix()
				retval = aDuration >= bDuration
			} else {
				aStr := obj.Ctrl.Model.ContractStatusInt2Str(a.Status)
				bStr := obj.Ctrl.Model.ContractStatusInt2Str(b.Status)
				retval = aStr >= bStr
			}
		case CTRT1Avail:
			aStr := obj.Ctrl.Model.ContractAvailInt2Str(a.Availability)
			bStr := obj.Ctrl.Model.ContractAvailInt2Str(b.Availability)
			retval = aStr >= bStr
		case CTRT2Items:
			aStr := "no items"
			bStr := "no items"
			if val, ok := obj.ctrID2ItemMap[a.Contract_id]; ok {
				aStr = val.ItemString
			}
			if val2, ok2 := obj.ctrID2ItemMap[b.Contract_id]; ok2 {
				bStr = val2.ItemString
			}
			retval = aStr >= bStr
		case CTRT3Price:
			retval = a.Price >= b.Price
		case CTRT4IssuedTS:
			retval = a.Date_issued >= b.Date_issued
		case CTRT5CCompTS:
			if a.Date_completed == 0 && b.Date_completed == 0 {
				retval = a.Date_expired >= b.Date_expired
			} else {
				retval = a.Date_completed >= b.Date_completed
			}
		case CTRT6Info:
			aStr, _ := obj.Ctrl.Model.GetStringEntry(a.Title)
			bStr, _ := obj.Ctrl.Model.GetStringEntry(b.Title)
			retval = aStr >= bStr
		case CTRT7Location:
			aStr := obj.Ctrl.GetStructureNameCached(a.End_location_id, obj.char)
			bStr := obj.Ctrl.GetStructureNameCached(b.End_location_id, obj.char)
			if aStr == bStr {
				aStr := "no items"
				bStr := "no items"
				if val, ok := obj.ctrID2ItemMap[a.Contract_id]; ok {
					aStr = val.ItemString
				}
				if val2, ok2 := obj.ctrID2ItemMap[b.Contract_id]; ok2 {
					bStr = val2.ItemString
				}
				retval = aStr >= bStr
			} else {
				retval = aStr >= bStr
			}

		}

		if obj.sortCount%2 == 0 {
			retval = !retval
		}
		return retval
	})
}

func (obj *OmipGui) createContractTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	var fullList []*model.DBContract
	if corp {
		fullList = obj.Ctrl.Model.GetCorpContracts(char.CharInfoExt.CooperationId)
	} else {
		fullList = obj.Ctrl.Model.GetContractsByIssuerId(char.CharInfoData.CharacterID, corp)
	}
	if len(fullList) == 0 {
		result = false
		return retTable, result
	} else {
		table := NewCTRTable(obj, fullList, char, corp)
		table.sortCount = 0
		table.SortCol(CTRT0Status)
		table.UpdateLists()

		retTable = obj.createGenTable2(table, false, true, "")
		result = true
	}

	return retTable, result
}

type ctrDetail struct {
	Name     string
	Quantity int
	AvgValue float64
}

type ctrDetailTable struct {
	fulllist     []*ctrDetail
	filteredList []*ctrDetail
	genericTable
}

const (
	CTRDTCol0Name = iota
	CTRDTCol1Quantity
	CTRDTCol2Value
	CTRDTCol4TotalValue
)

// NewCtrDT Contract Detail Table
func NewCtrDT(ctrl *ctrl.Ctrl) *ctrDetailTable {
	var ctrDT ctrDetailTable
	ctrDT.Ctrl = ctrl
	ctrDT.fulllist = make([]*ctrDetail, 0, 10)
	ctrDT.filteredList = make([]*ctrDetail, 0, 10)
	ctrDT.header = make([]string, 0, 10)
	ctrDT.header = append(ctrDT.header, "Item Name")
	ctrDT.header = append(ctrDT.header, "Quantity")
	ctrDT.header = append(ctrDT.header, "Unit Price")
	ctrDT.header = append(ctrDT.header, "Total Value")
	ctrDT.colWidth = make([]float32, 0, 10)
	ctrDT.colWidth = append(ctrDT.colWidth, 400)
	ctrDT.colWidth = append(ctrDT.colWidth, 100)
	ctrDT.colWidth = append(ctrDT.colWidth, 100)
	ctrDT.colWidth = append(ctrDT.colWidth, 100)
	ctrDT.filter = make([]string, 0, 10)
	for i := 0; i < len(ctrDT.header); i++ {
		ctrDT.filter = append(ctrDT.filter, "")
	}
	return &ctrDT
}
func (obj *ctrDetailTable) GetNumRows() int {
	return len(obj.filteredList)
}
func (obj *ctrDetailTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	return obj.getCellStrFull(rowIdx, colIdx, obj.filteredList), col
}
func (obj *ctrDetailTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	retval := ""
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	switch colIdx {
	case CTRDTCol0Name:
	case CTRDTCol1Quantity:
	case CTRDTCol2Value:
	case CTRDTCol4TotalValue:
		var sum float64
		for _, item := range obj.filteredList {
			total := float64(item.Quantity) * item.AvgValue
			sum += total
		}
		retval = util.HumanizeNumber(sum)
	}
	return retval, col
}

func (obj *ctrDetailTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	var retval string
	if colIdx != CTRDTCol2Value && colIdx != CTRDTCol4TotalValue {
		retval = obj.getCellStrFull(rowIdx, colIdx, obj.filteredList)
	} else { // return full float value for csv
		if colIdx == CTRDTCol2Value {
			retval = fmt.Sprintf("%3.2f", obj.filteredList[rowIdx].AvgValue)
		} else if colIdx == CTRDTCol4TotalValue {
			total := float64(obj.filteredList[rowIdx].Quantity) * obj.filteredList[rowIdx].AvgValue
			retval = fmt.Sprintf("%3.2f", total)
		}
	}
	return retval
}

func (obj *ctrDetailTable) getCellStrFull(rowIdx int, colIdx int, inputList []*ctrDetail) string {
	var retval string
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {
			case CTRDTCol0Name:
				retval = listElem.Name
			case CTRDTCol1Quantity:
				retval = fmt.Sprintf("%d", listElem.Quantity)
			case CTRDTCol2Value:
				retval = util.HumanizeNumber(listElem.AvgValue)
			case CTRDTCol4TotalValue:
				total := float64(listElem.Quantity) * listElem.AvgValue
				retval = util.HumanizeNumber(total)
			}
		}
	}
	return retval
}

func (obj *ctrDetailTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {

	}
}
func (obj *ctrDetailTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			if colIdx == CTRDTCol0Name {
				if currentFilter != "" {
					fMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", currentFilter), obj.fulllist[rowIdx].Name)
					if !fMatch {
						filterOK = false
						break
					}
				}
			}
			if colIdx == CTRDTCol1Quantity {
				if s, err := strconv.ParseInt(currentFilter, 10, 32); err == nil {
					if (obj.fulllist[rowIdx].Quantity) < int(s) {
						filterOK = false
						break
					}
				}
			}
			if colIdx == CTRDTCol2Value {
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (obj.fulllist[rowIdx].AvgValue / 1000000) < s {
						filterOK = false
						break
					}
				}
			}
			if colIdx == CTRDTCol4TotalValue {
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if ((obj.fulllist[rowIdx].AvgValue * float64(obj.fulllist[rowIdx].Quantity)) / 1000000) < s {
						filterOK = false
						break
					}
				}
			}
		}
		if filterOK {
			obj.filteredList = append(obj.filteredList, obj.fulllist[rowIdx])
		}
	}
}

func (obj *ctrDetailTable) SortCol(colIdx int) {
	obj.sortCount++
	sort.Slice(obj.fulllist, func(i, j int) bool {
		var retval bool
		switch colIdx {
		case CTRDTCol0Name:
			retval = obj.fulllist[i].Name >= obj.fulllist[j].Name
		case CTRDTCol1Quantity:
			retval = obj.fulllist[i].Quantity >= obj.fulllist[j].Quantity
		case CTRDTCol2Value:
			retval = obj.fulllist[i].AvgValue >= obj.fulllist[j].AvgValue
		case CTRDTCol4TotalValue:
			a := obj.fulllist[i].AvgValue * float64(obj.fulllist[i].Quantity)
			b := obj.fulllist[j].AvgValue * float64(obj.fulllist[j].Quantity)
			retval = a >= b
		}
		if obj.sortCount%2 == 0 {
			retval = !retval
		}
		return retval
	})
}

package view

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type miningDetail struct {
	altName   string
	dateStr   string
	oreType   string
	oreAmount int
	oreVolume int
	iskValue  float64
}

type miningDetailTable struct {
	fulllist     []*miningDetail
	filteredList []*miningDetail
	genericTable
}

const (
	MDTCol0AltName = iota
	MDTCol1DateString
	MDTCol2OreTypeStr
	MDTCol3OreAmount
	MDTCol4OreVolume
	MDTCol5IskValue
)

const (
	GROUP_SEL_CHAR = "Group By Char"
	GROUP_SEL_CORP = "Group By Corp"
)

func NewMDT(ctrl *ctrl.Ctrl) *miningDetailTable {
	var mDT miningDetailTable
	mDT.Ctrl = ctrl
	mDT.fulllist = make([]*miningDetail, 0, 10)
	mDT.filteredList = make([]*miningDetail, 0, 10)
	mDT.header = make([]string, 0, 10)
	mDT.header = append(mDT.header, "alt name")
	mDT.header = append(mDT.header, "date")
	mDT.header = append(mDT.header, "ore type")
	mDT.header = append(mDT.header, "ore amount")
	mDT.header = append(mDT.header, "ore volume")
	mDT.header = append(mDT.header, "isk value")
	mDT.colWidth = make([]float32, 0, 10)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.colWidth = append(mDT.colWidth, 150)
	mDT.filter = make([]string, 0, 10)
	for i := 0; i < len(mDT.header); i++ {
		mDT.filter = append(mDT.filter, "")
	}
	return &mDT
}
func (obj *miningDetailTable) GetNumRows() int {
	return len(obj.filteredList)
}

func (obj *miningDetailTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	return obj.getCellStrFull(rowIdx, colIdx, obj.filteredList), col
}
func (obj *miningDetailTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	txt := ""
	switch colIdx {
	case MDTCol0AltName:
	case MDTCol1DateString:
	case MDTCol2OreTypeStr:
	case MDTCol3OreAmount:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += float64(elem.oreAmount)
		}
		txt = util.HumanizeNumber(sum)
	case MDTCol4OreVolume:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += float64(elem.oreVolume)
		}
		txt = util.HumanizeNumber(sum)
	case MDTCol5IskValue:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += elem.iskValue
		}
		txt = util.HumanizeNumber(sum)
	}
	return txt, col
}
func (obj *miningDetailTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	var retval string

	retval = obj.getCellStrFull(rowIdx, colIdx, obj.filteredList)

	return retval
}
func (obj *miningDetailTable) getCellStrFull(rowIdx int, colIdx int, inputList []*miningDetail) string {
	var retval string
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {

			case MDTCol0AltName:
				//retval = listElem.altName

				//retval = "MDTCol0AltName " + fmt.Sprintf("%d %d", rowIdx, colIdx)
				retval = listElem.altName
			case MDTCol1DateString:

				//retval = listElem.dateStr

				//retval = "MDTCol1DateString"
				retval = listElem.dateStr
			case MDTCol2OreTypeStr:

				//retval = listElem.oreType
				retval = listElem.oreType
			case MDTCol3OreAmount:
				//retval = util.HumanizeNumber((float64)(listElem.oreAmount))
				retval = util.HumanizeNumber(float64(listElem.oreAmount))
			case MDTCol4OreVolume:
				//retval = util.HumanizeNumber((float64)(listElem.oreVolume))
				retval = util.HumanizeNumber(float64(listElem.oreVolume))
			case MDTCol5IskValue:
				//retval = util.HumanizeNumber(listElem.iskValue)
				retval = util.HumanizeNumber(listElem.iskValue)

			}
			//retval = fmt.Sprintf("%d %d", rowIdx, colIdx)
		}
	}
	return retval
}
func (obj *miningDetailTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {
		// nothing to do
	}
}
func (obj *miningDetailTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			if colIdx == MDTCol5IskValue {
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (obj.fulllist[rowIdx].iskValue / 1000000) < s {
						filterOK = false
						break
					}
				}
			} else {
				if currentFilter != "" {
					currentCellStr := obj.getCellStrFull(rowIdx, colIdx, obj.fulllist)
					fMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", currentFilter), currentCellStr)
					if !fMatch {
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
func (obj *miningDetailTable) SortCol(colIdx int) {
	obj.sortCount++
	sort.Slice(obj.fulllist, func(i, j int) bool {
		var retval bool
		a := obj.getCellStrFull(i, colIdx, obj.fulllist)
		b := obj.getCellStrFull(j, colIdx, obj.fulllist)
		if j <= len(obj.fulllist) {
			switch colIdx {
			case MDTCol3OreAmount:
				retval = obj.fulllist[i].oreAmount >= obj.fulllist[j].oreAmount
			case MDTCol4OreVolume:
				retval = obj.fulllist[i].oreVolume >= obj.fulllist[j].oreVolume
			case MDTCol5IskValue:
				retval = obj.fulllist[i].iskValue >= obj.fulllist[j].iskValue
			default:
				retval = strings.ToUpper(a) >= strings.ToUpper(b)
			}
		}

		if obj.sortCount%2 == 0 {
			retval = !retval
		}
		return retval
	})
}

func (obj *OmipGui) createMiningTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	// TODO add obervers to export table
	// todo add filter by observer
	// todo add alliance filter
	maxMonth := 12
	updateRunning := false

	groupSelectionStr := ""                     // used by groupSelect
	typeSelectionStr := "ISK"                   // used by typeSelect
	ColumnHdrCharNameBtnStr := "Character Name" // used to change column header
	lastUpdateTime := time.Now()
	origMiningData := obj.Ctrl.Model.GetCorpMiningData(char.CharInfoExt.CooperationId)
	membermap := make(map[int]int)
	for _, elem := range origMiningData {
		membermap[elem.MainID] = 1
	}
	obsList := obj.Ctrl.Model.GetCorpObservers(char.CharInfoExt.CooperationId)
	for _, obsId := range obsList {
		extList := obj.Ctrl.Model.GetExtMiningData(char.CharInfoExt.CooperationId, obsId)
		for _, elem := range extList {
			if _, ok := membermap[elem.MainID]; !ok {
				origMiningData = append(origMiningData, elem)
			}
		}
		//origMiningData = append(origMiningData, extList...)
	}
	if len(origMiningData) == 0 {
		return nil, false
	}

	var fullListOre *model.MonthlyTable
	var fullListIsk *model.MonthlyTable
	tickerMap := make(map[string]int)
	nameMapping := make(map[string]int)
	// rebuilding is necessary when switching from char view to corp view via groupSelect
	rebuildMTable := func() {
		normalizedListOre := make([]*model.DBTable, 0, 100)
		normalizedListIsk := make([]*model.DBTable, 0, 100)

		for _, elem := range origMiningData {
			var new_monthlyOre model.DBTable
			combinedMain := fmt.Sprintf("[%s] %s", elem.Ticker, elem.MainName)
			combinedAlt := fmt.Sprintf("[%s] %s", elem.Ticker, elem.AltName)
			if groupSelectionStr == GROUP_SEL_CORP {
				combinedMain = elem.Ticker
				combinedAlt = elem.MainName
			}
			tickerMap[elem.Ticker] = elem.RecordedCorporationID
			new_monthlyOre.MainName = combinedMain
			new_monthlyOre.AltName = combinedAlt
			new_monthlyOre.Time = elem.LastUpdated
			nameMapping[combinedMain] = elem.MainID
			var volume float64
			if props := obj.Ctrl.Model.GetSdePropsByID(elem.TypeID); props != nil {
				volume = props.GetVolume()
			}
			new_monthlyOre.Amount = (float64)(elem.Quantity) * volume

			var new_monthlyIsk model.DBTable
			new_monthlyIsk.MainName = combinedMain
			new_monthlyIsk.AltName = combinedAlt
			new_monthlyIsk.Time = elem.LastUpdated

			if value, err := obj.Ctrl.GetOreValueByAmount(elem.TypeID, elem.Quantity); err == nil {
				new_monthlyIsk.Amount = value
			} else {
				new_monthlyIsk.Amount = 0
			}

			normalizedListOre = append(normalizedListOre, &new_monthlyOre)
			normalizedListIsk = append(normalizedListIsk, &new_monthlyIsk)
		}
		fullListOre = obj.Ctrl.Model.GetMonthlyTable(char.CharInfoExt.CooperationId, normalizedListOre, maxMonth)
		fullListIsk = obj.Ctrl.Model.GetMonthlyTable(char.CharInfoExt.CooperationId, normalizedListIsk, maxMonth)
	}
	rebuildMTable()

	var tableObj *widget.Table
	var filterCharName *widget.Entry
	var filterAmount *widget.Entry
	sortByRow := "Character"
	var updateColumnWidth func()
	var filteredList model.MonthlyTable

	if len(fullListOre.ValCharPerMon) == 0 {
		return retTable, result
	} else {
		result = true
	}

	filteredCharList := make([]string, 0, 10)
	fullList := fullListOre
	filterReverse := false
	updateLists := func() {
		switch typeSelectionStr {
		case "ORE":
			fullList = fullListOre
		case "ISK":
			fullList = fullListIsk
		}
		filteredCharList = make([]string, 0, 10)
		filteredList.MaxAllTime = 0
		filteredList.SumInMonth = make(map[string]float64)
		filteredList.MaxInMonth = make(map[string]float64)
		filteredList.ValCharPerMon = make(map[string]map[string]float64)
		var minimum float64
		if s, err := strconv.ParseFloat(filterAmount.Text, 64); err == nil {
			minimum = s
		}

		keyList := util.GetSortKeysFromStrMap(fullList.ValCharPerMon, false)

		sort.Slice(keyList, func(i, j int) bool {
			name1 := keyList[i]
			name2 := keyList[j]
			compare1 := fullList.ValCharPerMon[name1][sortByRow]
			compare2 := fullList.ValCharPerMon[name2][sortByRow]
			if filterReverse {
				return compare1 < compare2
			}
			return compare1 >= compare2
		})

		for _, charName := range keyList {
			fNameMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterCharName.Text), charName)

			amountMatch := false

			for _, dateStr := range util.GetSortKeysFromStrMap(fullList.ValCharPerMon[charName], false) {
				amount := fullList.ValCharPerMon[charName][dateStr] / 1000000
				if amount >= minimum {
					amountMatch = true
					break
				}
			}

			if fNameMatch && amountMatch {
				filteredCharList = append(filteredCharList, charName)
				if _, ok := filteredList.ValCharPerMon[charName]; !ok {
					filteredList.ValCharPerMon[charName] = make(map[string]float64)
				}
				for _, dateStr := range util.GetSortKeysFromStrMap(fullList.ValCharPerMon[charName], false) {
					amount := fullList.ValCharPerMon[charName][dateStr]
					// only in ISK view the percentage is relevant
					if char.GuiSettings.CorpMining.SelType == "ISK" && char.GuiSettings.CorpMining.Percentage != 0 {
						amount = amount * char.GuiSettings.CorpMining.Percentage / 100
					}
					filteredList.ValCharPerMon[charName][dateStr] += amount
					filteredList.SumInMonth[dateStr] += amount
					if filteredList.MaxInMonth[dateStr] < filteredList.ValCharPerMon[charName][dateStr] {
						filteredList.MaxInMonth[dateStr] = filteredList.ValCharPerMon[charName][dateStr]
					}
				}
			}
		}
		for _, elem := range filteredList.SumInMonth {
			if elem > filteredList.MaxAllTime {
				filteredList.MaxAllTime = elem
			}
		}
		if updateColumnWidth != nil {
			updateColumnWidth()
		}
		if filterReverse {
			filterReverse = false
		} else {
			filterReverse = true
		}
		tableObj.Refresh()
	}
	tableObj = widget.NewTable(
		func() (int, int) { return len(filteredCharList), maxMonth + 1 },
		func() fyne.CanvasObject {
			newText := canvas.NewText("", color.NRGBA{0, 0x80, 0, 0xff})
			newText.Alignment = fyne.TextAlignCenter
			return newText
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			charName := filteredCharList[id.Row]
			text := cell.(*canvas.Text)
			outText := ""
			color := &color.NRGBA{0xff, 0xff, 0xff, 0xff}
			if id.Col == 0 {
				outText = charName
			} else {
				monthIdx := maxMonth - id.Col
				_, _, dateStr := ymStrPastMonth(monthIdx)
				if _, ok := filteredList.ValCharPerMon[charName]; ok {
					if _, ok2 := filteredList.ValCharPerMon[charName][dateStr]; ok2 {
						outText = util.HumanizeNumber(filteredList.ValCharPerMon[charName][dateStr])
						color = util.GetColor(filteredList.MaxInMonth[dateStr], filteredList.ValCharPerMon[charName][dateStr], false)
					} else {
						//obj.Ctrl.Model.LogObj.Printf("warning %s %s does not exist in ilteredList.ValCharPerMon", charName, dateStr)
					}
				} else {
					obj.Ctrl.Model.LogObj.Printf("warning %s does not exist in ilteredList.ValCharPerMon", charName)
				}
			}
			text.Text = outText
			text.Color = color
		})
	wList := make([]*WindowList, 0, 10)
	tableObj.OnSelected = func(id widget.TableCellID) {
		tableObj.UnselectAll()
		if id.Row < len(filteredCharList) {
			charName := filteredCharList[id.Row]
			if id.Col == 0 {
				obj.Ctrl.Model.LogObj.Printf("ERROR unexpected value %s", charName)
			} else {
				monthIdx := maxMonth - id.Col
				year, month, _ := ymStrPastMonth(monthIdx)
				startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				endTime := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
				windowTitle := fmt.Sprintf("Mining from %s in %s %d", charName, startTime.Month().String(), startTime.Year())
				windowFound := false
				for _, elem := range wList {
					if elem.Title == windowTitle {
						elem.wRef.Show()
						windowFound = true
						break
					}
				}

				if !windowFound {
					if _, ok := nameMapping[charName]; !ok {
						obj.Ctrl.Model.LogObj.Printf("ERROR %s not found in map", charName)
						return
					}
					var list []*model.ViewMiningData
					if groupSelectionStr == GROUP_SEL_CHAR {
						list = obj.Ctrl.Model.GetMiningFiltered(char.CharInfoExt.CooperationId, nameMapping[charName], startTime.Unix(), endTime.Unix())
						if len(list) == 0 {
							list = obj.Ctrl.Model.GetMiningFilteredExt(nameMapping[charName], startTime.Unix(), endTime.Unix())
						}
					} else {
						if corpId, ok := tickerMap[charName]; ok {
							list = obj.Ctrl.Model.GetMiningByCop(corpId, startTime.Unix(), endTime.Unix())
						} else {
							obj.Ctrl.Model.LogObj.Printf("ERROR %s not found in map", charName)
						}
					}

					if len(list) > 0 {
						//list = list[:10]
						miningDT := NewMDT(obj.Ctrl)
						for _, elem := range list {
							var newElem miningDetail
							newElem.altName = elem.AltName
							newElem.dateStr = util.ConvertUnixTimeToDateStr(elem.LastUpdated)
							newElem.oreType = obj.Ctrl.Model.GetTypeString(elem.TypeID)
							newElem.oreAmount = elem.Quantity
							var volume float64
							if props := obj.Ctrl.Model.GetSdePropsByID(elem.TypeID); props != nil {
								volume = props.GetVolume()
							}
							newElem.oreVolume = int((float64)(elem.Quantity) * volume)
							if oreIskValue, err := obj.Ctrl.GetOreValueByAmount(elem.TypeID, elem.Quantity); err == nil {
								newElem.iskValue = oreIskValue
							} else {
								newElem.iskValue = 0
							}

							miningDT.fulllist = append(miningDT.fulllist, &newElem)
						}
						miningDT.UpdateLists()
						w := fyne.CurrentApp().NewWindow(windowTitle)
						miningDetailWidget := obj.createGenTable2(miningDT, false, true, "")
						w.SetContent(miningDetailWidget)
						height := float32(700)
						if len(miningDT.fulllist) < 20 {
							height = float32(50*len(miningDT.fulllist)) + 100
						}
						w.Resize(fyne.NewSize(1024, height))
						wList = append(wList, &WindowList{windowTitle, w})
						w.Show()
						w.SetOnClosed(func() {
							found := false
							foundIdx := 0
							for idx, elem := range wList {
								if elem.Title == w.Title() {
									found = true
									foundIdx = idx
									break
								}
							}
							if found {
								wList = append(wList[:foundIdx], wList[foundIdx+1:]...)
								obj.Ctrl.Model.LogObj.Printf("removing %s", w.Title())
							}
						})
					}
				}
			}
		}
	}
	topRowTable := widget.NewTable(
		func() (int, int) { return 1, maxMonth + 1 },
		func() fyne.CanvasObject {
			topButton := widget.NewButton(ColumnHdrCharNameBtnStr, func() {
				sortByRow = "Character"
				updateLists()
			})
			return topButton
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			button := cell.(*widget.Button)
			if id.Col == 0 {
				button.SetText(ColumnHdrCharNameBtnStr)
				button.OnTapped = func() {
					sortByRow = "Character"
					updateLists()
				}
			} else {
				monthIdx := maxMonth - id.Col
				curYear, curMonth, _ := ymStrPastMonth(monthIdx)
				dateStr2 := fmt.Sprintf("%02d-%02d", curYear-2000, curMonth)
				button.SetText(dateStr2)
				button.OnTapped = func() {
					sortByRow = dateStr2
					updateLists()
				}
			}
		},
	)
	bottomRowTable := widget.NewTable(
		func() (int, int) { return 1, maxMonth + 1 },
		func() fyne.CanvasObject {
			newText := canvas.NewText("", color.NRGBA{0, 0x80, 0, 0xff})
			newText.Alignment = fyne.TextAlignCenter
			return newText
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			cellCanvasText := cell.(*canvas.Text)
			if id.Col == 0 {
				cellCanvasText.Text = "Monthly Sum"
				cellCanvasText.Color = color.NRGBA{0xff, 0xff, 0xff, 0xff}
			} else {
				monthIdx := maxMonth - id.Col
				curYear, curMonth, _ := ymStrPastMonth(monthIdx)
				dateStr2 := fmt.Sprintf("%02d-%02d", curYear-2000, curMonth)
				cellCanvasText.Color = util.GetColor(fullList.MaxAllTime, filteredList.SumInMonth[dateStr2], false)
				cellCanvasText.Text = util.HumanizeNumber(fullList.SumInMonth[dateStr2])
			}
		},
	)
	filterDelayFunc := func() {
		lastUpdateTime = time.Now()
		if !updateRunning {
			updateRunning = true
			go func() {
				// wait for no change for at least 2s
				for {
					if time.Since(lastUpdateTime).Milliseconds() > 500 {
						updateLists()
						updateRunning = false
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}()
		}
	}

	filterCharName = widget.NewEntry()
	filterCharName.PlaceHolder = "filter char name"
	filterCharName.OnChanged = func(s string) {
		filterDelayFunc()
		char.GuiSettings.CorpMining.CharName = s
	}
	filterCharName.SetText(char.GuiSettings.CorpMining.CharName)
	filterAmount = widget.NewEntry()
	filterAmount.PlaceHolder = "filter millions"
	filterAmount.OnChanged = func(s string) {
		char.GuiSettings.CorpMining.FilterAmount = s
		filterDelayFunc()
	}
	filterAmount.SetText(char.GuiSettings.CorpMining.FilterAmount)

	updateLists()

	updateColumnWidth = func() {
		leadingWidth := float32(300)
		tableObj.SetColumnWidth(0, leadingWidth)
		topRowTable.SetColumnWidth(0, leadingWidth)
		bottomRowTable.SetColumnWidth(0, leadingWidth)
		for i := 1; i < maxMonth+1; i++ {
			correction := float32(2.5)
			colwidth := 90 - correction
			tableObj.SetColumnWidth(i, colwidth-correction-2.5)
			topRowTable.SetColumnWidth(i, colwidth-correction-2.5)
			bottomRowTable.SetColumnWidth(i, colwidth-correction-2.5)
		}
	}
	updateColumnWidth()
	hintLabel := widget.NewLabel("Hint: Click cell to open character sheet")

	if char.GuiSettings.CorpMining.SelType == "" {
		char.GuiSettings.CorpMining.SelType = "ISK"
	}

	percentageEntry := widget.NewEntry()
	percentageEntry.SetPlaceHolder("Percentage")
	percentageBtn := widget.NewButton("Update", func() {
		if s, err := strconv.ParseFloat(percentageEntry.Text, 64); err == nil {
			if s > 0 && s <= 100 {
				char.GuiSettings.CorpMining.Percentage = s
				updateLists()
			} else {
				d := dialog.NewError(errors.New(fmt.Sprintf("percentage must be between 0..100 %s", filterAmount.Text)), obj.WindowPtr)
				d.Show()
			}
		} else {
			d := dialog.NewError(errors.New(fmt.Sprintf("no float number %s %s", filterAmount.Text, err.Error())), obj.WindowPtr)
			d.Show()
		}
	})
	if char.GuiSettings.CorpMining.Percentage != 0 {
		percentageEntry.SetText(fmt.Sprintf("%f", char.GuiSettings.CorpMining.Percentage))
	}
	percentageGrid := container.NewGridWithColumns(2, percentageEntry, percentageBtn)
	showHidePerc := func(s string) {
		switch s {
		case "ORE":
			percentageGrid.Hide()
		case "ISK":
			percentageGrid.Show()
		}
	}
	typeSelect := widget.NewSelect([]string{"ORE", "ISK"}, func(s string) {
		char.GuiSettings.CorpMining.SelType = s
		typeSelectionStr = s
		showHidePerc(s)
		updateLists()
		bottomRowTable.Refresh()
	})
	groupSelect := widget.NewSelect([]string{GROUP_SEL_CHAR, GROUP_SEL_CORP}, func(s string) {
		if s == GROUP_SEL_CHAR {
			ColumnHdrCharNameBtnStr = "Character Name"
		} else {
			ColumnHdrCharNameBtnStr = "Corp Ticker"
		}
		groupSelectionStr = s
		char.GuiSettings.CorpMining.GroupSelection = s
		rebuildMTable()
		updateLists()

		bottomRowTable.Refresh()
	})

	if char.GuiSettings.CorpMining.SelType != "" {
		typeSelectionStr = char.GuiSettings.CorpMining.SelType
		typeSelect.SetSelected(char.GuiSettings.CorpMining.SelType)
	} else {
		char.GuiSettings.CorpMining.SelType = "ISK"
		typeSelectionStr = "ISK"
	}
	if char.GuiSettings.CorpMining.GroupSelection != "" {
		groupSelectionStr = char.GuiSettings.CorpMining.GroupSelection
		groupSelect.SetSelected(groupSelectionStr)
	} else {
		char.GuiSettings.CorpMining.GroupSelection = GROUP_SEL_CHAR
		groupSelectionStr = GROUP_SEL_CHAR
		groupSelect.SetSelected(groupSelectionStr)
	}
	showHidePerc(char.GuiSettings.CorpMining.SelType)
	filtergrid := container.New(layout.NewGridLayout(6), filterCharName, filterAmount, typeSelect, percentageGrid, groupSelect, hintLabel)
	bottomGrid := container.New(layout.NewGridLayout(1), bottomRowTable, filtergrid)
	topGrid2 := container.New(layout.NewGridLayout(1), topRowTable)
	mainbox := container.NewBorder(topGrid2, bottomGrid, nil, nil, tableObj)
	return mainbox, result
}

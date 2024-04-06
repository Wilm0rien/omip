package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type kmDetail struct {
	altName string
	timeTS  int64
	value   float64
	kmId    int
	npc     bool
	zKillOK model.ZkStatus
}

type kmDetailTable struct {
	fulllist     []*kmDetail
	filteredList []*kmDetail
	genericTable
}

const (
	KMDTCol0AltName = iota
	KMDTCol1TS
	KMDTCol2KValue
	KMDTCol3KLink
)

func NewKMDT(ctrl *ctrl.Ctrl) *kmDetailTable {
	var kmDT kmDetailTable
	kmDT.Ctrl = ctrl
	kmDT.fulllist = make([]*kmDetail, 0, 10)
	kmDT.filteredList = make([]*kmDetail, 0, 10)
	kmDT.header = make([]string, 0, 10)
	kmDT.header = append(kmDT.header, "alt name")
	kmDT.header = append(kmDT.header, "time stamp")
	kmDT.header = append(kmDT.header, "value")
	kmDT.header = append(kmDT.header, "kill link")
	kmDT.colWidth = make([]float32, 0, 10)
	kmDT.colWidth = append(kmDT.colWidth, 300)
	kmDT.colWidth = append(kmDT.colWidth, 150)
	kmDT.colWidth = append(kmDT.colWidth, 100)
	kmDT.colWidth = append(kmDT.colWidth, 400)

	kmDT.filter = make([]string, 0, 10)
	for i := 0; i < len(kmDT.header); i++ {
		kmDT.filter = append(kmDT.filter, "")
	}
	return &kmDT
}

func (obj *kmDetailTable) GetNumRows() int {
	return len(obj.filteredList)
}

func (obj *kmDetailTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	if rowIdx < len(obj.filteredList) {
		if colIdx == KMDTCol3KLink {
			if obj.filteredList[rowIdx].zKillOK == model.ZkUnkown {
				col = color.NRGBA{128, 128, 128, 0xff}
			}
		}
	}
	return obj.getCellStrFull(rowIdx, colIdx, obj.filteredList), col
}
func (obj *kmDetailTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	return "", col
}
func (obj *kmDetailTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	var retval string
	if colIdx != KMDTCol2KValue {
		retval = obj.getCellStrFull(rowIdx, colIdx, obj.filteredList)
	} else { // return full float value for csv
		if rowIdx < len(obj.filteredList) {
			retval = fmt.Sprintf("%3.2f", obj.filteredList[rowIdx].value)
		}
	}
	return retval
}

func (obj *kmDetailTable) getCellStrFull(rowIdx int, colIdx int, inputList []*kmDetail) string {
	var retval string
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {
			case KMDTCol0AltName:
				retval = listElem.altName
			case KMDTCol1TS:
				retval = util.ConvertUnixTimeToStr(listElem.timeTS)
			case KMDTCol2KValue:
				retval = util.HumanizeNumber(listElem.value)
			case KMDTCol3KLink:
				switch listElem.zKillOK {
				case model.ZkOK:
					retval = fmt.Sprintf("https://zkillboard.com/kill/%d/", listElem.kmId)
				case model.ZkNOTOK:
					retval = fmt.Sprintf("%d not on zkill", listElem.kmId)
				case model.ZkUnkown:
					retval = fmt.Sprintf("%d unkown", listElem.kmId)
				}
			}
		} else {
			obj.Ctrl.Model.LogObj.Printf("getCellStrFull invalid colidx %d", colIdx)
		}

	} else {
		obj.Ctrl.Model.LogObj.Printf("getCellStrFull invalid row index %d", rowIdx)
	}

	return retval
}

func (obj *kmDetailTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {
		obj.Ctrl.Model.LogObj.Printf("selected row %d col %d", id.Row, id.Col)
		if id.Col == KMDTCol3KLink {
			if id.Row < len(obj.filteredList) {
				if obj.filteredList[id.Row].zKillOK == model.ZkUnkown {
					kmId := obj.filteredList[id.Row].kmId
					km := obj.Ctrl.Model.GetKillsMail(kmId)
					if obj.Ctrl.ZkillOk(kmId) {
						obj.filteredList[id.Row].zKillOK = model.ZkOK
						km.ZK_status = model.ZkOK
					} else {
						obj.filteredList[id.Row].zKillOK = model.ZkNOTOK
						km.ZK_status = model.ZkNOTOK
					}
					result := obj.Ctrl.Model.AddKillmailEntry(km)
					if result != model.DBR_Updated {
						obj.Ctrl.Model.LogObj.Printf("kmDetailTable warning unexepected result %d", result)
					}
				}
				if obj.filteredList[id.Row].zKillOK == model.ZkOK {
					obj.Ctrl.Model.LogObj.Printf("startink link for %d", obj.filteredList[id.Row].kmId)
					url := fmt.Sprintf("https://zkillboard.com/kill/%d/", obj.filteredList[id.Row].kmId)
					util.OpenUrl(url)
				}
			}

		}
	}
}

func (obj *kmDetailTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			if colIdx == KMDTCol2KValue {
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (obj.fulllist[rowIdx].value / 1000000) < s {
						filterOK = false
						break
					}
				}
			} else {
				currentCellStr := obj.getCellStrFull(rowIdx, colIdx, obj.fulllist)
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
func (obj *kmDetailTable) SortCol(colIdx int) {
	sort.Slice(obj.fulllist, func(i, j int) bool {
		var retval bool
		switch colIdx {
		case KMDTCol0AltName:
			retval = obj.fulllist[i].altName >= obj.fulllist[j].altName
		case KMDTCol1TS:
			retval = obj.fulllist[i].timeTS >= obj.fulllist[j].timeTS
		case KMDTCol2KValue:
			retval = obj.fulllist[i].value >= obj.fulllist[j].value
		case KMDTCol3KLink:
			retval = obj.fulllist[i].kmId >= obj.fulllist[j].kmId
		}
		return retval
	})

	obj.Ctrl.Model.LogObj.Printf("sorty by col %d", colIdx)
}

func (obj *OmipGui) createIskLossTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	maxMonth := 12
	var tableObj *widget.Table
	var filterCharName *widget.Entry
	var filterAmount *widget.Entry
	sortByRow := "Character"
	var updateColumnWidth func()
	fullList := obj.Ctrl.Model.GetKillTable(char.CharInfoExt.CooperationId, maxMonth, true)
	var filteredList model.MonthlyTable

	if len(fullList.ValCharPerMon) == 0 {
		return retTable, result
	} else {
		result = true
	}
	filteredCharList := make([]string, 0, 10)

	updateLists := func() {
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

		if sortByRow != "Character" {
			sort.Slice(keyList, func(i, j int) bool {
				name1 := keyList[i]
				name2 := keyList[j]
				compare1 := fullList.ValCharPerMon[name1][sortByRow]
				compare2 := fullList.ValCharPerMon[name2][sortByRow]
				return compare1 >= compare2
			})
		}

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
		if len(filteredCharList) == 0 {
			filteredCharList = append(filteredCharList, "N/A")
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
						color = util.GetColor(filteredList.MaxInMonth[dateStr], filteredList.ValCharPerMon[charName][dateStr], true)
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
	memberList := obj.Ctrl.Model.GetDBCorpMembers(char.CharInfoExt.CooperationId)
	nameMapping := make(map[string]int)
	for _, member := range memberList {
		name, _ := obj.Ctrl.Model.GetStringEntry(member.NameRef)
		nameMapping[name] = member.CharID
	}

	wList := make([]*WindowList, 0, 10)
	tableObj.OnSelected = func(id widget.TableCellID) {
		tableObj.UnselectAll()
		if id.Row < len(filteredCharList) {
			start := time.Now()
			charName := filteredCharList[id.Row]
			if id.Col == 0 {
				obj.Ctrl.Model.LogObj.Printf("%s", charName)
			} else {
				monthIdx := maxMonth - id.Col
				year, month, _ := ymStrPastMonth(monthIdx)
				startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				endTime := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
				windowTitle := fmt.Sprintf("lossed from %s in %s %d", charName, startTime.Month().String(), startTime.Year())
				windowFound := false
				for _, elem := range wList {
					if elem.Title == windowTitle {
						elem.wRef.Show()
						windowFound = true
						break
					}
				}
				if !windowFound {
					obj.Ctrl.Model.LogObj.Printf("%s %d start %d end %d", charName, nameMapping[charName], startTime.Unix(), endTime.Unix())
					list := obj.Ctrl.Model.GetVictimData(char.CharInfoExt.CooperationId, nameMapping[charName], startTime.Unix(), endTime.Unix())
					if len(list) > 0 {
						kmDT := NewKMDT(obj.Ctrl)

						for _, elem := range list {
							var newElem kmDetail
							newElem.timeTS = elem.Time
							newElem.kmId = elem.KMId
							newElem.value = elem.Amount
							newElem.altName = elem.AltName
							newElem.npc = obj.Ctrl.Model.CheckNPCKill(elem.KMId)
							newElem.zKillOK = elem.ZK_status
							kmDT.fulllist = append(kmDT.fulllist, &newElem)
						}
						kmDT.UpdateLists()

						w := fyne.CurrentApp().NewWindow(windowTitle)
						kmDetailWidget := obj.createGenTable2(kmDT, false, false, "")
						w.SetContent(kmDetailWidget)
						height := float32(700)
						if len(kmDT.fulllist) < 20 {
							height = float32(50*len(kmDT.fulllist)) + 100
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
			elapsed := time.Since(start)
			obj.Ctrl.Model.LogObj.Printf("tableObj.OnSelected took %s", elapsed)
		}

	}

	topRowTable := widget.NewTable(
		func() (int, int) { return 1, maxMonth + 1 },
		func() fyne.CanvasObject {
			topButton := widget.NewButton("Character Name", func() {
				sortByRow = "Character"
				updateLists()
			})
			return topButton
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			button := cell.(*widget.Button)
			if id.Col == 0 {
				button.SetText("Character Name")
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
				cellCanvasText.Color = util.GetColor(fullList.MaxAllTime, filteredList.SumInMonth[dateStr2], true)
				cellCanvasText.Text = util.HumanizeNumber(fullList.SumInMonth[dateStr2])
			}
		},
	)

	filterCharName = widget.NewEntry()
	filterCharName.PlaceHolder = "filter char name"
	filterCharName.OnChanged = func(s string) {
		updateLists()
	}

	filterAmount = widget.NewEntry()
	filterAmount.PlaceHolder = "filter millions"
	filterAmount.OnChanged = func(s string) {
		updateLists()
	}
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
	hintLabel := widget.NewLabel("Hint: Click cell to open kill list")
	filtergrid := container.New(layout.NewGridLayout(3), filterCharName, filterAmount, hintLabel)
	bottomGrid := container.New(layout.NewGridLayout(1), bottomRowTable, filtergrid)
	topGrid2 := container.New(layout.NewGridLayout(1), topRowTable)
	mainbox := container.NewBorder(topGrid2, bottomGrid, nil, nil, tableObj)
	return mainbox, result
}

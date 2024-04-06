package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"math"
	"regexp"
	"sort"
	"strconv"
)

func (obj *OmipGui) characterScreen() fyne.CanvasObject {
	obj.CharTabPtr = container.NewDocTabs(make([]*container.TabItem, 0, 5)...)
	charOverView, result := obj.createCharOViewTab()
	var overViewTab *container.TabItem
	if result {
		overViewTab = container.NewTabItem("Overview", charOverView)
		obj.CharTabPtr.Append(overViewTab)
		obj.CharTabPtr.OnClosed = func(item *container.TabItem) {
			if item.Text == overViewTab.Text {
				obj.CharTabPtr.Append(overViewTab)
			}
		}
	}
	// TODO remove workaround for https://github.com/fyne-io/fyne/issues/3169
	obj.CharTabPtr.OnSelected = func(item *container.TabItem) {
		obj.recurseRefresh(item)
	}
	return obj.CharTabPtr
}

const (
	COVCol0CharName = iota
	COVCol1Ticker
	COVCol2ISK
	COVCol3Balance
	COVNumCols
)

type charOverviewTable struct {
	Name    string
	Ticker  string
	Wallet  float64
	Balance float64
	char    *ctrl.EsiChar
	corp    *ctrl.EsiCorp
}

type charGuiTable struct {
	fulllist     []*charOverviewTable
	filteredList []*charOverviewTable
	gui          *OmipGui
	isCorp       bool
	genericTable
}

func (obj *charGuiTable) GetNumRows() int {
	return len(obj.filteredList)
}
func (obj *charGuiTable) getCellStrFromList(rowIdx int, colIdx int, inputList []*charOverviewTable) (retval string, col color.NRGBA) {
	col = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {
			case COVCol0CharName:
				retval = listElem.Name
				col = color.NRGBA{62, 160, 221, 0xff}
				if listElem.char.AuthValid == ctrl.AUTH_STATUS_INVALID {
					col = color.NRGBA{0xff, 0, 0, 0xff}
				}
			case COVCol1Ticker:
				retval = listElem.Ticker
			case COVCol2ISK:
				retval = util.HumanizeNumber(listElem.Wallet)
			case COVCol3Balance:
				retval = util.HumanizeNumber(listElem.Balance)
				if listElem.Balance > 0 {
					col = color.NRGBA{0, 0xff, 0, 0xff}
				} else if listElem.Balance != 0 {
					col = color.NRGBA{0xff, 0, 0, 0xff}
				}
			}
		}
	}
	return
}
func (obj *charGuiTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	txt := ""
	switch colIdx {
	case COVCol0CharName:
		txt = fmt.Sprintf("%d", len(obj.filteredList))
	case COVCol1Ticker:
		mapStr := make(map[string]int)
		count := 0
		for _, elem := range obj.filteredList {
			if elem.Ticker != "" {
				if _, ok := mapStr[elem.Ticker]; !ok {
					mapStr[elem.Ticker] = 1
					count++
				}
			}
		}
		txt = fmt.Sprintf("%d", count)
	case COVCol2ISK:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += elem.Wallet
		}
		txt = util.HumanizeNumber(sum)
	case COVCol3Balance:
		var sum float64
		for _, elem := range obj.filteredList {
			sum += elem.Balance
		}
		if sum > 0 {
			col = color.NRGBA{0, 0xff, 0, 0xff}
		} else if sum != 0 {
			col = color.NRGBA{0xff, 0, 0, 0xff}
		}
		txt = util.HumanizeNumber(sum)
	}

	return txt, col
}
func (obj *charGuiTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	retval := ""
	if rowIdx < len(obj.filteredList) {
		elem := obj.filteredList[rowIdx]
		switch colIdx {
		case COVCol0CharName:
			retval = elem.Name
		case COVCol1Ticker:
			retval = elem.Ticker
		case COVCol2ISK:
			retval = fmt.Sprintf("%3.2f", elem.Wallet)
		case COVCol3Balance:
			retval = fmt.Sprintf("%3.2f", elem.Balance)
		}
	}
	return retval
}

func (obj *charGuiTable) selectCol0(id widget.TableCellID) {
	elem := obj.filteredList[id.Row]
	found := false
	var tabPtr *container.DocTabs
	if obj.isCorp {
		tabPtr = obj.gui.CorpTabPtr
	} else {
		tabPtr = obj.gui.CharTabPtr
	}
	for idx, tab := range tabPtr.Items {
		if obj.isCorp {
			if tab.Text == elem.Ticker {
				tabPtr.Select(tabPtr.Items[idx])
				obj.MainTable.Unselect(id)
				found = true
				break
			}
		} else {
			if tab.Text == elem.Name {
				tabPtr.Select(tabPtr.Items[idx])
				obj.MainTable.Unselect(id)
				found = true
				break
			}
		}

	}
	if !found {
		if len(tabPtr.Items) > 8 {
			lastElemIdx := len(tabPtr.Items) - 1
			tabPtr.Remove(tabPtr.Items[lastElemIdx])
		}
		if obj.isCorp {
			director := obj.gui.getDirectorByCorpId(elem.corp.CooperationId)
			charTab := obj.gui.AddCorpTab(elem.corp, director)
			tabPtr.Select(charTab)
			obj.MainTable.Unselect(id)
		} else {
			charTab := obj.gui.AddCharTab(elem.char)
			tabPtr.Select(charTab)
			obj.MainTable.Unselect(id)
		}
	}
}

func (obj *charGuiTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {
		if id.Row < len(obj.filteredList) {
			switch id.Col {
			case COVCol0CharName:
				obj.selectCol0(id)
			}
		}
	}
}

func (obj *charGuiTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			filterString := true
			if colIdx == COVCol2ISK {
				filterString = false
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (obj.fulllist[rowIdx].Wallet / 1000000) < s {
						filterOK = false
						break
					}
				}
			}
			if colIdx == COVCol3Balance {
				filterString = false
				if s, err := strconv.ParseFloat(currentFilter, 64); err == nil {
					if (math.Abs(obj.fulllist[rowIdx].Balance) / 1000000) < s {
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
func (obj *charGuiTable) SortCol(colIdx int) {
	obj.sortCount++
	sort.Slice(obj.fulllist, func(i, j int) bool {
		var retval bool
		a := obj.fulllist[i]
		b := obj.fulllist[j]
		switch colIdx {
		case COVCol0CharName:
			retval = a.Name >= b.Name
		case COVCol1Ticker:
			retval = a.Ticker >= b.Ticker
		case COVCol2ISK:
			retval = a.Wallet >= b.Wallet
		case COVCol3Balance:
			retval = a.Balance >= b.Balance
		}
		if obj.sortCount%2 == 0 {
			retval = !retval
		}
		return retval
	})
}

func (obj *charGuiTable) init(ctrlObj *ctrl.Ctrl, gui *OmipGui) {
	obj.Ctrl = ctrlObj
	obj.gui = gui
	obj.fulllist = make([]*charOverviewTable, 0, 10)
	obj.header = make([]string, 0, 3)
	obj.header = append(obj.header, "Name")
	obj.header = append(obj.header, "Ticker")
	obj.header = append(obj.header, "Wallet")
	obj.header = append(obj.header, "30d Balance")
	obj.colWidth = make([]float32, 0, 10)
	obj.colWidth = append(obj.colWidth, 400)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 100)
	obj.filter = make([]string, 0, 10)
	for i := 0; i < len(obj.header); i++ {
		obj.filter = append(obj.filter, "")
	}
}

func (obj *charGuiTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	return obj.getCellStrFromList(rowIdx, colIdx, obj.filteredList)
}

func NewCorpTable(ctrlObj *ctrl.Ctrl, gui *OmipGui) *charGuiTable {
	var obj charGuiTable
	obj.isCorp = true
	obj.init(ctrlObj, gui)
	for _, corp := range obj.Ctrl.Esi.EsiCorpList {
		lcorp := corp
		var director *ctrl.EsiChar
		for _, char := range obj.Ctrl.Esi.EsiCharList {
			if char.CharInfoExt.CooperationId == lcorp.CooperationId &&
				char.CharInfoExt.Director {
				director = char
				break
			}
		}
		if director != nil {
			var newCorp charOverviewTable
			newCorp.Name = lcorp.Name
			newCorp.Ticker = lcorp.Ticker
			for i := 1; i < 8; i++ {
				newCorp.Wallet += obj.Ctrl.Model.GetLatestWallets(0, lcorp.CooperationId, i)
			}
			obj.fulllist = append(obj.fulllist, &newCorp)
			newCorp.char = director
			newCorp.corp = lcorp
			balance := obj.Ctrl.Model.GetBalanceOverTime(
				director.CharInfoData.CharacterID, director.CharInfoExt.CooperationId, true, 30)
			newCorp.Balance = balance
		}
	}
	return &obj
}

func NewCharTable(ctrlObj *ctrl.Ctrl, gui *OmipGui) *charGuiTable {
	var obj charGuiTable
	obj.init(ctrlObj, gui)
	for _, char := range obj.Ctrl.Esi.EsiCharList {
		var newChar charOverviewTable
		newChar.Name = char.CharInfoData.CharacterName
		if char.CharInfoExt.CooperationId > ctrl.EsiCorpIdLimit {
			corp := obj.Ctrl.GetCorp(char)
			if corp != nil {
				newChar.Ticker = corp.Ticker
			}
		}
		wallet := obj.Ctrl.Model.GetLatestWallets(char.CharInfoData.CharacterID, 0, 0)
		balance := obj.Ctrl.Model.GetBalanceOverTime(
			char.CharInfoData.CharacterID, char.CharInfoExt.CooperationId, false, 30)
		newChar.Wallet = wallet
		newChar.char = char
		newChar.Balance = balance
		obj.fulllist = append(obj.fulllist, &newChar)
	}
	copy(obj.filteredList, obj.fulllist)

	return &obj
}

func (obj *OmipGui) createCharOViewTab() (retTable fyne.CanvasObject, result bool) {
	if len(obj.Ctrl.Esi.EsiCharList) > 0 {
		table := NewCharTable(obj.Ctrl, obj)
		table.sortCount = 1
		table.SortCol(COVCol0CharName)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, false, true, "")
		result = true
	}

	return retTable, result
}

func (obj *OmipGui) AddCharTab(char *ctrl.EsiChar) *container.TabItem {
	charSubTabs := obj.CreateCharGui(char, true)
	newTab := container.NewTabItem(char.CharInfoData.CharacterName, charSubTabs)
	obj.CharTabPtr.Append(newTab)

	return newTab
}

type RefbutList struct {
	RefID     string
	RefButton *widget.Button
}

func (obj *OmipGui) CreateOverViewTable(char *ctrl.EsiChar, corp bool) (fyne.CanvasObject, []*RefbutList) {
	refList := make([]*RefbutList, 0, 10)
	vbox := container.NewVBox(make([]fyne.CanvasObject, 0, 10)...)
	if char.UpdateFlags.Journal {
		sum := obj.Ctrl.Model.GetBalanceOverTime(
			char.CharInfoData.CharacterID, char.CharInfoExt.CooperationId, corp, 30)
		if sum != 0 {
			var newRBL RefbutList
			text := fmt.Sprintf("wallet change last 30 days %s", util.HumanizeNumber(sum))
			newRBL.RefID = "Journal"
			newRBL.RefButton = widget.NewButton(text, func() {
			})
			refList = append(refList, &newRBL)
			vbox.Objects = append(vbox.Objects, newRBL.RefButton)
		}
	}
	if char.UpdateFlags.IndustryJobs {
		entityID := char.CharInfoData.CharacterID
		if corp {
			entityID = char.CharInfoExt.CooperationId
		}
		earliestEndDate, overdueSum := obj.Ctrl.Model.GetNextPendingJob(entityID, corp)
		if earliestEndDate != nil {
			retval2, _ := util.GetTimeDiffStringFromTS(earliestEndDate.EndDate)
			text := fmt.Sprintf("next industry job ends in %s", retval2)
			var newRBL RefbutList
			newRBL.RefID = "Industry"
			newRBL.RefButton = widget.NewButton(text, func() {
			})
			refList = append(refList, &newRBL)
			vbox.Objects = append(vbox.Objects, newRBL.RefButton)
		}
		if overdueSum != 0 {
			text := fmt.Sprintf("currently %d pending industry jobs", overdueSum)
			var newRBL RefbutList
			newRBL.RefID = "Industry"
			newRBL.RefButton = widget.NewButton(text, func() {
			})
			refList = append(refList, &newRBL)
			vbox.Objects = append(vbox.Objects, newRBL.RefButton)
		}
	}

	return vbox, refList
}

func (obj *OmipGui) CreateCharGui(char *ctrl.EsiChar, isTab bool) fyne.CanvasObject {
	localchar := char
	var portrait *canvas.Image
	portrait = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	portrait.File = localchar.ImageFile

	overViewtable, refButtonlist := obj.CreateOverViewTable(char, false)
	charMain := container.NewTabItem("Overview", container.NewBorder(portrait, nil, nil, nil, overViewtable))
	industryTabObj, needed := obj.createIndustryTab(localchar, false)
	charSubTabs := container.NewAppTabs(charMain)
	if obj.Ctrl.Model.DebugFlag {
		debugTab := obj.CreateCharDebugTab(char)
		charDebug := container.NewTabItem("debug", container.NewVScroll(debugTab))
		charSubTabs.Append(charDebug)
	}
	if needed {
		industryTab := container.NewTabItem("Industry", industryTabObj)
		obj.assignRefButton("Industry", refButtonlist, charSubTabs, len(charSubTabs.Items))
		charSubTabs.Append(industryTab)

	}
	ctrTabObj, ctrNeeded := obj.createContractTab(localchar, false)
	if ctrNeeded {
		ctrTab := container.NewTabItem("Contracts", ctrTabObj)
		charSubTabs.Append(ctrTab)
	}
	jourTabObj, jourNeeded := obj.createJournalTab(localchar, false)
	if jourNeeded {
		jourTab := container.NewTabItem("Journal", jourTabObj)
		obj.assignRefButton("Journal", refButtonlist, charSubTabs, len(charSubTabs.Items))
		charSubTabs.Append(jourTab)
	}

	charSubTabs.SetTabLocation(container.TabLocationTop)

	// TODO remove workaround for https://github.com/fyne-io/fyne/issues/3169
	charSubTabs.OnSelected = func(item *container.TabItem) {
		obj.recurseRefresh(item)
	}
	return charSubTabs
}

func (obj *OmipGui) assignRefButton(ID string, refList []*RefbutList, tab *container.AppTabs, idx int) {
	for _, elem := range refList {
		if ID == elem.RefID {
			elem.RefButton.OnTapped = func() {
				tab.SelectIndex(idx)
				tab.Refresh()
			}
		}
	}
}

func (obj *OmipGui) CreateCharDebugTab(localchar *ctrl.EsiChar) fyne.CanvasObject {
	statusBtn := widget.NewButton("status", func() {
		obj.Ctrl.CheckServerUp(localchar)
	})

	contractBtn := widget.NewButton("Contract", func() {
		obj.Ctrl.UpdateContracts(localchar, false)
	})
	contractItemsBtn := widget.NewButton("Contract Items", func() {
		obj.Ctrl.UpdateContractItems(localchar, false)
	})

	indudstryBtn := widget.NewButton("Industry", func() {
		obj.Ctrl.UpdateIndustry(localchar, false)
	})
	journalBtn := widget.NewButton("Journal", func() {
		obj.Ctrl.UpdateJournal(localchar, false, 0)
	})
	orderBtn := widget.NewButton("Orders", func() {
		obj.Ctrl.UpdateOrders(localchar, false)
	})
	transactionsBtn := widget.NewButton("Transactions", func() {
		obj.Ctrl.UpdateTransaction(localchar, false)
	})
	killmailsBtn := widget.NewButton("KillMails", func() {
		obj.Ctrl.UpdateKillMails(localchar, false)
	})
	walletBtn := widget.NewButton("Wallet", func() {
		obj.Ctrl.UpdateWallet(localchar, false)
	})
	notificationBtn := widget.NewButton("Notification", func() {
		obj.Ctrl.UpdateNotifications(localchar, false)
	})
	marketPriceBtn := widget.NewButton("Market Prices", func() {
		obj.Ctrl.UpdateMarket(localchar, false)
	})
	directorBtn := widget.NewButton("check director", func() {
		obj.Ctrl.CheckIfDirector(localchar)
	})
	mailButton := widget.NewButton("check mail", func() {
		obj.Ctrl.UpdateMailLabels(localchar, false)
	})

	debugVbox := container.NewVBox(statusBtn, contractBtn, contractItemsBtn, indudstryBtn, journalBtn, orderBtn, transactionsBtn, killmailsBtn, walletBtn, notificationBtn, marketPriceBtn, directorBtn, mailButton)
	return debugVbox
}

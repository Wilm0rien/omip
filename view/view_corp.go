package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"time"
)

func (obj *OmipGui) updateCorpScreen() {
	obj.CorpTabPtr = container.NewDocTabs(make([]*container.TabItem, 0, 5)...)
	corpOverView, result := obj.createCopOViewTab()
	if result {
		newTab := container.NewTabItem("Overview", corpOverView)
		obj.CorpTabPtr.Append(newTab)
		obj.CorpTabPtr.OnClosed = func(item *container.TabItem) {
			if item.Text == newTab.Text {
				obj.CorpTabPtr.Append(newTab)
			}
		}
		obj.CorpTabPtr.OnSelected = func(item *container.TabItem) {
			if item.Text == newTab.Text {
				corpOverView.Refresh()
			}
		}
		// TODO remove workaround for https://github.com/fyne-io/fyne/issues/3169
		obj.CorpTabPtr.OnSelected = func(item *container.TabItem) {
			obj.recurseRefresh(item)
		}
		obj.CorpTabPtr.Refresh()
	}
}

func (obj *OmipGui) createCopOViewTab() (retTable fyne.CanvasObject, result bool) {
	table := NewCorpTable(obj.Ctrl, obj)
	if len(table.fulllist) > 0 {
		table.sortCount = 1
		table.SortCol(COVCol0CharName)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, false, true, "")
		result = true
	}
	return retTable, result
}

// return year month string with n month in the past (n = 6 = "yy-dd" 6 month from now)
func ymStrPastMonth(n int) (int, int, string) {
	tm := time.Now()
	year, month, _ := tm.Date()
	curMonth := int(month) - n
	curYear := year
	if curMonth <= 0 {
		curMonth = 12 + curMonth
		curYear = year - 1
	}
	return curYear, curMonth, fmt.Sprintf("%02d-%02d", curYear-2000, curMonth)
}

func (obj *OmipGui) createPapTab2(corpID int, maxMonth int) (retTable fyne.CanvasObject, result bool) {
	fullList := obj.Ctrl.Model.GetPapTable(corpID)
	if fullList.MaxAllTime == 0 {
		result = false
		return retTable, result
	} else {
		table := NewMonthlyTable(obj.Ctrl, fullList, maxMonth, false)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, true, true, "Paps")
		result = true
	}
	return retTable, result
}

func (obj *OmipGui) createKillMailTab(corpID int, maxMonth int, losses bool) (retTable fyne.CanvasObject, result bool) {
	fullList := obj.Ctrl.Model.GetKillTable(corpID, maxMonth, losses)
	if fullList.MaxAllTime == 0 {
		result = false
		return retTable, result
	} else {
		table := NewMonthlyTable(obj.Ctrl, fullList, maxMonth, false)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, true, true, "Kills")
		result = true
	}
	return retTable, result
}
func (obj *OmipGui) createBountyTab(corpID int, maxMonth int) (retTable fyne.CanvasObject, result bool) {
	fullList := obj.Ctrl.Model.GetBountyTable(corpID)
	if fullList.MaxAllTime == 0 {
		result = false
		return retTable, result
	} else {
		table := NewMonthlyTable(obj.Ctrl, fullList, maxMonth, true)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, true, true, "Filter Millions")
		result = true
	}
	return retTable, result
}

func (obj *OmipGui) getDirectorByCorpId(corpID int) *ctrl.EsiChar {
	var director *ctrl.EsiChar
	var lcorp *ctrl.EsiCorp
	for _, corp := range obj.Ctrl.Esi.EsiCorpList {
		if corp.CooperationId == corpID {
			lcorp = corp
		}
	}
	if lcorp != nil {
		for _, char := range obj.Ctrl.Esi.EsiCharList {
			if char.CharInfoExt.CooperationId == lcorp.CooperationId &&
				char.CharInfoExt.Director {
				director = char
				break
			}
		}
	}

	return director
}

func (obj *OmipGui) CorpUpdate() {
	obj.updateCorpScreen()
	if len(obj.CorpTabPtr.Items) > 1 {
		obj.CorpTabPtr.Select(obj.CorpTabPtr.Items[0])
	}
}

func (obj *OmipGui) AddCorpTab(corp *ctrl.EsiCorp, director *ctrl.EsiChar) *container.TabItem {
	corpSubTabs := obj.CreateCorpGui(corp, director, true)
	newTab := container.NewTabItemWithIcon(corp.Ticker, obj.getIconResource(corp.ImageFile), corpSubTabs)
	obj.CorpTabPtr.Append(newTab)

	return newTab
}

func (obj *OmipGui) CreateCorpDebugTab(director *ctrl.EsiChar, corp *ctrl.EsiCorp) fyne.CanvasObject {
	contractBtn := widget.NewButton("Contract", func() {
		obj.Ctrl.UpdateContracts(director, true)
	})
	contractItemsBtn := widget.NewButton("Contract Items", func() {
		obj.Ctrl.UpdateContractItems(director, true)
	})
	indudstryBtn := widget.NewButton("Industry", func() {
		obj.Ctrl.UpdateIndustry(director, true)
	})
	journalBtn := widget.NewButton("Journal", func() {
		for i := 1; i <= 7; i++ {
			obj.Ctrl.UpdateJournal(director, true, i)
		}
	})
	orderBtn := widget.NewButton("Orders", func() {
		obj.Ctrl.UpdateOrders(director, true)
	})
	transactionsBtn := widget.NewButton("Transactions", func() {
		obj.Ctrl.UpdateTransaction(director, true)
	})
	killmailsBtn := widget.NewButton("KillMails", func() {
		obj.Ctrl.UpdateKillMails(director, true)
	})
	killmailsBtnLastMotn := widget.NewButton("KillMailsLastMonth", func() {
		obj.Ctrl.UpdateSkippListLastMonth(director, true)
		obj.Ctrl.UpdateKillMails(director, true)
	})
	structureBtn := widget.NewButton("Structure", func() {
		obj.Ctrl.UpdateStructures(director, true)
	})
	walletBtn := widget.NewButton("Wallet", func() {
		obj.Ctrl.UpdateWallet(director, true)
	})
	/** TODO remove paps
	papBtn := widget.NewButton("Paps", func() {
		if aDash, ok := obj.Ctrl.ADash[director.CharInfoExt.CooperationId]; ok {
			if aDash.Login() {
				aDash.GetPapLinks()
			} else {
				obj.AddLogEntry(fmt.Sprintf("Adash Login failed"))
			}
		}
	})*/
	crashbutton := widget.NewButton("test crash", func() {

	})
	membersBtn := widget.NewButton("Member", func() {
		obj.Ctrl.UpdateCorpMembers(director, true)
	})
	box := container.NewVBox(
		contractBtn, contractItemsBtn, indudstryBtn, journalBtn, orderBtn, transactionsBtn,
		killmailsBtn, killmailsBtnLastMotn, structureBtn, walletBtn, membersBtn, crashbutton)
	/** TODO remove paps
	if IsImperium(corp.AllianceId) {
		box.Objects = append(box.Objects, papBtn)
	}
	*/
	return box
}

func IsImperium(allianceID int) (retval bool) {
	/*
		alliances := []int{
			1354830081, 150097440, 99005518, 99003995, 99010664, 99009163, 499005583, 99009169, 131511956, 99010751,
			99009748, 99007203, 99009805, 99006751, 99009882, 1900696668, 99007362, 1220922756, 99007916, 99004425}
		for _, id := range alliances {
			if allianceID == id {
				retval = true
				break
			}
		}*/
	/* TODO remove PAP statistic or try to connect to dank fleet dash board:
	https://goonfleet.com/index.php/topic/329424-a-dank-fleet-dashboard-you/ */

	return false
}

func (obj *OmipGui) CreateCorpGui(corp *ctrl.EsiCorp, director *ctrl.EsiChar, isTab bool) fyne.CanvasObject {
	var corpLogo *canvas.Image
	corpLogo = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	corpLogo.File = corp.ImageFile

	vbox := container.NewVBox(make([]fyne.CanvasObject, 0, 10)...)
	/* TODO remove adash
	if IsImperium(corp.AllianceId) {
		aDashBtn := obj.aDashButton(corp.CooperationId)
		vbox.Objects = append(vbox.Objects, aDashBtn)
	}

	*/
	if len(obj.CorpTabPtr.Items) > 0 && isTab {
		vbox.Objects = append(vbox.Objects)
	}
	overViewtable, refButtonlist := obj.CreateOverViewTable(director, true)

	maintabLayout := container.NewBorder(corpLogo, vbox, nil, nil, overViewtable)
	corpTab := container.NewTabItem("Overview", maintabLayout)
	altmapTab := container.NewTabItem("ALTs", obj.createAltTab(director))
	corpSubTabs := container.NewAppTabs(corpTab)
	corpSubTabs.Append(altmapTab)
	if obj.DebugFlag {
		debugVbox := obj.CreateCorpDebugTab(director, corp)
		corpDebug := container.NewTabItem("debug", container.NewVScroll(debugVbox))
		corpSubTabs.Append(corpDebug)
	}

	killTab, needed := obj.createKillMailTab(corp.CooperationId, 12, false)
	if needed {
		killTab2 := container.NewTabItemWithIcon("Kill Count", killreportIcon, killTab)
		corpSubTabs.Append(killTab2)
	}
	bountyTab, needed := obj.createBountyTab(corp.CooperationId, 12)
	if needed {
		bountyTab2 := container.NewTabItemWithIcon("Bounties", killreportIcon, bountyTab)
		corpSubTabs.Append(bountyTab2)
	}
	/* TODO fix pap statistcs
	papTab2, needed := obj.createPapTab2(corp.CooperationId, 12)
	if needed {
		papTab3 := container.NewTabItemWithIcon("Paps", achievementsIcon, papTab2)
		corpSubTabs.Append(papTab3)
	}
	*/
	iskLossObj, needed := obj.createIskLossTab(director, true)
	if needed {
		killtab := container.NewTabItemWithIcon("ISK Loss", terminateIcon, iskLossObj)
		corpSubTabs.Append(killtab)
	}
	induTabObj, needed := obj.createIndustryTab(director, true)
	if needed {
		industryTab := container.NewTabItem("Industry", induTabObj)
		obj.assignRefButton("Industry", refButtonlist, corpSubTabs, len(corpSubTabs.Items))
		corpSubTabs.Append(industryTab)

	}
	ctrTabObj, ctrNeeded := obj.createContractTab(director, true)
	if ctrNeeded {
		ctrTab := container.NewTabItem("Contracts", ctrTabObj)
		corpSubTabs.Append(ctrTab)
	}
	strucTabObj, strucNeeded := obj.createStructureTab(director, true)
	if strucNeeded {
		strucTab := container.NewTabItem("Structures", strucTabObj)
		corpSubTabs.Append(strucTab)
	}
	jourTabObj, jourNeeded := obj.createJournalTab(director, true)
	if jourNeeded {
		jourTab := container.NewTabItem("Journal", jourTabObj)
		obj.assignRefButton("Journal", refButtonlist, corpSubTabs, len(corpSubTabs.Items))
		corpSubTabs.Append(jourTab)
	}

	// TODO remove workaround for https://github.com/fyne-io/fyne/issues/3169
	corpSubTabs.OnSelected = func(item *container.TabItem) {
		obj.recurseRefresh(item)
	}
	corpSubTabs.SetTabLocation(container.TabLocationTop)

	return corpSubTabs
}

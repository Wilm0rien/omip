package view

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"regexp"
)

func (obj *OmipGui) createAltTab(director *ctrl.EsiChar) fyne.CanvasObject {
	memberList := obj.Ctrl.Model.GetDBCorpMembers(director.CharInfoExt.CooperationId)
	var safeBtn *widget.Button
	var restoreBtn *widget.Button
	selctedLabel := widget.NewLabel("ALTs of selected main")
	nameMapping := make(map[int]string)
	altList := make([]*model.DBcorpMember, 0, 10)
	mainList := make([]*model.DBcorpMember, 0, 10)
	currenAltList := make([]*model.DBcorpMember, 0, 10)
	unassignedAltList := make([]*model.DBcorpMember, 0, 10)
	var lastSelectedMainID int
	var mainListWidget *widget.List
	lastSelectedMainID = 0xFFFF
	for _, member := range memberList {
		name, _ := obj.Ctrl.Model.GetStringEntry(member.NameRef)
		nameMapping[member.CharID] = name
	}
	var updateLists func()
	var currentAltListWidget *widget.List
	updateCurrentAltList := func(id int) {
		currenAltList = currenAltList[:0]
		for _, member := range memberList {
			if member.CharID != mainList[id].CharID && member.MainID == mainList[id].CharID {
				currenAltList = append(currenAltList, member)
			}
		}
		currentAltListWidget.Refresh()
	}
	currentAltListWidget = widget.NewList(
		func() int {
			return len(currenAltList)
		},
		func() fyne.CanvasObject {
			unassignBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {})
			return container.New(layout.NewHBoxLayout(), widget.NewLabel("Template Object"), unassignBtn)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			unassignBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				for _, member := range memberList {
					if member.CharID == currenAltList[id].CharID {
						member.MainID = 0
						break
					}
				}
				updateLists()
				updateCurrentAltList(lastSelectedMainID)
			})
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(nameMapping[currenAltList[id].CharID])
			item.(*fyne.Container).Objects[1] = unassignBtn
		})

	unassignedAltListWidget := widget.NewList(
		func() int {
			return len(unassignedAltList)
		},
		func() fyne.CanvasObject {
			assignBtn := widget.NewButtonWithIcon("ALT", theme.NavigateNextIcon(), func() {})
			mainBtn := widget.NewButtonWithIcon("MAIN", theme.NavigateBackIcon(), func() {})
			return container.New(layout.NewHBoxLayout(), widget.NewLabel("Template Object"), mainBtn, assignBtn)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(unassignedAltList) {
				id = len(unassignedAltList) - 1
			}
			assignBtn := widget.NewButtonWithIcon("ALT", theme.NavigateNextIcon(), func() {
				if lastSelectedMainID < len(mainList) {
					for _, member := range memberList {
						if member.CharID == unassignedAltList[id].CharID {
							member.MainID = mainList[lastSelectedMainID].CharID
							member.Updated = true
							break
						}
					}
					updateCurrentAltList(lastSelectedMainID)
				}
				updateLists()
			})
			mainBtn := widget.NewButtonWithIcon("MAIN", theme.NavigateBackIcon(), func() {
				currentCharID := unassignedAltList[id].CharID
				for _, member := range memberList {
					if member.CharID == currentCharID {
						member.MainID = member.CharID
						member.Updated = true
						break
					}
				}
				updateLists()
				for idx, member := range mainList {
					if member.CharID == currentCharID {
						mainListWidget.Select(idx)
						lastSelectedMainID = idx
						updateCurrentAltList(lastSelectedMainID)
						charname, _ := nameMapping[mainList[idx].CharID]
						selctedLabel.SetText(fmt.Sprintf("ALTs of %s", charname))
						break
					}
				}

			})
			name := "error"
			if id < len(unassignedAltList) {
				name = nameMapping[unassignedAltList[id].CharID]
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(name)
			item.(*fyne.Container).Objects[1] = mainBtn
			item.(*fyne.Container).Objects[2] = assignBtn
		})
	altListWidget := widget.NewList(
		func() int {
			return len(altList)
		},
		func() fyne.CanvasObject {
			unassignBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {})
			return container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Template Object"),
				widget.NewIcon(theme.NavigateNextIcon()),
				widget.NewLabel("Template Object"),
				unassignBtn)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(altList) {
				id = len(altList) - 1
			}
			unassignBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				for _, member := range memberList {
					if member.CharID == altList[id].CharID {
						member.MainID = 0
						member.Updated = true
						break
					}
				}
				if lastSelectedMainID != 0xFFFF {
					updateCurrentAltList(lastSelectedMainID)
				}
				updateLists()
			})
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(nameMapping[altList[id].CharID])
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(nameMapping[altList[id].MainID])
			item.(*fyne.Container).Objects[3] = unassignBtn
		})
	mainListWidget = widget.NewList(
		func() int {
			return len(mainList)
		},
		func() fyne.CanvasObject {
			altBtn := widget.NewButton("ALT", func() {})
			return container.New(layout.NewHBoxLayout(), widget.NewLabel("Template Object"), altBtn)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			numAlts := 0
			if id >= len(mainList) {
				id = len(mainList) - 1
			}
			for _, member := range memberList {
				if member.CharID != mainList[id].CharID && member.MainID == mainList[id].CharID {
					numAlts++
				}
			}
			outString := fmt.Sprintf("%s (%d)", nameMapping[mainList[id].CharID], numAlts)
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(outString)
			altBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				// unassign alts first
				for _, member := range memberList {
					if member.CharID != mainList[id].CharID && member.MainID == mainList[id].CharID {
						member.MainID = 0
						member.Updated = true
					}
				}

				for _, member := range memberList {
					if member.CharID == mainList[id].CharID {
						member.MainID = 0
						member.Updated = true
						break
					}
				}
				currenAltList = currenAltList[:0]
				currentAltListWidget.Refresh()
				updateLists()
				if id >= len(mainList) {
					if len(mainList) > 0 {
						newId := len(mainList) - 1
						mainListWidget.Select(newId)
						lastSelectedMainID = newId
						charname, _ := nameMapping[mainList[newId].CharID]
						selctedLabel.SetText(fmt.Sprintf("ALTs of %s", charname))
						updateCurrentAltList(lastSelectedMainID)
					} else {
						selctedLabel.SetText(fmt.Sprintf("ALTs of %s", "selected MAIN"))
					}
				} else {
					mainListWidget.Select(id)
					lastSelectedMainID = id
					charname, _ := nameMapping[mainList[id].CharID]
					selctedLabel.SetText(fmt.Sprintf("ALTs of %s", charname))
					updateCurrentAltList(lastSelectedMainID)
				}

			})
			item.(*fyne.Container).Objects[1] = altBtn
		})

	mainListWidget.OnSelected = func(id widget.ListItemID) {
		lastSelectedMainID = id
		updateCurrentAltList(id)
		charname, _ := nameMapping[mainList[id].CharID]
		selctedLabel.SetText(fmt.Sprintf("ALTs of %s", charname))
	}

	filterAlts := widget.NewEntry()
	filterAlts.PlaceHolder = "filter ALTs"
	filterUAlts := widget.NewEntry()
	filterUAlts.PlaceHolder = "filter unassigend ALTs"
	filterMains := widget.NewEntry()
	filterMains.PlaceHolder = "filter Mains"
	safeBtn = widget.NewButton("Safe to DB", func() {
		isUnassigned := false
		for _, member := range memberList {
			if _, ok := nameMapping[member.MainID]; !ok {
				isUnassigned = true
			}
		}
		if isUnassigned {
			err := errors.New("cannot safe with unassigend alts")
			dialog.ShowError(err, obj.WindowPtr)
		} else {
			obj.Ctrl.Model.SetDBCorpMembers(memberList)
			obj.UpdateGui()
		}
	})
	restoreBtn = widget.NewButton("Restore from DB", func() {
		memberList = memberList[:0]
		memberList2 := obj.Ctrl.Model.GetDBCorpMembers(director.CharInfoExt.CooperationId)
		for _, member := range memberList2 {
			memberList = append(memberList, member)
		}
		updateLists()
	})
	updateLists = func() {
		altList = altList[:0]
		mainList = mainList[:0]
		unassignedAltList = unassignedAltList[:0]
		dbChanged := false
		for _, member := range memberList {
			if member.Updated == true {
				dbChanged = true
			}
			if member.CharID != member.MainID && member.MainID != 0 {
				fAltsMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterAlts.Text), nameMapping[member.CharID])
				if fAltsMatch {
					altList = append(altList, member)
				}
			} else {
				if member.MainID == 0 {
					fuAltMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterUAlts.Text), nameMapping[member.CharID])
					if fuAltMatch {
						unassignedAltList = append(unassignedAltList, member)
					}
				} else {
					fMainMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterMains.Text), nameMapping[member.CharID])
					if fMainMatch {
						mainList = append(mainList, member)
					}
				}
			}
		}
		unassignedAltListWidget.Refresh()
		mainListWidget.Refresh()
		altListWidget.Refresh()
		if dbChanged == true {
			safeBtn.Show()
			restoreBtn.Show()
		} else {
			safeBtn.Hide()
			restoreBtn.Hide()
		}
	}
	updateLists()
	filterAlts.OnChanged = func(s string) {
		updateLists()
	}
	filterUAlts.OnChanged = func(s string) {
		updateLists()
	}
	filterMains.OnChanged = func(s string) {
		updateLists()
	}
	resetFilterBtn := widget.NewButton("Reset Filters", func() {
		filterAlts.SetText("")
		filterUAlts.SetText("")
		filterMains.SetText("")
	})

	topGrid := container.NewGridWithColumns(4)
	topGrid.Objects = append(topGrid.Objects, widget.NewLabel("All mains"))
	topGrid.Objects = append(topGrid.Objects, widget.NewLabel("All unassigned characters"))
	topGrid.Objects = append(topGrid.Objects, selctedLabel)
	topGrid.Objects = append(topGrid.Objects, widget.NewLabel("All assigned alts"))

	mainAltGrid := container.New(
		layout.NewGridLayout(4), mainListWidget, unassignedAltListWidget, currentAltListWidget, altListWidget)

	bottomGrid := container.NewGridWithColumns(4)
	bottomGrid.Objects = append(bottomGrid.Objects, filterMains)
	bottomGrid.Objects = append(bottomGrid.Objects, filterUAlts)
	bottomGrid.Objects = append(bottomGrid.Objects, filterAlts)

	bottomGrid.Objects = append(bottomGrid.Objects, container.NewHBox(resetFilterBtn, safeBtn, restoreBtn))

	return container.NewBorder(topGrid, bottomGrid, nil, nil, mainAltGrid)
}

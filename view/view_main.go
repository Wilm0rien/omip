package view

import (
	"bufio"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

//	obj.BlueColor = color.NRGBA{62, 160, 221, 0xff}

type WindowList struct {
	Title string
	wRef  fyne.Window
}

type OmipGui struct {
	Ctrl        *ctrl.Ctrl
	AppPtr      fyne.App
	WindowPtr   fyne.Window
	Progress    *widget.ProgressBar
	TabPtr      *container.AppTabs
	CharTabPtr  *container.DocTabs
	CorpTabPtr  *container.DocTabs
	NotifyEntry *widget.Entry
	NotifyText  string
	DebugFlag   bool
}

func NewOmipGui(ctrl *ctrl.Ctrl, app fyne.App, debug bool, version string) *OmipGui {
	var obj OmipGui
	obj.DebugFlag = debug
	obj.Ctrl = ctrl
	obj.Ctrl.AuthCb = obj.AddEsiKey
	obj.Ctrl.AddLogCB = obj.AddLogEntry
	obj.AppPtr = app
	obj.AppPtr.SetIcon(resourceLogoPng)
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Add Character...", func() {
				obj.Ctrl.OpenAuthInBrowser()
			}),
			fyne.NewMenuItem("Update all data", func() {
				obj.UpdateAllData()
			}),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("GFX Settings", func() {
				w := obj.AppPtr.NewWindow("OMIP GFX Settings")
				w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
				w.Resize(fyne.NewSize(480, 480))
				w.Show()
			}),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("LICENSE", func() {
				d := dialog.NewCustom("OMIP COPYRIGHT NOTICE", "OK", obj.makeLicenseDialog(), obj.WindowPtr)
				d.Show()
			}),
		))

	obj.TabPtr = container.NewAppTabs(make([]*container.TabItem, 0, 5)...)

	obj.TabPtr.SetTabLocation(container.TabLocationLeading)
	maintab := container.New(
		layout.NewGridLayout(1), obj.TabPtr)
	obj.WindowPtr = obj.AppPtr.NewWindow(fmt.Sprintf("OMIP - An Eve Online Data Aggregator %s", version))
	obj.WindowPtr.SetMainMenu(menu)
	obj.WindowPtr.SetContent(maintab)
	obj.WindowPtr.Resize(fyne.NewSize(1790, 800))
	obj.WindowPtr.SetMaster()
	obj.AppPtr.Settings().SetTheme(theme.DarkTheme())
	obj.AppPtr.Settings().Scale()
	return &obj
}

func (obj *OmipGui) makeLicenseDialog() (result fyne.CanvasObject) {
	omip := widget.NewLabel(
		`Copyright (C) 2021 Christian Wilmes

This program is free software: you can redistribute it and/or modify  it under the terms of the GNU General Public 
License as published by  the Free Software Foundation, either version 3 of the License, or (at your option) any later 
version. This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the 
implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more 
details.

For third party licenses see link below:
`)

	thirdPartyLink := "https://github.com/Wilm0rien/OMIP/third_party.md"
	link3, _ := url.Parse(thirdPartyLink)
	thirdPartyUrl := widget.NewHyperlink(thirdPartyLink, link3)
	return container.NewVBox(omip, thirdPartyUrl)
}

func (obj *OmipGui) UpdateAllData() {
	if len(obj.Ctrl.Esi.EsiCharList) == 0 {
		err := errors.New("no characters registered to update data")
		dialog.ShowError(err, obj.WindowPtr)
		return
	}
	if obj.Ctrl.Model.DebounceEntryExists("update_string") {
		err := errors.New("please wait 5 minutes between updates")
		dialog.ShowError(err, obj.WindowPtr)
		return
	}

	serverStatus := obj.Ctrl.CheckServerUp(obj.Ctrl.Esi.EsiCharList[0])
	if !serverStatus {
		err := errors.New("esi server is starting up or is not reachable. cannot update data")
		dialog.ShowError(err, obj.WindowPtr)
		return
	}

	obj.TabPtr.SelectIndex(0)
	obj.Ctrl.Model.AddDebounceEntry("update_string")
	prog := obj.Progress
	go func() {
		type UpdateFunc func(char *ctrl.EsiChar, corp bool)
		Updates := make([]UpdateFunc, 0, 5)
		Updates = append(Updates, obj.Ctrl.UpdateContracts)
		Updates = append(Updates, obj.Ctrl.UpdateContractItems)
		Updates = append(Updates, obj.Ctrl.UpdateIndustry)
		Updates = append(Updates, obj.Ctrl.UpdateKillMails)
		Updates = append(Updates, obj.Ctrl.UpdateWallet)
		Updates = append(Updates, obj.Ctrl.UpdateCorpMembers)
		Updates = append(Updates, obj.Ctrl.UpdateStructures)
		Updates = append(Updates, obj.Ctrl.UpdateNotifications)
		Updates = append(Updates, obj.Ctrl.UpdateTransaction)
		Updates = append(Updates, obj.Ctrl.UpdateOrders)

		totalItems := (len(obj.Ctrl.Esi.EsiCharList) + len(obj.Ctrl.Esi.EsiCorpList) + 1) * len(Updates)
		// add 1 journal request per character and 7 journal requests per corp
		totalItems += len(obj.Ctrl.Esi.EsiCharList) + (len(obj.Ctrl.Esi.EsiCorpList) * 7)
		var itemCount int
		if len(obj.Ctrl.Esi.EsiCharList) > 0 {
			obj.Ctrl.UpdateMarket(obj.Ctrl.Esi.EsiCharList[0], false)
		}
		for _, char := range obj.Ctrl.Esi.EsiCharList {
			// NOTE: the journal has to be updated first to update the journal_links table
			// this is because only contracts with journal links are identified as relevant for being stored
			obj.Ctrl.UpdateJournal(char, false, 0)
			itemCount++
			for _, updateFunc := range Updates {
				updateFunc(char, false)
				itemCount++
				prog.SetValue(float64(itemCount) / float64(totalItems))
			}

			prog.SetValue(float64(itemCount) / float64(totalItems))
		}
		for _, corp := range obj.Ctrl.Esi.EsiCorpList {
			director := obj.Ctrl.GetCorpDirector(corp.CooperationId)
			if director != nil {
				for _, updateFunc := range Updates {
					updateFunc(director, true)
					itemCount++
					prog.SetValue(float64(itemCount) / float64(totalItems))
				}
				for i := 1; i <= 7; i++ {
					obj.Ctrl.UpdateJournal(director, true, i)
					itemCount++
					prog.SetValue(float64(itemCount) / float64(totalItems))
				}
				if _, ok := obj.Ctrl.ADash[corp.CooperationId]; !ok {
					if obj.Ctrl.Model.ADashAuthExists(corp.CooperationId) {
						email, pw, _ := obj.Ctrl.Model.GetAuth(corp.CooperationId)
						ticker := obj.Ctrl.Model.GetCorpTicker(corp.CooperationId)
						obj.Ctrl.ADash[corp.CooperationId] = ctrl.NewADashClient(email, pw, ticker, obj.Ctrl.Model, corp.CooperationId)
						obj.Ctrl.ADash[corp.CooperationId].AddLogCB = obj.Ctrl.AddLogEntry
						obj.Ctrl.ADash[corp.CooperationId].Username = email
						obj.Ctrl.ADash[corp.CooperationId].Password = pw
					}
				}

				if aDash, ok := obj.Ctrl.ADash[corp.CooperationId]; ok {
					if aDash.Username != "test@example.com" && aDash.Password != "" {
						if aDash.Login() {
							aDash.GetPapLinks()
							itemCount++
							prog.SetValue(float64(itemCount) / float64(totalItems))
						} else {
							obj.AddLogEntry("Adash Login failed")
						}
					}
				}
			}
		}
		prog.SetValue(1)
		obj.UpdateGui()
		prog.Hide()
	}()
	prog.Show()

}

func (obj *OmipGui) UpdateGui() {
	obj.CorpUpdate()
	obj.TabPtr.Items = obj.TabPtr.Items[:0]
	obj.TabPtr.Items = append(obj.TabPtr.Items, container.NewTabItemWithIcon("Notifications", newspostIcon, obj.notifyScreen()))
	obj.TabPtr.Items = append(obj.TabPtr.Items, container.NewTabItemWithIcon("My Characters", charactersheetIcon, obj.characterScreen()))
	obj.TabPtr.Items = append(obj.TabPtr.Items, container.NewTabItemWithIcon("Corporations", corporationIcon, obj.CorpTabPtr))
	obj.TabPtr.Items = append(obj.TabPtr.Items, container.NewTabItemWithIcon("ESI Keys", dogtagsIcon, obj.keysScreen()))
	obj.CharTabPtr.Refresh()
	if len(obj.CharTabPtr.Items) > 0 {
		obj.CharTabPtr.Select(obj.CharTabPtr.Items[0])
	}
	obj.NotifyEntry.SetText(obj.NotifyText)
	obj.TabPtr.SelectIndex(0)
	obj.TabPtr.Refresh()
}

func (obj *OmipGui) getIconResource(imgFile string) fyne.Resource {
	iconFile, err := os.Open(imgFile)
	if err != nil {
		log.Printf(err.Error())
	}

	r := bufio.NewReader(iconFile)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf(err.Error())
	}

	return fyne.NewStaticResource("icon", b)
}

func (obj *OmipGui) keysScreen() fyne.CanvasObject {
	addButton := widget.NewButtonWithIcon("Add ESI Key", theme.ContentAddIcon(), func() {
		obj.Ctrl.OpenAuthInBrowser()
	})

	keyList := widget.NewList(
		func() int {
			return len(obj.Ctrl.Esi.EsiCharList)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(7),
				widget.NewLabel("name"),
				widget.NewCheck("Contracts", func(bool) {}),
				widget.NewCheck("Industry", func(bool) {}),
				widget.NewCheck("KillMails", func(bool) {}),
				widget.NewCheck("Structures", func(bool) {}),
				widget.NewCheck("Journal", func(bool) {}),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			char := obj.Ctrl.Esi.EsiCharList[id]

			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(char.CharInfoData.CharacterName)

			item.(*fyne.Container).Objects[1].(*widget.Check).SetChecked(obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Contracts)

			item.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(b bool) {
				obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Contracts = b
			}
			item.(*fyne.Container).Objects[2].(*widget.Check).SetChecked(obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.IndustryJobs)
			item.(*fyne.Container).Objects[2].(*widget.Check).OnChanged = func(b bool) {
				obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.IndustryJobs = b
			}
			item.(*fyne.Container).Objects[3].(*widget.Check).SetChecked(obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Killmails)
			item.(*fyne.Container).Objects[3].(*widget.Check).OnChanged = func(b bool) {
				obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Killmails = b
			}
			item.(*fyne.Container).Objects[4].(*widget.Check).SetChecked(obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Structures)
			item.(*fyne.Container).Objects[4].(*widget.Check).OnChanged = func(b bool) {
				obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Structures = b
			}
			item.(*fyne.Container).Objects[5].(*widget.Check).SetChecked(obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Journal)
			item.(*fyne.Container).Objects[5].(*widget.Check).OnChanged = func(b bool) {
				obj.Ctrl.Esi.EsiCharList[id].UpdateFlags.Journal = b
			}

			item.(*fyne.Container).Objects[6].(*widget.Button).OnTapped = func() {
				cnf := dialog.NewConfirm("Confirmation",
					fmt.Sprintf("REALLY delete %s", char.CharInfoData.CharacterName),
					func(confirmed bool) {
						if confirmed {
							obj.Ctrl.Esi.EsiCharList =
								append(obj.Ctrl.Esi.EsiCharList[:id], obj.Ctrl.Esi.EsiCharList[id+1:]...)
							obj.UpdateGui()
						}
					}, obj.WindowPtr)

				cnf.Show()
			}
		})

	return container.NewBorder(addButton, nil, nil, nil, keyList)
}

func (obj *OmipGui) AddEsiKey(char *ctrl.EsiChar) {
	obj.UpdateGui()
	if char.CharInfoExt.Director {
		corp := obj.Ctrl.GetCorp(char)
		if corp != nil {
			obj.AddLogEntry(fmt.Sprintf("added %s director of %s", char.CharInfoData.CharacterName, corp.Name))
		} else {
			obj.AddLogEntry(fmt.Sprintf("added %s", char.CharInfoData.CharacterName))
		}
	} else {
		obj.AddLogEntry(fmt.Sprintf("added %s", char.CharInfoData.CharacterName))
	}
}

func (obj *OmipGui) notifyScreen() fyne.CanvasObject {
	entryMultiLine := widget.NewMultiLineEntry()
	scroll := container.NewVScroll(entryMultiLine)
	obj.NotifyEntry = entryMultiLine
	obj.Progress = widget.NewProgressBar()
	obj.Progress.Hide()
	return container.NewBorder(nil, obj.Progress, nil, nil, scroll)
}

func (obj *OmipGui) AddLogEntry(newEntry string) {
	obj.NotifyText += newEntry + "\n"
	obj.NotifyEntry.SetText(obj.NotifyText)
}

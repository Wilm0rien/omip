package view

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/update"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"io/ioutil"
	"net/url"
	"os"
	"time"
)

const (
	updateUrl = "https://api.github.com/repos/Wilm0rien/omip/releases"
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
	Version     string
}

func NewOmipGui(ctrl *ctrl.Ctrl, app fyne.App, debug bool, version string) *OmipGui {
	var obj OmipGui
	obj.Version = version
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
			fyne.NewMenuItem("Open Log File", func() {
				util.OpenUrl(obj.Ctrl.Model.LocalLogFile)
			}),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("LICENSE", func() {
				d := dialog.NewCustom("OMIP COPYRIGHT NOTICE", "OK", obj.makeLicenseDialog(), obj.WindowPtr)
				d.Show()
			}),

			fyne.NewMenuItem("Check for Update", func() {
				asset, TagName := update.GetRelease(updateUrl, `omip\.zip$`)

				d := dialog.NewCustom("Update Status", "OK", obj.makeupdateDialog(TagName, version, asset), obj.WindowPtr)
				d.Show()

			}),
		))

	obj.TabPtr = container.NewAppTabs(make([]*container.TabItem, 0, 5)...)
	obj.TabPtr.SetTabLocation(container.TabLocationLeading)
	mainContainer := container.NewBorder(nil, obj.makeBottomStatus(), nil, nil, obj.TabPtr)

	obj.WindowPtr = obj.AppPtr.NewWindow(fmt.Sprintf("OMIP - An Eve Online Data Aggregator %s", version))
	obj.WindowPtr.SetMainMenu(menu)
	obj.WindowPtr.SetContent(mainContainer)
	obj.WindowPtr.Resize(fyne.NewSize(1920, 800))
	obj.WindowPtr.SetMaster()
	obj.AppPtr.Settings().SetTheme(theme.DarkTheme())
	obj.AppPtr.Settings().Scale()
	return &obj
}

func (obj *OmipGui) makeBottomStatus() (result fyne.CanvasObject) {
	statusCol1 := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	statusCol2 := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	statusLabel1 := canvas.NewText("", statusCol1)
	statusLabel2 := canvas.NewText("", statusCol2)
	newBottomGrid := container.NewGridWithColumns(2, statusLabel1, statusLabel2)

	lastUpdate := time.Now()
	status1Chan := make(chan string, 1000)
	status2Chan := make(chan string, 1000)
	updateStatusCB := func(entry string, fieldId int) {
		switch fieldId {
		case 1:
			status1Chan <- entry
		case 2:
			status2Chan <- entry
		}
	}

	go func() {
		for {
			elapsed := time.Since(lastUpdate)
			if elapsed.Milliseconds() > 50 {

				select {
				case <-time.After(50 * time.Millisecond):
					if statusCol1.A > 40 {
						statusCol1.A -= 40
					} else {
						statusCol1.A = 0
					}
					if statusCol2.A > 40 {
						statusCol2.A -= 40
					} else {
						statusCol2.A = 0
					}
					statusLabel1.Color = statusCol1
					statusLabel1.Refresh()
					statusLabel2.Color = statusCol2
					statusLabel2.Refresh()
				case status1 := <-status1Chan:
					statusCol1 = color.NRGBA{0xff, 0xff, 0xff, 0xff}
					statusLabel1.Text = status1
					statusLabel1.Color = statusCol1
					statusLabel1.Refresh()
				case status2 := <-status2Chan:
					statusCol2 = color.NRGBA{0xff, 0xff, 0xff, 0xff}
					statusLabel2.Text = status2
					statusLabel2.Color = statusCol2
					statusLabel2.Refresh()
				}
				lastUpdate = time.Now()
			}

			if obj.Ctrl.ServerCancelled() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	obj.Ctrl.GuiStatusCB = updateStatusCB
	return newBottomGrid
}
func (obj *OmipGui) makeupdateDialog(newVersion string, currentVersion string, asset *update.GitAssets) (result fyne.CanvasObject) {
	var msg string
	if newVersion == currentVersion {
		msg = fmt.Sprintf("software is up to date %s current %s", newVersion, currentVersion)
	} else {
		msg = fmt.Sprintf("downloaded %s", asset.Url)
		util.OpenUrl(asset.Url)
	}
	msgLabel := widget.NewLabel(msg)
	return container.NewVBox(msgLabel)
}

/*
	func (obj *OmipGui) makeupdateDialog(newVersion string, currentVersion string, asset *update.GitAssets) (result fyne.CanvasObject) {
		var msg string
		updaterExe := path.Join(obj.Ctrl.Model.LocalDir, "omip_updater.exe")
		updaterZip := path.Join(obj.Ctrl.Model.LocalDir, "omip_updater.zip")
		newOmipZip := path.Join(obj.Ctrl.Model.LocalDir, "omip.zip")
		if util.Exists(updaterExe) {
			obj.RemoveOldUpdater(updaterExe, newVersion)
		}
		if !util.Exists(updaterExe) {
			obj.DownLoadUpdater(updaterExe, updaterZip, asset)
		}

		if !util.Exists(updaterExe) {
			msg = fmt.Sprintf("ERROR updater executable not found at %s", updaterExe)
		} else {
			ex, _ := os.Executable()
			if newVersion == currentVersion {
				msg = fmt.Sprintf("software is up to date %s current %s", newVersion, currentVersion)
			} else {
				switch runtime.GOOS {
				case "linux":
					msg = fmt.Sprintf("TODO LINUX Update not implemented")
				case "windows":
					if dlErr := obj.DownloadUpdate(newOmipZip); dlErr != nil {
						msg = fmt.Sprintf("failed to download omip.zip %s", dlErr.Error())
					} else {
						arguments := fmt.Sprintf(`/k %s --target=%s --source=%s`, updaterExe, ex, newOmipZip)
						cmd := exec.Command("cmd", arguments)
						execErr2 := cmd.Start()
						if execErr2 != nil {
							msg = fmt.Sprintf("error starting process %s", updaterExe)
							obj.Ctrl.Model.LogObj.Printf(msg)
						} else {
							obj.WindowPtr.Close()
						}
					}
				}
			}
		}

		msgLabel := widget.NewLabel(msg)
		return container.NewVBox(msgLabel)
	}
*/
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
	if ok, err := obj.Ctrl.CheckUpdatePreCon(); !ok {
		dialog.ShowError(err, obj.WindowPtr)
		return
	}
	obj.TabPtr.SelectIndex(0)
	obj.Ctrl.NotifyInfo = make(map[int64]bool)
	prog := obj.Progress
	finishCb := func() {
		prog.SetValue(1)
		obj.UpdateGui()
		prog.Hide()
	}
	updateProg := func(c float64) {
		prog.SetValue(c)
	}
	go obj.Ctrl.UpdateAllDataCmd(updateProg, finishCb)
	prog.Show()
}

/*
func (obj *OmipGui) RemoveOldUpdater(updaterExe string, newVersion string) {
	switch runtime.GOOS {
	case "linux":
		obj.Ctrl.AddLogEntry("ERROR LINUX Update not implemented")
	case "windows":
		arguments := fmt.Sprintf(`--version`)
		cmd := exec.Command(updaterExe, arguments)
		output, execErr2 := cmd.Output()
		if execErr2 != nil {
			obj.Ctrl.Model.LogObj.Printf("error starting process %s", updaterExe)
		} else {
			versionStr := string(output)
			if versionStr != newVersion {
				obj.Ctrl.Model.LogObj.Printf("detected old version %s of updater, removing %s", versionStr, updaterExe)
				err := os.Remove(updaterExe)
				if err != nil {
					obj.Ctrl.AddLogEntry(fmt.Sprintf("%s", err.Error()))
				}
			}
		}
	}
}

func (obj *OmipGui) DownLoadUpdater(updaterExe string, updaterZip string, asset *update.GitAssets) {
	if asset == nil {
		obj.Ctrl.AddLogEntry("ERROR updater not found on github")
	} else {
		updateObj := update.NewUpdaterObj()
		downLoadErr := updateObj.DownloadFile(updaterZip, asset.Url, asset.FileSize)
		if downLoadErr != nil {
			obj.Ctrl.Model.LogObj.Printf("updated failed while downloading %s", downLoadErr.Error())
		} else {
			obj.Ctrl.Model.LogObj.Printf("downloaded %s, extracting to %s", updaterZip, updaterExe)
			extractErr := updateObj.ExtractExec(updaterZip, obj.Ctrl.Model.LocalDir)
			if extractErr != nil {
				obj.Ctrl.Model.LogObj.Printf("update failed while extracting %s", extractErr.Error())
			} else {
				if err := os.Remove(updaterZip); err != nil {
					obj.Ctrl.Model.LogObj.Printf(fmt.Sprintf("failed erasing zip file %s %s", updaterZip, err.Error()))
				}
				obj.Ctrl.Model.LogObj.Printf(fmt.Sprintf("downloaded %s", updaterExe))
			}
		}
	}
}

func (obj *OmipGui) DownloadUpdate(newOmipZip string) (err error) {
	updateObj := update.NewUpdaterObj()
	asset, _ := update.GetRelease(updateUrl, `omip\.zip$`)
	return updateObj.DownloadFile(newOmipZip, asset.Url, asset.FileSize)
}*/

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

	// TODO remove workaround for https://github.com/fyne-io/fyne/issues/3169
	obj.TabPtr.OnSelected = func(item *container.TabItem) {
		obj.recurseRefresh(item)
	}
}

func (obj *OmipGui) getIconResource(imgFile string) fyne.Resource {
	iconFile, err := os.Open(imgFile)
	if err != nil {
		obj.Ctrl.Model.LogObj.Printf(err.Error())
	}

	r := bufio.NewReader(iconFile)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		obj.Ctrl.Model.LogObj.Printf(err.Error())
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
	obj.TabPtr.SelectIndex(0)
	obj.AddLogEntry("Update Initiated! Please Wait!")
	obj.Ctrl.UpdateChar(char)

	if char.CharInfoExt.Director {
		corp := obj.Ctrl.GetCorp(char)
		if corp != nil {
			obj.AddLogEntry(fmt.Sprintf("added %s director of %s", char.CharInfoData.CharacterName, corp.Name))
			obj.Ctrl.UpdateCorp(char)
		} else {
			obj.AddLogEntry(fmt.Sprintf("added %s", char.CharInfoData.CharacterName))
		}
	} else {
		obj.AddLogEntry(fmt.Sprintf("added %s", char.CharInfoData.CharacterName))
	}
	obj.UpdateGui()
	obj.TabPtr.SelectIndex(0)
	obj.AddLogEntry("Update Finished!")
}

func (obj *OmipGui) notifyScreen() fyne.CanvasObject {
	entryMultiLine := widget.NewMultiLineEntry()
	scroll := container.NewVScroll(entryMultiLine)
	obj.NotifyEntry = entryMultiLine
	obj.Progress = widget.NewProgressBar()
	obj.Progress.Hide()
	asset, TagName := update.GetRelease(updateUrl, `omip\.zip$`)
	if asset == nil {
		obj.Ctrl.AddLogEntry(fmt.Sprintf("ERROR could not read update url %s", updateUrl))
	} else {
		if TagName != obj.Version {
			obj.Ctrl.AddLogEntry(fmt.Sprintf("NEW VERSION available %s. Update via Menu Bar Help --> Check for Update", TagName))
		}
	}
	return container.NewBorder(nil, obj.Progress, nil, nil, scroll)
}

func (obj *OmipGui) AddLogEntry(newEntry string) {
	obj.NotifyText += newEntry + "\n"
	obj.NotifyEntry.SetText(obj.NotifyText)
}

func (obj *OmipGui) recurseRefresh(item *container.TabItem) {
	//obj.Ctrl.Model.LogObj.Printf("recurse %s", item.Text)
	if _, ok := item.Content.(*fyne.Container); ok {
		obj.recurseRefreshTable(item)
	} else if cont3, ok2 := item.Content.(*container.DocTabs); ok2 {
		for _, subTabItem := range cont3.Items {
			//obj.Ctrl.Model.LogObj.Printf("DocTabs subtab %s", subTabItem.Text)
			obj.recurseRefreshTable(subTabItem)
		}
	} else if cont5, ok4 := item.Content.(*container.AppTabs); ok4 {
		for _, subTabItem := range cont5.Items {
			//obj.Ctrl.Model.LogObj.Printf("AppTabs subtab %s", subTabItem.Text)
			obj.recurseRefreshTable(subTabItem)
		}
	} else if tab, ok := item.Content.(*widget.Table); ok {
		//obj.Ctrl.Model.LogObj.Printf("refreshing table 1")
		tab.Hide()
		tab.Show()
	}
}

func (obj *OmipGui) recurseRefreshTable(item *container.TabItem) {
	if cont, ok := item.Content.(*fyne.Container); ok {
		obj.recurseRefreshContainer(cont)
	} else if _, ok2 := item.Content.(*container.DocTabs); ok2 {
		obj.recurseRefresh(item)
	} else if _, ok4 := item.Content.(*container.AppTabs); ok4 {
		obj.recurseRefresh(item)
	} else if tab, ok := item.Content.(*widget.Table); ok {
		//obj.Ctrl.Model.LogObj.Printf("refreshing table 2")
		tab.Hide()
		tab.Show()
	}
}

func (obj *OmipGui) recurseRefreshContainer(cont *fyne.Container) {
	for _, object := range cont.Objects {
		if cont2, ok := object.(*fyne.Container); ok {
			obj.recurseRefreshContainer(cont2)
		} else if cont3, ok2 := object.(*container.DocTabs); ok2 {
			for _, subTabItem := range cont3.Items {
				if cont4, ok3 := subTabItem.Content.(*fyne.Container); ok3 {
					obj.recurseRefreshContainer(cont4)
				}
			}
		} else if cont5, ok4 := object.(*container.AppTabs); ok4 {
			for _, subTabItem := range cont5.Items {
				if cont6, ok5 := subTabItem.Content.(*fyne.Container); ok5 {
					obj.recurseRefreshContainer(cont6)
				}
			}
		}
		if tab, ok := object.(*widget.Table); ok {
			//obj.Ctrl.Model.LogObj.Printf("refreshing table 3")
			tab.Hide()
			tab.Show()
		}
	}
}

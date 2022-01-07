package view

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
)

func (obj *OmipGui) aDashButton(corpID int) *widget.Button {
	email := "test@example.com"
	pw := ""
	if obj.Ctrl.Model.ADashAuthExists(corpID) {
		decryptOK := false
		email, pw, decryptOK = obj.Ctrl.Model.GetAuth(corpID)
		if !decryptOK {
			err := errors.New("could not retrieve password from database: decryption error")
			dialog.ShowError(err, obj.WindowPtr)
		}
	}
	ticker := obj.Ctrl.Model.GetCorpTicker(corpID)
	obj.Ctrl.ADash[corpID] = ctrl.NewADashClient(email, pw, ticker, obj.Ctrl.Model, corpID)
	obj.Ctrl.ADash[corpID].AddLogCB = obj.Ctrl.AddLogEntry
	aDash := widget.NewButton("Login to aDashBoard to get pap statistics", func() {
		aDashDialog, email, pw := obj.makeADashDialog(corpID)
		dialog.ShowCustomConfirm("aDashboard", "OK", "CANCEL", aDashDialog,
			func(response bool) {
				if response {
					obj.Ctrl.ADash[corpID].Username = *email
					obj.Ctrl.ADash[corpID].Password = *pw
					if obj.Ctrl.ADash[corpID].Login() {
						obj.Ctrl.Model.SetAuth(corpID, *email, *pw)
					} else {
						fmt.Printf("adash fail %s %s\n", *email, *pw)
					}
				}
			}, obj.WindowPtr)
	})
	return aDash
}

func (obj *OmipGui) makeADashDialog(corpID int) (co fyne.CanvasObject, email *string, pw *string) {
	emailLabel := widget.NewLabel("Email:")
	emailEntry := widget.NewEntry()
	emailEntry.Text = obj.Ctrl.ADash[corpID].Username

	pw1Label := widget.NewLabel("Password:                                        ")
	pw1Entry := widget.NewPasswordEntry()
	pw1Entry.Text = obj.Ctrl.ADash[corpID].Password
	testResResult := widget.NewLabel("not Tested")

	testButton := widget.NewButton("Test Login Data", func() {
		obj.Ctrl.ADash[corpID].Username = emailEntry.Text
		obj.Ctrl.ADash[corpID].Password = pw1Entry.Text
		if obj.Ctrl.ADash[corpID].Login() {
			testResult:= obj.Ctrl.ADash[corpID].CheckPapLinks()
			testResResult.SetText(testResult)
		} else {
			testResResult.SetText("FAIL")
		}
	})

	retval := container.New(layout.NewGridLayout(2))
	retval.Objects = append(retval.Objects, emailLabel)
	retval.Objects = append(retval.Objects, emailEntry)

	retval.Objects = append(retval.Objects, pw1Label)
	retval.Objects = append(retval.Objects, pw1Entry)

	retval.Objects = append(retval.Objects, testButton)
	retval.Objects = append(retval.Objects, testResResult)
	return retval, &emailEntry.Text, &pw1Entry.Text
}

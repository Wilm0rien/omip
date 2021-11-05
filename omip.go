package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"github.com/Wilm0rien/omip/view"
	"time"
)

var testEnableFlag = flag.Bool("test", false, "enable tests")
var debugEnableFlag = flag.Bool("debug", false, "enable debug tab")

func main() {
	flag.Parse()
	// kill existing instance
	urlStrShutdown := fmt.Sprintf("http://localhost:4716/callback?code=shutdown&state=0")
	if util.SendReq(urlStrShutdown) {
		time.Sleep(400 * time.Millisecond)
	}
	modelObj := model.NewModel(model.DbName, *testEnableFlag)
	ctrlObj := ctrl.NewCtrl(modelObj)
	ctrlObj.StartServer()
	loadErr := ctrlObj.Load(ctrl.ConfigFileName, *testEnableFlag)
	app := app.New()
	gui := view.NewOmipGui(ctrlObj, app, *debugEnableFlag)
	gui.WindowPtr.Show()
	gui.UpdateGui()
	if loadErr != nil {
		gui.AddLogEntry(loadErr.Error())
	}
	gui.WindowPtr.ShowAndRun()
	ctrlObj.Save(ctrl.ConfigFileName, *testEnableFlag)
	if !ctrlObj.ServerCancelled() {
		ctrlObj.HTTPShutdown()
	}
	modelObj.CloseDB()
}

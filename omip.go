package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"github.com/Wilm0rien/omip/view"
	"log"
	"time"
)

var testEnableFlag = flag.Bool("test", false, "enable tests")
var debugEnableFlag = flag.Bool("debug", false, "enable debug tab")

var guiEnableFlag = flag.Bool("gui", false, "enable gui (cmd mode)")
var cmdEnableFlag = flag.Bool("cmd", false, "enable cmd")
var CmdLineOpt string

func main() {
	flag.Parse()
	if *guiEnableFlag {
		CmdLineOpt = "default_gui"
	}
	// kill existing instance
	urlStrShutdown := fmt.Sprintf("http://localhost:4716/callback?code=shutdown&state=0")
	if util.SendReq(urlStrShutdown) {
		time.Sleep(400 * time.Millisecond)
	}
	modelObj := model.NewModel(model.DbName, *testEnableFlag)
	ctrlObj := ctrl.NewCtrl(modelObj)
	ctrlObj.StartServer()
	loadErr := ctrlObj.Load(ctrl.ConfigFileName, *testEnableFlag)
	if CmdLineOpt == "default_cmd" {
		ctrlObj.AddLogCB = func(entry string) {
			fmt.Printf("%s\n", entry)
		}
		ctrlObj.UpdateAllDataCmd(nil, nil)
	} else {
		appObj := app.New()
		gui := view.NewOmipGui(ctrlObj, appObj, *debugEnableFlag, util.OmipSoftwareVersion)
		gui.WindowPtr.Show()
		gui.UpdateGui()
		if loadErr != nil {
			gui.AddLogEntry(loadErr.Error())
		}
		gui.WindowPtr.ShowAndRun()
	}
	if !ctrlObj.ServerCancelled() {
		ctrlObj.HTTPShutdown()
	}
	safeErr := ctrlObj.Save(ctrl.ConfigFileName, *testEnableFlag)
	if safeErr != nil {
		log.Printf("%s", safeErr.Error())
	}
	modelObj.CloseDB()

}

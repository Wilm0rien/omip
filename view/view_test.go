package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"log"
	"testing"
	"time"
)

func TestGui(t *testing.T) {
	modelObj := model.NewModel(model.DbNameCtrlTest, false, false)
	ctrlObj := ctrl.NewCtrl(modelObj)
	ctrlObj.Load(ctrl.TstCfgJson, true)
	app := test.NewApp()
	gui := NewOmipGui(ctrlObj, app, true, "test version")
	start := time.Now()
	gui.UpdateGui()

	//test.TapCanvas(gui.WindowPtr.Canvas(), gui.TabPtr.Position())
	//test.TapCanvas(gui.WindowPtr.Canvas(), gui.CorpTabPtr.Position())
	log.Printf("%v", gui.TabPtr.Items[1])
	test.TapCanvas(gui.WindowPtr.Canvas(), fyne.Position{60, 170})
	test.TapCanvas(gui.WindowPtr.Canvas(), fyne.Position{314, 165})

	elapsed := time.Since(start)
	gui.TabPtr.SelectIndex(1)
	if gui.TabPtr.Selected().Text != "My Characters" {
		t.Fatalf("expected char tab but got %s", gui.TabPtr.Selected().Text)
	}

	log.Printf("UpdateGui took %s", elapsed)
	modelObj.CloseDB()
}

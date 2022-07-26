package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"log"
	"testing"
	"time"
)

func TestGui(t *testing.T) {
	modelObj := model.NewModel(model.DbNameCtrlTest, false)
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
	if tab, ok := gui.TabPtr.Selected().Content.(*container.DocTabs); ok {
		if len(tab.Items) != 1 {
			t.Fatalf("expected one tab")
		} else {
			if border, ok2 := tab.Selected().Content.(*fyne.Container); ok2 {
				var foundTable *widget.Table
				for _, obj := range border.Objects {
					if table, ok3 := obj.(*widget.Table); ok3 {
						foundTable = table
					}
				}
				if foundTable == nil {
					t.Fatalf("expected table ")
				}
				foundTable.Select(widget.TableCellID{Col: 0, Row: 0})
				tab.Refresh()

			} else {
				t.Fatalf("expected one border")
			}
		}
	} else {
		t.Error("expected tab in my characters")
	}

	log.Printf("UpdateGui took %s", elapsed)
	modelObj.CloseDB()
}

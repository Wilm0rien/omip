package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
)

func (obj *OmipGui) createMiningTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	fullList := obj.Ctrl.Model.GetMiningData(char.CharInfoExt.CooperationId)


	// todo check GetMonthlyTable() to calcate the m3 per month
	// todo add volume m3 to ViewMiningData and calcualte the m3 when populating the table
	/*

		var leftListWidget *widget.List
		var rightListWidget *widget.List


		leftListWidget = widget.NewList(
			func() int {
				return len(filteredAggList)
			},
			func() fyne.CanvasObject {
				return container.New(layout.NewHBoxLayout(),
					widget.NewLabel("Character"),
					canvas.NewText("Amount m3", color.NRGBA{0, 0x80, 0, 0xff}),
					widget.NewLabel("Date"))
			},
			func(id widget.ListItemID, item fyne.CanvasObject) {
				date := util.UnixTS2DateStr(filteredAggList[id].Date)
				refType := getJourTypeRef(filteredAggList[id])
				amount := jourGetAmount(filteredAggList[id])
				item.(*fyne.Container).Objects[0].(*widget.Label).SetText(date)
				item.(*fyne.Container).Objects[1].(*canvas.Text).Text = amount
				item.(*fyne.Container).Objects[1].(*canvas.Text).Color = getColor(filteredAggList[id].Amount)
				item.(*fyne.Container).Objects[2].(*widget.Label).SetText(refType)
			})
	*/
	mainGrid := widget.NewMultiLineEntry()
	topgrid := widget.NewButton("testbutton", func() {

		mainGrid.Text += fmt.Sprintf("%d\n", len(fullList))
		mainGrid.Refresh()
	})
	bottomgrid := widget.NewButton("clear", func() {
		mainGrid.Text = ""
		mainGrid.Refresh()
	})
	return container.NewBorder(topgrid, bottomgrid, nil, nil, mainGrid), true
}

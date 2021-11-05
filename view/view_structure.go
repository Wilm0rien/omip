package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"regexp"
)

func (obj *OmipGui) createStructureTab(director *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	structureList := obj.Ctrl.Model.GetCorpStructures(director.CharInfoExt.CooperationId)
	if len(structureList) > 0 {
		result = true
	} else {
		return retTable, result
	}

	var structureListWidget *widget.List
	var serviceListWidget *widget.List
	var lastSelectedStructure int64
	var filterStructures *widget.Entry
	var filterTypes *widget.Entry
	var filterServices *widget.Entry
	selectedLabel := widget.NewLabel("Services of selected Structure")

	filteredStructureList := obj.Ctrl.Model.GetCorpStructures(director.CharInfoExt.CooperationId)
	strucNameMapping := make(map[int64]string)
	svcNameMapping := make(map[int64]string)
	svcStateMapping := make(map[int64]string)
	svcMapping := make(map[int64][]*model.DBstructureService)

	for _, structure := range structureList {
		name := obj.Ctrl.Model.GetStructureNameStr(structure.StructureID)
		strucNameMapping[structure.StructureID] = name
		strSvcs := obj.Ctrl.Model.GetServiceEntries(structure.StructureID)
		svcMapping[structure.StructureID] = strSvcs
		for _, svc := range strSvcs {
			svcNameMapping[svc.Name], _ = obj.Ctrl.Model.GetStringEntry(svc.Name)
			svcStateMapping[svc.State], _ = obj.Ctrl.Model.GetStringEntry(svc.State)
		}
	}
	updateLists := func() {
		filteredStructureList = filteredStructureList[:0]
		for _, structure := range structureList {
			if structureName, ok := strucNameMapping[structure.StructureID]; ok {
				structype := obj.Ctrl.Model.GetTypeString(structure.TypeID)
				fStrMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterStructures.Text), structureName)
				fTypMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterTypes.Text), structype)

				fSvcMatch := false
				if filterServices.Text == "" {
					fSvcMatch = true
				} else {
					if v, ok := svcMapping[structure.StructureID]; ok {
						for _, svc := range v {
							svcName := svcNameMapping[svc.Name]
							match, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterServices.Text), svcName)
							if match {
								fSvcMatch = true
								break
							}
						}
					}
				}

				if fStrMatch && fTypMatch && fSvcMatch {
					filteredStructureList = append(filteredStructureList, structure)
					lastSelectedStructure = structure.StructureID
				}
			}
		}
		structureListWidget.Refresh()
		serviceListWidget.Refresh()
	}

	structureListWidget = widget.NewList(
		func() int {
			return len(filteredStructureList)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(4),
				widget.NewLabel("structure name"),
				widget.NewLabel("structure type"),
				widget.NewLabel("fuel expires"),
				widget.NewLabel("state"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			stateStr := "unkown"
			for k, v := range model.StructureStatus {
				if v == filteredStructureList[id].State {
					stateStr = k
				}
			}
			structureName := strucNameMapping[filteredStructureList[id].StructureID]
			fuelExpire := ""
			expireStr, future := util.GetTimeDiffStringFromTS(filteredStructureList[id].FuelExpires)
			if future {
				fuelExpire = fmt.Sprintf("%s", expireStr)
			} else {
				fuelExpire = fmt.Sprintf("%s (overdue)", expireStr)
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(structureName)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(obj.Ctrl.Model.GetTypeString(filteredStructureList[id].TypeID))
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(fuelExpire)
			item.(*fyne.Container).Objects[3].(*widget.Label).SetText(stateStr)
		})
	structureListWidget.OnSelected = func(id widget.ListItemID) {
		lastSelectedStructure = filteredStructureList[id].StructureID
		selectedLabel.SetText(fmt.Sprintf("Services of %s", strucNameMapping[lastSelectedStructure]))
		serviceListWidget.Refresh()
	}
	serviceListWidget = widget.NewList(
		func() int {
			var retval int
			if v, ok := svcMapping[lastSelectedStructure]; ok {
				retval = len(v)
			}
			return retval
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(2),
				widget.NewLabel("service name"),
				widget.NewLabel("service state"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if v, ok := svcMapping[lastSelectedStructure]; ok {
				if svcName, ok := svcNameMapping[v[id].Name]; ok {
					item.(*fyne.Container).Objects[0].(*widget.Label).SetText(svcName)
				}
				if svcState, ok := svcStateMapping[v[id].State]; ok {
					item.(*fyne.Container).Objects[1].(*widget.Label).SetText(svcState)
				}
			}
		})
	filterStructures = widget.NewEntry()
	filterStructures.PlaceHolder = "filter Structure Names"
	filterStructures.OnChanged = func(s string) {
		updateLists()
	}
	filterTypes = widget.NewEntry()
	filterTypes.PlaceHolder = "filter Structure Types"
	filterTypes.OnChanged = func(s string) {
		updateLists()
	}
	filterServices = widget.NewEntry()
	filterServices.PlaceHolder = "filter Services"
	filterServices.OnChanged = func(s string) {
		updateLists()
	}
	resetFilterBtn := widget.NewButton("Reset Filters", func() {
		filterStructures.SetText("")
		filterTypes.SetText("")
		filterServices.SetText("")
	})
	mainGrid := container.New(layout.NewGridLayout(2), structureListWidget, serviceListWidget)
	topGrid := container.New(layout.NewGridLayout(2),
		container.New(layout.NewGridLayout(4),
			widget.NewLabel("Name"),
			widget.NewLabel("Type"),
			widget.NewLabel("Fuel Expires"),
			widget.NewLabel("State")),
		selectedLabel)
	bottomGrid := container.New(layout.NewGridLayout(4), filterStructures, filterTypes, filterServices, resetFilterBtn)
	updateLists()
	selectedLabel.SetText(fmt.Sprintf("Services of %s", strucNameMapping[lastSelectedStructure]))
	return container.NewBorder(topGrid, bottomGrid, nil, nil, mainGrid), result
}

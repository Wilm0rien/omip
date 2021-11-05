package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"regexp"
	"sort"
	"time"
)

type induGuiTable struct {
	fulllist     []*model.DBJob
	filteredList []*model.DBJob
	genericTable
}

const (
	InduT0Status = iota
	InduT1Runs
	InduT2Activity
	InduT3Blueprint
	InduT4Facility
	InduT5Installer
	InduT6InstallDate
	InduT7EndDate
)

func NewInduTable(ctrl *ctrl.Ctrl, fullList []*model.DBJob) *induGuiTable {
	var obj induGuiTable
	obj.Ctrl = ctrl
	obj.fulllist = fullList
	obj.filteredList = make([]*model.DBJob, 0, 10)
	obj.header = make([]string, 0, 10)
	obj.header = append(obj.header, "Status")
	obj.header = append(obj.header, "Runs")
	obj.header = append(obj.header, "Activity")
	obj.header = append(obj.header, "Blueprint")
	obj.header = append(obj.header, "Facility")
	obj.header = append(obj.header, "Installer")
	obj.header = append(obj.header, "Install Date")
	obj.header = append(obj.header, "End Date")
	obj.colWidth = make([]float32, 0, 10)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 50)
	obj.colWidth = append(obj.colWidth, 100)
	obj.colWidth = append(obj.colWidth, 300)
	obj.colWidth = append(obj.colWidth, 300)
	obj.colWidth = append(obj.colWidth, 300)
	obj.colWidth = append(obj.colWidth, 150)
	obj.colWidth = append(obj.colWidth, 150)
	obj.filter = make([]string, 0, 10)
	for i := 0; i < len(obj.header); i++ {
		obj.filter = append(obj.filter, "")
	}
	return &obj
}
func (obj *induGuiTable) GetNumRows() int {
	return len(obj.filteredList)
}
func (obj *induGuiTable) getCellStrFromList(rowIdx int, colIdx int, inputList []*model.DBJob) string {
	var retval string
	if rowIdx < len(inputList) {
		listElem := inputList[rowIdx]
		if colIdx < len(obj.colWidth) {
			switch colIdx {
			case InduT0Status:
				if listElem.Status == model.Job_Stat_active {
					retval2, future := util.GetTimeDiffStringFromTS(listElem.EndDate)
					if future {
						duration := listElem.EndDate - listElem.StartDate
						partDuration := listElem.EndDate - time.Now().Unix()
						percentile := (1 - (float64(partDuration) / float64(duration))) * 100
						retval = fmt.Sprintf("%2.0f%% %s", percentile, retval2)
					} else {
						retval = fmt.Sprintf("-%s", retval2)
					}
				} else {
					retval = fmt.Sprintf("%s", obj.Ctrl.Model.JobStatusId2Str(listElem.Status))
				}
			case InduT1Runs:
				retval = fmt.Sprintf("%d", listElem.Runs)
			case InduT2Activity:
				retval = fmt.Sprintf("%s", obj.Ctrl.Model.JobActivityId2Str(listElem.ActivityId))
			case InduT3Blueprint:
				retval = fmt.Sprintf("%s", obj.Ctrl.Model.GetTypeString(listElem.ProductTypeId))
			case InduT4Facility:
				structureName := fmt.Sprintf("%d", listElem.FacilityId)
				nameStruct := obj.Ctrl.Model.GetStructureName(listElem.FacilityId)
				if nameStruct != nil {
					structureName, _ = obj.Ctrl.Model.GetStringEntry(nameStruct.NameRef)
				}
				retval = structureName
			case InduT5Installer:
				installerName := fmt.Sprintf("%d", listElem.InstallerId)
				installer := obj.Ctrl.Model.GetCharEntry(listElem.InstallerId)
				if installer.Name != "" {
					installerName = installer.Name
				} else {
					corpMember := obj.Ctrl.Model.GetDBCorpMember(listElem.InstallerId)
					if corpMember != nil {
						installerName, _ = obj.Ctrl.Model.GetStringEntry(corpMember.NameRef)
					}
				}
				retval = installerName
			case InduT6InstallDate:
				retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(listElem.StartDate))
			case InduT7EndDate:
				retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(listElem.EndDate))
			}
		}
	}
	return retval
}
func (obj *induGuiTable) GetCellStr(rowIdx int, colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	if rowIdx < len(obj.filteredList) {
		listElem := obj.filteredList[rowIdx]
		if colIdx == InduT0Status {
			if listElem.Status == model.Job_Stat_ready {
				col = color.NRGBA{0, 0xff, 0, 0xff}
			} else if listElem.Status == model.Job_Stat_active {
				partDuration := listElem.EndDate - time.Now().Unix()
				if partDuration > 0 {
					duration := listElem.EndDate - listElem.StartDate
					ratio := 1 - (float64(partDuration) / float64(duration))
					col = *util.GetColor(1, ratio, false)
				}
			}
		}
	}
	return obj.getCellStrFromList(rowIdx, colIdx, obj.filteredList), col
}
func (obj *induGuiTable) GetSumCellStr(colIdx int) (string, color.NRGBA) {
	col := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	txt := ""
	switch colIdx {
	case InduT0Status:
		txt = fmt.Sprintf("Cnt: %d", len(obj.filteredList))
	}
	return txt, col
}
func (obj *induGuiTable) GetCSVCellStr(rowIdx int, colIdx int) string {
	return obj.getCellStrFromList(rowIdx, colIdx, obj.filteredList)
}
func (obj *induGuiTable) SelectedFunc() func(id widget.TableCellID) {
	return func(id widget.TableCellID) {

	}
}

func (obj *induGuiTable) UpdateLists() {
	obj.filteredList = obj.filteredList[:0]
	for rowIdx, _ := range obj.fulllist {
		filterOK := true
		for colIdx, _ := range obj.header {
			currentFilter := obj.filter[colIdx]
			currentCellStr := obj.getCellStrFromList(rowIdx, colIdx, obj.fulllist)
			fMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", currentFilter), currentCellStr)
			if !fMatch {
				filterOK = false
				break
			}
		}
		if filterOK {
			obj.filteredList = append(obj.filteredList, obj.fulllist[rowIdx])
		}
	}
}
func (obj *induGuiTable) SortCol(colIdx int) {
	obj.sortCount++
	sort.Slice(obj.fulllist, func(i, j int) bool {
		a := obj.getCellStrFromList(i, colIdx, obj.fulllist)
		b := obj.getCellStrFromList(j, colIdx, obj.fulllist)
		var retval bool
		retval = a >= b
		if obj.sortCount%2 == 0 {
			retval = !retval
		}
		return retval
	})
}

func (obj *OmipGui) createIndustryTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	var fullList []*model.DBJob
	if corp {
		fullList = obj.Ctrl.Model.GetIndustryJobs(char.CharInfoExt.CooperationId, corp)
	} else {
		fullList = obj.Ctrl.Model.GetIndustryJobs(char.CharInfoData.CharacterID, corp)
	}

	if len(fullList) == 0 {
		result = false
		return retTable, result
	} else {
		table := NewInduTable(obj.Ctrl, fullList)
		table.UpdateLists()
		retTable = obj.createGenTable2(table, false, true, "")
		result = true
	}

	return retTable, result
}

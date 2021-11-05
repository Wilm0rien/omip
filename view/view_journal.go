package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wilm0rien/omip/ctrl"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"image/color"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"
)

type jourAggr struct {
	mainNode *model.DBJournal
	nodeList []*model.DBJournal
}

const (
	PERIOD_LAST_7_DAYS  = "last 7 days"
	PERIOD_LAST_30_DAYS = "last 30 days"
	PERIOD_LAST_60_DAYS = "last 60 days"
	PERIOD_LAST_90_DAYS = "last 90 days"
	PERIOD_THIS_MONTH   = "this month"
	PERIOD_LAST_MONTH   = "last month"
	PERIOD_THIS_YEAR    = "this year"
	FILTER_SIGN_PLUS    = "income"
	FILTER_SIGN_MINUS   = "expenses"
	FILTER_SIGN_ALL     = "all transactions"
)

func (obj *OmipGui) createJournalTab(char *ctrl.EsiChar, corp bool) (retTable fyne.CanvasObject, result bool) {
	sumAggregate := canvas.NewText("0", color.NRGBA{0, 0x80, 0, 0xff})
	signFilterStr := FILTER_SIGN_ALL
	var filterSign *widget.Select
	var sumAggrateValue float64
	fullList := obj.Ctrl.Model.GetJournal(char.CharInfoData.CharacterID, char.CharInfoExt.CooperationId, corp)
	if len(fullList) == 0 {
		return retTable, result
	}
	var filterRefTypes *widget.Entry
	var filterDate *widget.Entry
	var filterPeriod *widget.Select
	var filterAmount *widget.Entry

	contractDescmap := make(map[int64]string)
	industryJobDescmap := make(map[int64]string)
	aggregateRefmap := make(map[string]map[int]jourAggr)

	dateMap := make(map[int64]string)
	jourGetDate := func(jour *model.DBJournal) string {
		var retval string
		if v, ok := dateMap[jour.Date]; ok {
			retval = v
		} else {
			retval = fmt.Sprintf("%s", util.ConvertUnixTimeToStr(jour.Date))
			dateMap[jour.Date] = retval
		}
		return retval
	}
	descMap := make(map[int64]string)
	jourGetDesc := func(jour *model.DBJournal) string {
		var retval string
		if v, ok := descMap[jour.Description]; ok {
			retval = v
		} else {
			retval, _ = obj.Ctrl.Model.GetStringEntry(jour.Description)
			descMap[jour.Description] = retval
		}
		return retval
	}

	reasonMap := make(map[int64]string)
	jourGetReason := func(jour *model.DBJournal) string {
		var retval string
		if v, ok := reasonMap[jour.Reason]; ok {
			retval = v
		} else {
			retval, _ = obj.Ctrl.Model.GetStringEntry(jour.Reason)
			reasonMap[jour.Reason] = retval
		}
		return retval
	}
	jourGetAmount := func(jour *model.DBJournal) string {
		return fmt.Sprintf("%10.10s", util.HumanizeNumber(jour.Amount))
	}

	itemTypeMap := make(map[int64]string)
	jourGetItem := func(jour *model.DBJournal) string {
		var retval string
		if jour.Ref_type == 2 || jour.Ref_type == 42 {
			if v, ok := itemTypeMap[jour.ID]; ok {
				retval = v
			} else {
				transaction := obj.Ctrl.Model.GetTransactionEntry(jour.ID)
				if transaction != nil {
					retval = obj.Ctrl.Model.GetTypeString(transaction.TypeID)
				} else {
					retval = "N/A"
				}
				itemTypeMap[jour.ID] = retval
			}
		} else {
			retval = "N/A"
		}
		return retval
	}
	refIdTypeMap := make(map[int]string)
	getJourTypeRef := func(jour *model.DBJournal) string {
		var retval string
		if v, ok := refIdTypeMap[jour.Ref_type]; ok {
			retval = v
		} else {
			for value, k := range model.JournalRefType {
				retval = "N/A"
				if k == jour.Ref_type {
					retval = value
					refIdTypeMap[jour.Ref_type] = retval
				}
			}
		}
		return retval
	}

	transactionMap := make(map[int64]*model.DBTransaction)

	getTransactionQuant := func(jour *model.DBJournal) string {
		var retval string
		if jour.Ref_type == 2 || jour.Ref_type == 42 {
			if v, ok := transactionMap[jour.ID]; ok {
				retval = fmt.Sprintf("%d", v.Quantity)
			} else {
				transaction := obj.Ctrl.Model.GetTransactionEntry(jour.ID)
				if transaction == nil {
					var empty model.DBTransaction
					transactionMap[jour.ID] = &empty

				} else {
					transactionMap[jour.ID] = transaction
				}
				retval = fmt.Sprintf("%d", transactionMap[jour.ID].Quantity)
			}
		}
		return retval
	}

	getTransactionPrice := func(jour *model.DBJournal) string {
		var retval string
		if jour.Ref_type == 2 || jour.Ref_type == 42 {
			if v, ok := transactionMap[jour.ID]; ok {
				retval = util.HumanizeNumber(v.UnitPrice)
			} else {
				transaction := obj.Ctrl.Model.GetTransactionEntry(jour.ID)
				if transaction == nil {
					var empty model.DBTransaction
					transactionMap[jour.ID] = &empty

				} else {
					transactionMap[jour.ID] = transaction
				}
				retval = util.HumanizeNumber(v.UnitPrice)
			}
		}
		return retval
	}

	for _, row := range fullList {
		dateStr := util.UnixTS2DateStr(row.Date)
		addNode := false
		if _, ok := aggregateRefmap[dateStr]; ok {
			if aggregate, ok2 := aggregateRefmap[dateStr][row.Ref_type]; ok2 {
				aggregate.mainNode.Amount += row.Amount
				aggregate.nodeList = append(aggregate.nodeList, row)
				aggregateRefmap[dateStr][row.Ref_type] = aggregate
			} else {
				addNode = true
			}
		} else {
			aggregateRefmap[dateStr] = make(map[int]jourAggr)
			addNode = true
		}
		if addNode {
			var newJourAggregate jourAggr
			var newMainNode model.DBJournal
			newMainNode = *row
			newJourAggregate.mainNode = &newMainNode
			newJourAggregate.nodeList = make([]*model.DBJournal, 0, 10)
			newJourAggregate.nodeList = append(newJourAggregate.nodeList, row)
			aggregateRefmap[dateStr][row.Ref_type] = newJourAggregate
		}
	}
	fullAggrList := make([]*model.DBJournal, 0, 10)
	filteredAggList := make([]*model.DBJournal, 0, 10)
	keylist := util.GetSortKeysFromStrMap(aggregateRefmap, true)
	for _, dateStr := range keylist {
		keylist2 := util.GetSortKeysFromIntMap(aggregateRefmap[dateStr], true)
		for _, typeRef := range keylist2 {
			fullAggrList = append(fullAggrList, aggregateRefmap[dateStr][typeRef].mainNode)
			filteredAggList = append(filteredAggList, aggregateRefmap[dateStr][typeRef].mainNode)

		}
	}
	var lastSelectedDate string
	var lastSelectedRefType int
	var aggregateListWidget *widget.List
	var typeListWidget *widget.List

	var maxAbsPositive float64
	var maxAbsNegative float64
	updateMax := func(input float64, updateval *float64) {
		abs := math.Abs(input)
		if abs > *updateval {
			*updateval = abs
		}
	}
	getColor := func(input float64) *color.NRGBA {
		ratio := float64(1)
		returnColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
		if input > 0 {
			ratio = math.Abs(input) / maxAbsPositive
			green := 64*ratio + 191

			returnColor.G = uint8(green)
			returnColor.B = uint8(ratio * 127)
			returnColor.R = uint8(ratio * 127)
			//log.Printf("input %3.2f green: %d", input, uint8(green))

		} else {
			ratio = math.Abs(input) / maxAbsNegative
			red := 64*ratio + 191
			returnColor.R = uint8(red)
			returnColor.B = uint8(ratio * 127)
			returnColor.G = uint8(ratio * 127)
			//log.Printf("input %3.2f ratio %3.2f red: %v", input, ratio, returnColor)
		}
		return &returnColor
	}
	updateLists := func() {
		maxAbsNegative = 0
		maxAbsPositive = 0
		filteredAggList = filteredAggList[:0]
		sumAggrateValue = 0
		for _, listItem := range fullAggrList {
			date := util.UnixTS2DateStr(listItem.Date)
			refType := getJourTypeRef(listItem)
			amountMillion := math.Abs(listItem.Amount / 1000000)
			periodOK := checkPeriod(filterPeriod.Selected, listItem.Date)
			fDateMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterDate.Text), date)
			// date input overrides period
			if filterDate.Text != "" {
				periodOK = true
			}
			fTypMatch, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", filterRefTypes.Text), refType)
			var minimum float64
			if s, err := strconv.ParseFloat(filterAmount.Text, 64); err == nil {
				minimum = s
			}
			signCheck := true
			if signFilterStr == FILTER_SIGN_PLUS {
				signCheck = false
				if listItem.Amount >= 0 {
					signCheck = true
				}
			} else if signFilterStr == FILTER_SIGN_MINUS {
				signCheck = false
				if listItem.Amount < 0 {
					signCheck = true
				}
			}

			if periodOK && fDateMatch && fTypMatch && amountMillion >= minimum && signCheck {
				filteredAggList = append(filteredAggList, listItem)
				lastSelectedDate = util.UnixTS2DateStr(listItem.Date)
				lastSelectedRefType = listItem.Ref_type
				if listItem.Amount > 0 {
					updateMax(listItem.Amount, &maxAbsPositive)
				} else {
					updateMax(listItem.Amount, &maxAbsNegative)
				}
				sumAggrateValue += listItem.Amount
			}

		}
		sumAggregate.Text = util.HumanizeNumber(sumAggrateValue)
		if sumAggrateValue > 0 {
			sumAggregate.Color = color.NRGBA{0, 0xFF, 0, 0xff}
		} else if sumAggrateValue == 0 {
			sumAggregate.Color = color.NRGBA{0xff, 0xFF, 0xff, 0xff}
		} else {
			sumAggregate.Color = color.NRGBA{0xff, 0, 0, 0xff}
		}
		sumAggregate.Refresh()
		aggregateListWidget.Refresh()
		typeListWidget.Refresh()
	}

	selectedLabel := widget.NewLabel("selected journal entry")
	aggregateListWidget = widget.NewList(
		func() int {
			return len(filteredAggList)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewHBoxLayout(),
				widget.NewLabel("Date"),
				canvas.NewText("Amount", color.NRGBA{0, 0x80, 0, 0xff}),
				widget.NewLabel("Type"))
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

	aggregateListWidget.OnSelected = func(id widget.ListItemID) {
		lastSelectedDate = util.UnixTS2DateStr(filteredAggList[id].Date)
		lastSelectedRefType = filteredAggList[id].Ref_type
		refTypeStr := getJourTypeRef(filteredAggList[id])
		selectedLabel.SetText(fmt.Sprintf("%s %s", lastSelectedDate, refTypeStr))
		typeListWidget.Refresh()
	}
	parseForJobIDs := func(desc string) (info string) {
		jobRe := regexp.MustCompile(`[(]Job ID: ([0-9]+)[)]`)

		jobResult := jobRe.FindStringSubmatch(desc)
		if jobResult != nil {
			jobID, err := strconv.ParseInt(jobResult[1], 10, 64)
			if err == nil {
				if val, ok := industryJobDescmap[jobID]; ok {
					info = val
				} else {
					industryJobDescmap[jobID] = "N/A"
					job := obj.Ctrl.Model.GetIndustryJob(jobID)
					if job != nil {
						product := fmt.Sprintf("%s", obj.Ctrl.Model.GetTypeString(job.ProductTypeId))
						structureName := fmt.Sprintf("%d", job.FacilityId)
						nameStruct := obj.Ctrl.Model.GetStructureName(job.FacilityId)
						if nameStruct != nil {
							structureName, _ = obj.Ctrl.Model.GetStringEntry(nameStruct.NameRef)
						}

						activity := fmt.Sprintf("%s", obj.Ctrl.Model.JobActivityId2Str(job.ActivityId))
						installerName := fmt.Sprintf("%d", job.InstallerId)
						installer := obj.Ctrl.Model.GetCharEntry(job.InstallerId)
						if installer.Name != "" {
							installerName = installer.Name
						} else {
							corpMember := obj.Ctrl.Model.GetDBCorpMember(job.InstallerId)
							if corpMember != nil {
								installerName, _ = obj.Ctrl.Model.GetStringEntry(corpMember.NameRef)
							}
						}

						info = fmt.Sprintf("%s; %s; %s; %s", installerName, activity, product, structureName)
						industryJobDescmap[jobID] = info
					}
				}
			}
		}
		return
	}
	parseForContractIDs := func(desc string) (info string) {
		ctrRe := regexp.MustCompile(`^(.*)[(]contract ID: ([0-9]+)[)]`)
		ctrResult := ctrRe.FindStringSubmatch(desc)
		if ctrResult != nil {
			ctrID, err := strconv.ParseInt(ctrResult[2], 10, 64)
			if err == nil {
				if val, ok := contractDescmap[ctrID]; ok {
					info = val
				} else {
					contractDescmap[ctrID] = desc
					ctr := obj.Ctrl.Model.GetContractById(ctrID)
					if ctr != nil {
						collectedInfo := ""
						ctrDesc, ok2 := obj.Ctrl.Model.GetStringEntry(ctr.Title)
						if ok2 {
							collectedInfo = ctrDesc
						}

						items := obj.Ctrl.Model.GetContrItems(ctr.Contract_id)
						if len(items) == 1 {
							collectedInfo += " " + obj.Ctrl.Model.GetTypeString(items[0].Type_id)
						} else if len(items) > 1 {
							collectedInfo += " " + obj.GetMaxPricedItem(items)
						}
						contractDescmap[ctrID] = ctrResult[1] + ": " + collectedInfo
					}
				}
			}
		}
		return
	}
	typeListWidget = widget.NewList(
		func() int {
			var retval int
			if _, ok1 := aggregateRefmap[lastSelectedDate]; ok1 {
				if v2, ok2 := aggregateRefmap[lastSelectedDate][lastSelectedRefType]; ok2 {
					retval = len(v2.nodeList)
				}
			}
			return retval
		},
		func() fyne.CanvasObject {
			newDescriptionLabel := widget.NewLabel("Description")
			newDescriptionLabel.Wrapping = fyne.TextWrapWord
			labelDate := widget.NewLabel("Date")
			labelAmount := widget.NewLabel("Amount")

			dateAmount := container.New(layout.NewHBoxLayout(), labelDate, labelAmount)
			split := container.NewHSplit(dateAmount, newDescriptionLabel)
			split.Offset = 0.2
			return split
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if _, ok1 := aggregateRefmap[lastSelectedDate]; ok1 {
				if v2, ok2 := aggregateRefmap[lastSelectedDate][lastSelectedRefType]; ok2 {
					if id < len(v2.nodeList) {
						currentNode := v2.nodeList[id]
						date := jourGetDate(currentNode)
						amount := jourGetAmount(currentNode)
						desc := jourGetDesc(currentNode)
						if currentNode.Ref_type == model.Journal_market_transaction ||
							currentNode.Ref_type == model.Journal_market_escrow {
							itemStr := jourGetItem(currentNode)
							quantitiy := getTransactionQuant(currentNode)
							unitPrice := getTransactionPrice(currentNode)
							markeStr := fmt.Sprintf("%s x %s for %s each", quantitiy, itemStr, unitPrice)
							if itemStr != "N/A" {
								desc = markeStr
							}
						}

						reason := jourGetReason(currentNode)
						if reason != "" {
							desc += reason
						}
						info := parseForJobIDs(desc)
						if info != "" {
							desc = info
						} else {
							info2 := parseForContractIDs(desc)
							if info2 != "" {
								desc = info2
							}
						}

						//desc += fmt.Sprintf(" ref type = %d", currentNode.Ref_type)

						item.(*container.Split).Leading.(*fyne.Container).Objects[0].(*widget.Label).SetText(date)
						item.(*container.Split).Leading.(*fyne.Container).Objects[1].(*widget.Label).SetText(amount)
						item.(*container.Split).Trailing.(*widget.Label).SetText(desc)
					}
				}
			}
		})

	filterRefTypes = widget.NewEntry()
	filterRefTypes.PlaceHolder = "filter ref types"
	filterRefTypes.OnChanged = func(s string) {
		updateLists()
	}
	filterDate = widget.NewEntry()
	filterDate.PlaceHolder = "filter date"
	filterDate.OnChanged = func(s string) {
		updateLists()
	}
	filterPeriod = widget.NewSelect(
		[]string{PERIOD_THIS_MONTH, PERIOD_LAST_MONTH, PERIOD_LAST_7_DAYS, PERIOD_LAST_30_DAYS, PERIOD_LAST_60_DAYS,
			PERIOD_LAST_90_DAYS, PERIOD_THIS_YEAR}, func(s string) {
			updateLists()
			obj.Ctrl.Esi.NVConfig.PeriodFilter = s
		})

	filterAmount = widget.NewEntry()
	filterAmount.PlaceHolder = "filter amount millions"
	filterAmount.OnChanged = func(s string) {
		updateLists()
	}

	resetFilterBtn := widget.NewButton("Reset Filters", func() {
		filterRefTypes.SetText("")
		filterDate.SetText("")
		filterSign.SetSelected(FILTER_SIGN_ALL)
	})
	filterSign = widget.NewSelect(
		[]string{FILTER_SIGN_ALL, FILTER_SIGN_PLUS, FILTER_SIGN_MINUS}, func(s string) {
			signFilterStr = s
			updateLists()
		})
	sumLabel := widget.NewLabel("Sum:")
	//mainGrid := container.New(layout.NewGridLayout(2), aggregateListWidget, typeListWidget)
	mainGrid := container.NewHSplit(aggregateListWidget, typeListWidget)
	mainGrid.Offset = 0.3
	//mainGrid := container.NewBorder(nil, nil, aggregateListWidget, nil, typeListWidget)
	topGrid := container.New(layout.NewGridLayout(2),
		widget.NewLabel("Combined transaction types per day"),
		selectedLabel)
	filterSign.SetSelected(FILTER_SIGN_ALL)
	bottomGrid := container.New(layout.NewGridLayout(8),
		filterPeriod, filterDate, filterRefTypes, filterAmount, filterSign, resetFilterBtn, sumLabel, sumAggregate)
	updateLists()
	if obj.Ctrl.Esi.NVConfig.PeriodFilter != "" {
		filterPeriod.SetSelected(obj.Ctrl.Esi.NVConfig.PeriodFilter)
	} else {
		filterPeriod.SetSelected(PERIOD_THIS_MONTH)
	}
	result = true
	return container.NewBorder(topGrid, bottomGrid, nil, nil, mainGrid), result
}

func checkPeriod(periodType string, inputTimeStamp int64) (result bool) {
	tm := time.Unix(inputTimeStamp, 0)
	days := -7
	timeStart := time.Now()

	if periodType != PERIOD_LAST_MONTH {
		switch periodType {
		case PERIOD_LAST_7_DAYS:
			days = -7
		case PERIOD_LAST_30_DAYS:
			days = -30
		case PERIOD_LAST_60_DAYS:
			days = -60
		case PERIOD_LAST_90_DAYS:
			days = -90
		case PERIOD_THIS_MONTH:
			nowYear, nowMonth, _ := time.Now().Date()
			timeStart = getMonthStartTime(nowYear, nowMonth)
			days = 0
		case PERIOD_THIS_YEAR:
			nowYear, _, _ := time.Now().Date()
			timeStart = getMonthStartTime(nowYear, 1)
			days = 0
		}

		if tm.After(timeStart.AddDate(0, 0, days)) {
			result = true
		}
	} else {
		nowYear, nowMonth, _ := time.Now().Date()
		startYear := nowYear
		startMonth := nowMonth
		if startMonth == 1 {
			startMonth = 12
			startYear--
		} else {
			startMonth--
		}

		timeStart = getMonthStartTime(startYear, startMonth)
		timeEnd := getMonthStartTime(nowYear, nowMonth)
		afterOK := false
		beforeOK := false
		if tm.After(timeStart) {
			afterOK = true
		}
		if tm.Before(timeEnd) {
			beforeOK = true
		}
		if afterOK && beforeOK {
			result = true
		}

	}

	return result
}

func getMonthStartTime(year int, month time.Month) (t time.Time) {
	timeString := fmt.Sprintf("%04d-%02d-01T00:00:05Z", year, month)
	t, err := time.Parse("2006-01-02T15:04:05Z", timeString)
	if err != nil {
		fmt.Println(err)
		log.Printf("ConvertTimeStrToInt ERROR PARSING TIME %s", timeString)
	}
	return
}

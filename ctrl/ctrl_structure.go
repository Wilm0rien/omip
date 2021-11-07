package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"log"
	"time"
)

type structInfo struct {
	CorpId               int              `json:"corporation_id"`
	FuelExpires          string           `json:"fuel_expires"`
	NextReinforceApply   string           `json:"next_reinforce_apply"`
	NextReinforceHour    int              `json:"next_reinforce_hour"`
	NextReinforceWeekday int              `json:"next_reinforce_weekday"`
	ProfileId            int              `json:"profile_id"`
	ReinForceHour        int              `json:"reinforce_hour"`
	ReinforceWeekday     int              `json:"reinforce_weekday"`
	Services             []structServices `json:"services"`
	State                string           `json:"state"`
	StateTimerEnd        string           `json:"state_timer_end"`
	StateTimerStart      string           `json:"state_timer_start"`
	StructureID          int64            `json:"structure_id"`
	SystemID             int              `json:"system_id"`
	TypeID               int              `json:"type_id"`
	UnanchorsAt          string           `json:"unanchors_at"`
}
type structServices struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type structPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type structName struct {
	Name          string         `json:"name"`
	OwnerID       int            `json:"owner_id"`
	Position      structPosition `json:"position"`
	SolarSystemID int            `json:"solar_system_id"`
	TypeID        int            `json:"type_id"`
}

func (obj *Ctrl) UpdateStructures(char *EsiChar, corp bool) {
	if !char.UpdateFlags.Structures {
		return
	}
	if corp {
		url := fmt.Sprintf("https://esi.evetech.net/v4/corporations/%d/structures/?datasource=tranquility", char.CharInfoExt.CooperationId)
		bodyBytes, _ := obj.getSecuredUrl(url, char)
		var structinfo []structInfo
		//log.Printf("%s\n", string(bodyBytes))
		contentError := json.Unmarshal(bodyBytes, &structinfo)
		corpInfo2 := obj.GetCorp(char)
		if corpInfo2 != nil {
			if contentError != nil {
				obj.AddLogEntry(fmt.Sprintf("%s ERROR reading structures", corpInfo2.Name))
			} else {
				obj.processStructInfo(structinfo, corpInfo2, char)
			}
		} else {
			obj.AddLogEntry(fmt.Sprintf("UpdateStructures %s ERROR reading corpinfo", char.CharInfoData.CharacterName))
		}

	}
}

func (obj *Ctrl) GetStructureNameFromEsi(char *EsiChar, structureId int64) (retval string) {
	url := fmt.Sprintf("https://esi.evetech.net/v2/universe/structures/%d/?datasource=tranquility", structureId)
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var structName2 structName
	contentError := json.Unmarshal(bodyBytes, &structName2)
	corp := obj.GetCorp(char)
	if corp != nil {
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("%s ERROR reading structureName Code %s", corp.Name, contentError.Error()))
			retval = fmt.Sprintf("%d", structureId)
		} else {
			dbStruct := obj.convertEsiStructureName2DB(&structName2, structureId)
			obj.Model.AddStructureNameEntry(dbStruct)
			retval = structName2.Name
		}
	} else {
		obj.AddLogEntry(fmt.Sprintf("GetStructureNameFromEsi %s ERROR reading corp info",char.CharInfoData.CharacterName))
	}

	return retval
}

func (obj *Ctrl) processStructInfo(structinfo []structInfo, corp *EsiCorp, char *EsiChar) {
	// collect the list of structure IDs which actually exist
	structIDsFromEsi := make([]int64, 0, 4)
	structureExistMap := make(map[int64]int)
	for _, elem := range structinfo {
		structureExistMap[elem.StructureID] = 1
		structIDsFromEsi = append(structIDsFromEsi, elem.StructureID)
		dbStruct := obj.convertEsiStructureInfo2DB(&elem)
		oldInfo := obj.Model.GetStructureInfo(elem.StructureID)
		var newStructure bool
		if oldInfo == nil {
			structNameItem := obj.Model.GetStructureName(elem.StructureID)
			if structNameItem == nil {
				if elem.StructureID < 100000000 {
					obj.AddLogEntry(fmt.Sprintf("ERROR processStructInfo invalid station id %d", elem.StructureID))
				} else {
					name := obj.GetStructureNameFromEsi(char, elem.StructureID)
					dbStruct.NameRef = obj.Model.AddStringEntry(name)
				}
			} else {
				dbStruct.NameRef = structNameItem.NameRef
			}
			newStructure = true
		} else {
			dbStruct.NameRef = oldInfo.NameRef
			oldStatus := obj.Model.StrucStatID2Str(oldInfo.State)
			if oldStatus != elem.State {
				name, _ := obj.Model.GetStringEntry(oldInfo.NameRef)
				obj.AddLogEntry(fmt.Sprintf("[%s] %s : %s",
					corp.Ticker, name, elem.State))
			}
		}

		result := obj.Model.AddStructureInfoEntry(dbStruct)
		if result == model.DBR_Undefined {
			log.Printf("unexpetcted return code")
		}
		strucName, _ := obj.Model.GetStringEntry(dbStruct.NameRef)
		obj.processStructServices(elem.Services, elem.StructureID, corp, strucName, newStructure)

		diff := dbStruct.FuelExpires - time.Now().Unix()
		if diff < 24*60*60*7 {
			str, _ := util.GetTimeDiffStringFromTS(dbStruct.FuelExpires)
			obj.AddLogEntry(fmt.Sprintf("[%s] %s fuel expires in %s",
				corp.Ticker, strucName, str))
		}
	}
	corpStructures := obj.Model.GetCorpStructures(char.CharInfoExt.CooperationId)
	ticker:="N/A"
	corpObj:=obj.GetCorp(char)
	if corpObj!=nil {
		ticker = corpObj.Ticker
	}
	for _, corpStruc := range corpStructures {
		if _, ok := structureExistMap[corpStruc.StructureID]; !ok {
			name := obj.Model.GetStructureNameStr(corpStruc.StructureID)
			obj.AddLogEntry(fmt.Sprintf("%s: removing sturcute %s", ticker, name))
			obj.Model.DeleteStructureServiceEntries(corpStruc.StructureID)
			obj.Model.DeleteStructureNameEntries(corpStruc.StructureID)
			obj.Model.DeleteStructureInfoEntries(corpStruc.StructureID)
		}
	}
}

func (obj *Ctrl) processStructServices(services []structServices, structId int64, corp *EsiCorp, nameStr string, newStructure bool) {
	svcIDsFromEsi := make([]int64, 0, 4)

	for _, service := range services {
		var newService model.DBstructureService
		newService.StructureID = structId
		newService.Name = obj.Model.AddStringEntry(service.Name)
		newService.State = obj.Model.AddStringEntry(service.State)
		svcIDsFromEsi = append(svcIDsFromEsi, newService.Name)
		dbResult := obj.Model.AddStructureServiceEntry(&newService)
		if dbResult == model.DBR_Updated || dbResult == model.DBR_Inserted && !newStructure {
			obj.AddLogEntry(fmt.Sprintf("[%s] %s - %s %s", corp.Ticker, nameStr, service.Name, service.State))
		} else if dbResult == model.DBR_Undefined {
			log.Printf("WARNING: processEsiStructServices unexpected result")
		}
	}
	// remove deleted services from DB
	svcIDsFromDB := obj.Model.GetListOfServiceEntries(structId)
	for _, dbService := range svcIDsFromDB {
		serviceFound := false
		for _, existing := range svcIDsFromEsi {
			if dbService == existing {
				serviceFound = true
				break
			}
		}
		if serviceFound == false {
			serviceNameString, _ := obj.Model.GetStringEntry(dbService)
			result := obj.Model.DeleteStructureServiceEntry(structId, dbService)
			if result {
				obj.AddLogEntry(fmt.Sprintf("[%s] %s removed %s", corp.Ticker, nameStr, serviceNameString))
			} else {
				log.Printf("WARNING FAILED TO delete service %d from structure %s\n", structId, serviceNameString)
			}
		}
	}

}

func (obj *Ctrl) convertEsiStructureInfo2DB(esiStructInfo *structInfo) *model.DBstructureInfo {
	var newStructInfo model.DBstructureInfo
	newStructInfo.CorpId = esiStructInfo.CorpId
	newStructInfo.FuelExpires = util.ConvertTimeStrToInt(esiStructInfo.FuelExpires)
	newStructInfo.NextReinforceApply = util.ConvertTimeStrToInt(esiStructInfo.NextReinforceApply)
	newStructInfo.NextReinforceHour = esiStructInfo.NextReinforceHour
	newStructInfo.NextReinforceWeekday = esiStructInfo.NextReinforceWeekday
	newStructInfo.ProfileId = esiStructInfo.ProfileId
	newStructInfo.ReinForceHour = esiStructInfo.ReinForceHour
	newStructInfo.ReinforceWeekday = esiStructInfo.ReinforceWeekday
	newStructInfo.State = model.StructureStatus[esiStructInfo.State]
	newStructInfo.StateTimerEnd = util.ConvertTimeStrToInt(esiStructInfo.StateTimerEnd)
	newStructInfo.StateTimerStart = util.ConvertTimeStrToInt(esiStructInfo.StateTimerStart)
	newStructInfo.StructureID = esiStructInfo.StructureID
	newStructInfo.SystemID = esiStructInfo.SystemID
	newStructInfo.TypeID = esiStructInfo.TypeID
	newStructInfo.UnanchorsAt = util.ConvertTimeStrToInt(esiStructInfo.UnanchorsAt)
	return &newStructInfo
}

func (obj *Ctrl) convertEsiStructureName2DB(esiStructName *structName, structureId int64) *model.DBstructureName {
	var newStructInfo model.DBstructureName
	newStructInfo.StructureID = structureId
	newStructInfo.NameRef = obj.Model.AddStringEntry(esiStructName.Name)
	newStructInfo.OwnerID = esiStructName.OwnerID
	newStructInfo.PositionX = esiStructName.Position.X
	newStructInfo.PositionY = esiStructName.Position.Y
	newStructInfo.PositionZ = esiStructName.Position.Z
	newStructInfo.SolarSystemID = esiStructName.SolarSystemID
	newStructInfo.TypeID = esiStructName.TypeID
	return &newStructInfo
}

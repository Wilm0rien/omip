package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
)

type allyInfo struct {
	CreatorCorpID  int    `json:"creator_corporation_id"`
	CreatorID      int    `json:"creator_id"`
	DateFounded    string `json:"date_founded"`
	ExecutorCorpID int    `json:"executor_corporation_id"`
	FactionID      int    `json:"faction_id"`
	Name           string `json:"name"`
	Ticker         string `json:"ticker"`
}

func (obj *Ctrl) GetAllyInfoFromEsi(char *EsiChar, allyID int) (dbAllyInfo *model.DBallyInfo, result bool) {
	url := fmt.Sprintf("https://esi.evetech.net/v4/alliances/%d?datasource=tranquility", allyID)
	bodyBytes, _, _ := obj.getSecuredUrl(url, char)
	if bodyBytes != nil {
		var ainfo allyInfo
		err := json.Unmarshal(bodyBytes, &ainfo)
		if err == nil {
			dbAllyInfo = obj.convertEsiAllyInfo2DB(&ainfo, allyID)
			dbResult := obj.Model.AddAllyInfoEntry(dbAllyInfo)
			if dbResult == model.DBR_Inserted {
				result = true
			}
		}
	}
	return
}

func (obj *Ctrl) convertEsiAllyInfo2DB(esiAllyInfo *allyInfo, allyID int) *model.DBallyInfo {
	var newAllyInfo model.DBallyInfo
	newAllyInfo.AllianceID = allyID
	newAllyInfo.CreatorCorporationID = esiAllyInfo.CreatorCorpID
	newAllyInfo.CreatorID = esiAllyInfo.CreatorID
	newAllyInfo.DateFounded = util.ConvertTimeStrToInt(esiAllyInfo.DateFounded)
	newAllyInfo.ExecutorCorpID = esiAllyInfo.ExecutorCorpID
	newAllyInfo.FactionID = esiAllyInfo.FactionID
	newAllyInfo.NameStrRef = obj.Model.AddStringEntry(esiAllyInfo.Name)
	newAllyInfo.TickerStrRef = obj.Model.AddStringEntry(esiAllyInfo.Ticker)
	return &newAllyInfo
}

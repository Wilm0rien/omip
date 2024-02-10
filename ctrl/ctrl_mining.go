package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"time"
)

type MiningObservers struct {
	LastUpdated  string `json:"last_updated"`
	ObserverID   int64  `json:"observer_id"`
	ObserverType string `json:"observer_type"`
}

type MiningData struct {
	CharacterID    int32  `json:"character_id"`
	LastUpdated    string `json:"last_updated"`
	Quantity       int32  `json:"quantity"`
	RecordedCorpId int32  `json:"recorded_corporation_id"`
	TypeId         int32  `json:"type_id"`
}

// UpdateCorpMiningObs retrieve list of corp mining observers via /corporation/{corporation_id}/mining/observers/
func (obj *Ctrl) UpdateCorpMiningObs(char *EsiChar, _UnusedCorp bool) {
	// needs esi-industry.read_corporation_mining.v1
	pageID := 1
	for {
		url := fmt.Sprintf("https://esi.evetech.net/v1/corporation/%d/mining/observers?datasource=tranquility&page=%d",
			char.CharInfoExt.CooperationId, pageID)
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var miningObsList []*MiningObservers
		contentError := json.Unmarshal(bodyBytes, &miningObsList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			break
		}
		for _, miningObserver := range miningObsList {
			newObs := obj.convertEsiMOBS2DB(miningObserver)
			if lastTime, found := obj.Model.GetLastObsUpdateTime(miningObserver.ObserverID); found {
				if lastTime != newObs.LastUpdated {
					obj.getMiningData(char, miningObserver.ObserverID)
					db1R := obj.Model.AddMiningObsEntry(newObs)
					util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)
				}
			}
		}
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}
}

// getMiningData retrieve Moon mining data via /corporation/{corporation_id}/mining/observers/{observer_id}/
func (obj *Ctrl) getMiningData(char *EsiChar, observerID int64) {
	// needs esi-industry.read_corporation_mining.v1
	/*
		var url string
		pageID := 1
		for {

			url = fmt.Sprintf("https://esi.evetech.net/v1/corporation/%d/mining/observers/%d/?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, division, pageID)

		}

	*/
}

func (obj *Ctrl) convertEsiMOBS2DB(mObs *MiningObservers) *model.DBMiningObserver {
	var newMObs model.DBMiningObserver
	newMObs.LastUpdated = util.ConvertTimeStrToInt(mObs.LastUpdated)
	newMObs.ObserverID = mObs.ObserverID
	newMObs.ObserverType = obj.Model.AddStringEntry(mObs.ObserverType)
	return &newMObs
}

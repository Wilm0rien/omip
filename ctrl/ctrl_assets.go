package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"time"
)

type Asset struct {
	IsBlueprintCopy bool   `json:"is_blueprint_copy"`
	IsSingleton     bool   `json:"is_singleton"`
	ItemId          int64  `json:"item_id"`
	LocationFlag    string `json:"location_flag"`
	LocationId      int64  `json:"location_id"`
	LocationType    string `json:"location_type"`
	Quantity        int32  `json:"quantity"`
	TypeId          int32  `json:"type_id"`
}

func (obj *Ctrl) UpdateAssets(char *EsiChar, corp bool) {
	var url string
	timestamp := time.Now().Unix()
	pageID := 1
	addedNum := 0
	for {
		if corp {
			url = fmt.Sprintf("https://esi.evetech.net/v5/corporations/%d/assets/?datasource=tranquility&page=%d",
				char.CharInfoExt.CooperationId, pageID)
		} else {
			url = fmt.Sprintf("https://esi.evetech.net/v5/characters/%d/assets/?datasource=tranquility&page=%d",
				char.CharInfoData.CharacterID, pageID)
		}
		bodyBytes, Xpages := obj.getSecuredUrl(url, char)
		var assetList []Asset
		contentError := json.Unmarshal(bodyBytes, &assetList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			break
		}
		for _, asset := range assetList {
			dbAsset := obj.convertEsiAsset2DB(timestamp, &asset, char.CharInfoExt.CooperationId, char.CharInfoData.CharacterID)
			if dbResult := obj.Model.AddAssetEntry(dbAsset); dbResult == model.DBR_Inserted {
				addedNum++
			}
		}
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}
	if addedNum > 0 {
		namePrefix := fmt.Sprintf("[%s] %s: ", obj.GetCorpTicker(char), char.CharInfoData.CharacterName)
		obj.AddLogEntry(namePrefix + fmt.Sprintf("%d assets added", addedNum))
	}

}

func (obj *Ctrl) convertEsiAsset2DB(timestamp int64, assetEntry *Asset, corpId int,
	charId int) *model.DBAsset {
	var newAsset model.DBAsset
	newAsset.CharID = charId
	newAsset.CorpID = corpId
	if assetEntry.IsBlueprintCopy {
		newAsset.IsBlueprintCopy = 1
	}
	if assetEntry.IsSingleton {
		newAsset.IsBlueprintCopy = 1
	}
	newAsset.ItemId = assetEntry.ItemId
	newAsset.LocationFlag = obj.Model.AddStringEntry(assetEntry.LocationFlag)
	newAsset.LocationId = assetEntry.LocationId
	newAsset.LocationType = obj.Model.AddStringEntry(assetEntry.LocationFlag)
	newAsset.Quantity = int(assetEntry.Quantity)
	newAsset.TypeId = int(assetEntry.TypeId)
	newAsset.Timestamp = timestamp
	return &newAsset
}

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
	pageID := 1
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
		// todo
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}

}

func (obj *Ctrl) convertEsiAsset2DB(jourEntry *Journal, corpId int, charId int) *model.DBJournal {
	return nil
}

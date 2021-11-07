package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"net/http"
)

type corpInfo struct {
	AllianceId  int     `json:"alliance_id"`
	CeoId       int     `json:"ceo_id"`
	CreatorId   int     `json:"creator_id"`
	DateFounded string  `json:"date_founded"`
	Description string  `json:"description"`
	FactionId   int     `json:"faction_id"`
	HomeStation int     `json:"home_station"`
	MemberCount int     `json:"member_count"`
	Name        string  `json:"name"`
	Shares      int64   `json:"shares"`
	TaxRate     float64 `json:"tax_rate"`
	Ticker      string  `json:"ticker"`
	Url         string  `json:"url"`
	WarEligible bool    `json:"war_eligible"`
}

type universeNames struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

func (obj *Ctrl) getAuthCorpInfo(char *EsiChar) *EsiCorp {
	dbcorp, result := obj.Model.GetCorpInfoEntry(char.CharInfoExt.CooperationId)
	var esicorp EsiCorp

	if result == model.DBR_Success {
		esicorp.Name, _ = obj.Model.GetStringEntry(dbcorp.CorpNameStrRef)
		esicorp.CooperationId = char.CharInfoExt.CooperationId
		esicorp.AllianceId = dbcorp.AllianceId
		esicorp.Ticker, _ = obj.Model.GetStringEntry(dbcorp.TickerStrRef)
	} else {
		url := fmt.Sprintf("https://esi.evetech.net/v5/corporations/%d?datasource=tranquility", char.CharInfoExt.CooperationId)
		bodyBytes, _ := obj.getSecuredUrl(url, char)
		if bodyBytes != nil {
			var cinfo corpInfo
			err := json.Unmarshal(bodyBytes, &cinfo)
			if err != nil {
				corp := obj.GetCorp(char)
				corpName:="N/A"
				if corp != nil {
					corpName = corp.Name
				}
				obj.AddLogEntry(fmt.Sprintf("%s ERROR reading corporations", corpName))
			} else {
				dbcorp = obj.convertEsiCorpInfo2DB(&cinfo, char.CharInfoExt.CooperationId)
				dbResult := obj.Model.AddCorpInfoEntry(dbcorp)
				if dbResult == model.DBR_Inserted {

				}
				esicorp.Name, _ = obj.Model.GetStringEntry(dbcorp.CorpNameStrRef)
				esicorp.CooperationId = char.CharInfoExt.CooperationId
				esicorp.AllianceId = dbcorp.AllianceId
				esicorp.Ticker, _ = obj.Model.GetStringEntry(dbcorp.TickerStrRef)
			}
		}
	}
	util.Assert(esicorp.CooperationId != 0)
	imgFile := fmt.Sprintf("%d_128.jpg", char.CharInfoExt.CooperationId)
	imgPath := fmt.Sprintf("%s/%s", obj.Model.LocalImgDir, imgFile)
	if !util.Exists(imgPath) {
		fullUrlFile := "https://imageserver.eveonline.com/Corporation/" + imgFile
		util.GetImgFromUrl(fullUrlFile, imgPath)
	}
	esicorp.ImageFile = imgPath
	return &esicorp
}

func (obj *Ctrl) convertEsiCorpInfo2DB(esiCorpInfo *corpInfo, corpId int) *model.DBcorpInfo {
	var newCorpInfo model.DBcorpInfo
	newCorpInfo.CorpID = corpId
	newCorpInfo.AllianceId = esiCorpInfo.AllianceId
	newCorpInfo.CeoId = esiCorpInfo.CeoId
	newCorpInfo.CreatorId = esiCorpInfo.CreatorId
	newCorpInfo.DateFounded = util.ConvertTimeStrToInt(esiCorpInfo.DateFounded)
	newCorpInfo.DescriptionStrRef = 0
	newCorpInfo.FactionId = esiCorpInfo.FactionId
	newCorpInfo.HomeStationId = esiCorpInfo.HomeStation
	newCorpInfo.MemberCount = esiCorpInfo.MemberCount
	newCorpInfo.CorpNameStrRef = obj.Model.AddStringEntry(esiCorpInfo.Name)
	newCorpInfo.Shares = esiCorpInfo.Shares
	newCorpInfo.TaxRate = esiCorpInfo.TaxRate
	newCorpInfo.TickerStrRef = obj.Model.AddStringEntry(esiCorpInfo.Ticker)
	newCorpInfo.UrlStrRef = 0
	newCorpInfo.WarEligible = esiCorpInfo.WarEligible
	return &newCorpInfo
}

func (obj *Ctrl) UpdateCorpMembers(director *EsiChar, corp bool) {
	if corp {
		// esi-corporations.read_corporation_membership.v1
		url := fmt.Sprintf("https://esi.evetech.net/v4/corporations/%d/members/?datasource=tranquility", director.CharInfoExt.CooperationId)
		bodyBytes, _ := obj.getSecuredUrl(url, director)
		if bodyBytes != nil {

			obj.getMemberNames(string(bodyBytes), director)
		}
	}
}

func (obj *Ctrl) getMemberNames(memberIDs string, director *EsiChar) {
	url := fmt.Sprintf("https://esi.evetech.net/v3/universe/names/")
	bodyBytes2, resp := obj.getSecuredUrlPost(url, memberIDs, director)
	if resp.StatusCode == http.StatusOK {
		var members []universeNames
		err := json.Unmarshal(bodyBytes2, &members)
		if err != nil {
			obj.AddLogEntry(fmt.Sprintf(err.Error()))
		} else {
			obj.processMembers(members, director)
		}
	}
}

func (obj *Ctrl) processMembers(members []universeNames, director *EsiChar) {
	memberMap := make(map[int]int)
	updateCounter := 0
	insertCounter := 0
	for _, member := range members {
		var newMember model.DBcorpMember
		memberMap[member.ID] = member.ID
		newMember.CharID = member.ID
		newMember.CorpID = director.CharInfoExt.CooperationId
		newMember.MainID = member.ID
		newMember.NameRef = obj.Model.AddStringEntry(member.Name)

		result := obj.Model.AddCorpMemberEntry(&newMember)
		if result == model.DBR_Inserted {
			insertCounter++
		}
		if result == model.DBR_Updated {
			updateCounter++
		}
	}
	dbMemberList := obj.Model.GetCorpMemberList(director.CharInfoExt.CooperationId)
	for _, dbMemberId := range dbMemberList {
		_, exists := memberMap[dbMemberId]
		if !exists {
			obj.Model.RemoveCorpMemberEntry(dbMemberId, director.CharInfoExt.CooperationId)
		}
	}
	result := obj.Model.UpdateMemberCount(len(members), director.CharInfoExt.CooperationId)
	util.Assert(result == model.DBR_Updated)

}

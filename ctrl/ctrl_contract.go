package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"time"
)

const (
	maxContractItems int = 1000
)

type Contracts struct {
	Acceptor_id           int32   `json:"acceptor_id"`
	Assignee_id           int32   `json:"assignee_id"`
	Availability          string  `json:"availability"`
	Buyout                float64 `json:"buyout"`
	Collateral            float64 `json:"collateral"`
	Contract_id           int32   `json:"contract_id"`
	Date_accepted         string  `json:"date_accepted"`
	Date_completed        string  `json:"date_completed"`
	Date_expired          string  `json:"date_expired"`
	Date_issued           string  `json:"date_issued"`
	Days_to_complete      int32   `json:"days_to_complete"`
	End_location_id       int64   `json:"end_location_id"`
	For_corporation       bool    `json:"for_corporation"`
	Issuer_corporation_id int32   `json:"issuer_corporation_id"`
	Issuer_id             int32   `json:"issuer_id"`
	Price                 float64 `json:"price"`
	Reward                float64 `json:"reward"`
	Start_location_id     int64   `json:"start_location_id"`
	Status                string  `json:"status"`
	Title                 string  `json:"title"`
	Type                  string  `json:"type"`
	Volume                float64 `json:"volume"`
}

type CntrItems struct {
	Is_included  bool  `json:"is_included"`
	Is_singleton bool  `json:"is_singleton"`
	Quantity     int   `json:"quantity"`
	Raw_quantity int   `json:"raw_quantity"`
	Record_id    int64 `json:"record_id"`
	Type_id      int   `json:"type_id"`
}

func (obj *Ctrl) UpdateContracts(char *EsiChar, corp bool) {
	var url string
	pageID := 1
	if !char.UpdateFlags.Contracts {
		return
	}
	for {
		if corp {
			url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/contracts/?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, pageID)
		} else {
			url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/contracts/?datasource=tranquility&page=%d", char.CharInfoData.CharacterID, pageID)
		}
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var contractList []Contracts
		contentError := json.Unmarshal(bodyBytes, &contractList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			break
		}
		obj.UpdateContractsInDb(char, corp, contractList)
		obj.UpdateContractItems(char, corp)
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
		util.Assert(pageID < 30)

	}
}

func (obj *Ctrl) UpdateContractsInDb(char *EsiChar, corp bool, contractList []Contracts) {
	memberIdMap := obj.Model.GetCorpMemberIdMap(char.CharInfoExt.CooperationId)
	statusMap := make(map[string]int64)
	statusMap["requireAttention"] = 0
	statusMap["endMinSeconds"] = invalidTimeDiff
	unixTime := time.Now().Unix()
	ctrExistMap := make(map[int32]int)
	for _, contract := range contractList {
		var IssuerKnown bool
		var jourEntry *model.DBJournal
		if corp {
			_, IssuerKnown = memberIdMap[int(contract.Issuer_id)]
			if !IssuerKnown {
				if contract.Acceptor_id == int32(char.CharInfoExt.CooperationId) {
					IssuerKnown = true
				}
			}
		} else {
			IssuerKnown = char.CharInfoData.CharacterID == int(contract.Issuer_id) ||
				char.CharInfoData.CharacterID == int(contract.Acceptor_id)
			if char.CharInfoData.CharacterID == int(contract.Acceptor_id) {

			}
		}
		if IssuerKnown {
			if corp {
				jourEntry = obj.Model.GetJournalForContract(int(contract.Contract_id))
				// skip all corp contracts which do not have a journal entry
				// these are open contracts available for the corp or accepted by corp members privately
				if !contract.For_corporation && jourEntry == nil {
					continue
				}
			} else {
				if contract.For_corporation {
					continue
				}
			}
			ctrExistMap[contract.Contract_id] = 1
			dbCon := obj.convertEsiContract2DB(&contract)
			result := obj.Model.AddContractEntry(dbCon)
			if result == model.DBR_Updated {
				statusMap[contract.Status]++
			} else if result == model.DBR_Undefined {
				obj.AddLogEntry("UpdateContractsInDb ERROR adding contract entry")
			}
			if contract.Status != "finished" && contract.Status != "deleted" {
				diff := dbCon.Date_expired - unixTime
				if diff <= 0 {
					statusMap["requireAttention"] = 1
				} else {
					if diff < statusMap["endMinSeconds"] {
						statusMap["endMinSeconds"] = diff
					}
				}
			}
		}
	}
	var dbList []*model.DBContract
	if corp {
		dbList = obj.Model.GetContractsByIssuerId(char.CharInfoExt.CooperationId, true)
	} else {
		dbList = obj.Model.GetContractsByIssuerId(char.CharInfoData.CharacterID, false)
	}
	for _, ctr := range dbList {
		if _, ok := ctrExistMap[int32(ctr.Contract_id)]; !ok {
			if ctr.Status == model.Cntr_Stat_outstanding {

				ctr.Status = model.Cntr_Stat_deleted
				obj.Model.AddContractEntry(ctr)
			}
		}
	}

	obj.SummaryLogEntry("contract", char, corp, statusMap)
}

func (obj *Ctrl) convertEsiContract2DB(contr *Contracts) *model.DBContract {
	var newContract model.DBContract
	newContract.Acceptor_id = int(contr.Acceptor_id)
	newContract.Assignee_id = int(contr.Assignee_id)
	newContract.Availability = obj.Model.ContractAvailStr2Int(contr.Availability)
	newContract.Buyout = contr.Buyout
	newContract.Collateral = contr.Collateral
	newContract.Contract_id = int(contr.Contract_id)
	newContract.Date_accepted = util.ConvertTimeStrToInt(contr.Date_accepted)
	newContract.Date_completed = util.ConvertTimeStrToInt(contr.Date_completed)
	newContract.Date_expired = util.ConvertTimeStrToInt(contr.Date_expired)
	newContract.Date_issued = util.ConvertTimeStrToInt(contr.Date_issued)
	newContract.Days_to_complete = int(contr.Days_to_complete)
	newContract.End_location_id = contr.End_location_id
	newContract.For_corporation = contr.For_corporation
	newContract.Issuer_corporation_id = int(contr.Issuer_corporation_id)
	newContract.Issuer_id = int(contr.Issuer_id)
	newContract.Price = contr.Price
	newContract.Reward = contr.Reward
	newContract.Start_location_id = contr.Start_location_id
	newContract.Status = obj.Model.ContractStatusStr2Int(contr.Status)
	newContract.Title = obj.Model.AddStringEntry(contr.Title)
	newContract.Type = obj.Model.ContractTypeStr2Int(contr.Type)
	newContract.Volume = contr.Volume
	return &newContract
}

func (obj *Ctrl) UpdateContractItems(char *EsiChar, corp bool) {
	if !char.UpdateFlags.Contracts {
		return
	}
	DBContrList := make([]*model.DBContract, 0, 5)
	if corp {
		DBContrList = append(DBContrList, obj.Model.GetContractsByIssuerId(char.CharInfoExt.CooperationId, true)...)
	} else {
		DBContrList = obj.Model.GetContractsByIssuerId(char.CharInfoData.CharacterID, false)
	}
	for _, contract := range DBContrList {
		if !obj.Model.ContrItemsExist(contract.Contract_id) {
			obj.GetContractItems(char, corp, contract.Contract_id)
		}
	}
}

func (obj *Ctrl) GetContractItems(char *EsiChar, corp bool, contractID int) {
	var url string
	if corp {
		url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/contracts/%d/items/", char.CharInfoExt.CooperationId, contractID)
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/contracts/%d/items/", char.CharInfoData.CharacterID, contractID)
	}
	bodyBytes, _, resp := obj.getSecuredUrl(url, char)
	if resp == nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR receiving url %s", url))
	} else {
		if resp.StatusCode == 404 {
			obj.AddLogEntry(fmt.Sprintf("ERROR contract item markt as unkown url %d", contractID))
			var DbContrItem model.DBContrItem
			DbContrItem.Contract_id = contractID
			obj.Model.AddContrItemEntry(&DbContrItem)
		} else {
			var cntrItemLst []CntrItems
			contentError := json.Unmarshal(bodyBytes, &cntrItemLst)
			if contentError != nil {
				obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			} else {
				for _, cntrItem := range cntrItemLst {
					DbContrItem := obj.convertEsiCntrItem2DB(&cntrItem, contractID)
					obj.Model.AddContrItemEntry(DbContrItem)
				}
			}
		}

	}

}

func (obj *Ctrl) convertEsiCntrItem2DB(cntrItem *CntrItems, contractID int) *model.DBContrItem {
	var newContrItem model.DBContrItem
	newContrItem.Contract_id = contractID
	if cntrItem.Is_included {
		newContrItem.Is_included = 1
	}
	if cntrItem.Is_singleton {
		newContrItem.Is_singleton = 1
	}
	newContrItem.Quantity = cntrItem.Quantity
	newContrItem.Raw_quantity = cntrItem.Raw_quantity
	newContrItem.Record_id = cntrItem.Record_id
	newContrItem.Type_id = cntrItem.Type_id
	return &newContrItem
}

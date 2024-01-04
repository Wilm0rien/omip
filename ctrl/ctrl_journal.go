package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"regexp"
	"strconv"
	"time"
)

type Journal struct {
	Amount          float64 `json:"amount"`
	Balance         float64 `json:"balance"`
	Context_id      int64   `json:"context_id"`
	Context_id_type string  `json:"context_id_type"`
	Date            string  `json:"date"`
	Description     string  `json:"description"`
	First_party_id  int64   `json:"first_party_id"`
	Id              int64   `json:"id"`
	Reason          string  `json:"reason"`
	Ref_type        string  `json:"ref_type"`
	Second_party_id int64   `json:"second_party_id"`
	Tax             float64 `json:"tax"`
	Tax_receiver_id int64   `json:"tax_receiver_id"`
}

func (obj *Ctrl) UpdateJournal(char *EsiChar, corp bool, division int) {
	var url string

	pageID := 1
	for {
		if corp {
			url = fmt.Sprintf("https://esi.evetech.net/v4/corporations/%d/wallets/%d/journal?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, division, pageID)
		} else {
			url = fmt.Sprintf("https://esi.evetech.net/v6/characters/%d/wallet/journal/?datasource=tranquility&page=%d", char.CharInfoData.CharacterID, pageID)
			division = 0
		}
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var journalList []Journal
		contentError := json.Unmarshal(bodyBytes, &journalList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			break
		}
		for _, jourEntry := range journalList {
			jDBEntry := obj.convertEsiJournal2DB(&jourEntry, char.CharInfoExt.CooperationId, char.CharInfoData.CharacterID, division)
			jDBEntry.WalletDivID = division
			obj.ParseJournalDescription(&jourEntry)
			obj.Model.AddJournalEntry(jDBEntry)
		}
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}
}

func (obj *Ctrl) ParseJournalDescription(jourEntry *Journal) {
	jobRe := regexp.MustCompile(`[(]Job ID: ([0-9]+)[)]`)
	jobResult := jobRe.FindStringSubmatch(jourEntry.Description)
	if jobResult != nil {
		jobID, err := strconv.ParseInt(jobResult[1], 10, 64)
		if err == nil {
			obj.Model.AddJournalLinkEntry(
				&model.DBJournalLink{JournalID: jourEntry.Id, IndustryJobID: jobID, ContractID: 0})
		}
	}
	ctrRe := regexp.MustCompile(`[(]contract ID: ([0-9]+)[)]`)
	ctrResult := ctrRe.FindStringSubmatch(jourEntry.Description)
	if ctrResult != nil {
		ctrID, err := strconv.ParseInt(ctrResult[1], 10, 64)
		if err == nil {
			obj.Model.AddJournalLinkEntry(
				&model.DBJournalLink{JournalID: jourEntry.Id, IndustryJobID: 0, ContractID: ctrID})
		}
	}
}

func (obj *Ctrl) convertEsiJournal2DB(jourEntry *Journal, corpId int, charId int, divId int) *model.DBJournal {
	var newJournal model.DBJournal
	newJournal.CharID = charId
	newJournal.CorpID = corpId
	newJournal.WalletDivID = divId
	newJournal.Amount = jourEntry.Amount
	newJournal.Balance = jourEntry.Balance
	newJournal.Context_id = jourEntry.Context_id
	newJournal.Context_id_type, _ = model.JournalContextIdType[jourEntry.Context_id_type]
	newJournal.Date = util.ConvertTimeStrToInt(jourEntry.Date)
	newJournal.Description = obj.Model.AddStringEntry(jourEntry.Description)
	newJournal.First_party_id = jourEntry.First_party_id
	newJournal.ID = jourEntry.Id
	newJournal.Ref_type, _ = model.JournalRefType[jourEntry.Ref_type]

	if jourEntry.Ref_type == "bounty_prizes" {
		newJournal.Reason = 0
	} else {
		newJournal.Reason = obj.Model.AddStringEntry(jourEntry.Reason)
	}

	newJournal.Second_party_id = jourEntry.Second_party_id
	newJournal.Tax = jourEntry.Tax
	newJournal.Tax_receiver_id = jourEntry.Tax_receiver_id
	return &newJournal
}

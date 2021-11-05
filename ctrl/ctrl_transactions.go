package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"log"
)

type Transaction struct {
	ClientID      int     `json:"client_id"`
	Date          string  `json:"date"`
	IsBuy         bool    `json:"is_buy"`
	JournalRefID  int64   `json:"journal_ref_id"`
	LocationID    int64   `json:"location_id"`
	Quantity      int     `json:"quantity"`
	TransactionID int64   `json:"transaction_id"`
	TypeID        int     `json:"type_id"`
	UnitPrice     float64 `json:"unit_price"`
}

func (obj *Ctrl) UpdateTransaction(char *EsiChar, corp bool) {
	if !char.UpdateFlags.Journal {
		return
	}
	var url string
	division := 1
	var charId, corpId int
	var ticker string
	corpInfo2 := obj.GetCorp(corpId)
	if corpInfo2 != nil {
		ticker = corpInfo2.Ticker
	}
	if corp {
		url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/wallets/%d/transactions?datasource=tranquility", char.CharInfoExt.CooperationId, division)
		corpId = char.CharInfoExt.CooperationId
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/wallet/transactions/?datasource=tranquility", char.CharInfoData.CharacterID)
		charId = char.CharInfoData.CharacterID
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var transactionList []Transaction
	contentError := json.Unmarshal(bodyBytes, &transactionList)
	if contentError != nil {
		log.Printf("ERROR reading url %s", url)
		return
	}
	for _, trEntry := range transactionList {
		tDBEntry := obj.convertEsiTransaction2DB(&trEntry, corpId, charId, division)
		result := obj.Model.AddTransactionEntry(tDBEntry)
		if result == model.DBR_Inserted {
			var logEntry string
			typeName := obj.Model.GetTypeString(trEntry.TypeID)
			if typeName == "" {
				typeName = "Unknown Item Type"
			}
			if !corp {
				if ticker == "" {
					logEntry = fmt.Sprintf("%s", char.CharInfoData.CharacterName)
				} else {
					logEntry = fmt.Sprintf("[%s] %s", ticker, char.CharInfoData.CharacterName)
				}

			} else {
				logEntry = fmt.Sprintf("[%s] %s", ticker, corpInfo2.Name)
			}
			if trEntry.IsBuy {
				logEntry += fmt.Sprintf(" bought")
			} else {
				logEntry += fmt.Sprintf(" sold")
			}
			logEntry += fmt.Sprintf(" %d of units %s for %3.2fM", trEntry.Quantity, typeName,
				float64(trEntry.Quantity)*float64(trEntry.UnitPrice)/1000000)

			obj.AddLogEntry(logEntry)
		}
	}
}
func (obj *Ctrl) convertEsiTransaction2DB(trEntry *Transaction, corpId int, charId int, divId int) *model.DBTransaction {
	var newTr model.DBTransaction
	newTr.CharID = charId
	newTr.CorpID = corpId
	newTr.ClientID = trEntry.ClientID
	newTr.Date = util.ConvertTimeStrToInt(trEntry.Date)
	if trEntry.IsBuy {
		newTr.IsBuy = 1
	}
	newTr.JournalRefID = trEntry.JournalRefID
	newTr.LocationID = trEntry.LocationID
	newTr.Quantity = trEntry.Quantity
	newTr.TransactionID = trEntry.TransactionID
	newTr.TypeID = trEntry.TypeID
	newTr.UnitPrice = trEntry.UnitPrice
	return &newTr
}

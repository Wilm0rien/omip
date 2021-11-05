package ctrl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CorpWallet struct {
	Balance  float64 `json:"balance"`
	Division int     `json:"division"`
}

func (obj *Ctrl) UpdateWallet(char *EsiChar, corp bool) {
	var url string
	var balance float64
	if corp {
		url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/wallets/?datasource=tranquility", char.CharInfoExt.CooperationId)
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/wallet/?datasource=tranquility", char.CharInfoData.CharacterID)
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	if bodyBytes != nil {
		if corp {
			var corpWallets []CorpWallet
			err := json.Unmarshal(bodyBytes, &corpWallets)
			if err == nil {
				for _, wallet := range corpWallets {
					corpInfo := obj.GetCorp(char.CharInfoExt.CooperationId)
					obj.handleWallet(corpInfo.Name, wallet.Balance, 0, char.CharInfoExt.CooperationId, wallet.Division)
					balance += wallet.Balance
				}
			}
		} else {
			if newBalance, err := strconv.ParseFloat(string(bodyBytes), 64); err == nil {
				obj.handleWallet(char.CharInfoData.CharacterName, newBalance, char.CharInfoData.CharacterID, 0, 0)
				balance += newBalance
			}
		}
	}
}

func (obj *Ctrl) handleWallet(nameStr string, newBalance float64, charactr_id int, corporation_id int, division int) {
	oldBalance := obj.Model.GetLatestWallets(charactr_id, corporation_id, division)
	if newBalance != oldBalance {
		obj.Model.AddWalletEntry(charactr_id, corporation_id, division, newBalance)
		diff := newBalance - oldBalance
		if diff > 1000000 {
			var logEntry string
			if division == 0 {
				logEntry = fmt.Sprintf("%s Balance %3.3fM Change %3.3fM", nameStr, newBalance/1000000, diff/1000000)
			} else {
				logEntry = fmt.Sprintf("%s Division %d Balance %3.3fM Change %3.3fM", nameStr, division, newBalance/1000000, diff/1000000)
			}
			obj.AddLogEntry(logEntry)
		}
	}
}

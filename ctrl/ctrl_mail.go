package ctrl

import (
	"encoding/json"
	"fmt"
)

type MailLabel struct {
	Color       string `json:"color"`
	LabelId     int    `json:"label_id"`
	Name        string `json:"name"`
	UnreadCount int    `json:"unread_count"`
}

func (obj *Ctrl) UpdateMailLabels(char *EsiChar) {
	url := fmt.Sprintf("https://esi.evetech.net/v3/characters/%d/mail/labels/?datasource=tranquility", char.CharInfoExt.CooperationId)
	if !char.UpdateFlags.MailLabels {
		return
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var labels []MailLabel
	contentError := json.Unmarshal(bodyBytes, &labels)
	if bodyBytes == nil {
		return
	}
	if contentError != nil {
		for _, label := range labels {
			if label.UnreadCount > 0 {
				namePrefix := fmt.Sprintf("[%s] %s: ", obj.GetCorpTicker(char), char.CharInfoData.CharacterName)
				obj.AddLogEntry(namePrefix + fmt.Sprintf("unread mails in: %s", label.Name))
			}
		}
	}
}

package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"net/http"
)

type MailLabel struct {
	Color       string `json:"color"`
	LabelId     int    `json:"label_id"`
	Name        string `json:"name"`
	UnreadCount int    `json:"unread_count"`
}

// MailStatus struct for mail labels https://esi.evetech.net/v3/characters/%d/mail/labels/?datasource=tranquility
type MailStatus struct {
	Labels           []MailLabel `json:"labels"`
	TotalUnreadCount int         `json:"total_unread_count"`
}

type Recipients struct {
	RecipientId   int    `json:"recipient_id"`
	RecipientType string `json:"recipient_type"`
}

type Mail struct {
	From       int          `json:"from"`
	IsRead     bool         `json:"is_read"`
	Labels     []int        `json:"labels"`
	MailId     int          `json:"mail_id"`
	Recipients []Recipients `json:"recipients"`
	Subject    string       `json:"subject"`
	Timestamp  string       `json:"timestamp"`
}

func (obj *Ctrl) UpdateMailLabels(char *EsiChar, corp bool) {
	if corp {
		return
	}
	url := fmt.Sprintf("https://esi.evetech.net/v3/characters/%d/mail/labels/?datasource=tranquility", char.CharInfoData.CharacterID)

	if !char.UpdateFlags.MailLabels {
		return
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var mailStatus MailStatus
	fmt.Sprintf("%s", string(bodyBytes))
	contentError := json.Unmarshal(bodyBytes, &mailStatus)
	if bodyBytes == nil {
		return
	}
	if contentError == nil {
		unreadFound := false
		tooMuchFound := false
		for _, label := range mailStatus.Labels {
			if label.UnreadCount > 0 && label.UnreadCount <= 3 {
				//namePrefix := fmt.Sprintf("[%s] %s: ", obj.GetCorpTicker(char), char.CharInfoData.CharacterName)
				//obj.AddLogEntry(namePrefix + fmt.Sprintf("unread mails in: %s", label.Name))
				unreadFound = true
			} else if label.UnreadCount > 3 {
				namePrefix := fmt.Sprintf("[%s] %s: ", obj.GetCorpTicker(char), char.CharInfoData.CharacterName)
				obj.AddLogEntry(namePrefix + fmt.Sprintf("%d unread mails in: %s", label.UnreadCount, label.Name))
				tooMuchFound = true
			}
		}
		if unreadFound && !tooMuchFound {
			obj.UpdateMail(char)
		}
	}
}

func (obj *Ctrl) UpdateMail(char *EsiChar) {
	url := fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/mail/?datasource=tranquility", char.CharInfoData.CharacterID)
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var mails []Mail
	contentError := json.Unmarshal(bodyBytes, &mails)
	if bodyBytes == nil {
		return
	}
	if contentError == nil {
		obj.fetchNewNames(mails, char)
		for _, mail := range mails {
			if !mail.IsRead {
				namePrefix := fmt.Sprintf("[%s] %s: ", obj.GetCorpTicker(char), char.CharInfoData.CharacterName)
				from := obj.Model.GetNameByID(mail.From)
				obj.AddLogEntry(namePrefix + fmt.Sprintf("unread mail from: %s; Subject: %s", from, mail.Subject))
			}
		}
	}
}

func (obj *Ctrl) fetchNewNames(mails []Mail, char *EsiChar) {
	unknownIdList := make([]int, 0, 10)
	for _, mail := range mails {
		if !mail.IsRead {
			if !obj.Model.NameExists(mail.From) {
				unknownIdList = append(unknownIdList, mail.From)
			}
		}
	}
	if len(unknownIdList) > 0 {
		obj.getUniverseNames(unknownIdList, char)
	}
}

func (obj *Ctrl) getUniverseNames(unknownIdList []int, char *EsiChar) {
	url := fmt.Sprintf("https://esi.evetech.net/v3/universe/names/")
	reqMap := make(map[int]int)
	for _, elem := range unknownIdList {
		reqMap[elem] = 1
	}
	idx := 0
	reqLst := ""
	for key, _ := range reqMap {
		reqLst += fmt.Sprintf("%d", key)

		if idx != len(reqMap)-1 {
			reqLst += ", "
		}
		idx++
	}
	memberIDs := fmt.Sprintf("[%s]", reqLst)
	bodyBytes2, resp := obj.getSecuredUrlPost(url, memberIDs, char)
	if resp.StatusCode == http.StatusOK {
		var universEntries []universeNames
		err := json.Unmarshal(bodyBytes2, &universEntries)
		if err == nil {
			for _, uEntry := range universEntries {
				dbUEntry := obj.convertEsiUEntry2DB(&uEntry)
				obj.Model.AddNameEntry(dbUEntry)
			}
		}
	}
}

func (obj *Ctrl) convertEsiUEntry2DB(u *universeNames) *model.DBUniverseName {
	var newUEntry model.DBUniverseName
	newUEntry.ID = u.ID
	newUEntry.NameRef = obj.Model.AddStringEntry(u.Name)
	newUEntry.Category = obj.Model.AddStringEntry(u.Category)
	return &newUEntry
}

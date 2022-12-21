package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"time"
)

type Notification struct {
	Is_read         bool   `json:"is_read"`
	Notification_id int64  `json:"notification_id"`
	Sender_id       int32  `json:"sender_id"`
	Sender_type     string `json:"sender_type"`
	Text            string `json:"text"`
	Timestamp       string `json:"timestamp"`
	Type            string `json:"type"`
}

func (obj *Ctrl) UpdateNotifications(char *EsiChar, corp bool) {
	if !corp {
		corpInfo2 := obj.GetCorp(char)
		var ticker string
		if corpInfo2 != nil {
			ticker = corpInfo2.Ticker
		}

		last48h := time.Now().Unix() - 2*24*60*60

		url := fmt.Sprintf("https://esi.evetech.net/v6/characters/%d/notifications/?datasource=tranquility", char.CharInfoData.CharacterID)
		bodyBytes, _ := obj.getSecuredUrl(url, char)

		var notiList []Notification
		contentError := json.Unmarshal(bodyBytes, &notiList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
			return
		}
		for _, noti := range notiList {
			dbNoti := obj.convertEsiNoti2DB(&noti, char.CharInfoData.CharacterID)
			obj.Model.AddNotificationEntry(dbNoti)
			if _, ok := obj.NotifyInfo[dbNoti.NotificationId]; !ok {
				obj.NotifyInfo[dbNoti.NotificationId] = true
				if dbNoti.Type == model.NotiMsgTyp_StructureUnderAttack {
					if dbNoti.TimeStamp > last48h {
						obj.AddLogEntry(fmt.Sprintf("(Eve Time %s)\n[%s] %s\n ATTACK: %s", noti.Timestamp, ticker, char.CharInfoData.CharacterName, noti.Type))
					}
				}
				if dbNoti.Type == model.NotiMsgTyp_WarDeclared {
					if dbNoti.TimeStamp > last48h {
						obj.AddLogEntry(fmt.Sprintf("(Eve Time %s) [%s] %s\n WAR DECLARED: %s", noti.Timestamp, ticker, char.CharInfoData.CharacterName, noti.Type))
					}
				}
			}
		}
	}
}

func (obj *Ctrl) convertEsiNoti2DB(notiItem *Notification, charId int) *model.DBNotification {
	var newNoti model.DBNotification
	newNoti.CharId = charId
	if notiItem.Is_read {
		newNoti.IsRead = 1
	}
	newNoti.NotificationId = notiItem.Notification_id
	newNoti.SenderId = notiItem.Sender_id
	newNoti.SenderType = model.NotiSndTyp[notiItem.Sender_type]
	newNoti.TextRef = 0 // todo: store text efficiently
	newNoti.TimeStamp = util.ConvertTimeStrToInt(notiItem.Timestamp)
	newNoti.Type = model.NotiMsgTyp[notiItem.Type]
	return &newNoti
}

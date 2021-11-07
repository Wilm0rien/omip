package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
)

type KillMailsRecent_t struct {
	KillMailHash string `json:"killmail_hash"`
	KillmailID   int32  `json:"killmail_id"`
}

type KillMail_t struct {
	Attackers     []Attackers_t `json:"attackers"`
	KillMailID    int32         `json:"killmail_id"`
	KillMailTime  string        `json:"killmail_time"`
	MoonID        int32         `json:"moon_id"`
	SolarSystemID int32         `json:"solar_system_id"`
	Victim        Victim_t      `json:"victim"`
	WarID         int32         `json:"war_id"`
}

type Attackers_t struct {
	AllianceID     int32   `json:"alliance_id"`
	CharacterID    int32   `json:"character_id"`
	CorporationID  int32   `json:"corporation_id"`
	DamageDone     int32   `json:"damage_done"`
	FactionID      int32   `json:"faction_id"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float32 `json:"security_status"`
	ShipTypeID     int32   `json:"ship_type_id"`
	WeaponTypeID   int32   `json:"weapon_type_id"`
}

type Victim_t struct {
	AllianceID    int32      `json:"alliance_id"`
	CharacterID   int32      `json:"character_id"`
	CorporationID int32      `json:"corporation_id"`
	DamageTaken   int32      `json:"damage_taken"`
	FactionID     int32      `json:"faction_id"`
	Items         []Items_t  `json:"items"`
	Position      Position_t `json:"position"`
	ShipTypeID    int32      `json:"ship_type_id"`
}

type Position_t struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Items_t struct {
	Flag              int32     `json:"flag"`
	ItemTypeID        int32     `json:"item_type_id"`
	SubItems          []Items_t `json:"items"`
	QuantityDestroyed int32     `json:"quantity_destroyed"`
	QuantityDropped   int32     `json:"quantity_dropped"`
	Singleton         int32     `json:"singleton"`
}

func (obj *Ctrl) InitiateKMSkipList(char *EsiChar, corp bool) {
	var url string
	if corp {
		url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/killmails/recent/", char.CharInfoExt.CooperationId)
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/killmails/recent/", char.CharInfoData.CharacterID)
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var kmList []KillMailsRecent_t
	contentError := json.Unmarshal(bodyBytes, &kmList)
	if contentError != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
	} else {
		kmSkipList := make(map[int32]bool)
		for _, km := range kmList {
			kmSkipList[km.KillmailID] = true
		}
		if corp {
			corpObj := obj.GetCorp(char)
			if corpObj != nil {
				corpObj.KmSkipList = kmSkipList
			}
		} else {
			char.KmSkipList = kmSkipList
		}
	}
}

func (obj *Ctrl) UpdateKillMails(char *EsiChar, corp bool) {
	if !char.UpdateFlags.Killmails {
		return
	}
	var url string
	var killsBefore, killsLater, lossesBefore, lossesLater int
	memberList := obj.Model.GetDBCorpMembers(char.CharInfoExt.CooperationId)
	nameMapping := make(map[int]string)
	lossMapping := make(map[string]float64)
	for _, member := range memberList {
		name, _ := obj.Model.GetStringEntry(member.NameRef)
		nameMapping[member.CharID] = name
	}
	var kmSkipList map[int32]bool
	if corp {
		killsBefore, lossesBefore = obj.captureKmNumbers(char.CharInfoExt.CooperationId)
		url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/killmails/recent/", char.CharInfoExt.CooperationId)
		corpObj := obj.GetCorp(char)
		if corpObj != nil {
			kmSkipList = corpObj.KmSkipList
		}
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/killmails/recent/", char.CharInfoData.CharacterID)
		kmSkipList = char.KmSkipList
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var kmList []KillMailsRecent_t
	contentError := json.Unmarshal(bodyBytes, &kmList)
	if contentError != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
	} else {
		for _, km := range kmList {
			if _, ok := kmSkipList[km.KillmailID]; !ok {
				var newKM model.DBKillmail
				newKM.Killmail_id = km.KillmailID
				newKM.Killmail_hash = km.KillMailHash
				if !obj.Model.KillMailExists(newKM.Killmail_id) {
					obj.GetKillMail(char, &newKM, nameMapping, lossMapping)
				}
			}
		}
	}
	if corp {
		killsLater, lossesLater = obj.captureKmNumbers(char.CharInfoExt.CooperationId)
		obj.kmLogEntry(char, killsBefore, killsLater, lossesBefore, lossesLater, lossMapping)
	}
}

func (obj *Ctrl) captureKmNumbers(corpId int) (kills int, losses int) {
	kills = len(obj.Model.GetKillsDB(corpId, false))
	losses = len(obj.Model.GetKillsDB(corpId, true))
	return kills, losses
}
func (obj *Ctrl) kmLogEntry(char *EsiChar, kB int, kL int, lB int, lL int, lossMapping map[string]float64) {
	corpName := "N/A"
	corpTicker := "N/A"
	corpObj := obj.GetCorp(char)
	if corpObj != nil {
		corpName = corpObj.Name
		corpTicker = corpObj.Ticker
	}
	if kB != kL && lB != lL {
		obj.AddLogEntry(fmt.Sprintf("%s kills +%d losses +%d",
			corpName, kL-kB, lL-lB))
	} else {
		if kB != kL {
			obj.AddLogEntry(fmt.Sprintf("%s kills +%d",
				corpName, kL-kB))
		}
		if lB != lL {
			obj.AddLogEntry(fmt.Sprintf("%s losses +%d",
				corpName, lL-lB))
		}
	}
	for _, charName := range util.GetSortKeysFromStrMap(lossMapping, false) {
		obj.AddLogEntry(fmt.Sprintf("[%s] %s lost %s", corpTicker,
			charName, util.HumanizeNumber(lossMapping[charName])))
	}

}

func (obj *Ctrl) GetKillMail(char *EsiChar, dbKM *model.DBKillmail,
	nameMapping map[int]string, lossMapping map[string]float64) {

	url := fmt.Sprintf("https://esi.evetech.net/v1/killmails/%d/%s", dbKM.Killmail_id, dbKM.Killmail_hash)
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	if bodyBytes == nil {
		return
	}
	//fmt.Printf("%s", string(bodyBytes))
	var esiKM KillMail_t
	contentError := json.Unmarshal(bodyBytes, &esiKM)
	if contentError != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
	} else {
		newKM := convertEsiKM2DB(&esiKM, dbKM.Killmail_hash)

		for _, attacker := range esiKM.Attackers {
			newAttacker := convertEsiAttacker2DB(&attacker, dbKM.Killmail_id)
			obj.Model.AddKAttackerEntry(newAttacker)
		}
		newVictim := convertEsiVictim2DB(&esiKM.Victim, dbKM.Killmail_id)
		obj.Model.AddKVictimEntry(newVictim)

		var value float64
		for _, item := range esiKM.Victim.Items {
			newItem := convertEsiItem2DB(&item, dbKM.Killmail_id)
			value += obj.Model.GetItemValue(int(item.ItemTypeID))

			obj.Model.AddKItemEntry(newItem)
			for _, subItem := range item.SubItems {
				newSubItem := convertEsiItem2DB(&subItem, dbKM.Killmail_id)
				value += obj.Model.GetItemValue(int(newSubItem.Item_type_id))
				obj.Model.AddKItemEntry(newSubItem)
			}
		}
		newKM.Value = value
		db1R := obj.Model.AddKillmailEntry(newKM)
		util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)

		if charName, ok1 := nameMapping[int(newVictim.Character_id)]; ok1 {
			lossMapping[charName] += value
		}
	}
}

func convertEsiKM2DB(km *KillMail_t, km_hash string) *model.DBKillmail {
	var newKM model.DBKillmail
	newKM.Killmail_id = km.KillMailID
	newKM.Killmail_hash = km_hash
	newKM.Killmail_time = util.ConvertTimeStrToInt(km.KillMailTime)
	newKM.Moon_id = km.MoonID
	newKM.Solar_system_id = km.SolarSystemID
	newKM.War_id = km.WarID
	return &newKM
}

func convertEsiAttacker2DB(att *Attackers_t, km_id int32) *model.DBKAttacker {
	var newDBAtt model.DBKAttacker
	newDBAtt.Killmail_id = km_id
	newDBAtt.Alliance_id = att.AllianceID
	newDBAtt.Character_id = att.CharacterID
	newDBAtt.Corporation_id = att.CorporationID
	newDBAtt.Damage_done = att.DamageDone
	newDBAtt.Faction_id = att.FactionID
	if att.FinalBlow {
		newDBAtt.Final_blow = 1
	} else {
		newDBAtt.Final_blow = 0
	}
	newDBAtt.Security_status = att.SecurityStatus
	newDBAtt.Ship_type_id = att.ShipTypeID
	newDBAtt.Weapon_type_id = att.WeaponTypeID
	return &newDBAtt
}

func convertEsiVictim2DB(victim *Victim_t, km_id int32) *model.DBKVictim {
	var newDBVictim model.DBKVictim
	newDBVictim.Killmail_id = km_id
	newDBVictim.Alliance_id = victim.AllianceID
	newDBVictim.Character_id = victim.CharacterID
	newDBVictim.Corporation_id = victim.CorporationID
	newDBVictim.Damage_taken = victim.DamageTaken
	newDBVictim.Faction_id = victim.FactionID
	newDBVictim.Position_x = victim.Position.X
	newDBVictim.Position_y = victim.Position.Y
	newDBVictim.Position_z = victim.Position.Z
	newDBVictim.Ship_type_id = victim.ShipTypeID
	return &newDBVictim
}

func convertEsiItem2DB(item *Items_t, km_id int32) *model.DBKItem {
	var newDBItem model.DBKItem
	newDBItem.Killmail_id = km_id
	newDBItem.Flag = item.Flag
	newDBItem.Item_type_id = item.ItemTypeID
	newDBItem.Quantity_destroyed = item.QuantityDestroyed
	newDBItem.Quantity_dropped = item.QuantityDropped
	newDBItem.Singleton = item.Singleton
	return &newDBItem
}

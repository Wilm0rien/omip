package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"log"
	"time"
)

type ZkStatus int

const (
	ZkUnkown ZkStatus = iota
	ZkOK
	ZkNOTOK
)

type DBKillmail struct {
	Killmail_id     int32
	Killmail_hash   string
	Killmail_time   int64
	Moon_id         int32
	Solar_system_id int32
	War_id          int32
	Value           float64
	ZK_status       ZkStatus
}

type DBKAttacker struct {
	Killmail_id     int32
	Alliance_id     int32
	Character_id    int32
	Corporation_id  int32
	Damage_done     int32
	Faction_id      int32
	Final_blow      int
	Security_status float32
	Ship_type_id    int32
	Weapon_type_id  int32
}

type DBKVictim struct {
	Killmail_id    int32
	Alliance_id    int32
	Character_id   int32
	Corporation_id int32
	Damage_taken   int32
	Faction_id     int32
	Position_x     float64
	Position_y     float64
	Position_z     float64
	Ship_type_id   int32
}

type DBKItem struct {
	Killmail_id        int32
	Flag               int32
	Item_type_id       int32
	Quantity_destroyed int32
	Quantity_dropped   int32
	Singleton          int32
}

type DBKillTable struct {
	Time     int64
	MainName string
}

func (obj *Model) createKillmails() {
	if !obj.checkTableExists("killmails") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "killmails" (
			killmail_id     INT,
			killmail_hash   TEXT,
			killmail_time   INT,
			moon_id         INT,
			solar_system_id INT,
			war_id          INT,
			value           REAL,
			zk_status	    INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createKAttackersTable() {
	if !obj.checkTableExists("k_attackers") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "k_attackers" (
			killmail_id INT,
			alliance_id INT,
			character_id INT,
			corporation_id INT,
			damage_done INT,
			faction_id INT,
			final_blow INT,
			security_status REAL,
			ship_type_id INT,
			weapon_type_id INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createKVictimsTable() {
	if !obj.checkTableExists("k_victims") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "k_victims" (
			killmail_id    INT,
			alliance_id    INT,
			character_id   INT,
			corporation_id INT,
			damage_taken   INT,
			faction_id     INT,
			position_x     REAL,
			position_y     REAL,
			position_z     REAL,
			ship_type_id   INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createKItemsTable() {
	if !obj.checkTableExists("k_items") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "k_items" (
			killmail_id        INT,
			flag               INT,
			item_type_id       INT,
			quantity_destroyed INT,
			quantity_dropped   INT,
			singleton          INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createKMTables() {
	obj.createKillmails()
	obj.createKAttackersTable()
	obj.createKVictimsTable()
	obj.createKItemsTable()
}

func (obj *Model) AddKillmailEntry(km *DBKillmail) DBresult {
	whereClause := fmt.Sprintf(`killmail_id="%d"`, km.Killmail_id)
	num := obj.getNumEntries("killmails", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "killmails" (
			killmail_id,
			killmail_hash,
			killmail_time,
			moon_id,
			solar_system_id,
			war_id,
		    value,
		    zk_status)
			VALUES (?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			km.Killmail_id,
			km.Killmail_hash,
			km.Killmail_time,
			km.Moon_id,
			km.Solar_system_id,
			km.War_id,
			km.Value,
			km.ZK_status)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {

			retval = DBR_Inserted
		} else {
			log.Printf("insert failed! %d", km.Killmail_id)
		}
	} else {
		if km.Killmail_time != 0 {
			stmt, err := obj.DB.Prepare(`
				UPDATE "killmails" SET
					killmail_hash=?,
					killmail_time =?,
					moon_id =?,
					solar_system_id =?,
					war_id =?,
				    value=?,
				    zk_status=?					
				WHERE killmail_id=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				km.Killmail_hash,
				km.Killmail_time,
				km.Moon_id,
				km.Solar_system_id,
				km.War_id,
				km.Value,
				km.ZK_status,
				km.Killmail_id)
			util.CheckErr(err)
			affect, err := res.RowsAffected()
			util.CheckErr(err)
			if affect > 0 {
				retval = DBR_Updated
			} else {
				log.Printf("km update failed %d", km.Killmail_id)
			}
		} else {
			retval = DBR_Skipped
		}
	}
	return retval
}

func (obj *Model) AddKAttackerEntry(att *DBKAttacker) DBresult {
	num := 0
	if att.Character_id != 0 {
		whereClause := fmt.Sprintf(`killmail_id=%d AND character_id=%d`, att.Killmail_id, att.Character_id)
		num = obj.getNumEntries("k_attackers", whereClause)
	}
	retval := DBR_Undefined
	//wilm0rutil.Assert(num == 0)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "k_attackers" (
			killmail_id,
			alliance_id,
			character_id,
			corporation_id,
			damage_done,
			faction_id,
			final_blow,
			security_status,
			ship_type_id,
			weapon_type_id)
			VALUES (?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			att.Killmail_id,
			att.Alliance_id,
			att.Character_id,
			att.Corporation_id,
			att.Damage_done,
			att.Faction_id,
			att.Final_blow,
			att.Security_status,
			att.Ship_type_id,
			att.Weapon_type_id)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) AddKVictimEntry(victim *DBKVictim) DBresult {
	whereClause := fmt.Sprintf(`killmail_id=%d AND character_id=%d`, victim.Killmail_id, victim.Character_id)
	num := obj.getNumEntries("k_victims", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "k_victims" (
			killmail_id,
			alliance_id,
			character_id,
			corporation_id,
			damage_taken,
			faction_id,
			position_x,
			position_y,
			position_z,
			ship_type_id)
			VALUES (?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			victim.Killmail_id,
			victim.Alliance_id,
			victim.Character_id,
			victim.Corporation_id,
			victim.Damage_taken,
			victim.Faction_id,
			victim.Position_x,
			victim.Position_y,
			victim.Position_z,
			victim.Ship_type_id)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) AddKItemEntry(item *DBKItem) DBresult {
	whereClause := fmt.Sprintf(`killmail_id=%d and item_type_id=%d`, item.Killmail_id, item.Item_type_id)
	num := obj.getNumEntries("k_items", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		// cannot use num entries because killmail_id is not unique
		// every item belonging to this killmail_id will be present here
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "k_items" (
				killmail_id,
				flag,
				item_type_id,
				quantity_destroyed,
				quantity_dropped,
				singleton)
				VALUES (?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			item.Killmail_id,
			item.Flag,
			item.Item_type_id,
			item.Quantity_destroyed,
			item.Quantity_dropped,
			item.Singleton)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) KillMailExists(killmailID int32) bool {
	whereClause := fmt.Sprintf(`killmail_id="%d"`, killmailID)
	num := obj.getNumEntries("killmails", whereClause)
	return (num > 0)
}

func (obj *Model) GetKillsMails() []*DBKillmail {
	retval := make([]*DBKillmail, 0, 10)
	queryString := fmt.Sprintf(`
		SELECT 
			killmail_id,
			killmail_hash,
			killmail_time,
			moon_id,
			solar_system_id,
			war_id,
			value,
			zk_status
		FROM killmails;`)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var km DBKillmail
		rows.Scan(
			&km.Killmail_id,
			&km.Killmail_hash,
			&km.Killmail_time,
			&km.Moon_id,
			&km.Solar_system_id,
			&km.War_id,
			&km.Value,
			&km.ZK_status,
		)
		retval = append(retval, &km)
	}
	return retval
}

func (obj *Model) GetKillsMail(kmId int) *DBKillmail {
	var retval *DBKillmail
	queryString := fmt.Sprintf(`
		SELECT 
			killmail_id,
			killmail_hash,
			killmail_time,
			moon_id,
			solar_system_id,
			war_id,
			value,
			zk_status
		FROM killmails
		WHERE killmail_id=%d`, kmId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var km DBKillmail
		rows.Scan(
			&km.Killmail_id,
			&km.Killmail_hash,
			&km.Killmail_time,
			&km.Moon_id,
			&km.Solar_system_id,
			&km.War_id,
			&km.Value,
			&km.ZK_status,
		)

		retval = &km
	}
	return retval
}

func (obj *Model) GetKillsDB(corpId int, losses bool) []*DBTable {
	var retval []*DBTable
	srcTable := "k_attackers"
	if losses == true {
		srcTable = "k_victims"
	}
	queryString := fmt.Sprintf(`
		SELECT killmails.killmail_time, string_table.string as charname, killmails.killmail_id, killmails.value from %s 
		INNER JOIN 
				   corp_members ON corp_members.character_id = %s.character_id
		INNER JOIN
					(SELECT character_id, name FROM corp_members) corpRef2Main
					ON corpRef2Main.character_id = corp_members.main_id
		INNER JOIN
			string_table ON corpRef2Main.name = string_table.string_hash
		INNER JOIN
			killmails ON  %s.killmail_id = killmails.killmail_id
		WHERE %s.corporation_id=%d
		ORDER BY killmails.killmail_time DESC;`, srcTable, srcTable, srcTable, srcTable, corpId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		var elem DBTable
		var killmailId int
		var value float64
		rows.Scan(&elem.Time, &elem.MainName, &killmailId, &value)
		if losses {
			elem.Amount = value
		} else {
			elem.Amount = 1
		}
		if elem.Amount != 0 {
			retval = append(retval, &elem)
		}

	}
	return retval
}

func (obj *Model) CheckNPCKill(kmid int) bool {
	whereClause := fmt.Sprintf("alliance_id=0 and character_id=0 and killmail_id=%d", kmid)
	num := obj.getNumEntries("k_attackers", whereClause)

	whereClaus2 := fmt.Sprintf("character_id!=0 and killmail_id=%d", kmid)
	num2 := obj.getNumEntries("k_attackers", whereClaus2)
	return num != 0 && num2 == 0

}

func (obj *Model) GetVictimData(corpId int, mainId int, startTS int64, endTS int64) []*DBKMTable {
	var retval []*DBKMTable

	queryString := fmt.Sprintf(`
		SELECT killmails.killmail_time, killmails.zk_status, string_table.string as charname, stringRefAlt.string as alt, killmails.killmail_id,killmails.killmail_hash, killmails.value from k_victims 
		INNER JOIN 
				   corp_members ON corp_members.character_id = k_victims.character_id
		INNER JOIN
					(SELECT character_id, name FROM corp_members) corpRef2Main
					ON corpRef2Main.character_id = corp_members.main_id
		INNER JOIN
					(SELECT character_id, name FROM corp_members) corpRef2Alt
					ON corpRef2Alt.character_id = corp_members.character_id
		INNER JOIN
			string_table ON corpRef2Main.name = string_table.string_hash	
		INNER JOIN
					(SELECT string, string_hash FROM string_table) stringRefAlt
					ON stringRefAlt.string_hash = corpRef2Alt.name			
		INNER JOIN
			killmails ON  k_victims.killmail_id = killmails.killmail_id
		WHERE k_victims.corporation_id=%d and corpRef2Main.character_id = %d and killmails.killmail_time>=%d and killmails.killmail_time<%d
		ORDER BY killmails.killmail_time DESC;`, corpId, mainId, startTS, endTS)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		var elem DBKMTable
		rows.Scan(&elem.Time, &elem.ZK_status, &elem.MainName, &elem.AltName, &elem.KMId, &elem.KMHash, &elem.Amount)
		retval = append(retval, &elem)
	}
	return retval
}

func (obj *Model) GetKillValue(killMailId int) float64 {
	var retval float64
	killmail := obj.GetKillsMail(killMailId)
	if killmail != nil {
		retval = killmail.Value
	}
	return retval
}

func (obj *Model) GetKillsCurrentMonth(corpId int, losses bool) int {
	var retval int
	now_year, now_month, _ := time.Now().Date()
	killtable := obj.GetKillsDB(corpId, losses)
	for _, elem := range killtable {
		tm := time.Unix(elem.Time, 0)
		year, month, _ := tm.Date()
		if year == now_year && month == now_month {
			retval++
		}
	}
	return retval
}

func (obj *Model) GetKillTable(corpId int, maxMonth int, losses bool) *MonthlyTable {
	killTable := obj.GetKillsDB(corpId, losses)
	return obj.GetMonthlyTable(corpId, killTable, maxMonth)

}

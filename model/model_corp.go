package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"log"
)

type DBcorpInfo struct {
	CorpID            int
	AllianceId        int
	CeoId             int
	CreatorId         int
	DateFounded       int64
	DescriptionStrRef int64
	FactionId         int
	HomeStationId     int
	MemberCount       int
	CorpNameStrRef    int64
	Shares            int64
	TaxRate           float64
	TickerStrRef      int64
	UrlStrRef         int64
	WarEligible       bool
}

type DBcorpMember struct {
	CharID  int
	CorpID  int
	MainID  int
	NameRef int64
	Updated bool
}

type DBCorpNames struct {
	CorpID     int
	AllyID     int
	CorpName   string
	CorpTicker string
	AllyName   string
	AllyTicker string
}

func (obj *Model) createCorpMemberTable() {
	if !obj.checkTableExists("corp_members") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "corp_members" (
			"character_id" INT,
			"corporation_id" INT,
			"main_id" INT,
			"name" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createCorpInfoTable() {
	if !obj.checkTableExists("corp_info") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "corp_info" (
			"corporation_id" INT,
			"alliance_id" INT,
			"ceo_id" INT,
			"creator_id" INT,
			"date_founded" INT,
			"description" INT,
			"faction_id" INT,
			"home_station_id" INT,
			"member_count" INT,
			"name" INT,
			"shares" INT,
			"tax_rate" REAL,
			"ticker" INT,
			"url" INT,
			"war_eligible" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) GetCorpInfoEntry(corpId int) (retval *DBcorpInfo, result DBresult) {
	result = DBR_Undefined
	whereClause := fmt.Sprintf(`corporation_id="%d"`, corpId)
	num := obj.getNumEntries("corp_info", whereClause)

	if num != 0 {
		var newDBCorpInfo DBcorpInfo
		queryStr := fmt.Sprintf(`SELECT 
										alliance_id,
										ceo_id,
										creator_id,
										date_founded,
										description,
										faction_id,
										home_station_id,
										member_count,
										name,
										shares,
										tax_rate,
										ticker,
										url,
										war_eligible
									FROM corp_info WHERE corporation_id=?;`)
		stmt, err := obj.DB.Prepare(queryStr)
		util.CheckErr(err)
		defer stmt.Close()
		rows, err := stmt.Query(corpId)
		util.CheckErr(err)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(
				&newDBCorpInfo.AllianceId,
				&newDBCorpInfo.CeoId,
				&newDBCorpInfo.CreatorId,
				&newDBCorpInfo.DateFounded,
				&newDBCorpInfo.DescriptionStrRef,
				&newDBCorpInfo.FactionId,
				&newDBCorpInfo.HomeStationId,
				&newDBCorpInfo.MemberCount,
				&newDBCorpInfo.CorpNameStrRef,
				&newDBCorpInfo.Shares,
				&newDBCorpInfo.TaxRate,
				&newDBCorpInfo.TickerStrRef,
				&newDBCorpInfo.UrlStrRef,
				&newDBCorpInfo.WarEligible)
			break
		}
		if newDBCorpInfo.DateFounded != 0 {
			retval = &newDBCorpInfo
			result = DBR_Success
		} else {
			result = DBR_Failed
		}
	} else {
		result = DBR_Skipped
	}
	return retval, result
}

func (obj *Model) GetCorpTicker(corpID int) string {
	var retval string
	info, result := obj.GetCorpInfoEntry(corpID)
	if result == DBR_Success {
		retval, _ = obj.GetStringEntry(info.TickerStrRef)
	}
	return retval
}

func (obj *Model) GetCorpNames(corpID int) *DBCorpNames {
	var retval DBCorpNames
	info, result := obj.GetCorpInfoEntry(corpID)
	if result == DBR_Success {
		retval.CorpTicker, _ = obj.GetStringEntry(info.TickerStrRef)
		retval.CorpName, _ = obj.GetStringEntry(info.CorpNameStrRef)
		retval.AllyID = info.AllianceId
		retval.CorpID = corpID
		if ally, result2 := obj.GetAllyInfoEntry(info.AllianceId); result2 == DBR_Success {
			retval.AllyName, _ = obj.GetStringEntry(ally.NameStrRef)
			retval.AllyTicker, _ = obj.GetStringEntry(ally.TickerStrRef)
		}
	} else {
		return nil
	}
	return &retval
}

func (obj *Model) AddCorpInfoEntry(corpInfo *DBcorpInfo) DBresult {
	whereClause := fmt.Sprintf(`corporation_id="%d"`, corpInfo.CorpID)
	num := obj.getNumEntries("corp_info", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "corp_info" (
			corporation_id,
			alliance_id,
			ceo_id,
			creator_id,
			date_founded,
			description,
			faction_id,
			home_station_id,
			member_count,
			name,
			shares,
			tax_rate,
			ticker,
			url,
			war_eligible)
			values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			corpInfo.CorpID,
			corpInfo.AllianceId,
			corpInfo.CeoId,
			corpInfo.CreatorId,
			corpInfo.DateFounded,
			corpInfo.DescriptionStrRef,
			corpInfo.FactionId,
			corpInfo.HomeStationId,
			corpInfo.MemberCount,
			corpInfo.CorpNameStrRef,
			corpInfo.Shares,
			corpInfo.TaxRate,
			corpInfo.TickerStrRef,
			corpInfo.UrlStrRef,
			corpInfo.WarEligible)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		stmt, err := obj.DB.Prepare(`
		UPDATE "corp_info" SET
			alliance_id=?,
		    ceo_id=?,
		    description=?,
		    faction_id=?,
		    home_station_id=?,
		    member_count=?,
		    tax_rate=?,
		    url=?
			WHERE corporation_id=?;`)
		util.CheckErr(err)
		res, err := stmt.Exec(
			corpInfo.AllianceId,
			corpInfo.CeoId,
			corpInfo.DescriptionStrRef,
			corpInfo.FactionId,
			corpInfo.HomeStationId,
			corpInfo.MemberCount,
			corpInfo.TaxRate,
			corpInfo.UrlStrRef, corpInfo.CorpID)
		util.CheckErr(err)
		defer stmt.Close()
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Updated
		}
	}
	return retval
}

func (obj *Model) UpdateMemberCount(memberCount int, corpId int) DBresult {
	retval := DBR_Undefined
	stmt, err := obj.DB.Prepare(`
		UPDATE "corp_info" SET
			member_count=?
			WHERE corporation_id=?;`)
	util.CheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(
		memberCount,
		corpId)
	util.CheckErr(err)
	defer stmt.Close()
	affect, err := res.RowsAffected()
	util.CheckErr(err)
	if affect > 0 {
		retval = DBR_Updated
	}
	return retval
}

func (obj *Model) AddCorpMemberEntry(member *DBcorpMember) DBresult {
	whereClause := fmt.Sprintf(`character_id="%d"`, member.CharID)
	num := obj.getNumEntries("corp_members", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "corp_members" (
			character_id,
			corporation_id,
			main_id,
			name)
			values(?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			member.CharID,
			member.CorpID,
			member.MainID,
			member.NameRef)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		retval = DBR_Skipped
	}
	return retval
}

func (obj *Model) RemoveCorpMemberEntry(charId int, corpId int) DBresult {
	retval := DBR_Skipped
	queryStr := fmt.Sprintf("DELETE FROM corp_members WHERE character_id=? AND corporation_id=?")
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(charId, corpId)
	util.CheckErr(err)
	affect, err := res.RowsAffected()
	util.CheckErr(err)
	if affect > 0 {
		retval = DBR_Removed
	}
	return retval
}

func (obj *Model) GetCorpMemberList(corpId int) []int {
	retval := make([]int, 0, 5)
	queryStr := fmt.Sprintf("SELECT character_id FROM corp_members WHERE corporation_id=?;")
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var charId int
		rows.Scan(&charId)
		retval = append(retval, charId)
	}
	return retval
}

func (obj *Model) GetCorpMemberNames(corpId int) []string {
	retval := make([]string, 0, 5)
	queryStr := fmt.Sprintf(
		`SELECT 	string_table.string as CharacterName
				FROM string_table
				JOIN corp_members ON string_table.string_hash = corp_members.name
				WHERE corp_members.corporation_id=?
				ORDER BY CharacterName COLLATE NOCASE;`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var charName string
		rows.Scan(&charName)
		retval = append(retval, charName)
	}
	return retval
}

// GetCorpMemberIdMap returns map[characterID]=1
// this map is for probing if characterID is member of this corp
func (obj *Model) GetCorpMemberIdMap(corpId int) map[int]int {
	memberIdMap := make(map[int]int)
	memberList := obj.GetCorpMemberList(corpId)
	for _, memberId := range memberList {
		memberIdMap[memberId] = 1
	}
	return memberIdMap
}

func (obj *Model) GetDBCorpMembers(corpId int) []*DBcorpMember {
	retval := make([]*DBcorpMember, 0, 5)
	queryStr := fmt.Sprintf(
		`SELECT 
					character_id, 
					corporation_id, 
					main_id, 
					name, 
					string_table.string as CharacterName
				FROM string_table
				JOIN corp_members ON string_table.string_hash = corp_members.name
				WHERE corp_members.corporation_id=?
				ORDER BY CharacterName COLLATE NOCASE;`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var m DBcorpMember
		var name string
		rows.Scan(&m.CharID, &m.CorpID, &m.MainID, &m.NameRef, &name)
		retval = append(retval, &m)
	}
	return retval
}

func (obj *Model) GetDBCorpMember(charID int) *DBcorpMember {
	var retval *DBcorpMember
	var m DBcorpMember
	queryStr := fmt.Sprintf(
		`SELECT
				character_id,
				corporation_id, 
				main_id, 
				name 
			FROM corp_members
			WHERE character_id=?;`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(charID)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&m.CharID, &m.CorpID, &m.MainID, &m.NameRef)
		retval = &m
		break
	}
	return retval
}

func (obj *Model) SetDBCorpMembers(memberList []*DBcorpMember) {
	for _, member := range memberList {
		stmt, err := obj.DB.Prepare(`
		UPDATE "corp_members" SET
			main_id=?
			WHERE character_id=?;`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			member.MainID,
			member.CharID)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect == 0 {
			log.Printf("SetDBCorpMembers Error updating database")
		}
	}
}
func (obj *Model) GetAltMap(corpId int) (alt2main map[string]string) {
	alt2main = make(map[string]string)
	queryStr := fmt.Sprintf(
		`SELECT string_table.string as charName, stringMain.string as mainName from corp_members
				INNER JOIN 
					string_table ON corp_members.name = string_table.string_hash
				INNER JOIN
					(SELECT character_id, name FROM corp_members) corpRef2Main
					ON corpRef2Main.character_id = corp_members.main_id		
				INNER JOIN
					(SELECT string, string_hash FROM string_table) stringMain
					ON stringMain.string_hash = corpRef2Main.name
				WHERE corporation_id = ?
				ORDER BY charName COLLATE NOCASE`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var charName string
		var mainName string
		rows.Scan(&charName, &mainName)
		alt2main[charName] = mainName
	}
	return alt2main
}

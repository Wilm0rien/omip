package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBallyInfo struct {
	AllianceID           int
	CreatorCorporationID int
	CreatorID            int
	DateFounded          int64
	ExecutorCorpID       int
	FactionID            int
	NameStrRef           int64
	TickerStrRef         int64
}

func (obj *Model) createAllyInfoTable() {
	if !obj.checkTableExists("ally_info") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "ally_info" (
		    "alliance_id" INT,
			"creator_corporation_id" INT,
			"creator_id" INT,
			"date_founded" INT,
			"executor_corporation_id" INT,
			"faction_id" INT,
			"name" INT,
			"ticker" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) GetAllyInfoEntry(allyID int) (retval *DBallyInfo, result DBresult) {
	result = DBR_Undefined
	whereClause := fmt.Sprintf(`alliance_id="%d"`, allyID)
	num := obj.getNumEntries("ally_info", whereClause)
	if num != 0 {
		var newDBAllyInfo DBallyInfo
		queryStr := fmt.Sprintf(`SELECT 
			alliance_id,
			creator_corporation_id,
			creator_id,
			date_founded,
			executor_corporation_id,
			faction_id,
			name,
			ticker
			FROM ally_info WHERE alliance_id=?;`)
		stmt, err := obj.DB.Prepare(queryStr)
		util.CheckErr(err)
		defer stmt.Close()
		rows, err := stmt.Query(allyID)
		util.CheckErr(err)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(
				&newDBAllyInfo.AllianceID,
				&newDBAllyInfo.CreatorCorporationID,
				&newDBAllyInfo.CreatorID,
				&newDBAllyInfo.DateFounded,
				&newDBAllyInfo.ExecutorCorpID,
				&newDBAllyInfo.FactionID,
				&newDBAllyInfo.NameStrRef,
				&newDBAllyInfo.TickerStrRef)
			break
		}
		if newDBAllyInfo.DateFounded != 0 {
			retval = &newDBAllyInfo
			result = DBR_Success
		} else {
			result = DBR_Failed
		}
	} else {
		result = DBR_Skipped
	}
	return retval, result
}

func (obj *Model) GetAllyNames(allyID int) (ticker string, name string) {
	info, result := obj.GetAllyInfoEntry(allyID)
	if result == DBR_Success {
		ticker, _ = obj.GetStringEntry(info.TickerStrRef)
		name, _ = obj.GetStringEntry(info.NameStrRef)
	}
	return ticker, name
}

func (obj *Model) AddAllyInfoEntry(allyInfo *DBallyInfo) DBresult {
	whereClause := fmt.Sprintf(`alliance_id="%d"`, allyInfo.AllianceID)
	num := obj.getNumEntries("ally_info", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "ally_info" (
			alliance_id,
			creator_corporation_id,
			creator_id,
			date_founded,
			executor_corporation_id,
			faction_id,
			name,
			ticker)
			VALUES(?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			allyInfo.AllianceID,
			allyInfo.CreatorCorporationID,
			allyInfo.CreatorID,
			allyInfo.DateFounded,
			allyInfo.ExecutorCorpID,
			allyInfo.FactionID,
			allyInfo.NameStrRef,
			allyInfo.TickerStrRef)
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

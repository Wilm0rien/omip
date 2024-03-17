package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBMiningObserver struct {
	LastUpdated  int64
	ObserverID   int64
	ObserverType int64
	OwnerCorpID  int
}

type DBMiningData struct {
	LastUpdated           int64
	CharacterID           int
	RecordedCorporationID int
	TypeID                int
	Quantity              int
	ObserverID            int64
	OwnerCorpID           int
}

type ViewMiningData struct {
	MainID   int
	MainName string
	AltName  string
	Ticker   string
	DBMiningData
}

func (obj *Model) createMiningObserverTable() {
	if !obj.checkTableExists("mining_observers") {
		_, err2 := obj.DB.Exec(`
		CREATE TABLE "mining_observers" (
			"LastUpdated" INT,
			"ObserverID" INT,
			"ObserverType" INT,
			"OwnerCorpID" INT);`)
		util.CheckErr(err2)
	} else {
		// compatiblity chekc
		queryStr := fmt.Sprint(`SELECT * FROM mining_observers WHERE OwnerCorpID IS NOT NULL LIMIT 1;`)
		_, err := obj.DB.Prepare(queryStr)
		if err == nil {
			return
		}
		_, err3 := obj.DB.Exec(`DROP TABLE IF EXISTS mining_observers;`)
		util.CheckErr(err3)

		_, err2 := obj.DB.Exec(`
		CREATE TABLE "mining_observers" (
			"LastUpdated" INT,
			"ObserverID" INT,
			"ObserverType" INT,
			"OwnerCorpID" INT);`)
		util.CheckErr(err2)
	}
}
func (obj *Model) createMiningDataTable() {
	if !obj.checkTableExists("mining_data") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "mining_data" (
		    "ObserverID" INT,
			"CharID" INT,
			"LastUpdated" INT,
			"Quantity" INT,
			"RecordedCorpID" INT,
			"TypeID" INT,
		    "OwnerCorpID" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddMiningObsEntry(item *DBMiningObserver) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`ObserverID="%d"`, item.ObserverID)
	num := obj.getNumEntries("mining_observers", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "mining_observers" (
				LastUpdated,
				ObserverID,
				ObserverType,
				OwnerCorpID)
				VALUES(?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			item.LastUpdated,
			item.ObserverID,
			item.ObserverType,
			item.OwnerCorpID)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		// update service entry!
		stmt, err := obj.DB.Prepare(`
				UPDATE "mining_observers" SET 
				LastUpdated=?,
				OwnerCorpID=?
				WHERE ObserverID=?;`)
		util.CheckErr(err)
		defer stmt.Close()
		_, err = stmt.Exec(item.LastUpdated, item.OwnerCorpID, item.ObserverID)
		util.CheckErr(err)
		retval = DBR_Updated
	}
	return retval
}

func (obj *Model) AddMiningDataEntry(item *DBMiningData) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`ObserverID=%d AND LastUpdated=%d AND CharID=%d AND TypeID=%d`,
		item.ObserverID, item.LastUpdated, item.CharacterID, item.TypeID)
	num := obj.getNumEntries("mining_data", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "mining_data" (
			    ObserverID,
				CharID,
			    LastUpdated,
			    Quantity,
			    RecordedCorpID,
			    TypeID,
			    OwnerCorpID)
				VALUES(?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			item.ObserverID,
			item.CharacterID,
			item.LastUpdated,
			item.Quantity,
			item.RecordedCorporationID,
			item.TypeID,
			item.OwnerCorpID)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		// update service entry!
		stmt, err := obj.DB.Prepare(`
				UPDATE "mining_data" SET 
				Quantity=?,
				RecordedCorpID=?,
				OwnerCorpID=?
				WHERE ObserverID=? AND LastUpdated=? AND CharID=? AND TypeID=?;`)
		util.CheckErr(err)
		defer stmt.Close()
		_, err = stmt.Exec(item.Quantity, item.RecordedCorporationID, item.OwnerCorpID, item.ObserverID, item.LastUpdated, item.CharacterID, item.TypeID)
		util.CheckErr(err)
		retval = DBR_Updated
	}
	return retval
}

func (obj *Model) GetCorpMiningData(corpID int) (list []*ViewMiningData) {
	list = make([]*ViewMiningData, 0, 1000)
	queryStr := fmt.Sprint(`
SELECT ObserverID, CharID,corp_members.main_id as MainID,LastUpdated, Quantity, 
       RecordedCorpID, TypeID, OwnerCorpID,
	   string_table.string as AltName,
	   stringMain.string as MainName         
		FROM mining_data 
		Inner JOIN 
		   corp_members ON corp_members.character_id = CharID
		INNER JOIN		   
			(SELECT character_id, name FROM corp_members) corpRef2Main
			ON corpRef2Main.character_id = corp_members.main_id
		INNER JOIN		   
			string_table ON corp_members.name= string_table.string_hash
		INNER JOIN		   
			(SELECT string_hash, string FROM string_table) stringMain
			ON corpRef2Main.name= stringMain.string_hash							
								WHERE OwnerCorpID=?
                                ORDER BY LastUpdated DESC;
`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpID)
	util.CheckErr(err)
	defer rows.Close()
	tickerMap := make(map[int]string)
	for rows.Next() {
		var mininItem ViewMiningData
		rows.Scan(
			&mininItem.ObserverID,
			&mininItem.CharacterID,
			&mininItem.MainID,
			&mininItem.LastUpdated,
			&mininItem.Quantity,
			&mininItem.RecordedCorporationID,
			&mininItem.TypeID,
			&mininItem.OwnerCorpID,
			&mininItem.AltName,
			&mininItem.MainName,
		)
		if tickerStr, ok := tickerMap[mininItem.RecordedCorporationID]; !ok {
			mininItem.Ticker = obj.GetCorpTicker(mininItem.RecordedCorporationID)
			tickerMap[mininItem.RecordedCorporationID] = mininItem.Ticker
		} else {
			mininItem.Ticker = tickerStr
		}

		list = append(list, &mininItem)
	}
	return
}
func (obj *Model) GetExtMiningData(corpID int, obsId int) (list []*ViewMiningData) {
	list = make([]*ViewMiningData, 0, 1000)
	queryStr := fmt.Sprint(`
	SELECT ObserverID, CharID,LastUpdated, Quantity, 
		   RecordedCorpID, TypeID, OwnerCorpID
			FROM mining_data 
			WHERE RecordedCorpID!=? AND ObserverID=?
			ORDER BY LastUpdated DESC;
	`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpID, obsId)
	util.CheckErr(err)
	defer rows.Close()
	tickerMap := make(map[int]string)
	for rows.Next() {
		var mininItem ViewMiningData
		rows.Scan(
			&mininItem.ObserverID,
			&mininItem.CharacterID,
			&mininItem.LastUpdated,
			&mininItem.Quantity,
			&mininItem.RecordedCorporationID,
			&mininItem.TypeID,
			&mininItem.OwnerCorpID,
		)
		mininItem.MainID = mininItem.CharacterID
		mininItem.MainName = obj.GetNameByID(mininItem.CharacterID)
		mininItem.AltName = mininItem.MainName

		if tickerStr, ok := tickerMap[mininItem.RecordedCorporationID]; !ok {
			mininItem.Ticker = obj.GetCorpTicker(mininItem.RecordedCorporationID)
			tickerMap[mininItem.RecordedCorporationID] = mininItem.Ticker
		} else {
			mininItem.Ticker = tickerStr
		}

		list = append(list, &mininItem)
	}
	return
}

func (obj *Model) GetMiningFilteredExt(mainId int, startTS int64, endTS int64) (list []*ViewMiningData) {
	list = make([]*ViewMiningData, 0, 1000)
	queryStr := fmt.Sprint(`
	SELECT ObserverID, CharID,LastUpdated, Quantity, 
		   RecordedCorpID, TypeID, OwnerCorpID
			FROM mining_data 
			WHERE CharID=? AND LastUpdated>=? and LastUpdated<?
			ORDER BY LastUpdated DESC;
	`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(mainId, startTS, endTS)
	util.CheckErr(err)
	defer rows.Close()
	tickerMap := make(map[int]string)
	for rows.Next() {
		var mininItem ViewMiningData
		rows.Scan(
			&mininItem.ObserverID,
			&mininItem.CharacterID,
			&mininItem.LastUpdated,
			&mininItem.Quantity,
			&mininItem.RecordedCorporationID,
			&mininItem.TypeID,
			&mininItem.OwnerCorpID,
		)
		mininItem.MainID = mininItem.CharacterID
		mininItem.MainName = obj.GetNameByID(mininItem.CharacterID)
		mininItem.AltName = mininItem.MainName
		if tickerStr, ok := tickerMap[mininItem.RecordedCorporationID]; !ok {
			mininItem.Ticker = obj.GetCorpTicker(mininItem.RecordedCorporationID)
			tickerMap[mininItem.RecordedCorporationID] = mininItem.Ticker
		} else {
			mininItem.Ticker = tickerStr
		}
		list = append(list, &mininItem)
	}
	return
}

func (obj *Model) GetMiningFiltered(corpId int, mainId int, startTS int64, endTS int64) (list []*ViewMiningData) {
	list = make([]*ViewMiningData, 0, 1000)
	queryStr := fmt.Sprint(`
SELECT ObserverID, CharID,corp_members.main_id as MainID,LastUpdated, Quantity, 
       RecordedCorpID, TypeID, OwnerCorpID,
	   string_table.string as AltName,
	   stringMain.string as MainName         
		FROM mining_data 
		Inner JOIN 
		   corp_members ON corp_members.character_id = CharID
		INNER JOIN		   
			(SELECT character_id, name FROM corp_members) corpRef2Main
			ON corpRef2Main.character_id = corp_members.main_id
		INNER JOIN		   
			string_table ON corp_members.name= string_table.string_hash
		INNER JOIN		   
			(SELECT string_hash, string FROM string_table) stringMain
			ON corpRef2Main.name= stringMain.string_hash							
		WHERE OwnerCorpID=? AND corpRef2Main.character_id=? AND LastUpdated>=? and LastUpdated<?
        ORDER BY LastUpdated DESC;
`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId, mainId, startTS, endTS)
	util.CheckErr(err)
	defer rows.Close()
	tickerMap := make(map[int]string)
	for rows.Next() {
		var mininItem ViewMiningData
		rows.Scan(
			&mininItem.ObserverID,
			&mininItem.CharacterID,
			&mininItem.MainID,
			&mininItem.LastUpdated,
			&mininItem.Quantity,
			&mininItem.RecordedCorporationID,
			&mininItem.TypeID,
			&mininItem.OwnerCorpID,
			&mininItem.AltName,
			&mininItem.MainName,
		)
		if tickerStr, ok := tickerMap[mininItem.RecordedCorporationID]; !ok {
			mininItem.Ticker = obj.GetCorpTicker(mininItem.RecordedCorporationID)
			tickerMap[mininItem.RecordedCorporationID] = mininItem.Ticker
		} else {
			mininItem.Ticker = tickerStr
		}
		list = append(list, &mininItem)
	}
	return
}

func (obj *Model) GetMiningByCop(corpId int, startTS int64, endTS int64) (list []*ViewMiningData) {
	list = make([]*ViewMiningData, 0, 1000)
	queryStr := fmt.Sprint(`
	SELECT ObserverID, CharID,LastUpdated, Quantity, 
		   RecordedCorpID, TypeID, OwnerCorpID
			FROM mining_data 
			WHERE RecordedCorpID=? AND LastUpdated>=? and LastUpdated<?
			ORDER BY LastUpdated DESC;
	`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId, startTS, endTS)
	util.CheckErr(err)
	defer rows.Close()
	tickerMap := make(map[int]string)
	for rows.Next() {
		var mininItem ViewMiningData
		rows.Scan(
			&mininItem.ObserverID,
			&mininItem.CharacterID,
			&mininItem.LastUpdated,
			&mininItem.Quantity,
			&mininItem.RecordedCorporationID,
			&mininItem.TypeID,
			&mininItem.OwnerCorpID,
		)
		mininItem.MainID = mininItem.CharacterID
		mininItem.MainName = obj.GetNameByID(mininItem.CharacterID)
		mininItem.AltName = mininItem.MainName
		if tickerStr, ok := tickerMap[mininItem.RecordedCorporationID]; !ok {
			mininItem.Ticker = obj.GetCorpTicker(mininItem.RecordedCorporationID)
			tickerMap[mininItem.RecordedCorporationID] = mininItem.Ticker
		} else {
			mininItem.Ticker = tickerStr
		}
		list = append(list, &mininItem)
	}
	return
}

func (obj *Model) GetCorpObservers(corpId int) (result []int) {
	result = make([]int, 0, 10)
	queryStr := fmt.Sprint(`
		SELECT ObserverID FROM mining_observers WHERE OwnerCorpID=?
		`)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(corpId)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var obsId int
		rows.Scan(&obsId)
		result = append(result, obsId)
	}
	return

}

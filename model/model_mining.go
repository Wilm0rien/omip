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
	DBMiningData
}

func (obj *Model) createMiningObserverTable() {
	if !obj.checkTableExists("mining_observers") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "mining_observers" (
			"LastUpdated" INT,
			"ObserverID" INT,
			"ObserverType" INT,
			"OwnerCorpID" INT);`)
		util.CheckErr(err)
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

func (obj *Model) GetMiningData(corpID int) (list []*ViewMiningData) {
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
		list = append(list, &mininItem)
	}
	return
}

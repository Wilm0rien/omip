package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBMiningObserver struct {
	LastUpdated  int64
	ObserverID   int64
	ObserverType int64
}

type DBMiningData struct {
	LastUpdated           int64
	CharacterID           int
	RecordedCorporationID int
	TypeID                int
	Quantity              int
	ObserverID            int64
}

func (obj *Model) createMiningObserverTable() {
	if !obj.checkTableExists("mining_observers") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "mining_observers" (
			"LastUpdated" INT,
			"ObserverID" INT,
			"ObserverType" INT);`)
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
			"TypeID" INT);`)
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
				ObserverType)
				VALUES(?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			item.LastUpdated,
			item.ObserverID,
			item.ObserverType)
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
				LastUpdated=?
				WHERE ObserverID=?;`)
		util.CheckErr(err)
		defer stmt.Close()
		_, err = stmt.Exec(item.LastUpdated, item.ObserverID)
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
			    TypeID)
				VALUES(?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			item.ObserverID,
			item.CharacterID,
			item.LastUpdated,
			item.Quantity,
			item.RecordedCorporationID,
			item.TypeID)
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
				RecordedCorpID=?
				WHERE ObserverID=? AND LastUpdated=? AND CharID=? AND TypeID=?;`)
		util.CheckErr(err)
		defer stmt.Close()
		_, err = stmt.Exec(item.Quantity, item.RecordedCorporationID, item.ObserverID, item.LastUpdated, item.CharacterID, item.TypeID)
		util.CheckErr(err)
		retval = DBR_Updated
	}
	return retval
}

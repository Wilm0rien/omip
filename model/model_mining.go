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

func (obj *Model) createMiningObserverTable() {
	if !obj.checkTableExists("mining_observers") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "journal_links" (
			"LastUpdated" INT,
			"ObserverID" INT,
			"ObserverType" INT);`)
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

func (obj *Model) GetLastObsUpdateTime(ObserverID int64) (result int64, found bool) {
	queryString := fmt.Sprintf("SELECT LastUpdated FROM mining_observers WHERE ObserverID=%d;", ObserverID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&result)
		found = true
	}
	return
}

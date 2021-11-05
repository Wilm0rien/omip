package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBTransaction struct {
	CharID        int
	CorpID        int
	Division      int
	ClientID      int
	Date          int64
	IsBuy         int
	JournalRefID  int64
	LocationID    int64
	Quantity      int
	TransactionID int64
	TypeID        int
	UnitPrice     float64
}

func (obj *Model) createTransactionsTable() {
	if !obj.checkTableExists("transactions") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "transactions" (
			"charId" INT,
			"corpId" INT,
			"division" INT,
			"client_id" INT,
			"date" INT,
			"is_buy" INT,
			"journal_ref_id" INT,
			"location_id"  INT,
			"quantity" INT,
			"transaction_id" INT,
			"type_id" INT,
			"unit_price" REAL);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddTransactionEntry(trItem *DBTransaction) DBresult {
	whereClause := fmt.Sprintf(`transaction_id="%d"`, trItem.TransactionID)
	num := obj.getNumEntries("transactions", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "transactions" (
			charId,
			corpId,
			division,
			client_id,
			date,
			is_buy,
			journal_ref_id,
			location_id,
			quantity,
			transaction_id,
			type_id,
			unit_price)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			trItem.CharID,
			trItem.CorpID,
			trItem.Division,
			trItem.ClientID,
			trItem.Date,
			trItem.IsBuy,
			trItem.JournalRefID,
			trItem.LocationID,
			trItem.Quantity,
			trItem.TransactionID,
			trItem.TypeID,
			trItem.UnitPrice)
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

func (obj *Model) GetTransactionEntry(JournalRefID int64) *DBTransaction {
	var retval *DBTransaction
	queryString := fmt.Sprintf(`SELECT
			charId,
			corpId,
			division,
			client_id,
			date,
			is_buy,
			journal_ref_id,
			location_id,
			quantity,
			transaction_id,
			type_id,
			unit_price
		FROM transactions
		WHERE journal_ref_id=%d ORDER BY date ASC;`, JournalRefID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var TAitem DBTransaction
		rows.Scan(
			&TAitem.CharID,
			&TAitem.CorpID,
			&TAitem.Division,
			&TAitem.ClientID,
			&TAitem.Date,
			&TAitem.IsBuy,
			&TAitem.JournalRefID,
			&TAitem.LocationID,
			&TAitem.Quantity,
			&TAitem.TransactionID,
			&TAitem.TypeID,
			&TAitem.UnitPrice)
		retval = &TAitem
		break
	}
	return retval
}

package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"time"
)

func (obj *Model) createWalletHistoryTable() {
	if !obj.checkTableExists("wallet_history") {
		_, err := obj.DB.Exec(
			`CREATE TABLE "wallet_history" (
								"character_id" INT,
								"corporation_id" INT,  
								"division" INT,
								"timestamp" INT, 
								"balance" REAL);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddWalletEntry(character_id int, corporation_id int, division int, balance float64) {
	currentTime := time.Now().Unix()
	stmt, err := obj.DB.Prepare(
		`INSERT INTO wallet_history (
					character_id, 
					corporation_id,
					division,
					timestamp, 
					balance) 
				VALUES (?,?,?,?,?)`)
	util.CheckErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(character_id, corporation_id, division, currentTime, balance)
	util.CheckErr(err)
}

func (obj *Model) WalletEntryExists(character_id int, corporation_id int, division int) bool {
	var retval bool
	currentTime := time.Now().Unix()
	yesterDayTime := currentTime - 60*60*24
	whereClause := fmt.Sprintf(`
					character_id="%d" AND 
					corporation_id="%d" AND 
					division="%d" AND
					timestamp>"%d"`, character_id, corporation_id, division, yesterDayTime)
	num := obj.getNumEntries("wallet_history", whereClause)
	if num > 0 {
		retval = true
	}
	return retval
}

func (obj *Model) GetLatestWallets(character_id int, corporation_id int, division int) (balance float64) {
	queryString := fmt.Sprintf(`
			SELECT balance FROM wallet_history 
			WHERE character_id=%d AND  corporation_id="%d" AND division="%d"
			ORDER BY timestamp DESC;`, character_id, corporation_id, division)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&balance)
		break
	}
	return balance
}

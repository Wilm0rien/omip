package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

func (obj *Model) createStringTable() {
	if !obj.checkTableExists("string_table") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "string_table" (
			"string_hash" INT PRIMARY KEY UNIQUE,
			"string" TEXT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) GetStringEntry(id int64) (retStr string, result bool) {
	whereClause := fmt.Sprintf(`string_hash="%d"`, id)
	num := obj.getNumEntries("string_table", whereClause)
	util.Assert(num < 2)
	if num == 1 {
		queryString := fmt.Sprintf(`SELECT string from %s WHERE string_hash="%d";`, "string_table", id)
		rows, err := obj.DB.Query(queryString)
		util.CheckErr(err)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&retStr)
			result = true
		}
		defer rows.Close()
	}
	return retStr, result
}

func (obj *Model) AddStringEntry(entryString string) (stringhash int64) {
	stringhash = util.Get64BitMd5FromString(entryString)
	whereClause := fmt.Sprintf(`string_hash="%d"`, stringhash)
	num := obj.getNumEntries("string_table", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare("INSERT INTO string_table (string_hash, string) values (?,?)")
		util.CheckErr(err)
		_, err = stmt.Exec(stringhash, entryString)
		util.CheckErr(err)
		defer stmt.Close()
	}
	return stringhash
}

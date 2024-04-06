package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBUniverseName struct {
	Category int64
	ID       int
	NameRef  int64
}

func (obj *Model) createUniNamesTable() {
	if !obj.checkTableExists("universe_names") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "universe_names" (
			"CategoryStrRef" INT,
			"eve_id" INT PRIMARY KEY UNIQUE,
			"string_hash" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) NameExists(eve_id int) bool {
	whereClause := fmt.Sprintf(`eve_id="%d"`, eve_id)
	num := obj.getNumEntries("universe_names", whereClause)
	return num != 0
}

func (obj *Model) AddNameEntry(UNitem *DBUniverseName) (retval DBresult) {
	if !obj.NameExists(UNitem.ID) {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "universe_names" (
				string_hash,
				eve_id,
				CategoryStrRef)
				VALUES (?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			UNitem.NameRef,
			UNitem.ID,
			UNitem.Category)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) GetNameEntry(eve_id int) (UNitem *DBUniverseName) {
	queryString := fmt.Sprintf(`
		SELECT CategoryStrRef, eve_id, string_hash
		FROM universe_names
		WHERE eve_id=%d;`, eve_id)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var resultList []*DBUniverseName
	for rows.Next() {
		var elem DBUniverseName
		rows.Scan(&elem.Category, &elem.ID, &elem.NameRef)
		resultList = append(resultList, &elem)
	}
	if len(resultList) == 1 {
		UNitem = resultList[0]
	}
	return UNitem
}

func (obj *Model) GetNameByID(eve_id int) (name string) {
	UNitem := obj.GetNameEntry(eve_id)
	var ok bool
	if UNitem != nil {
		name, ok = obj.GetStringEntry(UNitem.NameRef)
		util.Assert(ok)
	} else {
		if member := obj.GetDBCorpMember(eve_id); member != nil {
			if nstr, ok2 := obj.GetStringEntry(member.NameRef); ok2 {
				name = nstr
			}
		}
		if name == "" {
			name = fmt.Sprintf("%d", eve_id)
		}
	}
	return name
}

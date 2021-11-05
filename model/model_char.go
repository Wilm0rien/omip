package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBcharacter struct {
	Character_id    int
	Alliance_id     int
	Birthday        int64
	Corporation_id  int
	Security_status float64
	Name            string
	Director        bool
}

func (obj *Model) createCharTable() {
	if !obj.checkTableExists("characters") {
		_, err := obj.DB.Exec("CREATE TABLE `characters` (`character_id` INTEGER PRIMARY KEY, `alliance_id` INT, `birthday` INT, `corporation_id` INT, `security_status` REAL, `name` TEXT, `director` INT)")
		util.CheckErr(err)
	}
}

func (obj *Model) AddCharEntry(entry *DBcharacter) {
	var director int
	if entry.Director {
		director = 1
	}
	stmt, err := obj.DB.Prepare("INSERT INTO characters (character_id, alliance_id, corporation_id, name, director) values (?,?,?,?,?)")
	util.CheckErr(err)
	_, err = stmt.Exec(entry.Character_id, entry.Alliance_id, entry.Corporation_id, entry.Name, director)
	util.CheckErr(err)
	defer stmt.Close()
}

func (obj *Model) CheckCharExists(entry *DBcharacter) bool {
	whereClause := fmt.Sprintf(`character_id="%d"`, entry.Character_id)
	num := obj.getNumEntries("characters", whereClause)
	var retval bool
	if num > 0 {
		retval = true
	}
	return retval
}

func (obj *Model) GetCharEntry(Character_id int) *DBcharacter {
	queryString := fmt.Sprintf(`SELECT alliance_id, corporation_id, director, name FROM characters WHERE character_id="%d";`, Character_id)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var retval DBcharacter
	var director int
	for rows.Next() {
		rows.Scan(&retval.Alliance_id, &retval.Corporation_id, &director, &retval.Name)
	}
	if director == 1 {
		retval.Director = true
	}
	retval.Character_id = Character_id
	return &retval
}

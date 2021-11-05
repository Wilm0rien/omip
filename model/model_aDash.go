package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"time"
)

type DBpap struct {
	CorpId     int
	PapLink    string
	ChName     int64
	CoTicker   int64
	AlShort    int64
	ShTypeName int64
	Sub        int64
	Loc        int64
	Timestamp  int64
}

type DBpapTable struct {
	AltName  string
	PapLink  string
	Time     int64
	MainName string
	MainID   int
}

func (obj *Model) createADashAuthTable() {
	if !obj.checkTableExists("adash_accounts") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "adash_accounts" (
			"corporation_id",
			"email" TEXT,
			"password" TEXT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createADashTable() {
	if !obj.checkTableExists("pap_status") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "pap_status" (
			"corporation_id" INT,
			"paplink" TEXT,
			"chName" INT,
			"coTicker" INT,
			"alShort" INT,
			"shTypeName" INT,
			"sub" INT,
			"loc" INT,
			"timestamp" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) ADashAuthExists(corpId int) bool {
	whereClause := fmt.Sprintf(`corporation_id=%d`, corpId)
	num := obj.getNumEntries("adash_accounts", whereClause)
	return (num != 0)
}

func (obj *Model) GetAuth(corpId int) (email string, pw string, success bool) {
	queryString := fmt.Sprintf(`SELECT email, password FROM adash_accounts WHERE corporation_id=%d;`, corpId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&email, &pw)
	}
	passPhrase := fmt.Sprintf(".#_%d%s", corpId, util.GenSysPassphrase())
	dencryptPW, success1 := util.Decrypt([]byte(pw), passPhrase)
	return email, string(dencryptPW), success1
}

func (obj *Model) SetAuth(corpId int, email string, pw string) {
	retval := DBR_Undefined
	passPhrase := fmt.Sprintf(".#_%d%s", corpId, util.GenSysPassphrase())
	encryptPW := util.Encrypt([]byte(pw), passPhrase)
	whereClause := fmt.Sprintf(`corporation_id=%d`, corpId)
	num := obj.getNumEntries("adash_accounts", whereClause)
	if num != 0 {
		queryStr := fmt.Sprintf("DELETE FROM adash_accounts WHERE corporation_id=%d;", corpId)
		stmt, err := obj.DB.Prepare(queryStr)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec()
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Removed
		}
		util.Assert(retval == DBR_Removed)
	}
	stmt, err := obj.DB.Prepare(`
		INSERT INTO "adash_accounts" (
			corporation_id,
			email,
			password)
			VALUES (?,?,?);`)
	util.CheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(corpId, email, string(encryptPW))
	affect, err := res.RowsAffected()
	util.CheckErr(err)
	if affect > 0 {
		retval = DBR_Inserted
	}
	util.Assert(retval == DBR_Inserted)
}

func (obj *Model) PapLinkExists(papLink string) bool {
	whereClause := fmt.Sprintf(`paplink="%s"`, papLink)
	num := obj.getNumEntries("pap_status", whereClause)
	return (num != 0)
}

func (obj *Model) AddCorpADashEntry(pap *DBpap) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`paplink="%s" and chName=%d`, pap.PapLink, pap.ChName)
	num := obj.getNumEntries("pap_status", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "pap_status" (
				corporation_id,
				paplink,
				chName,
				coTicker,
				alShort,
				shTypeName,
				sub,
				loc,
				timestamp)
				VALUES (?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			pap.CorpId,
			pap.PapLink,
			pap.ChName,
			pap.CoTicker,
			pap.AlShort,
			pap.ShTypeName,
			pap.Sub,
			pap.Loc,
			pap.Timestamp)
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

func (obj *Model) GetPapTable(corpId int) *MonthlyTable {
	paps := obj.GetDBPapTable(corpId)
	return obj.GetMonthlyTable(corpId, paps, 12)
}

func (obj *Model) GetDBPapTable(corpId int) []*DBTable {
	var retval []*DBTable
	queryString := fmt.Sprintf(`
		SELECT 
			string_table.string AS AltName, 
			pap_status.timestamp AS Time,
			string2.string AS MainName
		FROM pap_status
		INNER JOIN 
			string_table ON string_table.string_hash = pap_status.chName
		INNER JOIN 
			corp_members ON string_table.string_hash = corp_members.name
		JOIN 
		   (SELECT character_id, main_id, name FROM corp_members) corp2
		   ON corp_members.main_id = corp2.character_id
		JOIN 
			(SELECT string_hash, string FROM string_table) string2
			ON corp2.name=string2.string_hash
		WHERE pap_status.corporation_id=%d
		ORDER BY Time DESC;`, corpId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var elem DBTable
		rows.Scan(&elem.AltName, &elem.Time, &elem.MainName)
		elem.Amount = 1
		retval = append(retval, &elem)
	}
	return retval
}

func (obj *Model) GetCurrentPaps(corpId int) int {
	papTable := obj.GetDBPapTable(corpId)
	var participation int
	now_year, now_month, _ := time.Now().Date()
	for _, elem := range papTable {
		tm := time.Unix(elem.Time, 0)
		year, month, _ := tm.Date()
		if year == now_year && month == now_month {
			participation++
		}
	}
	return participation
}

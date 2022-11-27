package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

const (
	Ctnr_Avail_Public   = 0
	Ctnr_Avail_Personal = 1
	Ctnr_Avail_Corp     = 2
	Ctnr_Avail_Alli     = 3
)

var contractsAvailability = map[string]int{
	"public":      Ctnr_Avail_Public,
	"personal":    Ctnr_Avail_Personal,
	"corporation": Ctnr_Avail_Corp,
	"alliance":    Ctnr_Avail_Alli,
}

const (
	Cntr_Stat_outstanding         = 0
	Cntr_Stat_in_progress         = 1
	Cntr_Stat_finished_issuer     = 2
	Cntr_Stat_finished_contractor = 3
	Cntr_Stat_finished            = 4
	Cntr_Stat_cancelled           = 5
	Cntr_Stat_rejected            = 6
	Cntr_Stat_failed              = 7
	Cntr_Stat_deleted             = 8
	Cntr_Stat_reversed            = 9
)

var contractsStatus = map[string]int{
	"outstanding":         Cntr_Stat_outstanding,
	"in_progress":         Cntr_Stat_in_progress,
	"finished_issuer":     Cntr_Stat_finished_issuer,
	"finished_contractor": Cntr_Stat_finished_contractor,
	"finished":            Cntr_Stat_finished,
	"cancelled":           Cntr_Stat_cancelled,
	"rejected":            Cntr_Stat_rejected,
	"failed":              Cntr_Stat_failed,
	"deleted":             Cntr_Stat_deleted,
	"reversed":            Cntr_Stat_reversed,
}

const (
	Cntr_Type_unknown       = 0
	Cntr_Type_item_exchange = 1
	Cntr_Type_auction       = 2
	Cntr_Type_courier       = 3
	Cntr_Type_loan          = 4
)

var contractsType = map[string]int{
	"unknown":       Cntr_Type_unknown,
	"item_exchange": Cntr_Type_item_exchange,
	"auction":       Cntr_Type_auction,
	"courier":       Cntr_Type_courier,
	"loan":          Cntr_Type_loan,
}

type DBContract struct {
	Acceptor_id           int
	Assignee_id           int
	Availability          int
	Buyout                float64
	Collateral            float64
	Contract_id           int
	Date_accepted         int64
	Date_completed        int64
	Date_expired          int64
	Date_issued           int64
	Days_to_complete      int
	End_location_id       int64
	For_corporation       bool
	Issuer_corporation_id int
	Issuer_id             int
	Price                 float64
	Reward                float64
	Start_location_id     int64
	Status                int
	Title                 int64
	Type                  int
	Volume                float64
}

type DBContrItem struct {
	Contract_id  int
	Is_included  int
	Is_singleton int
	Quantity     int
	Raw_quantity int
	Record_id    int64
	Type_id      int
}

type DBcontractTable struct {
	Character_Id          int
	Contract_id           int
	Status                int
	CharName              string
	Date_issued           int64
	Date_completed        int64
	Price                 float64
	Title                 string
	Issuer_corporation_id int
	For_corporation       int
	Items                 []*DBContrItem
}

type contract_stats struct {
	avg_price float64
	min_price float64
	max_price float64
	sum_price float64
	count     int
}

func (obj *Model) createContrItemTable() {
	if !obj.checkTableExists("contr_items") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "contr_items" (
			contract_id INT,
			is_included INT,
			is_singleton INT,
			quantity INT,
			raw_quantity INT,
			record_id INT,
			type_id INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createContractTable() {
	if !obj.checkTableExists("contracts") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "contracts" (
			"acceptor_id" INT,
			"assignee_id" INT,
			"availability" INT,
			"buyout" REAL,
			"collateral" REAL,
			"contract_id" INT,
			"date_accepted" INT,
			"date_completed" INT,
			"date_expired" INT,
			"date_issued" INT,
			"days_to_complete" INT,
			"end_location_id" INT,
			"for_corporation" INT,
			"issuer_corporation_id" INT,
			"issuer_id" INT,
			"price" REAL,
			"reward" REAL,
			"start_location_id" INT,
			"status" INT,
			"title" INT,
			"type" INT,
			"volume" REAL
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) ContrItemsExist(contrID int) bool {
	whereClause := fmt.Sprintf(`contract_id="%d"`, contrID)
	num := obj.getNumEntries("contr_items", whereClause)
	return num != 0
}

func (obj *Model) AddContrItemEntry(contrItem *DBContrItem) DBresult {
	whereClause := fmt.Sprintf(`record_id="%d"`, contrItem.Record_id)
	num := obj.getNumEntries("contr_items", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "contr_items" (
				contract_id,
				is_included,
				is_singleton,
				quantity,
				raw_quantity,
				record_id,
				type_id)
				VALUES (?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			contrItem.Contract_id,
			contrItem.Is_included,
			contrItem.Is_singleton,
			contrItem.Quantity,
			contrItem.Raw_quantity,
			contrItem.Record_id,
			contrItem.Type_id)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) AddContractEntry(contract *DBContract) DBresult {
	whereClause := fmt.Sprintf(`contract_id="%d"`, contract.Contract_id)
	num := obj.getNumEntries("contracts", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "contracts" (
			acceptor_id,
			assignee_id,
			availability,
			buyout,
			collateral,
			contract_id,
			date_accepted,
			date_completed,
			date_expired,
			date_issued,
			days_to_complete,
			end_location_id,
			for_corporation,
			issuer_corporation_id,
			issuer_id,
			price,
			reward,
			start_location_id,
			status,
			title,
			type,
			volume) 
			values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			contract.Acceptor_id,
			contract.Assignee_id,
			contract.Availability,
			contract.Buyout,
			contract.Collateral,
			contract.Contract_id,
			contract.Date_accepted,
			contract.Date_completed,
			contract.Date_expired,
			contract.Date_issued,
			contract.Days_to_complete,
			contract.End_location_id,
			contract.For_corporation,
			contract.Issuer_corporation_id,
			contract.Issuer_id,
			contract.Price,
			contract.Reward,
			contract.Start_location_id,
			contract.Status,
			contract.Title,
			contract.Type,
			contract.Volume)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		// check if status changed before update
		whereClause := fmt.Sprintf(`contract_id="%d" AND status="%d"`, contract.Contract_id, contract.Status)
		num := obj.getNumEntries("contracts", whereClause)
		if num == 0 {
			stmt, err := obj.DB.Prepare(`
				UPDATE "contracts" SET
					acceptor_id=?,
					assignee_id=?,
					availability=?,
					buyout=?,
					collateral=?,
					date_accepted=?,
					date_completed=?,
					date_expired=?,
					date_issued=?,
					days_to_complete=?,
					end_location_id=?,
					for_corporation=?,
					issuer_corporation_id=?,
					issuer_id=?,
					price=?,
					reward=?,
					start_location_id=?,
					status=?,
					title=?,
					type=?,
					volume=? 
					WHERE contract_id=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				contract.Acceptor_id,
				contract.Assignee_id,
				contract.Availability,
				contract.Buyout,
				contract.Collateral,
				contract.Date_accepted,
				contract.Date_completed,
				contract.Date_expired,
				contract.Date_issued,
				contract.Days_to_complete,
				contract.End_location_id,
				contract.For_corporation,
				contract.Issuer_corporation_id,
				contract.Issuer_id,
				contract.Price,
				contract.Reward,
				contract.Start_location_id,
				contract.Status,
				contract.Title,
				contract.Type,
				contract.Volume,
				contract.Contract_id)
			util.CheckErr(err)
			affect, err := res.RowsAffected()
			util.CheckErr(err)
			if affect > 0 {
				retval = DBR_Updated
			}
		} else {
			retval = DBR_Skipped
		}
	}
	return retval
}

func (obj *Model) GetContrItems(contrID int) []*DBContrItem {
	queryString := fmt.Sprintf(`SELECT 
								contract_id,
								is_included,
								is_singleton,
								quantity,
								raw_quantity,
								record_id,
								type_id
								FROM contr_items
								WHERE contract_id=%d`, contrID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var retval []*DBContrItem
	for rows.Next() {
		var elem DBContrItem
		rows.Scan(&elem.Contract_id, &elem.Is_included, &elem.Is_singleton, &elem.Quantity, &elem.Raw_quantity, &elem.Record_id, &elem.Type_id)
		retval = append(retval, &elem)
	}
	return retval
}

func (obj *Model) GetCorpContracts(id int) []*DBContract {
	queryString := fmt.Sprintf(`
	SELECT 
		acceptor_id,
		assignee_id,
		availability,
		buyout,
		collateral,
		contract_id,
		date_accepted,
		date_completed,
		date_expired,
		date_issued,
		days_to_complete,
		end_location_id,
		for_corporation,
		issuer_corporation_id,
		issuer_id,
		price,
		reward,
		start_location_id,
		status,
		title,
		type,
		volume,
		journal.amount as jourPrice
	FROM contracts
	LEFT JOIN
		journal_links ON contracts.contract_id = journal_links.contractID
	LEFT JOIN
		journal ON journal_links.journalID = journal.id 
	WHERE 
		(issuer_corporation_id=%d AND for_corporation=1 and status!=%d) 
		or (contracts.contract_id = journal_links.contractID and journal.corpId = %d)
	GROUP BY contract_id;`, id, Cntr_Stat_deleted, id)

	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	retval := make([]*DBContract, 0, 10)
	for rows.Next() {
		var contract DBContract
		var jourPrice float64
		rows.Scan(&contract.Acceptor_id,
			&contract.Assignee_id,
			&contract.Availability,
			&contract.Buyout,
			&contract.Collateral,
			&contract.Contract_id,
			&contract.Date_accepted,
			&contract.Date_completed,
			&contract.Date_expired,
			&contract.Date_issued,
			&contract.Days_to_complete,
			&contract.End_location_id,
			&contract.For_corporation,
			&contract.Issuer_corporation_id,
			&contract.Issuer_id,
			&contract.Price,
			&contract.Reward,
			&contract.Start_location_id,
			&contract.Status,
			&contract.Title,
			&contract.Type,
			&contract.Volume,
			&jourPrice)

		if contract.Issuer_corporation_id != id {
			// if this is an accepted contract show the value from the journal entry rather than the contract

			contract.Price = jourPrice
		}
		retval = append(retval, &contract)

	}
	return retval
}

func (obj *Model) GetContractsByIssuerId(id int, corp bool) []*DBContract {
	var where string
	if corp {
		where = fmt.Sprintf("issuer_corporation_id=%d AND for_corporation=1", id)
	} else {
		where = fmt.Sprintf("(issuer_id=%d OR acceptor_id=%d) AND for_corporation=0", id, id)
	}
	queryString := fmt.Sprintf(`SELECT 
								acceptor_id,
								assignee_id,
								availability,
								buyout,
								collateral,
								contract_id,
								date_accepted,
								date_completed,
								date_expired,
								date_issued,
								days_to_complete,
								end_location_id,
								for_corporation,
								issuer_corporation_id,
								issuer_id,
								price,
								reward,
								start_location_id,
								status,
								title,
								type,
								volume
								FROM contracts
								WHERE %s and status!=%d
								ORDER BY date_issued DESC;`, where, Cntr_Stat_deleted)

	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	retval := make([]*DBContract, 0, 10)
	for rows.Next() {
		var contract DBContract
		rows.Scan(&contract.Acceptor_id,
			&contract.Assignee_id,
			&contract.Availability,
			&contract.Buyout,
			&contract.Collateral,
			&contract.Contract_id,
			&contract.Date_accepted,
			&contract.Date_completed,
			&contract.Date_expired,
			&contract.Date_issued,
			&contract.Days_to_complete,
			&contract.End_location_id,
			&contract.For_corporation,
			&contract.Issuer_corporation_id,
			&contract.Issuer_id,
			&contract.Price,
			&contract.Reward,
			&contract.Start_location_id,
			&contract.Status,
			&contract.Title,
			&contract.Type,
			&contract.Volume)

		if !corp {
			if contract.Issuer_id != id {
				contract.Price = -contract.Price
			}
		}

		retval = append(retval, &contract)

	}
	return retval
}
func (obj *Model) GetContractById(ContractID int64) (result *DBContract) {
	queryString := fmt.Sprintf(`SELECT 
								acceptor_id,
								assignee_id,
								availability,
								buyout,
								collateral,
								contract_id,
								date_accepted,
								date_completed,
								date_expired,
								date_issued,
								days_to_complete,
								end_location_id,
								for_corporation,
								issuer_corporation_id,
								issuer_id,
								price,
								reward,
								start_location_id,
								status,
								title,
								type,
								volume
								FROM contracts
								WHERE contract_id=%d and status!=%d
								ORDER BY date_issued DESC;`, ContractID, Cntr_Stat_deleted)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var contract DBContract
		rows.Scan(&contract.Acceptor_id,
			&contract.Assignee_id,
			&contract.Availability,
			&contract.Buyout,
			&contract.Collateral,
			&contract.Contract_id,
			&contract.Date_accepted,
			&contract.Date_completed,
			&contract.Date_expired,
			&contract.Date_issued,
			&contract.Days_to_complete,
			&contract.End_location_id,
			&contract.For_corporation,
			&contract.Issuer_corporation_id,
			&contract.Issuer_id,
			&contract.Price,
			&contract.Reward,
			&contract.Start_location_id,
			&contract.Status,
			&contract.Title,
			&contract.Type,
			&contract.Volume)
		result = &contract
		break
	}
	return
}

func (obj *Model) ContractStatusStr2Int(status string) int {
	retVal, _ := contractsStatus[status]
	return retVal
}
func (obj *Model) ContractStatusInt2Str(status int) (retval string) {
	retval = "unkown"
	for k, v := range contractsStatus {
		if v == status {
			retval = k
			break
		}
	}
	return retval
}
func (obj *Model) ContractTypeStr2Int(cntrType string) int {
	retVal, _ := contractsType[cntrType]
	return retVal
}
func (obj *Model) ContractAvailStr2Int(avail string) int {
	retVal, _ := contractsAvailability[avail]
	return retVal
}
func (obj *Model) ContractAvailInt2Str(avail int) (retval string) {
	retval = "unkown"
	for k, v := range contractsAvailability {
		if v == avail {
			retval = k
			break
		}
	}
	return retval
}

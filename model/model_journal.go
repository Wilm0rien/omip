package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"time"
)

var JournalContextIdType = map[string]int{
	"structure_id":          1,
	"station_id":            2,
	"market_transaction_id": 3,
	"character_id":          4,
	"corporation_id":        5,
	"alliance_id":           6,
	"eve_system":            7,
	"industry_job_id":       8,
	"contract_id":           9,
	"planet_id":             10,
	"system_id":             11,
	"type_id":               12,
}

const (
	Journal_market_transaction = 2
	Journal_market_escrow      = 42
)

// list taken from https://github.com/esi/eve-glue/blob/master/eve_glue/wallet_journal_ref.py
var JournalRefType = map[string]int{
	"public":                                          0,
	"player_trading":                                  1,
	"market_transaction":                              2,
	"gm_cash_transfer":                                3,
	"mission_reward":                                  7,
	"clone_activation":                                8,
	"inheritance":                                     9,
	"player_donation":                                 10,
	"corporation_payment":                             11,
	"docking_fee":                                     12,
	"office_rental_fee":                               13,
	"factory_slot_rental_fee":                         14,
	"repair_bill":                                     15,
	"bounty":                                          16,
	"bounty_prize":                                    17,
	"insurance":                                       19,
	"mission_expiration":                              20,
	"mission_completion":                              21,
	"shares":                                          22,
	"courier_mission_escrow":                          23,
	"mission_cost":                                    24,
	"agent_miscellaneous":                             25,
	"lp_store":                                        26,
	"agent_location_services":                         27,
	"agent_donation":                                  28,
	"agent_security_services":                         29,
	"agent_mission_collateral_paid":                   30,
	"agent_mission_collateral_refunded":               31,
	"agents_preward":                                  32,
	"agent_mission_reward":                            33,
	"agent_mission_time_bonus_reward":                 34,
	"cspa":                                            35,
	"cspaofflinerefund":                               36,
	"corporation_account_withdrawal":                  37,
	"corporation_dividend_payment":                    38,
	"corporation_registration_fee":                    39,
	"corporation_logo_change_cost":                    40,
	"release_of_impounded_property":                   41,
	"market_escrow":                                   42,
	"agent_services_rendered":                         43,
	"market_fine_paid":                                44,
	"corporation_liquidation":                         45,
	"brokers_fee":                                     46,
	"corporation_bulk_payment":                        47,
	"alliance_registration_fee":                       48,
	"war_fee":                                         49,
	"alliance_maintainance_fee":                       50,
	"contraband_fine":                                 51,
	"clone_transfer":                                  52,
	"acceleration_gate_fee":                           53,
	"transaction_tax":                                 54,
	"jump_clone_installation_fee":                     55,
	"manufacturing":                                   56,
	"researching_technology":                          57,
	"researching_time_productivity":                   58,
	"researching_material_productivity":               59,
	"copying":                                         60,
	"reverse_engineering":                             62,
	"contract_auction_bid":                            63,
	"contract_auction_bid_refund":                     64,
	"contract_collateral":                             65,
	"contract_reward_refund":                          66,
	"contract_auction_sold":                           67,
	"contract_reward":                                 68,
	"contract_collateral_refund":                      69,
	"contract_collateral_payout":                      70,
	"contract_price":                                  71,
	"contract_brokers_fee":                            72,
	"contract_sales_tax":                              73,
	"contract_deposit":                                74,
	"contract_deposit_sales_tax":                      75,
	"contract_auction_bid_corp":                       77,
	"contract_collateral_deposited_corp":              78,
	"contract_price_payment_corp":                     79,
	"contract_brokers_fee_corp":                       80,
	"contract_deposit_corp":                           81,
	"contract_deposit_refund":                         82,
	"contract_reward_deposited":                       83,
	"contract_reward_deposited_corp":                  84,
	"bounty_prizes":                                   85,
	"advertisement_listing_fee":                       86,
	"medal_creation":                                  87,
	"medal_issued":                                    88,
	"dna_modification_fee":                            90,
	"sovereignity_bill":                               91,
	"bounty_prize_corporation_tax":                    92,
	"agent_mission_reward_corporation_tax":            93,
	"agent_mission_time_bonus_reward_corporation_tax": 94,
	"upkeep_adjustment_fee":                           95,
	"planetary_import_tax":                            96,
	"planetary_export_tax":                            97,
	"planetary_construction":                          98,
	"corporate_reward_payout":                         99,
	"bounty_surcharge":                                101,
	"contract_reversal":                               102,
	"corporate_reward_tax":                            103,
	"store_purchase":                                  106,
	"store_purchase_refund":                           107,
	"datacore_fee":                                    112,
	"war_fee_surrender":                               113,
	"war_ally_contract":                               114,
	"bounty_reimbursement":                            115,
	"kill_right_fee":                                  116,
	"security_processing_fee":                         117,
	"industry_job_tax":                                120,
	"infrastructure_hub_maintenance":                  122,
	"asset_safety_recovery_tax":                       123,
	"opportunity_reward":                              124,
	"project_discovery_reward":                        125,
	"project_discovery_tax":                           126,
	"reprocessing_tax":                                127,
	"jump_clone_activation_fee":                       128,
	"operation_bonus":                                 129,
	"resource_wars_reward":                            131,
	"duel_wager_escrow":                               132,
	"duel_wager_payment":                              133,
	"duel_wager_refund":                               134,
	"reaction":                                        135,
	"structure_gate_jump":                             140,
	"skill_purchase":                                  141,
	"item_trader_payment":                             142,
	"ess_escrow_transfer":                             155,
}

type DBJournal struct {
	CharID          int
	CorpID          int
	WalletDivID     int
	Amount          float64
	Balance         float64
	Context_id      int64
	Context_id_type int
	Date            int64
	Description     int64
	First_party_id  int64
	ID              int64
	Reason          int64
	Ref_type        int
	Second_party_id int64
	Tax             float64
	Tax_receiver_id int64
}

type DBJournalLink struct {
	JournalID     int64
	IndustryJobID int64
	ContractID    int64
}

func (obj *Model) createJournalLinkTable() {
	if !obj.checkTableExists("journal_links") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "journal_links" (
			"journalID" INT,
			"industryJobID" INT,
			"contractID" INT);`)
		util.CheckErr(err)
	}
}
func (obj *Model) AddJournalLinkEntry(jLnkItem *DBJournalLink) DBresult {
	whereClause := fmt.Sprintf(`journalID="%d"`, jLnkItem.JournalID)
	num := obj.getNumEntries("journal_links", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "journal_links" (
			journalID,
			industryJobID,
			contractID)
			values(?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			jLnkItem.JournalID,
			jLnkItem.IndustryJobID,
			jLnkItem.ContractID)
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
func (obj *Model) GetJournalForContract(ctrID int) (result *DBJournal) {
	queryString := fmt.Sprintf(`
		SELECT 
			charId,
			corpId,
			walletDivId,
			amount,
			balance,
			context_id,
			context_id_type,
			date,
			description,
			first_party_id,
			id,
			reason,
			ref_type,
			second_party_id,
			tax,
			tax_receiver_id 
		FROM journal 
		Inner JOIN
			journal_links ON journal.id = journal_links.journalID
		WHERE journal_links.contractID = %d;`, ctrID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var jourItem DBJournal
		rows.Scan(
			&jourItem.CharID,
			&jourItem.CorpID,
			&jourItem.WalletDivID,
			&jourItem.Amount,
			&jourItem.Balance,
			&jourItem.Context_id,
			&jourItem.Context_id_type,
			&jourItem.Date,
			&jourItem.Description,
			&jourItem.First_party_id,
			&jourItem.ID,
			&jourItem.Reason,
			&jourItem.Ref_type,
			&jourItem.Second_party_id,
			&jourItem.Tax,
			&jourItem.Tax_receiver_id)
		result = &jourItem
		break
	}
	return
}

func (obj *Model) createJournalTable() {
	if !obj.checkTableExists("journal") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "journal" (
			"charId" INT,
			"corpId" INT,
			"walletDivId" INT,
			"amount" REAL,
			"balance" REAL,
			"context_id" INT,
			"context_id_type" INT,
			"date" INT,
			"description" INT,
			"first_party_id" INT,
			"id" INT,
			"reason"  INT,
			"ref_type" INT,
			"second_party_id" INT,
			"tax" REAL,
			"tax_receiver_id" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddJournalEntry(jouItem *DBJournal) DBresult {
	whereClause := fmt.Sprintf(`id="%d"`, jouItem.ID)
	num := obj.getNumEntries("journal", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "journal" (
			charId,
			corpId,
			walletDivId,
			amount,
			balance,
			context_id,
			context_id_type,
			date,
			description,
			first_party_id,
			id,
			reason,
			ref_type,
			second_party_id,
			tax,
			tax_receiver_id)
			values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			jouItem.CharID,
			jouItem.CorpID,
			jouItem.WalletDivID,
			jouItem.Amount,
			jouItem.Balance,
			jouItem.Context_id,
			jouItem.Context_id_type,
			jouItem.Date,
			jouItem.Description,
			jouItem.First_party_id,
			jouItem.ID,
			jouItem.Reason,
			jouItem.Ref_type,
			jouItem.Second_party_id,
			jouItem.Tax,
			jouItem.Tax_receiver_id)
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

func (obj *Model) GetBounties(corpId int) []*DBTable {
	var retval []*DBTable
	queryString := fmt.Sprintf(`
		SELECT 
			date, 
			amount,
			string_table.string as AltName,
			stringMain.string as MainName
		FROM journal 
		Inner JOIN 
		   corp_members ON corp_members.character_id = second_party_id
		INNER JOIN		   
			(SELECT character_id, name FROM corp_members) corpRef2Main
			ON corpRef2Main.character_id = corp_members.main_id
		INNER JOIN		   
			string_table ON corp_members.name= string_table.string_hash
		INNER JOIN		   
			(SELECT string_hash, string FROM string_table) stringMain
			ON corpRef2Main.name= stringMain.string_hash
		WHERE (ref_type=85 or ref_type=55) and corpId=%d
		ORDER BY date DESC;`, corpId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var elem DBTable
		rows.Scan(&elem.Time, &elem.Amount, &elem.AltName, &elem.MainName)
		retval = append(retval, &elem)
	}
	return retval
}

func (obj *Model) GetBountyTable(corpId int) *MonthlyTable {
	bountyTable := obj.GetBounties(corpId)
	return obj.GetMonthlyTable(corpId, bountyTable, 12)
}

func (obj *Model) GetJournal(charId int, corpId int, corp bool) []*DBJournal {
	retval := make([]*DBJournal, 0, 100)
	whereClause := fmt.Sprintf(`charId=%d and walletDivId==0`, charId)
	if corp {
		whereClause = fmt.Sprintf(`corpId=%d and walletDivId!=0`, corpId)
	}
	queryString := fmt.Sprintf(`SELECT
			charId,
			corpId,
			walletDivId,
			amount,
			balance,
			context_id,
			context_id_type,
			date,
			description,
			first_party_id,
			id,
			reason,
			ref_type,
			second_party_id,
			tax,
			tax_receiver_id
		FROM journal
		WHERE %s ORDER BY date DESC;`, whereClause)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var jourItem DBJournal
		rows.Scan(
			&jourItem.CharID,
			&jourItem.CorpID,
			&jourItem.WalletDivID,
			&jourItem.Amount,
			&jourItem.Balance,
			&jourItem.Context_id,
			&jourItem.Context_id_type,
			&jourItem.Date,
			&jourItem.Description,
			&jourItem.First_party_id,
			&jourItem.ID,
			&jourItem.Reason,
			&jourItem.Ref_type,
			&jourItem.Second_party_id,
			&jourItem.Tax,
			&jourItem.Tax_receiver_id)
		retval = append(retval, &jourItem)
	}
	return retval
}

func (obj *Model) GetBalanceOverTime(charId int, corpId int, corp bool, days int) (balance float64) {
	journal := obj.GetJournal(charId, corpId, corp)
	for _, j := range journal {
		if j.Date > time.Now().AddDate(0, 0, -days).Unix() {
			balance += j.Amount
		}
	}
	return
}

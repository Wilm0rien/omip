package model

import (
	"testing"
)

func TestModelMain(t *testing.T) {
	DeleteDb(DbNameModelTest)
	modelObj := NewModel(DbNameModelTest, true, false)
	tableList := []string{
		"characters",
		"debouncing",
		"adash_accounts",
		"pap_status",
		"contr_items",
		"contracts",
		"corp_info",
		"industry_jobs",
		"journal",
		"killmails",
		"k_attackers",
		"k_victims",
		"k_items",
		"string_table",
		"structure_services",
		"universe_names",
		"corp_members",
		"wallet_history",
		"structure_info",
		"notifications",
		"transactions",
		"market_prices",
		"structure_name",
		"mining_data",
		"mining_observers"}
	for _, tableName := range tableList {
		if !modelObj.checkTableExists(tableName) {
			t.Errorf("could not find table %s", tableName)
		}
	}
	modelObj.CloseDB()
}

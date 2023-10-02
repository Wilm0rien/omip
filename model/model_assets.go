package model

import (
	"github.com/Wilm0rien/omip/util"
)

type DBAsset struct {
	CharID          int
	CorpID          int
	IsBlueprintCopy int
	IsSingleton     int
	ItemId          int64
	LocationFlag    int64
	LocationId      int64
	LocationType    int64
	Quantity        int
	TypeId          int
}

func (obj *Model) createAssetTable() {
	if !obj.checkTableExists("assets") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "assets" (
			"charId" INT,
			"corpId" INT,
			"is_blueprint_copy" INT,
			"is_singleton" INT,
			"item_id" INT,
			"location_flag" INT,
			"location_id" INT,
			"location_type" INT,
			"quantity" INT,
			"type_id" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddAssetEntry(asset *DBAsset) DBresult {
	// todo every assets gets a time stamp value
	// item_id seems to be unique identifier
	return 0
}

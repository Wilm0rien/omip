package model

import (
	"fmt"
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
	Timestamp       int64
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
			"type_id" INT,
		    "timestamp" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddAssetEntry(asset *DBAsset) DBresult {
	whereClause := fmt.Sprintf(
		`item_id="%d" and`+`location_id="%d" and`+
			`quantity="%d" `, asset.ItemId, asset.LocationId, asset.Quantity)
	num := obj.getNumEntries("assets", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "assets" (
			charId,
			corpId,
			is_blueprint_copy,
			is_singleton,
			item_id,
			location_flag,
			location_id,
			location_type,
			quantity,
			type_id,
			timestamp) 
			values (?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			asset.CharID,
			asset.CorpID,
			asset.IsBlueprintCopy,
			asset.IsSingleton,
			asset.ItemId,
			asset.LocationFlag,
			asset.LocationId,
			asset.LocationType,
			asset.Quantity,
			asset.TypeId,
			asset.Timestamp)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

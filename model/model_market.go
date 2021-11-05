package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"strconv"
	"strings"
)

type DBMarketItem struct {
	AdjustedPrice float64
	AveragePrice  float64
	TypeId        int
}

const (
	BULK_ITEMS_MAX = 200
)

func (obj *Model) createPriceTable() {
	if !obj.checkTableExists("market_prices") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "market_prices" (
		    "adjusted_price" REAL,
		    "average_price" REAL,
		    "type_id" INT
		    );`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddMarketItems(mItems []DBMarketItem) {
	obj.DeleteMarketItems()
	running := true
	for running {
		localList := make([]DBMarketItem, 0, BULK_ITEMS_MAX)
		if len(mItems) > BULK_ITEMS_MAX {
			localList = append(localList, mItems[0:BULK_ITEMS_MAX]...)
			mItems = mItems[BULK_ITEMS_MAX:]
		} else {
			localList = mItems
		}
		obj.AddMarketItemsBatch(localList)
		if len(localList) < BULK_ITEMS_MAX {
			running = false
		}
	}

}
func (obj *Model) AddMarketItemsBatch(mItems []DBMarketItem) {
	util.Assert(len(mItems) > 0)
	sqlStr := `INSERT INTO "market_prices" (
				 adjusted_price,
				 average_price,
				 type_id)
				 VALUES  `
	vals := make([]interface{}, 0, 10000)
	for _, mItem := range mItems {
		obj.ItemAvgPrice[mItem.TypeId] = mItem.AveragePrice
		sqlStr += "(?, ?, ?),"
		vals = append(vals, mItem.AdjustedPrice, mItem.AveragePrice, mItem.TypeId)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	sqlStr = ReplaceSQL(sqlStr, "?")

	stmt, err := obj.DB.Prepare(sqlStr)
	util.CheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(vals...)
	util.CheckErr(err)
	affect, err := res.RowsAffected()
	util.Assert(affect > 0)
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func (obj *Model) DeleteMarketItems() {
	queryStr := fmt.Sprintf("DELETE FROM market_prices;")
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	_, err = stmt.Exec()
	util.CheckErr(err)
}

func (obj *Model) GetMarketItems() (retval []*DBMarketItem) {
	retval = make([]*DBMarketItem, 0, 10)
	queryString := fmt.Sprintf(`SELECT adjusted_price, average_price, type_id FROM market_prices;`)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var foundItem DBMarketItem
		rows.Scan(&foundItem.AdjustedPrice, &foundItem.AveragePrice, &foundItem.TypeId)
		retval = append(retval, &foundItem)
	}

	return retval

}

func (obj *Model) GetMarketItem(typeID int) (retval *DBMarketItem) {
	var foundItem DBMarketItem
	queryString := fmt.Sprintf(`SELECT adjusted_price, average_price FROM market_prices WHERE type_id=%d;`, typeID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&foundItem.AdjustedPrice, &foundItem.AveragePrice)
		break
	}
	if foundItem.AveragePrice > 0 || foundItem.AdjustedPrice > 0 {
		retval = &foundItem
	}
	return retval
}
func (obj *Model) GetItemValue(typeID int) (retval float64) {
	marketItemValue, ok := obj.ItemAvgPrice[typeID]
	if ok {
		retval = marketItemValue
	}
	return retval
}

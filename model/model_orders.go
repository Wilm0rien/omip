package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

const (
	Order_Range_1           = 0
	Order_Range_10          = 1
	Order_Range_2           = 2
	Order_Range_20          = 3
	Order_Range_3           = 4
	Order_Range_30          = 5
	Order_Range_4           = 6
	Order_Range_40          = 7
	Order_Range_5           = 8
	Order_Range_region      = 9
	Order_Range_solarsystem = 10
	Order_Range_station     = 11
)

var orderRange = map[string]int{
	"1":           Order_Range_1,
	"10":          Order_Range_10,
	"2":           Order_Range_2,
	"20":          Order_Range_20,
	"3":           Order_Range_3,
	"30":          Order_Range_30,
	"4":           Order_Range_4,
	"40":          Order_Range_40,
	"5":           Order_Range_5,
	"region":      Order_Range_region,
	"solarsystem": Order_Range_solarsystem,
	"station":     Order_Range_station,
}

type DBOrder struct {
	CharacterID   int
	CorporationID int
	Duration      int64
	Escrow        float64
	IsBuyOrder    bool
	IsCorporation bool
	Issued        int64
	LocationID    int64
	MinVolume     int
	OrderID       int64
	Price         float64
	Range         int
	RegionID      int
	TypeID        int
	VolumeRemain  int
	VolumeTotal   int
}

func (obj *Model) createOrderTable() {
	if !obj.checkTableExists("orders") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "orders" (		
		    character_id INT, 
			corporation_id INT,	
			duration INT,
			escrow REAL,
			is_buy_order INT,
			is_corporation INT,
			issued INT,
			location_id INT,
			min_volume INT,
			order_id INT,
			price REAL,
			range INT,
			region_id INT,
			type_id INT,
			volume_remain INT,
			volume_total INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddOrderEntry(order *DBOrder) DBresult {
	whereClause := fmt.Sprintf(`order_id="%d"`, order.OrderID)
	num := obj.getNumEntries("orders", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		//log.Printf("AddContractEntry\n")
		stmt, err := obj.DB.Prepare(`
		INSERT INTO "orders" (
				character_id, 
				corporation_id,		                         
		        duration,
				escrow,
				is_buy_order,
				is_corporation,
				issued,
				location_id,
				min_volume,
				order_id,
				price,
				range,
				region_id,
				type_id,
				volume_remain,
				volume_total) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			order.CharacterID,
			order.CorporationID,
			order.Duration,
			order.Escrow,
			order.IsBuyOrder,
			order.IsCorporation,
			order.Issued,
			order.LocationID,
			order.MinVolume,
			order.OrderID,
			order.Price,
			order.Range,
			order.RegionID,
			order.TypeID,
			order.VolumeRemain,
			order.VolumeTotal)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		whereClause2 := fmt.Sprintf(`order_id="%d" AND volume_remain="%d"`, order.OrderID, order.VolumeRemain)
		num2 := obj.getNumEntries("orders", whereClause2)
		if num2 == 0 {
			//log.Printf("UpdateContractEntry\n")
			stmt, err := obj.DB.Prepare(`
				UPDATE "orders" SET
						volume_remain=?
					WHERE order_id=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				order.VolumeRemain,
				order.OrderID)
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

func (obj *Model) GetOrdersByIssuerId(id int, corp bool) []*DBOrder {
	var where string
	if corp {
		where = fmt.Sprintf("corporation_id=%d AND is_corporation=1", id)
	} else {
		where = fmt.Sprintf("character_id=%d AND for_corporation=0", id)
	}
	queryString := fmt.Sprintf(`SELECT 
										character_id,
										corporation_id,
										duration,
										escrow,
										is_buy_order,
										is_corporation,
										issued,
										location_id,
										min_volume,
										order_id,
										price,
										range,
										region_id,
										type_id,
										volume_remain,
										volume_total,
										FROM orders
										WHERE %s
										ORDER BY issued ASC;`, where)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	retval := make([]*DBOrder, 0, 10)
	for rows.Next() {
		var order DBOrder
		rows.Scan(&order.CharacterID,
			&order.CorporationID,
			&order.Duration,
			&order.Escrow,
			&order.IsBuyOrder,
			&order.IsCorporation,
			&order.Issued,
			&order.LocationID,
			&order.MinVolume,
			&order.OrderID,
			&order.Price,
			&order.Range,
			&order.RegionID,
			&order.TypeID,
			&order.VolumeRemain,
			&order.VolumeTotal)
		retval = append(retval, &order)
	}
	return retval
}

func (obj *Model) OrderRangeStr2Int(rangeStr string) int {
	retVal, _ := orderRange[rangeStr]
	return retVal
}
func (obj *Model) OrderRangeInt2Str(rangeID int) (retval string) {
	retval = "unkown"
	for k, v := range orderRange {
		if v == rangeID {
			retval = k
			break
		}
	}
	return retval
}

package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
)

type Order struct {
	Duration      int64   `json:"duration"`
	Escrow        float64 `json:"escrow"`
	IsBuyOrder    bool    `json:"is_buy_order"`
	IsCorporation bool    `json:"is_corporation"`
	Issued        string  `json:"issued"`
	LocationID    int64   `json:"location_id"`
	MinVolume     int     `json:"min_volume"`
	OrderID       int64   `json:"order_id"`
	Price         float64 `json:"price"`
	Range         string  `json:"range"`
	RegionID      int     `json:"region_id"`
	TypeID        int     `json:"type_id"`
	VolumeRemain  int     `json:"volume_remain"`
	VolumeTotal   int     `json:"volume_total"`
}

func (obj *Ctrl) UpdateOrders(char *EsiChar, corp bool) {
	var url string
	if !char.UpdateFlags.Journal {
		return
	}
	if corp {
		url = fmt.Sprintf("https://esi.evetech.net/v3/corporations/%d/orders/?datasource=tranquility", char.CharInfoExt.CooperationId)
	} else {
		url = fmt.Sprintf("https://esi.evetech.net/v2/characters/%d/orders/?datasource=tranquility", char.CharInfoData.CharacterID)
	}
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var orderList []Order
	contentError := json.Unmarshal(bodyBytes, &orderList)
	if contentError != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
	}
	obj.UpdateOrdersInDb(char, corp, orderList)

}

func (obj *Ctrl) UpdateOrdersInDb(char *EsiChar, corp bool, orderList []Order) {
	for _, order := range orderList {
		dbOrder := obj.convertEsiOrder2DB(&order, char)
		result := obj.Model.AddOrderEntry(dbOrder)
		if result == model.DBR_Updated {
			issuer := char.CharInfoData.CharacterName
			if corp {
				corpObj := obj.GetCorp(char)
				if corpObj != nil {
					issuer = corpObj.Name
				} else {
					issuer = "unkown corp name"
				}
			}
			buySell := "sell"
			if order.IsBuyOrder {
				buySell = "buy"
			}
			orderItemStr := obj.Model.GetTypeString(order.TypeID)
			obj.AddLogEntry(fmt.Sprintf("%s %s order changed %s : %d/%d",
				issuer, buySell, orderItemStr, order.VolumeRemain, order.VolumeTotal))
		}
	}
}

func (obj *Ctrl) convertEsiOrder2DB(order *Order, char *EsiChar) *model.DBOrder {
	var newOrder model.DBOrder
	newOrder.CharacterID = char.CharInfoData.CharacterID
	newOrder.CorporationID = char.CharInfoExt.CooperationId
	newOrder.Duration = order.Duration
	newOrder.Escrow = order.Escrow
	newOrder.IsBuyOrder = order.IsBuyOrder
	newOrder.IsCorporation = order.IsCorporation
	newOrder.Issued = util.ConvertTimeStrToInt(order.Issued)
	newOrder.LocationID = order.LocationID
	newOrder.MinVolume = order.MinVolume
	newOrder.OrderID = order.OrderID
	newOrder.Price = order.Price
	newOrder.Range = obj.Model.OrderRangeStr2Int(order.Range)
	newOrder.RegionID = order.RegionID
	newOrder.TypeID = order.TypeID
	newOrder.VolumeRemain = order.VolumeRemain
	newOrder.VolumeTotal = order.VolumeTotal
	return &newOrder
}

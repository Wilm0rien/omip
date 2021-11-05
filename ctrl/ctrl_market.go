package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
)

type MarketPrice struct {
	Adjusted_price float64 `json:"adjusted_price"`
	Average_price  float64 `json:"average_price"`
	Type_id        int     `json:"type_id"`
}

func (obj *Ctrl) UpdateMarket(char *EsiChar, corp bool) {
	url := fmt.Sprintf("https://esi.evetech.net/v1/markets/prices/?datasource=tranquility")
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	var priceList []MarketPrice
	contentError := json.Unmarshal(bodyBytes, &priceList)
	if contentError != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
		return
	}
	var dbItemList []model.DBMarketItem
	for _, price := range priceList {
		var mItem model.DBMarketItem
		mItem.TypeId = price.Type_id
		mItem.AveragePrice = price.Average_price
		mItem.AdjustedPrice = price.Adjusted_price
		dbItemList = append(dbItemList, mItem)
	}
	obj.Model.AddMarketItems(dbItemList)
}

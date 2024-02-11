package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"log"
	"math/rand"
	"time"
)

// https://developers.eveonline.com/blog/article/esi-mining-ledger

type MiningObservers struct {
	LastUpdated  string `json:"last_updated"`
	ObserverID   int64  `json:"observer_id"`
	ObserverType string `json:"observer_type"`
}

type MiningData struct {
	CharacterID    int32  `json:"character_id"`
	LastUpdated    string `json:"last_updated"`
	Quantity       int32  `json:"quantity"`
	RecordedCorpId int32  `json:"recorded_corporation_id"`
	TypeId         int32  `json:"type_id"`
}

// UpdateCorpMiningObs retrieve list of corp mining observers via /corporation/{corporation_id}/mining/observers/
func (obj *Ctrl) UpdateCorpMiningObs(char *EsiChar, _UnusedCorp bool) {
	// needs esi-industry.read_corporation_mining.v1
	pageID := 1
	for {
		url := fmt.Sprintf("https://esi.evetech.net/v1/corporation/%d/mining/observers?datasource=tranquility&page=%d",
			char.CharInfoExt.CooperationId, pageID)
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var miningObsList []*MiningObservers
		contentError := json.Unmarshal(bodyBytes, &miningObsList)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR parsing url %s error %s", url, contentError.Error()))
			break
		}
		for _, miningObserver := range miningObsList {
			newObs := obj.convertEsiMOBS2DB(miningObserver)
			db1R := obj.Model.AddMiningObsEntry(newObs)
			util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)
			//TODO structureName := obj.GetStructureNameCached(miningObserver.ObserverID, char)
			obj.getMiningData(char, miningObserver.ObserverID)
		}
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}
}

// getMiningData retrieve Moon mining data via /corporation/{corporation_id}/mining/observers/{observer_id}/
func (obj *Ctrl) getMiningData(char *EsiChar, observerID int64) {
	// needs esi-industry.read_corporation_mining.v1

	pageID := 1
	for {
		url := fmt.Sprintf("https://esi.evetech.net/v1/corporation/%d/mining/observers/%d/?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, observerID, pageID)
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var miningData []*MiningData
		contentError := json.Unmarshal(bodyBytes, &miningData)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR parsing url %s error %s", url, contentError.Error()))
			break
		}
		for _, elem := range miningData {
			dbMiningData := obj.convertEsiMiningData2DB(elem, observerID)
			db1R := obj.Model.AddMiningDataEntry(dbMiningData)
			util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)
		}

		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
	}

}

func (obj *Ctrl) convertEsiMOBS2DB(mObs *MiningObservers) *model.DBMiningObserver {
	var newMObs model.DBMiningObserver
	newMObs.LastUpdated = util.ConvertDateStrToInt(mObs.LastUpdated)
	newMObs.ObserverID = mObs.ObserverID
	newMObs.ObserverType = obj.Model.AddStringEntry(mObs.ObserverType)
	return &newMObs
}

func (obj *Ctrl) convertEsiMiningData2DB(md *MiningData, obsID int64) *model.DBMiningData {
	var newMinDat model.DBMiningData
	newMinDat.LastUpdated = util.ConvertDateStrToInt(md.LastUpdated)
	newMinDat.CharacterID = int(md.CharacterID)
	newMinDat.RecordedCorporationID = int(md.RecordedCorpId)
	newMinDat.TypeID = int(md.TypeId)
	newMinDat.Quantity = int(md.Quantity)
	newMinDat.ObserverID = obsID
	return &newMinDat
}

func (obj *Ctrl) GenerateMiningData() {
	charIDs := []int32{95281762, 2115692519, 2115417359, 95067057, 2115636466, 2114367476, 2113199519, 2115448095, 2114908444, 2115714045, 2115692575}
	//common
	Cobalite := []string{"Cobaltite", "Copious Cobaltite", "Twinkling Cobaltite"}
	Euxenite := []string{"Euxenite", "Copious Euxenite", "Twinkling Euxenite"}
	Scheelite := []string{"Scheelite", "Copious Scheelite", "Twinkling Scheelite"}
	Titanite := []string{"Titanite", "Copious Titanite", "Twinkling Titanite"}

	// Exceptional
	Loparite := []string{"Loparite", "Bountiful Loparite", "Shining Loparite"}
	Monazite := []string{"Monazite", "Bountiful Monazite", "Shining Monazite"}
	Xenotime := []string{"Xenotime", "Bountiful Xenotime", "Shining Xenotime"}
	Ytterbite := []string{"Ytterbite", "Bountiful Ytterbite", "Shining Ytterbite"}

	// Rare
	Carnotite := []string{"Carnotite", "Glowing Carnotite", "Replete Carnotite"}
	Cinnabar := []string{"Cinnabar", "Glowing Cinnabar", "Replete Cinnabar"}
	Pollucite := []string{"Pollucite", "Glowing Pollucite", "Replete Pollucite"}
	Zircon := []string{"Zircon", "Glowing Zircon", "Replete Zircon"}

	// Ubiquitous
	Bitumens := []string{"Bitumens", "Brimful Bitumens", "Glistening Bitumens"}
	Coesite := []string{"Coesite", "Brimful Coesite", "Glistening Coesite"}
	Sylvite := []string{"Sylvite", "Brimful Sylvite", "Glistening Sylvite"}
	Zeolites := []string{"Zeolites", "Brimful Zeolites", "Glistening Zeolites"}

	// uncommon
	Chromite := []string{"Chromite", "Lavish Chromite", "Shimmering Chromite"}
	Otavite := []string{"Otavite", "Lavish Otavite", "Shimmering Otavite"}
	Sperrylite := []string{"Sperrylite", "Lavish Sperrylite", "Shimmering Sperrylite"}
	Vanadinite := []string{"Vanadinite", "Lavish Vanadinite", "Shimmering Vanadinite"}

	// normal ores
	Arkonor := []string{"Arkonor", "Crimson Arkonor", "Flawless Arkonor", "Prime Arkonor"}
	Bezdnacine := []string{"Bezdnacine", "Hadal Bezdnacine"}
	Bistot := []string{"Bistot", "Cubic Bistot", "Monoclinic Bistot", "Triclinic Bistot"}
	Crokite := []string{"Crokite", "Crystalline Crokite", "Pellucid Crokite", "Sharp Crokite"}
	DarkOchre := []string{"Dark Ochre", "Jet Ochre", "Obsidian Ochre", "Onyx Ochre"}
	Ducinium := []string{"Ducinium", "Imperial Ducinium", "Noble Ducinium", "Royal Ducinium"}
	Eifyrium := []string{"Eifyrium", "Augmented Eifyrium", "Boosted Eifyrium", "Doped Eifyrium"}
	Gneiss := []string{"Gneiss", "Brilliant Gneiss", "Iridescent Gneiss", "Prismatic Gneiss"}
	Hedbergite := []string{"Hedbergite", "Glazed Hedbergite", "Lustrous Hedbergite", "Vitric Hedbergite"}
	Hemorphite := []string{"Hemorphite", "Radiant Hemorphite", "Scintillating Hemorphite", "Vivid Hemorphite"}
	Jaspet := []string{"Jaspet", "Luminous Kernite", "Pure Jaspet", "Immaculate Jaspet"}
	Kernite := []string{"Kernite", "Pristine Jaspet", "Fiery Kernite", "Resplendant Kernite"}
	Mercoxit := []string{"Mercoxit", "Vitreous Mercoxit", "Magma Mercoxit", "Resplendant Kernite"}
	Mordunium := []string{"Mordunium", "Plum Mordunium", "Plunder Mordunium", "Prize Mordunium"}
	Omber := []string{"Omber", "Golden Omber", "Platinoid Omber", "Silvery Omber"}
	Plagioclase := []string{"Plagioclase", "Rich Plagioclase", "Sparkling Plagioclase", "Azure Plagioclase"}
	Pyroxeres := []string{"Pyroxeres", "Solid Pyroxeres", "Viscous Pyroxeres", "Opulent Pyroxeres"}
	Rakovene := []string{"Rakovene", "Nesosilicate Rakovene", "Hadal Rakovene", "Abyssal Rakovene"}
	Scordite := []string{"Scordite", "Massive Scordite", "Glossy Scordite", "Condensed Scordite"}
	Spodumain := []string{"Spodumain", "Gleaming Spodumain", "Dazzling Spodumain", "Bright Spodumain"}
	Talassonite := []string{"Talassonite", "Hadal Talassonite", "Abyssal Talassonite"}
	Veldspar := []string{"Veldspar", "Dense Veldspar", "Stable Veldspar", "Concentrated Veldspar"}

	ores := make([]string, 0, 100)
	ores = append(ores, Cobalite...)
	ores = append(ores, Euxenite...)
	ores = append(ores, Scheelite...)
	ores = append(ores, Titanite...)
	ores = append(ores, Loparite...)
	ores = append(ores, Monazite...)
	ores = append(ores, Xenotime...)
	ores = append(ores, Ytterbite...)
	ores = append(ores, Carnotite...)
	ores = append(ores, Cinnabar...)
	ores = append(ores, Pollucite...)
	ores = append(ores, Zircon...)
	ores = append(ores, Bitumens...)
	ores = append(ores, Coesite...)
	ores = append(ores, Sylvite...)
	ores = append(ores, Zeolites...)
	ores = append(ores, Chromite...)
	ores = append(ores, Otavite...)
	ores = append(ores, Sperrylite...)
	ores = append(ores, Vanadinite...)
	ores = append(ores, Arkonor...)
	ores = append(ores, Bezdnacine...)
	ores = append(ores, Bistot...)
	ores = append(ores, Crokite...)
	ores = append(ores, DarkOchre...)
	ores = append(ores, Ducinium...)
	ores = append(ores, Eifyrium...)
	ores = append(ores, Gneiss...)
	ores = append(ores, Hedbergite...)
	ores = append(ores, Hemorphite...)
	ores = append(ores, Jaspet...)
	ores = append(ores, Kernite...)
	ores = append(ores, Mercoxit...)
	ores = append(ores, Mordunium...)
	ores = append(ores, Omber...)
	ores = append(ores, Plagioclase...)
	ores = append(ores, Pyroxeres...)
	ores = append(ores, Rakovene...)
	ores = append(ores, Scordite...)
	ores = append(ores, Spodumain...)
	ores = append(ores, Talassonite...)
	ores = append(ores, Veldspar...)

	log.Printf("%v", ores)
	for month := int32(1); month <= 12; month++ {
		for day := int32(1); day <= 30; day++ {
			num := 10 + rand.Intn(100)
			for i := 0; i < num; i++ {
				randomChar := rand.Intn(len(charIDs))
				var newRandObvj MiningData
				newRandObvj.CharacterID = charIDs[randomChar]
				newRandObvj.LastUpdated = fmt.Sprintf("2024-%02d-%02d", month, day)
				newRandObvj.Quantity = int32(1 + rand.Intn(100000))

			}
		}
	}

}

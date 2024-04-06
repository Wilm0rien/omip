package ctrl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"log"
	"math/rand"
	"os"
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
func (obj *Ctrl) UpdateCorpMiningObs(char *EsiChar, corp bool) {
	if corp == false {
		// skip character update for mining observers
		return
	}
	// needs esi-industry.read_corporation_mining.v1
	if !char.UpdateFlags.Mining {
		//obj.AddLogEntry(fmt.Sprintf("UpdateCorpMiningObs disabled for %s %d", char.CharInfoData.CharacterName, char.CharInfoExt.CooperationId))
		return
	}
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
		obj.AddLogEntry(fmt.Sprintf("reading mining observers successful found %d entries", len(miningObsList)))
		for _, miningObserver := range miningObsList {
			newObs := obj.convertEsiMOBS2DB(miningObserver, char.CharInfoExt.CooperationId)
			db1R := obj.Model.AddMiningObsEntry(newObs)
			util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)

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
	corpCache := make(map[int32]int)
	charCache := make(map[int32]int)
	ownMemberIdMap := obj.Model.GetCorpMemberIdMap(char.CharInfoExt.CooperationId)
	newNameRequests := make([]int, 0, 10)
	pageID := 1
	// handle observer name
	strucName := obj.GetStructureNameCached(observerID, char)
	for {
		url := fmt.Sprintf("https://esi.evetech.net/v1/corporation/%d/mining/observers/%d/?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, observerID, pageID)
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var miningDataEntry []*MiningData
		contentError := json.Unmarshal(bodyBytes, &miningDataEntry)
		if contentError != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR parsing url %s error %s", url, contentError.Error()))
			break
		}
		obj.AddLogEntry(fmt.Sprintf("reading mining getMiningData %s successful found %d entries", strucName, len(miningDataEntry)))

		endIdx := len(miningDataEntry)
		lastTs := time.Now()
		for idx, elem := range miningDataEntry {
			if time.Since(lastTs).Milliseconds() > 500 {
				obj.GuiStatusCB("mining update", 1)
				obj.GuiStatusCB(fmt.Sprintf("%3.2f%%", float32(idx)/float32(endIdx)*100), 2)
				lastTs = time.Now()
			}
			dbMiningData := obj.convertEsiMiningData2DB(elem, observerID, char.CharInfoExt.CooperationId)
			db1R := obj.Model.AddMiningDataEntry(dbMiningData)
			util.Assert(db1R == model.DBR_Inserted || db1R == model.DBR_Updated)
			// check if corp exists in db
			if _, ok := corpCache[elem.RecordedCorpId]; !ok {
				corpCache[elem.RecordedCorpId] = 1
				// TODO check if corp infos are updated during update routine
				if _, result := obj.Model.GetCorpInfoEntry(int(elem.RecordedCorpId)); result != model.DBR_Success {
					if _, ok2 := obj.GetCorpInfoFromEsi(char, int(elem.RecordedCorpId)); !ok2 {
						obj.AddLogEntry(fmt.Sprintf("ERROR adding corp id %d", elem.RecordedCorpId))
					}
				}
			}
			if _, ok2 := ownMemberIdMap[int(elem.CharacterID)]; !ok2 {
				if _, ok := charCache[elem.CharacterID]; !ok {
					charCache[elem.RecordedCorpId] = 1
					if !obj.Model.NameExists(int(elem.CharacterID)) {
						newNameRequests = append(newNameRequests, int(elem.CharacterID))
					}
				}
			}
		}
		if pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}

	}
	if len(newNameRequests) > 0 {
		obj.getUniverseNames(newNameRequests, char)
	}

}

func (obj *Ctrl) convertEsiMOBS2DB(mObs *MiningObservers, corpID int) *model.DBMiningObserver {
	var newMObs model.DBMiningObserver
	newMObs.LastUpdated = util.ConvertDateStrToInt(mObs.LastUpdated)
	newMObs.ObserverID = mObs.ObserverID
	newMObs.ObserverType = obj.Model.AddStringEntry(mObs.ObserverType)
	newMObs.OwnerCorpID = corpID
	return &newMObs
}

func (obj *Ctrl) convertEsiMiningData2DB(md *MiningData, obsID int64, corpID int) *model.DBMiningData {
	var newMinDat model.DBMiningData
	newMinDat.LastUpdated = util.ConvertDateStrToInt(md.LastUpdated)
	newMinDat.CharacterID = int(md.CharacterID)
	newMinDat.RecordedCorporationID = int(md.RecordedCorpId)
	newMinDat.TypeID = int(md.TypeId)
	newMinDat.Quantity = int(md.Quantity)
	newMinDat.ObserverID = obsID
	newMinDat.OwnerCorpID = corpID
	return &newMinDat
}

func (obj *Ctrl) GetOreValueByM3(oreTypeID int, volumeM3 float64) (value float64, err error) {
	if props := obj.Model.GetSdePropsByID(oreTypeID); props != nil {
		amount := int(volumeM3 / props.GetVolume())
		value, err = obj.GetOreValueByAmount(oreTypeID, amount)
	}
	return
}

func (obj *Ctrl) GetOreValueByAmount(oreTypeID int, amount int) (totalValue float64, err error) {
	if props := obj.Model.GetSdePropsByID(oreTypeID); props != nil {
		// example Scordite oreTypeID 1228	amount 1000
		if numberOfBatches := amount / 100; numberOfBatches != 0 {
			// number of batches = 10
			for _, contMat := range props.Materials {
				// contained Materials in 100m Scordidte: 150 Tritanium + 90 Pyerite
				if contMatValue, ok := obj.Model.ItemAvgPrice[contMat.MaterialTypeID]; ok && contMatValue != 0 {
					totalValue += contMatValue * float64(numberOfBatches) * float64(contMat.Quantity)
				} else {
					err = errors.New(fmt.Sprintf("item not found ID %d value %3.2f", contMat.MaterialTypeID, contMatValue))
					totalValue = 0
					break
				}
			}
		} else {
			// no error if batch is too small
		}

	}
	return
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

	exores := make([][]string, 0, 100)
	exores = append(exores, Loparite)
	exores = append(exores, Monazite)
	exores = append(exores, Xenotime)
	exores = append(exores, Ytterbite)

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

	for _, ore := range ores {
		if id := obj.Model.GetItemID(ore); id == 0 {
			log.Printf("ERROR ID not found for %s", ore)
		}
	}

	list := make([]*MiningData, 0, 100)
	for month := int32(1); month <= 12; month++ {
		for day := int32(1); day <= 28; day++ {
			num := 10 + rand.Intn(10)
			filterMap := make(map[int]map[int]*MiningData)

			for i := 0; i < num; i++ {

				randomChar := rand.Intn(len(charIDs))
				var newRandObvj MiningData
				newRandObvj.CharacterID = charIDs[randomChar]
				newRandObvj.LastUpdated = fmt.Sprintf("2024-%02d-%02d", month, day)
				newRandObvj.Quantity = int32(1 + rand.Intn(1000))
				newRandObvj.RecordedCorpId = 98627127
				randOre := rand.Intn(3)
				randSubOre := rand.Intn(3)
				selected := exores[randOre][randSubOre]
				id := obj.Model.GetItemID(selected)
				newRandObvj.TypeId = int32(id)
				if val, ok := filterMap[int(newRandObvj.CharacterID)]; ok {
					if val2, ok2 := val[id]; ok2 {
						val2.Quantity += newRandObvj.Quantity
					} else {
						filterMap[int(newRandObvj.CharacterID)][id] = &newRandObvj
					}
				} else {
					filterMap[int(newRandObvj.CharacterID)] = make(map[int]*MiningData)
					filterMap[int(newRandObvj.CharacterID)][id] = &newRandObvj
				}

			}
			for _, val := range filterMap {
				for _, val2 := range val {
					list = append(list, val2)
				}
			}
		}
	}
	data, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		log.Printf("ERROR creating json %s", err.Error())
	} else {
		f2, _ := os.Create("output.json")
		f2.WriteString(string(data))
		f2.Close()
	}
}

// UpdateMiningMeta ensure we know all the names from all entities
func (obj *Ctrl) UpdateMiningMeta(char *EsiChar, corp bool) {
	// TODO do this only once because it is not relevant which char does this update
	obj.UpdateMiningChars(char)
	obj.UpdateMiningCorps(char)
}

func (obj *Ctrl) UpdateMiningChars(char *EsiChar) {
	charMap := obj.Model.GetMiningCharMap()
	newNameRequests := make([]int, 0, 10)
	for charID, _ := range charMap {
		if !obj.Model.NameExists(charID) {
			newNameRequests = append(newNameRequests, charID)
		}
	}
	if len(newNameRequests) > 0 {
		obj.getUniverseNames(newNameRequests, char)
	}
}
func (obj *Ctrl) UpdateMiningCorps(char *EsiChar) {
	obj.Model.DeleteCorpCache()
	coprpMap := obj.Model.GetMiningCorpMap()
	allyMap := make(map[int]int)
	for corpID, _ := range coprpMap {
		if dbCorp, ok := obj.GetCorpInfoFromEsi(char, corpID); ok {
			if dbCorp.AllianceId != 0 {
				allyMap[dbCorp.AllianceId] = 1
			}
		}

	}
	obj.UpdateMiningAllies(char, allyMap)
}

func (obj *Ctrl) UpdateMiningAllies(char *EsiChar, allyMap map[int]int) {
	for allyID, _ := range allyMap {
		if _, result := obj.Model.GetAllyInfoEntry(allyID); result != model.DBR_Success {
			obj.GetAllyInfoFromEsi(char, allyID)
		}
	}
}

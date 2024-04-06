package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"log"
	"net/http"
	"time"
)

const miningData_v1 = `
					[
					  {
						"last_updated": "2024-02-10",
						"character_id": 2115636466,
						"recorded_corporation_id": 98627127,
						"type_id": 17448,
						"quantity": 2292
					  },
					  {
						"last_updated": "2024-02-10",
						"character_id": 2115636466,
						"recorded_corporation_id": 98627127,
						"type_id": 17452,
						"quantity": 1250
					  },
					  {
						"last_updated": "2024-02-10",
						"character_id": 2115636466,
						"recorded_corporation_id": 98627127,
						"type_id": 20,
						"quantity": 1265
					  },
					  {
						"last_updated": "2024-02-10",
						"character_id": 2115636466,
						"recorded_corporation_id": 98627127,
						"type_id": 17449,
						"quantity": 6888
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 96227676,
						"recorded_corporation_id": 98179071,
						"type_id": 17449,
						"quantity": 6888
					  },
					  {
						"last_updated": "2024-03-10",
						"character_id": 2121966805,
						"recorded_corporation_id": 98760370,
						"type_id": 20,
						"quantity": 1500
					  }
					]
					`
const miningData_v2 = `
			[
					  {
						"last_updated": "2024-03-09",
						"character_id": 96227676,
						"recorded_corporation_id": 98179071,
						"type_id": 17449,
						"quantity": 20000
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 96338783,
						"recorded_corporation_id": 98179071,
						"type_id": 20,
						"quantity": 10000
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 2115714045,
						"recorded_corporation_id": 98627127,
						"type_id": 17452,
						"quantity": 5000
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 2115417359,
						"recorded_corporation_id": 98627127,
						"type_id": 17452,
						"quantity": 5000
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 2115417359,
						"recorded_corporation_id": 98627127,
						"type_id": 17452,
						"quantity": 5000
					  },
					  {
						"last_updated": "2024-03-09",
						"character_id": 2114367476,
						"recorded_corporation_id": 98627127,
						"type_id": 17452,
						"quantity": 5000
					  },
					  {
						"last_updated": "2024-03-10",
						"character_id": 92799763,
						"recorded_corporation_id": 98370861,
						"type_id": 17452,
						"quantity": 10000
					  },
					  {
						"last_updated": "2024-03-10",
						"character_id": 92799763,
						"recorded_corporation_id": 98370861,
						"type_id": 20,
						"quantity": 1000
					  },
					  {
						"last_updated": "2024-03-10",
						"character_id": 2121966805,
						"recorded_corporation_id": 98760370,
						"type_id": 20,
						"quantity": 1000
					  },
					  {
						"last_updated": "2024-03-10",
						"character_id": 2121966805,
						"recorded_corporation_id": 98760370,
						"type_id": 46681,
						"quantity": 1000
					  }
					]

`

const universeData = `
[
    {
        "category": "character",
        "id": 95281762,
        "name": "Zuberi Mwanajuma"
    },
    {
        "category": "character",
        "id": 2115692519,
        "name": "Rob Barrington"
    },
    {
        "category": "character",
        "id": 2115417359,
        "name": "Koriyi Chan"
    },
    {
        "category": "character",
        "id": 95067057,
        "name": "Gwen Facero"
    },
    {
        "category": "character",
        "id": 2115636466,
        "name": "Ion of Chios"
    },
    {
        "category": "character",
        "id": 2114367476,
        "name": "Koriyo -Skill1 Skill"
    },
    {
        "category": "character",
        "id": 2113199519,
        "name": "azullunes"
    },
    {
        "category": "character",
        "id": 2115448095,
        "name": "Koriyo -Skill2 Skill"
    },
    {
        "category": "character",
        "id": 2114908444,
        "name": "Gudrun Yassavi"
    },
    {
        "category": "character",
        "id": 2115714045,
        "name": "Luke Lovell"
    },
    {
        "category": "character",
        "id": 2115692575,
        "name": "Jill Kenton"
    },
	{
		"category": "character",
		"id": 96227676,
		"name": "Ares Aurelius"
	},
	{
		"category": "character",
		"id": 96338783,
		"name": "Ood Tau-5"
	},
	  {
		"category": "character",
		"id": 92799763,
		"name": "seth Yassavi"
	  },
	  {
		"category": "character",
		"id": 2121966805,
		"name": "The Zombie Overlord"
	  }
]

`

var HttpRequestMock func(req *http.Request) (bodyBytes []byte, err error, resp *http.Response)
var HttpPostDataMock string

func (obj *Ctrl) GetUniverseMock(inputReq string, inputUniData string) (result string) {
	uniData := []byte(inputUniData)
	var uniNames []universeNames
	err := json.Unmarshal(uniData, &uniNames)
	if err != nil {
		panic("must never happen!")
	}
	uniIdMap := make(map[int]*universeNames)
	for _, elem := range uniNames {
		var newElem universeNames
		newElem = elem
		uniIdMap[elem.ID] = &newElem
	}

	resultList := make([]universeNames, 0, 10)
	var reqIdList []int
	err = json.Unmarshal([]byte(inputReq), &reqIdList)
	if err != nil {
		panic("must never happen!")
	}
	for _, elem := range reqIdList {
		if val, ok := uniIdMap[elem]; ok {
			resultList = append(resultList, *val)
		}
	}
	data, err2 := json.MarshalIndent(resultList, "", "\t")
	if err2 != nil {
		panic("must never happen!")
	}

	return string(data)
}

func (obj *Ctrl) GetRequestMock() (result ReqMockFuncT) {
	dummyToken := `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJvbWlwIHRlc3QgdG9rZW4iLCJpYXQiOjE2MzU3NTE2NjMsImV4cCI6MTY2NzI4NzY2MywiYXVkIjoid3d3LmV2ZW9ubGluZS5jb20iLCJzdWIiOiJDSEFSQUNURVI6RVZFOjIxMTU2MzY0NjYiLCJuYW1lIjoiSW9uIG9mIENoaW9zIiwiRW1haWwiOiJqcm9ja2V0QGV4YW1wbGUuY29tIn0.kbAkeoDeGh3Hh5mVtKJNl-vJScbbkOOlTYTs1mR91ZY`
	expiresOn := util.UnixTS2DateTimeStr(time.Now().Add(1199 * time.Second).Unix())
	result = func(req *http.Request) (bodyBytes []byte, err error, resp *http.Response) {
		resp = &http.Response{
			StatusCode: http.StatusOK,
		}
		switch req.URL.String() {
		case "https://login.eveonline.com/v2/oauth/token":
			bodyBytes = []byte("{\"access_token\":\"" + dummyToken + "\",\"expires_in\":1199,\"token_type\":\"Bearer\",\"refresh_token\":\"refresh_token_dummytoken\"}")
		case "https://login.eveonline.com/oauth/verify":
			resultString := fmt.Sprintf("{\"CharacterID\":2115636466,\"CharacterName\":\"Ion of Chios\",\"ExpiresOn\":\"%s\",\"Scopes\":\"publicData esi-wallet.read_character_wallet.v1 esi-wallet.read_corporation_wallet.v1 esi-universe.read_structures.v1 esi-killmails.read_killmails.v1 esi-corporations.read_corporation_membership.v1 esi-corporations.read_structures.v1 esi-industry.read_character_jobs.v1 esi-contracts.read_character_contracts.v1 esi-killmails.read_corporation_killmails.v1 esi-corporations.track_members.v1 esi-wallet.read_corporation_wallets.v1 esi-characters.read_notifications.v1 esi-contracts.read_corporation_contracts.v1 esi-corporations.read_starbases.v1 esi-industry.read_corporation_jobs.v1\",\"TokenType\":\"Character\",\"CharacterOwnerHash\":\"dummyhash=\",\"IntellectualProperty\":\"EVE\"}",
				expiresOn)
			bodyBytes = []byte(resultString)
		case "https://esi.evetech.net/v5/characters/2115636466":
			bodyBytes = []byte("{\"ancestry_id\":9,\"birthday\":\"2019-08-28T18:48:02Z\",\"bloodline_id\":2,\"corporation_id\":98627127,\"description\":\"\",\"gender\":\"male\",\"name\":\"Ion of Chios\",\"race_id\":1,\"security_status\":0.0}")
		case "https://esi.evetech.net/v2/corporations/98627127/roles/?datasource=tranquility":
			bodyBytes = []byte("[{\"character_id\":95067057,\"grantable_roles\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"grantable_roles_at_base\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"grantable_roles_at_hq\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"grantable_roles_at_other\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"roles\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"roles_at_base\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"roles_at_hq\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"],\"roles_at_other\":[\"Director\",\"Personnel_Manager\",\"Accountant\",\"Security_Officer\",\"Factory_Manager\",\"Station_Manager\",\"Auditor\",\"Hangar_Take_1\",\"Hangar_Take_2\",\"Hangar_Take_3\",\"Hangar_Take_4\",\"Hangar_Take_5\",\"Hangar_Take_6\",\"Hangar_Take_7\",\"Hangar_Query_1\",\"Hangar_Query_2\",\"Hangar_Query_3\",\"Hangar_Query_4\",\"Hangar_Query_5\",\"Hangar_Query_6\",\"Hangar_Query_7\",\"Account_Take_1\",\"Account_Take_2\",\"Account_Take_3\",\"Account_Take_4\",\"Account_Take_5\",\"Account_Take_6\",\"Account_Take_7\",\"Diplomat\",\"Config_Equipment\",\"Container_Take_1\",\"Container_Take_2\",\"Container_Take_3\",\"Container_Take_4\",\"Container_Take_5\",\"Container_Take_6\",\"Container_Take_7\",\"Rent_Office\",\"Rent_Factory_Facility\",\"Rent_Research_Facility\",\"Junior_Accountant\",\"Config_Starbase_Equipment\",\"Trader\",\"Communications_Officer\",\"Contract_Manager\",\"Starbase_Defense_Operator\",\"Starbase_Fuel_Technician\",\"Fitting_Manager\"]},{\"character_id\":95281762,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2113199519,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2114367476,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2114908444,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115417359,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115448095,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115636466,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[\"Director\"],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115692519,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115692575,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]},{\"character_id\":2115714045,\"grantable_roles\":[],\"grantable_roles_at_base\":[],\"grantable_roles_at_hq\":[],\"grantable_roles_at_other\":[],\"roles\":[],\"roles_at_base\":[],\"roles_at_hq\":[],\"roles_at_other\":[]}]")
		case "https://esi.evetech.net/v5/corporations/98627127?datasource=tranquility":
			bodyBytes = []byte("{\"ceo_id\":95067057,\"creator_id\":2115636466,\"date_founded\":\"2020-01-09T17:27:50Z\",\"description\":\"Enter a description of your corporation here.\",\"home_station_id\":60011386,\"member_count\":11,\"name\":\"Feynman Electrodynamics\",\"shares\":1000,\"tax_rate\":0.0,\"ticker\":\"FYDYN\",\"url\":\"http:\\/\\/\"}")
		case "https://esi.evetech.net/v4/corporations/98627127/members/?datasource=tranquility":
			bodyBytes = []byte("[95281762,2115692519,2115417359,95067057,2115636466,2114367476,2113199519,2115448095,2114908444,2115714045,2115692575]")
		case "https://esi.evetech.net/v3/universe/names/":
			inputStr := universeData
			uniresult := obj.GetUniverseMock(HttpPostDataMock, inputStr)
			bodyBytes = []byte(uniresult)
		case "https://esi.evetech.net/v1/characters/2115636466/killmails/recent/":
			bodyBytes = []byte("[]")
		case "https://esi.evetech.net/v1/corporations/98627127/killmails/recent/":
			bodyBytes = []byte("[]")
		case "https://esi.evetech.net/v2/universe/structures/1000000000001/?datasource=tranquility":
			bodyBytes = []byte(`
			{
				  "name": "PhantomSystem - PhantomBase",
				  "owner_id": 98627127,
				  "position": {
					"x": 3362585191667,
					"y": 315435898402,
					"z": -1172591694720
				  },
				  "solar_system_id": 40001725,
				  "type_id": 35825
			}`)
		case "https://esi.evetech.net/v2/universe/structures/1000000000002/?datasource=tranquility":
			bodyBytes = []byte(`
			{
				"name": "ShadowSystem - ShadowBase",
				"owner_id": 98627127,
				  "position": {
					"x": -193810504896,
					"y": -4432966889720,
					"z": -1161937975641
				  },
				  "solar_system_id": 30002780,
				  "type_id": 35825
			}`)
		case "https://esi.evetech.net/v1/corporation/98627127/mining/observers?datasource=tranquility&page=1":
			bodyBytes = []byte(`
			[
				{
					"last_updated": "2024-03-10",
					"observer_id": 1000000000001,
					"observer_type": "structure"
				},
				{
					"last_updated": "2024-03-10",
					"observer_id": 1000000000002,
					"observer_type": "structure"
				}
			]
			`)
		case "https://esi.evetech.net/v1/corporation/98627127/mining/observers/1000000000001/?datasource=tranquility&page=1":
			bodyBytes = []byte(miningData_v1)
		case "https://esi.evetech.net/v1/corporation/98627127/mining/observers/1000000000002/?datasource=tranquility&page=1":
			bodyBytes = []byte(miningData_v2)
		case "https://esi.evetech.net/v5/corporations/98179071?datasource=tranquility":
			bodyBytes = []byte(`
				{
				  "alliance_id": 150097440,
				  "ceo_id": 743137360,
				  "creator_id": 90813985,
				  "date_founded": "2013-02-22T21:36:06Z",
				  "description": "<font size=\"13\" color=\"#99ffffff\"></font><font size=\"12\" color=\"#bfffffff\">Chat (public): </font><font size=\"12\" color=\"#ff6868e1\"><a href=\"joinChannel:-45630598//98179071//57114\">Omicron Pub</a></font><font size=\"12\" color=\"#bfffffff\"> </font>",
				  "home_station_id": 60003316,
				  "member_count": 22,
				  "name": "Omicron Project",
				  "shares": 1000,
				  "tax_rate": 0.10000000149011612,
				  "ticker": "OMIP",
				  "url": "",
				  "war_eligible": true
				}`)
		case "https://esi.evetech.net/v5/corporations/98760370?datasource=tranquility":
			bodyBytes = []byte(`
{
  "alliance_id": 99003214,
  "ceo_id": 2121966805,
  "creator_id": 2121966805,
  "date_founded": "2024-01-07T19:52:37Z",
  "description": "",
  "home_station_id": 60015157,
  "member_count": 1,
  "name": "Zombies are People Too",
  "shares": 1000,
  "tax_rate": 0.10000000149011612,
  "ticker": "Z0M8Y",
  "url": "http://",
  "war_eligible": true
}`)
		case "https://esi.evetech.net/v5/corporations/98370861?datasource=tranquility":
			bodyBytes = []byte(`
				{
				  "alliance_id": 1354830081,
				  "ceo_id": 1687616923,
				  "creator_id": 95195474,
				  "date_founded": "2015-01-12T22:01:46Z",
				  "description": "<font size=\"14\" color=\"#bfffffff\"></font><font size=\"12\" color=\"#ffffe400\"><a href=\"https://youtu.be/vJMvbhvj1bw\">Enlist Today!</a><br><br><br></font><font size=\"12\" color=\"#ffd98d00\"><a href=\"showinfo:16159//1354830081\">Goonswarm Federation's</a></font><font size=\"12\" color=\"#bfffffff\"> largest member Corporation where Eve's most knowledgeable veterans help teach pilots to fly.   New players come here to learn from the best, old players come here to revitalize their Eve experience.  KarmaFleet has an abundance of unique supportive services available only to our pilots like Corp level SRP(on top of Alliance SRP), ISK payouts to our FCs, and a PLEX for PVP program that rewards our pilots.  If you are looking for the perfect mixture of an extremely active and yet well organized experience this is the Corporation for you.   <br><br>Our recruitment application is back online in a beta function, link is above! Stop by our public channel to ask any questions regarding the recruitment process: </font><font size=\"12\" color=\"#ff6868e1\"><a href=\"joinChannel:player_46462f6ef10711e8b7609abe94f5aa9b//None//None\">Public Karma</a><a href=\"joinChannel:-66414639//None//None\">.</a></font><font size=\"12\" color=\"#bfffffff\"> A list of current recruiters can be found in the channel's MOTD. Applications are thoroughly reviewed and do take time to process.  Please be patient.<br><br>Applications that have not been first submitted through our recruitment website will be declined.<br><br>Karmafleet is an open corporation. What this means:<br><br>1. No one will ask you for isk to join.<br>2. No one will ask you to trade your assets to them in order to get you out to null.<br><br>If someone asks for these things, you are getting scammed and you need to notify a Recruiter immediately. <br><br>For any <b>Recruitment</b> matters please contact </font><font size=\"12\" color=\"#ffd98d00\"><loc><a href=\"showinfo:1376//2112859268\">Ross July</a></loc></font><font size=\"12\" color=\"#bfffffff\">, </font><font size=\"12\" color=\"#ffd98d00\"><loc><a href=\"showinfo:1378//477616994\">Wholly Thoc</a></loc></font><font size=\"12\" color=\"#bfffffff\">, </font><font size=\"12\" color=\"#ffd98d00\"><loc><a href=\"showinfo:1380//92771480\">Jiriha</a></loc></font><font size=\"12\" color=\"#bfffffff\">, </font><font size=\"12\" color=\"#ffd98d00\"><loc><a href=\"showinfo:1377//2115557767\">Mekosq</a></loc></font><font size=\"12\" color=\"#bfffffff\">, or hop into our public channel listed above.<br><br></font><font size=\"12\" color=\"#ffffe400\"><loc><a href=\"https://www.youtube.com/watch?v=FOv2yyqdI5Y\">Be a part of something big.</a><a href=\"https://youtu.be/vJMvbhvj1bw\"> </a></loc><br><br><a href=\"https://twitter.com/karmafleet\">@karmafleet</a><br><loc><a href=\"https://www.reddit.com/r/Karmafleet/\">/r/karmaflee</loc>t</a><br><br></font><font size=\"12\" color=\"#bfffffff\">If you are a High-Sec resident please check out our High-Sec Corp, </font><font size=\"12\" color=\"#ffd98d00\"><a href=\"showinfo:2//98615046\">KarmaFleet University</a></font><font size=\"12\" color=\"#bfffffff\">. Open invites, all are welcome.</font>",
				  "home_station_id": 60011008,
				  "member_count": 8136,
				  "name": "KarmaFleet",
				  "shares": 1000,
				  "tax_rate": 0.15000000596046448,
				  "ticker": "SNOOO",
				  "url": "https://discord.gg/karmafleet",
				  "war_eligible": true
				}
		`)
		case "https://esi.evetech.net/v4/alliances/150097440?datasource=tranquility":
			bodyBytes = []byte(`
			{
			  "creator_corporation_id": 147945511,
			  "creator_id": 740426190,
			  "date_founded": "2009-05-25T15:36:00Z",
			  "executor_corporation_id": 98391152,
			  "name": "Get Off My Lawn",
			  "ticker": "LAWN"
			}
`)
		case "https://esi.evetech.net/v4/alliances/99003214?datasource=tranquility":
			bodyBytes = []byte(`
{
  "creator_corporation_id": 98199293,
  "creator_id": 91432281,
  "date_founded": "2013-05-06T22:20:14Z",
  "executor_corporation_id": 98199293,
  "name": "Brave Collective",
  "ticker": "BRAVE"
}
`)
		case "https://esi.evetech.net/v4/alliances/1354830081?datasource=tranquility":
			bodyBytes = []byte(`
			{
			  "creator_corporation_id": 459299583,
			  "creator_id": 240070320,
			  "date_founded": "2010-06-01T05:36:00Z",
			  "executor_corporation_id": 1344654522,
			  "name": "Goonswarm Federation",
			  "ticker": "CONDI"
			}
		`)
		default:
			log.Printf("ERROR cannot find URL %s", req.URL.String())

		}

		return bodyBytes, err, resp
	}
	return
}

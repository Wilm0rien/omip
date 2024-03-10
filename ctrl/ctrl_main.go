package ctrl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	ConfigFileName = "omipCfg.dat"
	DebounceTime   = 60 * 60 * 4
	TstCfgJson     = "TestCtrlMain.json"
)

type ReqMockFuncT func(req *http.Request) (bodyBytes []byte, err error, resp *http.Response)

var TestMockReq ReqMockFuncT
var CtrlTestEnable bool

type AuthCallBack func(newChar *EsiChar)

type urlCache struct {
	timestamp int64
	cache     []byte
}

// NVConfigData non-volatile config
type NVConfigData struct {
	PeriodFilter string
}
type UpdateFunc func(char *EsiChar, corp bool)
type Ctrl struct {
	Model       *model.Model
	Esi         EsiData
	Svr         httpSrvData
	AuthCb      AuthCallBack
	ADash       map[int]*ADashClient
	urlCache    map[string]urlCache
	LogEntries  []string
	AddLogCB    func(entry string)
	GuiStatusCB func(entry string, fieldId int)
	structCache map[int64]string
	Up          EsiUpdate
	NotifyInfo  map[int64]bool // initialized in UpdateAllData()
}

type EsiUpdate struct {
	UpdateFuncList []UpdateFunc
	UpdateMutex    sync.Mutex
	JobList        []string
}

func NewCtrl(model *model.Model) *Ctrl {
	var obj Ctrl
	obj.Model = model
	obj.NotifyInfo = make(map[int64]bool)
	obj.LogEntries = make([]string, 0, 5)
	obj.Svr.cancelChan = make(chan struct{})
	obj.urlCache = make(map[string]urlCache)
	obj.ADash = make(map[int]*ADashClient)
	obj.structCache = make(map[int64]string)
	obj.populateUpdateFuncList()
	close(obj.Svr.cancelChan)
	return &obj
}

type EsiFileError struct {
	ErrorCode EsiFileErrCode
	FileName  string
	ExtErr    string
}

type EsiFileErrCode int

const (
	ESI_FILE_DECRYPT_ERROR EsiFileErrCode = iota
	ESI_FILE_OPEN_ERROR
	ESI_FILE_WRITE_ERROR
	ESI_FILE_JSON_ERROR
)

func (obj *Ctrl) populateUpdateFuncList() {
	obj.Up.JobList = []string{"Contracts", "ContractItems", "Industry", "KillMails", "Wallet", "CorpMembers", "Structures", "Notifications", "Transaction", "Orders", "Mails"}
	obj.Up.UpdateFuncList = make([]UpdateFunc, 0, 5)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateContracts)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateContractItems)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateIndustry)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateKillMails)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateWallet)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateCorpMembers)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateStructures)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateNotifications)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateTransaction)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateOrders)
	obj.Up.UpdateFuncList = append(obj.Up.UpdateFuncList, obj.UpdateMailLabels)

}

func (obj *Ctrl) UpdateChar(char *EsiChar) {
	obj.Up.UpdateMutex.Lock()
	defer obj.Up.UpdateMutex.Unlock()
	obj.UpdateJournal(char, false, 0)
	for idx, updateFunc := range obj.Up.UpdateFuncList {
		obj.UpdateGuiStatus1(fmt.Sprintf("%s %s", char.CharInfoData.CharacterName, obj.Up.JobList[idx]))
		updateFunc(char, false)
	}
}
func (obj *Ctrl) UpdateCorp(director *EsiChar) {
	obj.Up.UpdateMutex.Lock()
	corp := obj.GetCorp(director)
	defer obj.Up.UpdateMutex.Unlock()
	for idx, updateFunc := range obj.Up.UpdateFuncList {
		obj.UpdateGuiStatus1(fmt.Sprintf("%s %s", corp.Name, obj.Up.JobList[idx]))
		updateFunc(director, true)
	}
	for i := 1; i <= 7; i++ {
		obj.UpdateGuiStatus1(fmt.Sprintf("%s wallet division (%d)", corp.Name, i))
		obj.UpdateJournal(director, true, i)
	}
}

func (obj EsiFileError) Error() (result string) {
	switch obj.ErrorCode {
	case ESI_FILE_OPEN_ERROR:
		result = fmt.Sprintf("could not open file %s %s", obj.FileName, obj.ExtErr)
	case ESI_FILE_WRITE_ERROR:
		result = fmt.Sprintf("could not write file %s %s", obj.FileName, obj.ExtErr)
	case ESI_FILE_DECRYPT_ERROR:
		result = fmt.Sprintf("could not decrypt file %s", obj.FileName)
	case ESI_FILE_JSON_ERROR:
		result = fmt.Sprintf("could not parse jason data from file %s %s", obj.FileName, obj.ExtErr)
	}
	return result
}

func (obj *Ctrl) Save(cfgFileName string, testEnable bool) (retval error) {
	data, err := json.MarshalIndent(obj.Esi, "", "\t")
	if testEnable {
		cfgFileName = TstCfgJson
	}
	if err != nil {
		retval = &EsiFileError{ErrorCode: ESI_FILE_JSON_ERROR, FileName: cfgFileName, ExtErr: err.Error()}
		obj.AddLogEntry(fmt.Sprintf("ERROR write file %s", err))
		return
	} else {
		authData := obj.Model.LocalDir + "/" + cfgFileName
		f, err1 := os.Create(authData)
		if err1 != nil {
			retval = &EsiFileError{ErrorCode: ESI_FILE_WRITE_ERROR, FileName: cfgFileName, ExtErr: err1.Error()}
			return
		}
		passPhrase := util.GenSysPassphrase()
		_, err2 := f.WriteString(string(util.Encrypt(data, passPhrase)))
		if err2 != nil {
			retval = &EsiFileError{ErrorCode: ESI_FILE_WRITE_ERROR, FileName: cfgFileName, ExtErr: err2.Error()}
			f.Close()
			return
		}
		if testEnable {
			authDataUnenc := obj.Model.LocalDir + "/" + cfgFileName + ".clear"
			f2, _ := os.Create(authDataUnenc)
			f2.WriteString(string(data))
			log.Printf("writing clear json to %s", authDataUnenc)
		}
	}
	return retval
}

func (obj *Ctrl) Load(cfgFileName string, testEnable bool) (retval error) {
	if testEnable {
		cfgFileName = TstCfgJson
		CtrlTestEnable = true
		response := obj.GetRequestMock()
		HttpRequestMock = response
	}

	authData := obj.Model.LocalDir + "/" + cfgFileName
	if util.Exists(authData) {
		file, err := os.ReadFile(authData)
		if err != nil {
			retval = &EsiFileError{ErrorCode: ESI_FILE_OPEN_ERROR, FileName: authData, ExtErr: err.Error()}
		}
		if len(file) > 0 {
			passPhrase := util.GenSysPassphrase()
			fileData, success := util.Decrypt(file, passPhrase)
			if success {
				err = json.Unmarshal(fileData, &obj.Esi)
				if err != nil {
					retval = &EsiFileError{ErrorCode: ESI_FILE_JSON_ERROR, FileName: authData, ExtErr: err.Error()}
				}
			} else {
				retval = &EsiFileError{ErrorCode: ESI_FILE_DECRYPT_ERROR, FileName: authData}
			}
		}

	}
	if obj.Esi.ETags == nil {
		obj.Esi.ETags = make(map[string]string)
	}
	if obj.Esi.CacheEntries == nil {
		obj.Esi.CacheEntries = make(map[string]int64)
	}
	return retval
}

func (obj *Ctrl) GetCorpDirector(corpId int) *EsiChar {
	var director *EsiChar
	for _, char := range obj.Esi.EsiCharList {
		if char.CharInfoExt.CooperationId == corpId &&
			char.CharInfoExt.Director {
			director = char
			break
		}
	}
	return director
}

func (obj *Ctrl) GetCorp(char *EsiChar) *EsiCorp {
	var retval *EsiCorp
	for _, corp := range obj.Esi.EsiCorpList {
		if corp.CooperationId == char.CharInfoExt.CooperationId {
			retval = corp
			break
		}
	}
	if retval == nil {
		retval = obj.getAuthCorpInfo(char)
	}
	return retval
}

func (obj *Ctrl) AddLogEntry(entry string) {
	obj.Model.LogObj.Println(entry)
	timeStamp := util.ConvertUnixTimeToStr(time.Now().Unix())
	logEntry := fmt.Sprintf("%20.20s - %s", timeStamp, entry)
	if obj.AddLogCB != nil {
		obj.AddLogCB(logEntry)
	}
	obj.LogEntries = append(obj.LogEntries, logEntry)
}

func (obj *Ctrl) UpdateGuiStatus1(entry string) {
	if obj.GuiStatusCB != nil {
		obj.GuiStatusCB(entry, 1)
	}
}
func (obj *Ctrl) UpdateGuiStatus2(entry string) {
	if obj.GuiStatusCB != nil {
		obj.GuiStatusCB(entry, 2)
	}
}

func (obj *Ctrl) CheckUpdatePreCon() (ok bool, err error) {
	if len(obj.Esi.EsiCharList) == 0 {
		err = errors.New("no characters registered to update data")
		return
	}

	if obj.Model.DebounceEntryExists("update_string") && !obj.Model.DebugFlag {
		err = errors.New("please wait 5 minutes between updates")
		return
	}

	serverStatus := obj.CheckServerUp(obj.Esi.EsiCharList[0])
	if !serverStatus {
		err = errors.New("esi server is starting up or is not reachable. cannot update data")
		return
	}
	obj.Model.AddDebounceEntry("update_string")
	return true, nil
}

func (obj *Ctrl) UpdateAllDataCmd(updateProg func(c float64), finishCb func()) {
	obj.Up.UpdateMutex.Lock()
	defer obj.Up.UpdateMutex.Unlock()

	totalItems := (len(obj.Esi.EsiCharList) + len(obj.Esi.EsiCorpList) + 1) * len(obj.Up.UpdateFuncList)
	// add 1 journal request per character and 7 journal requests per corp
	totalItems += len(obj.Esi.EsiCharList) + (len(obj.Esi.EsiCorpList) * 7)
	var itemCount int
	if len(obj.Esi.EsiCharList) > 0 {
		obj.UpdateMarket(obj.Esi.EsiCharList[0], false)
	}
	for _, char := range obj.Esi.EsiCharList {
		obj.UpdateWallet(char, false)
		if char.AuthValid == AUTH_STATUS_INVALID {
			obj.AddLogEntry(fmt.Sprintf("skipping update for invalid auth for %s", char.CharInfoData.CharacterName))
			continue
		}
		itemCount++
		// NOTE: the journal has to be updated first to update the journal_links table
		// this is because only contracts with journal links are identified as relevant for being stored
		obj.UpdateJournal(char, false, 0)

		for idx, updateFunc := range obj.Up.UpdateFuncList {
			obj.UpdateGuiStatus1(fmt.Sprintf("%s %s", char.CharInfoData.CharacterName, obj.Up.JobList[idx]))
			updateFunc(char, false)
			itemCount++
			if updateProg != nil {
				updateProg(float64(itemCount) / float64(totalItems))
			}
		}

		if updateProg != nil {
			updateProg(float64(itemCount) / float64(totalItems))
		}
	}
	for _, corp := range obj.Esi.EsiCorpList {
		director := obj.GetCorpDirector(corp.CooperationId)
		if director != nil {
			for idx, updateFunc := range obj.Up.UpdateFuncList {
				obj.UpdateGuiStatus1(fmt.Sprintf("%s %s", corp.Name, obj.Up.JobList[idx]))
				updateFunc(director, true)
				itemCount++
				if updateProg != nil {
					updateProg(float64(itemCount) / float64(totalItems))
				}
			}
			for i := 1; i <= 7; i++ {
				obj.UpdateGuiStatus1(fmt.Sprintf("%s wallet division (%d)", corp.Name, i))
				obj.UpdateJournal(director, true, i)
				itemCount++
				if updateProg != nil {
					updateProg(float64(itemCount) / float64(totalItems))
				}
			}
		}
	}
	if finishCb != nil {
		finishCb()
	}
}

func (obj *Ctrl) GetCorpTicker(char *EsiChar) (corpTicker string) {
	corpObj := obj.GetCorp(char)
	if corpObj != nil {
		corpTicker = corpObj.Ticker
	}
	return
}

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
						"quantity": 6888
					  }
					]

`

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
			bodyBytes = []byte("[{\"category\":\"character\",\"id\":95281762,\"name\":\"Zuberi Mwanajuma\"},{\"category\":\"character\",\"id\":2115692519,\"name\":\"Rob Barrington\"},{\"category\":\"character\",\"id\":2115417359,\"name\":\"Koriyi Chan\"},{\"category\":\"character\",\"id\":95067057,\"name\":\"Gwen Facero\"},{\"category\":\"character\",\"id\":2115636466,\"name\":\"Ion of Chios\"},{\"category\":\"character\",\"id\":2114367476,\"name\":\"Koriyo -Skill1 Skill\"},{\"category\":\"character\",\"id\":2113199519,\"name\":\"azullunes\"},{\"category\":\"character\",\"id\":2115448095,\"name\":\"Koriyo -Skill2 Skill\"},{\"category\":\"character\",\"id\":2114908444,\"name\":\"Gudrun Yassavi\"},{\"category\":\"character\",\"id\":2115714045,\"name\":\"Luke Lovell\"},{\"category\":\"character\",\"id\":2115692575,\"name\":\"Jill Kenton\"}]")
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
				}
`)
		default:
			log.Printf("ERROR cannot find URL %s", req.URL.String())

		}

		return bodyBytes, err, resp
	}
	return
}

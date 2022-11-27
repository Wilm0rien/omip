package ctrl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"os"
	"sync"
	"time"
)

const (
	ConfigFileName = "omipCfg.dat"
	DebounceTime   = 60 * 60 * 4
	TstCfgJson     = "TestCtrlMain.json"
)

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
}

type EsiUpdate struct {
	UpdateFuncList []UpdateFunc
	UpdateMutex    sync.Mutex
	JobList        []string
}

func NewCtrl(model *model.Model) *Ctrl {
	var obj Ctrl
	obj.Model = model
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
	obj.Up.JobList = []string{"Contracts", "ContractItems", "Industry", "KillMails", "Wallet", "CorpMembers", "Structures", "Notifications", "Transaction", "Orders"}
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
	data, err := json.Marshal(obj.Esi)
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
	}
	return retval
}

func (obj *Ctrl) Load(cfgFileName string, testEnable bool) (retval error) {
	if testEnable {
		cfgFileName = TstCfgJson
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
	if obj.Model.DebounceEntryExists("update_string") {
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

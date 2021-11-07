package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"io/ioutil"
	"os"
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

type Ctrl struct {
	Model      *model.Model
	Esi        EsiData
	Svr        httpSrvData
	AuthCb     AuthCallBack
	ADash      map[int]*ADashClient
	urlCache   map[string]urlCache
	LogEntries []string
	AddLogCB   func(entry string)
}

func NewCtrl(model *model.Model) *Ctrl {
	var obj Ctrl
	obj.Model = model
	obj.LogEntries = make([]string, 0, 5)
	obj.Svr.cancelChan = make(chan struct{})
	obj.urlCache = make(map[string]urlCache)
	obj.ADash = make(map[int]*ADashClient)

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
		file, err := ioutil.ReadFile(authData)
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

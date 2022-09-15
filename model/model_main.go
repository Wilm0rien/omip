package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/util"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"regexp"
	"syscall"
	"time"
)

type Model struct {
	LocalDir        string
	LocalLogFile    string
	LocalImgDir     string
	LocalEtagCache  string
	LocalSqlDb      string
	LocalTypIDsFile string
	DB              *sql.DB
	LogObj          *log.Logger
	LogFileHandle   *os.File
	ItemIDs         map[int]string
	ItemAvgPrice    map[int]float64
}

const (
	AppName         = "omip"
	DbName          = "omip.db"
	DbNameModelTest = "omipModelTest.db"
	DbNameCtrlTest  = "omipCtrlTest.db"
	LogFileName     = "omip.log"
	LogFileNameTest = "omipTest.log"
	DebounceTimeSec = 60 * 5
	LogBufferSIze   = 5000
)

type DBresult int

const (
	DBR_Undefined DBresult = iota
	DBR_Updated
	DBR_Inserted
	DBR_Removed
	DBR_Failed
	DBR_Success
	DBR_Skipped
)

func NewModel(ldbName string, testEnable bool) *Model {
	var obj Model
	lLogFileName := LogFileName
	if testEnable {
		lLogFileName = LogFileNameTest
		ldbName = DbNameCtrlTest
	}
	appData := util.GetAppDataDir()
	obj.LocalDir = appData + "/" + AppName
	if !util.Exists(obj.LocalDir) {
		util.CreateDirectory(obj.LocalDir)
	}
	obj.LocalImgDir = obj.LocalDir + "/" + "images"
	if !util.Exists(obj.LocalImgDir) {
		util.CreateDirectory(obj.LocalImgDir)
	}
	obj.LocalEtagCache = obj.LocalDir + "/" + "etags"
	if !util.Exists(obj.LocalEtagCache) {
		util.CreateDirectory(obj.LocalEtagCache)
	}
	obj.LocalSqlDb = obj.LocalDir + "/" + ldbName

	obj.ItemIDs = make(map[int]string)
	obj.ItemAvgPrice = make(map[int]float64)

	fileBytes := []byte(SdeData)
	errJson := json.Unmarshal(fileBytes, &obj.ItemIDs)
	if errJson != nil {
		obj.LogObj.Printf("ERROR typeIds %s", errJson.Error())
	}

	if !util.Exists(obj.LocalSqlDb) {
		obj.createNewDb(obj.LocalSqlDb)
	} else {
		db, err := sql.Open("sqlite3", obj.LocalSqlDb)
		util.CheckErr(err)
		obj.DB = db
	}
	obj.LocalLogFile = obj.LocalDir + "/" + lLogFileName
	var err error
	obj.LogFileHandle, err = os.OpenFile(obj.LocalLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		obj.LogObj.Printf("error opening file: %v", err)
	}
	obj.LogObj = log.New(obj.LogFileHandle, "", log.Ldate|log.Ltime)

	obj.createADashAuthTable()
	obj.createADashTable()
	obj.createCharTable()
	obj.createDebounceTable()
	obj.createContrItemTable()
	obj.createContractTable()
	obj.createOrderTable()
	obj.createCorpMemberTable()
	obj.createCorpInfoTable()
	obj.createInduTable()
	obj.createJournalTable()
	obj.createJournalLinkTable()
	obj.createKMTables()
	obj.createStringTable()
	obj.createStructureInfoTable()
	obj.createStructureNameTable()
	obj.createStructureServiceTable()
	obj.createUniNamesTable()
	obj.createWalletHistoryTable()
	obj.createNotificationTable()
	obj.createPriceTable()
	obj.createTransactionsTable()

	mItems := obj.GetMarketItems()
	for _, item := range mItems {
		obj.ItemAvgPrice[item.TypeId] = item.AveragePrice
	}

	return &obj
}

type MonthlyTable struct {
	// ValCharPerMon --> map[characterName]map[year-month]float64
	// example ValCharPerMon["Wilm0rien"][21-06]=150000
	ValCharPerMon map[string]map[string]float64
	// MaxInMonth --> map[year-month]float64
	// example MaxInMonth["21-06"] = 3000000
	MaxInMonth map[string]float64
	SumInMonth map[string]float64
	MaxAllTime float64
}

type DBTable struct {
	Time     int64
	Amount   float64
	AltName  string
	MainName string
}

type DBKMTable struct {
	Time      int64
	Amount    float64
	MainName  string
	AltName   string
	KMId      int
	KMHash    string
	ZK_status ZkStatus
}

func (obj *Model) GetTypeString(ID int) string {
	var retval string
	if obj.ItemIDs != nil {
		if val, ok := obj.ItemIDs[ID]; ok {
			retval = val
		}
	}
	return retval
}

func (obj *Model) GetMonthlyTable(corpId int, inputTable []*DBTable, maxMonth int) *MonthlyTable {
	var omipTable MonthlyTable
	omipTable.ValCharPerMon = make(map[string]map[string]float64)
	omipTable.MaxInMonth = make(map[string]float64)
	omipTable.SumInMonth = make(map[string]float64)
	now := time.Now().Unix()
	maxSeconds := int64(maxMonth * 30 * 24 * 60 * 60)
	for _, elem := range inputTable {
		tm := time.Unix(elem.Time, 0)
		year, month, _ := tm.Date()
		diff := now - elem.Time
		// only accept kill from the last n month
		if diff < maxSeconds {
			dateStr := fmt.Sprintf("%02d-%02d", year-2000, month)
			if _, ok := omipTable.ValCharPerMon[elem.MainName]; !ok {
				omipTable.ValCharPerMon[elem.MainName] = make(map[string]float64)
			}
			omipTable.ValCharPerMon[elem.MainName][dateStr] += elem.Amount
			omipTable.SumInMonth[dateStr] += elem.Amount
			if omipTable.MaxInMonth[dateStr] < omipTable.ValCharPerMon[elem.MainName][dateStr] {
				omipTable.MaxInMonth[dateStr] = omipTable.ValCharPerMon[elem.MainName][dateStr]
			}
		}
	}
	for _, elem := range omipTable.SumInMonth {
		if elem > omipTable.MaxAllTime {
			omipTable.MaxAllTime = elem
		}
	}

	return &omipTable
}

func (obj *Model) createNewDb(dbname string) {
	os.Create(dbname)
	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		fmt.Println(err)
		panic("cannot open db")
	}
	obj.DB = db
	obj.createCharTable()
	obj.createDebounceTable()

}

func (obj *Model) createDebounceTable() {
	if !obj.checkTableExists("debouncing") {
		_, err := obj.DB.Exec("CREATE TABLE `debouncing` (`bouncing_id` INT, `timestamp` INT);")
		util.CheckErr(err)
	}
}

func (obj *Model) AddDebounceEntry(url string) {
	debounceID := util.Get32BitMd5FromString(url)
	currentTime := time.Now().Unix()
	stmt, err := obj.DB.Prepare("INSERT INTO debouncing (bouncing_id, timestamp) values (?,?)")
	util.CheckErr(err)
	_, err = stmt.Exec(debounceID, currentTime)
	util.CheckErr(err)
}

func (obj *Model) DebounceEntryExists(url string) bool {
	debounceID := util.Get32BitMd5FromString(url)
	debunceAbsTime := time.Now().Unix() - int64(DebounceTimeSec)
	whereClause := fmt.Sprintf(`timestamp>"%d" and bouncing_id="%d"`, debunceAbsTime, debounceID)
	num := obj.getNumEntries("debouncing", whereClause)
	var retval bool
	if num > 0 {
		retval = true
	}
	return retval
}

func (obj *Model) CloseDB() {
	obj.DB.Close()
	obj.LogFileHandle.Close()
}

func (obj *Model) getNumEntries(tableName string, whereClause string) int {
	queryString := fmt.Sprintf(`SELECT  COUNT(*)  from %s WHERE %s;`, tableName, whereClause)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var num int
	num = 0
	for rows.Next() {
		rows.Scan(&num)
	}
	return num
}

func (obj *Model) deleteEntries(tableName string, whereClause string) (retval int) {
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, whereClause)
	stmt, err := obj.DB.Prepare(queryStr)
	util.CheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec()
	util.CheckErr(err)
	affect, err := res.RowsAffected()
	util.CheckErr(err)
	if affect > 0 {
		retval = int(affect)
	}
	return retval
}

func (obj *Model) checkTableExists(tableName string) bool {
	whereClause := fmt.Sprintf("type='table' and tbl_name='%s'", tableName)
	num := obj.getNumEntries("sqlite_master", whereClause)
	return num == 1
}

func (obj *Model) LoadEtag(Etag string) (bodybytes []byte) {
	fileName := path.Join(obj.LocalEtagCache, Etag)
	content, err := os.ReadFile(fileName)
	if err == nil {
		bodybytes = content
	} else {
		log.Printf("error reading etag %s", err.Error())
	}
	return
}

func (obj *Model) StoreEtag(newEtag string, oldEtag string, bodybytes []byte) (retval bool) {
	NewFileName := path.Join(obj.LocalEtagCache, newEtag)
	OldFileName := path.Join(obj.LocalEtagCache, oldEtag)
	if !util.Exists(NewFileName) {
		err := os.WriteFile(NewFileName, bodybytes, 0644)
		if err != nil {
			log.Printf("error writing etag %s", err.Error())
		} else {
			retval = true
		}
		if oldEtag != "" && util.Exists(OldFileName) && oldEtag != newEtag {
			err2 := os.Remove(OldFileName)
			if err2 != nil {
				log.Printf("error writing etag %s", err2.Error())
			}
		}
	}
	return
}

func (obj *Model) ImportBounties() (retval bool) {
	dbname := "e:/upload/esi_static/esi_db.sqlite"
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println(err)
		panic("cannot open db")
	}
	type corp_journal_ext struct {
		TimeStamp   string
		RefType     string
		Description string
		Amount      float64
	}
	journalExt := make([]*corp_journal_ext, 0, 5)
	queryString := fmt.Sprintf(`SELECT  TimeStamp, RefType, Description, Amount FROM corp_journal;`)
	rows, dberr := db.Query(queryString)
	memberList := obj.GetDBCorpMembers(98179071)
	nameMapping := make(map[string]int)
	failmap := make(map[string]int)
	for _, member := range memberList {
		name, _ := obj.GetStringEntry(member.NameRef)
		nameMapping[name] = member.CharID
	}
	if dberr == nil {
		defer rows.Close()
		for rows.Next() {
			var elem corp_journal_ext
			if rowErr := rows.Scan(&elem.TimeStamp, &elem.RefType, &elem.Description, &elem.Amount); rowErr == nil {
				re := regexp.MustCompile(`^([a-zA-Z0-9 ]+)\sgot bounty prizes for killing pirates in`)
				result := re.FindStringSubmatch(elem.Description)
				if result != nil {
					charName := result[1]
					if charId, ok := nameMapping[charName]; ok {
						var newJournal DBJournal
						newJournal.CharID = charId
						newJournal.CorpID = 98179071
						newJournal.WalletDivID = 1
						newJournal.Amount = elem.Amount
						newJournal.Ref_type = 85
						newJournal.Date = util.ConvertTimeStrToInt(elem.TimeStamp)
						newJournal.Description = obj.AddStringEntry(elem.Description)
					} else {
						if _, ok2 := failmap[charName]; !ok2 {
							log.Printf("could not find %s", charName)
							failmap[charName] = 1
						}
						log.Printf("")
					}
				}
				journalExt = append(journalExt, &elem)
			}
		}
		closeErr := db.Close()
		if closeErr != nil {
			fmt.Println(err)
			panic("cannot close db")
		}
	}
	if len(journalExt) > 0 {
		log.Printf("%d entries found", len(journalExt))
	}

	return retval

}

func getAppHomeDir() string {
	appData := util.GetAppDataDir()
	return appData + "/" + AppName
}

func DeleteDb(dbName string) {
	LocalDir := getAppHomeDir()
	if util.Exists(LocalDir) {
		testDb := LocalDir + "/" + dbName
		if util.Exists(testDb) {
			err := syscall.Unlink(testDb)
			if err != nil {
				log.Printf("error deleting %s code %s", testDb, err.Error())
			}
		}
	}
}

package ctrl

import (
	"bufio"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	baseURL = "https://adashboard.info/"
)

var TestAdashFlag bool

type ADashClient struct {
	client     *http.Client
	Username   string
	Password   string
	CorpTicker string
	LoginOK    bool
	Model      *model.Model
	CorpID     int
	AddLogCB   func(entry string)
	PapLogMap  map[string]int
}

type PapElem struct {
	CorpId     int
	PapLink    string
	ChName     string
	CoTicker   string
	AlShort    string
	ShTypeName string
	Sub        string
	Loc        string
	Timestamp  string
}

var ADhttpGetMock func(url string, data url.Values) (bodyBytes []byte, err error, resp *http.Response)

func NewADashClient(username string, password string, ticker string, Model *model.Model, corpId int) *ADashClient {
	var obj ADashClient
	util.Assert(username != "")
	util.Assert(ticker != "")
	util.Assert(Model != nil)

	obj.Username = username
	obj.Password = password
	obj.CorpTicker = ticker
	obj.Model = Model
	obj.CorpID = corpId
	obj.PapLogMap = make(map[string]int)
	jar, _ := cookiejar.New(nil)
	obj.client = &http.Client{Jar: jar}
	return &obj
}

func (obj *ADashClient) aDHttpGet(url string) (bodyBytes []byte, err error, resp *http.Response) {
	if !TestAdashFlag {
		resp, err = obj.client.Get(url)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				bodyBytes, _ = ioutil.ReadAll(resp.Body)
				//log.Printf("url: %s\naDHttpGet\n%s\n", url, string(bodyBytes))
				resp.Body.Close()
			}
		}
	} else {
		if ADhttpGetMock != nil {
			bodyBytes, err, resp = ADhttpGetMock(url, nil)
		}
	}
	return bodyBytes, err, resp
}

func (obj *ADashClient) aDHttpPostForm(url string, data url.Values) (bodyBytes []byte, err error, resp *http.Response) {
	if !TestAdashFlag {
		resp, err = obj.client.PostForm(url, data)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				bodyBytes, _ = ioutil.ReadAll(resp.Body)
				//log.Printf("url: %s\naDHttpPostForm\n%s\n", url, string(bodyBytes))
				resp.Body.Close()
			}
		}
	} else {
		if ADhttpGetMock != nil {
			bodyBytes, err, resp = ADhttpGetMock(url, data)
		}
	}
	return bodyBytes, err, resp
}

func (obj *ADashClient) Login() bool {
	_, err1, response1 := obj.aDHttpGet(baseURL)
	if err1 != nil {
		obj.Model.LogObj.Printf("Error fetching response. %s", err1.Error())
		return false
	}
	if response1.StatusCode == http.StatusOK {
		data := url.Values{
			"Email address": {obj.Username},
			"Password":      {obj.Password},
		}
		bodyBytes, err2, response2 := obj.aDHttpPostForm(baseURL+"/login", data)
		if err2 != nil {
			obj.Model.LogObj.Printf("Error post form %s", err2.Error())
			return false
		}
		if response2.StatusCode == http.StatusOK {
			scanner := bufio.NewScanner(strings.NewReader(string(bodyBytes)))
			re := regexp.MustCompile(`<title>aD - (.*?)</title>`)
			for scanner.Scan() {
				result := re.FindStringSubmatch(scanner.Text())
				if result != nil {
					//fmt.Printf("LOGIN OK: %s\n", string(result[1]))
					obj.LoginOK = true
					break
				}
			}
		}
	}
	if !obj.LoginOK {
		fmt.Printf("LOGIN FAILED")
	}
	return obj.LoginOK
}

func (obj *ADashClient) CheckPapLinks() (result string) {
	var papsFound int
	if obj.LoginOK {
		listURL := fmt.Sprintf("https://adashboard.info/corporation/%s", obj.CorpTicker)
		htmlContent, err, _ := obj.aDHttpGet(listURL)
		if err != nil {
			return fmt.Sprintf("Error fetching response. %s", err.Error())
		} else {
			scanner := bufio.NewScanner(strings.NewReader(string(htmlContent)))
			re := regexp.MustCompile(`\/par\/view\/([0-9A-Za-z-]+)`)
			for scanner.Scan() {
				result2 := re.FindStringSubmatch(scanner.Text())
				if result2 != nil {
					papsFound ++
					if papsFound < 5 {
						papLink := result2[1]
						result +=papLink + ", "
					}
				}
			}
		}
	} else {
		result = "login failed"
	}
	if papsFound == 0 && result == "" {
		result = "NO PAPS FOUND"
	} else {
		result += fmt.Sprintf("%d PAPs found", papsFound)
	}
	return
}

func (obj *ADashClient) GetPapLinks() bool {
	if obj.LoginOK {
		listURL := fmt.Sprintf("https://adashboard.info/corporation/%s", obj.CorpTicker)
		bodyBytes, err, _ := obj.aDHttpGet(listURL)
		if err != nil {
			obj.Model.LogObj.Printf("Error fetching response. %s", err.Error())
			return false
		}
		alt2main := obj.Model.GetAltMap(obj.CorpID)
		obj.decodePap(string(bodyBytes), alt2main)
		//fmt.Printf("%s\n", string(bodyBytes))
	} else {
		obj.Model.LogObj.Printf("GetPapLinks ERROR: NOT LOGGED IN")
		return false
	}
	return true
}

func (obj *ADashClient) decodePap(htmlContent string, alt2main map[string]string) {
	scanner := bufio.NewScanner(strings.NewReader(htmlContent))
	re := regexp.MustCompile(`\/par\/view\/([0-9A-Za-z-]+)`)
	var papsFound int
	for scanner.Scan() {
		result := re.FindStringSubmatch(scanner.Text())
		if result != nil {
			papLink := result[1]

			if !obj.Model.PapLinkExists(papLink) {
				fmt.Printf("PAPLINK FOUND: %s\n", string(papLink))
				urlFleet := fmt.Sprintf("https://adashboard.info/par/export/%s", papLink)
				bodyBytes, err, _ := obj.aDHttpGet(urlFleet)
				if err != nil {
					obj.Model.LogObj.Printf("decodePap ERROR fetching response. %s", err.Error())
					return
				}
				papsFound += obj.getPapsFromCSV(string(bodyBytes), papLink, alt2main)
			}
		}
	}
	if papsFound != 0 {
		charNames := util.GetSortKeysFromStrMap(obj.PapLogMap, false)
		for _, charName := range charNames {
			if obj.AddLogCB != nil {
				if obj.PapLogMap[charName] == 1 {
					obj.AddLogCB(fmt.Sprintf("[%s]%s added %d PAP", obj.CorpTicker, charName, obj.PapLogMap[charName]))
				} else {
					obj.AddLogCB(fmt.Sprintf("[%s]%s added %d PAPs", obj.CorpTicker, charName, obj.PapLogMap[charName]))
				}

			}
		}
		obj.PapLogMap = make(map[string]int)
	}
}

func (obj *ADashClient) getPapsFromCSV(htmlContent string, papLink string, alt2main map[string]string) int {
	var papsFound int
	re := regexp.MustCompile(`"(.*?)","(.*?)","(.*?)","(.*?)","(.*?)","(.*?)","(.*?)"`)
	scanner := bufio.NewScanner(strings.NewReader(htmlContent))
	for scanner.Scan() {
		result := re.FindStringSubmatch(scanner.Text())
		if result != nil {
			var pap PapElem
			pap.CorpId = obj.CorpID
			pap.PapLink = papLink
			pap.ChName = result[1]
			pap.CoTicker = result[2]
			pap.AlShort = result[3]
			pap.ShTypeName = result[4]
			pap.Sub = result[5]
			pap.Loc = result[6]
			pap.Timestamp = result[7]
			//fmt.Printf("%s %s\n", papLink, pap.ChName)
			newPap := obj.ConvertAdashPap2DB(&pap)
			result := obj.Model.AddCorpADashEntry(newPap)
			if result != model.DBR_Inserted {
				// skipped
			} else {
				if main, ok := alt2main[pap.ChName]; ok {
					obj.PapLogMap[main]++
				} else {
					obj.PapLogMap[pap.ChName]++
				}
				papsFound++
			}
		}
	}
	return papsFound
}

func (obj *ADashClient) ConvertAdashPap2DB(pap *PapElem) *model.DBpap {
	var newPap model.DBpap
	newPap.CorpId = pap.CorpId
	newPap.PapLink = pap.PapLink
	newPap.ChName = obj.Model.AddStringEntry(pap.ChName)
	newPap.CoTicker = obj.Model.AddStringEntry(pap.CoTicker)
	newPap.AlShort = obj.Model.AddStringEntry(pap.AlShort)
	newPap.ShTypeName = obj.Model.AddStringEntry(pap.ShTypeName)
	newPap.Sub = obj.Model.AddStringEntry(pap.Sub)
	newPap.Loc = obj.Model.AddStringEntry(pap.Loc)
	newPap.Timestamp = obj.ConvertPapTimeStrToInt(pap.Timestamp)
	return &newPap
}

func (obj *ADashClient) ConvertPapTimeStrToInt(timeString string) int64 {
	var retval int64
	if timeString != "" {
		t, err := time.Parse("2006-01-02 15:04:05.0", timeString)
		if err != nil {
			fmt.Println(err)
		}
		retval = t.Unix()
	}
	return retval
}

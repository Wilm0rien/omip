package ctrl

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type httpSrvData struct {
	HTTPServerPtr *http.Server
	SrvExitDone   *sync.WaitGroup
	cancelChan    chan struct{}
}

type UpdateFlags struct {
	PapLinks     bool
	Contracts    bool
	Corpmembers  bool
	IndustryJobs bool
	Journal      bool
	Killmails    bool
	Structures   bool
	Wallet       bool
}

type EsiData struct {
	EsiCharList []*EsiChar
	EsiCorpList []*EsiCorp
	SecretCode  []byte
	NVConfig    NVConfigData
	ETags       map[string]string // map[url]=etag
}

type EsiChar struct {
	InitAuth          AuthResponse
	RefreshAuthData   AuthResponse
	CharInfoData      CharacterInfo
	CharInfoExt       CharacterInfoExt
	stateMagicNum     uint32
	ImageFile         string
	NextAuthTimeStamp int64
	UpdateFlags       UpdateFlags
	KmSkipList        map[int32]bool // map[KillmailID]bool
}
type EsiCorp struct {
	Name          string
	CooperationId int
	AllianceId    int
	Ticker        string
	ImageFile     string
	UpdateFlags   UpdateFlags
	KmSkipList    map[int32]bool // map[KillmailID]bool
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type CharacterInfo struct {
	CharacterID          int
	CharacterName        string
	ExpiresOn            string
	Scopes               string
	TokenType            string
	CharacterOwnerHash   string
	IntellectualProperty string
}

type CharacterInfoExt struct {
	AllianceId    int `json:"alliance_id"`
	CooperationId int `json:"corporation_id"`
	Director      bool
}

type ServerStatus struct {
	Players       int    `json:"players"`
	ServerVersion int    `json:"server_version"`
	StartTime     string `json:"start_time"`
	Vip           bool   `json:"vip"`
}

var NewChar *EsiChar

const (
	scopes   = "publicData esi-wallet.read_character_wallet.v1 esi-universe.read_structures.v1 esi-killmails.read_killmails.v1 esi-corporations.read_corporation_membership.v1 esi-corporations.read_structures.v1 esi-industry.read_character_jobs.v1 esi-markets.read_character_orders.v1 esi-characters.read_corporation_roles.v1 esi-contracts.read_character_contracts.v1 esi-killmails.read_corporation_killmails.v1 esi-wallet.read_corporation_wallets.v1 esi-characters.read_notifications.v1 esi-contracts.read_corporation_contracts.v1 esi-industry.read_corporation_jobs.v1 esi-markets.read_corporation_orders.v1"
	clientID = "41b2d654515d40b5a04e727a334c6358"
	callBack = "http://localhost:4716/callback"
	// id ranges https://gist.github.com/a-tal/5ff5199fdbeb745b77cb633b7f4400bb
	EsiCorpIdLimit = 98000000
)

func (obj *Ctrl) OpenAuthInBrowser() {
	var newChar EsiChar
	newChar.stateMagicNum = uint32(time.Now().Unix())
	NewChar = &newChar
	// register app at https://developers.eveonline.com/
	// https://docs.esi.evetech.net/docs/sso/native_sso_flow.html
	obj.Esi.SecretCode = []byte(base64.URLEncoding.EncodeToString([]byte(util.GenerateRandomString(48))))
	h := sha256.New()
	h.Write(obj.Esi.SecretCode)
	code_challenge := base64.URLEncoding.EncodeToString(h.Sum(nil))
	code_challenge = strings.Replace(code_challenge, "=", "", -1)
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("redirect_uri", callBack)
	params.Add("client_id", clientID)
	params.Add("scope", scopes)
	params.Add("code_challenge", code_challenge)
	params.Add("code_challenge_method", "S256")
	params.Add("state", fmt.Sprintf("%d", newChar.stateMagicNum))

	url := fmt.Sprintf(
		"https://login.eveonline.com/v2/oauth/authorize?%s",
		params.Encode())
	util.OpenUrl(url)
}

func (obj *Ctrl) StartServer() {
	obj.Svr.SrvExitDone = &sync.WaitGroup{}
	obj.Svr.SrvExitDone.Add(1)
	srv := &http.Server{Addr: ":4716"}
	obj.Svr.HTTPServerPtr = srv
	http.HandleFunc("/", obj.webSrv)
	obj.Svr.cancelChan = make(chan struct{})
	go func() {
		defer obj.Svr.SrvExitDone.Done() // let main know we are done cleaning up
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}
func (obj *Ctrl) HTTPShutdown() {
	if err := obj.Svr.HTTPServerPtr.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	obj.Svr.SrvExitDone.Wait()

	if !obj.ServerCancelled() {
		close(obj.Svr.cancelChan)
	}
}

func (obj *Ctrl) ServerCancelled() bool {
	select {
	case <-obj.Svr.cancelChan:
		return true
	default:
		return false
	}
}

func (obj *Ctrl) webSrv(w http.ResponseWriter, r *http.Request) {
	code := ""
	var receivedMagic uint32
	receivedMagic = 0

	param1s := r.URL.Query()["code"]
	if len(param1s) > 0 {

		code = param1s[0]
	}
	param2int := r.URL.Query()["state"]
	if len(param2int) > 0 {

		stringConvert, err := strconv.Atoi(param2int[0])
		if err != nil {
			fmt.Printf("ERROR converting string to number \"%s\"\n", param2int[0])
		} else {
			receivedMagic = uint32(stringConvert)
		}
	}
	header := "<html><head><style>body {  background-color: black;   color: white;}</style></head><body>\n"
	footer := "</body></html>\n"
	fmt.Fprintf(w, "%s", header)

	if code == "shutdown" {
		fmt.Fprintf(w, "shutdown<br>")
		go func() {
			time.Sleep(200 * time.Millisecond)
			os.Exit(0)
		}()
	}
	if NewChar != nil {
		if receivedMagic != 0 {
			if NewChar.stateMagicNum == receivedMagic {
				var lNewChar EsiChar
				lNewChar = *NewChar
				NewChar = nil
				fmt.Fprintf(w, "OMIP AUTH code accepted!<br>you can close this window now!<br>")
				go obj.initialAuth(code, &lNewChar)
			} else {
				fmt.Fprintf(w, "OMIP ERROR unexpected magic number received!<br>")
			}
		} else {
			fmt.Fprintf(w, "OMIP ERROR unexpected magic number received (0)!<br>")
		}
	}
	fmt.Fprintf(w, "%s", footer)
}

func (obj *Ctrl) GetCharInfo(char *EsiChar) {
	tokenString := char.InitAuth.AccessToken
	token2, _ := jwt.Parse(tokenString, nil)
	if token2 != nil {
		claims, _ := token2.Claims.(jwt.MapClaims)
		if name, ok := claims["name"]; ok {
			char.CharInfoData.CharacterName = name.(string)
		}
		if sub, ok := claims["sub"]; ok {
			re := regexp.MustCompile(`^CHARACTER:EVE:([0-9]+)$`)
			substring, ok2 := sub.(string)
			if ok2 {
				result := re.FindStringSubmatch(substring)
				if len(result) > 1 {
					charID, err := strconv.Atoi(result[1])
					if err == nil {
						char.CharInfoData.CharacterID = charID
					}
				}
			}
		}
		if char.CharInfoExt.CooperationId == 0 && char.CharInfoData.CharacterID != 0 {
			obj.GetCharInfoExt(char)
		}
	}
	util.Assert(char.CharInfoData.CharacterID != 0)
}

func (obj *Ctrl) GetCharInfoExt(char *EsiChar) {
	url := fmt.Sprintf("https://esi.evetech.net/v5/characters/%d",
		char.CharInfoData.CharacterID)
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	if bodyBytes != nil {

		err := json.Unmarshal(bodyBytes, &char.CharInfoExt)
		if err != nil {
			obj.AddLogEntry(fmt.Sprintf(err.Error()))
		}
	}
}

func (obj *Ctrl) initialAuth(token string, char *EsiChar) {
	body := fmt.Sprintf("grant_type=authorization_code&code=%s&client_id=%s&code_verifier=%s", token, clientID, obj.Esi.SecretCode)

	auth := obj.doAuthRequest(body)
	if auth != nil {
		char.InitAuth = *auth
		char.RefreshAuthData = *auth

		// claims are actually a map[string]interface{}
		// ESI BUG! if the time between the auth request and the refresh is too short the request will fail!
		time.Sleep(500 * time.Millisecond)
		util.Assert(len(char.InitAuth.AccessToken) != 0)

		obj.GetCharInfo(char)

		imgFile := fmt.Sprintf("%d_128.jpg", char.CharInfoData.CharacterID)
		imgPath := fmt.Sprintf("%s/%s", obj.Model.LocalImgDir, imgFile)
		if !util.Exists(imgPath) {
			fullUrlFile := "https://imageserver.eveonline.com/Character/" + imgFile
			util.GetImgFromUrl(fullUrlFile, imgPath)
		}
		char.ImageFile = imgPath
		obj.setUpdateFlags(char)
		if char.CharInfoData.CharacterID != 0 {
			NewChar = nil
			obj.InitiateKMSkipList(char, false)
			if obj.CheckIfDirector(char) {
				char.CharInfoExt.Director = true
				if !obj.corpExists(char) {
					newCorp := obj.getAuthCorpInfo(char)
					obj.Esi.EsiCorpList = append(obj.Esi.EsiCorpList, newCorp)
				}
				obj.InitiateKMSkipList(char, true)
			}
			obj.Esi.EsiCharList = append(obj.Esi.EsiCharList, char)
			if obj.AuthCb != nil {
				obj.AuthCb(char)
			}
		}
	}
}
func (obj *Ctrl) setUpdateFlags(char *EsiChar) {
	char.UpdateFlags.PapLinks = true
	char.UpdateFlags.Contracts = true
	char.UpdateFlags.Corpmembers = true
	char.UpdateFlags.IndustryJobs = true
	char.UpdateFlags.Journal = true
	char.UpdateFlags.Killmails = true
	char.UpdateFlags.Structures = true
	char.UpdateFlags.Wallet = true
}

func (obj *Ctrl) corpExists(char *EsiChar) bool {
	var found bool
	for _, corp := range obj.Esi.EsiCorpList {
		if corp.CooperationId == char.CharInfoExt.CooperationId {
			found = true
			break
		}
	}
	return found
}

func (obj *Ctrl) CheckIfDirector(char *EsiChar) bool {
	var retval bool
	url := fmt.Sprintf("https://esi.evetech.net/v2/corporations/%d/roles/?datasource=tranquility", char.CharInfoExt.CooperationId)
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	if bodyBytes != nil {
		retval = true
	}
	return retval
}

func (obj *Ctrl) CheckServerUp(char *EsiChar) (retval bool) {
	url := fmt.Sprintf("https://esi.evetech.net/v2/status/?datasource=tranquility")
	bodyBytes, _ := obj.getSecuredUrl(url, char)
	if bodyBytes != nil {
		var serverStatus ServerStatus
		err := json.Unmarshal(bodyBytes, &serverStatus)
		if err != nil {
			if serverStatus.Players == 0 {
				retval = false
			} else {
				retval = true
			}
		}
	}
	return
}

func (obj *Ctrl) RefreshAuth(char *EsiChar, enforce bool) {
	if time.Now().Unix() >= char.NextAuthTimeStamp || enforce {
		URLEncodedToken := url.QueryEscape(char.InitAuth.RefreshToken)
		body2 := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s",
			URLEncodedToken, clientID)
		var auth *AuthResponse
		auth = obj.doAuthRequest(body2)
		if auth != nil {
			char.RefreshAuthData = *auth
		}
		char.NextAuthTimeStamp =
			int64(time.Now().Unix()) + int64(char.RefreshAuthData.ExpiresIn-1)

	}
}

func (obj *Ctrl) doAuthRequest(body string) *AuthResponse {
	url := "https://login.eveonline.com/v2/oauth/token"
	req, err1 := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err1 != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR accessing login.eveonline.com %s", err1.Error()))
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "login.eveonline.com")

	bodyBytes, clientErr, resp := obj.httpClientRequest(req)
	if clientErr != nil {
		obj.AddLogEntry(fmt.Sprintf("ERROR accessing login.eveonline.com %s", clientErr.Error()))
		return nil
	}
	var retval *AuthResponse
	var authVal AuthResponse
	retval = nil
	if resp.StatusCode == http.StatusOK {
		err := json.Unmarshal(bodyBytes, &authVal)
		if err != nil {
			obj.AddLogEntry(fmt.Sprintf(err.Error()))
		} else {
			retval = &authVal
		}
	} else {
		obj.AddLogEntry(fmt.Sprintf("ERROR %s %s", resp.Status, string(bodyBytes)))
	}

	return retval
}

func (obj *Ctrl) getSecuredUrl(url string, char *EsiChar) (bodyBytes []byte, Xpages int) {
	obj.RefreshAuth(char, false)
	etagTrigger := false
	if len(char.RefreshAuthData.AccessToken) == 0 {
		obj.AddLogEntry("ERROR  no initial auth saved")
		return nil, 0
	} else {
		req, err1 := http.NewRequest("GET", url, nil)
		if err1 != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR %s", err1.Error()))
			return nil, 0
		}

		req.Header.Add("User-Agent", "Contact: Wilm0rien in game or on devfleet slack")
		req.Header.Add("Authorization", "Bearer "+char.RefreshAuthData.AccessToken)
		req.Header.Add("Host", "login.eveonline.com")
		oldEtag := ""
		if val, ok := obj.Esi.ETags[url]; ok {
			oldEtag = val
			bodyBytes = obj.Model.LoadEtag(oldEtag)
			if bodyBytes != nil {
				req.Header.Add("If-None-Match", oldEtag)
			}

		}

		var requestOK bool
		var noError bool
		var retrycounter int

		for !requestOK && retrycounter < 3 {
			bodyBytes2, clientErr, resp := obj.httpClientRequest(req)
			if resp.StatusCode == 304 { // https://developers.eveonline.com/blog/article/esi-etag-best-practices
				bodyBytes = obj.Model.LoadEtag(oldEtag)
				log.Printf("RESP (cached): %d bytes\n", len(bodyBytes))
				requestOK = true
				etagTrigger = true
			} else {
				if clientErr == nil {
					responseStr := string(bodyBytes2)
					matched, _ := regexp.MatchString(`"error":`, responseStr)
					if matched {
						matched2, _ := regexp.MatchString(`"error":"token is expired",`, responseStr)
						matched3, _ := regexp.MatchString(`"error":"Character does not have required role`, responseStr)
						matched4, _ := regexp.MatchString(`"error":"Unhandled internal error encountered!"`, responseStr)
						matched5, _ := regexp.MatchString(`"error":"ConStopSpamming`, responseStr)
						matched6, _ := regexp.MatchString(`"error":"The given character doesn't have the required role`, responseStr)
						if matched2 {
							obj.RefreshAuth(char, true)
						} else if matched5 {
							time.Sleep(800 * time.Millisecond)
						} else if matched3 || matched4 || matched6 {
							noError = true
							break
						}
					} else {
						if val, ok := resp.Header["X-Pages"]; ok {
							if len(resp.Header["X-Pages"]) == 1 {
								Xpages, _ = strconv.Atoi(val[0])
							}
						}

						if len(bodyBytes2) > 0 {
							requestOK = true
							if newEtag, ok := resp.Header["Etag"]; ok {
								re := regexp.MustCompile(`^"(.*)"$`)
								result := re.FindStringSubmatch(newEtag[0])
								if len(result) > 1 {
									obj.Esi.ETags[url] = result[1]
									obj.Model.StoreEtag(result[1], oldEtag, bodyBytes2)
								}

							}
						}
						bodyBytes = bodyBytes2
					}
				}
				if !requestOK {
					time.Sleep(200 * time.Millisecond)
				}
			}
			retrycounter++
		}

		if !requestOK {
			if !noError {
				matched, _ := regexp.MatchString(`roles`, url)
				if !matched {
					obj.AddLogEntry(fmt.Sprintf("URL FAILED: %s", url))
				}
			}
			bodyBytes = nil
		}

	}
	if !CtrlTestEnable && !etagTrigger {
		time.Sleep(200 * time.Millisecond)
	}

	return bodyBytes, Xpages
}

func (obj *Ctrl) getSecuredUrlPost(url string, body string, char *EsiChar) (bodyBytes []byte, resp *http.Response) {
	obj.RefreshAuth(char, false)
	if len(char.RefreshAuthData.AccessToken) == 0 {
		obj.AddLogEntry("no initial auth saved")
	} else {
		req, err1 := http.NewRequest("POST", url, bytes.NewBufferString(body))
		if err1 != nil {
			obj.AddLogEntry(err1.Error())
		}
		req.Header.Add("User-Agent", "Contact: Wilm0rien in game or on devfleet slack")
		req.Header.Add("Authorization", "Bearer "+char.RefreshAuthData.AccessToken)
		req.Header.Add("Host", "login.eveonline.com")
		var bodyBytes2 []byte
		var clientErr error
		bodyBytes2, clientErr, resp = obj.httpClientRequest(req)
		if clientErr == nil {
			bodyBytes = bodyBytes2
		} else {
			obj.AddLogEntry(clientErr.Error())
		}

	}
	return bodyBytes, resp
}

var HttpRequestMock func(req *http.Request) (bodyBytes []byte, err error, resp *http.Response)

func (obj *Ctrl) httpClientRequest(req *http.Request) (bodyBytes []byte, err error, resp *http.Response) {
	if !CtrlTestEnable {
		client := &http.Client{}
		log.Printf("REQ:\n%s\n", req.URL)
		resp, err = client.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				bodyBytes, _ = ioutil.ReadAll(resp.Body)
				log.Printf("RESP:\n%s\n", string(bodyBytes))
				resp.Body.Close()
			}
		}
	} else {
		if HttpRequestMock != nil {
			bodyBytes, err, resp = HttpRequestMock(req)
		}
	}

	return bodyBytes, err, resp
}

func (obj *Ctrl) ZkillOk(kmId int) (retval bool) {
	url := fmt.Sprintf("https://zkillboard.com/kill/%d/", kmId)

	req, err1 := http.NewRequest("GET", url, nil)

	if err1 != nil {
		return false
	}
	req.Header.Add("User-Agent", "Contact: Wilm0rien in game or on devfleet slack")
	bodyBytes, clientErr, _ := obj.httpClientRequest(req)
	if clientErr == nil && bodyBytes != nil {
		responseStr := string(bodyBytes)
		scanner := bufio.NewScanner(strings.NewReader(responseStr))
		retval = true
		for scanner.Scan() {
			curLine := scanner.Text()
			matched, _ := regexp.MatchString(`The content you're after isn't here`, curLine)
			if matched {
				retval = false
				break
			}
		}
	}
	return retval
}

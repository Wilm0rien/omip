package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/denisbrodbeck/machineid"
	"image/color"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

const (
	OmipSoftwareVersion = "1.0.6"
)

func Assert(value bool) {
	if !value {
		panic("assert")
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateDirectory(dirName string) bool {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			return false
		}
		return true
	}

	if src.Mode().IsRegular() {
		return false
	}

	return false
}

func GetImgFromUrl(srcUrl string, dstPath string) {
	url := srcUrl
	response, e := http.Get(url)
	if e != nil {
		log.Printf(e.Error())
	} else {
		defer response.Body.Close()
	}

	//open a file for writing
	file, err := os.Create(dstPath)
	if err != nil {
		log.Printf(err.Error())
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Printf(err.Error())
	}
	//log.Printf("Success! Downloading %s\n", srcUrl)
}

func Get32BitMd5FromString(input string) int32 {
	h := md5.New()
	io.WriteString(h, input)
	md5value := h.Sum(nil)
	md5int := int32(md5value[0])<<0 + int32(md5value[1])<<8 + int32(md5value[2])<<16 + int32(md5value[3])<<24
	return int32(md5int)
}

func Get64BitMd5FromString(input string) int64 {
	h := md5.New()
	io.WriteString(h, input)
	md5value := h.Sum(nil)
	md5int := int64(md5value[0])<<0 + int64(md5value[1])<<8 + int64(md5value[2])<<16 + int64(md5value[3])<<24 + int64(md5value[4])<<32 + int64(md5value[5])<<40 + int64(md5value[6])<<48 + int64(md5value[7])<<56
	md5int &= 0x7FFFFFFFFFFFFFFF
	return int64(md5int)
}
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, _ := GenerateRandomBytes(n)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HumanizeNumber(input float64) string {
	var retval string

	testNoDigits := func(x float64) bool {
		xint := int(x * 100)
		return xint%100 == 0
	}
	printFactor := func(input float64, magnitude string) string {
		if testNoDigits(input) {
			retval = fmt.Sprintf("%d%s", int(input), magnitude)
		} else {
			retval = fmt.Sprintf("%3.2f%s", input, magnitude)
		}
		return retval
	}

	if math.Abs(input) >= 1000 {
		if math.Abs(input) >= 1000000 {
			if math.Abs(input) >= 1000000000 {
				retval = printFactor(input/1000000000, "b")
			} else {
				retval = printFactor(input/1000000, "m")
			}
		} else {
			retval = printFactor(input/1000, "k")
		}
	} else {
		retval = printFactor(input, "")
	}
	return retval
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func sha1func(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	byteSha1 := h.Sum(nil)
	return hex.EncodeToString(byteSha1)
}

func GenSysPassphrase() (result string) {
	result = "VzWsunp7vSTFyWpOrZR7"
	id, err := machineid.ID()
	if err == nil {
		result += id
	}
	result = sha1func(result)
	return result
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) (plaintext []byte, success bool) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		success = false
		plaintext = []byte{}
	} else {
		success = true
	}
	return plaintext, success
}

func UnixTS2DateStr(ts int64) string {
	tm := time.Unix(ts, 0)
	year, month, day := tm.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func UnixTS2YMStr(ts int64) string {
	tm := time.Unix(ts, 0)
	year, month, _ := tm.Date()
	return fmt.Sprintf("%04d-%02d", year, month)
}

func UnixTS2DateTimeStr(ts int64) string {
	tm := time.Unix(ts, 0)
	year, month, day := tm.Date()
	hour, minute, second := tm.Clock()
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, minute, second)
}
func UnixTS2AdashDateTimeStr(ts int64) string {
	tm := time.Unix(ts, 0)
	year, month, day := tm.Date()
	hour, minute, second := tm.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.0", year, month, day, hour, minute, second)
}

func GetSortKeysFromStrMap(x interface{}, reverse bool) []string {
	var keyList []string
	rValue := reflect.ValueOf(x)
	Assert(rValue.Kind() == reflect.Map)

	if len(rValue.MapKeys()) > 0 {
		Assert(rValue.MapKeys()[0].Kind() == reflect.String)
		for _, key := range rValue.MapKeys() {
			keyList = append(keyList, key.String())
		}
		if reverse {
			sort.Slice(keyList, func(i, j int) bool { return strings.ToLower(keyList[i]) >= strings.ToLower(keyList[j]) })
		} else {
			sort.Slice(keyList, func(i, j int) bool { return strings.ToLower(keyList[i]) < strings.ToLower(keyList[j]) })
		}
	}
	return keyList
}

func GetSortKeysFromIntMap(x interface{}, reverse bool) []int {
	var keyList []int
	rValue := reflect.ValueOf(x)
	Assert(rValue.Kind() == reflect.Map)
	if len(rValue.MapKeys()) > 0 {
		Assert(rValue.MapKeys()[0].Kind() == reflect.Int)
		for _, key := range rValue.MapKeys() {
			keyList = append(keyList, int(key.Int()))
		}
		if reverse {
			sort.Sort(sort.Reverse(sort.IntSlice(keyList)))
		} else {
			sort.IntSlice(keyList).Sort()
		}
	}
	return keyList
}

func GetSortKeysFromFloat64Map(x interface{}, reverse bool) []float64 {
	var keyList []float64
	rValue := reflect.ValueOf(x)
	Assert(rValue.Kind() == reflect.Map)
	if len(rValue.MapKeys()) > 0 {
		Assert(rValue.MapKeys()[0].Kind() == reflect.Int)
		for _, key := range rValue.MapKeys() {
			keyList = append(keyList, key.Float())
		}
		if reverse {
			sort.Sort(sort.Reverse(sort.Float64Slice(keyList)))
		} else {
			sort.Float64Slice(keyList).Sort()
		}
	}
	return keyList
}

func CheckTimeStampIsInThisMonth(timestamp int64) (result bool) {
	nowYear, nowMonth, _ := time.Now().Date()
	tm := time.Unix(timestamp, 0)
	year, month, _ := tm.Date()
	if year == nowYear && month == nowMonth {
		result = true
	}
	return result
}

func GetTimeDiffStringFromTS(timestamp int64) (string, bool) {
	now := int64(time.Now().Unix())
	var timeDiff int64
	retval := "undefined"
	future := true
	if now > timestamp {
		timeDiff = now - timestamp
		future = false
	} else {
		timeDiff = timestamp - now
	}
	if timeDiff != 0 {
		if timeDiff > 24*60*60 {
			retval = fmt.Sprintf("% 3dd% 3dh", timeDiff/(24*60*60), timeDiff%(24*60*60)/(60*60))
		} else if timeDiff > 60*60 {
			retval = fmt.Sprintf("% 3dh% 3dm", timeDiff/(60*60), timeDiff%(60*60)/60)
		} else if timeDiff > 60 {
			retval = fmt.Sprintf("% 3dm% 3ds", timeDiff/(60), timeDiff%(60))
		} else {
			retval = fmt.Sprintf("% 3ds", timeDiff)
		}
	}

	return retval, future
}

func GetTimeDiffStringFromDiff(timeDiff int64) string {
	retval := "undefined"
	if timeDiff != 0 {
		if timeDiff > 24*60*60 {
			retval = fmt.Sprintf("% 3dd% 3dh", timeDiff/(24*60*60), timeDiff%(24*60*60)/(60*60))
		} else if timeDiff > 60*60 {
			retval = fmt.Sprintf("% 3dh% 3dm", timeDiff/(60*60), timeDiff%(60*60)/60)
		} else if timeDiff > 60 {
			retval = fmt.Sprintf("% 3dm% 3ds", timeDiff/(60), timeDiff%(60))
		} else {
			retval = fmt.Sprintf("% 3ds", timeDiff)
		}
	}
	return retval
}

func TimeDiffStrStarEnd(start int64, end int64) (string, bool) {
	now := start
	var timeDiff int64
	retval := "undefined"
	future := true
	if now > end {
		timeDiff = now - end
		future = false
	} else {
		timeDiff = end - now
	}
	if timeDiff != 0 {
		if timeDiff > 24*60*60 {
			retval = fmt.Sprintf("% 3dd% 3dh", timeDiff/(24*60*60), timeDiff%(24*60*60)/(60*60))
		} else if timeDiff > 60*60 {
			retval = fmt.Sprintf("% 3dh% 3dm", timeDiff/(60*60), timeDiff%(60*60)/60)
		} else if timeDiff > 60 {
			retval = fmt.Sprintf("% 3dm% 3ds", timeDiff/(60), timeDiff%(60))
		} else {
			retval = fmt.Sprintf("% 3ds", timeDiff)
		}
	}

	return retval, future
}

func OpenUrl(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("ERROR OpenUrl %s", err.Error())
	}

}

func CheckErr(err error) {
	if err != nil {
		log.Printf("%v", err)
		msg := fmt.Sprintf("%v\n", err)
		msg += fmt.Sprintf(string(debug.Stack()))
		os.WriteFile("err_stack_trace.log", []byte(msg), 0644)
	}
}

func ConvertTimeStrToInt(timeString string) int64 {
	var retval int64
	if timeString != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", timeString)
		if err != nil {
			fmt.Println(err)
			log.Printf("ConvertTimeStrToInt ERROR PARSING TIME %s", timeString)
		} else {
			retval = t.Unix()
		}
	}
	return retval
}

func ConvertServerTimeStrToInt(timeString string) int64 {
	var retval int64
	if timeString != "" {
		re := regexp.MustCompile(`^[A-Za-z]+, ([0-9]+ [a-zA-Z]+ [0-9]+ [0-9]+:[0-9]+:[0-9]+) GMT$`)
		result := re.FindStringSubmatch(timeString)
		if len(result) > 1 {
			subString := fmt.Sprintf("%s", result[1])
			t, err := time.Parse("02 Jan 2006 15:04:05", subString)
			if err != nil {
				fmt.Println(err)
				log.Printf("ConvertTimeStrToInt ERROR PARSING TIME %s", timeString)
			} else {
				retval = t.Unix()
			}
		} else {
			log.Printf("ConvertTimeStrToInt ERROR PARSING TIME %s", timeString)
		}

	}
	return retval
}

func ConvertUnixTimeToStr(timeValue int64) string {
	tm := time.Unix(timeValue, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func ClipboardCopy() string {
	text, _ := clipboard.ReadAll()
	return text
}

func ClipboardPaste(paste string) {
	clipboard.WriteAll(paste)
}
func GetColor(max float64, val float64, reverse bool) *color.NRGBA {
	ratio := val / max
	var red uint8
	var green uint8
	if reverse == false {
		if ratio <= 0.5 {
			red = 255
			green = uint8(255 * ratio * 2)
		} else {
			green = 255
			red = uint8(255 * (1 - ratio) * 2)
		}
	} else {
		if ratio <= 0.5 {
			green = 255
			red = uint8(255 * ratio * 2)

		} else {
			red = 255
			green = uint8(255 * (1 - ratio) * 2)
		}
	}
	return &color.NRGBA{R: red, G: green, B: 32, A: 255}
}

func GetAppDataDir() (appData string) {
	switch runtime.GOOS {
	case "linux":
		appData = os.Getenv("HOME")
	case "windows":
		appData = strings.Replace(os.Getenv("appdata"), "\\", "/", -1)
	}
	return appData
}

func SendReq(urlStr string) (retval bool) {
	client := &http.Client{}

	req, err1 := http.NewRequest("GET", urlStr, nil)
	if err1 != nil {
		log.Printf("failed creating request %s", err1.Error())
	}
	client.Timeout = 1000 * time.Millisecond
	resp, err2 := client.Do(req)

	if err2 != nil {
		//t.Errorf("failed sending request %s Error %s", urlStr, err1.Error())
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			retval = true
		}
	}
	return retval
}

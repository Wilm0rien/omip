package update

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Updater struct {
}

func NewUpdaterObj() (result *Updater) {
	var obj Updater
	return &obj
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
	Max   uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)

	fmt.Printf("\rDownloading... %3.2f %s complete", float64(wc.Total)/float64(wc.Max)*100, util.HumanizeNumber(float64(wc.Total)))
}

func (obj *Updater) DownloadFile(filepath string, url string, fileSize uint64) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{
		Max: fileSize,
	}

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func (obj *Updater) ExtractExec(zipfile string, destinationDir string) (err error) {
	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		return errors.New(fmt.Sprintf("OpenReader %s", err.Error()))
	}
	defer reader.Close()
	destinationDir, err = filepath.Abs(destinationDir)
	if err != nil {
		return errors.New(fmt.Sprintf("filepath.Abs %s", err.Error()))
	}
	for _, f := range reader.File {
		err2 := obj.unzipFile(f, destinationDir)
		if err2 != nil {
			return errors.New(fmt.Sprintf("unzipFile %s f %s destination %s", err2.Error(), f.Name, destinationDir))
		}
		break
	}
	return
}

func (obj *Updater) Copy(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (obj *Updater) unzipFile(f *zip.File, destinationDir string) error {
	filePath := filepath.Join(destinationDir, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destinationDir)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err3 := io.Copy(destinationFile, zippedFile); err != nil {
		return err3
	}
	return nil
}

type GitRelease struct {
	TagName   string `json:"tag_name"`
	Name      string `json:"name"`
	AssetsUrl string `json:"assets_url"`
	TimeStamp string `json:"published_at"`
}

type GitAssets struct {
	Url       string `json:"browser_download_url"`
	TimeStamp string `json:"updated_at"`
	FileSize  uint64 `json:"size"`
	TagName   string
}

func GetRelease(url string, pattern string) (asset *GitAssets, tag string) {
	releaseBytes := getUrl(url)
	var releaseList []GitRelease
	contentError := json.Unmarshal(releaseBytes, &releaseList)
	if contentError != nil {
		fmt.Printf("ERROR reading releases %s", contentError.Error())
		return
	}
	if len(releaseList) > 0 {
		latestRelease := getLatestRelease(releaseList)
		assetsBytes := getUrl(latestRelease.AssetsUrl)
		tag = latestRelease.TagName
		if len(assetsBytes) > 2 {
			assetList := handleAssets(assetsBytes)
			for _, elem := range assetList {
				matched, _ := regexp.MatchString(pattern, elem.Url)
				if matched {
					asset = &elem
					break
				}
			}
		}
	}
	return
}
func getLatestRelease(releaseList []GitRelease) (result GitRelease) {
	var latestTime int64
	for _, release := range releaseList {
		releaseTime := util.ConvertTimeStrToInt(release.TimeStamp)
		if releaseTime > latestTime {
			latestTime = releaseTime
			result = release
		}
	}
	return
}

func handleAssets(assetsBytes []byte) []GitAssets {
	var assetList []GitAssets
	contentError := json.Unmarshal(assetsBytes, &assetList)
	if contentError != nil {
		fmt.Printf("ERROR reading assets %s", contentError.Error())
		return nil
	}

	return assetList
}

func getUrl(url string) (result []byte) {
	client := &http.Client{}
	req, err1 := http.NewRequest("GET", url, nil)
	if err1 != nil {
		fmt.Printf("%s", err1.Error())
	}
	req.Header.Add("User-Agent", "https://github.com/Wilm0rien/omip")
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	//log.Printf("REQ:\n%s\n", req.URL)
	resp, err3 := client.Do(req)

	if err3 == nil {
		if resp.StatusCode == http.StatusOK {
			result, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	} else {
		fmt.Printf("%s", err1.Error())
	}
	return

}

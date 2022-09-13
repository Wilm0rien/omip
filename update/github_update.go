package update

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"io"
	"net/http"
	"regexp"
)

type GitRelease struct {
	TagName   string `json:"tag_name"`
	Name      string `json:"name"`
	AssetsUlr string `json:"assets_url"`
	TimeStamp string `json:"published_at"`
}

type GitAssets struct {
	Url       string `json:"browser_download_url"`
	TimeStamp string `json:"updated_at"`
	FileSize  uint64 `json:"size"`
	TagName   string
}

func GetRelease(url string) (asset *GitAssets) {
	releaseBytes := getUrl(url)
	var releaseList []GitRelease
	contentError := json.Unmarshal(releaseBytes, &releaseList)
	if contentError != nil {
		fmt.Printf("ERROR reading releases %s", contentError.Error())
		return
	}
	if len(releaseList) > 0 {
		latestRelease := getLatestRelease(releaseList)
		assetsBytes := getUrl(latestRelease.AssetsUlr)
		if len(assetsBytes) > 2 {
			asset = handleAssets(assetsBytes, latestRelease)
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

func handleAssets(assetsBytes []byte, release GitRelease) (result *GitAssets) {
	var assetList []GitAssets
	contentError := json.Unmarshal(assetsBytes, &assetList)
	if contentError != nil {
		fmt.Printf("ERROR reading assets %s", contentError.Error())
		return
	}
	if len(assetList) > 0 {
		for _, asset := range assetList {

			matched, _ := regexp.MatchString(`omip\.zip$`, asset.Url)
			if matched {
				fmt.Printf("%s %s\n", asset.Url, asset.TimeStamp)
				result = &asset
				result.TagName = release.TagName
				break
			}
		}
	}
	return
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

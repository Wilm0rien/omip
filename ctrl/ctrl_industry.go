package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/Wilm0rien/omip/model"
	"github.com/Wilm0rien/omip/util"
	"time"
)

const (
	maxJobItemsCorp int   = 1000
	invalidTimeDiff int64 = 0x0FFFFFFFFFFFFFFF
)

type JobInfos struct {
	ActivityId           int     `json:"activity_id"`
	BlueprintId          int64   `json:"blueprint_id"`
	BlueprintLocationId  int64   `json:"blueprint_location_id"`
	BlueprintTypeId      int     `json:"blueprint_type_id"`
	CompletedCharacterId int     `json:"completed_character_id"`
	CompletedDate        string  `json:"completed_date"`
	Cost                 float64 `json:"cost"`
	Duration             int     `json:"duration"`
	EndDate              string  `json:"end_date"`
	FacilityId           int64   `json:"facility_id"`
	InstallerId          int     `json:"installer_id"`
	JobId                int     `json:"job_id"`
	LicensedRuns         int     `json:"licensed_runs"`
	OutputLocationId     int64   `json:"output_location_id"`
	PauseDate            string  `json:"pause_date"`
	Probability          float64 `json:"probability"`
	ProductTypeId        int     `json:"product_type_id"`
	Runs                 int     `json:"runs"`
	StartDate            string  `json:"start_date"`
	StationId            int64   `json:"station_id"`
	Status               string  `json:"status"`
	SuccessfulRuns       int     `json:"successful_runs"`
}

func (obj *Ctrl) UpdateIndustry(char *EsiChar, corp bool) {
	if !char.UpdateFlags.IndustryJobs {
		return
	}
	var url string
	pageID := 1
	statusMap := make(map[string]int64)
	statusMap["requireAttention"] = 0
	statusMap["endMinSeconds"] = invalidTimeDiff
	for {
		if corp {
			url = fmt.Sprintf("https://esi.evetech.net/v1/corporations/%d/industry/jobs/?datasource=tranquility&page=%d", char.CharInfoExt.CooperationId, pageID)
		} else {
			url = fmt.Sprintf("https://esi.evetech.net/v1/characters/%d/industry/jobs/?datasource=tranquility", char.CharInfoData.CharacterID)
		}
		bodyBytes, Xpages, _ := obj.getSecuredUrl(url, char)
		var jobList []JobInfos
		err := json.Unmarshal(bodyBytes, &jobList)
		if err != nil {
			obj.AddLogEntry(fmt.Sprintf("ERROR reading url %s", url))
		} else {
			obj.UpdateIndustryInDb(char, corp, jobList, statusMap)
		}
		if corp && pageID < Xpages {
			time.Sleep(100 * time.Millisecond)
			pageID++
		} else {
			break
		}
		util.Assert(pageID < 30)
	}
	obj.SummaryLogEntry("industry job", char, corp, statusMap)
}

func (obj *Ctrl) UpdateIndustryInDb(char *EsiChar, corp bool, jobList []JobInfos, statusMap map[string]int64) {
	unixTime := time.Now().Unix()
	jobExistMap := make(map[int]int)
	for _, job := range jobList {
		var isCorp int
		if corp {
			isCorp = 1
		}
		jobExistMap[job.JobId] = 1
		dbJob := obj.convertEsiJob2DB(&job, char.CharInfoData.CharacterID, char.CharInfoExt.CooperationId, isCorp)
		diff := dbJob.EndDate - unixTime
		if diff <= 0 {
			dbJob.Status = model.JobStatus["ready"]
		}
		result := obj.Model.AddJobEntry(dbJob)
		if obj.Model.GetStructureName(dbJob.FacilityId) == nil {
			if dbJob.FacilityId < 100000000 {
				//obj.Model.LogObj.Printf("invalid station id %d for jobid %d", dbJob.StationId, dbJob.JobId)
			} else {
				obj.GetStructureNameFromEsi(char, dbJob.FacilityId)
			}
		}
		if result == model.DBR_Updated {
			statusMap[job.Status]++
		}

		if diff <= 0 {
			statusMap["requireAttention"]++
		} else {
			if diff < statusMap["endMinSeconds"] {
				statusMap["endMinSeconds"] = diff
			}
		}

	}
	var dbList []*model.DBJob
	if corp {
		dbList = obj.Model.GetIndustryJobs(char.CharInfoExt.CooperationId, corp)
	} else {
		dbList = obj.Model.GetIndustryJobs(char.CharInfoData.CharacterID, corp)
	}

	for _, job := range dbList {
		if _, ok := jobExistMap[job.JobId]; !ok {
			// job not in esi
			job.Status = model.Job_Stat_cancelled
			obj.Model.AddJobEntry(job)
		}
	}
}

func (obj *Ctrl) SummaryLogEntry(typeStr string, char *EsiChar, corp bool, statusMap map[string]int64) {
	var logEntry string
	var nameStr string
	//jobEndMinSeconds
	if corp {
		corpObj := obj.GetCorp(char)
		if corpObj != nil {
			nameStr = corpObj.Name
		} else {
			nameStr = "N/A"
		}
	} else {
		nameStr = char.CharInfoData.CharacterName
	}
	if statusMap["requireAttention"] != 0 {
		logEntry += fmt.Sprintf("%s require attention. ", typeStr)
	}
	if len(statusMap) > 0 {
		for status, count := range statusMap {
			if status != "requireAttention" && status != "endMinSeconds" {
				logEntry += fmt.Sprintf("%s %s:%d ", typeStr, status, count)
			}
		}
	}
	if statusMap["endMinSeconds"] != invalidTimeDiff {
		var days float64
		var dur2 time.Duration
		dur := time.Duration(statusMap["endMinSeconds"]) * time.Second
		if dur > 24*time.Hour {
			days = dur.Hours() / 24
			dur2 = dur - (24 * time.Duration(days) * time.Hour)
		}
		if days == 0 {
			logEntry += fmt.Sprintf("next %s ends in %s", typeStr,
				dur.String())
		} else if days < 4 {
			logEntry += fmt.Sprintf("next %s ends in %d days %s", typeStr,
				int(days), dur2.String())
		}
	}
	if logEntry != "" {
		namePrefix := fmt.Sprintf("%s: ", nameStr)
		obj.AddLogEntry(namePrefix + logEntry)
	}
}

func (obj *Ctrl) convertEsiJob2DB(esiJob *JobInfos, charId int, corpId int, isCorp int) *model.DBJob {
	var newJob model.DBJob
	newJob.CharId = charId
	newJob.CorpId = corpId
	newJob.IsCorp = isCorp
	newJob.ActivityId = esiJob.ActivityId
	newJob.BlueprintId = esiJob.BlueprintId
	newJob.BlueprintLocationId = esiJob.BlueprintLocationId
	newJob.BlueprintTypeId = esiJob.BlueprintTypeId
	newJob.CompletedCharacterId = esiJob.CompletedCharacterId
	newJob.CompletedDate = util.ConvertTimeStrToInt(esiJob.CompletedDate)
	newJob.Cost = esiJob.Cost
	newJob.Duration = esiJob.Duration
	newJob.EndDate = util.ConvertTimeStrToInt(esiJob.EndDate)
	newJob.FacilityId = esiJob.FacilityId
	newJob.InstallerId = esiJob.InstallerId
	newJob.JobId = esiJob.JobId
	newJob.LicensedRuns = esiJob.LicensedRuns
	newJob.OutputLocationId = esiJob.OutputLocationId
	newJob.PauseDate = util.ConvertTimeStrToInt(esiJob.PauseDate)
	newJob.Probability = esiJob.Probability
	newJob.ProductTypeId = esiJob.ProductTypeId
	newJob.Runs = esiJob.Runs
	newJob.StartDate = util.ConvertTimeStrToInt(esiJob.StartDate)
	newJob.StationId = esiJob.StationId
	newJob.Status = model.JobStatus[esiJob.Status]
	newJob.SuccessfulRuns = esiJob.SuccessfulRuns
	return &newJob
}

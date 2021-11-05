package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
	"time"
)

const (
	Job_Stat_active    = 0
	Job_Stat_cancelled = 1
	Job_Stat_delivered = 2
	Job_Stat_paused    = 3
	Job_Stat_ready     = 4
	Job_Stat_reverted  = 5
)

var JobStatus = map[string]int{
	"active":    Job_Stat_active,
	"cancelled": Job_Stat_cancelled,
	"delivered": Job_Stat_delivered,
	"paused":    Job_Stat_paused,
	"ready":     Job_Stat_ready,
	"reverted":  Job_Stat_reverted,
}

func (obj *Model) JobStatusId2Str(statusID int) (jobStatus string) {
	jobStatus = "unknownStatus"
	for key, value := range JobStatus {
		if value == statusID {
			jobStatus = key
			break
		}
	}
	return jobStatus
}

const (
	Job_Act_Manufacturing                   = 1
	Job_Act_ResearchingTechnology           = 2
	Job_Act_ResearchingTimeProductivity     = 3
	Job_Act_ResearchingMaterialProductivity = 4
	Job_Act_Copying                         = 5
	Job_Act_Duplicating                     = 6
	Job_Act_ReverseEngineering              = 7
	Job_Act_Invention                       = 8
)

//https://forums-archive.eveonline.com/topic/55011/
var JobActivity = map[string]int{
	"Manufacturing":   Job_Act_Manufacturing,
	"Tech":            Job_Act_ResearchingTechnology,
	"Time Efficiency": Job_Act_ResearchingTimeProductivity,
	"Mat Efficiency":  Job_Act_ResearchingMaterialProductivity,
	"Copying":         Job_Act_Copying,
	"Duplicating":     Job_Act_Duplicating,
	"Reverse Eng":     Job_Act_ReverseEngineering,
	"Invention":       Job_Act_Invention,
}

func (obj *Model) JobActivityId2Str(activityID int) (jobActivity string) {
	jobActivity = "unknownStatus"
	for key, value := range JobActivity {
		if value == activityID {
			jobActivity = key
			break
		}
	}
	return jobActivity
}

type DBJob struct {
	CharId               int
	CorpId               int
	IsCorp               int
	ActivityId           int
	BlueprintId          int64
	BlueprintLocationId  int64
	BlueprintTypeId      int
	CompletedCharacterId int
	CompletedDate        int64
	Cost                 float64
	Duration             int
	EndDate              int64
	FacilityId           int64
	InstallerId          int
	JobId                int
	LicensedRuns         int
	OutputLocationId     int64
	PauseDate            int64
	Probability          float64
	ProductTypeId        int
	Runs                 int
	StartDate            int64
	StationId            int64
	Status               int
	SuccessfulRuns       int
}

func (obj *Model) createInduTable() {
	if !obj.checkTableExists("industry_jobs") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "industry_jobs" (
			"charId" INT,
			"corpId" INT,
			"isCorpJob" INT,
			"activity_id" INT,
			"blueprint_id" INT,
			"blueprint_location_id" INT,
			"blueprint_type_id" INT,
			"completed_character_id" INT,
			"completed_date" INT,
			"cost" REAL,
			"duration" INT,
			"end_date" INT,
			"facility_id" INT,
			"installer_id" INT,
			"job_id" INT,
			"licensed_runs" INT,
			"output_location_id" INT,
			"pause_date" INT,
			"probability" REAL,
			"product_type_id" INT,
			"runs" INT,
			"start_date" INT,
			"station_id" INT,
			"status" INT,
			"successful_runs" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) JobItemExist(contrID int) bool {
	whereClause := fmt.Sprintf(`job_id="%d"`, contrID)
	num := obj.getNumEntries("industry_jobs", whereClause)
	return num != 0
}

func (obj *Model) AddJobEntry(jobItem *DBJob) DBresult {
	whereClause := fmt.Sprintf(`job_id="%d"`, jobItem.JobId)
	num := obj.getNumEntries("industry_jobs", whereClause)
	retval := DBR_Undefined
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "industry_jobs" (
				charId,
				corpId,
				isCorpJob,
				activity_id,
				blueprint_id,
				blueprint_location_id,
				blueprint_type_id,
				completed_character_id,
				completed_date,
				cost,
				duration,
				end_date,
				facility_id,
				installer_id,
				job_id,
				licensed_runs,
				output_location_id,
				pause_date,
				probability,
				product_type_id,
				runs,
				start_date,
				station_id,
				status,
				successful_runs)
				VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			jobItem.CharId,
			jobItem.CorpId,
			jobItem.IsCorp,
			jobItem.ActivityId,
			jobItem.BlueprintId,
			jobItem.BlueprintLocationId,
			jobItem.BlueprintTypeId,
			jobItem.CompletedCharacterId,
			jobItem.CompletedDate,
			jobItem.Cost,
			jobItem.Duration,
			jobItem.EndDate,
			jobItem.FacilityId,
			jobItem.InstallerId,
			jobItem.JobId,
			jobItem.LicensedRuns,
			jobItem.OutputLocationId,
			jobItem.PauseDate,
			jobItem.Probability,
			jobItem.ProductTypeId,
			jobItem.Runs,
			jobItem.StartDate,
			jobItem.StationId,
			jobItem.Status,
			jobItem.SuccessfulRuns)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		whereClause := fmt.Sprintf(`job_id="%d" AND status="%d"`, jobItem.JobId, jobItem.Status)
		num := obj.getNumEntries("industry_jobs", whereClause)
		if num == 0 {
			stmt, err := obj.DB.Prepare(`
				UPDATE "industry_jobs" SET
					charId=?,
					corpId=?,
					isCorpJob=?,
					activity_id=?,
					blueprint_id=?,
					blueprint_location_id=?,
					blueprint_type_id=?,
					completed_character_id=?,
					completed_date=?,
					cost=?,
					duration=?,
					end_date=?,
					facility_id=?,
					installer_id=?,
					licensed_runs=?,
					output_location_id=?,
					pause_date=?,
					probability=?,
					product_type_id=?,
					runs=?,
					start_date=?,
					station_id=?,
					status=?,
					successful_runs=?
					WHERE job_id=?;
			`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				jobItem.CharId,
				jobItem.CorpId,
				jobItem.IsCorp,
				jobItem.ActivityId,
				jobItem.BlueprintId,
				jobItem.BlueprintLocationId,
				jobItem.BlueprintTypeId,
				jobItem.CompletedCharacterId,
				jobItem.CompletedDate,
				jobItem.Cost,
				jobItem.Duration,
				jobItem.EndDate,
				jobItem.FacilityId,
				jobItem.InstallerId,
				jobItem.LicensedRuns,
				jobItem.OutputLocationId,
				jobItem.PauseDate,
				jobItem.Probability,
				jobItem.ProductTypeId,
				jobItem.Runs,
				jobItem.StartDate,
				jobItem.StationId,
				jobItem.Status,
				jobItem.SuccessfulRuns,
				jobItem.JobId)
			util.CheckErr(err)
			affect, err := res.RowsAffected()
			util.CheckErr(err)
			if affect > 0 {
				retval = DBR_Updated
			}
		} else {
			retval = DBR_Skipped
		}

	}
	return retval
}

func (obj *Model) GetNextJobEndTimeStamp(charId int) (timestamp int64) {
	queryString := fmt.Sprintf("SELECT end_date FROM industry_jobs WHERE charId=%d and status=0 ORDER BY end_date ASC;", charId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var value int64
		rows.Scan(&value)
		timestamp = value
		break
	}
	return timestamp
}

func (obj *Model) GetIndustryJobs(entityId int, corp bool) []*DBJob {
	corpJob := 0
	if corp {
		corpJob = 1
	}
	retval := make([]*DBJob, 0, 10)

	var selectOwner string
	if corp {
		selectOwner = fmt.Sprintf("corpId=%d", entityId)
	} else {
		selectOwner = fmt.Sprintf("charId=%d", entityId)
	}

	queryString := fmt.Sprintf(`SELECT 
										charId,
										corpId,
										isCorpJob,
										activity_id,
										blueprint_id,
										blueprint_location_id,
										blueprint_type_id,
										completed_character_id,
										completed_date,
										cost,
										duration,
										end_date,
										facility_id,
										installer_id,
										job_id,
										licensed_runs,
										output_location_id,
										pause_date,
										probability,
										product_type_id,
										runs,
										start_date,
										station_id,
										status,
										successful_runs
			FROM industry_jobs 
			WHERE %s and isCorpJob=%d and status!=%d 
			ORDER BY end_date ASC;`, selectOwner, corpJob, Job_Stat_cancelled)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var jobItem DBJob
		rows.Scan(
			&jobItem.CharId,
			&jobItem.CorpId,
			&jobItem.IsCorp,
			&jobItem.ActivityId,
			&jobItem.BlueprintId,
			&jobItem.BlueprintLocationId,
			&jobItem.BlueprintTypeId,
			&jobItem.CompletedCharacterId,
			&jobItem.CompletedDate,
			&jobItem.Cost,
			&jobItem.Duration,
			&jobItem.EndDate,
			&jobItem.FacilityId,
			&jobItem.InstallerId,
			&jobItem.JobId,
			&jobItem.LicensedRuns,
			&jobItem.OutputLocationId,
			&jobItem.PauseDate,
			&jobItem.Probability,
			&jobItem.ProductTypeId,
			&jobItem.Runs,
			&jobItem.StartDate,
			&jobItem.StationId,
			&jobItem.Status,
			&jobItem.SuccessfulRuns)
		retval = append(retval, &jobItem)
	}
	return retval
}

func (obj *Model) GetIndustryJob(jobId int64) (result *DBJob) {
	queryString := fmt.Sprintf(`SELECT 
										charId,
										corpId,
										isCorpJob,
										activity_id,
										blueprint_id,
										blueprint_location_id,
										blueprint_type_id,
										completed_character_id,
										completed_date,
										cost,
										duration,
										end_date,
										facility_id,
										installer_id,
										job_id,
										licensed_runs,
										output_location_id,
										pause_date,
										probability,
										product_type_id,
										runs,
										start_date,
										station_id,
										status,
										successful_runs
			FROM industry_jobs 
			WHERE job_id=%d
			ORDER BY end_date ASC;`, jobId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var jobItem DBJob
		rows.Scan(
			&jobItem.CharId,
			&jobItem.CorpId,
			&jobItem.IsCorp,
			&jobItem.ActivityId,
			&jobItem.BlueprintId,
			&jobItem.BlueprintLocationId,
			&jobItem.BlueprintTypeId,
			&jobItem.CompletedCharacterId,
			&jobItem.CompletedDate,
			&jobItem.Cost,
			&jobItem.Duration,
			&jobItem.EndDate,
			&jobItem.FacilityId,
			&jobItem.InstallerId,
			&jobItem.JobId,
			&jobItem.LicensedRuns,
			&jobItem.OutputLocationId,
			&jobItem.PauseDate,
			&jobItem.Probability,
			&jobItem.ProductTypeId,
			&jobItem.Runs,
			&jobItem.StartDate,
			&jobItem.StationId,
			&jobItem.Status,
			&jobItem.SuccessfulRuns)
		result = &jobItem
		break
	}

	return
}

func (obj *Model) GetNextPendingJob(entityId int, corp bool) (retval *DBJob, overdueCount int) {
	industry := obj.GetIndustryJobs(entityId, corp)
	if len(industry) > 0 {
		activeSum := 0
		minDuration := int64(0xFFFFFFFF)
		now := time.Now().Unix()
		for _, job := range industry {
			if job.Status == Job_Stat_active {
				activeSum++
				if job.EndDate > now {
					partDuration := job.EndDate - time.Now().Unix()
					if partDuration < minDuration {
						minDuration = partDuration
						retval = job
					}
				} else {
					overdueCount++
				}
			} else if job.Status == Job_Stat_ready {
				overdueCount++
			}
		}
	}
	return
}

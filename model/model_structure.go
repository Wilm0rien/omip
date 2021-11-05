package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

type DBstructureInfo struct {
	CorpId               int
	FuelExpires          int64
	NextReinforceApply   int64
	NextReinforceHour    int
	NextReinforceWeekday int
	ProfileId            int
	ReinForceHour        int
	ReinforceWeekday     int
	State                int
	StateTimerEnd        int64
	StateTimerStart      int64
	StructureID          int64
	SystemID             int
	TypeID               int
	UnanchorsAt          int64
	NameRef              int64
}

type DBstructureService struct {
	StructureID int64
	Name        int64
	State       int64
}

type DBstructureName struct {
	StructureID   int64
	NameRef       int64
	OwnerID       int
	PositionX     float64
	PositionY     float64
	PositionZ     float64
	SolarSystemID int
	TypeID        int
}

const (
	StructState_new                  = 0
	StructState_anchor_vulnerable    = 1
	StructState_anchoring            = 2
	StructState_armor_reinforce      = 3
	StructState_armor_vulnerable     = 4
	StructState_deploy_vulnerable    = 5
	StructState_fitting_invulnerable = 6
	StructState_hull_reinforce       = 7
	StructState_hull_vulnerable      = 8
	StructState_online_deprecated    = 9
	StructState_onlining_vulnerable  = 10
	StructState_shield_vulnerable    = 11
	StructState_unanchored           = 12
	StructState_unknown              = 13
)

var StructureStatus = map[string]int{
	"new_structure":        StructState_new,
	"anchor_vulnerable":    StructState_anchor_vulnerable,
	"anchoring":            StructState_anchoring,
	"armor_reinforce":      StructState_armor_reinforce,
	"armor_vulnerable":     StructState_armor_vulnerable,
	"deploy_vulnerable":    StructState_deploy_vulnerable,
	"fitting_invulnerable": StructState_fitting_invulnerable,
	"hull_reinforce":       StructState_hull_reinforce,
	"hull_vulnerable":      StructState_hull_vulnerable,
	"online_deprecated":    StructState_online_deprecated,
	"onlining_vulnerable":  StructState_onlining_vulnerable,
	"shield_vulnerable":    StructState_shield_vulnerable,
	"unanchored":           StructState_unanchored,
	"unknown":              StructState_unknown,
}

func (obj *Model) createStructureInfoTable() {
	if !obj.checkTableExists("structure_info") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "structure_info" (
			"corporation_id" INT,
			"fuel_expires" INT,
			"next_reinforce_apply" INT,
			"next_reinforce_hour" INT,
			"next_reinforce_weekday" INT,
			"profile_id" INT,
			"reinforce_hour" INT,
			"reinforce_weekday" INT,
			"state" INT,
			"state_timer_end" INT,
			"state_timer_start" INT,
			"structure_id" INT,
			"system_id" INT,
			"type_id" INT,
			"unanchors_at" INT,
			"name" INT
		);`)
		util.CheckErr(err)
	}
}
func (obj *Model) createStructureNameTable() {
	if !obj.checkTableExists("structure_name") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "structure_name" (
			"structure_id" INT,
			"name" INT,
			"owner_id" INT,
			"position_X" INT,
			"position_Y" INT,
			"position_Z" INT,
			"solar_system_id" INT,
			"type_id" INT);`)
		util.CheckErr(err)
	}
}

func (obj *Model) createStructureServiceTable() {
	if !obj.checkTableExists("structure_services") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "structure_services" (
			"structure_id" INT,
			"name" INT,
			"state" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) StrucStatID2Str(state int) string {
	retval := "unknown status"
	for key, value := range StructureStatus {
		if state == value {
			retval = key
			break
		}
	}
	return retval
}

func (obj *Model) AddStructureInfoEntry(structureInfo *DBstructureInfo) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`structure_id="%d"`, structureInfo.StructureID)
	num := obj.getNumEntries("structure_info", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "structure_info" (
				corporation_id,
				fuel_expires,
				next_reinforce_apply,
				next_reinforce_hour,
				next_reinforce_weekday,
				profile_id,
				reinforce_hour,
				reinforce_weekday,
				state,
				state_timer_end,
				state_timer_start,
				structure_id,
				system_id,
				type_id,
				unanchors_at,
	            name)
				VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			structureInfo.CorpId,
			structureInfo.FuelExpires,
			structureInfo.NextReinforceApply,
			structureInfo.NextReinforceHour,
			structureInfo.NextReinforceWeekday,
			structureInfo.ProfileId,
			structureInfo.ReinForceHour,
			structureInfo.ReinforceWeekday,
			structureInfo.State,
			structureInfo.StateTimerEnd,
			structureInfo.StateTimerStart,
			structureInfo.StructureID,
			structureInfo.SystemID,
			structureInfo.TypeID,
			structureInfo.UnanchorsAt,
			structureInfo.NameRef)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		whereClause2 := fmt.Sprintf(`structure_id="%d" AND state="%d" AND fuel_expires="%d" AND reinforce_hour="%d"`,
			structureInfo.StructureID,
			structureInfo.State,
			structureInfo.FuelExpires,
			structureInfo.ReinForceHour)
		num2 := obj.getNumEntries("structure_info", whereClause2)
		if num2 == 0 {
			stmt, err := obj.DB.Prepare(`
				UPDATE "structure_info" SET
					fuel_expires=?,
					reinforce_hour=?,
					reinforce_weekday=?,
					state=?,
					state_timer_end=?,
					state_timer_start=?,
					unanchors_at=?,
				    name=?
					WHERE structure_id=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				structureInfo.FuelExpires,
				structureInfo.ReinForceHour,
				structureInfo.ReinforceWeekday,
				structureInfo.State,
				structureInfo.StateTimerEnd,
				structureInfo.StateTimerStart,
				structureInfo.UnanchorsAt,
				structureInfo.NameRef,
				structureInfo.StructureID)
			util.CheckErr(err)
			affect, err := res.RowsAffected()
			util.CheckErr(err)
			if affect > 0 {
				retval = DBR_Updated
				//log.Printf("structure %d updated! state %s, fuel_expires:%s",
				//	structureInfo.StructureID,
				//	obj.StrucStatID2Str(structureInfo.State),
				//	util.UnixTS2DateTimeStr(structureInfo.FuelExpires))
			}
		} else {
			retval = DBR_Skipped
		}

	}
	return retval
}
func (obj *Model) AddStructureNameEntry(structureName *DBstructureName) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`structure_id="%d"`, structureName.StructureID)
	num := obj.getNumEntries("structure_name", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "structure_name" (
			structure_id,
			name,
			owner_id,
			position_X,
			position_Y,
			position_Z,
			solar_system_id,
			type_id)
			VALUES (?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			structureName.StructureID,
			structureName.NameRef,
			structureName.OwnerID,
			structureName.PositionX,
			structureName.PositionY,
			structureName.PositionZ,
			structureName.SolarSystemID,
			structureName.TypeID)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		whereClause2 := fmt.Sprintf(`structure_id="%d" AND name="%d" AND owner_id="%d"`,
			structureName.StructureID,
			structureName.NameRef,
			structureName.OwnerID)
		num2 := obj.getNumEntries("structure_info", whereClause2)
		if num2 == 0 {
			stmt, err := obj.DB.Prepare(`
				UPDATE "structure_name" SET
					name=?,
					owner_id=?)
					WHERE structure_id=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			res, err := stmt.Exec(
				structureName.NameRef,
				structureName.OwnerID,
				structureName.StructureID)
			util.CheckErr(err)
			util.CheckErr(err)
			affect, err := res.RowsAffected()
			util.CheckErr(err)
			if affect > 0 {
				retval = DBR_Updated
			}
		}
	}
	return retval
}

func (obj *Model) DeleteStructureServiceEntry(structureID int64, serviceName int64) (result bool) {
	whereClause := fmt.Sprintf("structure_id=%d AND name=%d;", structureID, serviceName)
	numDelete := obj.deleteEntries("structure_services", whereClause)
	if numDelete > 0 {
		result = true
	}
	return result
}

func (obj *Model) DeleteStructureServiceEntries(structureID int64) (result bool) {
	whereClause := fmt.Sprintf("structure_id=%d;", structureID)
	numDelete := obj.deleteEntries("structure_services", whereClause)
	if numDelete > 0 {
		result = true
	}
	return result
}

func (obj *Model) DeleteStructureInfoEntries(structureID int64) (result bool) {
	whereClause := fmt.Sprintf("structure_id=%d;", structureID)
	numDelete := obj.deleteEntries("structure_info", whereClause)
	if numDelete > 0 {
		result = true
	}
	return result
}
func (obj *Model) DeleteStructureNameEntries(structureID int64) (result bool) {
	whereClause := fmt.Sprintf("structure_id=%d;", structureID)
	numDelete := obj.deleteEntries("structure_name", whereClause)
	if numDelete > 0 {
		result = true
	}
	return result
}
func (obj *Model) AddStructureServiceEntry(svc *DBstructureService) DBresult {
	retval := DBR_Undefined
	whereClause := fmt.Sprintf(`structure_id="%d" AND name="%d"`, svc.StructureID, svc.Name)
	num := obj.getNumEntries("structure_services", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "structure_services" (
				structure_id,
				name,
				state)
				VALUES(?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			svc.StructureID,
			svc.Name,
			svc.State)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	} else {
		whereClause2 := fmt.Sprintf(`structure_id="%d" AND name="%d" AND state="%d"`, svc.StructureID, svc.Name, svc.State)
		num2 := obj.getNumEntries("structure_services", whereClause2)
		if num2 == 0 {
			// update service entry!
			stmt, err := obj.DB.Prepare(`
				UPDATE "structure_services" SET 
				state=?
				WHERE structure_id=? AND name=?;`)
			util.CheckErr(err)
			defer stmt.Close()
			_, err = stmt.Exec(svc.State, svc.StructureID, svc.Name)
			util.CheckErr(err)
			retval = DBR_Updated
		} else {
			retval = DBR_Skipped
		}
	}
	return retval
}

func (obj *Model) GetListOfServiceEntries(structureId int64) []int64 {
	dbServiceEntries := make([]int64, 0, 4)
	queryString := fmt.Sprintf("SELECT name FROM structure_services WHERE structure_id=%d;", structureId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var serviceName int64
		rows.Scan(&serviceName)
		dbServiceEntries = append(dbServiceEntries, serviceName)
	}
	return dbServiceEntries
}

func (obj *Model) GetServiceEntries(structureId int64) []*DBstructureService {
	dbServiceEntries := make([]*DBstructureService, 0, 4)
	queryString := fmt.Sprintf("SELECT name, state FROM structure_services WHERE structure_id=%d;", structureId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var strsvc DBstructureService
		rows.Scan(&strsvc.Name, &strsvc.State)
		dbServiceEntries = append(dbServiceEntries, &strsvc)
	}
	return dbServiceEntries
}

func (obj *Model) GetStructureStatus(structureId int64) (retval string) {
	queryString := fmt.Sprintf(`
		SELECT state FROM structure_info WHERE structure_id=%d;`, structureId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var stateINT int
	for rows.Next() {
		rows.Scan(&stateINT)
		break
	}
	return obj.StrucStatID2Str(stateINT)
}

func (obj *Model) GetStructureInfo(structureId int64) (retval *DBstructureInfo) {
	queryString := fmt.Sprintf(`
		SELECT
			corporation_id,
			fuel_expires,
			next_reinforce_apply,
			next_reinforce_hour,
			next_reinforce_weekday,
			profile_id,
			reinforce_hour,
			reinforce_weekday,
			state,
			state_timer_end,
			state_timer_start,
			structure_id,
			system_id,
			type_id,
			unanchors_at,
			name
		FROM structure_info 
		WHERE structure_id=%d;`, structureId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var newStructureInfo DBstructureInfo
	for rows.Next() {
		rows.Scan(
			&newStructureInfo.CorpId,
			&newStructureInfo.FuelExpires,
			&newStructureInfo.NextReinforceApply,
			&newStructureInfo.NextReinforceHour,
			&newStructureInfo.NextReinforceWeekday,
			&newStructureInfo.ProfileId,
			&newStructureInfo.ReinForceHour,
			&newStructureInfo.ReinforceWeekday,
			&newStructureInfo.State,
			&newStructureInfo.StateTimerEnd,
			&newStructureInfo.StateTimerStart,
			&newStructureInfo.StructureID,
			&newStructureInfo.SystemID,
			&newStructureInfo.TypeID,
			&newStructureInfo.UnanchorsAt,
			&newStructureInfo.NameRef)
		retval = &newStructureInfo
		break
	}
	return retval
}
func (obj *Model) GetCorpStructures(corpID int) (retval []*DBstructureInfo) {
	retval = make([]*DBstructureInfo, 0, 5)
	queryString := fmt.Sprintf(`
		SELECT
			corporation_id,
			fuel_expires,
			next_reinforce_apply,
			next_reinforce_hour,
			next_reinforce_weekday,
			profile_id,
			reinforce_hour,
			reinforce_weekday,
			state,
			state_timer_end,
			state_timer_start,
			structure_id,
			system_id,
			type_id,
			unanchors_at,
			name
		FROM structure_info 
		WHERE corporation_id=%d
		ORDER BY fuel_expires ASC;`, corpID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var newStructureInfo DBstructureInfo
		rows.Scan(
			&newStructureInfo.CorpId,
			&newStructureInfo.FuelExpires,
			&newStructureInfo.NextReinforceApply,
			&newStructureInfo.NextReinforceHour,
			&newStructureInfo.NextReinforceWeekday,
			&newStructureInfo.ProfileId,
			&newStructureInfo.ReinForceHour,
			&newStructureInfo.ReinforceWeekday,
			&newStructureInfo.State,
			&newStructureInfo.StateTimerEnd,
			&newStructureInfo.StateTimerStart,
			&newStructureInfo.StructureID,
			&newStructureInfo.SystemID,
			&newStructureInfo.TypeID,
			&newStructureInfo.UnanchorsAt,
			&newStructureInfo.NameRef)
		retval = append(retval, &newStructureInfo)
	}
	return retval
}

func (obj *Model) GetStructureName(structureId int64) (retval *DBstructureName) {
	queryString := fmt.Sprintf(`
		SELECT
			name,
			owner_id,
			position_X,
			position_Y,
			position_Z,
			solar_system_id,
			type_id
		FROM structure_name 
		WHERE structure_id=%d;`, structureId)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	var newStructureName DBstructureName
	for rows.Next() {
		rows.Scan(
			&newStructureName.NameRef,
			&newStructureName.OwnerID,
			&newStructureName.PositionX,
			&newStructureName.PositionY,
			&newStructureName.PositionZ,
			&newStructureName.SolarSystemID,
			&newStructureName.TypeID)
		retval = &newStructureName
		break
	}
	return retval
}

func (obj *Model) GetStructureNameStr(structureId int64) (retval string) {
	retval = "Unkown Structure Name"
	nameStruct := obj.GetStructureName(structureId)
	if nameStruct != nil {
		retval, _ = obj.GetStringEntry(nameStruct.NameRef)
	}
	return retval
}

package usecases

import (
	"bitbucket.com/bbd/unzip-parse/entities"
)

func GetMeasurementSyncJob(status string) *[]entities.MeasurementSyncJobs {
	stmntOut, err := entities.DB.Prepare("SELECT * FROM measurement_sync_jobs WHERE status = ?")
	CheckError(err)

	rows, err := stmntOut.Query(status)
	CheckError(err)

	var measurementSyncJob entities.MeasurementSyncJobs
	var measurementSyncJobs []entities.MeasurementSyncJobs

	for rows.Next() {
		rows.Scan(&measurementSyncJob.Id, &measurementSyncJob.FileName, &measurementSyncJob.Notes, &measurementSyncJob.RowsProcessed, &measurementSyncJob.UppkbId, &measurementSyncJob.Status, &measurementSyncJob.CreatedBy, &measurementSyncJob.TimeStamp.CreatedAt, &measurementSyncJob.TimeStamp.UpdatedAt)
		measurementSyncJobs = append(measurementSyncJobs, measurementSyncJob)
	}

	CheckError(rows.Close())

	return &measurementSyncJobs
}

func UpdateNoteAndStatus(query string, args []interface{}) {
	_, err := entities.DB.Exec(query, args...)
	CheckError(err)
}

func ById(measurementSyncJobsId uint, note, status string) (string, []interface{}) {
	return "UPDATE measurement_sync_jobs SET notes = ?, status = ? WHERE id = ?", []interface{}{note, status, measurementSyncJobsId}
}

func ByName(measurementSyncJobsName, note, status string) (string, []interface{}) {
	return "UPDATE measurement_sync_jobs SET notes = ?, status = ? WHERE file_name LIKE %?%", []interface{}{note, status, measurementSyncJobsName}
}

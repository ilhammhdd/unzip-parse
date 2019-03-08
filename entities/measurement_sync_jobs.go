package entities

import (
	"time"
)

type TimeStamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type MeasurementSyncJobs struct {
	Id            uint   `json:"id"`
	FileName      string `json:"file_name"`
	Notes         string `json:"notes"`
	RowsProcessed int    `json:"rows_processed"`
	UppkbId       int    `json:"uppkb_id"`
	Status        string `json:"status"`
	CreatedBy     int    `json:"created_by"`
	TimeStamp
}

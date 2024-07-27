package netdata

import "time"

type Alarm struct {
	Summary              string    `json:"summary"`
	Status               string    `json:"status"`
	ValueString          string    `json:"value_string"`
	Component            string    `json:"component"`
	Type                 string    `json:"type"`
	LastStatusChangeUnix int64     `json:"last_status_change"`
	LastStatusChange     time.Time `json:"last_status_change_human"`
}

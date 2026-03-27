package openarchiver

type IngestionSource struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Provider              string `json:"provider"`
	Status                string `json:"status"`
	LastSyncStartedAt     string `json:"lastSyncStartedAt"`
	LastSyncFinishedAt    string `json:"lastSyncFinishedAt"`
	LastSyncStatusMessage string `json:"lastSyncStatusMessage"`
}

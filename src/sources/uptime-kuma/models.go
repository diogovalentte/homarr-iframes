package uptimekuma

type Heartbeat struct {
	Status int `json:"status"`
}

type HeartbeatResponse struct {
	HeartbeatList map[string][]Heartbeat `json:"heartbeatList"`
}

type UpDownSites struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

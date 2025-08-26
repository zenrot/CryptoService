package schedule

type Request struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}

type Response struct {
	Enabled         bool   `json:"enabled"`
	IntervalSeconds string `json:"interval_seconds"`
	LastUpdated     string `json:"last_update"`
	NextUpdate      string `json:"next_update"`
}

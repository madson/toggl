package toggl

type TimeEntry struct {
	At          string   `json:"at"`
	Billable    bool     `json:"billable"`
	Description string   `json:"description"`
	Duration    int64    `json:"duration"`
	Guid        string   `json:"guid"`
	ID          int64    `json:"id"`
	PID         int64    `json:"pid"`
	Start       string   `json:"start"`
	Stop        string   `json:"stop"`
	Tags        []string `json:"tags"`
	UID         int64    `json:"uid"`
	WID         int64    `json:"wid"`
	CreatedWith string   `json:"created_with"`
}

type Tag struct {
	ID   int64  `json:"id"`
	WID  int64  `json:"wid"`
	Name string `json:"name"`
	At   string `json:"at"`
}

type tagCreationResult struct {
	Tag `json:"data"`
}

type timeEntryCreationResult struct {
	TimeEntry `json:"data"`
}

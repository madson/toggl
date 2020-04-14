package toggl

import "time"

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

type Workspace struct {
	ID                          int       `json:"id"`
	Name                        string    `json:"name"`
	Premium                     bool      `json:"premium"`
	Admin                       bool      `json:"admin"`
	DefaultHourlyRate           int       `json:"default_hourly_rate"`
	DefaultCurrency             string    `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	At                          time.Time `json:"at"`
	LogoURL                     string    `json:"logo_url,omitempty"`
}

type tagCreationResult struct {
	Tag `json:"data"`
}

type timeEntryCreationResult struct {
	TimeEntry `json:"data"`
}

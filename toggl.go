package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	ISO8601 = "2006-01-02T15:04:05-07:00"
)

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}

func (client *Client) GetTimeEntries(start, end time.Time) ([]TimeEntry, error) {
	s := url.QueryEscape(start.Format(ISO8601))
	e := url.QueryEscape(end.Format(ISO8601))
	query := fmt.Sprintf("start_date=%s&end_date=%s", s, e)

	URL := client.togglURLQuery("/api/v8/time_entries", query)

	var result []TimeEntry

	resp, err := http.Get(URL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error while getting time entries %s: %s\n", URL.String(), err.Error())
		return result, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading bytes %s", err.Error())
		return result, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("request error while getting time entries. statusCode %d. content: %s", resp.StatusCode, string(b))
		return result, errors.New(msg)
	}

	err = json.Unmarshal([]byte(b), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing body %s", err.Error())
		return result, err
	}

	return result, nil
}

func (client *Client) GetWorkspaces() ([]Workspace, error) {
	var result []Workspace

	URL := client.togglURL("/api/v8/workspaces")

	resp, err := http.Get(URL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error while getting workspaces %s: %s\n", URL.String(), err.Error())
		return result, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading bytes: %s\n", err.Error())
		return result, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("request error while getting workspaces. statusCode %d. content: %s", resp.StatusCode, string(b))
		return result, errors.New(msg)
	}

	err = json.Unmarshal([]byte(b), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing body: %s", err.Error())
		return result, err
	}

	return result, nil
}

func (client *Client) GetTags(workspaceID int64) ([]Tag, error) {
	var result []Tag

	path := fmt.Sprintf("/api/v8/workspaces/%d/tags", workspaceID)
	URL := client.togglURLQuery(path, "")

	resp, err := http.Get(URL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error while getting tags %s: %s\n", URL.String(), err.Error())
		return result, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading bytes %s", err.Error())
		return result, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("request error while getting tags. statusCode %d. content: %s", resp.StatusCode, string(b))
		return result, errors.New(msg)
	}

	err = json.Unmarshal([]byte(b), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing body %s", err.Error())
		return result, err
	}

	return result, nil
}

func (client *Client) AddTag(name string, workspaceID int64) (Tag, error) {
	var result tagCreationResult

	URL := client.togglURL("/api/v8/tags")
	params := map[string]Tag{
		"tag": Tag{
			WID:  workspaceID,
			Name: name,
		},
	}
	body, _ := json.Marshal(params)

	resp, err := http.Post(URL.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error while adding a tag %s: %s\n", URL.String(), err.Error())
		return result.Tag, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading bytes %s", err.Error())
		return result.Tag, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("request error while adding a tag. statusCode %d. content: %s", resp.StatusCode, string(b))
		return result.Tag, errors.New(msg)
	}

	err = json.Unmarshal([]byte(b), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing body %s", err.Error())
		return result.Tag, err
	}

	return result.Tag, nil
}

func (client *Client) AddTimeEntry(timeEntry TimeEntry) (TimeEntry, error) {
	var result timeEntryCreationResult

	URL := client.togglURL("/api/v8/time_entries")

	params := map[string]struct {
		Billable    bool     `json:"billable"`
		Description string   `json:"description"`
		Duration    int64    `json:"duration"`
		PID         int64    `json:"pid"`
		Start       string   `json:"start"`
		Tags        []string `json:"tags"`
		CreatedWith string   `json:"created_with"`
	}{
		"time_entry": {
			Billable:    timeEntry.Billable,
			Description: fmt.Sprintf("%s (%d)", timeEntry.Description, timeEntry.ID),
			Duration:    timeEntry.Duration,
			PID:         timeEntry.PID,
			Start:       timeEntry.Start,
			Tags:        timeEntry.Tags,
			CreatedWith: "toggl.go",
		},
	}
	body, _ := json.Marshal(params)

	resp, err := http.Post(URL.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error while creating time entry: %s", err.Error())
		return result.TimeEntry, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading bytes: %s", err.Error())
		return result.TimeEntry, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("request error while creating time entry. statusCode %d. content: %s", resp.StatusCode, string(b))
		return result.TimeEntry, errors.New(msg)
	}

	err = json.Unmarshal([]byte(b), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing body %s", err.Error())
		return result.TimeEntry, err
	}

	return result.TimeEntry, nil
}

func (client *Client) togglURL(path string) *url.URL {
	return client.togglURLQuery(path, "")
}

func (client *Client) togglURLQuery(path, query string) *url.URL {
	userInfo := url.UserPassword(client.token, "api_token")

	return &url.URL{
		Scheme:     "https",
		User:       userInfo,
		Host:       "www.toggl.com",
		Path:       path,
		ForceQuery: false,
		RawQuery:   query,
	}
}

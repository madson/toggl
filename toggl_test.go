package toggl

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToggl(t *testing.T) {
	token := os.Getenv("TOGGL_TOKEN")
	if token == "" {
		panic("TOGGL_TOKEN not defined")
	}

	loc, _ := time.LoadLocation("America/Los_Angeles")

	t.Run("get start of the day format", func(t *testing.T) {
		now := time.Now()
		expected := fmt.Sprintf("%d-%02d-%02dT00:00:00-07:00", now.Year(), now.Month(), now.Day())
		actual := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Format(ISO8601)
		assert.Equal(t, expected, actual)
	})

	t.Run("get end of the day format", func(t *testing.T) {
		now := time.Now()
		expected := fmt.Sprintf("%d-%02d-%02dT23:59:59-07:00", now.Year(), now.Month(), now.Day())
		actual := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc).Format(ISO8601)
		assert.Equal(t, expected, actual)
	})

	t.Run("get toggl URL string", func(t *testing.T) {
		now := time.Now()
		start := url.QueryEscape(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Format(ISO8601))
		end := url.QueryEscape(time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc).Format(ISO8601))
		localURL := "https://%s:api_token@www.toggl.com/api/v8/time_entries?start_date=%s&end_date=%s"
		query := fmt.Sprintf("start_date=%s&end_date=%s", start, end)

		expected := fmt.Sprintf(localURL, token, start, end)
		actual := NewClient(token).togglURLQuery("/api/v8/time_entries", query).String()

		assert.Equal(t, expected, actual)
	})

	t.Run("get user's workspace list", func(t *testing.T) {
		workspaces, err := NewClient(token).GetWorkspaces()

		assert.Nil(t, err)
		assert.NotNil(t, workspaces)
	})

	t.Run("get tags on a workspace", func(t *testing.T) {
		workspaceIDValue := os.Getenv("TOGGL_WORKSPACE_ID")
		workspaceID, _ := strconv.ParseInt(workspaceIDValue, 10, 64)

		tags, err := NewClient(token).GetTags(workspaceID)

		assert.Nil(t, err)
		assert.NotNil(t, tags)
	})

	t.Run("get time entries", func(t *testing.T) {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

		result, err := NewClient(token).GetTimeEntries(start, end)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
}

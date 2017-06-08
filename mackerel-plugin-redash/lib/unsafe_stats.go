package mpredash

import (
	"encoding/json"
	"net/http"
	"time"
)

// UnsafeRedashStats represents a redash stats
type UnsafeRedashStats struct {
	WaitTasks       []UnsafeTaskStats `json:"waiting"`
	DoneTasks       []UnsafeTaskStats `json:"done"`
	InProgressTasks []UnsafeTaskStats `json:"in_progress"`
}

// UnsafeTaskStats represents a task stats
type UnsafeTaskStats struct {
	State     string `json:"state"`
	Scheduled bool   `json:"scheduled"`
}

// UnsafeAllTaskStates represents task states
var UnsafeAllTaskStates = []string{
	"waiting",
	"finished",
	"executing_query",
	"failed",
	"processing",
	"checking_alerts",
	"other", // other state is used for comprehensiveness
}

func getUnsafeStats(p RedashPlugin) (*UnsafeRedashStats, error) {
	// get json data
	timeout := time.Duration(p.Timeout) * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(p.URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode the json data to UnsafeRedashStats struct
	var s UnsafeRedashStats
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func filterCount(ss []UnsafeTaskStats, test func(UnsafeTaskStats) bool) (count uint64) {
	for _, s := range ss {
		if test(s) {
			count++
		}
	}
	return
}

func isScheduled(s UnsafeTaskStats) bool { return s.Scheduled }
func isAdhoc(s UnsafeTaskStats) bool     { return !isScheduled(s) }
func isState(state string) func(UnsafeTaskStats) bool {
	return func(s UnsafeTaskStats) bool {
		if state == "other" {
			return isOtherState(s)
		}
		return s.State == state
	}
}
func isOtherState(s UnsafeTaskStats) bool {
	for _, state := range UnsafeAllTaskStates[0 : len(UnsafeAllTaskStates)-1] {
		if s.State == state {
			return false
		}
	}
	return true
}

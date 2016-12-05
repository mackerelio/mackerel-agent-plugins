package mpprocfd

import "testing"

func TestGraphDefinition(t *testing.T) {
	var fd ProcfdPlugin

	graph := fd.GraphDefinition()
	if actual := len(graph); actual != 1 {
		t.Errorf("GraphDefinition(): %d should be 1", actual)
	}
}

type TestOpenFd struct{}

func (o TestOpenFd) getNumOpenFileDesc() (map[string]uint64, error) {
	return map[string]uint64{
		"8273": 90,
		"8274": 100,
		"8275": 95,
	}, nil
}

func TestFetchMetrics(t *testing.T) {
	openFd = TestOpenFd{}
	var fd ProcfdPlugin
	stat, _ := fd.FetchMetrics()

	if actual := stat["max_fd"].(uint64); actual != 100 {
		t.Errorf("FetchMetrics(): max_fd(%d) should be 100", actual)
	}
}

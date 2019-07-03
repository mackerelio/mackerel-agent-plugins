package mpjvm

import (
	"reflect"
	"runtime"
	"testing"
)

func TestGenerateVmid(t *testing.T) {
	var expected string
	remote := "remotehost.local"
	lvmid := "12345"

	expected = "12345@remotehost.local"
	if id := generateVmid(remote, lvmid); id != expected {
		t.Errorf("vmid should be %s, but %v", expected, id)
	}

	expected = "remotehost.local"
	if id := generateVmid(remote, ""); id != expected {
		t.Errorf("vmid should be %s, but %v", expected, id)
	}

	expected = "12345"
	if id := generateVmid("", lvmid); id != expected {
		t.Errorf("vmid should be %s, but %v", expected, id)
	}

	expected = ""
	if id := generateVmid("", ""); id != expected {
		t.Errorf("vmid should be %s, but %v", expected, id)
	}
}

func TestGenerateRemote(t *testing.T) {
	var expected string
	remote := "remote.local:1099"
	host := "host.local"
	port := 12345

	expected = "remote.local:1099"
	if r := generateRemote(remote, "", 0); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}

	expected = "remote.local:1099"
	if r := generateRemote(remote, host, port); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}

	expected = "host.local:12345"
	if r := generateRemote("", host, port); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}

	expected = "host.local"
	if r := generateRemote("", host, 0); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}

	expected = "localhost:12345"
	if r := generateRemote("", "", port); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}

	expected = ""
	if r := generateRemote("", "", 0); r != expected {
		t.Errorf("remote should be %s, but %s", expected, r)
	}
}

func TestFetchMetrics(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skip()
	}
	m := JVMPlugin{
		Lvmid:     "1",
		JstatPath: "testdata/jstat_serialgc.sh",
	}
	actual, err := m.fetchJstatMetrics("-gc")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]float64{
		"S0C":  45184.0,
		"S1C":  45184.0,
		"S0U":  45184.0,
		"S1U":  0.0,
		"EC":   361728.0,
		"EU":   132414.7,
		"OC":   904068.0,
		"OU":   679249.5,
		"MC":   21248.0,
		"MU":   20787.3,
		"CCSC": 2304.0,
		"CCSU": 2105.8,
		"YGC":  22,
		"YGCT": 8.584,
		"FGC":  6,
		"FGCT": 2.343,
		//"CGC": no data
		//"CGCT": no data
		"GCT": 10.927,
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("fetchMetrics('-gc') = %v; want %v", actual, expected)
	}
}

package mpjvm

import (
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

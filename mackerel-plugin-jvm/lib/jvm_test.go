package mpjvm

import (
	"testing"
)

func TestGenerateVmid(t *testing.T) {
	var expected string
	remote := "remotehost.local"
	lvmid := "12345"

	expected = "12345@remotehost.local"
	if id := generateVmid(&remote, &lvmid); id == nil {
		t.Errorf("vmid should not be nil, but got nil")
	} else if *id != expected {
		t.Errorf("vmid should be %s, but %v", expected, *id)
	}

	expected = "remotehost.local"
	if id := generateVmid(&remote, nil); id == nil {
		t.Errorf("vmid should not be nil, but got nil")
	} else if *id != expected {
		t.Errorf("vmid should be %s, but %v", expected, *id)
	}

	expected = "12345"
	if id := generateVmid(nil, &lvmid); id == nil {
		t.Errorf("vmid should not be nil, but got nil")
	} else if *id != expected {
		t.Errorf("vmid should be %s, but %v", expected, *id)
	}

	if id := generateVmid(nil, nil); id != nil {
		t.Errorf("vmid should be nil, but non-nil %v", *id)
	}
}

func TestGenerateRemote(t *testing.T) {
	var expected string
	remote := "remote.local:1099"
	host := "host.local"
	port := 12345

	expected = "remote.local:1099"
	if r := generateRemote(remote, "", 0); r == nil {
		t.Errorf("remote should not be nil, but got nil")
	} else if *r != expected {
		t.Errorf("remote should be %s, but %v", expected, *r)
	}

	expected = "remote.local:1099"
	if r := generateRemote(remote, host, port); r == nil {
		t.Errorf("remote should not be nil, but got nil")
	} else if *r != expected {
		t.Errorf("remote should be %s, but %v", expected, *r)
	}

	expected = "host.local:12345"
	if r := generateRemote("", host, port); r == nil {
		t.Errorf("remote should not be nil, but got nil")
	} else if *r != expected {
		t.Errorf("remote should be %s, but %v", expected, *r)
	}

	expected = "host.local"
	if r := generateRemote("", host, 0); r == nil {
		t.Errorf("remote should not be nil, but got nil")
	} else if *r != expected {
		t.Errorf("remote should be %s, but %v", expected, *r)
	}

	expected = "localhost:12345"
	if r := generateRemote("", "", port); r == nil {
		t.Errorf("remote should not be nil, but got nil")
	} else if *r != expected {
		t.Errorf("remote should be %s, but %v", expected, *r)
	}

	if r := generateRemote("", "", 0); r != nil {
		t.Errorf("remote should be nil, but non-nil %v", *r)
	}
}

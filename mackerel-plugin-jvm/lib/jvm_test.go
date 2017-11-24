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

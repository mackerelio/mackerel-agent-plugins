package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProcVmstat(t *testing.T) {
	stub := `pgpgin 770294
pgpgout 31351354
pswpin 0
pswpout 113`
	stat := make(map[string]float64)

	err := parseProcVmstat(stub, &stat)
	assert.Nil(t, err)
	assert.Equal(t, stat["pgpgin"], 770294)
	assert.Equal(t, stat["pgpgout"], 31351354)
	assert.Equal(t, stat["pswpin"], 0)
	assert.Equal(t, stat["pswpout"], 113)
}

func TestGetProcVmstat(t *testing.T) {
	stub := "/proc/vmstat"

	_, err := os.Stat(stub)
	if err == nil {
		return
	}

	ret, err := getProcVmstat(stub)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Contains(t, ret, "pgpgout")
}

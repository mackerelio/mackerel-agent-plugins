package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSs(t *testing.T) {
	stub := `State      Recv-Q Send-Q                       Local Address:Port                         Peer Address:Port 
LISTEN     0      128                                     :::45103                                  :::*     
LISTEN     0      128                                     :::111                                    :::* 
TIME-WAIT  0      0                         ::ffff:127.0.0.1:80                       ::ffff:127.0.0.1:50082 
ESTAB      0      0                              10.0.25.101:60826                         10.0.25.104:5672  `
	stat := make(map[string]float64)

	err := parseSs(stub, &stat)
	assert.Nil(t, err)
	assert.Equal(t, stat["LISTEN"], 2)
	assert.Equal(t, stat["TIME-WAIT"], 1)
	assert.Equal(t, stat["ESTAB"], 1)
}

func TestGetSs(t *testing.T) {
	_, err := os.Stat("/usr/sbin/ss")
	if err == nil {
		return
	}

	ret, err := getSs()
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Contains(t, ret, "Stats")
}

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

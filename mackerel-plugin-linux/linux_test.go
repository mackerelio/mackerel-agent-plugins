package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProcVmstat(t *testing.T) {
	stub := ``
	stat := make(map[string]float64)

	err := parseProcVmstat(stub, &stat)
	assert.Nil(t, err)
	assert.Equal(t, stat["score-_"], 1)
}

func TestGetProcVmstat(t *testing.T) {
	stub := ``

	ret, err := getProcVmstat(host, uint16(port), path)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
}

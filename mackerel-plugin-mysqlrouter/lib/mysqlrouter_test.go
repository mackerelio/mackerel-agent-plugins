package mpmysqlrouter

import (
	"github.com/rluisr/mysqlrouter-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

// These credentials owner is @rluisr for developing mysqlrouter-go
// https://github.com/rluisr/mysqlrouter-go
var (
	testURL  = "https://mysqlrouter-test.xzy.pw"
	testUser = "luis"
	testPass = "luis"
)

func TestMRPlugin_FetchMetrics(t *testing.T) {
	mrClient, err := mysqlrouter.New(testURL, testUser, testPass)
	assert.NoError(t, err)

	mr := MRPlugin{
		MRClient: mrClient,
	}

	out, err := mr.FetchMetrics()
	assert.NoError(t, err)
	assert.NotEmpty(t, out)
}

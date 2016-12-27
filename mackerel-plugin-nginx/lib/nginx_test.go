package mpnginx

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var nginx NginxPlugin

	graphdef := nginx.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var nginx NginxPlugin
	stub := `Active connections: 123
server accepts handled requests
 1693613501 1693613501 7996986318
Reading: 66 Writing: 16 Waiting: 41
`

	nginxStats := bytes.NewBufferString(stub)

	stat, err := nginx.parseStats(nginxStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["writing"]).String(), "float64")
	assert.EqualValues(t, stat["writing"], 16)
	assert.EqualValues(t, reflect.TypeOf(stat["accepts"]).String(), "float64")
	assert.EqualValues(t, stat["accepts"], 1693613501)
}

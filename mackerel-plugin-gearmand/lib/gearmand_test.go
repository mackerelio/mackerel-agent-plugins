package mpgearmand

import (
	"bytes"
	"reflect"
	"testing"
	// "github.com/k0kubun/pp"

	"github.com/stretchr/testify/assert"
)

var stub = `Job::Foo	0	0	6
prefix1	Job::Bar	0	0	18
prefix2	Job::Baz	1	1	18
.
`

func TestGraphDefinition(t *testing.T) {
	var gearmand GearmandPlugin

	graphdef := gearmand.GraphDefinition()
	// pp.Print(graphdef)
	if len(graphdef) != 1 {
		t.Errorf("parseDefinition: %d should be 1", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var gearmand GearmandPlugin
	status := bytes.NewBufferString(stub)

	stat, err := gearmand.parseStats(status)
	// pp.Print(stat)
	assert.Nil(t, err)
	if len(stat) != 9 {
		t.Errorf("parseStats: %d should be 9", len(stat))
	}
	for _, val := range stat {
		assert.EqualValues(t, reflect.TypeOf(val).String(), "uint32")
	}
	assert.EqualValues(t, stat["gearmand.queue.Job--Foo.available"].(uint32), 6)
	assert.EqualValues(t, stat["gearmand.queue.Job--Foo.running"].(uint32), 0)
	assert.EqualValues(t, stat["gearmand.queue.Job--Foo.total"].(uint32), 0)
	assert.EqualValues(t, stat["gearmand.queue.prefix1-Job--Bar.available"].(uint32), 18)
	assert.EqualValues(t, stat["gearmand.queue.prefix1-Job--Bar.running"].(uint32), 0)
	assert.EqualValues(t, stat["gearmand.queue.prefix1-Job--Bar.total"].(uint32), 0)
	assert.EqualValues(t, stat["gearmand.queue.prefix2-Job--Baz.available"].(uint32), 18)
	assert.EqualValues(t, stat["gearmand.queue.prefix2-Job--Baz.running"].(uint32), 1)
	assert.EqualValues(t, stat["gearmand.queue.prefix2-Job--Baz.total"].(uint32), 1)
}

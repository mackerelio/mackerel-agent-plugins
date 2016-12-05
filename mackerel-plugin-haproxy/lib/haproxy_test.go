package mphaproxy

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var haproxy HAProxyPlugin

	graphdef := haproxy.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var haproxy HAProxyPlugin
	stub := `# pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,
hastats,FRONTEND,,,1,1,64,43,7061,15994,0,0,0,,,,,OPEN,,,,,,,,,1,1,0,,,,0,2,0,2,,,,0,10,0,15,17,0,,2,2,43,,,0,0,0,0,,,,,,,,
hastats,BACKEND,0,0,0,1,7,17,7061,15994,0,0,,17,0,0,0,UP,0,0,0,,0,1543,0,,1,1,0,,0,,1,0,,1,,,,0,0,0,0,17,0,,,,,0,0,0,0,0,0,0,,,0,0,0,0,
`

	haproxyStats := bytes.NewBufferString(stub)

	stat, err := haproxy.parseStats(haproxyStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	// HaProxy Stats
	assert.EqualValues(t, stat["sessions"], 17)
	assert.EqualValues(t, stat["bytes_in"], 7061)
	assert.EqualValues(t, stat["bytes_out"], 15994)
	assert.EqualValues(t, stat["connection_errors"], 17)
}

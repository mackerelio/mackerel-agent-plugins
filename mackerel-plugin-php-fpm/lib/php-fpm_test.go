package mpphpfpm

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	jsonStr := `{
    "pool":"www",
    "process manager":"dynamic",
    "start time":1461398921,
    "start since":1624,
    "accepted conn":664,
    "listen queue":1,
    "max listen queue":3,
    "listen queue len":2,
    "idle processes":40,
    "active processes":10,
    "total processes":50,
    "max active processes":100,
    "max children reached":200,
    "slow requests":1000
  }`

	httpmock.RegisterResponder("GET", "http://httpmock/status?json",
		httpmock.NewStringResponder(200, jsonStr))

	p := PhpFpmPlugin{
		URL:     "http://httpmock/status?json",
		Prefix:  "php-fpm",
		Timeout: 5,
	}
	status, _ := getStatus(p)

	assert.EqualValues(t, 50, status.TotalProcesses)
	assert.EqualValues(t, 10, status.ActiveProcesses)
	assert.EqualValues(t, 40, status.IdleProcesses)
	assert.EqualValues(t, 100, status.MaxActiveProcesses)
	assert.EqualValues(t, 200, status.MaxChildrenReached)
	assert.EqualValues(t, 1, status.ListenQueue)
	assert.EqualValues(t, 2, status.ListenQueueLen)
	assert.EqualValues(t, 3, status.MaxListenQueue)
	assert.EqualValues(t, 1000, status.SlowRequests)
}

func TestSocketFlag_Set(t *testing.T) {
	tests := []struct {
		Name string
		URL  string

		Network string
		Address string
		Err     string
	}{
		{
			Name:    "parse absolute filepath",
			URL:     "/dev/null",
			Network: "unix",
			Address: "/dev/null",
		},
		{
			Name:    "parse relative filepath",
			URL:     "a/b/c",
			Network: "unix",
			Address: "a/b/c",
		},
		{
			Name:    "parse filename",
			URL:     "file",
			Network: "unix",
			Address: "file",
		},
		{
			Name:    "parse Go's address style",
			URL:     "localhost:9000",
			Network: "tcp",
			Address: "localhost:9000",
		},
		{
			Name:    "parse ipv6 address",
			URL:     "[::1]:1000",
			Network: "tcp",
			Address: "[::1]:1000",
		},
		{
			Name:    "parse unix:// scheme",
			URL:     "unix:///path/to",
			Network: "unix",
			Address: "/path/to",
		},
		{
			Name:    "parse tcp:// scheme",
			URL:     "tcp://localhost",
			Network: "tcp",
			Address: "localhost:9000",
		},
		{
			Name:    "parse tcp:// scheme hostport",
			URL:     "tcp://localhost:9000/",
			Network: "tcp",
			Address: "localhost:9000",
		},
		{
			Name: "parse error: empty scheme",
			URL:  "://test/",
			Err:  "parse",
		},
		{
			Name: "parse error: unknown scheme",
			URL:  "aaa://test/",
			Err:  "parse",
		},
		{
			Name: "parse error: syntax",
			URL:  ":@:",
			Err:  "parse",
		},
		{
			Name: "parse error: no host or port",
			URL:  "?aaa",
			Err:  "parse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var p SocketFlag
			err := p.Set(tt.URL)
			assert.Equal(t, tt.Network, p.Network)
			assert.Equal(t, tt.Address, p.Address)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.Err)
			}
		})
	}
}

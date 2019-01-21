package mpphpfpm

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// FastCGIServer is a FastCGI server, for use in tests.
type FastCGIServer struct {
	lis     net.Listener
	Address string
	URL     string
}

// NewFastCGIServer returns a server that is listening on address.
func NewFastCGIServer(network, address string, handler http.Handler) (*FastCGIServer, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	s := &FastCGIServer{
		lis:     l,
		Address: l.Addr().String(),
		URL:     "http://localhost/status?json",
	}
	go fcgi.Serve(s.lis, handler)
	return s, nil
}

// Close shutdown the server.
func (s *FastCGIServer) Close() error {
	return s.lis.Close()
}

func TestFCGITransport(t *testing.T) {
	dir, err := ioutil.TempDir("", "mackerel-plugin-php-fpm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Name    string
		Network string
		Address string
	}{
		{
			Name:    "listening on TCP-IPv4 address",
			Network: "tcp",
			Address: "127.0.0.1:0",
		},
		{
			Name:    "listening on TCP-hostname",
			Network: "tcp",
			Address: "localhost:0",
		},
		{
			Name:    "listening Unix socket",
			Network: "unix",
			Address: filepath.Join(dir, "php-fpm.sock"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			testFCGITransport(t, tt.Network, tt.Address)
		})
	}
}

func testFCGITransport(t *testing.T, network, address string) {
	tests := []struct {
		Status int
	}{
		{
			Status: http.StatusOK,
		},
		{
			Status: http.StatusBadRequest,
		},
		{
			Status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		msg := http.StatusText(tt.Status)
		t.Run(msg, func(t *testing.T) {
			ts, err := NewFastCGIServer(network, address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.Status)
				w.Write([]byte(msg))
			}))
			if err != nil {
				assert.Fail(t, "failed to launch FastCGI server", err)
				return
			}
			defer ts.Close()

			c := http.Client{
				Transport: &FastCGITransport{
					Network: network,
					Address: ts.Address,
				},
			}
			resp, err := c.Get(ts.URL)
			if err != nil {
				assert.Fail(t, "failed to request a resource", err)
				return
			}
			defer resp.Body.Close()

			assert.Equal(t, tt.Status, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				assert.Fail(t, "failed to read response", err)
				return
			}
			assert.Equal(t, msg, string(b))
		})
	}
}

func TestFCGITransportDialTimeout(t *testing.T) {
	ts, err := NewFastCGIServer("tcp", "127.0.0.1:0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "shouldn't reach here")
	}))
	if err != nil {
		assert.FailNow(t, "failed to launch FastCGI server", err)
	}
	defer ts.Close()

	c := http.Client{
		Transport: &FastCGITransport{
			Network: "tcp",
			Address: ts.Address,
		},
	}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		assert.FailNow(t, "failed to create a request", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "timeout")
	}
	if resp != nil {
		resp.Body.Close()
	}
}

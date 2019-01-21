package mpphpfpm

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/tomasen/fcgi_client"
)

// FastCGITransport is an implementation of RoundTripper that supports FastCGI.
type FastCGITransport struct {
	Network string
	Address string
}

func (*FastCGITransport) timeout(req *http.Request) time.Duration {
	t, ok := req.Context().Deadline()
	if !ok {
		return 0 // no timeout
	}
	return t.Sub(time.Now())
}

// RoundTrip implements the RoundTripper interface.
func (t *FastCGITransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c, err := fcgiclient.DialTimeout(t.Network, t.Address, t.timeout(req))
	if err != nil {
		return nil, err
	}
	defer c.Close()

	params := make(map[string]string)
	params["REQUEST_METHOD"] = req.Method
	if req.ContentLength >= 0 {
		params["CONTENT_LENGTH"] = strconv.FormatInt(req.ContentLength, 10)
	}

	// https://github.com/dreamcat4/php-fpm/blob/master/cgi/cgi_main.c#L781
	params["PATH_INFO"] = req.URL.Path       // TODO(lufia): correct?
	params["SCRIPT_NAME"] = req.URL.Path     // TODO(lufia): correct?
	params["SCRIPT_FILENAME"] = req.URL.Path // TODO(lufia): correct?
	params["REQUEST_URI"] = req.URL.RequestURI()
	params["QUERY_STRING"] = req.URL.RawQuery
	params["SERVER_NAME"] = req.URL.Hostname()
	params["SERVER_ADDR"] = req.URL.Port()
	if req.URL.Scheme == "https" {
		params["HTTPS"] = "on"
	}
	params["SERVER_PROTOCOL"] = req.Proto
	if ctype := req.Header.Get("Content-Type"); ctype != "" {
		params["CONTENT_TYPE"] = ctype
	}
	if ua := req.Header.Get("User-Agent"); ua != "" {
		params["USER_AGENT"] = ua
	}
	resp, err := c.Request(params, req.Body)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	// Var c can disconnect before end of reading of resp.Body.
	// So contents of resp.Body should copy in memory.
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(&buf)
	return resp, nil
}

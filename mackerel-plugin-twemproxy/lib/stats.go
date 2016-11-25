package mptwemproxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// TwemproxyStats represents a twemproxy stats
type TwemproxyStats struct {
	TotalConnections *uint64
	CurrConnections  *uint64
	Pools            map[string]*PoolStats
}

// PoolStats represents a pool stats
type PoolStats struct {
	ClientEOF         *uint64
	ClientErr         *uint64
	ClientConnections *uint64
	ServerEjects      *uint64
	ForwardError      *uint64
	Servers           map[string]*ServerStats
}

// ServerStats represents a server stats
type ServerStats struct {
	ServerEOF         *uint64
	ServerErr         *uint64
	ServerTimedout    *uint64
	ServerConnections *uint64
	OutQueueBytes     *uint64
	InQueueBytes      *uint64
	OutQueue          *uint64
	InQueue           *uint64
	RequestBytes      *uint64
	ResponseBytes     *uint64
	Requests          *uint64
	Responses         *uint64
}

func getStats(p TwemproxyPlugin) (*TwemproxyStats, error) {
	// get json data
	address := p.Address
	timeout := time.Duration(p.Timeout) * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}
	res := bufio.NewReader(conn)

	// decode the json data to TwemproxyStats struct
	var t TwemproxyStats
	decoder := json.NewDecoder(res)
	err = decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UnmarshalJSON interface for json.Unmarshaler
func (t *TwemproxyStats) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	t.Pools = make(map[string]*PoolStats)

L:
	for k, v := range raw {
		switch v.(type) {
		case float64:
			cv := uint64(v.(float64))
			switch k {
			case "total_connections":
				t.TotalConnections = &cv
			case "curr_connections":
				t.CurrConnections = &cv
			case "uptime", "timestamp":
				// do not use these parameters. skip.
			default:
				err = fmt.Errorf("invalid key: %v in rawTwemproxy: %v", k, raw)
				break L
			}
		case map[string]interface{}:
			pool, perr := decodePoolStats(v.(map[string]interface{}))
			if perr != nil {
				err = perr
				break L
			}
			t.Pools[k] = pool
		case string:
			// do not use parameters(service, source, version). skip.
		default:
			err = fmt.Errorf("invalid type in rawTwemproxy: %v", raw)
			break L
		}
	}

	return err
}

func decodePoolStats(rawStats map[string]interface{}) (*PoolStats, error) {
	pool := new(PoolStats)
	pool.Servers = make(map[string]*ServerStats)

	var err error

L:
	for k, v := range rawStats {
		switch v.(type) {
		case float64:
			cv := uint64(v.(float64))
			switch k {
			case "client_eof":
				pool.ClientEOF = &cv
			case "client_err":
				pool.ClientErr = &cv
			case "client_connections":
				pool.ClientConnections = &cv
			case "server_ejects":
				pool.ServerEjects = &cv
			case "forward_error":
				pool.ForwardError = &cv
			case "fragments":
				// do not use this parameter. skip.
			default:
				err = fmt.Errorf("invalid key: %v in rawPool: %v", k, rawStats)
				break L
			}
		case map[string]interface{}:
			server, serr := decodeServerStats(v.(map[string]interface{}))
			if serr != nil {
				err = serr
				break L
			}
			pool.Servers[k] = server
		default:
			err = fmt.Errorf("invalid type in rawPool: %v", rawStats)
			break L
		}
	}

	if err != nil {
		return nil, err
	}
	return pool, nil
}

func decodeServerStats(rawStats map[string]interface{}) (*ServerStats, error) {
	server := new(ServerStats)

	var err error

L:
	for k, v := range rawStats {
		cv := uint64(v.(float64))
		switch k {
		case "server_eof":
			server.ServerEOF = &cv
		case "server_err":
			server.ServerErr = &cv
		case "server_timedout":
			server.ServerTimedout = &cv
		case "server_connections":
			server.ServerConnections = &cv
		case "out_queue_bytes":
			server.OutQueueBytes = &cv
		case "in_queue_bytes":
			server.InQueueBytes = &cv
		case "out_queue":
			server.OutQueue = &cv
		case "in_queue":
			server.InQueue = &cv
		case "request_bytes":
			server.RequestBytes = &cv
		case "response_bytes":
			server.ResponseBytes = &cv
		case "requests":
			server.Requests = &cv
		case "responses":
			server.Responses = &cv
		case "server_ejected_at":
			// do not use this parameter. skip.
		default:
			err = fmt.Errorf("invalid key: %v in rawServer: %v", k, rawStats)
			break L
		}
	}

	if err != nil {
		return nil, err
	}
	return server, nil
}

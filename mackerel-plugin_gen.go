package main

import (
	"fmt"

	"github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-inode/lib"
	"github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-memcached/lib"
)

func runPlugin(plug string) error {
	switch plug {
	case "inode":
		mpinode.Do()
	case "memcached":
		mpmemcached.Do()
	default:
		return fmt.Errorf("unknown plugin: %s", plug)
	}
	return nil
}

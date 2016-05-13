package main

var relayPrefix = "graphite-carbon.relay."

var relayMeta = map[string]meta{
	"cpuUsage": meta{
		label: " CPU Usage",
		unit:  "float",
	},
	"memUsage": meta{
		label: " Memory Usage",
		unit:  "integer",
	},
	"metricsRecieved": meta{
		label: " Metrics Recieved",
		unit:  "integer",
	},
	"destinations_attemptedRelays": meta{
		label: " Attempted Relays",
		unit:  "integer",
	},
	"destinations_queuedUntilConnected": meta{
		label: " Queued Until Connected",
		unit:  "integer",
	},
	"destinations_queuedUntilReady": meta{
		label: " Queued Until Ready",
		unit:  "integer",
	},
	"destinations_sent": meta{
		label: " Sent",
		unit:  "integer",
	},
}

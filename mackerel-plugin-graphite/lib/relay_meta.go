package mpgraphite

var relayPrefix = "graphite-carbon.relay."

var relayMeta = map[string]meta{
	"cpuUsage": {
		label: " CPU Usage",
		unit:  "float",
	},
	"memUsage": {
		label: " Memory Usage",
		unit:  "integer",
	},
	"metricsRecieved": {
		label: " Metrics Recieved",
		unit:  "integer",
	},
	"destinations_attemptedRelays": {
		label: " Attempted Relays",
		unit:  "integer",
	},
	"destinations_queuedUntilConnected": {
		label: " Queued Until Connected",
		unit:  "integer",
	},
	"destinations_queuedUntilReady": {
		label: " Queued Until Ready",
		unit:  "integer",
	},
	"destinations_sent": {
		label: " Sent",
		unit:  "integer",
	},
}

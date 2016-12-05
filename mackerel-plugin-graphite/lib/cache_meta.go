package mpgraphite

var cachePrefix = "graphite-carbon.cache."

var cacheMeta = map[string]meta{
	"avgUpdateTime": {
		label: " Average Update Time",
		unit:  "float",
	},
	"committedPoints": {
		label: " Committed Points",
		unit:  "integer",
	},
	"cpuUsage": {
		label: " CPU Usage",
		unit:  "float",
	},
	"creates": {
		label: " Creates",
		unit:  "integer",
	},
	"errors": {
		label: " Errors",
		unit:  "integer",
	},
	"memUsage": {
		label: " Memory Usage",
		unit:  "integer",
	},
	"metricsReceived": {
		label: " Metrics Received",
		unit:  "integer",
	},
	"pointsPerUpdate": {
		label: " Points Per Update",
		unit:  "float",
	},
	"updateOperations": {
		label: " Update Operations",
		unit:  "integer",
	},
	"cache_overflow": {
		label: " Overflow",
		unit:  "integer",
	},
	"cache_queries": {
		label: " Queries",
		unit:  "integer",
	},
	"cache_queues": {
		label: " Queues",
		unit:  "integer",
	},
	"cache_size": {
		label: " Size",
		unit:  "integer",
	},
}

package main

var cachePrefix = "graphite-carbon.cache."

var cacheMeta = map[string]meta{
	"avgUpdateTime": meta{
		label: " Average Update Time",
		unit:  "float",
	},
	"committedPoints": meta{
		label: " Committed Points",
		unit:  "integer",
	},
	"cpuUsage": meta{
		label: " CPU Usage",
		unit:  "float",
	},
	"creates": meta{
		label: " Creates",
		unit:  "integer",
	},
	"errors": meta{
		label: " Errors",
		unit:  "integer",
	},
	"memUsage": meta{
		label: " Memory Usage",
		unit:  "integer",
	},
	"metricsReceived": meta{
		label: " Metrics Received",
		unit:  "integer",
	},
	"pointsPerUpdate": meta{
		label: " Points Per Update",
		unit:  "float",
	},
	"updateOperations": meta{
		label: " Update Operations",
		unit:  "integer",
	},
	"cache_overflow": meta{
		label: " Overflow",
		unit:  "integer",
	},
	"cache_queries": meta{
		label: " Queries",
		unit:  "integer",
	},
	"cache_queues": meta{
		label: " Queues",
		unit:  "integer",
	},
	"cache_size": meta{
		label: " Size",
		unit:  "integer",
	},
}

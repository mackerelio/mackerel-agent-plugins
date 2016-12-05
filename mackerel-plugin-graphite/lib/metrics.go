package mpgraphite

import (
	"regexp"
	"strings"
)

type metrics struct {
	Target     string
	Datapoints [][]interface{}
}

// carbon.agents.host_name-{instance|a}.{metric}
// carbon.agents.host_name-{instance|a}.cache.{metric}
var cacheRegexp = regexp.MustCompile(`carbon\.agents\..*-(.*?)\.(.*)`)

// carbon.relays.host_name-a.{metric}
// carbon.relays.host_name-a.destinations.{instance|127_0_0_1:3004:a}.{metric}
var relayRegexp = regexp.MustCompile(`carbon\.relays\..*-.*?\.(.*)`)

func (m metrics) getInstanceName() string {
	matched := cacheRegexp.FindStringSubmatch(m.Target)
	if matched != nil {
		return matched[1]
	}
	return ""
}

func (m metrics) getDestinationName() string {
	matched := relayRegexp.FindStringSubmatch(m.Target)
	if matched == nil {
		return ""
	}

	if !strings.Contains(matched[1], ".") {
		// If metric is {cpuUsage,memUsage,metricsRecieved}
		return ""
	}

	return strings.Split(matched[1], ".")[1]
}

func (m metrics) getMetricName() string {
	// In case of carbon-cache
	matched := cacheRegexp.FindStringSubmatch(m.Target)
	if matched != nil {
		metric := matched[2]
		if strings.Contains(matched[2], ".") {
			metric = strings.Replace(metric, ".", "_", -1)
		}
		return metric
	}

	// In case of carbon-relay
	matched = relayRegexp.FindStringSubmatch(m.Target)
	if matched != nil {
		metric := matched[1]
		if !strings.Contains(metric, ".") {
			return metric
		}
		return "destinations_" + strings.Split(metric, ".")[2]
	}

	return ""
}

func (m metrics) getUnitType() string {
	name := m.getMetricName()
	if m, ok := cacheMeta[name]; ok {
		return m.unit
	}
	if m, ok := relayMeta[name]; ok {
		return m.unit
	}
	return ""
}

func (m metrics) isDataAllNil() bool {
	for _, d := range m.Datapoints {
		if d[0] != nil {
			return false
		}
	}
	return true
}

func (m metrics) getMetricKey() string {
	// In case of carbon-cache
	matched := cacheRegexp.FindStringSubmatch(m.Target)
	if matched != nil {
		instance := matched[1]
		metric := matched[2]
		if strings.Contains(matched[2], ".") {
			metric = strings.Replace(metric, ".", "_", -1)
		}
		return cachePrefix + metric + "." + instance
	}

	// In case of carbon-relay
	matched = relayRegexp.FindStringSubmatch(m.Target)
	if matched != nil {
		metric := matched[1]
		if !strings.Contains(metric, ".") {
			// If metric is {cpuUsage,memUsage,metricsRecieved}
			return relayPrefix + metric + "." + metric
		}
		split := strings.Split(metric, ".")
		dest := split[1]
		metric = split[2]
		dest = strings.Replace(dest, ":", "-", -1)
		return relayPrefix + "destinations_" + metric + "." + dest
	}

	return ""
}

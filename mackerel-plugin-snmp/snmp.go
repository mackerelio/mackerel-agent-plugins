package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alouca/gosnmp"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// SNMPMetrics metrics
type SNMPMetrics struct {
	OID     string
	Metrics mp.Metrics
}

// SNMPPlugin mackerel plugin for snmp
type SNMPPlugin struct {
	GraphName        string
	GraphUnit        string
	Host             string
	Community        string
	Tempfile         string
	SNMPMetricsSlice []SNMPMetrics
}

// FetchMetrics interface for mackerelplugin
func (m SNMPPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	s, err := gosnmp.NewGoSNMP(m.Host, m.Community, gosnmp.Version2c, 30)
	if err != nil {
		return nil, err
	}

	for _, sm := range m.SNMPMetricsSlice {
		resp, err := s.Get(sm.OID)
		if err != nil {
			log.Println("SNMP get failed: ", err)
			continue
		}

		ret, err := strconv.ParseFloat(fmt.Sprint(resp.Variables[0].Value), 64)
		if err != nil {
			log.Println(err)
			continue
		}

		stat[sm.Metrics.Name] = ret
	}

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (m SNMPPlugin) GraphDefinition() map[string](mp.Graphs) {
	metrics := []mp.Metrics{}
	for _, sm := range m.SNMPMetricsSlice {
		metrics = append(metrics, sm.Metrics)
	}

	return map[string](mp.Graphs){
		m.GraphName: mp.Graphs{
			Label:   m.GraphName,
			Unit:    m.GraphUnit,
			Metrics: metrics,
		},
	}
}

func main() {
	optGraphName := flag.String("name", "snmp", "Graph name")
	optGraphUnit := flag.String("unit", "float", "Graph unit")

	optHost := flag.String("host", "localhost", "Hostname")
	optCommunity := flag.String("community", "public", "SNMP V2c Community")

	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var snmp SNMPPlugin
	snmp.Host = *optHost
	snmp.Community = *optCommunity
	snmp.GraphName = *optGraphName
	snmp.GraphUnit = *optGraphUnit

	sms := []SNMPMetrics{}
	for _, arg := range flag.Args() {
		vals := strings.Split(arg, ":")
		if len(vals) < 2 {
			continue
		}

		mpm := mp.Metrics{Name: vals[1], Label: vals[1]}
		if len(vals) >= 3 {
			mpm.Diff, _ = strconv.ParseBool(vals[2])
		}
		if len(vals) >= 4 {
			mpm.Stacked, _ = strconv.ParseBool(vals[3])
		}

		sms = append(sms, SNMPMetrics{OID: vals[0], Metrics: mpm})
	}
	snmp.SNMPMetricsSlice = sms

	helper := mp.NewMackerelPlugin(snmp)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

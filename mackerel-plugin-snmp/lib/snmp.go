package mpsnmp

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
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

	gosnmp.Default.Target = m.Host
	gosnmp.Default.Community = m.Community
	gosnmp.Default.Version = gosnmp.Version2c
	gosnmp.Default.Timeout = time.Duration(30) * time.Second

	err := gosnmp.Default.Connect()
	if err != nil {
		return nil, err
	}
	defer gosnmp.Default.Conn.Close()

	for _, sm := range m.SNMPMetricsSlice {
		resp, err := gosnmp.Default.Get([]string{sm.OID})
		if err != nil {
			log.Println("SNMP get failed: ", err)
			continue
		}

		ret, err := strconv.ParseFloat(fmt.Sprint(resp.Variables[0].Value), 64)
		if err != nil {
			// NOTE: Cannot assume strconv.ParseFloat("%s", resp.Variables[0].Value)
			// first, as resp.Variables[0].Value may be an int class, etc.
			// Normally, an object values such as INTEGER or Counter are
			// successfully accepted in the above conversions.
			// However, STRING object values are passed as byte arrays, so the above
			// conversion will result in an error.
			ret, err = strconv.ParseFloat(fmt.Sprintf("%s", resp.Variables[0].Value), 64)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		stat[sm.Metrics.Name] = ret
	}

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (m SNMPPlugin) GraphDefinition() map[string]mp.Graphs {
	metrics := []mp.Metrics{}
	for _, sm := range m.SNMPMetricsSlice {
		metrics = append(metrics, sm.Metrics)
	}

	return map[string]mp.Graphs{
		m.GraphName: {
			Label:   m.GraphName,
			Unit:    m.GraphUnit,
			Metrics: metrics,
		},
	}
}

// Do the plugin
func Do() {
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
		if len(vals) >= 5 {
			switch vals[4] {
			case "uint64":
				mpm.Type = "uint64"
			case "uint32":
				mpm.Type = "uint32"
			default:
				// do nothing
			}
		}

		sms = append(sms, SNMPMetrics{OID: vals[0], Metrics: mpm})
	}
	snmp.SNMPMetricsSlice = sms

	helper := mp.NewMackerelPlugin(snmp)
	helper.Tempfile = *optTempfile

	helper.Run()
}

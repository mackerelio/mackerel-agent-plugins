package mpgraphite

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// GraphitePlugin is plugin for Graphite
type GraphitePlugin struct {
	Host        string
	WebHost     string
	WebPort     string
	Type        string
	Instance    string
	LabelPrefix string
	URL         string
}

// GraphDefinition interface for mackerelplugin
func (p GraphitePlugin) GraphDefinition() map[string]mp.Graphs {
	var graphdef map[string]mp.Graphs
	switch p.Type {
	case "cache":
		graphdef = p.cacheGraphDefinition()
	case "relay":
		graphdef = p.relayGraphDefinition()
	}
	return graphdef
}

// cacheGraphDefinition returns Graphs for carbon-cache
func (p GraphitePlugin) cacheGraphDefinition() map[string]mp.Graphs {
	data, err := p.fetchData()
	if err != nil {
		log.Fatalln(err)
	}

	set := make(map[string]struct{})
	for _, m := range data {
		instance := m.getInstanceName()
		if !m.isDataAllNil() {
			set[instance] = struct{}{}
		}
	}

	graphdef := make(map[string]mp.Graphs)
	for key, m := range cacheMeta {
		var ms []mp.Metrics

		var thetype string
		if m.unit == "float" {
			thetype = "float64"
		} else {
			thetype = "uint64"
		}

		for instance := range set {
			ms = append(ms, mp.Metrics{
				Name:  instance,
				Label: instance,
				Type:  thetype,
			})
		}

		graphdef[cachePrefix+key] = mp.Graphs{
			Label:   p.LabelPrefix + m.label,
			Unit:    m.unit,
			Metrics: ms,
		}
	}

	return graphdef
}

// relayGraphDefinition returns Graphs for carbon-relay
func (p GraphitePlugin) relayGraphDefinition() map[string]mp.Graphs {
	data, err := p.fetchData()
	if err != nil {
		log.Fatalln(err)
	}

	set := make(map[string]struct{})
	for _, m := range data {
		dest := m.getDestinationName()
		if dest != "" && !m.isDataAllNil() {
			set[dest] = struct{}{}
		}
	}

	graphdef := make(map[string]mp.Graphs)
	for key, m := range relayMeta {
		var unit string
		if m.unit == "float" {
			unit = "float64"
		} else if m.unit == "integer" {
			unit = "uint64"
		}

		if !strings.Contains(key, "destinations_") {
			graphdef[relayPrefix+key] = mp.Graphs{
				Label: p.LabelPrefix + m.label,
				Unit:  m.unit,
				Metrics: []mp.Metrics{
					{Name: key, Label: key, Type: unit},
				},
			}
		} else {
			var ms []mp.Metrics
			for dest := range set {
				ms = append(ms, mp.Metrics{
					Name:  strings.Replace(dest, ":", "-", -1),
					Label: dest,
					Type:  unit,
				})
			}

			graphdef[relayPrefix+key] = mp.Graphs{
				Label:   p.LabelPrefix + m.label,
				Unit:    m.unit,
				Metrics: ms,
			}
		}
	}
	return graphdef
}

// FetchMetrics interface for mackerelplugin
// But, don't use this
func (p GraphitePlugin) FetchMetrics() (map[string]interface{}, error) {
	return nil, nil
}

func (p GraphitePlugin) outputValues(w io.Writer) {
	data, err := p.fetchData()
	if err != nil {
		log.Fatalln("fetchData():", err)
	}

	for _, m := range data {
		key := m.getMetricKey()
		if key == "" {
			continue
		}
		unit := m.getUnitType()
		for _, point := range m.Datapoints {
			if point[0] != nil {
				printValue(w, key, point[0], uint64(point[1].(float64)), unit)
			}
		}
	}
}

func printValue(w io.Writer, key string, value interface{}, now uint64, unit string) {
	switch unit {
	case "integer":
		fmt.Fprintf(w, "%s\t%d\t%d\n", key, uint64(value.(float64)), now)
	case "float":
		fmt.Fprintf(w, "%s\t%f\t%d\n", key, value.(float64), now)
	}
}

// fetchData fetches metrics data from -15min
func (p GraphitePlugin) fetchData() ([]metrics, error) {
	res, err := http.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	var d []metrics
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Initialize plugin
func newGraphitePlugin(host, webHost, webPort, thetype, instance, labelPrefix string) GraphitePlugin {
	plugin := GraphitePlugin{}

	// If a hostname is not specified, we get a name reported by the kernel.
	if host == "" {
		h, err := os.Hostname()
		if err != nil {
			log.Fatalln(err)
		}
		plugin.Host = strings.Replace(h, ".", "_", -1)
	} else {
		plugin.Host = strings.Replace(host, ".", "_", -1)
	}

	plugin.WebHost = webHost
	plugin.WebPort = webPort

	switch thetype {
	case "cache":
		plugin.Type = thetype
		if instance == "" {
			plugin.Instance = "*"
		} else {
			plugin.Instance = instance
		}
	case "relay":
		plugin.Type = thetype
		if instance == "" || instance == "*" {
			log.Fatalln("You mush specify concrete instance name in case of relay")
		} else {
			plugin.Instance = instance
		}
	default:
		log.Fatalln("Not accept such a type")
	}

	plugin.LabelPrefix = labelPrefix

	var targets string
	switch plugin.Type {
	case "cache":
		targets = fmt.Sprintf("target=carbon.agents.%s-%s.*&target=carbon.agents.%s-%s.*.*", plugin.Host, plugin.Instance, plugin.Host, plugin.Instance)
	case "relay":
		targets = fmt.Sprintf("target=carbon.relays.%s-%s.*&target=carbon.relays.%s-%s.destinations.*.*", plugin.Host, plugin.Instance, plugin.Host, plugin.Instance)
	}
	plugin.URL = fmt.Sprintf("http://%s:%s/render/?%s&from=-15min&format=json", plugin.WebHost, plugin.WebPort, targets)

	return plugin
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "", "Hostname")
	optWebHost := flag.String("webhost", "", "Graphite-web hostname")
	optWebPort := flag.String("webport", "", "Graphite-web port")
	optType := flag.String("type", "", "Carbon type (cache or relay)")
	optInstance := flag.String("instance", "", "Instance name")
	optLabelPrefix := flag.String("metric-label-prefix", "Carbon", "Metric Label Prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	plugin := newGraphitePlugin(*optHost, *optWebHost, *optWebPort, *optType, *optInstance, *optLabelPrefix)

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		// Not using go-mackerel-plugin-helper method
		// bacause we want to post multiple metrics with arbitrary timestamp
		plugin.outputValues(os.Stdout)
	}
}

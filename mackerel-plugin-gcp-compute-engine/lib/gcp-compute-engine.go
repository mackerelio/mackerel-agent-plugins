package gcpce

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/monitoring/v3"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

const zuluFormat string = "2006-01-02T15:4:05Z"
const duration string = "3m0s"
const computeDomain string = "compute.googleapis.com"

// ComputeEnginePlugin is mackerel plugin for Google Compute Engine
type ComputeEnginePlugin struct {
	Project           string
	InstanceID        string
	InstanceName      string
	MonitoringService *monitoring.Service
	Option            *Option
}

// Option is optional argument to an API call
type Option struct {
	Key string
}

// Get returns key and value
func (c Option) Get() (string, string) {
	return "key", c.Key
}

var graphdef = map[string]mp.Graphs{
	"firewall.dropped_bytes_count": mp.Graphs{
		Label: "FireWall Dropped Bytes Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/firewall/dropped_bytes_count", Label: "Dropped Bytes Count"},
		},
	},
	"firewall.dropped_packets_count": mp.Graphs{
		Label: "FireWall Dropped Packets Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/firewall/dropped_packets_count", Label: "Dropped Packets Count"},
		},
	},
	"cpu.utilization": mp.Graphs{
		Label: "CPU Utilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "/instance/cpu/utilization", Label: "Utilization"},
		},
	},
	"disk.bytes_count": mp.Graphs{
		Label: "Disk Read Bytes Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/instance/disk/read_bytes_count", Label: "Read Bytes Count"},
			{Name: "/instance/disk/write_bytes_count", Label: "Write Bytes Count"},
		},
	},
	"disk.ops_count": mp.Graphs{
		Label: "Disk Read Ops Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/instance/disk/read_ops_count", Label: "Read Ops Count"},
			{Name: "/instance/disk/write_ops_count", Label: "Write Ops Count"},
		},
	},
	"network.bytes_count": mp.Graphs{
		Label: "Network Received Bytes Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/instance/network/received_bytes_count", Label: "Received Bytes Count"},
			{Name: "/instance/network/sent_bytes_count", Label: "Sent Bytes Count"},
		},
	},
	"network.packets_count": mp.Graphs{
		Label: "Network Received Packets Count",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "/instance/network/received_packets_count", Label: "Received Packets Count"},
			{Name: "/instance/network/sent_packets_count", Label: "Sent Packets Count"},
		},
	},
}

// GraphDefinition is return graphdef
func (p ComputeEnginePlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func getLatestValue(listCall *monitoring.ProjectsTimeSeriesListCall, filter string, startTime string, endTime string, opts *Option) (float64, error) {
	res, err := listCall.Filter(filter).IntervalEndTime(endTime).IntervalStartTime(startTime).Do(*opts)

	if err != nil || res == nil {
		return 0, err
	}

	ps := res.TimeSeries[0].Points
	valuePtr := ps[len(ps)-1].Value

	var value float64
	if valuePtr.Int64Value != nil {
		value = float64(*valuePtr.Int64Value)
	} else if valuePtr.DoubleValue != nil {
		value = *valuePtr.DoubleValue
	}

	return value, nil
}

func mkFilter(domain string, metricName string, instance string) string {
	filter := `metric.type = "` + domain + metricName + `"`
	switch domain {
	case computeDomain:
		filter += " AND metric.label.instance_name = " + instance
	}

	return filter
}

// FetchMetrics fetches metrics from Google Monitoring API
func (p ComputeEnginePlugin) FetchMetrics() (map[string]float64, error) {
	now := time.Now()
	formattedEnd := now.Format(zuluFormat)
	m, _ := time.ParseDuration(duration)
	formattedStart := now.Add(-m).Format(zuluFormat)
	listCall := p.MonitoringService.Projects.TimeSeries.List(p.Project)

	stat := map[string]float64{}
	for _, metricName := range []string{
		"/firewall/dropped_bytes_count",
		"/firewall/dropped_packets_count",
		"/instance/cpu/utilization",
		"/instance/disk/read_bytes_count",
		"/instance/disk/read_ops_count",
		"/instance/disk/write_bytes_count",
		"/instance/disk/write_ops_count",
		"/instance/network/received_bytes_count",
		"/instance/network/received_packets_count",
		"/instance/network/sent_bytes_count",
		"/instance/network/sent_packets_count",
	} {
		value, err := getLatestValue(listCall, mkFilter(computeDomain, metricName, p.InstanceName), formattedStart, formattedEnd, p.Option)
		if err != nil {
			continue
		}
		stat[metricName] = value
	}

	return stat, nil
}

// Do the plugin
func Do() {
	optProject := flag.String("project", "", "Project No")
	optInstanceName := flag.String("instance-name", "", "Instance Name")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optAPIKey := flag.String("api-key", "", "API key")
	flag.Parse()

	if *optProject == "" || *optInstanceName == "" || *optInstanceID == "" || *optAPIKey == "" {
		fmt.Println("Errors:", errors.New("Not enough arguments"))
		return
	}

	ctx := context.Background()

	client, err := google.DefaultClient(ctx, monitoring.CloudPlatformScope)
	if err != nil {
		return
	}

	service, err := monitoring.New(client)
	if err != nil {
		return
	}

	var computeEngine = ComputeEnginePlugin{
		MonitoringService: service,
		Project:           "projects/" + *optProject,
		InstanceName:      *optInstanceName,
		InstanceID:        *optInstanceID,
		Option:            &Option{Key: *optAPIKey},
	}

	helper := mp.NewMackerelPlugin(computeEngine)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

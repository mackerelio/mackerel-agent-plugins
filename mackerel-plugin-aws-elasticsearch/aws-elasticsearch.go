package main

import (
	"errors"
	"flag"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
	"os"
	"time"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"es.Nodes": mp.Graphs{
		Label: "AWS ES Nodes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Nodes", Label: "Nodes"},
		},
	},
	"es.CPUUtilization": mp.Graphs{
		Label: "AWS ES CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"es.JVMMemoryPressure": mp.Graphs{
		Label: "AWS ES JVMMemoryPressure",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "JVMMemoryPressure", Label: "JVMMemoryPressure"},
		},
	},
	"es.FreeStorageSpace": mp.Graphs{
		Label: "AWS ES Free Storage Space",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
		},
	},
	"es.SearchableDocuments": mp.Graphs{
		Label: "AWS ES SearchableDocuments",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SearchableDocuments", Label: "SearchableDocuments"},
		},
	},
	"es.DeletedDocuments": mp.Graphs{
		Label: "AWS ES DeletedDocuments",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "DeletedDocuments", Label: "DeletedDocuments"},
		},
	},
	"es.IOPS": mp.Graphs{
		Label: "AWS ES IOPS",
		Unit:  "iops",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReadIOPS", Label: "ReadIOPS"},
			mp.Metrics{Name: "WriteIOPS", Label: "WriteIOPS"},
		},
	},
	"es.Throughput": mp.Graphs{
		Label: "AWS ES Throughput",
		Unit:  "bytes/sec",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReadThroughput", Label: "ReadThroughput"},
			mp.Metrics{Name: "WriteThroughput", Label: "WriteThroughput"},
		},
	},
	"es.DiskQueueDepth": mp.Graphs{
		Label: "AWS ES DiskQueueDepth",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "DiskQueueDepth", Label: "DiskQueueDepth"},
		},
	},
	"es.AutomatedSnapshotFailure": mp.Graphs{
		Label: "AWS ES AutomatedSnapshotFailure",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "AutomatedSnapshotFailure", Label: "AutomatedSnapshotFailure"},
		},
	},
	"es.MasterCPUUtilization": mp.Graphs{
		Label: "AWS ES MasterCPUUtilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "MasterCPUUtilization", Label: "MasterCPUUtilization"},
		},
	},
	"es.Latency": mp.Graphs{
		Label: "AWS ES Latency",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReadLatency", Label: "ReadLatency"},
			mp.Metrics{Name: "WriteLatency", Label: "WriteLatency"},
		},
	},
	"es.MasterJVMMemoryPressure": mp.Graphs{
		Label: "AWS ES MasterJVMMemoryPressure",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "MasterJVMMemoryPressure", Label: "MasterJVMMemoryPressure"},
		},
	},
	"es.MasterFreeStorageSpace": mp.Graphs{
		Label: "AWS ES MasterFreeStorageSpace",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "MasterFreeStorageSpace", Label: "MasterFreeStorageSpace"},
		},
	},
	"es.ClusterStatus": mp.Graphs{
		Label: "AWS ES ClusterStatus",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ClusterStatus.green", Label: "green"},
			mp.Metrics{Name: "ClusterStatus.yellow", Label: "yellow"},
			mp.Metrics{Name: "ClusterStatus.red", Label: "red"},
		},
	},
}

type ESPlugin struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	Domain          string
	ClientId        string
	CloudWatch      *cloudwatch.CloudWatch
}

const ESNameSpace = "AWS/ES"

func (p *ESPlugin) Prepare() error {
	auth, err := aws.GetAuth(p.AccessKeyId, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return err
	}

	p.CloudWatch, err = cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return err
	}
	return nil
}

func (p ESPlugin) GetLastPoint(dimensions *[]cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: *dimensions,
		StartTime:  now.Add(time.Duration(180) * time.Second * -1),
		EndTime:    now,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{"Average"},
		Namespace:  ESNameSpace,
	})

	if err != nil {
		return 0, err
	}

	datapoints := response.GetMetricStatisticsResult.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := time.Unix(0, 0)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(latest) {
			continue
		}

		latest = dp.Timestamp
		latestVal = dp.Average
	}

	return latestVal, nil
}

func (p ESPlugin) FetchMetrics() (map[string]float64, error) {
	dimensions := []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  "DomainName",
			Value: p.Domain,
		},
		cloudwatch.Dimension{
			Name:  "ClientId",
			Value: p.ClientId,
		},
	}

	ret, err := p.CloudWatch.ListMetrics(&cloudwatch.ListMetricsRequest{
		Namespace:  ESNameSpace,
		Dimensions: dimensions,
	})
	if err != nil {
		log.Printf("%s", err)
	}

	stat := make(map[string]float64)

	for _, met := range ret.ListMetricsResult.Metrics {
		v, err := p.GetLastPoint(&dimensions, met.MetricName)
		if err == nil {
			if met.MetricName == "MasterFreeStorageSpace" || met.MetricName == "FreeStorageSpace" {
				// MBytes -> Bytes
				v = v * 1024 * 1024
			}
			stat[met.MetricName] = v
		} else {
			log.Printf("%s: %s", met.MetricName, err)
		}
	}

	return stat, nil
}

func (p ESPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyId := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optClientId := flag.String("client-id", "", "AWS Client ID")
	optDomain := flag.String("domain", "", "ES domain name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var es ESPlugin

	if *optRegion == "" {
		es.Region = aws.InstanceRegion()
	} else {
		es.Region = *optRegion
	}

	es.Domain = *optDomain
	es.ClientId = *optClientId
	es.AccessKeyId = *optAccessKeyId
	es.SecretAccessKey = *optSecretAccessKey

	err := es.Prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(es)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-aws-elasticsearch"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

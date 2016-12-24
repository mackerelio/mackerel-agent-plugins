package mpawselasticsearch

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string]mp.Graphs{
	"es.Nodes": {
		Label: "AWS ES Nodes",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Nodes", Label: "Nodes"},
		},
	},
	"es.CPUUtilization": {
		Label: "AWS ES CPU Utilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"es.JVMMemoryPressure": {
		Label: "AWS ES JVMMemoryPressure",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "JVMMemoryPressure", Label: "JVMMemoryPressure"},
		},
	},
	"es.FreeStorageSpace": {
		Label: "AWS ES Free Storage Space",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
		},
	},
	"es.SearchableDocuments": {
		Label: "AWS ES SearchableDocuments",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "SearchableDocuments", Label: "SearchableDocuments"},
		},
	},
	"es.DeletedDocuments": {
		Label: "AWS ES DeletedDocuments",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "DeletedDocuments", Label: "DeletedDocuments"},
		},
	},
	"es.IOPS": {
		Label: "AWS ES IOPS",
		Unit:  "iops",
		Metrics: []mp.Metrics{
			{Name: "ReadIOPS", Label: "ReadIOPS"},
			{Name: "WriteIOPS", Label: "WriteIOPS"},
		},
	},
	"es.Throughput": {
		Label: "AWS ES Throughput",
		Unit:  "bytes/sec",
		Metrics: []mp.Metrics{
			{Name: "ReadThroughput", Label: "ReadThroughput"},
			{Name: "WriteThroughput", Label: "WriteThroughput"},
		},
	},
	"es.DiskQueueDepth": {
		Label: "AWS ES DiskQueueDepth",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "DiskQueueDepth", Label: "DiskQueueDepth"},
		},
	},
	"es.AutomatedSnapshotFailure": {
		Label: "AWS ES AutomatedSnapshotFailure",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "AutomatedSnapshotFailure", Label: "AutomatedSnapshotFailure"},
		},
	},
	"es.MasterCPUUtilization": {
		Label: "AWS ES MasterCPUUtilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "MasterCPUUtilization", Label: "MasterCPUUtilization"},
		},
	},
	"es.Latency": {
		Label: "AWS ES Latency",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "ReadLatency", Label: "ReadLatency"},
			{Name: "WriteLatency", Label: "WriteLatency"},
		},
	},
	"es.MasterJVMMemoryPressure": {
		Label: "AWS ES MasterJVMMemoryPressure",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "MasterJVMMemoryPressure", Label: "MasterJVMMemoryPressure"},
		},
	},
	"es.MasterFreeStorageSpace": {
		Label: "AWS ES MasterFreeStorageSpace",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "MasterFreeStorageSpace", Label: "MasterFreeStorageSpace"},
		},
	},
	"es.ClusterStatus": {
		Label: "AWS ES ClusterStatus",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ClusterStatus.green", Label: "green"},
			{Name: "ClusterStatus.yellow", Label: "yellow"},
			{Name: "ClusterStatus.red", Label: "red"},
		},
	},
}

// ESPlugin mackerel plugin for aws elasticsearch
type ESPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Domain          string
	ClientID        string
	CloudWatch      *cloudwatch.CloudWatch
}

const esNameSpace = "AWS/ES"

func (p *ESPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	p.CloudWatch = cloudwatch.New(sess, config)
	return nil
}

func (p ESPlugin) getLastPoint(dimensions []*cloudwatch.Dimension, metricName *string) (float64, error) {
	now := time.Now()

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)),
		EndTime:    aws.Time(now),
		MetricName: metricName,
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
		Namespace:  aws.String(esNameSpace),
	})

	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := new(time.Time)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(*latest) {
			continue
		}

		latest = dp.Timestamp
		latestVal = *dp.Average
	}

	return latestVal, nil
}

// FetchMetrics interface for mackerelplugin
func (p ESPlugin) FetchMetrics() (map[string]float64, error) {
	dimensionFilters := []*cloudwatch.DimensionFilter{
		{
			Name:  aws.String("DomainName"),
			Value: aws.String(p.Domain),
		},
		{
			Name:  aws.String("ClientId"),
			Value: aws.String(p.ClientID),
		},
	}

	ret, err := p.CloudWatch.ListMetrics(&cloudwatch.ListMetricsInput{
		Namespace:  aws.String(esNameSpace),
		Dimensions: dimensionFilters,
	})
	if err != nil {
		log.Printf("%s", err)
	}

	stat := make(map[string]float64)

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DomainName"),
			Value: aws.String(p.Domain),
		},
		{
			Name:  aws.String("ClientId"),
			Value: aws.String(p.ClientID),
		},
	}
	for _, met := range ret.Metrics {
		v, err := p.getLastPoint(dimensions, met.MetricName)
		if err == nil {
			if *met.MetricName == "MasterFreeStorageSpace" || *met.MetricName == "FreeStorageSpace" {
				// MBytes -> Bytes
				v = v * 1024 * 1024
			}
			stat[*met.MetricName] = v
		} else {
			log.Printf("%s: %s", met.MetricName, err)
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p ESPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optClientID := flag.String("client-id", "", "AWS Client ID")
	optDomain := flag.String("domain", "", "ES domain name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var es ESPlugin

	if *optRegion == "" {
		ec2metadata := ec2metadata.New(session.New())
		if ec2metadata.Available() {
			es.Region, _ = ec2metadata.Region()
		}
	} else {
		es.Region = *optRegion
	}

	es.Region = *optRegion
	es.Domain = *optDomain
	es.ClientID = *optClientID
	es.AccessKeyID = *optAccessKeyID
	es.SecretAccessKey = *optSecretAccessKey

	err := es.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(es)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

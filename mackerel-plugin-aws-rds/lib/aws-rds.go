package mpawsrds

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// RDSPlugin mackerel plugin for amazon RDS
type RDSPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Identifier      string
	Engine          string
	Prefix          string
	LabelPrefix     string
}

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{dimension},
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)), // 3 min (to fetch at least 1 data-point)
		EndTime:    aws.Time(now),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
		Namespace:  aws.String("AWS/RDS"),
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

// FetchMetrics interface for mackerel-plugin
func (p RDSPlugin) FetchMetrics() (map[string]float64, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	cloudWatch := cloudwatch.New(sess, config)

	stat := make(map[string]float64)

	perInstance := &cloudwatch.Dimension{
		Name:  aws.String("DBInstanceIdentifier"),
		Value: aws.String(p.Identifier),
	}

	for _, met := range p.rdsMetrics() {
		v, err := getLastPoint(cloudWatch, perInstance, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

func (p RDSPlugin) baseGraphDefs() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		p.Prefix + ".DiskQueueDepth": {
			Label: p.LabelPrefix + " BinLogDiskUsage",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "DiskQueueDepth", Label: "Depth"},
			},
		},
		p.Prefix + ".CPUUtilization": {
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		// .CPUCreditBalance ...Only valid for T2 instances
		p.Prefix + ".CPUCreditBalance": {
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		// .CPUCreditUsage ...Only valid for T2 instances
		p.Prefix + ".CPUCreditUsage": {
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
		p.Prefix + ".DatabaseConnections": {
			Label: p.LabelPrefix + " Database Connections",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "DatabaseConnections", Label: "DatabaseConnections"},
			},
		},
		p.Prefix + ".FreeableMemory": {
			Label: p.LabelPrefix + " Freeable Memory",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "FreeableMemory", Label: "FreeableMemory"},
			},
		},
		p.Prefix + ".FreeStorageSpace": {
			Label: p.LabelPrefix + " Free Storage Space",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
			},
		},
		p.Prefix + ".SwapUsage": {
			Label: p.LabelPrefix + " Swap Usage",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "SwapUsage", Label: "SwapUsage"},
			},
		},
		p.Prefix + ".IOPS": {
			Label: p.LabelPrefix + " IOPS",
			Unit:  "iops",
			Metrics: []mp.Metrics{
				{Name: "ReadIOPS", Label: "Read"},
				{Name: "WriteIOPS", Label: "Write"},
			},
		},
		p.Prefix + ".Latency": {
			Label: p.LabelPrefix + " Latency in second",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "ReadLatency", Label: "Read"},
				{Name: "WriteLatency", Label: "Write"},
			},
		},
		p.Prefix + ".Throughput": {
			Label: p.LabelPrefix + " Throughput",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "ReadThroughput", Label: "Read"},
				{Name: "WriteThroughput", Label: "Write"},
			},
		},
		p.Prefix + ".NetworkThroughput": {
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
	}
}

// GraphDefinition interface for mackerel plugin
func (p RDSPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := p.baseGraphDefs()
	switch p.Engine {
	case "mysql", "mariadb":
		graphdef = mergeGraphDefs(graphdef, p.mySQLGraphDefinition())
	case "postgresql":
		graphdef = mergeGraphDefs(graphdef, p.postgreSQLGraphDefinition())
	case "aurora":
		graphdef = p.auroraGraphDefinition()
	}
	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "DB Instance Identifier")
	optEngine := flag.String("engine", "", "RDS Engine")
	optPrefix := flag.String("metric-key-prefix", "rds", "Metric key prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Metric Label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	rds := RDSPlugin{
		Prefix: *optPrefix,
	}
	if *optLabelPrefix == "" {
		if *optPrefix == "rds" {
			rds.LabelPrefix = "RDS"
		} else {
			rds.LabelPrefix = strings.Title(*optPrefix)
		}
	} else {
		rds.LabelPrefix = *optLabelPrefix
	}

	if *optRegion == "" {
		ec2metadata := ec2metadata.New(session.New())
		if ec2metadata.Available() {
			rds.Region, _ = ec2metadata.Region()
		}
	} else {
		rds.Region = *optRegion
	}

	rds.Identifier = *optIdentifier
	rds.AccessKeyID = *optAccessKeyID
	rds.SecretAccessKey = *optSecretAccessKey
	rds.Engine = *optEngine

	helper := mp.NewMackerelPlugin(rds)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

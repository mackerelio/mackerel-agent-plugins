package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
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
	RDSMetrics      []string
}

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: []cloudwatch.Dimension{*dimension},
		StartTime:  now.Add(time.Duration(180) * time.Second * -1), // 3 min (to fetch at least 1 data-point)
		EndTime:    now,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{"Average"},
		Namespace:  "AWS/RDS",
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

// FetchMetrics interface for mackerel-plugin
func (p RDSPlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	cloudWatch, err := cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)

	perInstance := &cloudwatch.Dimension{
		Name:  "DBInstanceIdentifier",
		Value: p.Identifier,
	}

	for _, met := range p.RDSMetrics {
		v, err := getLastPoint(cloudWatch, perInstance, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerel plugin
func (p RDSPlugin) GraphDefinition() map[string](mp.Graphs) {
	var graphdef map[string](mp.Graphs)
	switch p.Engine {
	case "mysql":
		graphdef = p.MySQLGraphDefinition()
	case "aurora":
		graphdef = p.AuroraGraphDefinition()
	case "mariadb":
		graphdef = p.MariaDBGraphDefinition()
	case "postgresql":
		graphdef = p.PostgreSQLGraphDefinition()
	default:
		log.Printf("RDS Engine is 'mysql' or 'aurora' or 'mariadb' or 'postgresql'.")
		os.Exit(1)
	}
	return graphdef
}

func main() {
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
		rds.Region = aws.InstanceRegion()
	} else {
		rds.Region = *optRegion
	}

	rds.Identifier = *optIdentifier
	rds.AccessKeyID = *optAccessKeyID
	rds.SecretAccessKey = *optSecretAccessKey
	rds.Engine = *optEngine

	switch rds.Engine {
	case "mysql":
		rds.RDSMetrics = MetricsdefMySQL
	case "aurora":
		rds.RDSMetrics = MetricsdefAurora
	case "mariadb":
		rds.RDSMetrics = MetricsdefMariaDB
	case "postgresql":
		rds.RDSMetrics = MetricsdefPostgreSQL
	default:
		log.Printf("RDS Engine is 'mysql' or 'aurora' or 'mariadb' or 'postgresql'.")
		os.Exit(1)
	}

	helper := mp.NewMackerelPlugin(rds)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-rds"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

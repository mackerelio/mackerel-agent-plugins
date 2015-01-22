package main

import (
//	"errors"
	"flag"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	//"github.com/crowdmob/goamz/cloudwatch"
	ses "github.com/naokibtn/go-ses"
	//"github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-aws-ses/ses"
	mp "github.com/mackerelio/go-mackerel-plugin"
//	"log"
	"os"
	"time"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
//	"ses.CPUUtilization": mp.Graphs{
//		Label: "SES CPU Utilization",
//		Unit:  "percentage",
//		Metrics: [](mp.Metrics){
//			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
//		},
//	},
}

type SESPlugin struct {
	Region          string
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
}

//func GetLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
//	now := time.Now()
//
//	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
//		Dimensions: []cloudwatch.Dimension{*dimension},
//		StartTime:  now.Add(time.Duration(180) * time.Second * -1), // 3 min (to fetch at least 1 data-point)
//		EndTime:    now,
//		MetricName: metricName,
//		Period:     60,
//		Statistics: []string{"Average"},
//		Namespace:  "AWS/SES",
//	})
//	if err != nil {
//		return 0, err
//	}
//
//	datapoints := response.GetMetricStatisticsResult.Datapoints
//	if len(datapoints) == 0 {
//		return 0, errors.New("fetched no datapoints")
//	}
//
//	latest := time.Unix(0, 0)
//	var latestVal float64
//	for _, dp := range datapoints {
//		if dp.Timestamp.Before(latest) {
//			continue
//		}
//
//		latest = dp.Timestamp
//		latestVal = dp.Average
//	}
//
//	return latestVal, nil
//}

func (p SESPlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyId, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	//region := aws.ServiceInfo{
	//    p.Endpoint,
	//    aws.V2Signature,
	//}

	sescfg := ses.Config{
	    AccessKeyID: auth.AccessKey,
	    SecretAccessKey: auth.SecretKey,
	    Endpoint: p.Endpoint,
	}

	stat := make(map[string]float64)
	quota, err := sescfg.GetSendQuota()
	fmt.Printf("%+v\n", quota)
	fmt.Printf("%+v\n", err)
	stat["hoge"] = 1

	//for _, met := range [...]string{
	//	"BinLogDiskUsage", "CPUUtilization", "DatabaseConnections", "DiskQueueDepth", "FreeableMemory",
	//	"FreeStorageSpace", "ReplicaLag", "SwapUsage", "ReadIOPS", "WriteIOPS", "ReadLatency",
	//	"WriteLatency",
	//} {
	//	v, err := GetLastPoint(ses, perInstance, met)
	//	if err == nil {
	//		stat[met] = v
	//	} else {
	//		log.Printf("%s: %s", met, err)
	//	}
	//}

	return stat, nil
}

func (p SESPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optEndpoint := flag.String("endpoint", "", "AWS Endpoint")
	optAccessKeyId := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ses SESPlugin

	if *optRegion == "" {
		//ses.Region = aws.InstanceRegion()
	} else {
		ses.Region = *optRegion
	}

	ses.Endpoint = *optEndpoint
	ses.AccessKeyId = *optAccessKeyId
	ses.SecretAccessKey = *optSecretAccessKey

	helper := mp.NewMackerelPlugin(ses)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-ses"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

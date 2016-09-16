package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"latency": mp.Graphs{
		Label: "Whole ELB Latency",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Latency", Label: "Latency"},
		},
	},
	"http_backend": mp.Graphs{
		Label: "Whole ELB HTTP Backend Count",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "HTTPCode_Backend_2XX", Label: "2XX", Stacked: true},
			mp.Metrics{Name: "HTTPCode_Backend_3XX", Label: "3XX", Stacked: true},
			mp.Metrics{Name: "HTTPCode_Backend_4XX", Label: "4XX", Stacked: true},
			mp.Metrics{Name: "HTTPCode_Backend_5XX", Label: "5XX", Stacked: true},
		},
	},
	// "healthy_host_count", "unhealthy_host_count" will be generated dynamically
}

type statType int

const (
	stAve statType = iota
	stSum
)

func (s statType) String() string {
	switch s {
	case stAve:
		return "Average"
	case stSum:
		return "Sum"
	}
	return ""
}

// ELBPlugin elb plugin for mackerel
type ELBPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	AZs             []string
	CloudWatch      *cloudwatch.CloudWatch
	Lbname          string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p ELBPlugin) MetricKeyPrefix() string {
	return "elb"
}

func (p *ELBPlugin) prepare() error {
	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return err
	}

	p.CloudWatch, err = cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return err
	}

	ret, err := p.CloudWatch.ListMetrics(&cloudwatch.ListMetricsRequest{
		Namespace: "AWS/ELB",
		Dimensions: []cloudwatch.Dimension{
			{
				Name: "AvailabilityZone",
			},
		},
		MetricName: "HealthyHostCount",
	})

	if err != nil {
		return err
	}

	p.AZs = make([]string, 0, len(ret.ListMetricsResult.Metrics))
	for _, met := range ret.ListMetricsResult.Metrics {
		if len(met.Dimensions) > 1 {
			continue
		} else if met.Dimensions[0].Name != "AvailabilityZone" {
			continue
		}

		p.AZs = append(p.AZs, met.Dimensions[0].Value)
	}

	return nil
}

func (p ELBPlugin) getLastPoint(dimensions *[]cloudwatch.Dimension, metricName string, sTyp statType) (float64, error) {
	now := time.Now()

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: *dimensions,
		StartTime:  now.Add(time.Duration(120) * time.Second * -1), // 2 min (to fetch at least 1 data-point)
		EndTime:    now,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{sTyp.String()},
		Namespace:  "AWS/ELB",
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
		switch sTyp {
		case stAve:
			latestVal = dp.Average
		case stSum:
			latestVal = dp.Sum
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch elb metrics
func (p ELBPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	// HostCount per AZ
	for _, az := range p.AZs {
		d := []cloudwatch.Dimension{
			{
				Name:  "AvailabilityZone",
				Value: az,
			},
		}
		if p.Lbname != "" {
			d2 := cloudwatch.Dimension{
				Name:  "LoadBalancerName",
				Value: p.Lbname,
			}
			d = append(d, d2)
		}
		for _, met := range []string{"HealthyHostCount", "UnHealthyHostCount"} {
			v, err := p.getLastPoint(&d, met, stAve)
			if err == nil {
				stat[met+"_"+az] = v
			}
		}
	}

	glb := []cloudwatch.Dimension{
		{
			Name:  "Service",
			Value: "ELB",
		},
	}
	if p.Lbname != "" {
		g2 := cloudwatch.Dimension{
			Name:  "LoadBalancerName",
			Value: p.Lbname,
		}
		glb = append(glb, g2)
	}

	v, err := p.getLastPoint(&glb, "Latency", stAve)
	if err == nil {
		stat["Latency"] = v
	}

	for _, met := range [...]string{"HTTPCode_Backend_2XX", "HTTPCode_Backend_3XX", "HTTPCode_Backend_4XX", "HTTPCode_Backend_5XX"} {
		v, err := p.getLastPoint(&glb, met, stSum)
		if err == nil {
			stat[met] = v
		}
	}

	return stat, nil
}

// GraphDefinition for Mackerel
func (p ELBPlugin) GraphDefinition() map[string](mp.Graphs) {
	for _, grp := range [...]string{"healthy_host_count", "unhealthy_host_count"} {
		var namePre string
		var label string
		switch grp {
		case "healthy_host_count":
			namePre = "HealthyHostCount_"
			label = "ELB Healthy Host Count"
		case "unhealthy_host_count":
			namePre = "UnHealthyHostCount_"
			label = "ELB Unhealthy Host Count"
		}

		var metrics [](mp.Metrics)
		for _, az := range p.AZs {
			metrics = append(metrics, mp.Metrics{Name: namePre + az, Label: az, Stacked: true})
		}
		graphdef[grp] = mp.Graphs{
			Label:   label,
			Unit:    "integer",
			Metrics: metrics,
		}
	}

	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optLbname := flag.String("lbname", "", "ELB Name")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var elb ELBPlugin

	if *optRegion == "" {
		elb.Region = aws.InstanceRegion()
	} else {
		elb.Region = *optRegion
	}

	elb.AccessKeyID = *optAccessKeyID
	elb.SecretAccessKey = *optSecretAccessKey
	elb.Lbname = *optLbname

	err := elb.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(elb)
	helper.Tempfile = *optTempfile
	helper.Run()
}

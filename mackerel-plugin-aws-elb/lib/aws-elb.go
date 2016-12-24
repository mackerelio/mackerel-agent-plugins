package mpawselb

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
	"elb.latency": {
		Label: "Whole ELB Latency",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Latency", Label: "Latency"},
		},
	},
	"elb.http_backend": {
		Label: "Whole ELB HTTP Backend Count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "HTTPCode_Backend_2XX", Label: "2XX", Stacked: true},
			{Name: "HTTPCode_Backend_3XX", Label: "3XX", Stacked: true},
			{Name: "HTTPCode_Backend_4XX", Label: "4XX", Stacked: true},
			{Name: "HTTPCode_Backend_5XX", Label: "5XX", Stacked: true},
		},
	},
	// "elb.healthy_host_count", "elb.unhealthy_host_count" will be generated dynamically
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
	AZs             []*string
	CloudWatch      *cloudwatch.CloudWatch
	Lbname          string
}

func (p *ELBPlugin) prepare() error {
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

	ret, err := p.CloudWatch.ListMetrics(&cloudwatch.ListMetricsInput{
		Namespace: aws.String("AWS/ELB"),
		Dimensions: []*cloudwatch.DimensionFilter{
			{
				Name: aws.String("AvailabilityZone"),
			},
		},
		MetricName: aws.String("HealthyHostCount"),
	})

	if err != nil {
		return err
	}

	p.AZs = make([]*string, 0, len(ret.Metrics))
	for _, met := range ret.Metrics {
		if len(met.Dimensions) > 1 {
			continue
		} else if *met.Dimensions[0].Name != "AvailabilityZone" {
			continue
		}

		p.AZs = append(p.AZs, met.Dimensions[0].Value)
	}

	return nil
}

func (p ELBPlugin) getLastPoint(dimensions []*cloudwatch.Dimension, metricName string, sTyp statType) (float64, error) {
	now := time.Now()

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(120) * time.Second * -1)), // 2 min (to fetch at least 1 data-point)
		EndTime:    aws.Time(now),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String(sTyp.String())},
		Namespace:  aws.String("AWS/ELB"),
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
		switch sTyp {
		case stAve:
			latestVal = *dp.Average
		case stSum:
			latestVal = *dp.Sum
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch elb metrics
func (p ELBPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	// HostCount per AZ
	for _, az := range p.AZs {
		d := []*cloudwatch.Dimension{
			{
				Name:  aws.String("AvailabilityZone"),
				Value: az,
			},
		}
		if p.Lbname != "" {
			d2 := &cloudwatch.Dimension{
				Name:  aws.String("LoadBalancerName"),
				Value: aws.String(p.Lbname),
			}
			d = append(d, d2)
		}
		for _, met := range []string{"HealthyHostCount", "UnHealthyHostCount"} {
			v, err := p.getLastPoint(d, met, stAve)
			if err == nil {
				stat[met+"_"+*az] = v
			}
		}
	}

	glb := []*cloudwatch.Dimension{}
	if p.Lbname != "" {
		g2 := &cloudwatch.Dimension{
			Name:  aws.String("LoadBalancerName"),
			Value: aws.String(p.Lbname),
		}
		glb = append(glb, g2)
	}

	v, err := p.getLastPoint(glb, "Latency", stAve)
	if err == nil {
		stat["Latency"] = v
	}

	for _, met := range [...]string{"HTTPCode_Backend_2XX", "HTTPCode_Backend_3XX", "HTTPCode_Backend_4XX", "HTTPCode_Backend_5XX"} {
		v, err := p.getLastPoint(glb, met, stSum)
		if err == nil {
			stat[met] = v
		}
	}

	return stat, nil
}

// GraphDefinition for Mackerel
func (p ELBPlugin) GraphDefinition() map[string]mp.Graphs {
	for _, grp := range [...]string{"elb.healthy_host_count", "elb.unhealthy_host_count"} {
		var namePre string
		var label string
		switch grp {
		case "elb.healthy_host_count":
			namePre = "HealthyHostCount_"
			label = "ELB Healthy Host Count"
		case "elb.unhealthy_host_count":
			namePre = "UnHealthyHostCount_"
			label = "ELB Unhealthy Host Count"
		}

		var metrics []mp.Metrics
		for _, az := range p.AZs {
			metrics = append(metrics, mp.Metrics{Name: namePre + *az, Label: *az, Stacked: true})
		}
		graphdef[grp] = mp.Graphs{
			Label:   label,
			Unit:    "integer",
			Metrics: metrics,
		}
	}

	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optLbname := flag.String("lbname", "", "ELB Name")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var elb ELBPlugin

	if *optRegion == "" {
		ec2metadata := ec2metadata.New(session.New())
		if ec2metadata.Available() {
			elb.Region, _ = ec2metadata.Region()
		}
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

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

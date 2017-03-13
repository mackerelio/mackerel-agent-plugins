package mpawswaf

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/waf"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string]mp.Graphs{
	"waf.Requests.#": {
		Label: "AWS WAF Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "AllowedRequests", Label: "AllowedRequests"},
			{Name: "BlockedRequests", Label: "BlockedRequests"},
			{Name: "CountedRequests", Label: "CountedRequests"},
		},
	},
}

// WafPlugin mackerel plugin for aws waf
type WafPlugin struct {
	AccessKeyID     string
	SecretAccessKey string
	WebACLID        string
	WebACL          string
	Rules           []string
	CloudWatch      *cloudwatch.CloudWatch
}

func (p *WafPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	config = config.WithRegion("us-east-1")

	svc := waf.New(sess, config)
	response, err := svc.GetWebACL(&waf.GetWebACLInput{
		WebACLId: aws.String(p.WebACLID),
	})
	if err != nil {
		return err
	}
	p.WebACL = *response.WebACL.MetricName

	rules := []string{"ALL", "Default_Action"}
	for _, rule := range response.WebACL.Rules {
		response, err := svc.GetRule(&waf.GetRuleInput{
			RuleId: aws.String(*rule.RuleId),
		})
		if err != nil {
			continue
		}
		rules = append(rules, *response.Rule.MetricName)
	}
	p.Rules = rules

	p.CloudWatch = cloudwatch.New(sess, config)

	return nil
}

func (p WafPlugin) getLastPoint(dimensions []*cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("WAF"),
		MetricName: aws.String(metricName),
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)),
		EndTime:    aws.Time(now),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})

	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := time.Unix(0, 0)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(latest) {
			continue
		}

		latest = *dp.Timestamp
		latestVal = *dp.Sum
	}

	return latestVal, nil
}

// FetchMetrics interface for mackerelplugin
func (p WafPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	for _, rule := range p.Rules {
		dimensions := []*cloudwatch.Dimension{
			{
				Name:  aws.String("Rule"),
				Value: aws.String(rule),
			},
			{
				Name:  aws.String("WebACL"),
				Value: aws.String(p.WebACL),
			},
		}

		for _, met := range [...]string{"AllowedRequests", "BlockedRequests", "CountedRequests"} {
			v, err := p.getLastPoint(dimensions, met)
			if err == nil {
				stat[fmt.Sprintf("waf.Requests.%s.%s", rule, met)] = v
			} else {
				log.Printf("%s.%s: %s", rule, met, err)
			}
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p WafPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optWebACLID := flag.String("web-acl-id", "", "AWS Web ACL ID")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var waf WafPlugin

	waf.WebACLID = *optWebACLID
	waf.AccessKeyID = *optAccessKeyID
	waf.SecretAccessKey = *optSecretAccessKey

	err := waf.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(waf)
	helper.Tempfile = *optTempfile
	helper.Run()
}

package mpawsses

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string]mp.Graphs{
	"ses.send24h": {
		Label: "SES Send (last 24h)",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Max24HourSend", Label: "Max"},
			{Name: "SentLast24Hours", Label: "Sent"},
		},
	},
	"ses.max_send_rate": {
		Label: "SES Max Send Rate",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "MaxSendRate", Label: "MaxRate"},
		},
	},
	"ses.stats": {
		Label: "SES Stats",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Complaints", Label: "Complaints"},
			{Name: "DeliveryAttempts", Label: "DeliveryAttempts"},
			{Name: "Bounces", Label: "Bounces"},
			{Name: "Rejects", Label: "Rejects"},
		},
	},
}

// SESPlugin mackerel plugin for Amazon SES
type SESPlugin struct {
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string

	Svc *ses.SES
}

// prepare creates SES instance
func (p *SESPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()

	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" && p.Endpoint != "" {
		return errors.New("--region and --endpoint are exclusive")
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}
	// for backward compatibility
	if p.Endpoint != "" {
		u, err := url.Parse(p.Endpoint)
		if err != nil {
			return err
		}
		hosts := strings.Split(u.Host, ".")
		if len(hosts) != 4 || hosts[2] != "amazonaws" {
			return errors.New("--endpoint is invalid")
		}
		config = config.WithRegion(hosts[1])
	}

	p.Svc = ses.New(sess, config)

	return nil
}

// FetchMetrics interface for mackerel plugin
func (p SESPlugin) FetchMetrics() (map[string]float64, error) {

	stat := make(map[string]float64)
	quota, err := p.Svc.GetSendQuota(&ses.GetSendQuotaInput{})
	if err != nil {
		return nil, err
	}

	if quota.SentLast24Hours != nil {
		stat["SentLast24Hours"] = *quota.SentLast24Hours
	}

	if quota.Max24HourSend != nil {
		stat["Max24HourSend"] = *quota.Max24HourSend
	}

	if quota.MaxSendRate != nil {
		stat["MaxSendRate"] = *quota.MaxSendRate
	}

	result, err := p.Svc.GetSendStatistics(nil)
	if err == nil {
		t := time.Unix(0, 0)
		latest := &ses.SendDataPoint{
			Timestamp: &t,
		}

		datapoints := result.SendDataPoints

		if len(datapoints) > 0 {
			for _, dp := range datapoints {
				if latest.Timestamp.Before(*dp.Timestamp) {
					latest = dp
				}
			}

			stat["Complaints"] = float64(*latest.Complaints)
			stat["DeliveryAttempts"] = float64(*latest.DeliveryAttempts)
			stat["Bounces"] = float64(*latest.Bounces)
			stat["Rejects"] = float64(*latest.Rejects)
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerel plugin
func (p SESPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optEndpoint := flag.String("endpoint", "", "AWS Endpoint (deprecated)")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ses SESPlugin

	ses.Region = *optRegion
	ses.Endpoint = *optEndpoint
	ses.AccessKeyID = *optAccessKeyID
	ses.SecretAccessKey = *optSecretAccessKey

	err := ses.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(ses)
	helper.Tempfile = *optTempfile

	helper.Run()
}

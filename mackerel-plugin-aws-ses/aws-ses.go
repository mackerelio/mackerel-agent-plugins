package main

import (
	"errors"
	"flag"
	"time"

	"github.com/crowdmob/goamz/aws"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	ses "github.com/naokibtn/go-ses"
)

var graphdef = map[string](mp.Graphs){
	"send24h": mp.Graphs{
		Label: "SES Send (last 24h)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Max24HourSend", Label: "Max"},
			mp.Metrics{Name: "SentLast24Hours", Label: "Sent"},
		},
	},
	"max_send_rate": mp.Graphs{
		Label: "SES Max Send Rate",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "MaxSendRate", Label: "MaxRate"},
		},
	},
	"stats": mp.Graphs{
		Label: "SES Stats",
		Unit:  "int",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Complaints", Label: "Complaints"},
			mp.Metrics{Name: "DeliveryAttempts", Label: "DeliveryAttempts"},
			mp.Metrics{Name: "Bounces", Label: "Bounces"},
			mp.Metrics{Name: "Rejects", Label: "Rejects"},
		},
	},
}

// SESPlugin mackerel plugin for Amazon SES
type SESPlugin struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p SESPlugin) MetricKeyPrefix() string {
	return "ses"
}

// FetchMetrics interface for mackerel plugin
func (p SESPlugin) FetchMetrics() (map[string]interface{}, error) {
	if p.Endpoint == "" {
		return nil, errors.New("no endpoint")
	}

	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	sescfg := ses.Config{
		AccessKeyID:     auth.AccessKey,
		SecretAccessKey: auth.SecretKey,
		SecurityToken:   auth.Token(),
		Endpoint:        p.Endpoint,
	}

	stat := make(map[string]interface{})
	quota, err := sescfg.GetSendQuota()
	if err == nil {
		stat["SentLast24Hours"] = quota.SentLast24Hours
		stat["Max24HourSend"] = quota.Max24HourSend
		stat["MaxSendRate"] = quota.MaxSendRate
	}

	datapoints, err := sescfg.GetSendStatistics()
	if err == nil {
		latest := ses.SendDataPoint{
			Timestamp: time.Unix(0, 0),
		}

		for _, dp := range datapoints {
			if latest.Timestamp.Before(dp.Timestamp) {
				latest = dp
			}
		}

		stat["Complaints"] = float64(latest.Complaints)
		stat["DeliveryAttempts"] = float64(latest.DeliveryAttempts)
		stat["Bounces"] = float64(latest.Bounces)
		stat["Rejects"] = float64(latest.Rejects)
	}

	return stat, nil
}

// GraphDefinition interface for mackerel plugin
func (p SESPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optEndpoint := flag.String("endpoint", "", "AWS Endpoint")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ses SESPlugin

	ses.Endpoint = *optEndpoint
	ses.AccessKeyID = *optAccessKeyID
	ses.SecretAccessKey = *optSecretAccessKey

	helper := mp.NewMackerelPlugin(ses)
	helper.Tempfile = *optTempfile
	helper.Run()
}

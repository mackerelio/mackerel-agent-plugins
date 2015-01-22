package main

import (
	"errors"
	"flag"
	"github.com/crowdmob/goamz/aws"
	mp "github.com/mackerelio/go-mackerel-plugin"
	ses "github.com/naokibtn/go-ses"
	"os"
	"time"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"ses.send24h": mp.Graphs{
		Label: "SES Send (last 24h)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Max24HourSend", Label: "Max"},
			mp.Metrics{Name: "SentLast24Hours", Label: "Sent"},
		},
	},
	"ses.max_send_rate": mp.Graphs{
		Label: "SES Max Send Rate",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "MaxSendRate", Label: "MaxRate"},
		},
	},
	"ses.stats": mp.Graphs{
		Label: "SES Stats",
		Unit:  "int",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Complaints", Label: "Complaints"},
			mp.Metrics{Name: "DeliveryAttempts", Label: "DeliveryAttempts"},
			mp.Metrics{Name: "Bounces int", Label: "Bounces"},
			mp.Metrics{Name: "Rejects int", Label: "Rejects"},
		},
	},
}

type SESPlugin struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
}

func (p SESPlugin) FetchMetrics() (map[string]float64, error) {
	if p.Endpoint == "" {
		return nil, errors.New("no endpoint")
	}

	auth, err := aws.GetAuth(p.AccessKeyId, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	sescfg := ses.Config{
		AccessKeyID:     auth.AccessKey,
		SecretAccessKey: auth.SecretKey,
		SecurityToken:   auth.Token(),
		Endpoint:        p.Endpoint,
	}

	stat := make(map[string]float64)
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

func (p SESPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optEndpoint := flag.String("endpoint", "", "AWS Endpoint")
	optAccessKeyId := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ses SESPlugin

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

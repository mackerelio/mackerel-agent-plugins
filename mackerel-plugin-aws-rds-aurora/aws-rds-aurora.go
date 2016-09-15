package main

import (
	"flag"

	"github.com/crowdmob/goamz/aws"
)

// AuroraPlugin mackerel plugin for amazon Aurora
type AuroraPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Identifier      string
	Prefix          string
	LabelPrefix     string
	Metrics         []string
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "DB Instance Identifier")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Metric Label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var aurora AuroraPlugin

	if *optLabelPrefix == "" {
		aurora.LabelPrefix = "Aurora"
	} else {
		aurora.LabelPrefix = *optLabelPrefix
	}

	if *optRegion == "" {
		aurora.Region = aws.InstanceRegion()
	} else {
		aurora.Region = *optRegion
	}

	aurora.Identifier = *optIdentifier
	aurora.AccessKeyID = *optAccessKeyID
	aurora.SecretAccessKey = *optSecretAccessKey
	aurora.Metrics = metricsdefAurora
}

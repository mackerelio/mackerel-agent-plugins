package main

import (
	"./lib"
	"flag"
	"fmt"
)

func main() {
	now := time.Now()
	optDest := flag.String("dest", "", "Please Specify ServiceMetric or Host (※default Host)")
	optApiKey := flag.String("api-key", "", "API Key must have read and write authority")
	optServiceName := flag.String("service-name", "", "target serviceName")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optCurrency := flag.String("currency", "USD", "Unit of currency")
	optTarget := flag.String("target", "", "Target AWS Service. if no specific Aws Service, get metrics list from cloudwatch and draw all available metrics")
	flag.Parse()

	if now.Minute() == 0 {
		mpawsbilling.WriteCache(*optAccessKeyID, *optSecretAccessKey, *optCurrency, *optTarget)
	}

	if *optDest == "ServiceMetric" {
		mpawsbilling.SendServiceMetric(*optApiKey, *optServiceName)
	} else if *optDest == "Host" {
		mpawsbilling.OutputData()
	}
}

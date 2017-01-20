package mpawsbilling

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type AWSBilling struct {
	Region      string
	Currency    string
	Target      []string
	Credentials *credentials.Credentials
	CloudWatch  *cloudwatch.CloudWatch
}

func getLatestValue(cloudWatch *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, error) {
	now := time.Now()

	startTime := time.Unix(now.Unix()-86400, int64(now.Nanosecond()))

	statistics := []*string{aws.String("Maximum")}

	in := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(now),
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Period:     aws.Int64(3600),
		Statistics: statistics,
	}

	out, err := cloudWatch.GetMetricStatistics(in)

	if err != nil {
		return 0, nil
	}

	datapoints := out.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	var latest time.Time
	var latestIndex int

	for i, datapoint := range datapoints {
		if datapoint.Timestamp.After(latest) {
			latest = *datapoint.Timestamp
			latestIndex = i
		}
	}

	return *datapoints[latestIndex].Maximum, nil
}

func toJson(data map[string]float64) string {
	jsonStr := `{`
	var jsonInner []string
	for k, v := range data {
		jsonInner = append(jsonInner, []string{`"` + k + `"` + `:` + `"` + fmt.Sprint(v) + `"`}...)
	}
	jsonStr += strings.Join(jsonInner[:], ",") + `}`

	return jsonStr
}

func (p AWSBilling) WriteLatestValue() {
	data := make(map[string]float64)
	file, err := os.OpenFile("/tmp/aws-billing-cache", os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}

	defer func() {
		file.Close()
	}()

	baseDimension := []*cloudwatch.Dimension{&cloudwatch.Dimension{
		Name:  aws.String("Currency"),
		Value: aws.String(p.Currency),
	}}

	goroutines := len(p.Target)
	c := make(chan map[string]float64)
	for _, metricName := range p.Target {
		go func(s chan<- map[string]float64, metricName string) {
			var dimensions []*cloudwatch.Dimension
			if metricName == "All" {
				dimensions = baseDimension
			} else {
				dimensions = append(
					[]*cloudwatch.Dimension{baseDimension[0]},
					[]*cloudwatch.Dimension{&cloudwatch.Dimension{
						Name:  aws.String("ServiceName"),
						Value: aws.String(metricName),
					}}...)

			}

			v, err := getLatestValue(p.CloudWatch, dimensions)
			if err == nil {
				s <- map[string]float64{metricName: v}
			} else {
				log.Printf("%s: %s", metricName, err)
			}

		}(c, metricName)
	}

	for i := 0; i < goroutines; i++ {
		for k, v := range <-c {
			data[k] = v
		}
	}

	close(c)

	file.WriteString(toJson(data))
}

func getServiceNameList(metrics *cloudwatch.ListMetricsOutput) (target []string) {
	for _, metric := range metrics.Metrics {
		for _, dimension := range metric.Dimensions {
			if *dimension.Name == "ServiceName" {
				target = append(target, []string{*dimension.Value}...)
			}
		}
	}

	return target
}

func WriteCache(optAccessKeyID string, optSecretAccessKey string, optCurrency string, optTarget string) {

	var billing AWSBilling

	if optAccessKeyID != "" && optSecretAccessKey != "" {
		billing.Credentials = credentials.NewStaticCredentials(optAccessKeyID, optSecretAccessKey, "")
	}

	billing.Region = "us-east-1"

	billing.Currency = optCurrency

	billing.CloudWatch = cloudwatch.New(session.New(
		&aws.Config{
			Credentials: billing.Credentials,
			Region:      aws.String(billing.Region),
		}))

	var target []string
	if optTarget == "" {
		metrics, _ := billing.CloudWatch.ListMetrics(&cloudwatch.ListMetricsInput{Namespace: aws.String("AWS/Billing")})
		target = getServiceNameList(metrics)
	} else {
		target = strings.Split(optTarget, ",")
	}

	billing.Target = append([]string{"All"}, target...)

	billing.WriteLatestValue()
}

type BillingCachePlugin struct {
	Data map[string]interface{}
}

func (p BillingCachePlugin) GraphDefinition() map[string]mp.Graphs {
	metrics := func() []mp.Metrics {
		var metrics []mp.Metrics
		for target, _ := range p.Data {
			metrics = append(metrics, []mp.Metrics{mp.Metrics{Name: target, Label: target}}...)
		}

		return metrics
	}()

	return map[string]mp.Graphs{
		"billing.all": mp.Graphs{
			Label:   "AWSBilling",
			Unit:    "float",
			Metrics: metrics,
		},
	}
}

func (p BillingCachePlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	for k, v := range p.Data {
		f, _ := strconv.ParseFloat(v.(string), 64)
		stat[k] = f
	}

	return stat, nil
}

func readData() (interface{}, error) {
	data, err := ioutil.ReadFile(`/tmp/aws-billing-cache`)
	if err != nil {
		return nil, err
	}

	str := string(data)
	var f interface{}
	err = json.Unmarshal([]byte(str), &f)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func OutputData() {
	var billingCache BillingCachePlugin

	f, err := readData()

	if err != nil {
		panic("failed to parse /tmp/aws-billing-cache")
	}

	billingCache.Data = f.(map[string]interface{})

	helper := mp.NewMackerelPlugin(billingCache)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}

}

type MetricValue struct {
	Name  string  `json:"name"`
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

func SendServiceMetric(optApiKey string, optServiceName string) {

	f, err := readData()

	if err != nil {
		panic("failed to parse /tmp/aws-billing-cache")
	}

	mapObj := f.(map[string]interface{})

	var metricValues []MetricValue

	now := time.Now().Unix()
	for name, value := range mapObj {
		f64, _ := strconv.ParseFloat(value.(string), 64)
		metricValue := MetricValue{Name: name, Value: f64, Time: now}

		metricValues = append(metricValues, metricValue)
	}

	jsonStr, _ := json.Marshal(metricValues)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://mackerel.io/api/v0/services/%s/tsdb", optServiceName), bytes.NewBuffer([]byte(string(jsonStr))))
	req.Header.Add("X-Api-Key", optApiKey)
	req.Header.Set("Content-Type", "application/json")
	client.Do(req)
}

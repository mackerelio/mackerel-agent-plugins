package mpawsec2ebs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (p EBSPlugin) getLastPointGetMetricData(ctx context.Context, now time.Time, vol types.Volume, queries []cloudwatchTypes.MetricDataQuery) (map[string]float64, error) {
	period := metricPeriodDefault
	if tmp, ok := metricPeriodByVolumeType[vol.VolumeType]; ok {
		period = tmp
	}
	start := now.Add(time.Duration(period) * 3 * time.Second * -1)

	resp, err := p.CloudWatch.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
		StartTime:         &start,
		EndTime:           &now,
		MetricDataQueries: queries,
		ScanBy:            cloudwatchTypes.ScanByTimestampDescending,
	})
	if err != nil {
		return nil, err
	}

	var values = make(map[string]float64, 0)

	for _, res := range resp.MetricDataResults {
		// GetMetricData API Timestamps are sorted by newest first.
		// ScanBy: TimestampDescending
		for i := range res.Values {
			values[*res.Id] = res.Values[i]
			break
		}
	}

	return values, nil
}

func metricDataID(metricName string, statType cloudwatchTypes.Statistic) string {
	return fmt.Sprintf("v_%s_%s", metricName, statType)
}

func createMetricDataQuery(id string, vol types.Volume, metricName string, statType cloudwatchTypes.Statistic) cloudwatchTypes.MetricDataQuery {
	return cloudwatchTypes.MetricDataQuery{
		Id: aws.String(id),
		MetricStat: &cloudwatchTypes.MetricStat{
			Metric: &cloudwatchTypes.Metric{
				MetricName: &metricName,
				Dimensions: []cloudwatchTypes.Dimension{
					{
						Name:  aws.String("VolumeId"),
						Value: vol.VolumeId,
					},
				},
				Namespace: aws.String("AWS/EBS"),
			},
			Period: aws.Int32(aggregationPeriod),
			Stat:   aws.String(string(statType)),
		},
	}
}

func (p EBSPlugin) fetchMetrics_GetMetricData() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	now := time.Now()

	// Override when Nitro instance.
	if p.Hypervisor == types.InstanceTypeHypervisorNitro {
		for i := range cloudwatchdefsNitro {
			cloudwatchdefs[i] = cloudwatchdefsNitro[i]
		}
	}

	for _, vol := range p.Volumes {
		queryMap := make(map[string]bool) // for duplicate check
		var queries []cloudwatchTypes.MetricDataQuery

		var graphs []string
		if vol.VolumeType == types.VolumeTypeIo1 {
			graphs = io1Graphs
		} else {
			graphs = defaultGraphs
		}
		for _, graphName := range graphs {
			for _, metric := range graphdef[graphName].Metrics {
				metricKey := graphName + "." + metric.Name
				cloudwatchdef := cloudwatchdefs[metricKey]

				if id := metricDataID(cloudwatchdef.MetricName, cloudwatchdef.Statistics); !queryMap[id] {

					queryMap[id] = true
					queries = append(queries, createMetricDataQuery(id, vol, cloudwatchdef.MetricName, cloudwatchdef.Statistics))
				}
				if cloudwatchdef.Additional != nil {
					if id := metricDataID(cloudwatchdef.Additional.MetricName, cloudwatchdef.Additional.Statistics); !queryMap[id] {
						queryMap[id] = true
						queries = append(queries, createMetricDataQuery(id, vol, cloudwatchdef.Additional.MetricName, cloudwatchdef.Additional.Statistics))
					}
				}
			}
		}

		idValues, err := p.getLastPointGetMetricData(context.TODO(), now, vol, queries)
		if err != nil {
			return nil, err
		}

		volumeID := normalizeVolumeID(*vol.VolumeId)

		for _, graphName := range graphs {
			for _, metric := range graphdef[graphName].Metrics {
				metricKey := graphName + "." + metric.Name
				cloudwatchdef := cloudwatchdefs[metricKey]
				val, err := fetch_GetMetricData(idValues, vol, cloudwatchdef)
				if err != nil {
					if errors.Is(err, errNoDataPoint) {
						// nop
					} else {
						return nil, err
					}
				} else {
					stat[strings.ReplaceAll(metricKey, "#", volumeID)] = val
				}
			}
		}
	}

	return stat, nil
}

func fetch_GetMetricData(idValues map[string]float64, volume types.Volume, setting cloudWatchSetting) (float64, error) {
	val, ok := idValues[metricDataID(setting.MetricName, setting.Statistics)]
	if !ok {
		return 0, errNoDataPoint
	}

	if setting.Additional == nil {
		return setting.CalcFunc(val), nil
	}

	val2, ok2 := idValues[metricDataID(setting.Additional.MetricName, setting.Additional.Statistics)]
	if !ok2 {
		return 0, errNoDataPoint
	}

	return setting.Additional.CalcFunc(val, val2), nil
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var metricsdefMemcached = []string{
	"CPUUtilization", "SwapUsage", "FreeableMemory", "NetworkBytesIn", "NetworkBytesOut",
	"BytesUsedForCacheItems", "BytesReadIntoMemcached", "BytesWrittenOutFromMemcached", "CasBadval", "CasHits",
	"CasMisses", "CmdFlush", "CmdGet", "CmdSet", "CurrConnections",
	"CurrItems", "DecrHits", "DecrMisses", "DeleteHits", "DeleteMisses",
	"Evictions", "GetHits", "GetMisses", "IncrHits", "IncrMisses",
	"Reclaimed", "BytesUsedForHash", "CmdConfigGet", "CmdConfigSet", "CmdTouch",
	"CurrConfig", "EvictedUnfetched", "ExpiredUnfetched", "SlabsMoved", "TouchHits",
	"TouchMisses", "NewConnections", "NewItems", "UnusedMemory",
}

var metricsdefRedis = []string{
	"CPUUtilization", "SwapUsage", "FreeableMemory", "NetworkBytesIn", "NetworkBytesOut",
	"CurrConnections", "Evictions", "Reclaimed", "NewConnections", "BytesUsedForCache",
	"CacheHits", "CacheMisses", "ReplicationLag", "GetTypeCmds", "SetTypeCmds",
	"KeyBasedCmds", "StringBasedCmds", "HashBasedCmds", "ListBasedCmds", "SetBasedCmds",
	"SortedSetBasedCmds", "CurrItems",
}

var graphdefMemcached = map[string](mp.Graphs){
	"ecache.CPUUtilization": mp.Graphs{
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"ecache.SwapUsage": mp.Graphs{
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
		},
	},
	"ecache.FreeableMemory": mp.Graphs{
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory"},
		},
	},
	"ecache.NetworkTraffic": mp.Graphs{
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "NetworkBytesIn", Label: "NetworkBytesIn"},
			mp.Metrics{Name: "NetworkBytesOut", Label: "NetworkBytesOut"},
		},
	},
	"ecache.Command": mp.Graphs{
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CmdGet", Label: "CmdGet"},
			mp.Metrics{Name: "CmdSet", Label: "CmdSet"},
			mp.Metrics{Name: "CmdFlush", Label: "CmdFlush"},
			mp.Metrics{Name: "CmdTouch", Label: "CmdTouch"},
			mp.Metrics{Name: "CmdConfigGet", Label: "CmdConfigGet"},
			mp.Metrics{Name: "CmdConfigSet", Label: "CmdConfigSet"},
		},
	},
	"ecache.CacheHitAndMiss": mp.Graphs{
		Label: "ECache Hits/Misses",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "GetHits", Label: "GetHits"},
			mp.Metrics{Name: "GetMisses", Label: "GetMisses"},
			mp.Metrics{Name: "DeleteHits", Label: "DeleteHits"},
			mp.Metrics{Name: "DeleteMisses", Label: "DeleteMisses"},
			mp.Metrics{Name: "IncrHits", Label: "IncrHits"},
			mp.Metrics{Name: "IncrMisses", Label: "IncrMisses"},
			mp.Metrics{Name: "DecrHits", Label: "DecrHits"},
			mp.Metrics{Name: "DecrMisses", Label: "DecrMisses"},
			mp.Metrics{Name: "CasBadval", Label: "CasBadval"},
			mp.Metrics{Name: "CasHits", Label: "CasHits"},
			mp.Metrics{Name: "CasMisses", Label: "CasMisses"},
			mp.Metrics{Name: "TouchHits", Label: "TouchHits"},
			mp.Metrics{Name: "TouchMisses", Label: "TouchMisses"},
		},
	},
	"ecache.Evictions": mp.Graphs{
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Evictions", Label: "Evictions"},
		},
	},
	"ecache.Unfetched": mp.Graphs{
		Label: "ECache Unfetched",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "EvictedUnfetched", Label: "EvictedUnfetched"},
			mp.Metrics{Name: "ExpiredUnfetched", Label: "ExpiredUnfetched"},
		},
	},
	"ecache.Bytes": mp.Graphs{
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "BytesReadIntoMemcached", Label: "BytesReadIntoMemcached"},
			mp.Metrics{Name: "BytesWrittenOutFromMemcached", Label: "BytesWrittenOutFromMemcached"},
		},
	},
	"ecache.Connections": mp.Graphs{
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrConnections", Label: "CurrConnections"},
			mp.Metrics{Name: "NewConnections", Label: "NewConnections"},
		},
	},
	"ecache.Items": mp.Graphs{
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrItems", Label: "CurrItems"},
			mp.Metrics{Name: "NewItems", Label: "NewItems"},
			mp.Metrics{Name: "Reclaimed", Label: "Reclaimed"},
			mp.Metrics{Name: "CurrConfig", Label: "CurrConfig"},
			mp.Metrics{Name: "SlabsMoved", Label: "SlabsMoved"},
		},
	},
	"ecache.MemoryUsage": mp.Graphs{
		Label: "ECache Memory Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "UnusedMemory", Label: "UnusedMemory"},
			mp.Metrics{Name: "BytesUsedForHash", Label: "BytesUsedForHash"},
			mp.Metrics{Name: "BytesUsedForCacheItems", Label: "BytesUsedForCacheItems"},
		},
	},
}

var graphdefRedis = map[string](mp.Graphs){
	"ecache.CPUUtilization": mp.Graphs{
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"ecache.SwapUsage": mp.Graphs{
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
		},
	},
	"ecache.FreeableMemory": mp.Graphs{
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory"},
		},
	},
	"ecache.NetworkTraffic": mp.Graphs{
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "NetworkBytesIn", Label: "NetworkBytesIn"},
			mp.Metrics{Name: "NetworkBytesOut", Label: "NetworkBytesOut"},
		},
	},
	"ecache.Command": mp.Graphs{
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "GetTypeCmds", Label: "GetTypeCmds"},
			mp.Metrics{Name: "SetTypeCmds", Label: "SetTypeCmds"},
			mp.Metrics{Name: "KeyBasedCmds", Label: "KeyBasedCmds"},
			mp.Metrics{Name: "StringBasedCmds", Label: "StringBasedCmds"},
			mp.Metrics{Name: "HashBasedCmds", Label: "HashBasedCmds"},
			mp.Metrics{Name: "ListBasedCmds", Label: "ListBasedCmds"},
			mp.Metrics{Name: "SetBasedCmds", Label: "SetBasedCmds"},
			mp.Metrics{Name: "SortedSetBasedCmds", Label: "SortedSetBasedCmds"},
		},
	},
	"ecache.CacheHitAndMiss": mp.Graphs{
		Label: "ECache Hits/Misses",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CacheHits", Label: "CacheHits"},
			mp.Metrics{Name: "CacheMisses", Label: "CacheMisses"},
		},
	},
	"ecache.Evictions": mp.Graphs{
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Evictions", Label: "Evictions"},
		},
	},
	"ecache.Bytes": mp.Graphs{
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "BytesUsedForCache", Label: "BytesUsedForCache"},
		},
	},
	"ecache.Connections": mp.Graphs{
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrConnections", Label: "CurrConnections"},
			mp.Metrics{Name: "NewConnections", Label: "NewConnections"},
		},
	},
	"ecache.Items": mp.Graphs{
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrItems", Label: "CurrItems"},
			mp.Metrics{Name: "Reclaimed", Label: "Reclaimed"},
		},
	},
}

// ECachePlugin mackerel plugin for elasticache
type ECachePlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	CacheClusterID  string
	CacheNodeID     string
	ElastiCacheType string
	CacheMetrics    []string
}

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimensions *[]cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: *dimensions,
		StartTime:  now.Add(time.Duration(180) * time.Second * -1), // 3 mins (to fetch at least 1 data-point)
		EndTime:    now,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{"Average"},
		Namespace:  "AWS/ElastiCache",
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
		latestVal = dp.Average
	}

	return latestVal, nil
}

// FetchMetrics fetch elasticache values
func (p ECachePlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	cloudWatch, err := cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)

	perInstances := &[]cloudwatch.Dimension{
		{
			Name:  "CacheClusterId",
			Value: p.CacheClusterID,
		},
		{
			Name:  "CacheNodeId",
			Value: p.CacheNodeID,
		},
	}

	for _, met := range p.CacheMetrics {
		v, err := getLastPoint(cloudWatch, perInstances, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// GraphDefinition graph definition
func (p ECachePlugin) GraphDefinition() map[string](mp.Graphs) {
	switch p.ElastiCacheType {
	case "memcached":
		return graphdefMemcached
	case "redis":
		return graphdefRedis
	default:
		log.Printf("elasticache-type is 'memcached' or 'redis'.")
		return nil
	}
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optCacheClusterID := flag.String("cache-cluster-id", "", "Cache Cluster Id")
	optCacheNodeID := flag.String("cache-node-id", "0001", "Cache Node Id")
	optElastiCacheType := flag.String("elasticache-type", "", "ElastiCache type")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ecache ECachePlugin

	if *optRegion == "" {
		ecache.Region = aws.InstanceRegion()
	} else {
		ecache.Region = *optRegion
	}

	ecache.AccessKeyID = *optAccessKeyID
	ecache.SecretAccessKey = *optSecretAccessKey
	ecache.CacheClusterID = *optCacheClusterID
	ecache.CacheNodeID = *optCacheNodeID
	ecache.ElastiCacheType = *optElastiCacheType
	switch ecache.ElastiCacheType {
	case "memcached":
		ecache.CacheMetrics = metricsdefMemcached
	case "redis":
		ecache.CacheMetrics = metricsdefRedis
	default:
		log.Printf("elasticache-type is 'memcached' or 'redis'.")
		os.Exit(1)
	}

	helper := mp.NewMackerelPlugin(ecache)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-aws-elasticache-%s-%s", *optCacheClusterID, *optCacheNodeID)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

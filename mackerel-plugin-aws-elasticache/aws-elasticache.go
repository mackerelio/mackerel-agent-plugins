package main

import (
	"errors"
	"flag"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
	"os"
	"time"
)

var metricsdefMemcached []string = []string{
	"CPUUtilization", "SwapUsage", "FreeableMemory", "NetworkBytesIn", "NetworkBytesOut",
	"BytesUsedForCacheItems", "BytesReadIntoMemcached", "BytesWrittenOutFromMemcached", "CasBadval", "CasHits",
	"CasMisses", "CmdFlush", "CmdGet", "CmdSet", "CurrConnections",
	"CurrItems", "DecrHits", "DecrMisses", "DeleteHits", "DeleteMisses",
	"Evictions", "GetHits", "GetMisses", "IncrHits", "IncrMisses",
	"Reclaimed", "BytesUsedForHash", "CmdConfigGet", "CmdConfigSet", "CmdTouch",
	"CurrConfig", "EvictedUnfetched", "ExpiredUnfetched", "SlabsMoved", "TouchHits",
	"TouchMisses", "NewConnections", "NewItems", "UnusedMemory",
}

var metricsdefRedis []string = []string{
	"CPUUtilization", "SwapUsage", "FreeableMemory", "NetworkBytesIn", "NetworkBytesOut",
	"CurrConnections", "Evictions", "Reclaimed", "NewConnections", "BytesUsedForCache",
	"CacheHits", "CacheMisses", "ReplicationLag", "GetTypeCmds", "SetTypeCmds",
	"KeyBasedCmds", "StringBasedCmds", "HashBasedCmds", "ListBasedCmds", "SetBasedCmds",
	"SortedSetBasedCmds", "CurrItems",
}

var graphdefMemcached map[string](mp.Graphs) = map[string](mp.Graphs){
	"ecache.CPUUtilization": mp.Graphs{
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization", Diff: false},
		},
	},
	"ecache.SwapUsage": mp.Graphs{
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SwapUsage", Label: "SwapUsage", Diff: false},
		},
	},
	"ecache.FreeableMemory": mp.Graphs{
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory", Diff: false},
		},
	},
	"ecache.NetworkTraffic": mp.Graphs{
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "NetworkBytesIn", Label: "NetworkBytesIn", Diff: false},
			mp.Metrics{Name: "NetworkBytesOut", Label: "NetworkBytesOut", Diff: false},
		},
	},
	"ecache.Command": mp.Graphs{
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CmdGet", Label: "CmdGet", Diff: true},
			mp.Metrics{Name: "CmdSet", Label: "CmdSet", Diff: true},
			mp.Metrics{Name: "CmdFlush", Label: "CmdFlush", Diff: true},
			mp.Metrics{Name: "CmdTouch", Label: "CmdTouch", Diff: true},
			mp.Metrics{Name: "CmdConfigGet", Label: "CmdConfigGet", Diff: true},
			mp.Metrics{Name: "CmdConfigSet", Label: "CmdConfigSet", Diff: true},
		},
	},
	"ecache.CacheHitAndMiss": mp.Graphs{
		Label: "ECache Hits/Misses",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "GetHits", Label: "GetHits", Diff: true},
			mp.Metrics{Name: "GetMisses", Label: "GetMisses", Diff: true},
			mp.Metrics{Name: "DeleteHits", Label: "DeleteHits", Diff: true},
			mp.Metrics{Name: "DeleteMisses", Label: "DeleteMisses", Diff: true},
			mp.Metrics{Name: "IncrHits", Label: "IncrHits", Diff: true},
			mp.Metrics{Name: "IncrMisses", Label: "IncrMisses", Diff: true},
			mp.Metrics{Name: "DecrHits", Label: "DecrHits", Diff: true},
			mp.Metrics{Name: "DecrMisses", Label: "DecrMisses", Diff: true},
			mp.Metrics{Name: "CasBadval", Label: "CasBadval", Diff: true},
			mp.Metrics{Name: "CasHits", Label: "CasHits", Diff: true},
			mp.Metrics{Name: "CasMisses", Label: "CasMisses", Diff: true},
			mp.Metrics{Name: "TouchHits", Label: "TouchHits", Diff: true},
			mp.Metrics{Name: "TouchMisses", Label: "TouchMisses", Diff: true},
		},
	},
	"ecache.Evictions": mp.Graphs{
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Evictions", Label: "Evictions", Diff: true},
		},
	},
	"ecache.Unfetched": mp.Graphs{
		Label: "ECache Unfetched",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "EvictedUnfetched", Label: "EvictedUnfetched", Diff: true},
			mp.Metrics{Name: "ExpiredUnfetched", Label: "ExpiredUnfetched", Diff: true},
		},
	},
	"ecache.Bytes": mp.Graphs{
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "BytesReadIntoMemcached", Label: "BytesReadIntoMemcached", Diff: true},
			mp.Metrics{Name: "BytesWrittenOutFromMemcached", Label: "BytesWrittenOutFromMemcached", Diff: true},
		},
	},
	"ecache.Connections": mp.Graphs{
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrConnections", Label: "CurrConnections", Diff: false},
			mp.Metrics{Name: "NewConnections", Label: "NewConnections", Diff: false},
		},
	},
	"ecache.Items": mp.Graphs{
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrItems", Label: "CurrItems", Diff: false},
			mp.Metrics{Name: "NewItems", Label: "NewItems", Diff: false},
			mp.Metrics{Name: "Reclaimed", Label: "Reclaimed", Diff: false},
			mp.Metrics{Name: "CurrConfig", Label: "CurrConfig", Diff: false},
			mp.Metrics{Name: "SlabsMoved", Label: "SlabsMoved", Diff: false},
		},
	},
	"ecache.MemoryUsage": mp.Graphs{
		Label: "ECache Memory Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "UnusedMemory", Label: "UnusedMemory", Diff: false},
			mp.Metrics{Name: "BytesUsedForHash", Label: "BytesUsedForHash", Diff: false},
			mp.Metrics{Name: "BytesUsedForCacheItems", Label: "BytesUsedForCacheItems", Diff: false},
		},
	},
}

var graphdefRedis map[string](mp.Graphs) = map[string](mp.Graphs){
	"ecache.CPUUtilization": mp.Graphs{
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization", Diff: false},
		},
	},
	"ecache.SwapUsage": mp.Graphs{
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SwapUsage", Label: "SwapUsage", Diff: false},
		},
	},
	"ecache.FreeableMemory": mp.Graphs{
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory", Diff: false},
		},
	},
	"ecache.NetworkTraffic": mp.Graphs{
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "NetworkBytesIn", Label: "NetworkBytesIn", Diff: false},
			mp.Metrics{Name: "NetworkBytesOut", Label: "NetworkBytesOut", Diff: false},
		},
	},
	"ecache.Command": mp.Graphs{
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "GetTypeCmds", Label: "GetTypeCmds", Diff: true},
			mp.Metrics{Name: "SetTypeCmds", Label: "SetTypeCmds", Diff: true},
			mp.Metrics{Name: "KeyBasedCmds", Label: "KeyBasedCmds", Diff: true},
			mp.Metrics{Name: "StringBasedCmds", Label: "StringBasedCmds", Diff: true},
			mp.Metrics{Name: "HashBasedCmds", Label: "HashBasedCmds", Diff: true},
			mp.Metrics{Name: "ListBasedCmds", Label: "ListBasedCmds", Diff: true},
			mp.Metrics{Name: "SetBasedCmds", Label: "SetBasedCmds", Diff: true},
			mp.Metrics{Name: "SortedSetBasedCmds", Label: "SortedSetBasedCmds", Diff: true},
		},
	},
	"ecache.CacheHitAndMiss": mp.Graphs{
		Label: "ECache Hits/Misses",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CacheHits", Label: "CacheHits", Diff: true},
			mp.Metrics{Name: "CacheMisses", Label: "CacheMisses", Diff: true},
		},
	},
	"ecache.Evictions": mp.Graphs{
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Evictions", Label: "Evictions", Diff: true},
		},
	},
	"ecache.Bytes": mp.Graphs{
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "BytesUsedForCache", Label: "BytesUsedForCache", Diff: true},
		},
	},
	"ecache.Connections": mp.Graphs{
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrConnections", Label: "CurrConnections", Diff: false},
			mp.Metrics{Name: "NewConnections", Label: "NewConnections", Diff: false},
		},
	},
	"ecache.Items": mp.Graphs{
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CurrItems", Label: "CurrItems", Diff: false},
			mp.Metrics{Name: "Reclaimed", Label: "Reclaimed", Diff: false},
		},
	},
}

type ECachePlugin struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	CacheClusterId  string
	CacheNodeId     string
	ElastiCacheType string
	CacheMetrics    []string
}

func GetLastPoint(cloudWatch *cloudwatch.CloudWatch, dimensions *[]cloudwatch.Dimension, metricName string) (float64, error) {
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

func (p ECachePlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyId, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	cloudWatch, err := cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)

	perInstances := &[]cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  "CacheClusterId",
			Value: p.CacheClusterId,
		},
		cloudwatch.Dimension{
			Name:  "CacheNodeId",
			Value: p.CacheNodeId,
		},
	}

	for _, met := range p.CacheMetrics {
		v, err := GetLastPoint(cloudWatch, perInstances, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

func (p ECachePlugin) GraphDefinition() map[string](mp.Graphs) {
	if p.ElastiCacheType == "memcached" {
		return graphdefMemcached
	} else {
		return graphdefRedis
	}
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyId := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optCacheClusterId := flag.String("cache-cluster-id", "", "Cache Cluster Id")
	optCacheNodeId := flag.String("cache-node-id", "0001", "Cache Node Id")
	optElastiCacheType := flag.String("elasticache-type", "", "ElastiCache type")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ecache ECachePlugin

	if *optRegion == "" {
		ecache.Region = aws.InstanceRegion()
	} else {
		ecache.Region = *optRegion
	}

	ecache.AccessKeyId = *optAccessKeyId
	ecache.SecretAccessKey = *optSecretAccessKey
	ecache.CacheClusterId = *optCacheClusterId
	ecache.CacheNodeId = *optCacheNodeId
	ecache.ElastiCacheType = *optElastiCacheType
	if ecache.ElastiCacheType == "memcached" {
		ecache.CacheMetrics = metricsdefMemcached
	} else {
		ecache.CacheMetrics = metricsdefRedis
	}

	helper := mp.NewMackerelPlugin(ecache)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-aws-elasticache"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

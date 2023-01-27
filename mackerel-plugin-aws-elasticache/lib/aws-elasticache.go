package mpawselasticache

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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

var graphdefMemcached = map[string]mp.Graphs{
	"ecache.CPUUtilization": {
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"ecache.SwapUsage": {
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "SwapUsage", Label: "SwapUsage"},
		},
	},
	"ecache.FreeableMemory": {
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "FreeableMemory", Label: "FreeableMemory"},
		},
	},
	"ecache.NetworkTraffic": {
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NetworkBytesIn", Label: "NetworkBytesIn"},
			{Name: "NetworkBytesOut", Label: "NetworkBytesOut"},
		},
	},
	"ecache.Command": {
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "CmdGet", Label: "CmdGet"},
			{Name: "CmdSet", Label: "CmdSet"},
			{Name: "CmdFlush", Label: "CmdFlush"},
			{Name: "CmdTouch", Label: "CmdTouch"},
			{Name: "CmdConfigGet", Label: "CmdConfigGet"},
			{Name: "CmdConfigSet", Label: "CmdConfigSet"},
		},
	},
	"ecache.CacheHitAndMiss": {
		Label: "ECache Hits/Misses",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "GetHits", Label: "GetHits"},
			{Name: "GetMisses", Label: "GetMisses"},
			{Name: "DeleteHits", Label: "DeleteHits"},
			{Name: "DeleteMisses", Label: "DeleteMisses"},
			{Name: "IncrHits", Label: "IncrHits"},
			{Name: "IncrMisses", Label: "IncrMisses"},
			{Name: "DecrHits", Label: "DecrHits"},
			{Name: "DecrMisses", Label: "DecrMisses"},
			{Name: "CasBadval", Label: "CasBadval"},
			{Name: "CasHits", Label: "CasHits"},
			{Name: "CasMisses", Label: "CasMisses"},
			{Name: "TouchHits", Label: "TouchHits"},
			{Name: "TouchMisses", Label: "TouchMisses"},
		},
	},
	"ecache.Evictions": {
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Evictions", Label: "Evictions"},
		},
	},
	"ecache.Unfetched": {
		Label: "ECache Unfetched",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "EvictedUnfetched", Label: "EvictedUnfetched"},
			{Name: "ExpiredUnfetched", Label: "ExpiredUnfetched"},
		},
	},
	"ecache.Bytes": {
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "BytesReadIntoMemcached", Label: "BytesReadIntoMemcached"},
			{Name: "BytesWrittenOutFromMemcached", Label: "BytesWrittenOutFromMemcached"},
		},
	},
	"ecache.Connections": {
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "CurrConnections", Label: "CurrConnections"},
			{Name: "NewConnections", Label: "NewConnections"},
		},
	},
	"ecache.Items": {
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "CurrItems", Label: "CurrItems"},
			{Name: "NewItems", Label: "NewItems"},
			{Name: "Reclaimed", Label: "Reclaimed"},
			{Name: "CurrConfig", Label: "CurrConfig"},
			{Name: "SlabsMoved", Label: "SlabsMoved"},
		},
	},
	"ecache.MemoryUsage": {
		Label: "ECache Memory Usage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "UnusedMemory", Label: "UnusedMemory"},
			{Name: "BytesUsedForHash", Label: "BytesUsedForHash"},
			{Name: "BytesUsedForCacheItems", Label: "BytesUsedForCacheItems"},
		},
	},
}

var graphdefRedis = map[string]mp.Graphs{
	"ecache.CPUUtilization": {
		Label: "ECache CPU Utilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"ecache.SwapUsage": {
		Label: "ECache Swap Usage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "SwapUsage", Label: "SwapUsage"},
		},
	},
	"ecache.FreeableMemory": {
		Label: "ECache Freeable Memory",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "FreeableMemory", Label: "FreeableMemory"},
		},
	},
	"ecache.NetworkTraffic": {
		Label: "ECache Network Traffic",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NetworkBytesIn", Label: "NetworkBytesIn"},
			{Name: "NetworkBytesOut", Label: "NetworkBytesOut"},
		},
	},
	"ecache.Command": {
		Label: "ECache Command",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "GetTypeCmds", Label: "GetTypeCmds"},
			{Name: "SetTypeCmds", Label: "SetTypeCmds"},
			{Name: "KeyBasedCmds", Label: "KeyBasedCmds"},
			{Name: "StringBasedCmds", Label: "StringBasedCmds"},
			{Name: "HashBasedCmds", Label: "HashBasedCmds"},
			{Name: "ListBasedCmds", Label: "ListBasedCmds"},
			{Name: "SetBasedCmds", Label: "SetBasedCmds"},
			{Name: "SortedSetBasedCmds", Label: "SortedSetBasedCmds"},
		},
	},
	"ecache.CacheHitAndMiss": {
		Label: "ECache Hits/Misses",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "CacheHits", Label: "CacheHits"},
			{Name: "CacheMisses", Label: "CacheMisses"},
		},
	},
	"ecache.Evictions": {
		Label: "ECache Evictions",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Evictions", Label: "Evictions"},
		},
	},
	"ecache.Bytes": {
		Label: "ECache Traffics",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "BytesUsedForCache", Label: "BytesUsedForCache"},
		},
	},
	"ecache.Connections": {
		Label: "ECache Connections",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "CurrConnections", Label: "CurrConnections"},
			{Name: "NewConnections", Label: "NewConnections"},
		},
	},
	"ecache.Items": {
		Label: "ECache Items",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "CurrItems", Label: "CurrItems"},
			{Name: "Reclaimed", Label: "Reclaimed"},
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

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)), // 3 mins (to fetch at least 1 data-point)
		EndTime:    aws.Time(now),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
		Namespace:  aws.String("AWS/ElastiCache"),
	})
	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := new(time.Time)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(*latest) {
			continue
		}

		latest = dp.Timestamp
		latestVal = *dp.Average
	}

	return latestVal, nil
}

// FetchMetrics fetch elasticache values
func (p ECachePlugin) FetchMetrics() (map[string]float64, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	cloudWatch := cloudwatch.New(sess, config)

	stat := make(map[string]float64)

	perInstances := []*cloudwatch.Dimension{
		{
			Name:  aws.String("CacheClusterId"),
			Value: aws.String(p.CacheClusterID),
		},
		{
			Name:  aws.String("CacheNodeId"),
			Value: aws.String(p.CacheNodeID),
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
func (p ECachePlugin) GraphDefinition() map[string]mp.Graphs {
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

// Do the plugin
func Do() {
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
		sess, err := session.NewSession()
		if err != nil {
			log.Fatalln(err)
		}
		ec2metadata := ec2metadata.New(sess)
		if ec2metadata.Available() {
			ecache.Region, _ = ec2metadata.Region()
		}
	} else {
		ecache.Region = *optRegion
	}

	ecache.Region = *optRegion
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
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-aws-elasticache-%s-%s", *optCacheClusterID, *optCacheNodeID))
	}

	helper.Run()
}

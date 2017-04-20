package mpawskinesisstreams

import (
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

func ExampleKinesisStreamsPlugin_GraphDefinition() {
	var kinesis KinesisStreamsPlugin
	kinesis.Prefix = "my-stream"
	helper := mp.NewMackerelPlugin(kinesis)
	helper.OutputDefinitions()

	// Output:
	// # mackerel-agent-plugin
	// {"graphs":{"my-stream.bytes":{"label":"My Stream Bytes","unit":"integer","metrics":[{"name":"GetRecordsBytes","label":"GetRecords","stacked":false},{"name":"IncomingBytes","label":"Total Incoming","stacked":false},{"name":"PutRecordBytes","label":"PutRecord","stacked":false},{"name":"PutRecordsBytes","label":"PutRecords","stacked":false}]},"my-stream.iteratorage":{"label":"My Stream Read Delay","unit":"integer","metrics":[{"name":"GetRecordsDelayAverageMilliseconds","label":"Average","stacked":false},{"name":"GetRecordsDelayMaxMilliseconds","label":"Max","stacked":false},{"name":"GetRecordsDelayMinMilliseconds","label":"Min","stacked":false}]},"my-stream.latency":{"label":"My Stream Operation Latency","unit":"integer","metrics":[{"name":"GetRecordsLatency","label":"GetRecords","stacked":false},{"name":"PutRecordLatency","label":"PutRecord","stacked":false},{"name":"PutRecordsLatency","label":"PutRecords","stacked":false}]},"my-stream.pending":{"label":"My Stream Pending Operations","unit":"integer","metrics":[{"name":"ReadThroughputExceeded","label":"Read","stacked":false},{"name":"WriteThroughputExceeded","label":"Write","stacked":false}]},"my-stream.records":{"label":"My Stream Records","unit":"integer","metrics":[{"name":"GetRecordsRecords","label":"GetRecords","stacked":false},{"name":"IncomingRecords","label":"Total Incoming","stacked":false},{"name":"PutRecordsRecords","label":"PutRecords","stacked":false}]},"my-stream.success":{"label":"My Stream Operation Success","unit":"integer","metrics":[{"name":"GetRecordsSuccess","label":"GetRecords","stacked":false},{"name":"PutRecordSuccess","label":"PutRecord","stacked":false},{"name":"PutRecordsSuccess","label":"PutRecords","stacked":false}]}}}
}

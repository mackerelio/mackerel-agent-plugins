package mpmongodb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var mongodb MongoDBPlugin

	graphdef := mongodb.GraphDefinition()
	if len(graphdef) != 4 {
		t.Errorf("GetTempfilename: %d should be 4", len(graphdef))
	}
}

func TestParse22(t *testing.T) {
	var mongodb MongoDBPlugin
	stub2_2_7 := `
{"asserts":{"msg":0,"regular":0,"rollovers":0,"user":0,"warning":0},"backgroundFlushing":{"average_ms":0,"flushes":0,"last_finished":"1970-01-01T09:00:00+09:00","last_ms":0,"total_ms":0},"connections":{"available":818,"current":1},"cursors":{"clientCursors_size":0,"timedOut":0,"totalOpen":0},"dur":{"commits":30,"commitsInWriteLock":0,"compression":0,"earlyCommits":0,"journaledMB":0,"timeMs":{"dt":3074,"prepLogBuffer":0,"remapPrivateView":0,"writeToDataFiles":0,"writeToJournal":0},"writeToDataFilesMB":0},"extra_info":{"heap_usage_bytes":25585584,"note":"fields vary by platform","page_faults":136},"globalLock":{"activeClients":{"readers":0,"total":0,"writers":0},"currentQueue":{"readers":0,"total":0,"writers":0},"lockTime":1638,"totalTime":35489000},"host":"58a1c98acba3","indexCounters":{"btree":{"accesses":0,"hits":5,"missRatio":0,"misses":0,"resets":0}},"localTime":"2015-08-17T15:08:02.677+09:00","locks":{".":{"timeAcquiringMicros":{"R":1593,"W":279},"timeLockedMicros":{"R":1906,"W":1638}},"admin":{"timeAcquiringMicros":{},"timeLockedMicros":{}},"local":{"timeAcquiringMicros":{"r":9,"w":0},"timeLockedMicros":{"r":44,"w":0}}},"mem":{"bits":64,"mapped":0,"mappedWithJournal":0,"resident":30,"supported":true,"virtual":128},"network":{"bytesIn":510,"bytesOut":2319,"numRequests":9},"ok":1,"opcounters":{"command":10,"delete":0,"getmore":0,"insert":0,"query":0,"update":0},"pid":1,"process":"mongod","recordStats":{"accessesNotInMemory":0,"local":{"accessesNotInMemory":0,"pageFaultExceptionsThrown":0},"pageFaultExceptionsThrown":0},"uptime":35,"uptimeEstimate":34,"uptimeMillis":35489,"version":"2.2.7","writeBacksQueued":false}
`

	var v interface{}
	err := json.Unmarshal([]byte(stub2_2_7), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["btree_hits"]).String(), "float64")
	assert.EqualValues(t, stat["btree_hits"], 5.0)
}

func TestParse24(t *testing.T) {
	var mongodb MongoDBPlugin
	stub2_4_14 := `
{"asserts":{"msg":0,"regular":0,"rollovers":0,"user":0,"warning":0},"backgroundFlushing":{"average_ms":0.6153846153846154,"flushes":26,"last_finished":"2015-08-17T14:55:58.622+09:00","last_ms":0,"total_ms":16},"connections":{"available":818,"current":1,"totalCreated":10},"cursors":{"clientCursors_size":0,"timedOut":0,"totalOpen":0},"dur":{"commits":30,"commitsInWriteLock":0,"compression":0,"earlyCommits":0,"journaledMB":0,"timeMs":{"dt":3074,"prepLogBuffer":0,"remapPrivateView":0,"writeToDataFiles":0,"writeToJournal":0},"writeToDataFilesMB":0},"extra_info":{"heap_usage_bytes":62256840,"note":"fields vary by platform","page_faults":181},"globalLock":{"activeClients":{"readers":0,"total":0,"writers":0},"currentQueue":{"readers":0,"total":0,"writers":0},"lockTime":143869,"totalTime":1603601000},"host":"bcd5355930ff","indexCounters":{"accesses":0,"hits":5,"missRatio":0,"misses":0,"resets":0},"localTime":"2015-08-17T14:56:42.209+09:00","locks":{".":{"timeAcquiringMicros":{"R":66884,"W":12244},"timeLockedMicros":{"R":86058,"W":143869}},"admin":{"timeAcquiringMicros":{},"timeLockedMicros":{}},"local":{"timeAcquiringMicros":{"r":513,"w":0},"timeLockedMicros":{"r":11886,"w":0}}},"mem":{"bits":64,"mapped":80,"mappedWithJournal":160,"resident":38,"supported":true,"virtual":341},"metrics":{"document":{"deleted":0,"inserted":1,"returned":0,"updated":0},"getLastError":{"wtime":{"num":0,"totalMillis":0},"wtimeouts":0},"operation":{"fastmod":0,"idhack":0,"scanAndOrder":0},"queryExecutor":{"scanned":0},"record":{"moves":0},"repl":{"apply":{"batches":{"num":0,"totalMillis":0},"ops":0},"buffer":{"count":0,"maxSizeBytes":268435456,"sizeBytes":0},"network":{"bytes":0,"getmores":{"num":0,"totalMillis":0},"ops":0,"readersCreated":0},"oplog":{"insert":{"num":0,"totalMillis":0},"insertBytes":0},"preload":{"docs":{"num":0,"totalMillis":0},"indexes":{"num":0,"totalMillis":0}}},"ttl":{"deletedDocuments":0,"passes":26}},"network":{"bytesIn":1940,"bytesOut":18064,"numRequests":33},"ok":1,"opcounters":{"command":35,"delete":0,"getmore":0,"insert":1,"query":26,"update":0},"opcountersRepl":{"command":0,"delete":0,"getmore":0,"insert":0,"query":0,"update":0},"pid":1,"process":"mongod","recordStats":{"accessesNotInMemory":0,"local":{"accessesNotInMemory":0,"pageFaultExceptionsThrown":0},"pageFaultExceptionsThrown":0},"uptime":1604,"uptimeEstimate":1581,"uptimeMillis":1603600,"version":"2.4.14","writeBacksQueued":false}
`

	var v interface{}
	err := json.Unmarshal([]byte(stub2_4_14), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["btree_hits"]).String(), "float64")
	assert.EqualValues(t, stat["btree_hits"], 5.0)
}

func TestParse26(t *testing.T) {
	var mongodb MongoDBPlugin
	stub2_6_11 := `
{"asserts":{"msg":0,"regular":0,"rollovers":0,"user":0,"warning":0},"backgroundFlushing":{"average_ms":0,"flushes":0,"last_finished":"1970-01-01T09:00:00+09:00","last_ms":0,"total_ms":0},"connections":{"available":818,"current":1,"totalCreated":1},"cursors":{"clientCursors_size":0,"note":"deprecated, use server status metrics","pinned":0,"timedOut":0,"totalNoTimeout":0,"totalOpen":0},"dur":{"commits":30,"commitsInWriteLock":0,"compression":0,"earlyCommits":0,"journaledMB":0,"timeMs":{"dt":3074,"prepLogBuffer":0,"remapPrivateView":0,"writeToDataFiles":0,"writeToJournal":0},"writeToDataFilesMB":0},"extra_info":{"heap_usage_bytes":62512008,"note":"fields vary by platform","page_faults":228},"globalLock":{"activeClients":{"readers":0,"total":0,"writers":0},"currentQueue":{"readers":0,"total":0,"writers":0},"lockTime":68622,"totalTime":6583000},"host":"08ea07b5a8fd","indexCounters":{"accesses":2,"hits":5,"missRatio":0,"misses":0,"resets":0},"localTime":"2015-08-17T15:08:44.187+09:00","locks":{".":{"timeAcquiringMicros":{"R":254,"W":94},"timeLockedMicros":{"R":520,"W":68622}},"admin":{"timeAcquiringMicros":{"r":28,"w":0},"timeLockedMicros":{"r":338,"w":0}},"local":{"timeAcquiringMicros":{"r":22,"w":0},"timeLockedMicros":{"r":46,"w":0}}},"mem":{"bits":64,"mapped":80,"mappedWithJournal":160,"resident":36,"supported":true,"virtual":342},"metrics":{"cursor":{"open":{"noTimeout":0,"pinned":0,"total":0},"timedOut":0},"document":{"deleted":0,"inserted":1,"returned":0,"updated":0},"getLastError":{"wtime":{"num":0,"totalMillis":0},"wtimeouts":0},"operation":{"fastmod":0,"idhack":0,"scanAndOrder":0},"queryExecutor":{"scanned":0,"scannedObjects":0},"record":{"moves":0},"repl":{"apply":{"batches":{"num":0,"totalMillis":0},"ops":0},"buffer":{"count":0,"maxSizeBytes":268435456,"sizeBytes":0},"network":{"bytes":0,"getmores":{"num":0,"totalMillis":0},"ops":0,"readersCreated":0},"preload":{"docs":{"num":0,"totalMillis":0},"indexes":{"num":0,"totalMillis":0}}},"storage":{"freelist":{"search":{"bucketExhausted":0,"requests":6,"scanned":11}}},"ttl":{"deletedDocuments":0,"passes":0}},"network":{"bytesIn":224,"bytesOut":380,"numRequests":4},"ok":1,"opcounters":{"command":6,"delete":0,"getmore":0,"insert":1,"query":1,"update":0},"opcountersRepl":{"command":0,"delete":0,"getmore":0,"insert":0,"query":0,"update":0},"pid":1,"process":"mongod","recordStats":{"accessesNotInMemory":0,"admin":{"accessesNotInMemory":0,"pageFaultExceptionsThrown":0},"local":{"accessesNotInMemory":0,"pageFaultExceptionsThrown":0},"pageFaultExceptionsThrown":0},"uptime":7,"uptimeEstimate":6,"uptimeMillis":6581,"version":"2.6.11","writeBacksQueued":false}
`

	var v interface{}
	err := json.Unmarshal([]byte(stub2_6_11), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["btree_hits"]).String(), "float64")
	assert.EqualValues(t, stat["btree_hits"], 5.0)
}

func TestParse30(t *testing.T) {
	var mongodb MongoDBPlugin
	stub3_0_5 := `
{"asserts":{"msg":0,"regular":0,"rollovers":0,"user":0,"warning":0},"backgroundFlushing":{"average_ms":0,"flushes":0,"last_finished":"1970-01-01T09:00:00+09:00","last_ms":0,"total_ms":3},"connections":{"available":818,"current":1,"totalCreated":1},"cursors":{"clientCursors_size":0,"note":"deprecated, use server status metrics","pinned":0,"timedOut":0,"totalNoTimeout":0,"totalOpen":0},"dur":{"commits":31,"commitsInWriteLock":0,"compression":1.1555931725208068,"earlyCommits":0,"journaledMB":0.024576,"timeMs":{"commits":1,"commitsInWriteLock":0,"dt":3037,"prepLogBuffer":0,"remapPrivateView":0,"writeToDataFiles":0,"writeToJournal":3},"writeToDataFilesMB":0.020941},"extra_info":{"heap_usage_bytes":62891464,"note":"fields vary by platform","page_faults":189},"globalLock":{"activeClients":{"readers":0,"total":9,"writers":0},"currentQueue":{"readers":0,"total":0,"writers":0},"totalTime":5609000},"host":"db625bac64b5","localTime":"2015-08-17T15:09:12.821+09:00","locks":{"Collection":{"acquireCount":{"R":7}},"Database":{"acquireCount":{"W":2,"r":7}},"Global":{"acquireCount":{"W":5,"r":21,"w":2}},"MMAPV1Journal":{"acquireCount":{"R":54,"r":5,"w":10}},"Metadata":{"acquireCount":{"W":4}}},"mem":{"bits":64,"mapped":80,"mappedWithJournal":160,"resident":52,"supported":true,"virtual":358},"metrics":{"commands":{"\u003cUNKNOWN\u003e":0,"_getUserCacheGeneration":{"failed":0,"total":0},"_isSelf":{"failed":0,"total":0},"_mergeAuthzCollections":{"failed":0,"total":0},"_migrateClone":{"failed":0,"total":0},"_recvChunkAbort":{"failed":0,"total":0},"_recvChunkCommit":{"failed":0,"total":0},"_recvChunkStart":{"failed":0,"total":0},"_recvChunkStatus":{"failed":0,"total":0},"_transferMods":{"failed":0,"total":0},"aggregate":{"failed":0,"total":0},"appendOplogNote":{"failed":0,"total":0},"applyOps":{"failed":0,"total":0},"authSchemaUpgrade":{"failed":0,"total":0},"authenticate":{"failed":0,"total":0},"availableQueryOptions":{"failed":0,"total":0},"buildInfo":{"failed":0,"total":0},"checkShardingIndex":{"failed":0,"total":0},"cleanupOrphaned":{"failed":0,"total":0},"clone":{"failed":0,"total":0},"cloneCollection":{"failed":0,"total":0},"cloneCollectionAsCapped":{"failed":0,"total":0},"collMod":{"failed":0,"total":0},"collStats":{"failed":0,"total":0},"compact":{"failed":0,"total":0},"connPoolStats":{"failed":0,"total":0},"connPoolSync":{"failed":0,"total":0},"connectionStatus":{"failed":0,"total":0},"convertToCapped":{"failed":0,"total":0},"copydb":{"failed":0,"total":0},"copydbgetnonce":{"failed":0,"total":0},"copydbsaslstart":{"failed":0,"total":0},"count":{"failed":0,"total":0},"create":{"failed":0,"total":0},"createIndexes":{"failed":0,"total":0},"createRole":{"failed":0,"total":0},"createUser":{"failed":0,"total":0},"currentOpCtx":{"failed":0,"total":0},"cursorInfo":{"failed":0,"total":0},"dataSize":{"failed":0,"total":0},"dbHash":{"failed":0,"total":0},"dbStats":{"failed":0,"total":0},"delete":{"failed":0,"total":0},"diagLogging":{"failed":0,"total":0},"distinct":{"failed":0,"total":0},"driverOIDTest":{"failed":0,"total":0},"drop":{"failed":0,"total":0},"dropAllRolesFromDatabase":{"failed":0,"total":0},"dropAllUsersFromDatabase":{"failed":0,"total":0},"dropDatabase":{"failed":0,"total":0},"dropIndexes":{"failed":0,"total":0},"dropRole":{"failed":0,"total":0},"dropUser":{"failed":0,"total":0},"eval":{"failed":0,"total":0},"explain":{"failed":0,"total":0},"features":{"failed":0,"total":0},"filemd5":{"failed":0,"total":0},"find":{"failed":0,"total":0},"findAndModify":{"failed":0,"total":0},"forceerror":{"failed":0,"total":0},"fsync":{"failed":0,"total":0},"geoNear":{"failed":0,"total":0},"geoSearch":{"failed":0,"total":0},"getCmdLineOpts":{"failed":0,"total":0},"getLastError":{"failed":0,"total":0},"getLog":{"failed":0,"total":0},"getParameter":{"failed":0,"total":0},"getPrevError":{"failed":0,"total":0},"getShardMap":{"failed":0,"total":0},"getShardVersion":{"failed":0,"total":0},"getnonce":{"failed":0,"total":1},"grantPrivilegesToRole":{"failed":0,"total":0},"grantRolesToRole":{"failed":0,"total":0},"grantRolesToUser":{"failed":0,"total":0},"group":{"failed":0,"total":0},"handshake":{"failed":0,"total":0},"hostInfo":{"failed":0,"total":0},"insert":{"failed":0,"total":0},"invalidateUserCache":{"failed":0,"total":0},"isMaster":{"failed":0,"total":1},"listCollections":{"failed":0,"total":0},"listCommands":{"failed":0,"total":0},"listDatabases":{"failed":0,"total":0},"listIndexes":{"failed":0,"total":0},"logRotate":{"failed":0,"total":0},"logout":{"failed":0,"total":0},"mapReduce":{"failed":0,"total":0},"mapreduce":{"shardedfinish":{"failed":0,"total":0}},"medianKey":{"failed":0,"total":0},"mergeChunks":{"failed":0,"total":0},"moveChunk":{"failed":0,"total":0},"parallelCollectionScan":{"failed":0,"total":0},"ping":{"failed":0,"total":2},"planCacheClear":{"failed":0,"total":0},"planCacheClearFilters":{"failed":0,"total":0},"planCacheListFilters":{"failed":0,"total":0},"planCacheListPlans":{"failed":0,"total":0},"planCacheListQueryShapes":{"failed":0,"total":0},"planCacheSetFilter":{"failed":0,"total":0},"profile":{"failed":0,"total":0},"reIndex":{"failed":0,"total":0},"renameCollection":{"failed":0,"total":0},"repairCursor":{"failed":0,"total":0},"repairDatabase":{"failed":0,"total":0},"replSetElect":{"failed":0,"total":0},"replSetFreeze":{"failed":0,"total":0},"replSetFresh":{"failed":0,"total":0},"replSetGetConfig":{"failed":0,"total":0},"replSetGetRBID":{"failed":0,"total":0},"replSetGetStatus":{"failed":0,"total":0},"replSetHeartbeat":{"failed":0,"total":0},"replSetInitiate":{"failed":0,"total":0},"replSetMaintenance":{"failed":0,"total":0},"replSetReconfig":{"failed":0,"total":0},"replSetStepDown":{"failed":0,"total":0},"replSetSyncFrom":{"failed":0,"total":0},"replSetUpdatePosition":{"failed":0,"total":0},"resetError":{"failed":0,"total":0},"resync":{"failed":0,"total":0},"revokePrivilegesFromRole":{"failed":0,"total":0},"revokeRolesFromRole":{"failed":0,"total":0},"revokeRolesFromUser":{"failed":0,"total":0},"rolesInfo":{"failed":0,"total":0},"saslContinue":{"failed":0,"total":0},"saslStart":{"failed":0,"total":0},"serverStatus":{"failed":0,"total":1},"setParameter":{"failed":0,"total":0},"setShardVersion":{"failed":0,"total":0},"shardConnPoolStats":{"failed":0,"total":0},"shardingState":{"failed":0,"total":0},"shutdown":{"failed":0,"total":0},"splitChunk":{"failed":0,"total":0},"splitVector":{"failed":0,"total":0},"top":{"failed":0,"total":0},"touch":{"failed":0,"total":0},"unsetSharding":{"failed":0,"total":0},"update":{"failed":0,"total":0},"updateRole":{"failed":0,"total":0},"updateUser":{"failed":0,"total":0},"usersInfo":{"failed":0,"total":0},"validate":{"failed":0,"total":0},"whatsmyuri":{"failed":0,"total":0},"writebacklisten":{"failed":0,"total":0}},"cursor":{"open":{"noTimeout":0,"pinned":0,"total":0},"timedOut":0},"document":{"deleted":0,"inserted":0,"returned":0,"updated":0},"getLastError":{"wtime":{"num":0,"totalMillis":0},"wtimeouts":0},"operation":{"fastmod":0,"idhack":0,"scanAndOrder":0,"writeConflicts":0},"queryExecutor":{"scanned":0,"scannedObjects":0},"record":{"moves":0},"repl":{"apply":{"batches":{"num":0,"totalMillis":0},"ops":0},"buffer":{"count":0,"maxSizeBytes":268435456,"sizeBytes":0},"network":{"bytes":0,"getmores":{"num":0,"totalMillis":0},"ops":0,"readersCreated":0},"preload":{"docs":{"num":0,"totalMillis":0},"indexes":{"num":0,"totalMillis":0}}},"storage":{"freelist":{"search":{"bucketExhausted":0,"requests":8,"scanned":0}}},"ttl":{"deletedDocuments":0,"passes":0}},"network":{"bytesIn":224,"bytesOut":381,"numRequests":4},"ok":1,"opcounters":{"command":5,"delete":0,"getmore":0,"insert":0,"query":1,"update":0},"opcountersRepl":{"command":0,"delete":0,"getmore":0,"insert":0,"query":0,"update":0},"pid":1,"process":"mongod","storageEngine":{"name":"mmapv1"},"uptime":5,"uptimeEstimate":5,"uptimeMillis":5616,"version":"3.0.5","writeBacksQueued":false}
`

	var v interface{}
	err := json.Unmarshal([]byte(stub3_0_5), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["duration_ms"]).String(), "float64")
	assert.EqualValues(t, stat["duration_ms"], 3.0)
}

func TestParse32(t *testing.T) {
	var mongodb MongoDBPlugin
	stub3_2_0 := `
{"host":"mbp13.local","advisoryHostFQDNs":["mbp13.local"],"version":"3.2.0","process":"mongod","pid":20843,"uptime":725,"uptimeMillis":724266,"uptimeEstimate":625,"localTime":"2016-01-06T08:02:12.187Z","asserts":{"regular":0,"warning":0,"msg":0,"user":0,"rollovers":0},"connections":{"current":1,"available":2047,"totalCreated":26},"extra_info":{"note":"fieldsvarybyplatform","page_faults":457},"globalLock":{"totalTime":724259000,"currentQueue":{"total":0,"readers":0,"writers":0},"activeClients":{"total":8,"readers":0,"writers":0}},"locks":{"Global":{"acquireCount":{"r":400,"w":2,"W":4}},"Database":{"acquireCount":{"r":197,"W":2}},"Collection":{"acquireCount":{"r":197}},"Metadata":{"acquireCount":{"w":1}}},"network":{"bytesIn":11962,"bytesOut":419980,"numRequests":174},"opcounters":{"insert":0,"query":1,"update":0,"delete":0,"getmore":0,"command":175},"opcountersRepl":{"insert":0,"query":0,"update":0,"delete":0,"getmore":0,"command":0},"storageEngine":{"name":"wiredTiger","supportsCommittedReads":true},"wiredTiger":{"uri":"statistics:","LSM":{"sleepforLSMcheckpointthrottle":0,"sleepforLSMmergethrottle":0,"rowsmergedinanLSMtree":0,"applicationworkunitscurrentlyqueued":0,"mergeworkunitscurrentlyqueued":0,"treequeuehitmaximum":0,"switchworkunitscurrentlyqueued":0,"treemaintenanceoperationsscheduled":0,"treemaintenanceoperationsdiscarded":0,"treemaintenanceoperationsexecuted":0},"async":{"numberofallocationstateraces":0,"numberofoperationslotsviewedforallocation":0,"currentworkqueuelength":0,"numberofflushcalls":0,"numberoftimesoperationallocationfailed":0,"maximumworkqueuelength":0,"numberoftimesworkerfoundnowork":0,"totalallocations":0,"totalcompactcalls":0,"totalinsertcalls":0,"totalremovecalls":0,"totalsearchcalls":0,"totalupdatecalls":0},"block-manager":{"mappedbytesread":0,"bytesread":4096,"byteswritten":86016,"mappedblocksread":0,"blockspre-loaded":0,"blocksread":1,"blockswritten":19},"cache":{"trackeddirtybytesinthecache":0,"trackedbytesbelongingtointernalpagesinthecache":1398,"bytescurrentlyinthecache":16143,"trackedbytesbelongingtoleafpagesinthecache":14745,"maximumbytesconfigured":9663676416,"trackedbytesbelongingtooverflowpagesinthecache":0,"bytesreadintocache":0,"byteswrittenfromcache":13061,"pagesevictedbyapplicationthreads":0,"checkpointblockedpageeviction":0,"unmodifiedpagesevicted":0,"pagesplitduringevictiondeepenedthetree":0,"modifiedpagesevicted":0,"pagesselectedforevictionunabletobeevicted":0,"pagesevictedbecausetheyexceededthein-memorymaximum":0,"pagesevictedbecausetheyhadchainsofdeleteditems":0,"failedevictionofpagesthatexceededthein-memorymaximum":0,"hazardpointerblockedpageeviction":0,"internalpagesevicted":0,"maximumpagesizeateviction":0,"evictionservercandidatequeueemptywhentoppingup":0,"evictionservercandidatequeuenotemptywhentoppingup":0,"evictionserverevictingpages":0,"evictionserverpopulatingqueue,butnotevictingpages":0,"evictionserverunabletoreachevictiongoal":0,"internalpagessplitduringeviction":0,"leafpagessplitduringeviction":0,"pageswalkedforeviction":0,"evictionworkerthreadevictingpages":0,"in-memorypagesplits":0,"in-memorypagepassedcriteriatobesplit":0,"lookasidetableinsertcalls":0,"lookasidetableremovecalls":0,"percentageoverhead":8,"trackeddirtypagesinthecache":0,"pagescurrentlyheldinthecache":11,"pagesreadintocache":0,"pagesreadintocacherequiringlookasideentries":0,"pageswrittenfromcache":12,"pagewrittenrequiringlookasiderecords":0,"pageswrittenrequiringin-memoryrestoration":0},"connection":{"pthreadmutexconditionwaitcalls":90553,"filescurrentlyopen":9,"memoryallocations":12590,"memoryfrees":12013,"memoryre-allocations":2360,"totalreadI/Os":8,"pthreadmutexsharedlockread-lockcalls":951,"pthreadmutexsharedlockwrite-lockcalls":706,"totalwriteI/Os":44},"cursor":{"cursorcreatecalls":131,"cursorinsertcalls":25,"cursornextcalls":13,"cursorprevcalls":3,"cursorremovecalls":0,"cursorresetcalls":56,"cursorrestartedsearches":0,"cursorsearchcalls":145,"cursorsearchnearcalls":1,"truncatecalls":0,"cursorupdatecalls":0},"data-handle":{"connectiondatahandlescurrentlyactive":6,"sessiondhandlesswept":0,"sessionsweepattempts":20,"connectionsweepdhandlesclosed":0,"connectionsweepcandidatebecamereferenced":0,"connectionsweepdhandlesremovedfromhashlist":1,"connectionsweeptime-of-deathsets":1,"connectionsweeps":72},"log":{"totallogbuffersize":33554432,"logbytesofpayloaddata":7566,"logbyteswritten":9856,"yieldswaitingforpreviouslogfileclose":0,"totalsizeofcompressedrecords":5927,"totalin-memorysizeofcompressedrecords":10292,"logrecordstoosmalltocompress":10,"logrecordsnotcompressed":7,"logrecordscompressed":8,"logflushoperations":7037,"maximumlogfilesize":104857600,"pre-allocatedlogfilesprepared":2,"numberofpre-allocatedlogfilestocreate":2,"pre-allocatedlogfilesnotreadyandmissed":1,"pre-allocatedlogfilesused":0,"logreleaseadvanceswriteLSN":8,"recordsprocessedbylogscan":0,"logscanrecordsrequiringtworeads":0,"logscanoperations":0,"consolidatedslotclosures":20847,"writtenslotscoalesced":0,"loggingbytesconsolidated":9472,"consolidatedslotjoins":25,"consolidatedslotjoinraces":0,"busyreturnsattemptingtoswitchslots":0,"consolidatedslotjointransitions":20847,"consolidatedslotunbufferedwrites":0,"logsyncoperations":12,"logsync_diroperations":1,"logserverthreadadvanceswriteLSN":5,"logwriteoperations":25,"logfilesmanuallyzero-filled":0},"reconciliation":{"pagesdeleted":0,"fast-pathpagesdeleted":0,"pagereconciliationcalls":12,"pagereconciliationcallsforeviction":0,"splitbytescurrentlyawaitingfree":0,"splitobjectscurrentlyawaitingfree":0},"session":{"opencursorcount":13,"opensessioncount":15},"thread-yield":{"pageacquirebusyblocked":0,"pageacquireevictionblocked":0,"pageacquirelockedblocked":0,"pageacquirereadblocked":0,"pageacquiretimesleeping(usecs)":0},"transaction":{"transactionbegins":16,"transactioncheckpoints":12,"transactioncheckpointgeneration":12,"transactioncheckpointcurrentlyrunning":0,"transactioncheckpointmaxtime(msecs)":77,"transactioncheckpointmintime(msecs)":0,"transactioncheckpointmostrecenttime(msecs)":0,"transactioncheckpointtotaltime(msecs)":138,"transactionscommitted":3,"transactionfailuresduetocacheoverflow":0,"transactionrangeofIDscurrentlypinnedbyacheckpoint":0,"transactionrangeofIDscurrentlypinned":0,"transactionrangeofIDscurrentlypinnedbynamedsnapshots":0,"transactionsrolledback":13,"numberofnamedsnapshotscreated":0,"numberofnamedsnapshotsdropped":0,"transactionsynccalls":0},"concurrentTransactions":{"write":{"out":0,"available":128,"totalTickets":128},"read":{"out":0,"available":128,"totalTickets":128}}},"writeBacksQueued":false,"mem":{"bits":64,"resident":52,"virtual":2511,"supported":true,"mapped":0,"mappedWithJournal":0},"metrics":{"commands":{"buildInfo":{"failed":0,"total":23},"getLog":{"failed":0,"total":2},"getnonce":{"failed":0,"total":3},"isMaster":{"failed":0,"total":71},"ping":{"failed":0,"total":6},"replSetGetStatus":{"failed":23,"total":23},"serverStatus":{"failed":0,"total":24},"whatsmyuri":{"failed":0,"total":23}},"cursor":{"timedOut":0,"open":{"noTimeout":0,"pinned":0,"total":0}},"document":{"deleted":0,"inserted":0,"returned":0,"updated":0},"getLastError":{"wtime":{"num":0,"totalMillis":0},"wtimeouts":0},"operation":{"fastmod":0,"idhack":0,"scanAndOrder":0,"writeConflicts":0},"queryExecutor":{"scanned":0,"scannedObjects":0},"record":{"moves":0},"repl":{"executor":{"counters":{"eventCreated":0,"eventWait":0,"cancels":0,"waits":0,"scheduledNetCmd":0,"scheduledDBWork":0,"scheduledXclWork":0,"scheduledWorkAt":0,"scheduledWork":0,"schedulingFailures":0},"queues":{"networkInProgress":0,"dbWorkInProgress":0,"exclusiveInProgress":0,"sleepers":0,"ready":0,"free":0},"unsignaledEvents":0,"eventWaiters":0,"shuttingDown":false,"networkInterface":"NetworkInterfaceASIOinShutdown:0"},"apply":{"batches":{"num":0,"totalMillis":0},"ops":0},"buffer":{"count":0,"maxSizeBytes":268435456,"sizeBytes":0},"network":{"bytes":0,"getmores":{"num":0,"totalMillis":0},"ops":0,"readersCreated":0},"preload":{"docs":{"num":0,"totalMillis":0},"indexes":{"num":0,"totalMillis":0}}},"storage":{"freelist":{"search":{"bucketExhausted":0,"requests":0,"scanned":0}}},"ttl":{"deletedDocuments":0,"passes":12}},"ok":1}
`

	var v interface{}
	err := json.Unmarshal([]byte(stub3_2_0), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["opcounters_command"]).String(), "float64")
	assert.EqualValues(t, stat["opcounters_command"], 175)
}

//Check in version 3.6.

func TestParse36(t *testing.T) {
	var mongodb MongoDBPlugin
	stub3_6_4 := `
{"asserts":{"msg":0,"regular":0,"rollovers":0,"user":218,"warning":0},"connections":{"available":818,"current":1,"totalCreated":78970},"extra_info":{"note":"fields vary by platform","page_faults":261675},"globalLock":{"activeClients":{"readers":0,"total":10,"writers":0},"currentQueue":{"readers":0,"total":0,"writers":0},"totalTime":4299432115000},"host":"foobar-12810","localTime":"2019-02-03T03:34:03.4Z","locks":{"Collection":{"acquireCount":{"r":1700856134,"w":2089228}},"Database":{"acquireCount":{"R":56,"W":17510,"r":1700899731,"w":2089185},"acquireWaitCount":{"W":62,"r":119},"timeAcquiringMicros":{"W":34661,"r":133960011}},"Global":{"acquireCount":{"W":12,"r":3408208489,"w":2106695},"acquireWaitCount":{"r":1,"w":1},"timeAcquiringMicros":{"r":322174,"w":480349}},"Metadata":{"acquireCount":{"W":8531},"acquireWaitCount":{"W":3},"timeAcquiringMicros":{"W":11506}}},"logicalSessionRecordCache":{"activeSessionsCount":0,"lastSessionsCollectionJobCursorsClosed":0,"lastSessionsCollectionJobDurationMillis":0,"lastSessionsCollectionJobEntriesEnded":0,"lastSessionsCollectionJobEntriesRefreshed":0,"lastSessionsCollectionJobTimestamp":"2019-02-03T03:32:02.768Z","lastTransactionReaperJobDurationMillis":0,"lastTransactionReaperJobEntriesCleanedUp":0,"lastTransactionReaperJobTimestamp":"2018-12-15T09:16:58.857Z","sessionsCollectionJobCount":14331,"transactionReaperJobCount":0},"mem":{"bits":64,"mapped":0,"mappedWithJournal":0,"resident":1326,"supported":true,"virtual":2873},"metrics":{"commands":{"\u003cUNKNOWN\u003e":0,"_configsvrAddShard":{"failed":0,"total":0},"_configsvrAddShardToZone":{"failed":0,"total":0},"_configsvrBalancerStart":{"failed":0,"total":0},"_configsvrBalancerStatus":{"failed":0,"total":0},"_configsvrBalancerStop":{"failed":0,"total":0},"_configsvrCommitChunkMerge":{"failed":0,"total":0},"_configsvrCommitChunkMigration":{"failed":0,"total":0},"_configsvrCommitChunkSplit":{"failed":0,"total":0},"_configsvrCreateDatabase":{"failed":0,"total":0},"_configsvrEnableSharding":{"failed":0,"total":0},"_configsvrMoveChunk":{"failed":0,"total":0},"_configsvrMovePrimary":{"failed":0,"total":0},"_configsvrRemoveShard":{"failed":0,"total":0},"_configsvrRemoveShardFromZone":{"failed":0,"total":0},"_configsvrShardCollection":{"failed":0,"total":0},"_configsvrUpdateZoneKeyRange":{"failed":0,"total":0},"_flushRoutingTableCacheUpdates":{"failed":0,"total":0},"_getNextSessionMods":{"failed":0,"total":0},"_getUserCacheGeneration":{"failed":0,"total":0},"_isSelf":{"failed":0,"total":0},"_mergeAuthzCollections":{"failed":0,"total":0},"_migrateClone":{"failed":0,"total":0},"_recvChunkAbort":{"failed":0,"total":0},"_recvChunkCommit":{"failed":0,"total":0},"_recvChunkStart":{"failed":0,"total":0},"_recvChunkStatus":{"failed":0,"total":0},"_transferMods":{"failed":0,"total":0},"aggregate":{"failed":0,"total":8202},"appendOplogNote":{"failed":0,"total":0},"applyOps":{"failed":0,"total":0},"authSchemaUpgrade":{"failed":0,"total":0},"authenticate":{"failed":0,"total":0},"availableQueryOptions":{"failed":0,"total":0},"buildInfo":{"failed":0,"total":313},"checkShardingIndex":{"failed":0,"total":0},"cleanupOrphaned":{"failed":0,"total":0},"clone":{"failed":0,"total":0},"cloneCollection":{"failed":0,"total":0},"cloneCollectionAsCapped":{"failed":0,"total":0},"collMod":{"failed":0,"total":0},"collStats":{"failed":0,"total":0},"compact":{"failed":0,"total":0},"connPoolStats":{"failed":0,"total":0},"connPoolSync":{"failed":0,"total":0},"connectionStatus":{"failed":0,"total":0},"convertToCapped":{"failed":0,"total":0},"copydb":{"failed":0,"total":0},"copydbgetnonce":{"failed":0,"total":0},"copydbsaslstart":{"failed":0,"total":0},"count":{"failed":0,"total":6762},"create":{"failed":0,"total":0},"createIndexes":{"failed":0,"total":17311},"createRole":{"failed":0,"total":0},"createUser":{"failed":0,"total":0},"currentOp":{"failed":0,"total":0},"dataSize":{"failed":0,"total":0},"dbHash":{"failed":0,"total":0},"dbStats":{"failed":0,"total":0},"delete":{"failed":0,"total":2770},"distinct":{"failed":0,"total":0},"driverOIDTest":{"failed":0,"total":0},"drop":{"failed":7,"total":36},"dropAllRolesFromDatabase":{"failed":0,"total":0},"dropAllUsersFromDatabase":{"failed":0,"total":0},"dropDatabase":{"failed":0,"total":7},"dropIndexes":{"failed":0,"total":0},"dropRole":{"failed":0,"total":0},"dropUser":{"failed":0,"total":0},"endSessions":{"failed":0,"total":35175},"eval":{"failed":0,"total":0},"explain":{"failed":0,"total":0},"features":{"failed":0,"total":0},"filemd5":{"failed":0,"total":0},"find":{"failed":0,"total":2327339},"findAndModify":{"failed":0,"total":0},"forceerror":{"failed":0,"total":0},"fsync":{"failed":0,"total":0},"fsyncUnlock":{"failed":0,"total":0},"geoNear":{"failed":0,"total":0},"geoSearch":{"failed":0,"total":0},"getCmdLineOpts":{"failed":0,"total":0},"getDiagnosticData":{"failed":0,"total":0},"getLastError":{"failed":0,"total":7},"getLog":{"failed":0,"total":1},"getMore":{"failed":0,"total":38615},"getParameter":{"failed":0,"total":0},"getPrevError":{"failed":0,"total":0},"getShardMap":{"failed":0,"total":0},"getShardVersion":{"failed":0,"total":0},"getnonce":{"failed":0,"total":1306},"grantPrivilegesToRole":{"failed":0,"total":0},"grantRolesToRole":{"failed":0,"total":0},"grantRolesToUser":{"failed":0,"total":1},"group":{"failed":0,"total":0},"handshake":{"failed":0,"total":0},"hostInfo":{"failed":0,"total":0},"insert":{"failed":0,"total":1244505},"invalidateUserCache":{"failed":0,"total":0},"isMaster":{"failed":0,"total":266847},"killAllSessions":{"failed":0,"total":0},"killAllSessionsByPattern":{"failed":0,"total":0},"killCursors":{"failed":0,"total":6},"killOp":{"failed":0,"total":0},"killSessions":{"failed":0,"total":0},"listCollections":{"failed":0,"total":34},"listCommands":{"failed":0,"total":0},"listDatabases":{"failed":0,"total":103},"listIndexes":{"failed":0,"total":10},"lockInfo":{"failed":0,"total":0},"logRotate":{"failed":0,"total":0},"logout":{"failed":0,"total":8},"mapReduce":{"failed":0,"total":0},"mapreduce":{"shardedfinish":{"failed":0,"total":0}},"mergeChunks":{"failed":0,"total":0},"moveChunk":{"failed":0,"total":0},"parallelCollectionScan":{"failed":0,"total":0},"ping":{"failed":0,"total":9677},"planCacheClear":{"failed":0,"total":0},"planCacheClearFilters":{"failed":0,"total":0},"planCacheListFilters":{"failed":0,"total":0},"planCacheListPlans":{"failed":0,"total":0},"planCacheListQueryShapes":{"failed":0,"total":0},"planCacheSetFilter":{"failed":0,"total":0},"profile":{"failed":0,"total":5},"reIndex":{"failed":0,"total":0},"refreshSessions":{"failed":0,"total":0},"refreshSessionsInternal":{"failed":0,"total":0},"renameCollection":{"failed":0,"total":0},"repairCursor":{"failed":0,"total":0},"repairDatabase":{"failed":0,"total":0},"replSetAbortPrimaryCatchUp":{"failed":0,"total":0},"replSetElect":{"failed":0,"total":0},"replSetFreeze":{"failed":0,"total":0},"replSetFresh":{"failed":0,"total":0},"replSetGetConfig":{"failed":0,"total":0},"replSetGetRBID":{"failed":0,"total":0},"replSetGetStatus":{"failed":1,"total":1},"replSetHeartbeat":{"failed":0,"total":0},"replSetInitiate":{"failed":0,"total":0},"replSetMaintenance":{"failed":0,"total":0},"replSetReconfig":{"failed":0,"total":0},"replSetRequestVotes":{"failed":0,"total":0},"replSetResizeOplog":{"failed":0,"total":0},"replSetStepDown":{"failed":0,"total":0},"replSetStepUp":{"failed":0,"total":0},"replSetSyncFrom":{"failed":0,"total":0},"replSetUpdatePosition":{"failed":0,"total":0},"resetError":{"failed":0,"total":0},"resync":{"failed":0,"total":0},"revokePrivilegesFromRole":{"failed":0,"total":0},"revokeRolesFromRole":{"failed":0,"total":0},"revokeRolesFromUser":{"failed":0,"total":0},"rolesInfo":{"failed":0,"total":0},"saslContinue":{"failed":2,"total":73701},"saslStart":{"failed":2,"total":36856},"serverStatus":{"failed":0,"total":1290},"setFeatureCompatibilityVersion":{"failed":0,"total":0},"setParameter":{"failed":0,"total":0},"setShardVersion":{"failed":0,"total":0},"shardConnPoolStats":{"failed":0,"total":0},"shardingState":{"failed":0,"total":0},"shutdown":{"failed":0,"total":0},"splitChunk":{"failed":0,"total":0},"splitVector":{"failed":0,"total":0},"startSession":{"failed":0,"total":0},"top":{"failed":0,"total":25},"touch":{"failed":0,"total":0},"unsetSharding":{"failed":0,"total":0},"update":{"failed":0,"total":558349},"updateRole":{"failed":0,"total":0},"updateUser":{"failed":0,"total":0},"usersInfo":{"failed":0,"total":0},"validate":{"failed":0,"total":0},"whatsmyuri":{"failed":0,"total":82}},"cursor":{"open":{"noTimeout":0,"pinned":0,"total":0},"timedOut":45},"document":{"deleted":12746,"inserted":4235118,"returned":19830101,"updated":556421},"getLastError":{"wtime":{"num":0,"totalMillis":0},"wtimeouts":0},"operation":{"scanAndOrder":41394,"writeConflicts":0},"queryExecutor":{"scanned":170287541,"scannedObjects":998071417},"record":{"moves":0},"repl":{"apply":{"attemptsToBecomeSecondary":0,"batches":{"num":0,"totalMillis":0},"ops":0},"buffer":{"count":0,"maxSizeBytes":0,"sizeBytes":0},"executor":{"networkInterface":"\nNetworkInterfaceASIO Operations' Diagnostic:\nOperation:    Count:   \nConnecting    0        \nIn Progress   0        \nSucceeded     0        \nCanceled      0        \nFailed        0        \nTimed Out     0        \n\n","pool":{"inProgressCount":0},"queues":{"networkInProgress":0,"sleepers":0},"shuttingDown":false,"unsignaledEvents":0},"initialSync":{"completed":0,"failedAttempts":0,"failures":0},"network":{"bytes":0,"getmores":{"num":0,"totalMillis":0},"ops":0,"readersCreated":0},"preload":{"docs":{"num":0,"totalMillis":0},"indexes":{"num":0,"totalMillis":0}}},"storage":{"freelist":{"search":{"bucketExhausted":0,"requests":0,"scanned":0}}},"ttl":{"deletedDocuments":243,"passes":71656}},"network":{"bytesIn":63621967856,"bytesOut":55969180955,"compression":{"snappy":{"compressor":{"bytesIn":0,"bytesOut":0},"decompressor":{"bytesIn":0,"bytesOut":0}}},"numRequests":4570241,"physicalBytesIn":63621967856,"physicalBytesOut":55969180955,"serviceExecutorTaskStats":{"executor":"passthrough","threadsRunning":1}},"ok":1,"opLatencies":{"commands":{"latency":319770842,"ops":428888},"reads":{"latency":254611347099,"ops":2342194},"writes":{"latency":1912656540,"ops":1799158}},"opcounters":{"command":457766,"delete":40030,"getmore":38615,"insert":4235117,"query":2327494,"update":561179},"opcountersRepl":{"command":0,"delete":0,"getmore":0,"insert":0,"query":0,"update":0},"pid":14120,"process":"mongod","storageEngine":{"name":"wiredTiger","persistent":true,"readOnly":false,"supportsCommittedReads":true},"tcmalloc":{"generic":{"current_allocated_bytes":1263414592,"heap_size":2036281344},"tcmalloc":{"aggressive_memory_decommit":0,"central_cache_free_bytes":75923056,"current_total_thread_cache_bytes":10614000,"formattedString":"------------------------------------------------\nMALLOC:     1263415168 ( 1204.9 MiB) Bytes in use by application\nMALLOC: +    463212544 (  441.8 MiB) Bytes in page heap freelist\nMALLOC: +     75923056 (   72.4 MiB) Bytes in central cache freelist\nMALLOC: +       421728 (    0.4 MiB) Bytes in transfer cache freelist\nMALLOC: +     10613424 (   10.1 MiB) Bytes in thread cache freelists\nMALLOC: +     10281216 (    9.8 MiB) Bytes in malloc metadata\nMALLOC:   ------------\nMALLOC: =   1823867136 ( 1739.4 MiB) Actual memory used (physical + swap)\nMALLOC: +    222695424 (  212.4 MiB) Bytes released to OS (aka unmapped)\nMALLOC:   ------------\nMALLOC: =   2046562560 ( 1951.8 MiB) Virtual address space used\nMALLOC:\nMALLOC:          83527              Spans in use\nMALLOC:             20              Thread heaps in use\nMALLOC:           4096              Tcmalloc page size\n------------------------------------------------\nCall ReleaseFreeMemory() to release freelist memory to the OS (via madvise()).\nBytes released to the OS take up virtual address space but no physical memory.\n","max_total_thread_cache_bytes":517996544,"pageheap_commit_count":390808,"pageheap_committed_bytes":1813585920,"pageheap_decommit_count":59639,"pageheap_free_bytes":463212544,"pageheap_reserve_count":1081,"pageheap_scavenge_count":34355,"pageheap_total_commit_bytes":19607371776,"pageheap_total_decommit_bytes":17793785856,"pageheap_total_reserve_bytes":2036281344,"pageheap_unmapped_bytes":222695424,"spinlock_total_delay_ns":36809653132,"thread_cache_free_bytes":10614000,"total_free_bytes":86958784,"transfer_cache_free_bytes":421728}},"transactions":{"retriedCommandsCount":0,"retriedStatementsCount":0,"transactionsCollectionWriteCount":0},"uptime":4299432,"uptimeEstimate":4299432,"uptimeMillis":4299432133,"version":"3.6.4"}
`
	var v interface{}
	err := json.Unmarshal([]byte(stub3_6_4), &v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	bsonStats, err := bson.Marshal(v)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var m bson.M
	err = bson.Unmarshal(bsonStats, &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	stat, err := mongodb.parseStatus(m)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Mongodb Stats
	assert.EqualValues(t, reflect.TypeOf(stat["opcounters_command"]).String(), "float64")
	assert.EqualValues(t, stat["opcounters_command"], 457766)
}

func TestMetricKeyPrefix(t *testing.T) {
	var m MongoDBPlugin
	prefix := m.MetricKeyPrefix()
	assert.Equal(t, "mongodb", prefix)

	m.KeyPrefix = "test"
	prefix = m.MetricKeyPrefix()
	assert.Equal(t, "test", prefix)
}

func TestLabelPrefix(t *testing.T) {
	var m MongoDBPlugin
	label := m.LabelPrefix()
	assert.Equal(t, "MongoDB", label)

	m.KeyPrefix = "test"
	label = m.LabelPrefix()
	assert.Equal(t, "Test", label)
}

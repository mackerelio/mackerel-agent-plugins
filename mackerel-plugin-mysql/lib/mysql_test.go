package mpmysql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition_DisableInnoDB(t *testing.T) {
	var mysql MySQLPlugin

	mysql.DisableInnoDB = true
	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 8 {
		t.Errorf("GetTempfilename: %d should be 7", len(graphdef))
	}
}

func TestGraphDefinition(t *testing.T) {
	var mysql MySQLPlugin

	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 29 {
		t.Errorf("GetTempfilename: %d should be 28", len(graphdef))
	}
}

func TestGraphDefinition_DisableInnoDB_EnableExtended(t *testing.T) {
	var mysql MySQLPlugin

	mysql.DisableInnoDB = true
	mysql.EnableExtended = true
	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 18 {
		t.Errorf("GetTempfilename: %d should be 18", len(graphdef))
	}
}

func TestGraphDefinition_EnableExtended(t *testing.T) {
	var mysql MySQLPlugin

	mysql.EnableExtended = true
	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 39 {
		t.Errorf("GetTempfilename: %d should be 39", len(graphdef))
	}
}

func TestParseProcStat56(t *testing.T) {
	stub := `=====================================
2015-03-09 20:11:22 7f6c0c845700 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 6 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 178 srv_active, 0 srv_shutdown, 1244368 srv_idle
srv_master_thread log flush and writes: 1244546
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 227
OS WAIT ARRAY INFO: signal count 220
Mutex spin waits 923, rounds 9442, OS waits 193
RW-shared spins 19, rounds 538, OS waits 16
RW-excl spins 5, rounds 476, OS waits 13
Spin rounds per wait: 10.23 mutex, 28.32 RW-shared, 95.20 RW-excl
------------
TRANSACTIONS
------------
Trx id counter 1093821584
Purge done for trx's n:o < 1093815563 undo n:o < 0 state: running but idle
History list length 649
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started
MySQL thread id 27954, OS thread handle 0x7f6c0c845700, query id 90345 localhost root init
SHOW /*!50000 ENGINE*/ INNODB STATUS
---TRANSACTION 1093821554, not started
MySQL thread id 27893, OS thread handle 0x7f6c0c886700, query id 90144 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093821583, not started
MySQL thread id 27888, OS thread handle 0x7f6c0c8c7700, query id 90175 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093811214, not started
MySQL thread id 27887, OS thread handle 0x7f6c0c98a700, query id 80071 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093820819, not started
MySQL thread id 27886, OS thread handle 0x7f6c0c949700, query id 89403 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093811160, not started
MySQL thread id 27885, OS thread handle 0x7f6c0c908700, query id 80015 127.0.0.1 cactiuser cleaning up
--------
FILE I/O
--------
I/O thread 0 state: waiting for completed aio requests (insert buffer thread)
I/O thread 1 state: waiting for completed aio requests (log thread)
I/O thread 2 state: waiting for completed aio requests (read thread)
I/O thread 3 state: waiting for completed aio requests (read thread)
I/O thread 4 state: waiting for completed aio requests (read thread)
I/O thread 5 state: waiting for completed aio requests (read thread)
I/O thread 6 state: waiting for completed aio requests (write thread)
I/O thread 7 state: waiting for completed aio requests (write thread)
I/O thread 8 state: waiting for completed aio requests (write thread)
I/O thread 9 state: waiting for completed aio requests (write thread)
Pending normal aio reads: 0 [0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0] ,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
124669 OS file reads, 4457 OS file writes, 3498 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 0.00 writes/s, 0.00 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 63, seg size 65, 2 merges
merged operations:
 insert 48, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 34679, node heap has 1 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
---
LOG
---
Log sequence number 53339891261
Log flushed up to   53339891261
Pages flushed up to 53339891261
Last checkpoint at  53339891261
10 pending log writes, 20 pending chkp writes
3395 log i/o's done, 0.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 17170432; in additional pool allocated 0
Dictionary memory allocated 318159
Buffer pool size   1024
Free buffers       755
Database pages     256
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 6, not young 751793
0.00 youngs/s, 0.00 non-youngs/s
Pages read 124617, created 40, written 1020
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
No buffer pool page gets since the last printout
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 256, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
0 read views open inside InnoDB
Main thread process no. 1968, id 140101998331648, state: sleeping
Number of rows inserted 3089, updated 220, deleted 212, read 2099881
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 0.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT`
	stat := make(map[string]float64)

	parseInnodbStatus(stub, false, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 947)
	assert.EqualValues(t, stat["spin_rounds"], 9442)
	assert.EqualValues(t, stat["os_waits"], 222)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 1093821584)
	assert.EqualValues(t, stat["unpurged_txns"], 6021)
	assert.EqualValues(t, stat["history_list"], 649)
	assert.EqualValues(t, stat["current_transactions"], 6)
	assert.EqualValues(t, stat["active_transactions"], 0)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 0) // empty
	assert.EqualValues(t, stat["read_views"], 0)            // empty
	assert.EqualValues(t, stat["innodb_tables_in_use"], 0)  // empty
	assert.EqualValues(t, stat["innodb_locked_tables"], 0)  // empty
	assert.EqualValues(t, stat["locked_transactions"], 0)   // empty
	assert.EqualValues(t, stat["innodb_lock_structs"], 0)   // empty
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 63)
	assert.EqualValues(t, stat["ibuf_cell_count"], 65)
	assert.EqualValues(t, stat["ibuf_inserts"], 48)
	assert.EqualValues(t, stat["ibuf_merges"], 2)
	assert.EqualValues(t, stat["ibuf_merged"], 48)
	assert.EqualValues(t, stat["hash_index_cells_total"], 34679)
	assert.EqualValues(t, stat["hash_index_cells_used"], 0) // empty
	// Log
	assert.EqualValues(t, stat["log_writes"], 3395)
	assert.EqualValues(t, stat["pending_log_writes"], 10)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 53339891261)
	assert.EqualValues(t, stat["log_bytes_flushed"], 53339891261)
	assert.EqualValues(t, stat["last_checkpoint"], 53339891261)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 0)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 17170432)
	assert.EqualValues(t, stat["additional_pool_alloc"], 0)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 0)     // empty
	assert.EqualValues(t, stat["page_hash_memory"], 0)         // empty
	assert.EqualValues(t, stat["dictionary_cache_memory"], 0)  // empty
	assert.EqualValues(t, stat["file_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["lock_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 0)

}

func TestParseProcStat55(t *testing.T) {

	stub := `=====================================
150310 10:40:23 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 19 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 19237002 1_second, 19236988 sleeps, 1923209 10_second, 6607 background, 6605 flush
srv_master_thread log flush and writes: 19327347
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 51338456, signal count 76067518
Mutex spin waits 4968902217, rounds 3687067031, OS waits 18668882
RW-shared spins 28966474, rounds 745089322, OS waits 23123092
RW-excl spins 22696709, rounds 329125903, OS waits 7388425
Spin rounds per wait: 0.74 mutex, 25.72 RW-shared, 14.50 RW-excl
------------------------
LATEST FOREIGN KEY ERROR
------------------------
140804 16:06:30 Transaction:
TRANSACTION 74D88599, ACTIVE 58 sec inserting, thread declared inside InnoDB 500
mysql tables in use 1, locked 1
14 lock struct(s), heap size 3112, 21 row lock(s), undo log entries 8
MySQL thread id 3244964, OS thread handle 0x7f7bcaecb700, query id 258109451 172.19.66.170 core update
------------
TRANSACTIONS
------------
Trx id counter C76C862D
Purge done for trx's n:o < C76C856A undo n:o < 0
History list length 3102
--------
FILE I/O
--------
I/O thread 0 state: waiting for completed aio requests (insert buffer thread)
I/O thread 1 state: waiting for completed aio requests (log thread)
I/O thread 2 state: waiting for completed aio requests (read thread)
I/O thread 3 state: waiting for completed aio requests (read thread)
I/O thread 4 state: waiting for completed aio requests (read thread)
I/O thread 5 state: waiting for completed aio requests (read thread)
I/O thread 6 state: waiting for completed aio requests (read thread)
I/O thread 7 state: waiting for completed aio requests (read thread)
I/O thread 8 state: waiting for completed aio requests (read thread)
I/O thread 9 state: waiting for completed aio requests (read thread)
I/O thread 10 state: waiting for completed aio requests (write thread)
I/O thread 11 state: waiting for completed aio requests (write thread)
I/O thread 12 state: waiting for completed aio requests (write thread)
I/O thread 13 state: waiting for completed aio requests (write thread)
I/O thread 14 state: waiting for completed aio requests (write thread)
I/O thread 15 state: waiting for completed aio requests (write thread)
I/O thread 16 state: waiting for completed aio requests (write thread)
I/O thread 17 state: waiting for completed aio requests (write thread)
Pending normal aio reads: 0 [0, 0, 0, 0, 0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0, 0, 0, 0, 0] ,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
80654072 OS file reads, 816873637 OS file writes, 575117750 OS fsyncs
3.58 reads/s, 16384 avg bytes/read, 20.74 writes/s, 9.53 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 9714, seg size 9716, 6224456 merges
merged operations:
 insert 8206050, delete mark 156570, delete 1983
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 42499631, node heap has 103815 buffer(s)
1329.14 hash searches/s, 338.14 non-hash searches/s
---
LOG
---
Log sequence number 1737766297992
Log flushed up to   1737766297992
Last checkpoint at  1737766159992
10 pending log writes, 20 pending chkp writes
532375066 log i/o's done, 7.79 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 21978152960; in additional pool allocated 0
Dictionary memory allocated 1592986
Buffer pool size   1310719
Free buffers       1
Database pages     1206903
Old database pages 445496
Modified db pages  180
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 222286179, not young 0
8.21 youngs/s, 0.00 non-youngs/s
Pages read 80651165, created 15602833, written 276352840
3.58 reads/s, 0.21 creates/s, 12.63 writes/s
Buffer pool hit rate 1000 / 1000, young-making rate 1 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 1206903, unzip_LRU len: 0
I/O sum[1126]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
12 read views open inside InnoDB
Main thread process no. 2510, id 140169706182400, state: sleeping
Number of rows inserted 686919123, updated 623703731, deleted 24439131, read 13570264742306
6.05 inserts/s, 1.84 updates/s, 0.00 deletes/s, 1960.21 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================`
	stat := make(map[string]float64)

	parseInnodbStatus(stub, true, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 5020565400)
	assert.EqualValues(t, stat["spin_rounds"], 3687067031)
	assert.EqualValues(t, stat["os_waits"], 49180399)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 3345778221)
	assert.EqualValues(t, stat["unpurged_txns"], 195)
	assert.EqualValues(t, stat["history_list"], 3102)
	assert.EqualValues(t, stat["current_transactions"], 0)
	assert.EqualValues(t, stat["active_transactions"], 0)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 0) // empty
	assert.EqualValues(t, stat["read_views"], 12)
	assert.EqualValues(t, stat["innodb_tables_in_use"], 0) // empty
	assert.EqualValues(t, stat["innodb_locked_tables"], 0) // empty
	assert.EqualValues(t, stat["locked_transactions"], 0)  // empty
	assert.EqualValues(t, stat["innodb_lock_structs"], 0)  // empty
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 9714)
	assert.EqualValues(t, stat["ibuf_cell_count"], 9716)
	assert.EqualValues(t, stat["ibuf_inserts"], 8206050)
	assert.EqualValues(t, stat["ibuf_merges"], 6224456)
	assert.EqualValues(t, stat["ibuf_merged"], 8364603)
	assert.EqualValues(t, stat["hash_index_cells_total"], 42499631)
	assert.EqualValues(t, stat["hash_index_cells_used"], 0)
	// Log
	assert.EqualValues(t, stat["log_writes"], 532375066)
	assert.EqualValues(t, stat["pending_log_writes"], 10)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 1737766297992)
	assert.EqualValues(t, stat["log_bytes_flushed"], 1737766297992)
	assert.EqualValues(t, stat["last_checkpoint"], 1737766159992)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 0)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 21978152960)
	assert.EqualValues(t, stat["additional_pool_alloc"], 0)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 0)     // empty
	assert.EqualValues(t, stat["page_hash_memory"], 0)         // empty
	assert.EqualValues(t, stat["dictionary_cache_memory"], 0)  // empty
	assert.EqualValues(t, stat["file_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["lock_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 138000)
}

func TestParseProcStat51(t *testing.T) {

	stub := `=====================================
150310 10:34:58 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 21 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 15513788 1_second, 15513624 sleeps, 1551102 10_second, 2807 background, 2807 flush
srv_master_thread log flush and writes: 15526310
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 2951389, signal count 41536793
Mutex spin waits 158882785, rounds 142931556, OS waits 1214105
RW-shared spins 9360396, OS waits 1636457; RW-excl spins 12223552, OS waits 76746
Spin rounds per wait: 0.90 mutex, 6.38 RW-shared, 5.10 RW-excl
------------
TRANSACTIONS
------------
Trx id counter 39009CDC7
Purge done for trx's n:o < 39009CD1B undo n:o < 0
History list length 9
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started, process no 3794, OS thread id 140154700429056
MySQL thread id 125413339, query id 8845328767 localhost root
SHOW /*!50000 ENGINE*/ INNODB STATUS
---TRANSACTION 39009CD65, not started, process no 3794, OS thread id 140154778973952
MySQL thread id 125412426, query id 8845326939 localhost test
---TRANSACTION 39009CD35, not started, process no 3794, OS thread id 140154804532992
MySQL thread id 125412424, query id 8845326190 localhost test
---TRANSACTION 39009CD60, not started, process no 3794, OS thread id 140154746492672
MySQL thread id 125412423, query id 8845326929 localhost test
---TRANSACTION 39009CD30, not started, process no 3794, OS thread id 140154749953792
MySQL thread id 125412420, query id 8845326179 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154784298752
MySQL thread id 125412417, query id 8845326923 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154708150016
MySQL thread id 125412415, query id 8845326548 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154680993536
MySQL thread id 125412413, query id 8845326928 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154684188416
MySQL thread id 125412412, query id 8845326893 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154674337536
MySQL thread id 125412411, query id 8845326479 localhost test
---TRANSACTION 0, not started, process no 3794, OS thread id 140154686318336
MySQL thread id 125412410, query id 8845326477 localhost test
--------
FILE I/O
--------
I/O thread 0 state: waiting for i/o request (insert buffer thread)
I/O thread 1 state: waiting for i/o request (log thread)
I/O thread 2 state: waiting for i/o request (read thread)
I/O thread 3 state: waiting for i/o request (read thread)
I/O thread 4 state: waiting for i/o request (read thread)
I/O thread 5 state: waiting for i/o request (read thread)
I/O thread 6 state: waiting for i/o request (write thread)
I/O thread 7 state: waiting for i/o request (write thread)
I/O thread 8 state: waiting for i/o request (write thread)
I/O thread 9 state: waiting for i/o request (write thread)
Pending normal aio reads: 0, aio writes: 0,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
613992 OS file reads, 134400134 OS file writes, 83130666 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 4.67 writes/s, 2.10 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 5, seg size 7,
18849 inserts, 18849 merged recs, 17834 merges
Hash table size 14874907, node heap has 6180 buffer(s)
171.90 hash searches/s, 328.17 non-hash searches/s
---
LOG
---
Log sequence number 7220257512009
Log flushed up to   7220257512009
Last checkpoint at  7220257512009
10 pending log writes, 20 pending chkp writes
78358216 log i/o's done, 1.81 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 7685013504; in additional pool allocated 0
Dictionary memory allocated 5255181
Buffer pool size   458751
Free buffers       1
Database pages     452570
Old database pages 167041
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 1360770, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 1203250, created 1230474, written 83593763
0.00 reads/s, 0.38 creates/s, 4.29 writes/s
Buffer pool hit rate 1000 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 452570, unzip_LRU len: 0
I/O sum[130]:cur[25], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
1 read views open inside InnoDB
Main thread process no. 3794, id 140154864322304, state: sleeping
Number of rows inserted 24090641, updated 8332796, deleted 18513402, read 139771797310
0.71 inserts/s, 0.00 updates/s, 0.10 deletes/s, 236.42 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================`
	stat := make(map[string]float64)

	parseInnodbStatus(stub, true, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 180466733)
	assert.EqualValues(t, stat["spin_rounds"], 142931556)
	assert.EqualValues(t, stat["os_waits"], 2927308)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 15301463495)
	assert.EqualValues(t, stat["unpurged_txns"], 172)
	assert.EqualValues(t, stat["history_list"], 9)
	assert.EqualValues(t, stat["current_transactions"], 11)
	assert.EqualValues(t, stat["active_transactions"], 0)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 0) // empty
	assert.EqualValues(t, stat["read_views"], 1)
	assert.EqualValues(t, stat["innodb_tables_in_use"], 0) // empty
	assert.EqualValues(t, stat["innodb_locked_tables"], 0) // empty
	assert.EqualValues(t, stat["locked_transactions"], 0)  // empty
	assert.EqualValues(t, stat["innodb_lock_structs"], 0)  // empty
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 5)
	assert.EqualValues(t, stat["ibuf_cell_count"], 7)
	assert.EqualValues(t, stat["ibuf_inserts"], 18849)
	assert.EqualValues(t, stat["ibuf_merges"], 17834)
	assert.EqualValues(t, stat["ibuf_merged"], 18849)
	assert.EqualValues(t, stat["hash_index_cells_total"], 14874907)
	assert.EqualValues(t, stat["hash_index_cells_used"], 0)
	// Log
	assert.EqualValues(t, stat["log_writes"], 78358216)
	assert.EqualValues(t, stat["pending_log_writes"], 10)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 7220257512009)
	assert.EqualValues(t, stat["log_bytes_flushed"], 7220257512009)
	assert.EqualValues(t, stat["last_checkpoint"], 7220257512009)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 0)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 7685013504)
	assert.EqualValues(t, stat["additional_pool_alloc"], 0)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 0)     // empty
	assert.EqualValues(t, stat["page_hash_memory"], 0)         // empty
	assert.EqualValues(t, stat["dictionary_cache_memory"], 0)  // empty
	assert.EqualValues(t, stat["file_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["lock_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 0)
}

func TestParseProcStat50(t *testing.T) {

	stub := `=====================================
150515 18:25:10 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 3 seconds
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 781, signal count 781
Mutex spin waits 0, rounds 30300, OS waits 1
RW-shared spins 1755, OS waits 778; RW-excl spins 7, OS waits 2
------------
TRANSACTIONS
------------
Trx id counter 0 2369392
Purge done for trx's n:o < 0 2368227 undo n:o < 0 0
History list length 1
Total number of lock structs in row lock hash table 0
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0 0, not started, process no 28986, OS thread id 3032255376
MySQL thread id 31989, query id 288360 localhost root
SHOW /*!50000 ENGINE*/ INNODB STATUS
--------
FILE I/O
--------
I/O thread 0 state: waiting for i/o request (insert buffer thread)
I/O thread 1 state: waiting for i/o request (log thread)
I/O thread 2 state: waiting for i/o request (read thread)
I/O thread 3 state: waiting for i/o request (write thread)
Pending normal aio reads: 0, aio writes: 0,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
332 OS file reads, 7564 OS file writes, 4398 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 0.00 writes/s, 0.00 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2,
2 inserts, 2 merged recs, 2 merges
Hash table size 34679, used cells 23275, node heap has 39 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
---
LOG
---
Log sequence number 0 51296721
Log flushed up to   0 51296721
Last checkpoint at  0 51296721
10 pending log writes, 20 pending chkp writes
2158 log i/o's done, 0.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 17874468; in additional pool allocated 1048576
Buffer pool size   512
Free buffers       1
Database pages     472
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages read 467, created 9, written 5185
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
No buffer pool page gets since the last printout
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
1 read views open inside InnoDB
Main thread process no. 28986, id 2996738960, state: waiting for server activity
Number of rows inserted 835, updated 104, deleted 2, read 226461457
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 0.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================`
	stat := make(map[string]float64)

	parseInnodbStatus(stub, true, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 1762)
	assert.EqualValues(t, stat["spin_rounds"], 30300)
	assert.EqualValues(t, stat["os_waits"], 781)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 2369392)
	assert.EqualValues(t, stat["unpurged_txns"], 1165)
	assert.EqualValues(t, stat["history_list"], 1)
	assert.EqualValues(t, stat["current_transactions"], 1)
	assert.EqualValues(t, stat["active_transactions"], 0)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 0) // empty
	assert.EqualValues(t, stat["read_views"], 1)
	assert.EqualValues(t, stat["innodb_tables_in_use"], 0) // empty
	assert.EqualValues(t, stat["innodb_locked_tables"], 0) // empty
	assert.EqualValues(t, stat["locked_transactions"], 0)  // empty
	assert.EqualValues(t, stat["innodb_lock_structs"], 0)  // empty
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 0)
	assert.EqualValues(t, stat["ibuf_cell_count"], 2)
	assert.EqualValues(t, stat["ibuf_inserts"], 2)
	assert.EqualValues(t, stat["ibuf_merges"], 2)
	assert.EqualValues(t, stat["ibuf_merged"], 2)
	assert.EqualValues(t, stat["hash_index_cells_total"], 34679)
	assert.EqualValues(t, stat["hash_index_cells_used"], 23275)
	// Log
	assert.EqualValues(t, stat["log_writes"], 2158)
	assert.EqualValues(t, stat["pending_log_writes"], 10)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 51296721)
	assert.EqualValues(t, stat["log_bytes_flushed"], 51296721)
	assert.EqualValues(t, stat["last_checkpoint"], 51296721)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 0)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 17874468)
	assert.EqualValues(t, stat["additional_pool_alloc"], 1048576)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 0)     // empty
	assert.EqualValues(t, stat["page_hash_memory"], 0)         // empty
	assert.EqualValues(t, stat["dictionary_cache_memory"], 0)  // empty
	assert.EqualValues(t, stat["file_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["lock_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 0)
}

func TestParseProcStat57(t *testing.T) {
	stub := `
=====================================
2016-02-22 19:08:31 0x700000eda000 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 4 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 1 srv_active, 0 srv_shutdown, 2 srv_idle
srv_master_thread log flush and writes: 3
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 63
OS WAIT ARRAY INFO: signal count 111
RW-shared spins 0, rounds 85, OS waits 22
RW-excl spins 0, rounds 4705, OS waits 17
RW-sx spins 70, rounds 70, OS waits 70
Spin rounds per wait: 85.00 RW-shared, 4705.00 RW-excl, 0.00 RW-sx
------------
TRANSACTIONS
------------
Trx id counter 49154
Purge done for trx's n:o < 44675 undo n:o < 0 state: running but idle
History list length 775
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 281479529875248, not started
0 lock struct(s), heap size 1136, 0 row lock(s)
--------
FILE I/O
--------
I/O thread 0 state: waiting for i/o request (insert buffer thread)
I/O thread 1 state: waiting for i/o request (log thread)
I/O thread 2 state: waiting for i/o request (read thread)
I/O thread 3 state: waiting for i/o request (read thread)
I/O thread 4 state: waiting for i/o request (read thread)
I/O thread 5 state: waiting for i/o request (read thread)
I/O thread 6 state: waiting for i/o request (write thread)
I/O thread 7 state: waiting for i/o request (write thread)
I/O thread 8 state: waiting for i/o request (write thread)
I/O thread 9 state: waiting for i/o request (write thread)
Pending normal aio reads: [0, 0, 0, 0] , aio writes: [0, 0, 0, 0] ,
 ibuf aio reads:, log i/o's:, sync i/o's:
Pending flushes (fsync) log: 0; buffer pool: 0
516 OS file reads, 55 OS file writes, 9 OS fsyncs
128.97 reads/s, 20393 avg bytes/read, 13.75 writes/s, 2.25 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2, 0 merges
merged operations:
 insert 0, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 276671, node heap has 2 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 1 buffer(s)
Hash table size 276671, node heap has 1 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 4 buffer(s)
276.93 hash searches/s, 835.29 non-hash searches/s
---
LOG
---
Log sequence number 379575319
Log flushed up to   379575319
Pages flushed up to 379575319
Last checkpoint at  379575310
10 pending log flushes, 20 pending chkp writes
12 log i/o's done, 3.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total large memory allocated 1099431936
Dictionary memory allocated 312184
Buffer pool size   65528
Free buffers       64999
Database pages     521
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 487, created 34, written 36
121.72 reads/s, 8.50 creates/s, 9.00 writes/s
Buffer pool hit rate 974 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 521, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
----------------------
INDIVIDUAL BUFFER POOL INFO
----------------------
---BUFFER POOL 0
Buffer pool size   16382
Free buffers       16228
Database pages     152
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 152, created 0, written 2
37.99 reads/s, 0.00 creates/s, 0.50 writes/s
Buffer pool hit rate 976 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 152, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 1
Buffer pool size   16382
Free buffers       16244
Database pages     136
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 136, created 0, written 0
33.99 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 978 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 136, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 2
Buffer pool size   16382
Free buffers       16313
Database pages     67
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 67, created 0, written 0
16.75 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 975 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 67, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 3
Buffer pool size   16382
Free buffers       16214
Database pages     166
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 132, created 34, written 34
32.99 reads/s, 8.50 creates/s, 8.50 writes/s
Buffer pool hit rate 963 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 166, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
0 read views open inside InnoDB
Process ID=28837, Main thread ID=123145312497664, state: sleeping
Number of rows inserted 0, updated 0, deleted 0, read 8
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 2.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================
`
	stat := make(map[string]float64)
	parseInnodbStatus(stub, false, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 70)
	assert.EqualValues(t, stat["spin_rounds"], 0) // empty
	assert.EqualValues(t, stat["os_waits"], 109)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 49154) // empty
	assert.EqualValues(t, stat["unpurged_txns"], 4479)
	assert.EqualValues(t, stat["history_list"], 775)
	assert.EqualValues(t, stat["current_transactions"], 1)
	assert.EqualValues(t, stat["active_transactions"], 0)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 0) // empty
	assert.EqualValues(t, stat["read_views"], 0)
	assert.EqualValues(t, stat["innodb_tables_in_use"], 0) // empty
	assert.EqualValues(t, stat["innodb_locked_tables"], 0) // empty
	assert.EqualValues(t, stat["locked_transactions"], 0)  // empty
	assert.EqualValues(t, stat["innodb_lock_structs"], 0)  // empty
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 0)
	assert.EqualValues(t, stat["ibuf_cell_count"], 2)
	assert.EqualValues(t, stat["ibuf_inserts"], 0)
	assert.EqualValues(t, stat["ibuf_merges"], 0)
	assert.EqualValues(t, stat["ibuf_merged"], 0)
	assert.EqualValues(t, stat["hash_index_cells_total"], 276671)
	assert.EqualValues(t, stat["hash_index_cells_used"], 0)
	// Log
	assert.EqualValues(t, stat["log_writes"], 12)
	assert.EqualValues(t, stat["pending_log_writes"], 0)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 379575319)
	assert.EqualValues(t, stat["log_bytes_flushed"], 379575319)
	assert.EqualValues(t, stat["last_checkpoint"], 379575310)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 10)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 1099431936)
	assert.EqualValues(t, stat["additional_pool_alloc"], 0)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 0)     // empty
	assert.EqualValues(t, stat["page_hash_memory"], 0)         // empty
	assert.EqualValues(t, stat["dictionary_cache_memory"], 0)  // empty
	assert.EqualValues(t, stat["file_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["lock_system_memory"], 0)       // empty
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 9)

}

func TestParseLockedTransactions(t *testing.T) {

	stub := `=====================================
170829 11:50:33 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 18 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 26 1_second, 26 sleeps, 2 10_second, 9 background, 9 flush
srv_master_thread log flush and writes: 28
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 12, signal count 11
Mutex spin waits 6, rounds 180, OS waits 6
RW-shared spins 6, rounds 180, OS waits 6
RW-excl spins 0, rounds 0, OS waits 0
Spin rounds per wait: 30.00 mutex, 30.00 RW-shared, 0.00 RW-excl
--------
FILE I/O
--------
I/O thread 0 state: waiting for completed aio requests (insert buffer thread)
I/O thread 1 state: waiting for completed aio requests (log thread)
I/O thread 2 state: waiting for completed aio requests (read thread)
I/O thread 3 state: waiting for completed aio requests (read thread)
I/O thread 4 state: waiting for completed aio requests (read thread)
I/O thread 5 state: waiting for completed aio requests (read thread)
I/O thread 6 state: waiting for completed aio requests (write thread)
I/O thread 7 state: waiting for completed aio requests (write thread)
I/O thread 8 state: waiting for completed aio requests (write thread)
I/O thread 9 state: waiting for completed aio requests (write thread)
Pending normal aio reads: 0 [0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0] ,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
310 OS file reads, 174 OS file writes, 22 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 0.00 writes/s, 0.00 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2, 0 merges
merged operations:
 insert 0, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 276671, node heap has 1 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
---
LOG
---
Log sequence number 1602283
Log flushed up to   1602283
Last checkpoint at  1602283
Max checkpoint age    7782360
Checkpoint age target 7539162
Modified age          0
Checkpoint age        0
10 pending log writes, 20 pending chkp writes
40 log i/o's done, 0.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 137756672; in additional pool allocated 0
Total memory allocated by read views 88
Internal hash tables (constant factor + variable factor)
    Adaptive hash index 2233968 	(2213368 + 20600)
    Page hash           139112 (buffer pool 0 only)
    Dictionary cache    597886 	(554768 + 43118)
    File system         83536 	(82672 + 864)
    Lock system         334000 	(332872 + 1128)
    Recovery system     0 	(0 + 0)
Dictionary memory allocated 43118
Buffer pool size        8191
Buffer pool size, bytes 134201344
Free buffers            8039
Database pages          151
Old database pages      0
Modified db pages       0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 147, created 4, written 156
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
No buffer pool page gets since the last printout
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 151, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
1 read views open inside InnoDB
2 transactions active inside InnoDB
2 out of 1000 descriptors used
---OLDEST VIEW---
Normal read view
Read view low limit trx n:o 505
Read view up limit trx id 505
Read view low limit trx id 505
Read view individually stored trx ids:
-----------------
Main thread process no. 458, id 139631366485760, state: waiting for server activity
Number of rows inserted 2, updated 0, deleted 1, read 2
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 0.00 reads/s
------------
TRANSACTIONS
------------
Trx id counter 507
Purge done for trx's n:o < 505 undo n:o < 0
History list length 1
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started
MySQL thread id 8, OS thread handle 0x7efe7cb12700, query id 52 localhost root
SHOW ENGINE INNODB STATUS
---TRANSACTION 506, ACTIVE 804 sec starting index read
mysql tables in use 1, locked 1
LOCK WAIT 2 lock struct(s), heap size 376, 1 row lock(s)
MySQL thread id 3, OS thread handle 0x7efe7cb5b700, query id 47 localhost root statistics
SELECT * FROM test WHERE id = 1 LOCK IN SHARE MODE
------- TRX HAS BEEN WAITING 22 SEC FOR THIS LOCK TO BE GRANTED:
RECORD LOCKS space id 0 page no 307 n bits 72 index ` + "`PRIMARY` of table `test`.`test`" + ` trx id 506 lock mode S locks rec but not gap waiting
------------------
---TRANSACTION 505, ACTIVE 815 sec
2 lock struct(s), heap size 376, 1 row lock(s), undo log entries 1
MySQL thread id 2, OS thread handle 0x7efe7cba4700, query id 35 localhost root
----------------------------
END OF INNODB MONITOR OUTPUT
============================`
	stat := make(map[string]float64)
	parseInnodbStatus(stub, false, stat)
	// Innodb Semaphores
	assert.EqualValues(t, stat["spin_waits"], 12)
	assert.EqualValues(t, stat["spin_rounds"], 180)
	assert.EqualValues(t, stat["os_waits"], 12)
	assert.EqualValues(t, stat["innodb_sem_wait"], 0)         // empty
	assert.EqualValues(t, stat["innodb_sem_wait_time_ms"], 0) // empty
	// Innodb Transactions
	assert.EqualValues(t, stat["innodb_transactions"], 507)
	assert.EqualValues(t, stat["unpurged_txns"], 2)
	assert.EqualValues(t, stat["history_list"], 1)
	assert.EqualValues(t, stat["current_transactions"], 3)
	assert.EqualValues(t, stat["active_transactions"], 2)
	assert.EqualValues(t, stat["innodb_lock_wait_secs"], 22)
	assert.EqualValues(t, stat["read_views"], 1)
	assert.EqualValues(t, stat["innodb_tables_in_use"], 1)
	assert.EqualValues(t, stat["innodb_locked_tables"], 1)
	assert.EqualValues(t, stat["locked_transactions"], 1)
	assert.EqualValues(t, stat["innodb_lock_structs"], 4)
	// File I/O
	assert.EqualValues(t, stat["pending_normal_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_normal_aio_writes"], 0)
	assert.EqualValues(t, stat["pending_ibuf_aio_reads"], 0)
	assert.EqualValues(t, stat["pending_aio_log_ios"], 0)
	assert.EqualValues(t, stat["pending_aio_sync_ios"], 0)
	assert.EqualValues(t, stat["pending_log_flushes"], 0)
	assert.EqualValues(t, stat["pending_buf_pool_flushes"], 0)
	// Insert Buffer and Adaptive Hash Index
	assert.EqualValues(t, stat["ibuf_used_cells"], 1)
	assert.EqualValues(t, stat["ibuf_free_cells"], 0)
	assert.EqualValues(t, stat["ibuf_cell_count"], 2)
	assert.EqualValues(t, stat["ibuf_inserts"], 0)
	assert.EqualValues(t, stat["ibuf_merges"], 0)
	assert.EqualValues(t, stat["ibuf_merged"], 0)
	assert.EqualValues(t, stat["hash_index_cells_total"], 276671)
	assert.EqualValues(t, stat["hash_index_cells_used"], 0)
	// Log
	assert.EqualValues(t, stat["log_writes"], 40)
	assert.EqualValues(t, stat["pending_log_writes"], 10)
	assert.EqualValues(t, stat["pending_chkp_writes"], 20)
	assert.EqualValues(t, stat["log_bytes_written"], 1602283)
	assert.EqualValues(t, stat["log_bytes_flushed"], 1602283)
	assert.EqualValues(t, stat["last_checkpoint"], 1602283)
	assert.EqualValues(t, stat["log_pending_log_flushes"], 0)
	// Buffer Pool and Memory
	assert.EqualValues(t, stat["total_mem_alloc"], 137756672)
	assert.EqualValues(t, stat["additional_pool_alloc"], 0)
	assert.EqualValues(t, stat["adaptive_hash_memory"], 2233968)
	assert.EqualValues(t, stat["page_hash_memory"], 139112)
	assert.EqualValues(t, stat["dictionary_cache_memory"], 597886)
	assert.EqualValues(t, stat["file_system_memory"], 83536)
	assert.EqualValues(t, stat["lock_system_memory"], 334000)
	assert.EqualValues(t, stat["recovery_system_memory"], 0)   // empty
	assert.EqualValues(t, stat["thread_hash_memory"], 0)       // empty
	assert.EqualValues(t, stat["innodb_io_pattern_memory"], 0) // empty
	// etc
	assert.EqualValues(t, stat["unflushed_log"], 0)
	assert.EqualValues(t, stat["uncheckpointed_bytes"], 0)

}

func TestParseProcesslist1(t *testing.T) {
	stat := make(map[string]float64)
	pattern := []string{"NULL"}

	for _, val := range pattern {
		parseProcesslist(val, stat)
	}
	assert.EqualValues(t, 0, stat["State_closing_tables"])
	assert.EqualValues(t, 0, stat["State_copying_to_tmp_table"])
	assert.EqualValues(t, 0, stat["State_end"])
	assert.EqualValues(t, 0, stat["State_freeing_items"])
	assert.EqualValues(t, 0, stat["State_init"])
	assert.EqualValues(t, 0, stat["State_locked"])
	assert.EqualValues(t, 0, stat["State_login"])
	assert.EqualValues(t, 0, stat["State_preparing"])
	assert.EqualValues(t, 0, stat["State_reading_from_net"])
	assert.EqualValues(t, 0, stat["State_sending_data"])
	assert.EqualValues(t, 0, stat["State_sorting_result"])
	assert.EqualValues(t, 0, stat["State_statistics"])
	assert.EqualValues(t, 0, stat["State_updating"])
	assert.EqualValues(t, 0, stat["State_writing_to_net"])
	assert.EqualValues(t, 0, stat["State_none"])
	assert.EqualValues(t, 1, stat["State_other"])
}

func TestParseProcesslist2(t *testing.T) {
	stat := make(map[string]float64)

	// https://dev.mysql.com/doc/refman/5.6/en/general-thread-states.html
	pattern := []string{
		"",
		"After create",
		"altering table",
		"Analyzing",
		"checking permissions",
		"Checking table",
		"cleaning up",
		"closing tables",
		"committing alter table to storage engine",
		"converting HEAP to MyISAM",
		"MEMORY",
		"MyISAM",
		"copy to tmp table",
		"Copying to group table",
		"GROUP BY",
		"Copying to tmp table",
		"Copying to tmp table on disk",
		"Creating index",
		"Creating sort index",
		"creating table",
		"Creating tmp table",
		"deleting from main table",
		"deleting from reference tables",
		"discard_or_import_tablespace",
		"end",
		"executing",
		"Execution of init_command",
		"freeing items",
		"FULLTEXT initialization",
		"init",
		"Killed",
		"logging slow query",
		"login",
		"manage keys",
		"NULL",
		"Opening tables",
		"Opening table",
		"optimizing",
		"preparing",
		"preparing for alter table",
		"Purging old relay logs",
		"query end",
		"Reading from net",
		"Removing duplicates",
		"removing tmp table",
		"rename",
		"rename result table",
		"Reopen tables",
		"Repair by sorting",
		"Repair done",
		"Repair with keycache",
		"Rolling back",
		"Saving state",
		"Searching rows for update",
		"Sending data",
		"setup",
		"Sorting for group",
		"Sorting for order",
		"Sorting index",
		"Sorting result",
		"statistics",
		"System lock",
		"update",
		"Updating",
		"updating main table",
		"updating reference tables",
		"User lock",
		"User sleep",
		"Waiting for commit lock",
		"Waiting for global read lock",
		"Waiting for tables",
		"Waiting for table flush",
		"Waiting for lock_type lock",
		"Waiting for table level lock",
		"Waiting for event metadata lock",
		"Waiting for global read lock",
		"Waiting for schema metadat lock",
		"Waiting for stored function metadata  lock",
		"Waiting for stored procedure metadata lock",
		"Waiting for table metadata lock",
		"Waiting for trigger metadata lock",
		"Waiting on cond",
		"Writing to net",
		"Table lock",
	}

	for _, val := range pattern {
		parseProcesslist(val, stat)
	}
	assert.EqualValues(t, 1, stat["State_closing_tables"])
	assert.EqualValues(t, 1, stat["State_copying_to_tmp_table"])
	assert.EqualValues(t, 1, stat["State_end"])
	assert.EqualValues(t, 1, stat["State_freeing_items"])
	assert.EqualValues(t, 1, stat["State_init"])
	assert.EqualValues(t, 12, stat["State_locked"])
	assert.EqualValues(t, 1, stat["State_login"])
	assert.EqualValues(t, 1, stat["State_preparing"])
	assert.EqualValues(t, 1, stat["State_reading_from_net"])
	assert.EqualValues(t, 1, stat["State_sending_data"])
	assert.EqualValues(t, 1, stat["State_sorting_result"])
	assert.EqualValues(t, 1, stat["State_statistics"])
	assert.EqualValues(t, 1, stat["State_updating"])
	assert.EqualValues(t, 1, stat["State_writing_to_net"])
	assert.EqualValues(t, 1, stat["State_none"])
	assert.EqualValues(t, 58, stat["State_other"])
}

type TestCaseAio struct {
	stub   string
	reads  int
	writes int
}

func TestParseAio(t *testing.T) {
	pattern := []TestCaseAio{
		{"Pending normal aio reads: [1, 3, 5, 7] , aio writes: [3, 5, 7, 9] ,", 16, 24},
		{"Pending normal aio reads: [1, 3, 5, 7] ", 16, 0},
		{"Pending normal aio reads: 10 [4, 6] , aio writes: 20 [2, 4, 6, 8] ,", 10, 20},
		{"Pending normal aio reads: 10 [4, 6] ", 10, 0},
		{"Pending normal aio reads: 10, aio writes: 20,", 10, 20},
		{"Pending normal aio reads: 10", 10, 0},
		{"Pending normal aio reads:, aio writes: [1, 3, 5, 7],", 0, 16},
		{"Pending normal aio reads:, aio writes:,", 0, 0},
	}

	for _, tt := range pattern {
		stat := make(map[string]float64)
		parseInnodbStatus(tt.stub, false, stat)
		assert.EqualValues(t, stat["pending_normal_aio_reads"], tt.reads)
		assert.EqualValues(t, stat["pending_normal_aio_writes"], tt.writes)
	}
}

func TestMetricNamesShouldUniqueAndConst(t *testing.T) {
	m := MySQLPlugin{
		DisableInnoDB:  false,
		EnableExtended: true,
	}
	defs := m.GraphDefinition()
	keys := make(map[string]string) // metricName: graphDefName
	for name, g := range defs {
		for _, v := range g.Metrics {
			if v.Name == "Threads_connected" {
				if name != "connections" && name != "threads" {
					t.Errorf(`%q are duplicated in "connections", "threads" and %q`, v.Name, name)
				}
				continue
			}
			if v.Name == "Qcache_hits" {
				if name != "cmd" && name != "query_cache" {
					t.Errorf(`%q are duplicated in "cmd", "query_cache" and %q`, v.Name, name)
				}
				continue
			}

			if strings.ContainsAny(v.Name, "#*") {
				t.Errorf("%q should not contains wildcards", v.Name)
			}
			if s, ok := keys[v.Name]; ok {
				t.Errorf("%q are defined in both %q and %q", v.Name, s, name)
			}
			keys[v.Name] = name
		}
	}
}

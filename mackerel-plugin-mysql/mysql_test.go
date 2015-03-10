package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var mysql MySQLPlugin

	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 10 {
		t.Errorf("GetTempfilename: %d should be 10", len(graphdef))
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
0 pending log writes, 0 pending chkp writes
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

	err := parseInnodbStatus(stub, &stat)
	// fmt.Println(stat)
	assert.Nil(t, err)
	// Innodb Semaphores
	assert.Equal(t, stat["spin_waits"], 947)
	assert.Equal(t, stat["spin_rounds"], 9442)
	assert.Equal(t, stat["os_waits"], 222)
	// Innodb Transactions
	assert.Equal(t, stat["innodb_transactions"], 71194252676)
	assert.Equal(t, stat["unpurged_txns"], 49185)
	assert.Equal(t, stat["history_list"], 649)
	assert.Equal(t, stat["current_transactions"], 6)
	assert.Equal(t, stat["active_transactions"], 0)
	assert.Equal(t, stat["innodb_lock_wait_secs"], 0)
	assert.Equal(t, stat["read_views"], 0)
	assert.Equal(t, stat["innodb_tables_in_use"], 0)
	assert.Equal(t, stat["innodb_locked_tables"], 0)
	assert.Equal(t, stat["innodb_lock_structs"], 0)
	assert.Equal(t, stat["locked_transactions"], 0)
	assert.Equal(t, stat["innodb_lock_structs"], 0)
	// File I/O
	assert.Equal(t, stat["file_reads"], 124669)
	assert.Equal(t, stat["file_writes"], 4457)
	assert.Equal(t, stat["file_fsyncs"], 3498)
	assert.Equal(t, stat["pending_normal_aio_reads"], 0)
	assert.Equal(t, stat["pending_normal_aio_writes"], 0)
	assert.Equal(t, stat["pending_ibuf_aio_reads"], 0)
	assert.Equal(t, stat["pending_aio_log_ios"], 0)
	assert.Equal(t, stat["pending_aio_sync_ios"], 0)
	assert.Equal(t, stat["pending_log_flushes"], 0)
	assert.Equal(t, stat["pending_buf_pool_flushes"], 0)
	//assert.Equal(t, stat[""], )

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
0 pending log writes, 0 pending chkp writes
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

	err := parseInnodbStatus(stub, &stat)
	// fmt.Println(stat)
	assert.Nil(t, err)
	// Innodb Semaphores
	assert.Equal(t, stat["spin_waits"], 5020565400)
	assert.Equal(t, stat["spin_rounds"], 3687067031)
	assert.Equal(t, stat["os_waits"], 49180399)
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
0 pending log writes, 0 pending chkp writes
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

	err := parseInnodbStatus(stub, &stat)
	// fmt.Println(stat)
	assert.Nil(t, err)
	// Innodb Semaphores
	assert.Equal(t, stat["spin_waits"], 180466733)
	assert.Equal(t, stat["spin_rounds"], 142931556)
	assert.Equal(t, stat["os_waits"], 2927308)
}

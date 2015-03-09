package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestGraphDefinition(t *testing.T) {
	var mysql MySQLPlugin

	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 10 {
		t.Errorf("GetTempfilename: %d should be 10", len(graphdef))
	}
}

func TestParseProcStat(t *testing.T) {
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
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.Equal(t, stat["spin_waits"], 947)
	assert.Equal(t, stat["spin_rounds"], 9442)
	assert.Equal(t, stat["os_waits"], 222)
}
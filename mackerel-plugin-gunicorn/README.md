mackerel-plugin-gunicorn
=====================

Gunicorn custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-gunicorn [-status-file=<path>]
```

## Requirements

This plugin requires that gunicorn include the following code in its `config.py`.

``` python
import fcntl
import json
import os
import stat

SERVER_STATUS = "/dev/shm/gunicorn_status.json"


def on_starting(servers):
    if os.path.exists(SERVER_STATUS):
        os.remove(SERVER_STATUS)
    with open(SERVER_STATUS, mode="w") as f:
        obj = {
            "TotalAccesses": "0",
            "IdleWorkers": "0",
            "BusyWorkers": "0",
            "stats": [],
        }
        json.dump(obj, f)
    os.chown(SERVER_STATUS, os.getuid(), servers.cfg.gid)
    statinfo = os.stat(SERVER_STATUS)
    mode = statinfo.st_mode + stat.S_IWGRP
    os.chmod(SERVER_STATUS, mode=mode)


def post_fork(server, worker):
    stat = {
        "pid": worker.pid,
        "host": "",
        "method": "",
        "uri": "",
        "status": "_"
    }

    with open(SERVER_STATUS, mode="r+") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            obj = json.load(f)

            stats = [(i, v) for i, v in enumerate(obj["stats"]) if v["pid"] == worker.pid]
            if len(stats) == 0:
                obj["stats"].append(stat)
            else:
                for i, _ in stats:
                    obj["stats"][i] = stat

            obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) + 1)

            f.seek(0)
            f.truncate(0)
            json.dump(obj, f)
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)


def pre_request(worker, req):
    headers = dict(req.headers)
    stat = {
        "pid": worker.pid,
        "host": headers["HOST"],
        "method": req.method,
        "uri": req.uri,
        "status": "A"
    }

    with open(SERVER_STATUS, mode="r+") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            obj = json.load(f)

            stats = [(i, v) for i, v in enumerate(obj["stats"]) if v["pid"] == worker.pid]
            for i, _ in stats:
                obj["stats"][i] = stat

            obj["TotalAccesses"] = str(int(obj["TotalAccesses"]) + 1)
            obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) - 1)
            obj["BusyWorkers"] = str(int(obj["BusyWorkers"]) + 1)

            f.seek(0)
            f.truncate(0)
            json.dump(obj, f)
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)


def post_request(worker, req):
    headers = dict(req.headers)
    stat = {
        "pid": worker.pid,
        "host": headers["HOST"],
        "method": req.method,
        "uri": req.uri,
        "status": "_"
    }

    with open(SERVER_STATUS, mode="r+") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            obj = json.load(f)

            stats = [(i, v) for i, v in enumerate(obj["stats"]) if v["pid"] == worker.pid]
            if len(stats) == 0:
                worker.log.warn("not find worker.pid: %d in stats object", worker.pid)
            else:
                for i, _ in stats:
                    obj["stats"][i] = stat

            obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) + 1)
            obj["BusyWorkers"] = str(int(obj["BusyWorkers"]) - 1)

            f.seek(0)
            f.truncate(0)
            json.dump(obj, f)
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)


def worker_int(worker):
    stat = {
        "pid": worker.pid,
        "host": "",
        "method": "",
        "uri": "",
        "status": "SIGINT",
    }

    with open(SERVER_STATUS, mode="r+") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            obj = json.load(f)

            stats = [(i, v) for i, v in enumerate(obj["stats"]) if v["pid"] == worker.pid]
            if len(stats) == 0:
                obj["stats"].append(stat)
                obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) - 1)
            else:
                for i, v in stats:
                    obj["stats"][i] = stat
                    if v["status"] == "A":
                        obj["BusyWorkers"] = str(int(obj["BusyWorkers"]) - 1)
                    else:
                        obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) - 1)

            f.seek(0)
            f.truncate(0)
            json.dump(obj, f)
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)


def worker_abort(worker):
    stat = {
        "pid": worker.pid,
        "host": "",
        "method": "",
        "uri": "",
        "status": "SIGABRT",
    }

    with open(SERVER_STATUS, mode="r+") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            obj = json.load(f)

            stats = [(i, v) for i, v in enumerate(obj["stats"]) if v["pid"] == worker.pid]
            if len(stats) == 0:
                obj["stats"].append(stat)
                obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) - 1)
            else:
                for i, v in stats:
                    obj["stats"][i] = stat
                    if v["status"] == "A":
                        obj["BusyWorkers"] = str(int(obj["BusyWorkers"]) - 1)
                    else:
                        obj["IdleWorkers"] = str(int(obj["IdleWorkers"]) - 1)

            f.seek(0)
            f.truncate(0)
            json.dump(obj, f)
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)
```

## Example of mackerel-agent.conf

```
[plugin.metrics.gunicorn]
command = "/path/to/mackerel-plugin-gunicorn -status-file /dev/shm/gunicorn_status.json"
```

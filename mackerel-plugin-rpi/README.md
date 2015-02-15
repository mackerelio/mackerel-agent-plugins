mackerel-plugin-rpi
============================

## Description

Collecting Raspberry PI system metrics for Mackerel.

## Usage

### Build

First, you might want to build this plugin.

```
go get github.com/mackerelio/mackerel-agent-plugins
cd $GO_HOME/src/github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-rpi
go build
```

### Run

Then, you can run it.

```
./mackerel-plugin-rpi
```

### Configuration

If you want to get metrics via Mackerel, you have to edit your `mackerel-agent.conf`.
Below is an example.

```
[plugin.metrics.rpi]
command = "/path/to/mackerel-plugin-rpi"
```

Then, you need to restart `mackerel-agent` to apply your new settings.

## Author

[Takuya Arita](https://github.com/ariarijp)

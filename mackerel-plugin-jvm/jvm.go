package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.jvm")

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"jvm.gc_events": mp.Graphs{
		Label: "Number of GC events",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "YGC", Label: "Number of young generation GC events", Diff: true},
			mp.Metrics{Name: "FGC", Label: "Number of stop the world events", Diff: true},
		},
	},
	"jvm.gc_time": mp.Graphs{
		Label: "Garbage collection time",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "YGCT", Label: "Young generation garbage collection time", Diff: true},
			mp.Metrics{Name: "FGCT", Label: "Full garbage collection time", Diff: true},
		},
	},
	"jvm.survivor_space": mp.Graphs{
		Label: "Survivor Space (KB)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "S0C", Label: "Survivor space 0 capacity", Diff: false},
			mp.Metrics{Name: "S1C", Label: "Survivor space 1 capacity", Diff: false},
			mp.Metrics{Name: "S0U", Label: "Survivor space 0 utilization", Diff: false},
			mp.Metrics{Name: "S1C", Label: "Survivor space 1 utilization", Diff: false},
			mp.Metrics{Name: "DSS", Label: "Adequate size of survivor", Diff: false},
		},
	},
	"jvm.eden_space": mp.Graphs{
		Label: "Eden Space (KB)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "EC", Label: "Eden space capacity", Diff: false},
			mp.Metrics{Name: "EU", Label: "Eden space utilication", Diff: false},
		},
	},
	"jvm.old_space": mp.Graphs{
		Label: "Old Space (KB)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "OC", Label: "Old space capacity", Diff: false},
			mp.Metrics{Name: "OU", Label: "Old space utilication", Diff: false},
		},
	},
	"jvm.permanent_space": mp.Graphs{
		Label: "Permanent Space (KB)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "PC", Label: "Permanent space capacity", Diff: false},
			mp.Metrics{Name: "PU", Label: "Permanent space utilication", Diff: false},
		},
	},
	"jvm.new_area": mp.Graphs{
		Label: "New area",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "NGCMN", Label: "Minimum size", Diff: false},
			mp.Metrics{Name: "NGCMX", Label: "Maximum size", Diff: false},
			mp.Metrics{Name: "NGC", Label: "Current size", Diff: false},
		},
	},
	"jvm.old_area": mp.Graphs{
		Label: "Old area",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "OGCMN", Label: "Minimum size", Diff: false},
			mp.Metrics{Name: "OGCMX", Label: "Maximum size", Diff: false},
			mp.Metrics{Name: "OGC", Label: "Current size", Diff: false},
		},
	},
	"jvm.permanent_area": mp.Graphs{
		Label: "Permanent area",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "PGCMN", Label: "Minimum size", Diff: false},
			mp.Metrics{Name: "PGCMX", Label: "Maximum size", Diff: false},
			mp.Metrics{Name: "PGC", Label: "Current size", Diff: false},
		},
	},
	"jvm.tenuring_threshold": mp.Graphs{
		Label: "Tenuring threshold",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "TT", Label: "Tenuring threhold", Diff: false},
			mp.Metrics{Name: "PGCMX", Label: "Maximum tenuring threshold", Diff: false},
		},
	},
}

type JVMPlugin struct {
	Target    string
	Lvmid     string
	JstatPath string
	Tempfile  string
}

// # jps
// 26547 NettyServer
// 6438 Jps
func FetchLvmidByAppname(appname, target, jpsPath string) (string, error) {
	out, err := exec.Command(jpsPath, target).Output()
	if err != nil {
		logger.Errorf("Failed to run exec jps. %s", err)
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		words := strings.Split(line, " ")
		if len(words) != 2 {
			continue
		}
		lvmid, name := words[0], words[1]
		if name == appname {
			return lvmid, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Cannot get lvmid from %s", appname))
}

func fetchJstatMetrics(lvmid, option, jstatPath string) (map[string]float64, error) {
	out, err := exec.Command(jstatPath, option, lvmid).Output()
	if err != nil {
		logger.Errorf("Failed to run exec jstat. %s", err)
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	keys := strings.Fields(lines[0])
	values := strings.Fields(lines[1])

	stat := make(map[string]float64)
	for i, key := range keys {
		value, err := strconv.ParseFloat(values[i], 64)
		if err != nil {
			logger.Warningf("Failed to parse value. %s", err)
		}
		stat[key] = value
	}

	return stat, nil
}

func mergeStat(dst, src map[string]float64) {
	for k, v := range src {
		dst[k] = v
	}
}

// # jstat -gc <vmid>
//  S0C    S1C    S0U    S1U      EC       EU        OC         OU       PC     PU    YGC     YGCT    FGC    FGCT     GCT
// 3584.0 3584.0 2528.0  0.0   692224.0 19062.4  1398272.0   485450.1  72704.0 72611.3   3152   30.229   0      0.000   30.229

// # jstat -gccapacity  <vmid>
//  NGCMN    NGCMX     NGC     S0C   S1C       EC      OGCMN      OGCMX       OGC         OC      PGCMN    PGCMX     PGC       PC     YGC    FGC
// 699392.0 699392.0 699392.0 4096.0 4096.0 691200.0  1398272.0  1398272.0  1398272.0  1398272.0  21504.0 524288.0  72704.0  72704.0   4212     0

// # jstat -gcnew  <vmid>
//  S0C    S1C    S0U    S1U   TT MTT  DSS      EC       EU     YGC     YGCT
// 3072.0 3072.0    0.0 2848.0  1  15 3072.0 693248.0 626782.2   3463   33.658

func (m JVMPlugin) FetchMetrics() (map[string]float64, error) {
	gcStat, err := fetchJstatMetrics(m.Lvmid, "-gc", m.JstatPath)
	if err != nil {
		return nil, err
	}
	gcCapacityStat, err := fetchJstatMetrics(m.Lvmid, "-gccapacity", m.JstatPath)
	if err != nil {
		return nil, err
	}
	gcNewStat, err := fetchJstatMetrics(m.Lvmid, "-gcnew", m.JstatPath)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	mergeStat(stat, gcStat)
	mergeStat(stat, gcCapacityStat)
	mergeStat(stat, gcNewStat)

	return stat, nil
}

func (m JVMPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "1099", "Port")
	optJstatPath := flag.String("jstatpath", "/usr/bin/jstat", "jstat path")
	optJpsPath := flag.String("jpspath", "/usr/bin/jps", "jps path")
	optJavaName := flag.String("javaname", "", "Java app name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	if *optJavaName == "" {
		logger.Errorf("javaname is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var jvm JVMPlugin
	jvm.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	lvmid, err := FetchLvmidByAppname(*optJavaName, jvm.Target, *optJpsPath)
	if err != nil {
		logger.Errorf("Failed to fetch lvmid. %s", err)
		os.Exit(1)
	}
	jvm.Lvmid = lvmid
	jvm.JstatPath = *optJstatPath

	helper := mp.NewMackerelPlugin(jvm)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-jvm-%s", *optHost)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

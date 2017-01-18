package mpjvm

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.jvm")

// JVMPlugin plugin for JVM
type JVMPlugin struct {
	Target    string
	Lvmid     string
	JstatPath string
	JinfoPath string
	JavaName  string
	Tempfile  string
}

// # jps
// 26547 NettyServer
// 6438 Jps
func fetchLvmidByAppname(appname, target, jpsPath string) (string, error) {
	stdout, _, exitStatus, err := runTimeoutCommand(jpsPath, target)

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("jps command timed out")
	}
	if err != nil {
		logger.Errorf("Failed to run exec jps. %s. Please run with the java process user.", err)
		return "", err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		words := strings.Split(line, " ")
		if len(words) != 2 {
			continue
		}
		lvmid, name := words[0], words[1]
		if name == appname {
			return lvmid, nil
		}
	}
	return "", fmt.Errorf("cannot get lvmid from %s (please run with the java process user)", appname)
}

func fetchJstatMetrics(lvmid, option, jstatPath string) (map[string]float64, error) {
	stdout, _, exitStatus, err := runTimeoutCommand(jstatPath, option, lvmid)

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("jstat command timed out")
	}
	if err != nil {
		logger.Errorf("Failed to run exec jstat. %s. Please run with the java process user.", err)
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")
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

func calculateMemorySpaceRate(gcStat map[string]float64, m JVMPlugin) (map[string]float64, error) {
	ret := make(map[string]float64)
	ret["oldSpaceRate"] = gcStat["OU"] / gcStat["OC"] * 100
	ret["newSpaceRate"] = (gcStat["S0U"] + gcStat["S1U"] + gcStat["EU"]) / (gcStat["S0C"] + gcStat["S1C"] + gcStat["EC"]) * 100
	if checkCMSGC(m.Lvmid, m.JinfoPath) {
		ret["CMSInitiatingOccupancyFraction"] = fetchCMSInitiatingOccupancyFraction(m.Lvmid, m.JinfoPath)
	}

	return ret, nil
}

func checkCMSGC(lvmid, JinfoPath string) bool {
	stdout, _, exitStatus, err := runTimeoutCommand(JinfoPath, "-flag", "UseConcMarkSweepGC", lvmid)

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("jinfo command timed out")
		os.Exit(1)
	}
	if err != nil {
		logger.Errorf("Failed to run exec jinfo. %s. Please run with the java process user.", err)
		os.Exit(1)
	}
	return strings.Index(string(stdout), "+UseConcMarkSweepGC") != -1
}

func fetchCMSInitiatingOccupancyFraction(lvmid, JinfoPath string) float64 {
	var fraction float64

	stdout, _, exitStatus, err := runTimeoutCommand(JinfoPath, "-flag", "CMSInitiatingOccupancyFraction", lvmid)

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("jinfo command timed out")
	}
	if err != nil {
		logger.Errorf("Failed to run exec jinfo. %s. Please run with the java process user.", err)
		os.Exit(1)
	}

	out := strings.Trim(string(stdout), "\n")
	tmp := strings.Split(out, "=")
	fraction, _ = strconv.ParseFloat(tmp[1], 64)

	return fraction
}

func mergeStat(dst, src map[string]float64) {
	for k, v := range src {
		dst[k] = v
	}
}

func runTimeoutCommand(Path string, Args ...string) (string, string, timeout.ExitStatus, error) {
	var TimeoutDuration = 10 * time.Second
	var TimeoutKillAfter = 5 * time.Second
	tio := &timeout.Timeout{
		Cmd:       exec.Command(Path, Args...),
		Duration:  TimeoutDuration,
		KillAfter: TimeoutKillAfter,
	}
	exitStatus, stdout, stderr, err := tio.Run()
	return stdout, stderr, exitStatus, err
}

// <Java8> https://docs.oracle.com/javase/8/docs/technotes/tools/unix/jstat.html
// # jstat -gc <vmid>
//  S0C    S1C    S0U    S1U      EC       EU        OC         OU       MC     MU    CCSC   CCSU   YGC     YGCT    FGC    FGCT     GCT
// 1024.0 1024.0  0.0    0.0    8256.0   8256.0   20480.0     453.4    4864.0 2776.2 512.0  300.8       0    0.000   1      0.003    0.003

// # jstat -gccapacity <vmid>
//  NGCMN    NGCMX     NGC     S0C   S1C       EC      OGCMN      OGCMX       OGC         OC       MCMN     MCMX      MC     CCSMN    CCSMX     CCSC    YGC    FGC
//  10240.0 160384.0  10304.0 1024.0 1024.0   8256.0    20480.0   320896.0    20480.0    20480.0      0.0 1056768.0   4864.0      0.0 1048576.0    512.0      0     1

// # jstat -gcnew <vmid>
//  S0C    S1C    S0U    S1U   TT MTT  DSS      EC       EU     YGC     YGCT
// 1024.0 1024.0    0.0    0.0 15  15    0.0   8256.0   8256.0      0    0.000

// <Java7>
// # jstat -gc <vmid>
//  S0C    S1C    S0U    S1U      EC       EU        OC         OU       PC     PU    YGC     YGCT    FGC    FGCT     GCT
// 3584.0 3584.0 2528.0  0.0   692224.0 19062.4  1398272.0   485450.1  72704.0 72611.3   3152   30.229   0      0.000   30.229

// # jstat -gccapacity  <vmid>
//  NGCMN    NGCMX     NGC     S0C   S1C       EC      OGCMN      OGCMX       OGC         OC      PGCMN    PGCMX     PGC       PC     YGC    FGC
// 699392.0 699392.0 699392.0 4096.0 4096.0 691200.0  1398272.0  1398272.0  1398272.0  1398272.0  21504.0 524288.0  72704.0  72704.0   4212     0

// # jstat -gcnew  <vmid>
//  S0C    S1C    S0U    S1U   TT MTT  DSS      EC       EU     YGC     YGCT
// 3072.0 3072.0    0.0 2848.0  1  15 3072.0 693248.0 626782.2   3463   33.658

// FetchMetrics interface for mackerelplugin
func (m JVMPlugin) FetchMetrics() (map[string]interface{}, error) {
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
	gcOldStat, err := fetchJstatMetrics(m.Lvmid, "-gcold", m.JstatPath)
	if err != nil {
		return nil, err
	}
	gcSpaceRate, err := calculateMemorySpaceRate(gcStat, m)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	mergeStat(stat, gcStat)
	mergeStat(stat, gcCapacityStat)
	mergeStat(stat, gcNewStat)
	mergeStat(stat, gcOldStat)
	mergeStat(stat, gcSpaceRate)

	result := make(map[string]interface{})
	for k, v := range stat {
		result[k] = v
	}
	return result, nil
}

// GraphDefinition interface for mackerelplugin
func (m JVMPlugin) GraphDefinition() map[string]mp.Graphs {
	rawJavaName := m.JavaName
	lowerJavaName := strings.ToLower(m.JavaName)
	return map[string]mp.Graphs{
		fmt.Sprintf("jvm.%s.gc_events", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s GC events", rawJavaName),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "YGC", Label: "Young GC event", Diff: true},
				{Name: "FGC", Label: "Full GC event", Diff: true},
			},
		},
		fmt.Sprintf("jvm.%s.gc_time", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s GC time (sec)", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "YGCT", Label: "Young GC time", Diff: true},
				{Name: "FGCT", Label: "Full GC time", Diff: true},
			},
		},
		fmt.Sprintf("jvm.%s.gc_time_percentage", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s GC time percentage", rawJavaName),
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				// gc_time_percentage is the percentage of gc time to 60 sec.
				{Name: "YGCT", Label: "Young GC time", Diff: true, Scale: (100.0 / 60)},
				{Name: "FGCT", Label: "Full GC time", Diff: true, Scale: (100.0 / 60)},
			},
		},
		fmt.Sprintf("jvm.%s.new_space", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s New Space memory", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "NGCMX", Label: "New max", Diff: false, Scale: 1024},
				{Name: "NGC", Label: "New current", Diff: false, Scale: 1024},
				{Name: "EU", Label: "Eden used", Diff: false, Scale: 1024},
				{Name: "S0U", Label: "Survivor0 used", Diff: false, Scale: 1024},
				{Name: "S1U", Label: "Survivor1 used", Diff: false, Scale: 1024},
			},
		},
		fmt.Sprintf("jvm.%s.old_space", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s Old Space memory", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "OGCMX", Label: "Old max", Diff: false, Scale: 1024},
				{Name: "OGC", Label: "Old current", Diff: false, Scale: 1024},
				{Name: "OU", Label: "Old used", Diff: false, Scale: 1024},
			},
		},
		fmt.Sprintf("jvm.%s.perm_space", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s Permanent Space", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "PGCMX", Label: "Perm max", Diff: false, Scale: 1024},
				{Name: "PGC", Label: "Perm current", Diff: false, Scale: 1024},
				{Name: "PU", Label: "Perm used", Diff: false, Scale: 1024},
			},
		},
		fmt.Sprintf("jvm.%s.metaspace", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s Metaspace", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "MCMX", Label: "Metaspace capacity max", Diff: false, Scale: 1024},
				{Name: "MCMN", Label: "Metaspace capacity min", Diff: false, Scale: 1024},
				{Name: "MC", Label: "Metaspace capacity", Diff: false, Scale: 1024},
				{Name: "MU", Label: "Metaspace utilization ", Diff: false, Scale: 1024},
				{Name: "CCSC", Label: "Compressed Class Space Capacity", Diff: false, Scale: 1024},
				{Name: "CCSU", Label: "Compressed Class Space Used", Diff: false, Scale: 1024},
			},
		},
		fmt.Sprintf("jvm.%s.memorySpace", lowerJavaName): {
			Label: fmt.Sprintf("JVM %s MemorySpace", rawJavaName),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "oldSpaceRate", Label: "GC Old Memory Space", Diff: false},
				{Name: "newSpaceRate", Label: "GC New Memory Space", Diff: false},
				{Name: "CMSInitiatingOccupancyFraction", Label: "CMS Initiating Occupancy Fraction", Diff: false},
			},
		},
	}
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "1099", "Port")
	optJstatPath := flag.String("jstatpath", "/usr/bin/jstat", "jstat path")
	optJinfoPath := flag.String("jinfopath", "/usr/bin/jinfo", "jinfo path")
	optJpsPath := flag.String("jpspath", "/usr/bin/jps", "jps path")
	optJavaName := flag.String("javaname", "", "Java app name")
	optPidFile := flag.String("pidfile", "", "pidfile path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var jvm JVMPlugin
	jvm.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	jvm.JstatPath = *optJstatPath
	jvm.JinfoPath = *optJinfoPath

	if *optJavaName == "" {
		logger.Errorf("javaname is required (if you use 'pidfile' option, 'javaname' is used as just a prefix of graph label)")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *optPidFile == "" {
		lvmid, err := fetchLvmidByAppname(*optJavaName, jvm.Target, *optJpsPath)
		if err != nil {
			logger.Errorf("Failed to fetch lvmid. %s. Please run with the java process user.", err)
			os.Exit(1)
		}
		jvm.Lvmid = lvmid
	} else {
		// https://docs.oracle.com/javase/7/docs/technotes/tools/share/jps.html
		// `The lvmid is typically, but not necessarily, the operating system's process identifier for the JVM process.`
		pid, err := ioutil.ReadFile(*optPidFile)
		if err != nil {
			logger.Errorf("Failed to load pid. %s", err)
			os.Exit(1)
		}
		jvm.Lvmid = strings.Replace(string(pid), "\n", "", 1)
	}

	jvm.JavaName = *optJavaName

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

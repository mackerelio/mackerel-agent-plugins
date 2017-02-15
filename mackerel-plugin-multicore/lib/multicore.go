package mpmulticore

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"path/filepath"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphDef = map[string]mp.Graphs{
	"multicore.cpu.#": {
		Label: "MultiCore CPU",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "user", Label: "user", Diff: false, Stacked: true},
			{Name: "nice", Label: "nice", Diff: false, Stacked: true},
			{Name: "system", Label: "system", Diff: false, Stacked: true},
			{Name: "idle", Label: "idle", Diff: false, Stacked: true},
			{Name: "iowait", Label: "ioWait", Diff: false, Stacked: true},
			{Name: "irq", Label: "irq", Diff: false, Stacked: true},
			{Name: "softirq", Label: "softirq", Diff: false, Stacked: true},
			{Name: "steal", Label: "steal", Diff: false, Stacked: true},
			{Name: "guest", Label: "guest", Diff: false, Stacked: true},
			{Name: "guest_nice", Label: "guest_nice", Diff: false, Stacked: true},
		},
	},
	"multicore.loadavg_per_core": {
		Label: "MultiCore loadavg5 per core",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "loadavg5", Label: "loadavg5", Diff: false, Stacked: false},
		},
	},
}

type saveItem struct {
	LastTime       time.Time
	ProcStatsByCPU map[string]*procStats
}

type procStats struct {
	User      *float64 `json:"user"`
	Nice      *float64 `json:"nice"`
	System    *float64 `json:"system"`
	Idle      *float64 `json:"idle"`
	IoWait    *float64 `json:"iowait"`
	Irq       *float64 `json:"irq"`
	SoftIrq   *float64 `json:"softirq"`
	Steal     *float64 `json:"steal"`
	Guest     *float64 `json:"guest"`
	GuestNice *float64 `json:"guest_nice"`
	Total     float64  `json:"total"`
}

type cpuPercentages struct {
	GroupName string
	User      *float64
	Nice      *float64
	System    *float64
	Idle      *float64
	IoWait    *float64
	Irq       *float64
	SoftIrq   *float64
	Steal     *float64
	Guest     *float64
	GuestNice *float64
}

func getProcStat() (string, error) {
	contentbytes, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return "", err
	}
	return string(contentbytes), nil
}

func parseFloats(values []string) ([]float64, error) {
	var result []float64
	for _, v := range values {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, nil
}

func fill(arr []float64, elementCount int) []*float64 {
	var filled []*float64
	for _, v := range arr {
		copy := v
		filled = append(filled, &copy)
	}

	if len(arr) < elementCount {
		emptyArray := make([]*float64, elementCount-len(arr))
		filled = append(filled, emptyArray...)
	}
	return filled
}

func parseProcStat(str string) (map[string]*procStats, error) {
	var result = make(map[string]*procStats)
	for _, line := range strings.Split(str, "\n") {
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			key := fields[0]
			values := fields[1:]

			// skip total cpu usage
			if key == "cpu" {
				continue
			}

			floatValues, err := parseFloats(values)
			if err != nil {
				return nil, err
			}

			total := 0.0
			for _, v := range floatValues {
				total += v
			}

			filledValues := fill(floatValues, 10)

			ps := &procStats{
				User:      filledValues[0],
				Nice:      filledValues[1],
				System:    filledValues[2],
				Idle:      filledValues[3],
				IoWait:    filledValues[4],
				Irq:       filledValues[5],
				SoftIrq:   filledValues[6],
				Steal:     filledValues[7],
				Guest:     filledValues[8],
				GuestNice: filledValues[9],
				Total:     total,
			}
			result[key] = ps
		} else {
			break
		}
	}
	return result, nil
}

func collectProcStatValues() (map[string]*procStats, error) {
	procStats, err := getProcStat()
	if err != nil {
		return nil, err
	}
	return parseProcStat(procStats)
}

func saveValues(tempFileName string, values map[string]*procStats, now time.Time) error {
	f, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	s := saveItem{
		LastTime:       now,
		ProcStatsByCPU: values,
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(s)
	if err != nil {
		return err
	}

	return nil
}

func fetchSavedItem(tempFileName string) (*saveItem, error) {
	f, err := os.Open(tempFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var stat saveItem
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&stat)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func calcCPUUsage(currentValues map[string]*procStats, now time.Time, savedItem *saveItem) ([]*cpuPercentages, error) {
	lastValues := savedItem.ProcStatsByCPU
	lastTime := savedItem.LastTime

	var result []*cpuPercentages
	for key, current := range currentValues {
		last, ok := lastValues[key]
		if ok {
			user, err := calcPercentage(current.User, last.User, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			nice, err := calcPercentage(current.Nice, last.Nice, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			system, err := calcPercentage(current.System, last.System, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			idle, err := calcPercentage(current.Idle, last.Idle, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			iowait, err := calcPercentage(current.IoWait, last.IoWait, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			irq, err := calcPercentage(current.Irq, last.Irq, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			softirq, err := calcPercentage(current.SoftIrq, last.SoftIrq, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			steal, err := calcPercentage(current.Steal, last.Steal, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			guest, err := calcPercentage(current.Guest, last.Guest, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}
			guestNice, err := calcPercentage(current.GuestNice, last.GuestNice, current.Total, last.Total, now, lastTime)
			if err != nil {
				return nil, err
			}

			p := &cpuPercentages{
				GroupName: key,
				User:      user,
				Nice:      nice,
				System:    system,
				Idle:      idle,
				IoWait:    iowait,
				Irq:       irq,
				SoftIrq:   softirq,
				Steal:     steal,
				Guest:     guest,
				GuestNice: guestNice,
			}
			result = append(result, p)
		}

	}

	return result, nil
}

func calcPercentage(currentValue *float64, lastValue *float64, currentTotal float64, lastTotal float64, now time.Time, lastTime time.Time) (*float64, error) {

	if currentValue == nil || lastValue == nil {
		return nil, nil
	}

	value, err := calcDiff(*currentValue, now, *lastValue, lastTime)
	if err != nil {
		return nil, err
	}

	total, err := calcDiff(currentTotal, now, lastTotal, lastTime)
	if err != nil {
		return nil, err
	}

	ret := value / total * 100.0
	return &ret, nil
}

func calcDiff(value float64, now time.Time, lastValue float64, lastTime time.Time) (float64, error) {
	diffTime := now.Unix() - lastTime.Unix()
	if diffTime > 600 {
		return 0.0, fmt.Errorf("Too long duration")
	}

	diff := (value - lastValue) * 60 / float64(diffTime)

	if lastValue <= value {
		return diff, nil
	}
	return 0.0, fmt.Errorf("lastValue > value")
}

func fetchLoadavg5() (float64, error) {
	contentbytes, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0.0, err
	}
	content := string(contentbytes)
	cols := strings.Fields(content)

	if len(cols) > 2 {
		f, err := strconv.ParseFloat(cols[1], 64)
		if err != nil {
			return 0.0, err
		}
		return f, nil
	}
	return 0.0, fmt.Errorf("cannot fetch loadavg5")
}

func printValue(key string, value *float64, time time.Time) {
	if value != nil {
		fmt.Printf("%s\t%f\t%d\n", key, *value, time.Unix())
	}
}

func outputCPUUsage(cpuUsage []*cpuPercentages, now time.Time) {
	for _, u := range cpuUsage {
		printValue(fmt.Sprintf("multicore.cpu.%s.user", u.GroupName), u.User, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.nice", u.GroupName), u.Nice, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.system", u.GroupName), u.System, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.idle", u.GroupName), u.Idle, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.iowait", u.GroupName), u.IoWait, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.irq", u.GroupName), u.Irq, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.softirq", u.GroupName), u.SoftIrq, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.steal", u.GroupName), u.Steal, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.guest", u.GroupName), u.Guest, now)
		printValue(fmt.Sprintf("multicore.cpu.%s.guest_nice", u.GroupName), u.GuestNice, now)
	}
}

func outputLoadavgPerCore(loadavgPerCore float64, now time.Time) {
	printValue("multicore.loadavg_per_core.loadavg5", &loadavgPerCore, now)
}

func outputDefinitions() {
	fmt.Println("# mackerel-agent-plugin")
	var graphs mp.GraphDef
	graphs.Graphs = graphDef

	b, err := json.Marshal(graphs)
	if err != nil {
		log.Fatalln("OutputDefinitions: ", err)
	}
	fmt.Println(string(b))
}

func outputMulticore(tempFileName string) {
	now := time.Now()

	currentValues, err := collectProcStatValues()
	if err != nil {
		log.Fatalln("collectProcStatValues: ", err)
	}

	savedItem, err := fetchSavedItem(tempFileName)
	saveValues(tempFileName, currentValues, now)
	if err != nil {
		log.Fatalln("fetchLastValues: ", err)
	}

	// maybe first time run
	if savedItem == nil {
		return
	}

	cpuUsage, err := calcCPUUsage(currentValues, now, savedItem)
	if err != nil {
		log.Fatalln("calcCPUUsage: ", err)
	}

	loadavg5, err := fetchLoadavg5()
	if err != nil {
		log.Fatalln("fetchLoadavg5: ", err)
	}
	loadPerCPUCount := loadavg5 / (float64(len(cpuUsage)))

	outputCPUUsage(cpuUsage, now)
	outputLoadavgPerCore(loadPerCPUCount, now)
}

func generateTempfilePath() string {
	dir := os.Getenv("MACKEREL_PLUGIN_WORKDIR")
	if dir == "" {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "mackerel-plugin-multicore")
}

// Do the plugin
func Do() {
	var tempFileName string
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	tempFileName = *optTempfile
	if tempFileName == "" {
		tempFileName = generateTempfilePath()
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		outputDefinitions()
	} else {
		outputMulticore(tempFileName)
	}
}

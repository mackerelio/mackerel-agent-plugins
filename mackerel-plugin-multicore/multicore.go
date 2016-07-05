package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphDef = map[string](mp.Graphs){
	"multicore.cpu.#": mp.Graphs{
		Label: "MultiCore CPU",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "user", Label: "user", Diff: false, Stacked: true},
			mp.Metrics{Name: "nice", Label: "nice", Diff: false, Stacked: true},
			mp.Metrics{Name: "system", Label: "system", Diff: false, Stacked: true},
			mp.Metrics{Name: "idle", Label: "idle", Diff: false, Stacked: true},
			mp.Metrics{Name: "iowait", Label: "ioWait", Diff: false, Stacked: true},
			mp.Metrics{Name: "irq", Label: "irq", Diff: false, Stacked: true},
			mp.Metrics{Name: "softirq", Label: "softirq", Diff: false, Stacked: true},
			mp.Metrics{Name: "steal", Label: "steal", Diff: false, Stacked: true},
			mp.Metrics{Name: "guest", Label: "guest", Diff: false, Stacked: true},
		},
	},
	"multicore.loadavg_per_core": mp.Graphs{
		Label: "MultiCore loadavg5 per core",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "loadavg5", Label: "loadavg5", Diff: false, Stacked: true},
		},
	},
}

type saveItem struct {
	LastTime       time.Time
	ProcStatsByCPU map[string]*procStats
}

type procStats struct {
	User    float64 `json:"user"`
	Nice    float64 `json:"nice"`
	System  float64 `json:"system"`
	Idle    float64 `json:"idle"`
	IoWait  float64 `json:"iowait"`
	Irq     float64 `json:"irq"`
	SoftIrq float64 `json:"softirq"`
	Steal   float64 `json:"steal"`
	Guest   float64 `json:"guest"`
	Total   float64 `json:"total"`
}

type cpuPercentages struct {
	GroupName string
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	IoWait    float64
	Irq       float64
	SoftIrq   float64
	Steal     float64
	Guest     float64
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

func fill(arr []float64, elementCount int) []float64 {
	arraySize := len(arr)
	if arraySize < elementCount {
		for i := arraySize; i < elementCount; i++ {
			arr = append(arr, 0.0)
		}
	}
	return arr
}

func parseProcStat(str string) (map[string]*procStats, error) {
	var result = make(map[string]*procStats)
	for _, line := range strings.Split(str, "\n") {
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			key := fields[0]
			values := fields[1:]

			floatValues, err := parseFloats(values)
			if err != nil {
				return nil, err
			}
			filledValues := fill(floatValues, 9)

			total := 0.0
			for _, v := range floatValues {
				total += v
			}

			ps := &procStats{
				User:    filledValues[0],
				Nice:    filledValues[1],
				System:  filledValues[2],
				Idle:    filledValues[3],
				IoWait:  filledValues[4],
				Irq:     filledValues[5],
				SoftIrq: filledValues[6],
				Steal:   filledValues[7],
				Guest:   filledValues[8],
				Total:   total,
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
		LastTime:       time.Now(),
		ProcStatsByCPU: values,
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(s)
	if err != nil {
		return err
	}

	return nil
}

func fetchLastValues(tempFileName string) (map[string]*procStats, time.Time, error) {
	f, err := os.Open(tempFileName)
	if err != nil {
		return nil, time.Now(), err
	}
	defer f.Close()

	var stat saveItem
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&stat)
	if err != nil {
		return stat.ProcStatsByCPU, stat.LastTime, err
	}
	return stat.ProcStatsByCPU, stat.LastTime, nil
}

func calcCPUUsage(currentValues map[string]*procStats, now time.Time, lastValues map[string]*procStats, lastTime time.Time) ([]*cpuPercentages, error) {

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
			}
			result = append(result, p)
		}

	}

	return result, nil
}

func calcPercentage(currentValue float64, lastValue float64, currentTotal float64, lastTotal float64, now time.Time, lastTime time.Time) (float64, error) {
	value, err := calcDiff(currentValue, now, lastValue, lastTime)
	if err != nil {
		return 0.0, err
	}

	total, err := calcDiff(currentTotal, now, lastTotal, lastTime)
	if err != nil {
		return 0.0, err
	}

	return (value / total * 100.0), nil
}

func calcDiff(value float64, now time.Time, lastValue float64, lastTime time.Time) (float64, error) {
	diffTime := now.Unix() - lastTime.Unix()
	if diffTime > 600 {
		return 0.0, errors.New("Too long duration")
	}

	diff := (value - lastValue) * 60 / float64(diffTime)

	if lastValue <= value {
		return diff, nil
	}
	return 0.0, nil
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
	return 0.0, errors.New("cannot fetch loadavg5.")
}

func printValue(key string, value float64, time time.Time) {
	fmt.Printf("%s\t%f\t%d\n", key, value, time.Unix())
}

func outputCPUUsage(cpuUsage []*cpuPercentages, now time.Time) {
	if cpuUsage != nil {
		for _, u := range cpuUsage {
			if u.GroupName != "cpu" {
				printValue(fmt.Sprintf("multicore.cpu.%s.user", u.GroupName), u.User, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.nice", u.GroupName), u.Nice, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.system", u.GroupName), u.System, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.idle", u.GroupName), u.Idle, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.iowait", u.GroupName), u.IoWait, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.irq", u.GroupName), u.Irq, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.softirq", u.GroupName), u.SoftIrq, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.steal", u.GroupName), u.Steal, now)
				printValue(fmt.Sprintf("multicore.cpu.%s.guest", u.GroupName), u.Guest, now)
			}
		}
	}
}

func outputLoadavgPerCore(loadavgPerCore float64, now time.Time) {
	printValue("multicore.loadavg_per_core.loadavg5", loadavgPerCore, now)
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

// main function
func main() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	var tempFileName string
	if *optTempfile != "" {
		tempFileName = *optTempfile
	} else {
		tempFileName = "/tmp/mackerel-plugin-multicore"
	}
	now := time.Now()

	currentValues, _ := collectProcStatValues()
	lastValues, lastTime, err := fetchLastValues(tempFileName)
	saveValues(tempFileName, currentValues, now)
	if err != nil {
		log.Fatalln("fetchLastValues: ", err)
	}

	var cpuUsage []*cpuPercentages
	if lastValues != nil {
		var err error
		cpuUsage, err = calcCPUUsage(currentValues, now, lastValues, lastTime)
		if err != nil {
			log.Fatalln("calcCPUUsage: ", err)
		}
	}

	loadavg5, err := fetchLoadavg5()
	if err != nil {
		log.Fatalln("fetchLoadavg5: ", err)
	}
	loadPerCPUCount := loadavg5 / (float64(len(cpuUsage) - 1))

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		outputDefinitions()
	} else {
		outputCPUUsage(cpuUsage, now)
		outputLoadavgPerCore(loadPerCPUCount, now)
	}
}

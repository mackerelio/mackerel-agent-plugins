package mplinux

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/urfave/cli"
)

const (
	pathVmstat = "/proc/vmstat"
	pathStat   = "/proc/stat"
	pathSysfs  = "/sys"
)

var collectVirtualDevice = regexp.MustCompile("^fio[a-z]+$") // ioDrive(FusionIO)

// metric value structure
// note: all metrics are add dynamic at collect*().
var graphdef = map[string]mp.Graphs{}

// LinuxPlugin mackerel plugin for linux
type LinuxPlugin struct {
	Tempfile string
	Typemap  map[string]bool
}

// GraphDefinition interface for mackerelplugin
func (c LinuxPlugin) GraphDefinition() map[string]mp.Graphs {
	var err error

	p := make(map[string]interface{})

	if c.Typemap["all"] || c.Typemap["swap"] {
		err = collectProcVmstat(pathVmstat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["netstat"] {
		err = collectNetworkStat(&p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["diskstats"] {
		err = collectDiskStats(pathSysfs, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["proc_stat"] {
		err = collectProcStat(pathStat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["users"] {
		err = collectWho(&p)
		if err != nil {
			return nil
		}
	}

	return graphdef
}

// main function
func doMain(c *cli.Context) error {
	var linux LinuxPlugin

	typemap := map[string]bool{}
	types := c.StringSlice("type")
	// If no `type` is specified, fetch all metrics
	if len(types) == 0 {
		typemap["all"] = true
	} else {
		for _, t := range types {
			typemap[t] = true
		}
	}
	linux.Typemap = typemap
	helper := mp.NewMackerelPlugin(linux)
	helper.Tempfile = c.String("tempfile")

	helper.Run()
	return nil
}

// FetchMetrics interface for mackerelplugin
func (c LinuxPlugin) FetchMetrics() (map[string]interface{}, error) {
	var err error

	p := make(map[string]interface{})

	if c.Typemap["all"] || c.Typemap["swap"] {
		err = collectProcVmstat(pathVmstat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["netstat"] {
		err = collectNetworkStat(&p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["diskstats"] {
		err = collectDiskStats(pathSysfs, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["proc_stat"] {
		err = collectProcStat(pathStat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["users"] {
		err = collectWho(&p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// collect who
func collectWho(p *map[string]interface{}) error {
	var err error
	var data string

	graphdef["linux.users"] = mp.Graphs{
		Label: "Linux Users",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "users", Label: "Users", Diff: false},
		},
	}

	data, err = getWho()
	if err != nil {
		return err
	}
	err = parseWho(data, p)
	if err != nil {
		return err
	}

	return nil
}

// parsing metrics from /proc/stat
func parseWho(str string, p *map[string]interface{}) error {
	str = strings.TrimSpace(str)
	if str == "" {
		(*p)["users"] = 0
		return nil
	}
	line := strings.Split(str, "\n")
	(*p)["users"] = float64(len(line))

	return nil
}

// Getting who
func getWho() (string, error) {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// collect /proc/stat
func collectProcStat(path string, p *map[string]interface{}) error {
	graphdef["linux.interrupts"] = mp.Graphs{
		Label: "Linux Interrupts",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "interrupts", Label: "Interrupts", Diff: true},
		},
	}
	graphdef["linux.context_switches"] = mp.Graphs{
		Label: "Linux Context Switches",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "context_switches", Label: "Context Switches", Diff: true},
		},
	}
	graphdef["linux.forks"] = mp.Graphs{
		Label: "Linux Forks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "forks", Label: "Forks", Diff: true},
		},
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return parseProcStat(file, p)
}

// parsing metrics from /proc/stat
func parseProcStat(r io.Reader, p *map[string]interface{}) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Fields(line)
		if len(record) < 2 {
			continue
		}
		name := record[0]
		value, errParse := atof(record[1])
		if errParse != nil {
			return errParse
		}

		if name == "intr" {
			(*p)["interrupts"] = value
		} else if name == "ctxt" {
			(*p)["context_switches"] = value
		} else if name == "processes" {
			(*p)["forks"] = value
		}
	}

	return nil
}

// collect /sys/block/<device>/stat
// See also. http://man7.org/linux/man-pages/man5/sysfs.5.html
func collectDiskStats(path string, p *map[string]interface{}) error {
	var elapsedData []mp.Metrics
	var rwtimeData []mp.Metrics

	sysBlockDir := filepath.Join(path, "block")

	devices, err := ioutil.ReadDir(sysBlockDir)
	if err != nil {
		return err
	}

	for _, d := range devices {
		if d.Mode()&os.ModeSymlink != os.ModeSymlink {
			continue
		}

		name := d.Name()

		// /sys/block/<device> is a symbolic link for block device
		realPath, err := filepath.EvalSymlinks(filepath.Join(sysBlockDir, name))
		if err != nil {
			return err
		}

		// exclude virtual device
		if strings.Index(realPath, "/devices/virtual/") != -1 {
			if !collectVirtualDevice.Match([]byte(name)) {
				continue
			}
		}

		// exclude removable device
		content, err := ioutil.ReadFile(filepath.Join(realPath, "removable"))
		if err != nil {
			return err
		}
		if len(content) > 0 && string(content[0]) == "1" {
			continue
		}

		content, err = ioutil.ReadFile(filepath.Join(realPath, "stat"))
		if err != nil {
			return err
		}

		err = parseDiskStat(name, string(content), p)
		if err != nil {
			return err
		}

		elapsedData = append(elapsedData, mp.Metrics{Name: fmt.Sprintf("iotime_%s", name), Label: fmt.Sprintf("%s IO Time", name), Diff: true})
		elapsedData = append(elapsedData, mp.Metrics{Name: fmt.Sprintf("iotime_weighted_%s", name), Label: fmt.Sprintf("%s IO Time Weighted", name), Diff: true})

		rwtimeData = append(rwtimeData, mp.Metrics{Name: fmt.Sprintf("tsreading_%s", name), Label: fmt.Sprintf("%s Read", name), Diff: true})
		rwtimeData = append(rwtimeData, mp.Metrics{Name: fmt.Sprintf("tswriting_%s", name), Label: fmt.Sprintf("%s Write", name), Diff: true})
	}

	graphdef["linux.disk.elapsed"] = mp.Graphs{
		Label:   "Disk Elapsed IO Time",
		Unit:    "integer",
		Metrics: elapsedData,
	}

	graphdef["linux.disk.rwtime"] = mp.Graphs{
		Label:   "Disk Read/Write Time",
		Unit:    "integer",
		Metrics: rwtimeData,
	}

	return nil
}

func parseDiskStat(name, stat string, p *map[string]interface{}) error {
	fields := strings.Fields(stat)
	if len(fields) != 11 {
		return nil
	}

	// See also. https://www.kernel.org/doc/Documentation/block/stat.txt
	(*p)[fmt.Sprintf("iotime_%s", name)], _ = atof(fields[9])           // io_ticks
	(*p)[fmt.Sprintf("iotime_weighted_%s", name)], _ = atof(fields[10]) // time_in_queue
	(*p)[fmt.Sprintf("tsreading_%s", name)], _ = atof(fields[3])        // read ticks
	(*p)[fmt.Sprintf("tswriting_%s", name)], _ = atof(fields[7])        // write ticks

	return nil
}

// collect ss
func collectNetworkStat(p *map[string]interface{}) error {
	graphdef["linux.ss"] = mp.Graphs{
		Label: "Linux Network Connection States",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ESTAB", Label: "Established", Diff: false, Stacked: true},
			{Name: "SYN-SENT", Label: "Syn Sent", Diff: false, Stacked: true},
			{Name: "SYN-RECV", Label: "Syn Received", Diff: false, Stacked: true},
			{Name: "FIN-WAIT-1", Label: "Fin Wait 1", Diff: false, Stacked: true},
			{Name: "FIN-WAIT-2", Label: "Fin Wait 2", Diff: false, Stacked: true},
			{Name: "TIME-WAIT", Label: "Time Wait", Diff: false, Stacked: true},
			{Name: "UNCONN", Label: "Close", Diff: false, Stacked: true},
			{Name: "CLOSE-WAIT", Label: "Close Wait", Diff: false, Stacked: true},
			{Name: "LAST-ACK", Label: "Last Ack", Diff: false, Stacked: true},
			{Name: "LISTEN", Label: "Listen", Diff: false, Stacked: true},
			{Name: "CLOSING", Label: "Closing", Diff: false, Stacked: true},
			{Name: "UNKNOWN", Label: "Unknown", Diff: false, Stacked: true},
		},
	}

	cmd := exec.Command("ss", "-na")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := parseSs(out, p); err != nil {
		return err
	}
	return cmd.Wait()
}

// parsing metrics from ss
func parseSs(r io.Reader, p *map[string]interface{}) error {
	var (
		status      = 0
		first       = true
		overstuffed = false
	)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Fields(line)
		if len(record) < 5 {
			continue
		}
		if first {
			first = false
			if record[0] == "State" {
				// for RHEL6
				status = 0
			} else if record[1] == "State" {
				// for RHEL7
				status = 1
			} else if record[0] == "NetidState" {
				status = 1
				overstuffed = true
			}
			continue
		}
		key := record[status]
		if overstuffed && len(record[0]) > 5 {
			key = record[0][5:]
		}
		v, _ := (*p)[key].(float64)
		(*p)[key] = v + 1
	}

	return nil
}

// collect /proc/vmstat
func collectProcVmstat(path string, p *map[string]interface{}) error {
	graphdef["linux.swap"] = mp.Graphs{
		Label: "Linux Swap Usage",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pswpin", Label: "Swap In", Diff: true},
			{Name: "pswpout", Label: "Swap Out", Diff: true},
		},
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return parseProcVmstat(file, p)
}

// parsing metrics from /proc/vmstat
func parseProcVmstat(r io.Reader, p *map[string]interface{}) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Fields(line)
		if len(record) != 2 {
			continue
		}
		var errParse error
		(*p)[record[0]], errParse = atof(record[1])
		if errParse != nil {
			return errParse
		}
	}

	return nil
}

// atof
func atof(str string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(str, " "), 64)
}

// Do the plugin
func Do() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-linux"
	app.Version = version
	app.Usage = "Get metrics from Linux."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = flags
	app.Action = doMain

	app.Run(os.Args)
}

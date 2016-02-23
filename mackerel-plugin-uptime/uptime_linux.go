package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func fetchMetrics() (map[string]interface{}, error) {
	contentbytes, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch uptime metrics: %s", err)
	}
	return calcMetrics(string(contentbytes))
}

func calcMetrics(str string) (map[string]interface{}, error) {
	cols := strings.Split(str, " ")
	f, err := strconv.ParseFloat(cols[0], 64)
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch uptime metrics: %s", err)
	}
	return map[string]interface{}{"seconds": f}, nil
}

package mpsolr

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var (
	logger       = logging.GetLogger("metrics.plugin.solr")
	coreStatKeys = []string{"numDocs", "deletedDocs", "indexHeapUsageBytes", "version", "segmentCount", "sizeInBytes"}
	handlerPaths = []string{"/update/json", "/select", "/update/json/docs", "/get", "/update/csv", "/replication", "/update", "/dataimport"}
	// Solr5 ... "5minRateReqsPerSecond", "15minRateReqsPerSecond"
	// Solr6 ... "5minRateRequestsPerSecond", "15minRateRequestsPerSecond"
	handlerStatKeys = []string{"requests", "errors", "timeouts", "avgRequestsPerSecond", "5minRateReqsPerSecond", "5minRateRequestsPerSecond",
		"15minRateReqsPerSecond", "15minRateRequestsPerSecond", "avgTimePerRequest", "medianRequestTime",
		"75thPcRequestTime", "95thPcRequestTime", "99thPcRequestTime", "999thPcRequestTime"}
	cacheTypes    = []string{"filterCache", "perSegFilter", "queryResultCache", "documentCache", "fieldValueCache"}
	cacheStatKeys = []string{"lookups", "hits", "hitratio", "inserts", "evictions", "size", "warmupTime"}
)

// SolrPlugin mackerel plugin for Solr
type SolrPlugin struct {
	Protocol string
	Host     string
	Port     string
	BaseURL  string
	Cores    []string
	Prefix   string
	Stats    map[string](map[string]float64)
	Tempfile string
}

func (s *SolrPlugin) getStats(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mackerel-plugin-solr")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Failed to %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var stats map[string]interface{}
	err = dec.Decode(&stats)
	if err != nil {
		logger.Errorf("Failed to %s", err)
		return nil, err
	}
	return stats, nil
}

func (s *SolrPlugin) loadStatsCore(core string, values interface{}) error {
	coreStats := values.(map[string]interface{})["index"].(map[string]interface{})
	for _, k := range coreStatKeys {
		v, ok := coreStats[k].(float64)
		if !ok {
			logger.Errorf("Failed to cast from %s to %s", coreStats[k], "float64")
			return errors.New("type assersion error")
		}
		s.Stats[core][k] = v
	}
	return nil
}

func (s *SolrPlugin) setStatsMbean(core string, stats map[string]interface{}, allowKeys []string) {
	for _, values := range stats["solr-mbeans"].([]interface{}) {
		switch values.(type) {
		case string:
			continue
		default:
			for key, value := range values.(map[string]interface{}) {
				for k, v := range value.(map[string]interface{}) {
					if k != "stats" {
						continue
					}
					if v == nil {
						continue
					}
					v2 := v.(map[string]interface{})
					for _, allowKey := range allowKeys {
						if v2[allowKey] == nil {
							continue
						}
						s.Stats[core][allowKey+"_"+escapeSlash(key)] = v2[allowKey].(float64)
					}
				}
			}
		}
		break // if QUERYHANDLER and QUERY or UPDATEHANDLER and UPDATE
	}
}

func (s *SolrPlugin) loadStatsMbeanHandler(core string, cat string) error {
	uri := s.BaseURL + "/" + core + "/admin/mbeans?stats=true&wt=json&cat=" + cat
	for _, path := range handlerPaths {
		uri += fmt.Sprintf("&key=%s", url.QueryEscape(path))
	}
	stats, err := s.getStats(uri)
	if err != nil {
		return err
	}
	s.setStatsMbean(core, stats, handlerStatKeys)
	return nil
}

func (s *SolrPlugin) loadStatsMbeanCache(core string) error {
	stats, err := s.getStats(s.BaseURL + "/" + core + "/admin/mbeans?stats=true&wt=json&cat=CACHE&key=filterCache&key=perSegFilter&key=queryResultCache&key=documentCache&key=fieldValueCache")
	if err != nil {
		return err
	}
	s.setStatsMbean(core, stats, cacheStatKeys)
	return nil
}

func (s *SolrPlugin) loadStats() error {
	s.Stats = map[string](map[string]float64){}

	stats, err := s.getStats(s.BaseURL + "/admin/cores?wt=json")
	if err != nil {
		return err
	}
	s.Cores = []string{}
	for core, values := range stats["status"].(map[string]interface{}) {
		s.Cores = append(s.Cores, core)
		s.Stats[core] = map[string]float64{}
		err := s.loadStatsCore(core, values)
		if err != nil {
			return err
		}
		err = s.loadStatsMbeanHandler(core, "QUERYHANDLER")
		if err != nil {
			return err
		}
		err = s.loadStatsMbeanHandler(core, "UPDATEHANDLER")
		if err != nil {
			return err
		}
		err = s.loadStatsMbeanHandler(core, "REPLICATION")
		if err != nil {
			return err
		}
		err = s.loadStatsMbeanCache(core)
		if err != nil {
			return err
		}
	}
	return nil
}

func escapeSlash(slashIncludedString string) (str string) {
	str = strings.Replace(slashIncludedString, "/", "", -1)
	return
}

// FetchMetrics interface for mackerelplugin
func (s SolrPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	for core, stats := range s.Stats {
		for k, v := range stats {
			stat[core+"_"+k] = v
		}
	}
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (s SolrPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := make(map[string]mp.Graphs)

	for _, core := range s.Cores {
		graphdef[fmt.Sprintf("%s.%s.docsCount", s.Prefix, core)] = mp.Graphs{
			Label: fmt.Sprintf("%s DocsCount", core),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: core + "_numDocs", Label: "NumDocs"},
				{Name: core + "_deletedDocs", Label: "DeletedDocs"},
			},
		}

		for _, key := range []string{"indexHeapUsageBytes", "segmentCount", "sizeInBytes"} {
			metricLabel := strings.Title(key)
			graphdef[fmt.Sprintf("%s.%s.%s", s.Prefix, core, key)] = mp.Graphs{
				Label: fmt.Sprintf("%s %s", core, metricLabel),
				Unit:  "integer",
				Metrics: []mp.Metrics{
					{Name: core + "_" + key, Label: metricLabel},
				},
			}
		}

		for _, key := range handlerStatKeys {
			var metrics []mp.Metrics
			for _, path := range handlerPaths {
				path = escapeSlash(path)
				metricLabel := strings.Title(path)
				diff := false
				if key == "requests" {
					diff = true
				}
				metrics = append(metrics,
					mp.Metrics{Name: fmt.Sprintf("%s_%s_%s", core, key, path), Label: metricLabel, Diff: diff},
				)
			}
			unit := "float"
			if key == "requests" || key == "errors" || key == "timeouts" {
				unit = "integer"
			}
			graphLabel := fmt.Sprintf("%s %s", core, strings.Title(key))
			graphdef[fmt.Sprintf("%s.%s.%s", s.Prefix, core, key)] = mp.Graphs{
				Label:   graphLabel,
				Unit:    unit,
				Metrics: metrics,
			}
		}

		for _, key := range cacheStatKeys {
			var metrics []mp.Metrics
			for _, cacheType := range cacheTypes {
				metricLabel := strings.Title(cacheType)
				metrics = append(metrics,
					mp.Metrics{Name: fmt.Sprintf("%s_%s_%s", core, key, cacheType), Label: metricLabel},
				)
			}
			unit := "integer"
			if key == "hitratio" {
				unit = "float"
			}
			graphLabel := fmt.Sprintf("%s %s", core, strings.Title(key))
			graphdef[fmt.Sprintf("%s.%s.%s", s.Prefix, core, key)] = mp.Graphs{
				Label:   graphLabel,
				Unit:    unit,
				Metrics: metrics,
			}
		}
	}
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8983", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	solr := SolrPlugin{
		Protocol: "http",
		Host:     *optHost,
		Port:     *optPort,
		Prefix:   "solr",
	}

	solr.BaseURL = fmt.Sprintf("%s://%s:%s/solr", solr.Protocol, solr.Host, solr.Port)
	solr.loadStats()

	helper := mp.NewMackerelPlugin(solr)
	helper.Tempfile = *optTempfile

	helper.Run()
}

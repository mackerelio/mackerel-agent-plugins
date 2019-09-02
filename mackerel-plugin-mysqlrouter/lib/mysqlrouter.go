package mpmysqlrouter

import (
	"flag"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/rluisr/mysqlrouter-go"
)

var (
	url      = os.Getenv("MYSQLROUTER_URL")
	user     = os.Getenv("MYSQLROUTER_USER")
	pass     = os.Getenv("MYSQLROUTER_PASS")
	graphdef map[string]mp.Graphs
)

// MRPlugin is the prefix of struct of graph
type MRPlugin struct {
	Prefix   string
	MRClient *mysqlrouter.Client
}

// MetricKeyPrefix is set prefix of metrics
func (mr MRPlugin) MetricKeyPrefix() string {
	if mr.Prefix == "" {
		mr.Prefix = "mysqlrouter"
	}
	return mr.Prefix
}

// GraphDefinition return struct of graph
func (mr MRPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// FetchMetrics set metrics from MySQL Router
func (mr MRPlugin) FetchMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)

	routes, err := mr.MRClient.GetAllRoutes()
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		routeStatus, err := mr.MRClient.GetRouteStatus(route.Name)
		if err != nil {
			return nil, err
		}
		metrics[route.Name+".active_connections"] = float64(routeStatus.ActiveConnections)
		metrics[route.Name+".total_connections"] = float64(routeStatus.TotalConnections)
		metrics[route.Name+".blocked_host"] = float64(routeStatus.BlockedHosts)

		routeHealth, err := mr.MRClient.GetRouteHealth(route.Name)
		if err != nil {
			return nil, err
		}
		if routeHealth.IsAlive {
			metrics[route.Name+".health"] = float64(1)
		} else {
			metrics[route.Name+".health"] = float64(0)
		}

		// Todo
		//routeConnections, err := mr.MRClient.GetRouteConnections(route.Name)
	}

	return metrics, nil
}

// Prepare define struct of metrics
func (mr MRPlugin) Prepare() {
	g := map[string]mp.Graphs{}

	routes, err := mr.MRClient.GetAllRoutes()
	if err != nil {
		panic(err)
	}

	var metrics []mp.Metrics
	for _, route := range routes {
		metrics = append(metrics, mp.Metrics{Name: route.Name + ".active_connections", Label: "Active connection", Diff: false})
		metrics = append(metrics, mp.Metrics{Name: route.Name + ".total_connections", Label: "Total connection", Diff: false})
		metrics = append(metrics, mp.Metrics{Name: route.Name + ".blocked_host", Label: "Blocked Host", Diff: false})
		metrics = append(metrics, mp.Metrics{Name: route.Name + ".health", Label: "Health of the route. 0 is No health", Diff: false})
	}

	g[mr.Prefix+"route"] = mp.Graphs{
		Label:   "Connections of route",
		Unit:    "integer",
		Metrics: metrics,
	}

	graphdef = g
}

func Do() {
	if url == "" || user == "" || pass == "" {
		panic("The environment missing.\n" +
			"MYSQLROUTER_URL, MYSQLROUTER_USER and MYSQLROUTER_PASS is required.")
	}

	mrr, err := mysqlrouter.New(url, user, pass)
	if err != nil {
		panic(err)
	}

	prefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	flag.Parse()

	mr := MRPlugin{
		Prefix:   *prefix,
		MRClient: mrr,
	}

	mr.Prepare()

	plugin := mp.NewMackerelPlugin(mr)
	plugin.Run()
}

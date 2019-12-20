package mpopenldap

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var (
	logger   = logging.GetLogger("metrics.plugin.openldap")
	cnRegexp = regexp.MustCompile("^cn=([^,]+),")
)

// OpenLDAPPlugin plugin for OpenLDAP
type OpenLDAPPlugin struct {
	Prefix             string
	UseTLS             bool
	InsecureSkipVerify bool
	TargetHost         string
	BindDn             string
	ReplBase           string
	ReplMasterHost     string
	ReplMasterUseTLS   bool
	ReplMasterBind     string
	ReplMasterPass     string
	ReplLocalBind      string
	ReplLocalPass      string
	BindPasswd         string
	Tempfile           string
	l                  *ldap.Conn
}

func transformKeyName(key string) string {
	results := cnRegexp.FindStringSubmatch(key)
	if len(results) < 2 {
		return ""
	}
	return strings.Replace(strings.ToLower(strings.TrimSpace(results[1])), " ", "_", -1)
}

func getStats(sr *ldap.SearchResult, prefix string) map[string]float64 {
	stat := make(map[string]float64)
	for _, entry := range sr.Entries {
		for _, attr := range entry.Attributes {
			key := prefix + transformKeyName(entry.DN) + "_" + attr.Name
			value, err := strconv.ParseFloat(entry.GetAttributeValue(attr.Name), 64)
			if err != nil {
				logger.Warningf("Failed to parse value. %s", err)
			}
			stat[key] = value
		}
	}
	return stat
}

func fetchOpenldapMetrics(l *ldap.Conn, base, prefix string, attrs []string) (map[string]float64, error) {
	searchRequest := ldap.NewSearchRequest(base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&(objectClass=*))", attrs, nil)
	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Errorf("Failed to ldap search. %s.", err)
		return nil, err
	}
	stat := getStats(sr, prefix)
	return stat, nil
}

func mergeStat(dst, src map[string]float64) {
	for k, v := range src {
		dst[k] = v
	}
}

func latestCSN(sr *ldap.SearchResult) (time.Time, error) {
	var res time.Time
	if len(sr.Entries) == 0 {
		return res, errors.New("not found CSN")
	}
	entry := sr.Entries[0]
	if len(entry.Attributes) == 0 {
		return res, errors.New("not found CSN")
	}
	attr := entry.Attributes[0]
	vs := entry.GetAttributeValues(attr.Name)
	csns := make([]time.Time, len(vs))
	for i, v := range vs {
		t, err := time.Parse("20060102150405.999999Z", v[0:strings.Index(v, "#")])
		if err != nil {
			return res, err
		}
		csns[i] = t
	}
	sort.Slice(csns, func(i, j int) bool {
		return csns[i].After(csns[j])
	})
	res = csns[0]
	return res, nil
}

func getLatestCSN(host, base, bind, passwd string, useTLS, insecureSkipVerify bool) (time.Time, error) {
	var l *ldap.Conn
	var err error
	var res time.Time
	if useTLS {
		l, err = ldap.DialTLS("tcp", host, &tls.Config{InsecureSkipVerify: insecureSkipVerify})
	} else {
		l, err = ldap.Dial("tcp", host)
	}
	err = l.Bind(bind, passwd)
	if err != nil {
		logger.Errorf("Failed to Bind %s, err: %s", bind, err)
		return res, err
	}
	searchRequest := ldap.NewSearchRequest(base, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(&(objectClass=*))", []string{"ContextCSN"}, nil)
	sr, err := l.Search(searchRequest)
	l.Close()
	if err != nil {
		logger.Errorf("Failed to ldap search. %s.", err)
		return res, err
	}
	return latestCSN(sr)

}

// FetchMetrics interface for mackerelplugin
func (m OpenLDAPPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]float64)
	if m.ReplBase != "" {
		masterTime, err := getLatestCSN(m.ReplMasterHost, m.ReplBase, m.ReplMasterBind, m.ReplMasterPass, m.ReplMasterUseTLS, m.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		localTime, err := getLatestCSN(m.TargetHost, m.ReplBase, m.ReplLocalBind, m.ReplLocalPass, m.UseTLS, m.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		stat["replication_delay"] = masterTime.Sub(localTime).Seconds()
	}

	ldapOpes, err := fetchOpenldapMetrics(m.l, "cn=Operations,cn=Monitor", "", []string{"monitorOpInitiated", "monitorOpCompleted"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapOpes)
	ldapWaiters, err := fetchOpenldapMetrics(m.l, "cn=Waiters,cn=Monitor", "waiters_", []string{"monitorCounter"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapWaiters)

	ldapMaxThreads, err := fetchOpenldapMetrics(m.l, "cn=Max,cn=Threads,cn=Monitor", "threads_", []string{"monitoredInfo"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapMaxThreads)
	ldapOpenThreads, err := fetchOpenldapMetrics(m.l, "cn=Open,cn=Threads,cn=Monitor", "threads_", []string{"monitoredInfo"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapOpenThreads)
	ldapActiveThreads, err := fetchOpenldapMetrics(m.l, "cn=Active,cn=Threads,cn=Monitor", "threads_", []string{"monitoredInfo"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapActiveThreads)

	ldapStatistics, err := fetchOpenldapMetrics(m.l, "cn=Statistics,cn=Monitor", "statistics_", []string{"monitorCounter"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapStatistics)
	ldapTotalConns, err := fetchOpenldapMetrics(m.l, "cn=Total,cn=Connections,cn=Monitor", "connections_", []string{"monitorCounter"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapTotalConns)
	ldapCurrentConns, err := fetchOpenldapMetrics(m.l, "cn=Current,cn=Connections,cn=Monitor", "connections_", []string{"monitorCounter"})
	if err != nil {
		return nil, err
	}
	mergeStat(stat, ldapCurrentConns)

	result := make(map[string]interface{})
	for k, v := range stat {
		result[k] = v
	}
	return result, nil
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m OpenLDAPPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "openldap"
	}
	return m.Prefix
}

// GraphDefinition interface for mackerelplugin
func (m OpenLDAPPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(m.Prefix)
	graphs := map[string]mp.Graphs{
		"operations_initiated": {
			Label: (labelPrefix + " operations Initiated"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "operations_monitorOpInitiated", Label: "All", Diff: true},
				{Name: "bind_monitorOpInitiated", Label: "Bind", Diff: true},
				{Name: "unbind_monitorOpInitiated", Label: "Unbind", Diff: true},
				{Name: "search_monitorOpInitiated", Label: "Search", Diff: true},
				{Name: "compare_monitorOpInitiated", Label: "Compare", Diff: true},
				{Name: "modify_monitorOpInitiated", Label: "Modify", Diff: true},
				{Name: "modrdn_monitorOpInitiated", Label: "Modrdn", Diff: true},
				{Name: "add_monitorOpInitiated", Label: "Add", Diff: true},
				{Name: "delete_monitorOpInitiated", Label: "Delete", Diff: true},
				{Name: "abandon_monitorOpInitiated", Label: "Abandon", Diff: true},
				{Name: "extended_monitorOpInitiated", Label: "Extended", Diff: true},
			},
		},
		"operations_Completed": {
			Label: (labelPrefix + " operations Completed"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "operations_monitorOpCompleted", Label: "All", Diff: true},
				{Name: "bind_monitorOpCompleted", Label: "Bind", Diff: true},
				{Name: "unbind_monitorOpCompleted", Label: "Unbind", Diff: true},
				{Name: "search_monitorOpCompleted", Label: "Search", Diff: true},
				{Name: "compare_monitorOpCompleted", Label: "Compare", Diff: true},
				{Name: "modify_monitorOpCompleted", Label: "Modify", Diff: true},
				{Name: "modrdn_monitorOpCompleted", Label: "Modrdn", Diff: true},
				{Name: "add_monitorOpCompleted", Label: "Add", Diff: true},
				{Name: "delete_monitorOpCompleted", Label: "Delete", Diff: true},
				{Name: "abandon_monitorOpCompleted", Label: "Abandon", Diff: true},
				{Name: "extended_monitorOpCompleted", Label: "Extended", Diff: true},
			},
		},
		"waiters": {
			Label: (labelPrefix + " waiters"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "waiters_read_monitorCounter", Label: "read", Diff: false},
				{Name: "waiters_write_monitorCounter", Label: "write", Diff: false},
			},
		},
		"threads": {
			Label: (labelPrefix + " threads"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "threads_max_monitoredInfo", Label: "max", Diff: false},
				{Name: "threads_open_monitoredInfo", Label: "open", Diff: false},
				{Name: "threads_active_monitoredInfo", Label: "active", Diff: false},
			},
		},
		"statistics_bytes": {
			Label: (labelPrefix + " statistics bytes"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "statistics_bytes_monitorCounter", Label: "bytes", Diff: true},
			},
		},
		"statistics_pdu": {
			Label: (labelPrefix + " statistics pdu"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "statistics_pdu_monitorCounter", Label: "pdu", Diff: true},
			},
		},
		"statistics_entries": {
			Label: (labelPrefix + " statistics entries"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "statistics_entries_monitorCounter", Label: "entries", Diff: true},
			},
		},
		"statistics_referrals": {
			Label: (labelPrefix + " statistics referrals"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "statistics_referrals_monitorCounter", Label: "referrals", Diff: true},
			},
		},
		"connections": {
			Label: (labelPrefix + " connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connections_total_monitorCounter", Label: "total connection", Diff: true},
				{Name: "connections_current_monitorCounter", Label: "current connection", Diff: false},
			},
		},
	}
	if m.ReplBase != "" {
		graphs["replications"] = mp.Graphs{
			Label: (labelPrefix + " replication delay"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "replication_delay", Label: "replication delay sec", Diff: false},
			},
		}
	}
	return graphs
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "389", "Port")
	optTLS := flag.Bool("tls", false, "TLS(ldaps)")
	optInsecureSkipVerify := flag.Bool("insecureSkipVerify", false, "TLS accepts any certificate.")
	optReplBase := flag.String("replBase", "", "replication base dn")
	optReplMasterBind := flag.String("replMasterBind", "", "replication master bind dn")
	optReplMasterHost := flag.String("replMasterHost", "", "replication master hostname")
	optReplMasterTLS := flag.Bool("replMasterTLS", false, "replication master TLS(ldaps)")
	optReplMasterPort := flag.String("replMasterPort", "389", "replication master port")
	optReplMasterPass := flag.String("replMasterPW", os.Getenv("OPENLDAP_REPL_MASTER_PASSWORD"), "replication master bind password")
	optReplLocalBind := flag.String("replLocalBind", "", "replicationlocalmaster bind dn")
	optReplLocalPass := flag.String("replLocalPW", os.Getenv("OPENLDAP_REPL_LOCAL_PASSWORD"), "replication local bind password")
	optBindDn := flag.String("bind", "", "bind dn")
	optBindPasswd := flag.String("pw", os.Getenv("OPENLDAP_PASSWORD"), "bind password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "openldap", "Metric key prefix")
	flag.Parse()

	var m OpenLDAPPlugin
	m.TargetHost = fmt.Sprintf("%s:%s", *optHost, *optPort)
	m.ReplMasterHost = fmt.Sprintf("%s:%s", *optReplMasterHost, *optReplMasterPort)
	m.UseTLS = *optTLS
	m.InsecureSkipVerify = *optInsecureSkipVerify
	m.ReplBase = *optReplBase
	m.ReplMasterUseTLS = *optReplMasterTLS
	m.ReplMasterBind = *optReplMasterBind
	m.ReplMasterPass = *optReplMasterPass
	m.ReplLocalBind = *optReplLocalBind
	m.ReplLocalPass = *optReplLocalPass
	if m.InsecureSkipVerify {
		m.UseTLS = true
	}
	m.BindDn = *optBindDn
	m.BindPasswd = *optBindPasswd
	m.Prefix = *optPrefix

	if *optBindDn == "" {
		logger.Errorf("bind is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	var err error
	if m.UseTLS {
		m.l, err = ldap.DialTLS("tcp", m.TargetHost, &tls.Config{InsecureSkipVerify: m.InsecureSkipVerify})
	} else {
		m.l, err = ldap.Dial("tcp", m.TargetHost)
	}
	if err != nil {
		logger.Errorf("Failed to Dial %s, err: %s", m.TargetHost, err)
		os.Exit(1)
	}
	err = m.l.Bind(m.BindDn, m.BindPasswd)
	if err != nil {
		logger.Errorf("Failed to Bind %s, err: %s", m.BindDn, err)
		os.Exit(1)
	}
	helper := mp.NewMackerelPlugin(m)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-openldap-%s", *optHost))
	}

	helper.Run()
}

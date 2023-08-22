//go:build linux

package mpdocker

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

func TestNormalizeMetricName(t *testing.T) {
	testSets := [][]string{
		{"foo/bar", "foo_bar"},
		{"foo:bar", "foo_bar"},
	}

	for _, testSet := range testSets {
		if normalizeMetricName(testSet[0]) != testSet[1] {
			t.Errorf("normalizeMetricName: '%s' should be normalized to '%s', but '%s'", testSet[0], testSet[1], normalizeMetricName(testSet[0]))
		}
	}
}

func TestGraphDefinition(t *testing.T) {
	var docker DockerPlugin

	graphdef := docker.GraphDefinition()
	if len(graphdef) != 6 {
		t.Errorf("GraphDefinition: %d should be 6", len(graphdef))
	}
}

func TestGenerateName(t *testing.T) {
	stub := docker.APIContainers{
		ID:      "bab2b03c736de41ecba6470eba736c5109436f706eedca4f3e0d93d6530eccd4",
		Image:   "tutum/mongodb",
		Command: "/run.sh",
		Created: 1456995574,
		Status:  "Up 4 days",
		Ports: []docker.APIPort{
			{PrivatePort: 28017, Type: "tcp"},
			{PrivatePort: 27017, Type: "tcp"},
		},
		Names:  []string{"/my-mongodb"},
		Labels: map[string]string{"foo": "bar"},
	}
	/* {"Id":"5b963f266d609d2b02aee8f57d664e04d35aa8c23afcbc6bb73bc4a5b2e7c44d",
	   "Image":"memcached",
	   "Command":"/entrypoint.sh memcached",
	   "Created":1456994862,
	   "Status":"Up 4 days",
	   "Ports":[{"PrivatePort":11211,
	   "Type":"tcp"}],
	   "Names":["/my-memcache"]}]`
	*/
	var docker DockerPlugin
	docker.NameFormat = "name_id"
	if docker.generateName(stub) != "my-mongodb_bab2b0" {
		t.Errorf("generateName(name): %s should be 'my-mongodb_bab2b0'", docker.generateName(stub))
	}
	docker.NameFormat = "name"
	if docker.generateName(stub) != "my-mongodb" {
		t.Errorf("generateName(name): %s should be 'my-mongodb'", docker.generateName(stub))
	}
	docker.NameFormat = "id"
	if docker.generateName(stub) != "bab2b03c736de41ecba6470eba736c5109436f706eedca4f3e0d93d6530eccd4" {
		t.Errorf("generateName(name): %s should be 'bab2b03c736de41ecba6470eba736c5109436f706eedca4f3e0d93d6530eccd4'", docker.generateName(stub))
	}
	docker.NameFormat = "image"
	if docker.generateName(stub) != "tutum/mongodb" {
		t.Errorf("generateName(name): %s should be 'tutum/mongodb'", docker.generateName(stub))
	}
	docker.NameFormat = "image_id"
	if docker.generateName(stub) != "tutum/mongodb_bab2b0" {
		t.Errorf("generateName(name): %s should be 'tutum/mongodb_bab2b0'", docker.generateName(stub))
	}
	docker.NameFormat = "image_name"
	if docker.generateName(stub) != "tutum/mongodb_my-mongodb" {
		t.Errorf("generateName(name): %s should be 'tutum/mongodb_my-mongodb'", docker.generateName(stub))
	}
	docker.NameFormat = "label"
	docker.Label = "foo"
	if docker.generateName(stub) != "bar" {
		t.Errorf("generateName(name): %s should be 'bar'", docker.generateName(stub))
	}

}

func TestAddCPUPercentageStats(t *testing.T) {
	stats := map[string]interface{}{
		"docker._internal.cpuacct.containerA.user":       uint64(3000),
		"docker._internal.cpuacct.containerA.system":     uint64(2000),
		"docker._internal.cpuacct.containerA.host":       uint64(100000),
		"docker._internal.cpuacct.containerA.onlineCPUs": int(2),
		"docker._internal.cpuacct.containerB.host":       uint64(100000),
		"docker._internal.cpuacct.containerB.user":       uint64(3500),
		"docker._internal.cpuacct.containerC.user":       uint64(3300),
		"docker._internal.cpuacct.containerC.system":     uint64(2300),
		"docker._internal.cpuacct.containerD.host":       uint64(100000),
		"docker._internal.cpuacct.containerD.user":       uint64(3000),
		"docker._internal.cpuacct.containerD.system":     uint64(2000),
		"docker._internal.cpuacct.containerF.user":       uint64(3000), // it has been reset
		"docker._internal.cpuacct.containerF.system":     uint64(1000), // it has been reset
		"docker._internal.cpuacct.containerF.host":       uint64(100000100000),
		"docker._internal.cpuacct.containerF.onlineCPUs": int(2),
	}
	oldStats := map[string]interface{}{
		"docker._internal.cpuacct.containerA.host":   float64(90000),
		"docker._internal.cpuacct.containerA.user":   float64(1000),
		"docker._internal.cpuacct.containerA.system": float64(1500),
		"docker._internal.cpuacct.containerB.host":   float64(90000),
		"docker._internal.cpuacct.containerB.user":   float64(3000),
		"docker._internal.cpuacct.containerC.user":   float64(3000),
		"docker._internal.cpuacct.containerC.system": float64(2000),
		"docker._internal.cpuacct.containerE.host":   float64(100000),
		"docker._internal.cpuacct.containerE.user":   float64(3000),
		"docker._internal.cpuacct.containerE.system": float64(2000),
		"docker._internal.cpuacct.containerF.user":   float64(40000000000),
		"docker._internal.cpuacct.containerF.system": float64(20000000000),
		"docker._internal.cpuacct.containerF.host":   float64(100000000000),
	}
	addCPUPercentageStats(&stats, oldStats)

	if stat, ok := stats["docker.cpuacct_percentage.containerA.user"]; !ok {
		t.Errorf("docker.cpuacct_percentage.containerA.user should be calculated")
	} else if stat != float64(40.0) {
		t.Errorf("docker.cpuacct_percentage.containerA.user should be %f, but %f", stat, float64(40.0))
	}

	if _, ok := stats["docker.cpuacct_percentage.containerC.user"]; ok {
		t.Errorf("docker.cpuacct_percentage.containerC.user should not be calculated")
	}

	if _, ok := stats["docker.cpuacct_percentage.containerB.user"]; ok {
		t.Errorf("docker.cpuacct_percentage.containerB.user should not be calculated")
	}

	if _, ok := stats["docker.cpuacct_percentage.containerD.user"]; ok {
		t.Errorf("docker.cpuacct_percentage.containerD.user should not be calculated")
	}

	if _, ok := stats["docker.cpuacct_percentage.containerE.user"]; ok {
		t.Errorf("docker.cpuacct_percentage.containerE.user should not be calculated")
	}

	if stat, ok := stats["docker.cpuacct_percentage.containerF.user"]; !ok {
		t.Errorf("docker.cpuacct_percentage.containerF.user should be calculated")
	} else if stat != float64(6.0) {
		t.Errorf("docker.cpuacct_percentage.containerF.user should be %f, but %f", float64(6.0), stat)
	}
}

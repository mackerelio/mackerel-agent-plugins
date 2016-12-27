package mpdocker

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
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
	if len(graphdef) != 5 {
		t.Errorf("GetTempfilename: %d should be 5", len(graphdef))
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

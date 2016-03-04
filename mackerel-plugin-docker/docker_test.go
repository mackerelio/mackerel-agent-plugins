package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeMetricName(t *testing.T) {
	testSets := [][]string{
		[]string{"foo/bar", "foo_bar"},
		[]string{"foo:bar", "foo_bar"},
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

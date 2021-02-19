package mpawsses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare_Region(t *testing.T) {
	var p SESPlugin
	p.Region = "us-west-2"
	p.prepare()
	assert.Equal(t, "us-west-2", *p.Svc.Config.Region, "Specified region is used")
}

func TestPrepare_Endpoint(t *testing.T) {
	var p SESPlugin
	p.Endpoint = "https://email.us-west-2.amazonaws.com"
	p.prepare()
	assert.Equal(t, "us-west-2", *p.Svc.Config.Region, "Convert from Endpoint to Specified Region")
}

func TestPrepare_InvalidEndpoint(t *testing.T) {
	var p SESPlugin
	p.Endpoint = "https://email.us-west-2.foobar.com"
	actual := p.prepare()
	expected := "--endpoint is invalid"

	if actual.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func TestPrepare_InvalidEndpoint2(t *testing.T) {
	var p SESPlugin
	p.Endpoint = "https://email.us-west-2.foobar.bazqux.com"
	actual := p.prepare()
	expected := "--endpoint is invalid"

	if actual.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

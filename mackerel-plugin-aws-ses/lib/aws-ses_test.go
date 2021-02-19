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

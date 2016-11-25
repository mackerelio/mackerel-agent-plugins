package mpmailq

import (
	"os"
	"strings"
	"testing"
)

func TestGetNthLine(t *testing.T) {
	{
		s := `0000`
		l, _ := getNthLine(strings.NewReader(s), 0)
		if l != "0000" {
			t.Errorf("0-th line is expected to be 0000")
		}
	}

	{
		s := `
1111`
		l, _ := getNthLine(strings.NewReader(s), 0)
		if l != "" {
			t.Errorf("0-th line is expected to be empty")
		}
	}

	{
		s := `0000
1111
`
		l, _ := getNthLine(strings.NewReader(s), 0)
		if l != "0000" {
			t.Errorf("0-th line is expected to be 0000")
		}
	}

	{
		s := ``
		l, _ := getNthLine(strings.NewReader(s), 0)
		if l != "" {
			t.Errorf("0-th line is expected to be empty")
		}
	}

	{
		s := `0000`
		l, _ := getNthLine(strings.NewReader(s), 1)
		if l != "" {
			t.Errorf("1-st line is expected to be empty")
		}
	}

	{
		s := `
1111`
		l, _ := getNthLine(strings.NewReader(s), 1)
		if l != "1111" {
			t.Errorf("1-st line is expected to be 1111")
		}
	}

	{
		s := `0000
1111
`
		l, _ := getNthLine(strings.NewReader(s), 1)
		if l != "1111" {
			t.Errorf("1-st line is expected to be 1111")
		}
	}

	{
		s := ``
		l, _ := getNthLine(strings.NewReader(s), 1)
		if l != "" {
			t.Errorf("1-st line is expected to be empty")
		}
	}

	{
		s := `0000`
		l, _ := getNthLine(strings.NewReader(s), -1)
		if l != "0000" {
			t.Errorf("-1-st line is expected to be 0000")
		}
	}

	{
		s := `
1111`
		l, _ := getNthLine(strings.NewReader(s), -1)
		if l != "1111" {
			t.Errorf("1-st line is expected to be 1111")
		}
	}

	{
		s := `0000
1111
`
		l, _ := getNthLine(strings.NewReader(s), -1)
		if l != "1111" {
			t.Errorf("1-st line is expected to be 1111")
		}
	}

	{
		s := ``
		l, _ := getNthLine(strings.NewReader(s), -1)
		if l != "" {
			t.Errorf("-1-st line is expected to be empty")
		}
	}

	{
		s := `0000`
		l, _ := getNthLine(strings.NewReader(s), -2)
		if l != "" {
			t.Errorf("-2-nd line is expected to be empty")
		}
	}

	{
		s := `
1111`
		l, _ := getNthLine(strings.NewReader(s), -2)
		if l != "" {
			t.Errorf("-2-nd line is expected to be empty")
		}
	}

	{
		s := `0000
1111
`
		l, _ := getNthLine(strings.NewReader(s), -2)
		if l != "0000" {
			t.Errorf("-2-nd line is expected to be 0000")
		}
	}

	{
		s := ``
		l, _ := getNthLine(strings.NewReader(s), -2)
		if l != "" {
			t.Errorf("-2-nd line is expected to be empty")
		}
	}
}

func TestParseMailqPostfix(t *testing.T) {
	mailq := mailqFormats["postfix"]

	{
		output := `-Queue ID- --Size-- ----Arrival Time---- -Sender/Recipient-------
DD0C740001C      274 Thu Mar  3 23:52:37  foobar@example.com
          (connect to mail.invalid[192.0.2.100]:25: Connection timed out)
                                         nyao@mail.invalid

-- 15 Kbytes in 42 Requests.
`

		count, err := mailq.parse(strings.NewReader(output))
		if err != nil {
			t.Errorf("Error in parseMailq: %s", err.Error())
		}
		if count != 42 {
			t.Errorf("Incorrect parse result %d", count)
		}
	}

	{
		output := `Mail queue is empty
`

		count, err := mailq.parse(strings.NewReader(output))
		if err != nil {
			t.Errorf("Error in parseMailq: %s", err.Error())
		}
		if count != 0 {
			t.Errorf("Incorrect parse result %d", count)
		}
	}
}

func TestParseMailqQmain(t *testing.T) {
	mailq := mailqFormats["qmail"]

	{
		output := `messages in queue: 42
messages in queue but not yet preprocessed: 3
`

		count, err := mailq.parse(strings.NewReader(output))
		if err != nil {
			t.Errorf("Error in parseMailq: %s", err.Error())
		}
		if count != 42 {
			t.Errorf("Incorrect parse result %d", count)
		}
	}
}

func TestParseMailqExim(t *testing.T) {
	mailq := mailqFormats["exim"]

	{
		output := `42
`

		count, err := mailq.parse(strings.NewReader(output))
		if err != nil {
			t.Errorf("Error in parseMailq: %s", err.Error())
		}
		if count != 42 {
			t.Errorf("Incorrect parse result %d", count)
		}
	}
}

func TestGraphDefinition(t *testing.T) {
	plugin := plugin{
		mailq:       mailqFormats["postfix"],
		keyPrefix:   "mailq",
		labelPrefix: "Mailq",
	}

	{
		graphs := plugin.GraphDefinition()
		graphMailq, ok := graphs["mailq"]
		if !ok {
			t.Errorf("No graph definition for mailq")
		}

		if graphMailq.Unit != "integer" {
			t.Errorf("Mailq is expected to be an integral graph")
		}

		if graphMailq.Label == "" {
			t.Errorf("Mailq is expected to have a label")
		}

		if len(graphMailq.Metrics) != 1 {
			t.Errorf("Mailq is expected to have one definition of metrics")
		}

		if graphMailq.Metrics[0].Name != "count" {
			t.Errorf("Mailq is expected to have count metric")
		}

		if graphMailq.Metrics[0].Type != "uint64" {
			t.Errorf("Mailq is expected to have type uint64")
		}
	}
}

func TestFetchMetricsPostfix(t *testing.T) {
	plugin := plugin{
		mailq:       mailqFormats["postfix"],
		keyPrefix:   "mailq",
		labelPrefix: "Mailq",
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "./fixtures:/bin:/usr/bin")
	defer os.Setenv("PATH", origPath)

	{
		os.Setenv("TEST_MAILQ_COUNT", "42")
		defer os.Unsetenv("TEST_MAILQ_COUNT")

		metrics, err := plugin.FetchMetrics()
		if err != nil {
			t.Errorf("Error %s", err.Error())
		}
		if metrics["count"].(uint64) != 42 {
			t.Errorf("Incorrect value: %d", metrics["count"].(uint64))
		}
	}
}

func TestFetchMetricsQmail(t *testing.T) {
	cwd, _ := os.Getwd()

	plugin := plugin{
		mailq:       mailqFormats["qmail"],
		path:        cwd + "/fixtures/qmail/qmail-qstat",
		keyPrefix:   "mailq",
		labelPrefix: "Mailq",
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "/bin:/usr/bin")
	defer os.Setenv("PATH", origPath)

	{
		os.Setenv("TEST_MAILQ_COUNT", "42")
		defer os.Unsetenv("TEST_MAILQ_COUNT")

		metrics, err := plugin.FetchMetrics()
		if err != nil {
			t.Errorf("Error %s", err.Error())
		}
		if metrics["count"].(uint64) != 42 {
			t.Errorf("Incorrect value: %d", metrics["count"].(uint64))
		}
	}
}

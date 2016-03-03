package main

import (
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

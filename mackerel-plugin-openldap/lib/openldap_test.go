package mpopenldap

import (
	"testing"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
)

func TestTransformKeyName(t *testing.T) {
	s := transformKeyName("cn=Operations,cn=Monitor")
	want := "operations"
	if s != want {
		t.Errorf("transformKeyName() = %q, want %q", s, want)
	}
}

func TestGetStats(t *testing.T) {
	sr := &ldap.SearchResult{
		Entries: []*ldap.Entry{
			&ldap.Entry{
				DN: "cn=Extended,cn=Operations,cn=Monitor",
				Attributes: []*ldap.EntryAttribute{
					&ldap.EntryAttribute{
						Name:       "monitorOpInitiated",
						Values:     []string{"0"},
						ByteValues: [][]uint8{[]uint8{0x30}},
					},
					&ldap.EntryAttribute{
						Name:       "monitorOpCompleted",
						Values:     []string{"0"},
						ByteValues: [][]uint8{[]uint8{0x30}},
					},
				},
			},
		},
	}
	stats := getStats(sr, "prefix_")
	key := "prefix_extended_monitorOpCompleted"
	v, ok := stats[key]
	if !ok {
		t.Errorf("not found key:%s", key)
	}
	want := float64(0)
	if v != want {
		t.Errorf("stats[%s] = %f, want %f", key, v, want)
	}
	key = "prefix_extended_monitorOpInitiated"
	v, ok = stats[key]
	if !ok {
		t.Errorf("not found key:%s", key)
	}
	want = float64(0)
	if v != want {
		t.Errorf("stats[%s] = %f, want %f", key, v, want)
	}
}

func TestLatestCSN(t *testing.T) {
	sr := &ldap.SearchResult{
		Entries: []*ldap.Entry{
			&ldap.Entry{
				DN: "dc=example,dc=net",
				Attributes: []*ldap.EntryAttribute{
					&ldap.EntryAttribute{
						Name: "contextCSN",
						Values: []string{
							"20161205033538.343893Z#000000#001#000000",
							"20140128082749.962641Z#000000#002#000000",
							"20170713094701.963361Z#000000#05b#000000",
						},
						ByteValues: [][]uint8{
							[]uint8{
								0x32, 0x30, 0x31, 0x36, 0x31, 0x32, 0x30, 0x35, 0x30, 0x33, 0x33, 0x35, 0x33, 0x38, 0x2e, 0x33,
								0x34, 0x33, 0x38, 0x39, 0x33, 0x5a, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x23, 0x30, 0x30,
								0x31, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30,
							},
							[]uint8{
								0x32, 0x30, 0x31, 0x34, 0x30, 0x31, 0x32, 0x38, 0x30, 0x38, 0x32, 0x37, 0x34, 0x39, 0x2e, 0x39,
								0x36, 0x32, 0x36, 0x34, 0x31, 0x5a, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x23, 0x30, 0x30,
								0x32, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30,
							},
							[]uint8{
								0x32, 0x30, 0x31, 0x37, 0x30, 0x37, 0x31, 0x33, 0x30, 0x39, 0x34, 0x37, 0x30, 0x31, 0x2e, 0x39,
								0x36, 0x33, 0x33, 0x36, 0x31, 0x5a, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x23, 0x30, 0x35,
								0x62, 0x23, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30,
							},
						},
					},
				},
			},
		},
	}
	tm, _ := latestCSN(sr)
	want, _ := time.Parse("2006-01-02T15:04:05.999999Z07:00", "2017-07-13T09:47:01.963361Z")
	if !tm.Equal(want) {
		t.Errorf("latestCSN = %q, want %q", tm, want)
	}
}

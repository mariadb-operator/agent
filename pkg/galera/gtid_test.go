package galera_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mariadb-operator/agent/pkg/galera"
)

func TestGTIDUnmarshal(t *testing.T) {
	var gtid galera.GTID
	err := gtid.UnmarshalText([]byte("0-1-2"))
	if err != nil {
		t.Fatal(err)
	}

	want := galera.GTID{
		DomainID:       0,
		ServerID:       1,
		SequenceNumber: 2,
	}
	if diff := cmp.Diff(gtid, want); diff != "" {
		t.Errorf("GTID.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
}

func TestGTIDUnmarshalTextInvalidValues(t *testing.T) {
	tt := []struct {
		name string
		in   string
	}{
		{
			name: "invalid domain id",
			in:   "a-1-2",
		},
		{
			name: "invalid server id",
			in:   "0-a-2",
		},
		{
			name: "invalid seqno",
			in:   "0-1-a",
		},
		{
			name: "missing part",
			in:   "0-1",
		},
		{
			name: "too many parts",
			in:   "0-1-2-3",
		},
		{
			name: "out of range domain id",
			in:   "4294967296-1-2",
		},
		{
			name: "out of range server id",
			in:   "0-4294967296-2",
		},
		{
			name: "out of range seqno",
			in:   "0-1-18446744073709551616",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var gtid galera.GTID
			err := gtid.UnmarshalText([]byte(tc.in))
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestGTIDMarshalText(t *testing.T) {
	gtid := galera.GTID{
		DomainID:       0,
		ServerID:       1,
		SequenceNumber: 2,
	}
	text, err := gtid.MarshalText()
	if err != nil {
		t.Fatalf("failed to marshal GTID to text: %v", err)
	}

	want := "0-1-2"
	if string(text) != want {
		t.Errorf("GTID.MarshalText() = %q, want %q", text, want)
	}
}

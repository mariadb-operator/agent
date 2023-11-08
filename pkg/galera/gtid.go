package galera

import (
	"fmt"
	"strconv"
	"strings"
)

// GTID represents a MariaDB global transaction identifier.
// See: https://mariadb.com/kb/en/gtid/
type GTID struct {
	DomainID       uint32 `json:"domainId"`
	ServerID       uint32 `json:"serverId"`
	SequenceNumber uint64 `json:"sequenceNumber"`
}

func (gtid *GTID) String() string {
	return fmt.Sprintf("%d-%d-%d", gtid.DomainID, gtid.ServerID, gtid.SequenceNumber)
}

func (gtid *GTID) MarshalText() ([]byte, error) {
	return []byte(gtid.String()), nil
}

func (gtid *GTID) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid gtid: %s", text)
	}
	domainID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid domain id: %v", err)
	}
	serverID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid server id: %v", err)
	}
	seqno, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid seqno: %v", err)
	}
	gtid.DomainID = uint32(domainID)
	gtid.ServerID = uint32(serverID)
	gtid.SequenceNumber = seqno
	return nil
}

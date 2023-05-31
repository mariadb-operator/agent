package galerastate

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type GaleraState struct {
	Version         string `json:"version"`
	UUID            string `json:"uuid"`
	Seqno           int    `json:"seqno"`
	SafeToBootstrap bool   `json:"safeToBootstrap"`
}

func (g *GaleraState) Unmarshal(b []byte) error {
	fileScanner := bufio.NewScanner(bytes.NewReader(b))
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		parts := strings.Split(fileScanner.Text(), ":")
		if len(parts) != 2 {
			return fmt.Errorf("error unmarshalling galera state: invalid '%s'", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "version":
			g.Version = value
		case "uuid":
			g.UUID = value
		case "seqno":
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing seqno: %v", err)
			}
			g.Seqno = i
		case "safe_to_bootstrap":
			b, err := parseBool(value)
			if err != nil {
				return fmt.Errorf("error parsing safe_to_bootstrap: %v", err)
			}
			g.SafeToBootstrap = b
		}
	}
	return nil
}

func parseBool(s string) (bool, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, fmt.Errorf("error parsing integer bool: %v", err)
	}
	if i != 0 && i != 1 {
		return false, fmt.Errorf("invalid integer bool: %d", i)
	}
	return i == 1, nil
}

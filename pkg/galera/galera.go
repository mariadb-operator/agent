package galera

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

var (
	GaleraStateFileName = "grastate.dat"
	BootstrapFileName   = "1-bootstrap.cnf"
	BootstrapFile       = `[galera]
wsrep_new_cluster="ON"`
	RecoveryFileName = "2-recovery.cnf"
	RecoveryFile     = `[galera]
log_error=mariadb.err
wsrep_recover="ON"`
)

type GaleraState struct {
	Version         string `json:"version"`
	UUID            string `json:"uuid"`
	Seqno           int    `json:"seqno"`
	SafeToBootstrap bool   `json:"safeToBootstrap"`
}

func (g *GaleraState) MarshalText() ([]byte, error) {
	type tplOpts struct {
		Version         string
		UUID            string
		Seqno           int
		SafeToBootstrap int
	}
	tpl := createTpl("grastate.dat", `version: {{ .Version }}
uuid: {{ .UUID }}
seqno: {{ .Seqno }}
safe_to_bootstrap: {{ .SafeToBootstrap }}`)
	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, tplOpts{
		Version: g.Version,
		UUID:    g.UUID,
		Seqno:   g.Seqno,
		SafeToBootstrap: func() int {
			if g.SafeToBootstrap {
				return 1
			}
			return 0
		}(),
	})
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %v", err)
	}
	return buf.Bytes(), nil
}

func (g *GaleraState) UnmarshalText(text []byte) error {
	fileScanner := bufio.NewScanner(bytes.NewReader(text))
	fileScanner.Split(bufio.ScanLines)

	var version *string
	var uuid *string
	var seqno *int
	var safeToBootstrap *bool

	for fileScanner.Scan() {
		parts := strings.Split(fileScanner.Text(), ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "version":
			version = &value
		case "uuid":
			uuid = &value
		case "seqno":
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing seqno: %v", err)
			}
			seqno = &i
		case "safe_to_bootstrap":
			b, err := parseBool(value)
			if err != nil {
				return fmt.Errorf("error parsing safe_to_bootstrap: %v", err)
			}
			safeToBootstrap = &b
		}
	}

	if version == nil || uuid == nil || seqno == nil || safeToBootstrap == nil {
		return fmt.Errorf(
			"invalid galera state file: version=%v uuid=%v seqno=%v safeToBootstrap=%v",
			version, uuid, seqno, safeToBootstrap,
		)
	}
	g.Version = *version
	g.UUID = *uuid
	g.Seqno = *seqno
	g.SafeToBootstrap = *safeToBootstrap
	return nil
}

func createTpl(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
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

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
	GaleraStateFile = "grastate.dat"
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

	for fileScanner.Scan() {
		line := fileScanner.Text()
		parts := strings.Split(fileScanner.Text(), ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid galera state line: '%s'", line)
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

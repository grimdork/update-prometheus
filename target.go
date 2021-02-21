package main

import (
	"net/url"
	"sort"
	"strings"

	"github.com/Urethramancer/signor/files"
)

// Target represents one group of targets to monitor as part of one job.
type Target struct {
	// Job this belongs to.
	Job string
	// Scheme to use for every target
	Scheme string
	// URLs to check.
	URLs []string
	// Metrics path
	Metrics string
	// Params are the parameters after the path.
	Params []Parameter
	// ScrapeInterval is the time between checks.
	ScrapeInterval string
	// EvaluationInterval is the time between rule checks.
	EvaluationInterval string
	// ScrapeTimeout is the time to wait for a scrape to respond before declaring a failure.
	ScrapeTimeout string
}

// Parameter for targets.
type Parameter struct {
	Key   string
	Value string
}

// LoadTarget from INI file.
func LoadTarget(fn string) (Target, error) {
	t := Target{}
	ini, err := files.LoadINI(fn)
	if err != nil {
		return t, err
	}

	sec, ok := ini.Sections["target"]
	if !ok {
		return t, ErrNoTarget
	}
	t.Job = sec.GetString("job", "")
	if t.Job == "" {
		return t, ErrNoJob
	}

	t.ScrapeInterval = sec.GetString("scrape_interval", "30s")
	t.EvaluationInterval = sec.GetString("evaluation_interval", "30s")
	t.ScrapeTimeout = sec.GetString("scrape_timeout", "10s")

	urlstring := sec.GetString("urls", "")
	if urlstring == "" {
		return t, ErrNoURLs
	}

	urls := strings.Split(urlstring, " ")
	for _, x := range urls {
		u, err := url.Parse(x)
		if err != nil {
			return t, err
		}

		t.Scheme = u.Scheme
		if t.Scheme == "" {
			t.Scheme = "http"
		}

		t.URLs = append(t.URLs, u.Host)
		t.Metrics = u.Path
		if t.Metrics == "" {
			t.Metrics = "/metrics"
		}

		for k, v := range u.Query() {
			p := Parameter{k, v[0]}
			t.Params = append(t.Params, p)
		}
	}

	sort.Slice(t.Params, func(i, j int) bool {
		return t.Params[i].Key < t.Params[j].Key
	})
	return t, nil
}

// YAML format string.
func (t Target) YAML(indent int) string {
	b := strings.Builder{}
	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString("metrics_path: '")
	b.WriteString(t.Metrics)
	b.WriteString("'\n")

	if t.Scheme != "http" {
		b.WriteString(strings.Repeat(" ", indent))
		b.WriteString("scheme: ")
		b.WriteString(t.Scheme)
		b.WriteString("\n")
	}

	if len(t.Params) > 0 {
		b.WriteString(strings.Repeat(" ", indent))
		b.WriteString("params:\n")
	}
	for _, p := range t.Params {
		b.WriteString(strings.Repeat(" ", indent+2))
		b.WriteString(p.Key)
		b.WriteString(": ['")
		b.WriteString(p.Value)
		b.WriteString("'")
		b.WriteString("']\n")
	}

	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString("static_configs:\n")
	b.WriteString(strings.Repeat(" ", indent+2))
	b.WriteString("- targets: [")
	count := 0
	for _, x := range t.URLs {
		if count > 0 {
			b.WriteString(", ")
		}
		b.WriteByte('"')
		b.WriteString(x)
		b.WriteByte('"')
		count++
	}
	b.WriteString("]\n")
	return b.String()
}

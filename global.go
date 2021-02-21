package main

import (
	"strings"

	"github.com/Urethramancer/signor/files"
)

// Global settings.
type Global struct {
	// ScrapeInterval is the time between checks.
	ScrapeInterval string
	// EvaluationInterval is the time between rule checks.
	EvaluationInterval string
	// ScrapeTimeout is the time to wait for a scrape to respond before declaring a failure.
	ScrapeTimeout string
	// ExternalLabels are optional.
	ExternalLabels []Parameter
	// Rules point to alert rules.
	Rules         []string
	AlertManagers []string
}

// LoadGlobal from INI file.
func LoadGlobal(fn string) (Global, error) {
	g := Global{}
	ini, err := files.LoadINI(fn)
	if err != nil {
		return g, err
	}

	sec, ok := ini.Sections["global"]
	if !ok {
		return g, ErrNoTarget
	}

	labels := sec.GetString("external_labels", "")
	if labels != "" {
		list := strings.Split(labels, ",")
		for _, x := range list {
			a := strings.Split(x, ":")
			if len(a) < 2 {
				return g, ErrKV
			}
			p := Parameter{a[0], a[1]}
			g.ExternalLabels = append(g.ExternalLabels, p)
		}
	}

	g.ScrapeInterval = sec.GetString("scrape_interval", "30s")
	g.EvaluationInterval = sec.GetString("scrape_interval", "30s")
	g.ScrapeTimeout = sec.GetString("scrape_interval", "10s")

	rules := sec.GetString("rules", "")
	if rules != "" {
		g.Rules = strings.Split(rules, " ")
	}

	am := sec.GetString("alertmanagers", "")
	if am != "" {
		list := strings.Split(am, " ")
		for _, x := range list {
			g.AlertManagers = append(g.AlertManagers, x)
		}
	}

	return g, nil
}

// YAML format string.
func (g Global) YAML(indent int) string {
	b := strings.Builder{}
	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString("global:\n")

	b.WriteString(strings.Repeat(" ", indent+2))
	b.WriteString("scrape_interval: ")
	b.WriteString(g.ScrapeInterval)
	b.WriteString("\n")

	b.WriteString(strings.Repeat(" ", indent+2))
	b.WriteString("evaluation_interval: ")
	b.WriteString(g.EvaluationInterval)
	b.WriteString("\n")

	b.WriteString(strings.Repeat(" ", indent+2))
	b.WriteString("scrape_timeout: ")
	b.WriteString(g.ScrapeInterval)
	b.WriteString("\n\n")

	if len(g.ExternalLabels) > 0 {
		b.WriteString(strings.Repeat(" ", indent+2))
		b.WriteString("external_labels:\n")
		for _, l := range g.ExternalLabels {
			b.WriteString(strings.Repeat(" ", indent+4))
			b.WriteString(l.Key)
			b.WriteString(": '")
			b.WriteString(l.Value)
			b.WriteString("'\n")
		}
		b.WriteString("\n")
	}

	if len(g.AlertManagers) > 0 {
		b.WriteString(strings.Repeat(" ", indent))
		b.WriteString("alerting:\n")
		b.WriteString(strings.Repeat(" ", indent+2))
		b.WriteString("alertmanagers:\n")
		b.WriteString(strings.Repeat(" ", indent+2))
		b.WriteString("- static_configs:\n")
		b.WriteString(strings.Repeat(" ", indent+4))
		b.WriteString("- targets: [")
		count := 0
		for _, am := range g.AlertManagers {
			if count > 0 {
				b.WriteString(", ")
			}
			b.WriteString("'")
			b.WriteString(am)
			b.WriteString("'")
			count++
		}
		b.WriteString("]\n\n")
	}

	if len(g.Rules) > 0 {
		b.WriteString(strings.Repeat(" ", indent))
		b.WriteString("rule_files:\n")
		for _, r := range g.Rules {
			b.WriteString(strings.Repeat(" ", indent+2))
			b.WriteString("- ")
			b.WriteString(r)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

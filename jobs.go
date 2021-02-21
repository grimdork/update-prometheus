package main

import "strings"

// Jobs list
type Jobs struct {
	Targets map[string][]Target
}

// AddTarget URLs to monitor to a job.
func (jobs *Jobs) AddTarget(t Target) {
	list := jobs.Targets[t.Job]
	list = append(list, t)
	jobs.Targets[t.Job] = list
}

// YAML format string.
func (jobs *Jobs) YAML(indent int) string {
	b := strings.Builder{}
	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString("scrape_configs:\n")
	for j, list := range jobs.Targets {
		b.WriteString(strings.Repeat(" ", indent+2))
		b.WriteString("- job_name: ")
		b.WriteString(j)
		b.WriteString("\n")
		for _, t := range list {
			b.WriteString(t.YAML(indent + 4))
			b.WriteString("\n")
		}
	}
	return b.String()
}

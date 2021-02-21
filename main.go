package main

import (
	"os"
	"path/filepath"

	"github.com/Urethramancer/signor/cfmt"
	"github.com/Urethramancer/signor/opt"
)

var o struct {
	opt.DefaultHelp
	Config string `short:"c" long:"config" placeholder:"PATH" help:"Path with configuration INI files to assemble from." default:"prometheus.d"`
	Out    string `short:"o" long:"output" placeholder:"FILE" help:"Output Prometheus configuration file to create." default:"prometheus.yml"`
}

func main() {
	a := opt.Parse(&o)
	if o.Help {
		a.Usage()
		return
	}

	dir, err := os.ReadDir(o.Config)
	if err != nil {
		cfmt.Printf("%red Error reading '%s':%reset %s", o.Config, err.Error())
		os.Exit(2)
	}

	var g Global
	jobs := Jobs{make(map[string][]Target)}
	for _, de := range dir {
		if !de.IsDir() && filepath.Ext(de.Name()) == ".ini" && de.Name() != "global.ini" {
			fn := filepath.Join(o.Config, de.Name())
			t, err := LoadTarget(fn)
			if err != nil {
				cfmt.Printf("%red Error:%reset Couldn't load '%s': %s", de.Name(), err.Error())
				os.Exit(2)
			}

			jobs.AddTarget(t)
		}

		if de.Name() == "global.ini" {
			fn := filepath.Join(o.Config, de.Name())
			g, err = LoadGlobal(fn)
			if err != nil {
				cfmt.Printf("%red Error:%reset Couldn't load '%s'; %s", de.Name(), err.Error())
				os.Exit(2)
			}
		}
	}
	out := g.YAML(0) + jobs.YAML(0)
	err = os.WriteFile(o.Out, []byte(out), 0600)
	if err != nil {
		cfmt.Printf("%red Error saving output:%reset %s", err.Error())
		os.Exit(2)
	}
}

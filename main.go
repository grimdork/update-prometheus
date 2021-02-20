package main

import (
	"fmt"
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
		pr("Error reading '%s': %s", o.Config, err.Error())
		os.Exit(2)
	}

	jobs := Jobs{make(map[string][]Target)}
	for _, de := range dir {
		if !de.IsDir() && filepath.Ext(de.Name()) == ".ini" {
			fn := filepath.Join(o.Config, de.Name())
			t, err := LoadTarget(fn)
			if err != nil {
				cfmt.Printf("%red Error:%reset  Couldn't load '%s': %s", de.Name(), err.Error())
				os.Exit(2)
			}

			jobs.AddTarget(t)
		}
	}
	pr("%s", jobs.YAML(4))
}

func pr(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Urethramancer/signor/cfmt"
	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/opt"
)

var o struct {
	opt.DefaultHelp
	Rec     string `short:"r" long:"receivers" placeholder:"PATH" help:"Path for e-mail users to be defined." default:"receivers.d"`
	Targets string `short:"t" long:"targets" placeholder:"PATH" help:"Path monitoring target configurations." default:"targets.d"`
	Out     string `short:"o" long:"output" placeholder:"FILE" help:"Output Prometheus configuration file to create." default:"prometheus.yml"`
}

func main() {
	a := opt.Parse(&o)
	if o.Help {
		a.Usage()
		return
	}

	dir, err := os.ReadDir(o.Targets)
	if err != nil {
		pr("Error reading '%s': %s", o.Targets, err.Error())
		os.Exit(2)
	}

	for _, de := range dir {
		if !de.IsDir() && filepath.Ext(de.Name()) == ".ini" {
			fn := filepath.Join(o.Targets, de.Name())
			s, err := loadTarget(fn)
			if err != nil {
				cfmt.Printf("%red Error:%reset  Couldn't load '%s': %s", de.Name(), err.Error())
				os.Exit(2)
			}

			pr("%s", s)
		}
	}
}

func pr(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func loadTarget(fn string) (string, error) {
	ini, err := files.LoadINI(fn)
	if err != nil {
		return "", err
	}

	pr("%s", ini.Sections["main"])
	return "", nil
}

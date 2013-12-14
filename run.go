package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

type runner interface {
	run(args ...string)
	flags(string) *flag.FlagSet
}

type shell struct {
	user   *string
	pass   *string
	Handle func(user, pass string, commands ...string)
}

func (cmd *shell) flags(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	cmd.user = fs.String("user", "vagrant", "user to run as")
	cmd.pass = fs.String("pass", "vagrant", "user's password")
	return fs
}

func (cmd *shell) run(commands ...string) {
	if len(commands) == 0 {
		fmt.Fprint(os.Stderr, "specify a command to run\n")
		fail()
	}

	if cmd.Handle != nil {
		cmd.Handle(*cmd.user, *cmd.pass)
	}
}

func Run(runners ...runner) {
	specs := make(map[string]*runspec, len(runners))
	for _, r := range runners {
		v := reflect.ValueOf(r).Elem()
		name := v.Type().Name()
		specs[name] = &runspec{r, r.flags(name)}
	}

	usage := flag.Usage
	flag.Usage = func() {
		usage()
		for name, spec := range specs {
			fmt.Fprintf(os.Stderr, "\n%s %s [options] [arguments]\n", os.Args[0], name)
			spec.fs.PrintDefaults()
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	flag.Parse()

	if flag.NArg() < 1 {
		fail()
	}

	name := flag.Arg(0)
	if name == "help" {
		help()
	}

	args := flag.Args()[1:]
	if spec, ok := specs[name]; ok {
		spec.fs.Parse(args)
		spec.r.run(spec.fs.Args()...)
	} else {
		fmt.Fprintf(os.Stderr, "%s is not a runner command\n", name)
		fail()
	}
}

type runspec struct {
	r  runner
	fs *flag.FlagSet
}

func help() {
	flag.Usage()
	os.Exit(0)
}

func fail() {
	flag.Usage()
	os.Exit(1)
}

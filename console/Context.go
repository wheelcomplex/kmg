package console

import (
	"flag"
	"io"
)

type Context struct {
	Args    []string //origin args,if you want args after command name and flag->FlagSet().Args()
	flagSet *flag.FlagSet
	Stdin   io.ReadCloser  // in
	Stdout  io.WriteCloser // in
	Stderr  io.WriteCloser // in
	// main execution name normal is args[0] in/out
	// if you set this value,Manager will use it
	ExecutionName string
	// command name normal is args[1]  in/out,
	// if you set this value,Manager will not parse it from args,and leave args[1:] to FlagSet()
	// this value is case insensitive,always lower case
	CommandName string
	exitCode    int //exit code
}

func (c *Context) FlagSet() *flag.FlagSet {
	return c.flagSet
}
func (c *Context) ExitCode() int {
	return c.exitCode
}
func (c *Context) PrintUsage(msg string) {
	c.Stderr.Write([]byte(msg))
	c.flagSet.PrintDefaults()
}

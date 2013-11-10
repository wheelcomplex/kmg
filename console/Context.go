package console

import (
	"flag"
	"io"
)

type Context struct {
	Args          []string //origin args include execute name
	FlagSet       *flag.FlagSet
	Stdin         io.ReadCloser
	Stdout        io.WriteCloser
	Stderr        io.WriteCloser
	ExecutionName string
}

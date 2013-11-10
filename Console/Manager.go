package console

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

type Manager struct {
	Map  map[string]Command
	Name string
}

func NewManager() *Manager {
	return &Manager{Map: make(map[string]Command)}
}
func (manager *Manager) Add(command Command) error {
	name := getNameFromCommand(command)
	_, ok := manager.Map[name]
	if ok {
		return errors.Sprintf("command %s already exist", name)
	}
	manager.Map[name] = command
	return nil
}
func (manager *Manager) MustAdd(command Command) {
	err := manager.Add(command)
	if err != nil {
		panic(err)
	}
}

func (manager *Manager) ExecuteGlobal() {
	manager.Execute(&Context{Args: os.Args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	},
	)
}
func (manager *Manager) Execute(context *Context) {
	args := context.Args
	manager.Name = filepath.Base(args[0])
	context.ExecutionName = manager.Name
	if len(args) < 2 {
		manager.usage(context.Stderr)
	}
	commandName := args[1]
	command, ok := manager.Map[commandName]
	if !ok {
		fmt.Fprintf(context.Stderr, "unknown subcommand %q\n", commandName)
		os.Exit(2)
	}
	flagSet := flag.NewFlagSet(manager.Name+" "+commandName, flag.ExitOnError)
	if command, ok := command.(FlagSetAwareInterface); ok {
		command.ConfigFlagSet(flagSet)
	}
	flagSet.Parse(args[2:])
	err := command.Execute(context)
	if err != nil {
		fmt.Fprintf(context.Stderr, "error happend: %s\n", err.Error())
		os.Exit(2)
	}
	os.Exit(0)
}

func (manager *Manager) usage(w io.Writer) {
	Commands := make([]NameConfig, 0, len(manager.Map))
	for _, v := range manager.Map {
		Commands = append(Commands, NameConfig{Name: getNameFromCommand(v), Short: getShortFromCommand(v)})
	}
	tmpl(w, usageTemplate,
		&struct {
			ExecuteName string
			Commands    []NameConfig
		}{
			ExecuteName: manager.Name,
			Commands:    Commands,
		},
	)
	os.Exit(2)
}

var usageTemplate = `Usage:

	{{.ExecuteName }} command [arguments]

The commands are:
{{range .Commands}}
    {{.Name | printf "%-11s"}} {{.Short }}{{end}}
`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func getNameFromCommand(command Command) string {
	if nameCommand, ok := command.(NameAwareInterface); ok {
		return nameCommand.GetNameConfig().Name
	}

	typeName := getTypeShortName(reflect.TypeOf(command))
	return strings.TrimSuffix(typeName, "Command")
}
func getShortFromCommand(command Command) string {
	if nameCommand, ok := command.(NameAwareInterface); ok {
		return nameCommand.GetNameConfig().Short
	}
	typeName := getTypeShortName(reflect.TypeOf(command))
	return strings.TrimSuffix(typeName, "Command")
}

func getTypeShortName(t reflect.Type) string {
	typeName := t.Name()
	if typeName != "" {
		return typeName
	}
	switch t.Kind() {
	case reflect.Ptr:
		return getTypeShortName(t.Elem())
	}
	return ""
}

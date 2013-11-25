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
	name := strings.ToLower(getNameFromCommand(command))
	_, ok := manager.Map[name]
	if ok {
		return errors.Sprintf("command %s already exist,command name is case insensitive.", name)
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
	context := &Context{Args: os.Args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	manager.Execute(context)
	os.Exit(context.exitCode)
}
func (manager *Manager) ExecuteGlobalByName(name string) {
	context := &Context{Args: os.Args,
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		CommandName: name,
	}
	manager.Execute(context)
	os.Exit(context.exitCode)
}
func (manager *Manager) Execute(context *Context) {
	args := context.Args
	if context.ExecutionName == "" {
		context.ExecutionName = filepath.Base(args[0])
	}
	var command Command
	//arguments
	var actualArgs []string
	if context.CommandName == "" {
		if len(args) < 2 {
			manager.usage(context.Stderr)
		}
		context.CommandName = strings.ToLower(args[1])
		actualArgs = args[2:]

	} else {
		context.CommandName = strings.ToLower(context.CommandName)
		actualArgs = args[1:]
	}

	//exec command
	var ok bool
	command, ok = manager.Map[context.CommandName]
	if !ok {
		fmt.Fprintf(context.Stderr, "unknown subcommand %q\n", context.CommandName)
		context.exitCode = 2
		return
	}
	context.flagSet = flag.NewFlagSet(context.ExecutionName+" "+context.CommandName, flag.ExitOnError)
	if command, ok := command.(FlagSetAwareInterface); ok {
		command.ConfigFlagSet(context.flagSet)
	}
	context.flagSet.Parse(actualArgs)

	err := command.Execute(context)
	if err != nil {
		fmt.Fprintln(context.Stderr, err.Error())
		context.exitCode = 1
		return
	}
	return
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

The commands are: (command name is case insensitive)
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

package ajkApi

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/dependencyInjection"
	"net"
	"net/http"
)

//start a golang http api server
type GoHttpApiServerCommand struct {
	Container *dependencyInjection.Container
	http      string
	https     string
	randPort  bool
}

func (command *GoHttpApiServerCommand) SetContainer(Container *dependencyInjection.Container) {
	command.Container = Container
}
func (command *GoHttpApiServerCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoHttpApiServer", Short: `start a golang http api server `}
}
func (command *GoHttpApiServerCommand) ConfigFlagSet(f *flag.FlagSet) {
	f.StringVar(&command.http, "http", ":18080", "listen http port of the server")
	f.StringVar(&command.https, "https", "", "listen https port of the server")
	f.BoolVar(&command.randPort, "randPort", false, "if can not listen on default port ,will listen on random port")
}

func (command *GoHttpApiServerCommand) Execute(context *console.Context) error {
	c := command.Container
	handler, err := c.GetByType((*JsonHttpHandler)(nil))
	if err != nil {
		return err
	}
	http.Handle("/api", handler.(http.Handler))
	l, err := command.listen()
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, "Listen on %s\n", l.Addr().String())
	return http.Serve(l, nil)
}

//first try addr,if err happened try random addrss.
func (command *GoHttpApiServerCommand) listen() (l net.Listener, err error) {
	l, err = net.Listen("tcp", *command.http)
	if err == nil {
		return
	}
	if *command.randPort {
		l, err = net.Listen("tcp", ":0")
		return
	}
	return
}

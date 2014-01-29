package kernel

import (
	"fmt"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Kernel struct {
	Bundles   []*Bundle
	Container *dependencyInjection.Container
	// app env "dev" or "test" or "prod",default "dev"
	// This value priority is
	// 1.Kernel.Env
	// 2.Parameter.Env
	Env string
	// set this value to tell kernel where some path is,or it will guess from work dir.
	Context *kmgContext.Context
	// default load from Context.ConfigPath/parameters.yml
	Parameters map[string]interface{}
}

func NewKernel() *Kernel {
	return &Kernel{}
}

func (kernel *Kernel) Boot() (err error) {
	builder := dependencyInjection.NewContainerBuilder()
	err = kernel.handleConfig(builder)
	if err != nil {
		return err
	}
	for _, bundle := range kernel.Bundles {
		bundle.Build(builder)
	}
	kernel.Container, err = builder.Compile()
	return
}

func (kernel *Kernel) MustBoot() *Kernel {
	err := kernel.Boot()
	if err != nil {
		panic(err)
	}
	return kernel
}

func (kernel *Kernel) AddBundle(bundle *Bundle) {
	kernel.Bundles = append(kernel.Bundles, bundle)
}

//parameter should in ./app/config/parameters.yml
func (kernel *Kernel) handleConfig(builder *dependencyInjection.ContainerBuilder) (err error) {
	// already set
	// TODO pass in kernel.Parameters problem
	err = kernel.guessContext()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filepath.Join(kernel.Context.ConfigPath, "parameters.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found!\nyou should put config into ./app/config/parameters.yml")
		}
		return err
	}

	err = kmgYaml.Unmarshal(data, &kernel.Parameters)
	if err != nil {
		return fmt.Errorf("config file parse fail! error:%s", err)
	}

	kernel.Parameters["AppPath"] = kernel.Context.AppPath
	kernel.Parameters["DataPath"] = kernel.Context.DataPath
	kernel.Parameters["TmpPath"] = kernel.Context.TmpPath
	kernel.Parameters["ConfigPath"] = kernel.Context.ConfigPath

	if kernel.Env == "" {
		kernel.Env = "dev"
	}
	kernel.Parameters["Env"] = kernel.Env

	for k, v := range kernel.Parameters {
		k = "Parameter." + k
		err = builder.Set(k, v, "")
		if err != nil {
			return
		}
	}
	return
}

//guess path find a base path with config
func (kernel *Kernel) guessContext() (err error) {
	// already set
	if kernel.Context != nil {
		return nil
	}
	kernel.Context, err = kmgContext.FindFromWd()
	return
}

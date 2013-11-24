package kernel

import (
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/errors"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
)

type Kernel struct {
	Bundles   []*Bundle
	Container *dependencyInjection.Container
	// app env "dev" or "test" or "prod"
	Env string
	// set this value to tell kernel where app path is,or it will guess.
	AppPath    string
	Parameters map[string]string
}

func NewKernel() *Kernel {
	return &Kernel{}
}

func (kernel *Kernel) Boot() (err error) {
	err = kernel.guessParameter()
	if err != nil {
		return err
	}
	kernel.Parameters["AppPath"] = kernel.AppPath
	kernel.Parameters["DataPath"] = filepath.Join(kernel.AppPath, "data")
	kernel.Parameters["TmpPath"] = filepath.Join(kernel.AppPath, "tmp")

	builder := dependencyInjection.NewContainerBuilder()
	builder.Parameters = kernel.Parameters
	for _, bundle := range kernel.Bundles {
		bundle.Build(builder)
	}
	kernel.Container, err = builder.Compile()
	return
}

func (kernel *Kernel) MustBoot() {
	err := kernel.Boot()
	if err != nil {
		panic(err)
	}
}

func (kernel *Kernel) AddBundle(bundle *Bundle) {
	kernel.Bundles = append(kernel.Bundles, bundle)
}

//parameter should in ./app/config/parameters.yml
func (kernel *Kernel) guessParameter() (err error) {
	// already set
	if len(kernel.Parameters) != 0 {
		return nil
	}
	parameter := make(map[string]string)
	err = kernel.guessAppPath()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filepath.Join(kernel.AppPath, "config", "parameters.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("config file not found!\nyou should put config into ./app/config/parameters.yml")
		}
		return err
	}
	err = goyaml.Unmarshal(data, parameter)
	if err != nil {
		return err
	}
	kernel.Parameters = parameter
	return
}

//guess path find a base path with config
func (kernel *Kernel) guessAppPath() (err error) {
	// already set
	if kernel.AppPath != "" {
		return nil
	}
	// find a directory named app from current directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	thisPath := wd
	for {
		fi, err := os.Stat(filepath.Join(thisPath, "app"))
		if err != nil && (!os.IsNotExist(err)) {
			return err
		}
		if err == nil && fi.IsDir() {
			kernel.AppPath = filepath.Join(thisPath, "app")
			return nil
		}
		lastPath := thisPath
		thisPath = filepath.Dir(thisPath)
		if lastPath == thisPath {
			break
		}
	}
	return errors.New("can not guess app path")
}

/*
type KernelInterface interface {
	//register bundles to this kernel
	//all register bundles name should be unique
	RegisterBundles() []BundleInterface
	RegisterContainerConfiguration()
	Boot()
	//use for functional testing
	Shutdown()
	GetBundles() []BundleInterface
	GetBundle(name string) BundleInterface
	// Returns the file path for a given resource.
	// A Resource can be a file or a directory.
	// The resource name must follow the following pattern:
	//   @BundleName/path/to/a/file.something
	// where BundleName is the name of the bundle
	// and the remaining part is the relative path in the bundle.
	LocateResource(name string) string
	// the environment name of kernel : prod,dev,test
	GetEnvironment() string
	// Checks if debug mode is enabled.
	IsDebug() bool
	// It is "app/xxx" in standard way
	GetRootDir() string
	GetCacheDir() string
	GetLogDir() string
	// Gets the current container.
	GetContainer() dependencyInjection.ContainerInterface
}
*/

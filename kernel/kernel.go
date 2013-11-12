package kernel

import "github.com/bronze1man/kmg/dependencyInjection"

type Kernel struct {
	Bundles   []*Bundle
	Container *dependencyInjection.Container
}

func NewKernel() *Kernel {
	return &Kernel{}
}
func (kernel *Kernel) Boot() (err error) {
	builder := dependencyInjection.NewContainerBuilder()
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

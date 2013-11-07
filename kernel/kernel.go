package kernel

import "kmg/dependencyInjection"

type BundleInterface interface {
	dependencyInjection.ContainerAwareInterface
	//can use dic here
	Boot() error
	//use for functional testing
	Shutdown() error
	Build(container dependencyInjection.ContainerBuilder)
	GetParent() string
	GetName() string
	//full file path of this bundle
	//GetFilePath() string
	//
	RegisterCommands()
	GetPackageName() string
}

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

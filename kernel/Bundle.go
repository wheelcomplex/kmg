package kernel

import "github.com/bronze1man/kmg/dependencyInjection"

type Bundle struct {
	Extensions    []dependencyInjection.ExtensionInterface
	CompliePasses []dependencyInjection.CompilePassInterface
	Bootes        []dependencyInjection.BootInterface
}

func (bundle *Bundle) Build(builder *dependencyInjection.ContainerBuilder) {
	for _, extension := range bundle.Extensions {
		builder.AddExtension(extension)
	}
	for _, compilePass := range bundle.CompliePasses {
		builder.AddCompilePass(compilePass)
	}
	for _, boot := range bundle.Bootes {
		builder.AddBoot(boot)
	}
}

func (bundle *Bundle) AddExtension(extenstion dependencyInjection.ExtensionInterface) {
	bundle.Extensions = append(bundle.Extensions, extenstion)
}
func (bundle *Bundle) AddCompliePass(compliePass dependencyInjection.CompilePassInterface) {
	bundle.CompliePasses = append(bundle.CompliePasses, compliePass)
}
func (bundle *Bundle) AddBoot(boot dependencyInjection.BootInterface) {
	bundle.Bootes = append(bundle.Bootes, boot)
}

package kmgBundle

import "github.com/bronze1man/kmg/kernel"

func NewBundle() (bundle *kernel.Bundle) {
	bundle = &kernel.Bundle{}
	bundle.AddExtension(&KmgExtension{})
	return bundle
}

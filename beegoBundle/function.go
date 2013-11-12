package beegoBundle

import "github.com/bronze1man/kmg/kernel"

func NewBundle() (bundle *kernel.Bundle) {
	bundle = &kernel.Bundle{}
	bundle.AddExtension(&BeegoExtension{})
	return bundle
}

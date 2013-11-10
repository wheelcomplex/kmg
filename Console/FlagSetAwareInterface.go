package console

import "flag"

type FlagSetAwareInterface interface {
	ConfigFlagSet(flag *flag.FlagSet)
}

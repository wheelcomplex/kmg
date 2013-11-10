package console

// see FlagSetAwareInterface for flag
// see NameAwareInterface for name
// name come from type name by default
type Command interface {
	Execute(*Context) error
}

package console

type NameConfig struct {
	//name for this command,used for command dispatch
	Name string
	//short description,a phrase to describe the role of command,used for list
	Short string
	//long description,used for help
	Detail string
}

package webTypeAdmin

//path -> key(Key type)
type mapType struct {
	commonType
	key  typeInterface
	elem typeInterface
}

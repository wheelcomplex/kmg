package sessionStore

//all call should thread safe
//not any transaction about single session
type Provider interface {
	//caller should not modify Value
	Get(Id string) (Value []byte, Exist bool, err error)
	//caller should not modify Value
	Set(Id string, Value []byte) (err error)
	//when id not found,it should not report err ,just return silently
	Delete(Id string) (err error)
}

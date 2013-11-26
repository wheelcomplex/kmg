package ajkApi

/*

 ApiFunc must be something like(follow golang rpc protocol)(can not reflect parameter name)
 	1.func Add(apiInput *ReqisterRequest,apiOutput *RegisterResponse)(error){
    2.func Add(apiInput *ReqisterRequest)(error){
    3.func Add()(error)
    4.func Add()
*/
type ApiManagerInterface interface {
	/*
		must pass a session object
		input and output can be
		*struct{xxx}
		map[string]interface
	*/
	RpcCall(session *Session, name string, caller func(*ApiFuncMeta) error) error
}

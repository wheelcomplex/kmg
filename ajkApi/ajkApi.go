package ajkApi

type ApiManagerInterface interface{
	RpcCall(session *Session,name string,input interface{},output interface{})error
}

type Session struct{
	Self *Peer
	Other *Peer
	guid string
}

type Peer struct{

}


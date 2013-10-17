package kmgHttp

import (
	"net/url"
//	"fmt"
)

type Url struct{
	*url.URL
}

func (u *Url)SetQueryByStringMap(input map[string]string){
	value:=&url.Values{}
	for k,v:=range input{
		value.Add(k,v)
	}
	u.URL.RawQuery=value.Encode()
	return
}
func NewUrlByString(raw_url string)(u *Url,err error){
	origin_url,err := url.Parse(raw_url)
	if err!=nil{
		return
	}
	u=&Url{URL:origin_url}
	return
}

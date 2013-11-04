package ajkApi
import (
	"encoding/json"
)
type AjkOutput struct{
	Data interface{} `json:"data"`
	Err string    `json:"err"`
	ErrLocal string `json:"err_local"`// not implement always "",for compatible with php ajk api
}
func OutputToBytes(data interface{},err error)(output []byte){
	outputObject:= OutputToObj(data,err)
	output,errJson:=json.Marshal(outputObject)
	if errJson!=nil{
		return []byte(`{data:null,err:"error happened in jsonEncode.",err_local:""}`)
	}
	return
}

func OutputToString(data interface{},err error)(output string){
	return string(OutputToBytes(data,err))
}

func OutputToObj(data interface{},err error)(output *AjkOutput){
	output= &AjkOutput{
		Data:data,
	}
	if err!=nil{
		output.Err = err.Error()
	}
	return
}

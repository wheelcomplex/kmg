package ajkApi

import (
	"encoding/json"
	"net/http"
	"reflect"
	"kmg/session"
)

type httpInput struct {
	Name string
	Guid string //
	Data json.RawMessage
}
type httpOutput struct {
	Err  string
	Guid string // "" as not set guid to peer
	Data interface{}
}
type JsonHttpHandler struct {
	ApiManager ApiManagerInterface
	SessionManager *session.Manager
}

func (handler *JsonHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	defer req.Body.Close()
	rawInput := &httpInput{}
	err = json.NewDecoder(req.Body).Decode(rawInput)
	if err != nil {
		handler.returnError(w, err)
		return
	}
	var apiOutputValue reflect.Value
	session:=
	err = handler.ApiManager.RpcCall(nil, rawInput.Name, func(meta *ApiFuncMeta) error {
		funcType := meta.Func.Type()
		var inValues []reflect.Value
		serviceValue := meta.AttachObject
		switch funcType.NumIn() {
		case 1:
			inValues = []reflect.Value{serviceValue}
		case 2:
			apiInputValue, err := jsonUnmarshalFromPtrReflectType(funcType.In(1), []byte(rawInput.Data))
			if err != nil {
				return err
			}
			inValues = []reflect.Value{serviceValue, apiInputValue}
		case 3:
			apiInputValue, err := jsonUnmarshalFromPtrReflectType(funcType.In(1), []byte(rawInput.Data))
			if err != nil {
				return err
			}
			apiOutputValue = reflect.New(funcType.In(2).Elem())
			inValues = []reflect.Value{serviceValue, apiInputValue, apiOutputValue}
		default:
			return &ApiFuncArgumentError{Reason: "only accept function input argument num 0,1,2", ApiName: rawInput.Name}
		}
		switch funcType.NumOut() {
		case 0:
		case 1:
			if funcType.Out(0).Kind() != reflect.Interface {
				return &ApiFuncArgumentError{
					Reason:  "only accept function output one argument with error",
					ApiName: rawInput.Name,
				}
			}
		default:
			return &ApiFuncArgumentError{Reason: "only accept function output argument num 0,1", ApiName: rawInput.Name}
		}
		outValues := meta.Func.Call(inValues)

		if len(outValues) == 1 {
			if outValues[0].IsNil() {
				return nil
			}
			err, ok := outValues[0].Interface().(error)
			if ok == false {
				return &ApiFuncArgumentError{
					Reason:  "only accept function output one argument with error",
					ApiName: rawInput.Name,
				}
			}
			return err
		}
		return nil
	})
	if err != nil {
		handler.returnError(w, err)
		return
	}
	handler.returnSuccess(w, apiOutputValue.Interface())
}
func jsonUnmarshalFromPtrReflectType(inputType reflect.Type, data []byte) (reflect.Value, error) {
	var apiInputValue = reflect.New(inputType.Elem())
	apiInput := apiInputValue.Interface()
	err := json.Unmarshal(data, apiInput)
	if err != nil {
		return reflect.Value{}, err
	}
	return apiInputValue, nil
}
func (handler *JsonHttpHandler) returnError(w http.ResponseWriter, err error) {
	rawOutput := &httpOutput{Err: err.Error()}
	err = json.NewEncoder(w).Encode(rawOutput)
	if err != nil {
		//TODO log error
		panic(err)
	}
}
func (handler *JsonHttpHandler) returnSuccess(w http.ResponseWriter, data interface{}) {
	rawOutput := &httpOutput{Data: data}
	err := json.NewEncoder(w).Encode(rawOutput)
	if err != nil {
		//TODO log error
		panic(err)
	}
}

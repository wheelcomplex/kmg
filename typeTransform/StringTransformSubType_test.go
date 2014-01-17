package typeTransform

import "testing"
import (
	"github.com/bronze1man/kmg/test"
)

type StringTranT1 struct {
	T2 StringTranT2
}
type StringTranT2 string

func TestStringTransformSubType(ot *testing.T) {
	t := test.NewTestTools(ot)
	in := &StringTranT1{
		T2: "6",
	}
	err := StringTransformSubType(in, map[string]map[string]string{
		"github.com/bronze1man/kmg/typeTransform.StringTranT2": {
			"6": "Fire",
		},
	})
	t.Equal(err, nil)
	t.Equal(in.T2, StringTranT2("Fire"))
}

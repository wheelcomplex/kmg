package webTypeAdmin

import (
	//"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

type T struct {
	String1 string
	Map1    map[string]string
	Map2    map[string]*string
	Map3    map[string]T2
	Map4    map[string]map[string]string
	Slice1  []string
	Ptr1    *string
	Array1  [5]string
}
type T2 struct {
	A string
	B string
}

func TestPtrType(ot *testing.T) {
	/*
		t := kmgTest.NewTestTools(ot)
		var data **string
		data = new(*string)
		m, err := NewManagerFromPtr(data)
		t.Equal(err, nil)

		err = m.save(Path{"ptr", "ptr"}, "")
		t.Equal(err, nil)
		t.Ok(data != nil)
		t.Ok(*data != nil)
		t.Equal(**data, "")
	*/
}

/*
func TestStringType(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	var data *string
	data = new(string)
	m, err := NewManagerFromPtr(data)
	t.Equal(err, nil)

	err = m.save(Path{"ptr"}, "123")
	t.Equal(err, nil)
	t.Ok(data != nil)
	t.Equal(*data, "123")
}

func TestStructType(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	data := &struct {
		A string
	}{}
	m, err := NewManagerFromPtr(data)
	t.Equal(err, nil)

	err = m.save(Path{"ptr", "A"}, "123")
	t.Equal(err, nil)
	t.Ok(data != nil)
	t.Equal(data.A, "123")
}

func TestType(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	data := &T{}
	m, err := NewManagerFromPtr(data)
	t.Equal(err, nil)

	err = m.save(Path{"ptr", "String1"}, "B")
	t.Equal(err, nil)
	t.Equal(data.String1, "B")

	m.save(Path{"ptr", "Map1", "A"}, "1123")
	_, ok := data.Map1["A"]
	t.Equal(ok, true)
	t.Equal(data.Map1["A"], "1123")

	err = m.save(Path{"ptr", "Map2", "B", "ptr"}, "1")
	t.Equal(err, nil)
	rpString, ok := data.Map2["B"]
	t.Equal(ok, true)
	t.Equal(*rpString, "1")

	err = m.save(Path{"ptr", "Map3", "C", "A"}, "1")
	t.Equal(err, nil)
	t.Equal(data.Map3["C"].A, "1")

	err = m.save(Path{"ptr", "Map4", "D", "F"}, "1234")
	t.Equal(err, nil)
	t.Equal(data.Map4["D"]["F"], "1234")
}
*/

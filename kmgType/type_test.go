package kmgType

import (
	"github.com/bronze1man/kmg/test"
	"testing"
)

type T struct {
	String1 string
	Map1    map[string]string
	Map2    map[string]*string
	Map3    map[string]T2
	Map4    map[string]map[string]string
	Map5    map[string][]string
	Slice1  []string
	Ptr1    *string
	Ptr2    *T2
	Array1  [5]string
}
type T2 struct {
	A string
	B string
}

func TestPtrType(ot *testing.T) {
	t := test.NewTestTools(ot)
	var data **string
	data = new(*string)
	m, err := NewContext(data)
	t.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "ptr"}, "")
	t.Equal(err, nil)
	t.Ok(data != nil)
	t.Ok(*data != nil)
	t.Equal(**data, "")
}

func TestStringType(ot *testing.T) {
	t := test.NewTestTools(ot)
	var data *string
	data = new(string)
	m, err := NewContext(data)
	t.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr"}, "123")
	t.Equal(err, nil)
	t.Ok(data != nil)
	t.Equal(*data, "123")
}

func TestStructType(ot *testing.T) {
	t := test.NewTestTools(ot)
	data := &struct {
		A string
	}{}
	m, err := NewContext(data)
	t.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "A"}, "123")
	t.Equal(err, nil)
	t.Ok(data != nil)
	t.Equal(data.A, "123")
}

func TestType(ot *testing.T) {
	t := test.NewTestTools(ot)
	data := &T{}
	m, err := NewContext(data)
	t.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "String1"}, "B")
	t.Equal(err, nil)
	t.Equal(data.String1, "B")

	m.SaveByPath(Path{"ptr", "Map1", "A"}, "1123")
	_, ok := data.Map1["A"]
	t.Equal(ok, true)
	t.Equal(data.Map1["A"], "1123")

	err = m.SaveByPath(Path{"ptr", "Map1", "A"}, "1124")
	t.Equal(err, nil)
	t.Equal(data.Map1["A"], "1124")

	err = m.DeleteByPath(Path{"ptr", "Map1", "A"})
	t.Equal(err, nil)
	_, ok = data.Map1["A"]
	t.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "1")
	t.Equal(err, nil)
	rpString, ok := data.Map2["B"]
	t.Equal(ok, true)
	t.Equal(*rpString, "1")

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "2")
	t.Equal(err, nil)
	t.Equal(*rpString, "2")

	err = m.DeleteByPath(Path{"ptr", "Map2", "B", "ptr"})
	t.Equal(err, nil)
	_, ok = data.Map2["B"]
	t.Equal(ok, true)
	t.Equal(data.Map2["B"], nil)

	err = m.DeleteByPath(Path{"ptr", "Map2", "B"})
	t.Equal(err, nil)
	_, ok = data.Map2["B"]
	t.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map3", "C", "A"}, "1")
	t.Equal(err, nil)
	t.Equal(data.Map3["C"].A, "1")

	err = m.DeleteByPath(Path{"ptr", "Map3", "C"})
	t.Equal(err, nil)
	t.Ok(data.Map3 != nil)
	_, ok = data.Map3["C"]
	t.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "F"}, "1234")
	t.Equal(err, nil)
	t.Equal(data.Map4["D"]["F"], "1234")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12345")
	t.Equal(err, nil)
	t.Equal(data.Map4["D"]["H"], "12345")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12346")
	t.Equal(err, nil)
	t.Equal(data.Map4["D"]["H"], "12346")

	err = m.DeleteByPath(Path{"ptr", "Map4", "D", "F"})
	t.Equal(err, nil)
	t.Ok(data.Map4["D"] != nil)
	_, ok = data.Map4["D"]["F"]
	t.Equal(ok, false)

	_, ok = data.Map4["D"]["H"]
	t.Equal(ok, true)

	err = m.SaveByPath(Path{"ptr", "Map5", "D", ""}, "1234")
	t.Equal(err, nil)
	t.Equal(len(data.Map5["D"]), 1)
	t.Equal(data.Map5["D"][0], "1234")

	err = m.DeleteByPath(Path{"ptr", "Map5", "D", "0"})
	t.Equal(err, nil)
	t.Equal(len(data.Map5["D"]), 0)

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "1234")
	t.Equal(err, nil)
	t.Equal(len(data.Slice1), 1)
	t.Equal(data.Slice1[0], "1234")

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "12345")
	t.Equal(err, nil)
	t.Equal(data.Slice1[1], "12345")
	t.Equal(len(data.Slice1), 2)

	err = m.DeleteByPath(Path{"ptr", "Slice1", "0"})
	t.Equal(err, nil)
	t.Equal(len(data.Slice1), 1)
	t.Equal(data.Slice1[0], "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr1", "ptr"}, "12345")
	t.Equal(err, nil)
	t.Equal(*data.Ptr1, "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr2", "ptr"}, "")
	t.Equal(err, nil)
	t.Equal(data.Ptr2.A, "")

	err = m.SaveByPath(Path{"ptr", "Array1", "1"}, "12345")
	t.Equal(err, nil)
	t.Equal(data.Array1[1], "12345")
}

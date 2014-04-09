package kmgYaml

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var marshalIntTest = 123

var marshalTests = []struct {
	value interface{}
	data  string
}{
	{
		&struct{}{},
		"{}\n",
	}, {
		map[string]string{"v": "hi"},
		"v: hi\n",
	}, {
		map[string]interface{}{"v": "hi"},
		"v: hi\n",
	}, {
		map[string]string{"v": "true"},
		"v: \"true\"\n",
	}, {
		map[string]string{"v": "false"},
		"v: \"false\"\n",
	}, {
		map[string]interface{}{"v": true},
		"v: true\n",
	}, {
		map[string]interface{}{"v": false},
		"v: false\n",
	}, {
		map[string]interface{}{"v": 10},
		"v: 10\n",
	}, {
		map[string]interface{}{"v": -10},
		"v: -10\n",
	}, {
		map[string]uint{"v": 42},
		"v: 42\n",
	}, {
		map[string]interface{}{"v": int64(4294967296)},
		"v: 4294967296\n",
	}, {
		map[string]int64{"v": int64(4294967296)},
		"v: 4294967296\n",
	}, {
		map[string]uint64{"v": 4294967296},
		"v: 4294967296\n",
	}, {
		map[string]interface{}{"v": "10"},
		"v: \"10\"\n",
	}, {
		map[string]interface{}{"v": 0.1},
		"v: 0.1\n",
	}, {
		map[string]interface{}{"v": float64(0.1)},
		"v: 0.1\n",
	}, {
		map[string]interface{}{"v": -0.1},
		"v: -0.1\n",
	}, {
		map[string]interface{}{"v": math.Inf(+1)},
		"v: .inf\n",
	}, {
		map[string]interface{}{"v": math.Inf(-1)},
		"v: -.inf\n",
	}, {
		map[string]interface{}{"v": math.NaN()},
		"v: .nan\n",
	}, {
		map[string]interface{}{"v": nil},
		"v: null\n",
	}, {
		map[string]interface{}{"v": ""},
		"v: \"\"\n",
	}, {
		map[string][]string{"v": []string{"A", "B"}},
		"v:\n- A\n- B\n",
	}, {
		map[string][]string{"v": []string{"A", "B\nC"}},
		"v:\n- A\n- 'B\n\n  C'\n",
	}, {
		map[string][]interface{}{"v": []interface{}{"A", 1, map[string][]int{"B": []int{2, 3}}}},
		"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
	}, {
		map[string]interface{}{"a": map[interface{}]interface{}{"b": "c"}},
		"a:\n  b: c\n",
	}, {
		time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC),
		"2001-02-03T04:05:06Z\n",
	},

	// Simple values.
	{
		&marshalIntTest,
		"123\n",
	},

	// Structures
	{
		&struct{ Hello string }{"world"},
		"Hello: world\n",
	}, {
		&struct {
			A struct {
				B string
			}
		}{struct{ B string }{"c"}},
		"A:\n  B: c\n",
	}, {
		&struct {
			A *struct {
				B string
			}
		}{&struct{ B string }{"c"}},
		"A:\n  B: c\n",
	}, {
		&struct {
			A *struct {
				B string
			}
		}{},
		"A: null\n",
	}, {
		&struct{ A int }{1},
		"A: 1\n",
	}, {
		&struct{ A []int }{[]int{1, 2}},
		"A:\n- 1\n- 2\n",
	}, {
		&struct {
			B int "a"
		}{1},
		"a: 1\n",
	}, {
		&struct{ A bool }{true},
		"A: true\n",
	},

	// Conditional flag
	{
		&struct {
			A int "a,omitempty"
			B int "b,omitempty"
		}{1, 0},
		"a: 1\n",
	}, {
		&struct {
			A int "a,omitempty"
			B int "b,omitempty"
		}{0, 0},
		"{}\n",
	}, {
		&struct {
			A *struct{ X int } "a,omitempty"
			B int              "b,omitempty"
		}{nil, 0},
		"{}\n",
	},

	// Flow flag
	{
		&struct {
			A []int "a,flow"
		}{[]int{1, 2}},
		"a: [1, 2]\n",
	}, {
		&struct {
			A map[string]string "a,flow"
		}{map[string]string{"b": "c", "d": "e"}},
		"a: {b: c, d: e}\n",
	}, {
		&struct {
			A struct {
				B, D string
			} "a,flow"
		}{struct{ B, D string }{"c", "e"}},
		"a: {B: c, D: e}\n",
	},

	// Unexported field
	{
		&struct {
			u int
			A int
		}{0, 1},
		"A: 1\n",
	},

	// Ignored field
	{
		&struct {
			A int
			B int "-"
		}{1, 2},
		"A: 1\n",
	},

	// Struct inlining
	{
		&struct {
			A int
			C inlineB `yaml:",inline"`
		}{1, inlineB{2, inlineC{3}}},
		"A: 1\nB: 2\nC: 3\n",
	},
}

func (t *S) TestMarshal() {
	for _, item := range marshalTests {
		data, err := Marshal(item.value)
		t.Equal(err, nil)
		t.Equal(string(data), item.data)
	}
}

var marshalErrorTests = []struct {
	value interface{}
	error string
}{
	{
		&struct {
			B       int
			inlineB ",inline"
		}{1, inlineB{2, inlineC{3}}},
		`Duplicated key 'B' in struct struct { B int; kmgYaml.inlineB ",inline" }`,
	},
}

func (t *S) TestMarshalErrors() {
	for _, item := range marshalErrorTests {
		_, err := Marshal(item.value)
		t.Ok(err != nil)
		t.Equal(err.Error(), item.error)
	}
}

var marshalTaggedIfaceTest interface{} = &struct{ A string }{"B"}

var getterTests = []struct {
	data, tag string
	value     interface{}
}{
	{"_:\n  hi: there\n", "", map[interface{}]interface{}{"hi": "there"}},
	{"_:\n- 1\n- A\n", "", []interface{}{1, "A"}},
	{"_: 10\n", "", 10},
	{"_: null\n", "", nil},
	{"_: !foo BAR!\n", "!foo", "BAR!"},
	{"_: !foo 1\n", "!foo", "1"},
	{"_: !foo '\"1\"'\n", "!foo", "\"1\""},
	{"_: !foo 1.1\n", "!foo", 1.1},
	{"_: !foo 1\n", "!foo", 1},
	{"_: !foo 1\n", "!foo", uint(1)},
	{"_: !foo true\n", "!foo", true},
	{"_: !foo\n- A\n- B\n", "!foo", []string{"A", "B"}},
	{"_: !foo\n  A: B\n", "!foo", map[string]string{"A": "B"}},
	{"_: !foo\n  A: B\n", "!foo", &marshalTaggedIfaceTest},
}

func (t *S) TestMarshalTypeCache() {
	var data []byte
	var err error
	func() {
		type T struct{ A int }
		data, err = Marshal(&T{})
		t.Equal(err, nil)
	}()
	func() {
		type T struct{ B int }
		data, err = Marshal(&T{})
		t.Equal(err, nil)
	}()
	t.Equal(string(data), "B: 0\n")
}

type typeWithGetter struct {
	tag   string
	value interface{}
}

func (o typeWithGetter) GetYAML() (tag string, value interface{}) {
	return o.tag, o.value
}

type typeWithGetterField struct {
	Field typeWithGetter "_"
}

func (t *S) TestMashalWithGetter() {
	for _, item := range getterTests {
		obj := &typeWithGetterField{}
		obj.Field.tag = item.tag
		obj.Field.value = item.value
		data, err := Marshal(obj)
		t.Equal(err, nil)
		t.Equal(string(data), string(item.data))
	}
}

func (t *S) TestUnmarshalWholeDocumentWithGetter() {
	obj := &typeWithGetter{}
	obj.tag = ""
	obj.value = map[string]string{"hello": "world!"}
	data, err := Marshal(obj)
	t.Equal(err, nil)
	t.Equal(string(data), "hello: world!\n")
}

func (t *S) TestSortedOutput() {
	order := []interface{}{
		false,
		true,
		1,
		uint(1),
		1.0,
		1.1,
		1.2,
		2,
		uint(2),
		2.0,
		2.1,
		"",
		".1",
		".2",
		".a",
		"1",
		"2",
		"a!10",
		"a/10",
		"a/2",
		"ab/1",
		"a~10",
		"b/01",
		"b/02",
		"b/03",
		"b/1",
		"b/2",
		"b/3",
		"b01",
		"b1",
		"b3",
		"c10.2",
		"c2.10",
		"d1",
		"d12",
		"d12a",
	}
	m := make(map[interface{}]int)
	for _, k := range order {
		m[k] = 1
	}
	data, err := Marshal(m)
	t.Equal(err, nil)
	out := "\n" + string(data)
	last := 0
	for i, k := range order {
		repr := fmt.Sprint(k)
		if s, ok := k.(string); ok {
			if _, err = strconv.ParseFloat(repr, 32); s == "" || err == nil {
				repr = `"` + repr + `"`
			}
		}
		index := strings.Index(out, "\n"+repr+":")
		if index == -1 {
			t.Fatalf("%#v is not in the output: %#v", k, out)
		}
		if index < last {
			t.Fatalf("%#v was generated before %#v: %q", k, order[i-1], out)
		}
		last = index
	}
}

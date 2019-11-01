package conv

import (
	"encoding/json"
	"testing"
)

type A struct {
	A string
	B int64
	C int
	D float64
	E float32
	F *C
	G C
	H *C
	I C
}

type B struct {
	A string
	B int64
	D float64
	E float64
	F *D
	G D
	H D
	I *D
}

type C struct {
	A string
}

type D struct {
	A string
}

var a = &A{
	A: "a",
	B: 1,
	C: 2,
	D: 1.1,
	E: 2.1,
	F: &C{A: "aa"},
	G: C{A: "bb"},
	H: &C{A: "cc"},
	I: C{A: "dd"},
}
var b = &B{}

func TestStructToStruct(t *testing.T) {
	StructToStruct(a, b)
	j, _ := json.Marshal(b)
	want := `{"A":"a","B":1,"D":1.1,"E":0,"F":{"A":"aa"},"G":{"A":"bb"},"H":{"A":"cc"},"I":{"A":"dd"}}`
	if string(j) != want {
		t.Errorf("have: %v \n want: %v", string(j), want)
	} else {
		t.Logf("B: %v", string(j))
	}
}

func BenchmarkStructToStruct(tb *testing.B) {
	StructToStruct(a, b)
}

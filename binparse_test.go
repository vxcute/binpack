package binparse_test

import (
	"bytes"
	"testing"

	"github.com/vxcute/binparse"
)

type User struct {
	Name   string
	Age    int8
	Gender string
}

func TestEncodeAndDecode(t *testing.T) {

	usr := &User{
		Name:   "AAA",
		Age:    1,
		Gender: "BBB",
	}

	buf := new(bytes.Buffer)

	if err := binparse.NewEncoder(buf).Encode(usr); err != nil {
		t.Fatal("encoder error: ", err)
	}

	t.Log("Bytes: ", buf.Bytes())

	var dusr User

	if err := binparse.NewDecoder(buf).Decode(&dusr); err != nil {
		t.Fatal(err)
	}

	t.Logf("Name: %s | Age: %d | Gender: %s\n", dusr.Name, dusr.Age, dusr.Gender)
}

package binpack_test

import (
	"bytes"
	"testing"

	"github.com/vxcute/binpack"
)

type User struct {
	Name   string
	Age    int8
	Gender string
}

func TestEncodeAndDecode(t *testing.T) {

	usr := &User{
		Name:   "Ahmed",
		Age:    18,
		Gender: "Male",
	}

	buf := new(bytes.Buffer)

	err := binpack.Pack(buf, usr) 

	if err != nil {
		t.Fatal(err)
	}

	t.Log(buf.Bytes())

	var user User 

	if err := binpack.Unpack(buf.Bytes(), &user); err != nil {
		t.Fatal(err)
	}

	t.Log(usr)
	t.Log(user)
}

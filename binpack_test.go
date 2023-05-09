package binpack_test

import (
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

	buf, err := binpack.Pack(usr) 

	if err != nil {
		t.Fatal(err)
	}

	t.Log(buf)

	var user User 

	if err := binpack.Unpack(buf, &user); err != nil {
		t.Fatal(err)
	}

	t.Log(usr)
	t.Log(user)
}

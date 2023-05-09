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
		Name:   "AAA",
		Age:    1,
		Gender: "BBB",
	}

	b, err := binpack.Pack(usr) 

	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)

	var user User 

	if err := binpack.Unpack(b, &user); err != nil {
		t.Fatal(err)
	}

	t.Log(usr)
	t.Log(user)
}

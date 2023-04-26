package binparse_test

import (
	"bytes"
	"testing"

	"github.com/vxcute/binparse"
)

func TestEncode(t *testing.T) {

	usr := struct {
		Name   string
		Age    int
		Gender string
	}{
		Name:   "Ahmed",
		Age:    18,
		Gender: "Male",
	}

	buf := new(bytes.Buffer)

	if err := binparse.NewEncoder(buf).Encode(usr); err != nil {
		t.Fatal("encoder error: ", err)
	}

	t.Log("Bytes: ", buf.Bytes())
}

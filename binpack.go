package binpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
)


var byteOrder binary.ByteOrder = binary.BigEndian

func Pack(v any) ([]byte,error) {

	val := reflect.Indirect(reflect.ValueOf(v))

	buf := new(bytes.Buffer)

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i) 

		if field.CanInterface() {
			switch field.Kind() {
			case reflect.String:
				if err := binary.Write(buf, byteOrder, []byte(field.String() + "\x00")); err != nil {
					return nil, err
				}
			case reflect.Struct:
				b, err := Pack(field.Interface())
				if err != nil {
					return nil, err
				}
				if err := binary.Write(buf, byteOrder, b); err != nil {
					return nil, err
				}
			default:
				if err := binary.Write(buf, byteOrder, field.Interface()); err != nil {
					return nil,err
				}
			}
		}
	}

	return buf.Bytes(), nil
}

func Unpack(buf []byte, v any) error {
	
	iv := reflect.ValueOf(v)

	if iv.Kind() != reflect.Ptr {
		return errors.New("binpack: not a pointer")
	} else if iv.IsNil() {
		return errors.New("binpack: nil ptr")
	}

	iv = iv.Elem()
	it := iv.Type()

	if it.Kind() != reflect.Struct {
		return errors.New("binpack: not a struct")
	}

	var (
		terminator int              = 0
		rest       bool             = false
	)

	for i := 0; i < it.NumField(); i++ {

		fv := iv.Field(i)
		ft := it.Field(i)
		tag := ft.Tag.Get("binpack")

		if !fv.CanSet() || tag == "-" {
			continue
		}

		tagOpts := strings.Split(tag, ",")

		for _, opt := range tagOpts {
			switch opt {
			case "little":
				byteOrder = binary.LittleEndian
			case "rest":
				rest = true
			}
		}

		switch fv.Kind() {

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:

			if len(buf) == 0 {
				return io.EOF
			}

			err := binary.Read(bytes.NewReader(buf), byteOrder, fv.Addr().Interface())

			if err != nil {
				return err
			}

			buf = buf[fv.Type().Size():]
		case reflect.String:
			n := bytes.IndexByte(buf, byte(terminator))

			if n == -1 {
				return errors.New("binpack: missing terminator")
			}

			fv.SetString(string(buf[:n]))

			buf = buf[n+1:]
		case reflect.Slice:
			if rest {
				fv.SetBytes(buf)
				buf = nil
			}
		default:
			return errors.New("binpack: invalid type")
		}
	}

	return nil
}
package binparse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
)

type Decoder interface {
	Decode(v any) error
}

type Encoder interface {
	Encode(v any) error
}

type decoder struct {
	buf []byte
}

type encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) Encoder {
	return &encoder{w: w}
}

func (e *encoder) Encode(v any) error {

	val := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)

		if field.CanInterface() {

			switch field.Kind() {
			case reflect.String:
				if err := binary.Write(e.w, binary.BigEndian, []byte(field.String()+"\x00")); err != nil {
					return err
				}
			default:
				if err := binary.Write(e.w, binary.BigEndian, field.Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func NewDecoder(r io.Reader) Decoder {

	d := &decoder{
		buf: make([]byte, 0),
	}

	buf, err := io.ReadAll(r)

	if err != nil {
		return nil
	}

	d.buf = buf

	return d
}

func (d *decoder) Decode(v any) error {

	iv := reflect.ValueOf(v)

	if iv.Kind() != reflect.Ptr {
		return errors.New("bin: not a pointer")
	} else if iv.IsNil() {
		return errors.New("bin: nil ptr")
	}

	iv = iv.Elem()
	it := iv.Type()

	if it.Kind() != reflect.Struct {
		return errors.New("bin: not a struct")
	}

	var (
		byteOrder  binary.ByteOrder = binary.BigEndian
		terminator int              = 0
		rest       bool             = false
	)

	for i := 0; i < it.NumField(); i++ {

		fv := iv.Field(i)
		ft := it.Field(i)
		tag := ft.Tag.Get("bin")

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
			reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint:

			err := binary.Read(bytes.NewReader(d.buf), byteOrder, fv.Addr().Interface())

			if err != nil {
				return err
			}

			d.buf = d.buf[fv.Type().Size():]
		case reflect.String:
			n := bytes.IndexByte(d.buf, byte(terminator))

			if n == -1 {
				return errors.New("bin: missing terminator")
			}

			fv.SetString(string(d.buf[:n]))

			d.buf = d.buf[n+1:]
		case reflect.Slice:
			if rest {
				fv.SetBytes(d.buf)
				d.buf = nil
			}

		case reflect.Struct:
			err := d.Decode(fv.Addr().Interface())
			if err != nil {
				return errors.New("bin: failed to decode inner struct")
			}
		default:
			return errors.New("bin: invalid type")
		}
	}

	return nil
}

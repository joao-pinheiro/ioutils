package ioutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func Write(dest io.Writer, src interface{}, order binary.ByteOrder) (int, error) {

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int32:
		buf := make([]byte, 4)
		order.PutUint32(buf, uint32(v.Int()))
		return dest.Write(buf)

	case reflect.Uint32, reflect.Uintptr:
		buf := make([]byte, 4)
		order.PutUint32(buf, uint32(v.Uint()))
		return  dest.Write(buf)

	case reflect.Int16:
		buf := make([]byte, 2)
		order.PutUint16(buf, uint16(v.Int()))
		return dest.Write(buf)

	case reflect.Uint16:
		buf := make([]byte, 2)
		order.PutUint16(buf, uint16(v.Uint()))
		return dest.Write(buf)

	case reflect.Uint64:
		buf := make([]byte, 8)
		order.PutUint64(buf, v.Uint())
		return  dest.Write(buf)

	case reflect.Int64:
		buf := make([]byte, 8)
		order.PutUint64(buf, uint64(v.Int()))
		return dest.Write(buf)

	case reflect.Int8, reflect.Uint8:
		return dest.Write([]byte{byte(v.Uint())})

	case reflect.Struct:
		w := 0
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			if field.CanInterface() {
				tmp, err := Write(dest, field, order)
				if err != nil {
					return w, err
				}
				w+=tmp
			}
		}
		return w, nil

	case reflect.Slice, reflect.Array:
		type_len := 0
		w := 0
		if v.Kind() == reflect.Slice {
			type_len = v.Len()
		} else {
			type_len = v.Type().Len()
		}
		for i := 0; i < type_len; i++ {
			tmp, err := Write(dest, v.Index(i).Interface(), order)
			if err != nil {
				return w, err
			}
			w+= tmp
		}
		return w, nil
	}

	return 0, errors.New(fmt.Sprintf("datatype for field %s is not supported", v.Type().Name()))
}

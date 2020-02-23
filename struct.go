package ioutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

/**
 * Reads struct dest{} from reader
 *  Data types supported:
 *    - int8/16/32/64,
 *    - uint8/16,32/64
 *    - slices/arrays
 *    - nested structs
 */
func ReadStruct(reader io.Reader, dest interface{}, order binary.ByteOrder) error {

	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		return fmt.Errorf("dest must be a pointer")
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			ftype := field.Type()
			size := int(ftype.Size())

			if !field.CanSet() {
				return fmt.Errorf("field %s is not settable", ftype.Name())
			}

			if field.CanInterface() {
				switch field.Kind() {
				// Recurse structures
				case reflect.Struct:
					if err := ReadStruct(reader, field.Addr().Interface(), order); err != nil {
						return err
					}

				// Read signed ints
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					n, err := ReadInt(reader, size, order)
					if err != nil {
						return err
					}
					if field.OverflowInt(n) {
						return fmt.Errorf("value %s overflow on field %s", field.Kind(), ftype.Name())
					}
					field.SetInt(n)

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					n, err := ReadUint(reader, size, order)
					if err != nil {
						return err
					}
					if field.OverflowUint(n) {
						return fmt.Errorf("value %s overflow on field %s", field.Kind(), ftype.Name())
					}
					field.SetUint(n)

				case reflect.Slice, reflect.Array:
					if err := ReadArray(reader, field.Addr().Interface(), order); err != nil {
						return err
					}

				default:
					return fmt.Errorf("datatype for field %s is not supported", ftype.Name())
				}
			}
		}

	case reflect.Slice, reflect.Array:
		if err := ReadArray(reader, reflect.ValueOf(dest).Interface(), order); err != nil {
			return err
		}

	default:
		return errors.New(fmt.Sprintf("unsupported data type %s in %s", v.Kind(), v.Type().Name()))
	}
	return nil
}

/**
 * Reads array or slice from reader
 */
func ReadArray(reader io.Reader, dest interface{}, order binary.ByteOrder) error {

	field := reflect.ValueOf(dest)
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	kind := field.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		var type_len int
		ftype := field.Type()
		if kind == reflect.Slice {
			type_len = field.Len()
		} else {
			type_len = ftype.Len()
		}

		for idx := 0; idx < type_len; idx++ {
			cell := field.Index(idx)
			cell_size := int(cell.Type().Size())

			switch cell.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				n, err := ReadInt(reader, cell_size, order)
				if err != nil {
					return err
				}
				if cell.OverflowInt(n) {
					return fmt.Errorf("value %s overflow on field %s", cell.Kind(), ftype.Name())
				}
				cell.SetInt(n)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				n, err := ReadUint(reader, cell_size, order)
				if err != nil {
					return err
				}
				if cell.OverflowUint(n) {
					return fmt.Errorf("value %s overflow on field %s", cell.Kind(), ftype.Name())
				}
				cell.SetUint(n)

			case reflect.Struct:
				if err := ReadStruct(reader, cell.Interface(), order); err != nil {
					return err
				}

			default:
				return fmt.Errorf("datatype for array %s is not supported", ftype.Name())
			}
		}
	} else {
		return fmt.Errorf("datatype %s not supported in %s", field.Kind(), field.Type().Name())
	}
	return nil
}

/***
 * Reads count bytes from reader
 */
func ReadBytes(reader io.Reader, count int) ([]byte, error) {
	result := make([]byte, count)
	total, err := reader.Read(result)
	if total == count && err == nil {
		return result, err
	}
	return nil, err
}

/***
 * Reads size bytes from reader as int64
 */
func ReadInt(reader io.Reader, size int, order binary.ByteOrder) (int64, error) {
	tmp, err := ReadBytes(reader, size)
	if err != nil {
		return 0, err
	}
	return int64(order.Uint64(PadBytes(tmp, 8, order))), nil
}

/***
 * Reads size bytes from reader as uint64
 */
func ReadUint(reader io.Reader, size int, order binary.ByteOrder) (uint64, error) {
	tmp, err := ReadBytes(reader, size)
	if err != nil {
		return 0, err
	}
	return uint64(order.Uint64(PadBytes(tmp, 8, order))), nil
}

/***
 * Generates a byte buffer to padded the desired size according to endianess
 */
func PadBytes(src []byte, desired_size int, order binary.ByteOrder) []byte {
	l := desired_size - len(src)

	if l <= 0 {
		return src
	}
	i := make([]byte, l)
	if (order == binary.LittleEndian) {
		return append(src, i...)
	}
	return append(i, src...)
}

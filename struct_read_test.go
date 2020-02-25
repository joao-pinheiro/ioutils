package ioutils

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

var buffer01 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

type uint_array [8]byte

// simple uint struct
type uintstruct01 struct {
	Val_16 uint16
	Val_8  uint8
	Val_32 uint32
	Val_64 uint64
	Pad    uint8
}

type UVAL16 struct {
	Val_16 uint16
}

type UVAL32 struct {
	Val_32 uint32
}

type UVAL64 struct {
	Val_64 uint64
}

// uint struct with nested structs
type uintstruct02 struct {
	U16  UVAL16
	Val8 uint8
	U32  UVAL32
	U64  UVAL64
	Pad  uint8
}

// uint struct with anonymous nested structs
type uintstruct03 struct {
	UVAL16
	Val8 uint8
	UVAL32
	UVAL64
	Pad uint8
}

type URecord [2]uint32

// uint struct with arrays
type uintstruct04 struct {
	R1 URecord
	R2 URecord
}

// uint struct with slice
// slice is preinitialized
type uintstruct05 struct {
	R1 []uint32
	R2 []uint32
}

//******************************************************************************
// Byte slice/array
//******************************************************************************
//@todo: array test

func TestByteSlice(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	// read buffer01 into 2 slices
	slice1 := make([]byte, 8)
	slice2 := make([]byte, 8)
	err := ReadStruct(reader, &slice1, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, buffer01[:8], slice1)

	err = ReadStruct(reader, &slice2, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, buffer01[8:], slice2)
}

func TestArray(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	// read buffer01 into 2 slices
	slice1 := &uint_array{}
	slice2 := &uint_array{}
	err := ReadStruct(reader, slice1, binary.BigEndian)
	assert.Nil(t, err)
	for i:=0; i < 8; i++ {
		assert.Equal(t, buffer01[i], slice1[i])
	}
	err = ReadStruct(reader, slice2, binary.BigEndian)
	assert.Nil(t, err)
	for i:=0; i < 8; i++ {
		assert.Equal(t, buffer01[i+8], slice2[i])
	}
}

//******************************************************************************
// Uint Big Endian
//******************************************************************************

func TestUIntStruct_bigendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct01{}
	err := ReadStruct(reader, s01, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x01), s01.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val_8)
	assert.Equal(t, uint32(0x03040506), s01.Val_32)
	assert.Equal(t, uint64(0x0708090a0b0c0d0e), s01.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUintStruct_Nested_bigendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct02{}
	err := ReadStruct(reader, s01, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x01), s01.U16.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val8)
	assert.Equal(t, uint32(0x03040506), s01.U32.Val_32)
	assert.Equal(t, uint64(0x0708090a0b0c0d0e), s01.U64.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUIntStruct_anonymous_bigendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct03{}
	err := ReadStruct(reader, s01, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x01), s01.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val8)
	assert.Equal(t, uint32(0x03040506), s01.Val_32)
	assert.Equal(t, uint64(0x0708090a0b0c0d0e), s01.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUIntStruct_array_bigendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct04{}
	err := ReadStruct(reader, s01, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0x00010203), s01.R1[0])
	assert.Equal(t, uint32(0x04050607), s01.R1[1])
	assert.Equal(t, uint32(0x08090a0b), s01.R2[0])
	assert.Equal(t, uint32(0x0c0d0e0f), s01.R2[1])
}

func TestUIntStruct_slice_bigendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct05{
		R1: make([]uint32, 2),
		R2: make([]uint32, 2),
	}
	err := ReadStruct(reader, s01, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0x00010203), s01.R1[0])
	assert.Equal(t, uint32(0x04050607), s01.R1[1])
	assert.Equal(t, uint32(0x08090a0b), s01.R2[0])
	assert.Equal(t, uint32(0x0c0d0e0f), s01.R2[1])
}

//******************************************************************************
// Uint Little Endian
//******************************************************************************

func TestUIntStruct_littleendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct01{}
	err := ReadStruct(reader, s01, binary.LittleEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x0100), s01.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val_8)
	assert.Equal(t, uint32(0x06050403), s01.Val_32)
	assert.Equal(t, uint64(0x0e0d0c0b0a090807), s01.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUintStruct_Nested_littleendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct02{}
	err := ReadStruct(reader, s01, binary.LittleEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x0100), s01.U16.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val8)
	assert.Equal(t, uint32(0x06050403), s01.U32.Val_32)
	assert.Equal(t, uint64(0x0e0d0c0b0a090807), s01.U64.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUIntStruct_anonymous_littleendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct03{}
	err := ReadStruct(reader, s01, binary.LittleEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x0100), s01.Val_16)
	assert.Equal(t, uint8(0x02), s01.Val8)
	assert.Equal(t, uint32(0x06050403), s01.Val_32)
	assert.Equal(t, uint64(0x0e0d0c0b0a090807), s01.Val_64)
	assert.Equal(t, uint8(0x0F), s01.Pad)
}

func TestUIntStruct_array_littleendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct04{}
	err := ReadStruct(reader, s01, binary.LittleEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0x03020100), s01.R1[0])
	assert.Equal(t, uint32(0x07060504), s01.R1[1])
	assert.Equal(t, uint32(0x0b0a0908), s01.R2[0])
	assert.Equal(t, uint32(0x0f0e0d0c), s01.R2[1])
}

func TestUIntStruct_slice_littleendian(t *testing.T) {
	reader := bytes.NewBuffer(buffer01)
	s01 := &uintstruct05{
		R1: make([]uint32, 2),
		R2: make([]uint32, 2),
	}
	err := ReadStruct(reader, s01, binary.LittleEndian)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0x03020100), s01.R1[0])
	assert.Equal(t, uint32(0x07060504), s01.R1[1])
	assert.Equal(t, uint32(0x0b0a0908), s01.R2[0])
	assert.Equal(t, uint32(0x0f0e0d0c), s01.R2[1])
}

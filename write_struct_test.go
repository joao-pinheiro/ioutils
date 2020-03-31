package ioutils

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

//******************************************************************************
// Byte slice/array
//******************************************************************************
//@todo: array test

func TestWriteByteSlice(t *testing.T) {
	writer := bytes.NewBuffer(make([]byte, 0))
	slen := 32
	slice1 := make([]byte, slen)
	for i:= 0; i < slen; i++ {
		slice1[i] = uint8(i)
 	}
	w, err := Write(writer, &slice1, binary.BigEndian)
	assert.Nil(t, err)
	assert.Equal(t, slen, w)

	reader := bytes.NewBuffer(writer.Bytes())
	for i := 0; i < slen; i++ {
		b, err := reader.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, b, slice1[i])
	}
}


ioutils - "C"-style binary read & write with structs
====================================================

Golang package that provides reading & writing of structs from/to binary files, similar to C and other languages. This is
achieved by reflection, and can read streams that conform to io.Reader interface.

Supported data types:
- uint8/16/32/64 (littleEndian and bigEndian)
- int8/16/32/64 (littleEndian and bigEndian)
- byte
- slices
- arrays
- nested structs


Example:
```go
package png

import (
	"io"
	"encoding/binary"
	"github.com/joao-pinheiro/ioutils"
)

type PNGSignature [8]byte

type Chunk struct {
	Size uint32
	Type [4]byte
}

type IHDR struct {
	Width       uint32
	Height      uint32
	Bit         uint8
	ColorType   uint8
	Compression uint8
	Filter      uint8
	Interlace   uint8
}

type PNGHeader struct {
	PNGSignature
	Chunk
	IHDR
	CRC uint32
}

func ReadPNGHeader(file io.Reader) (*PNGHeader, error) {
	header := &PNGHeader{}
	if err := ReadStruct(file, header, binary.BigEndian); err != nil {
		return nil, err
	}
	return header, nil
}
```
 
 
 #todo:
 type yy struct {
   A string
   B string
   C string
 }
 type xxx struct {
   Abc [10]yyy
  } 
  
  # todo
  // check read-write with 64 bit unsigned values
  
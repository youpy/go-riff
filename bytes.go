package riff

import (
	"io"
	"io/ioutil"
)

type Bytes struct {
	Offset uint32
	Data   []byte
}

func newBytes(r io.Reader) (bytes *Bytes, err error) {
	raw_bytes, err := ioutil.ReadAll(r)
	bytes = &Bytes{0, raw_bytes}

	return
}

func (bytes *Bytes) readLEUint32() uint32 {
	offset := bytes.Offset
	data := bytes.Data

	defer func() {
		bytes.Offset += 4
	}()

	return uint32(data[offset+3])<<24 +
		uint32(data[offset+2])<<16 +
		uint32(data[offset+1])<<8 +
		uint32(data[offset])
}

func (bytes *Bytes) readLEUint16() uint16 {
	offset := bytes.Offset
	data := bytes.Data

	defer func() {
		bytes.Offset += 2
	}()

	return uint16(data[offset+1])<<8 + uint16(data[offset])
}

func (bytes *Bytes) readLEInt16() int16 {
	offset := bytes.Offset
	data := bytes.Data

	defer func() {
		bytes.Offset += 2
	}()

	return int16(data[offset+1])<<8 + int16(data[offset])
}

func (bytes *Bytes) readBytes(size uint32) []byte {
	defer func() {
		bytes.Offset += size
	}()

	return bytes.Data[bytes.Offset : bytes.Offset+size]
}

func (bytes *Bytes) eof() bool {
	return bytes.Offset >= uint32(len(bytes.Data))
}

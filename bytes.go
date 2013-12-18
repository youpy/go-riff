package riff

import (
	"io"
)

type Bytes struct {
	Offset uint32
	Reader io.ReaderAt
}

func newBytes(r io.ReaderAt) (bytes *Bytes, err error) {
	bytes = &Bytes{0, r}

	return
}

func (bytes *Bytes) readLEUint32() uint32 {
	offset := bytes.Offset
	data := make([]byte, 4)

	n, err := bytes.Reader.ReadAt(data, int64(offset))

	if err != nil || n < 4 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.Offset += 4
	}()

	return uint32(data[3])<<24 +
		uint32(data[2])<<16 +
		uint32(data[1])<<8 +
		uint32(data[0])
}

func (bytes *Bytes) readLEUint16() uint16 {
	offset := bytes.Offset
	data := make([]byte, 2)

	n, err := bytes.Reader.ReadAt(data, int64(offset))

	if err != nil || n < 2 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.Offset += 2
	}()

	return uint16(data[1])<<8 + uint16(data[0])
}

func (bytes *Bytes) readLEInt16() int16 {
	offset := bytes.Offset
	data := make([]byte, 2)

	n, err := bytes.Reader.ReadAt(data, int64(offset))

	if err != nil || n < 2 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.Offset += 2
	}()

	return int16(data[offset+1])<<8 + int16(data[offset])
}

func (bytes *Bytes) readBytes(size uint32) []byte {
	offset := bytes.Offset
	data := make([]byte, size)

	n, err := bytes.Reader.ReadAt(data, int64(offset))

	if err != nil || n < int(size) {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.Offset += size
	}()

	return data[0:size]
}

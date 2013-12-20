package riff

import (
	"encoding/binary"
	"io"
)

type Encoder struct {
	io.Writer
}

type writeCallback func(w io.Writer)

func NewEncoder(w io.Writer, fileType []byte, fileSize uint32) *Encoder {
	w.Write([]byte("RIFF"))
	binary.Write(w, binary.LittleEndian, fileSize)
	w.Write(fileType)

	return &Encoder{w}
}

func (e *Encoder) WriteChunk(chunkID []byte, chunkSize uint32, cb writeCallback) (err error) {
	_, err = e.Write(chunkID)

	if err != nil {
		return
	}

	err = binary.Write(e, binary.LittleEndian, chunkSize)

	if err != nil {
		return
	}

	cb(e)

	return
}

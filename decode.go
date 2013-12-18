package riff

import (
	"errors"
	"io"
)

type RIFFChunk struct {
	FileSize uint32
	FileType []byte
	Chunks   []*Chunk
}

type Chunk struct {
	ChunkID   []byte
	ChunkSize uint32
	Data      []byte
}

func Decode(r io.Reader) (chunk *RIFFChunk, err error) {
	chunk, err = decodeRIFFChunk(r)

	return
}

func decodeRIFFChunk(r io.Reader) (chunk *RIFFChunk, err error) {
	bytes, err := newBytes(r)

	if err != nil {
		err = errors.New("Can't read RIFF file")
		return
	}

	chunkId := bytes.readBytes(4)

	if string(chunkId[:]) != "RIFF" {
		err = errors.New("Given bytes is not a RIFF format")
		return
	}

	fileSize := bytes.readLEUint32()
	fileType := bytes.readBytes(4)

	chunk = &RIFFChunk{fileSize, fileType, make([]*Chunk, 0)}

	for !bytes.eof() {
		chunkId = bytes.readBytes(4)
		chunkSize := bytes.readLEUint32()

		if chunkSize%2 == 1 {
			chunkSize += 1
		}

		data := bytes.readBytes(chunkSize)

		chunk.Chunks = append(chunk.Chunks, &Chunk{chunkId, chunkSize, data})
	}

	return
}

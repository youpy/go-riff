package riff

import (
	"errors"
	"io"
)

type RIFFReader interface {
	io.Reader
	io.ReaderAt
}

type RIFFChunk struct {
	FileSize uint32
	FileType []byte
	Chunks   []*Chunk
}

type Chunk struct {
	ChunkID   []byte
	ChunkSize uint32
	RIFFReader
}

func Decode(r RIFFReader) (chunk *RIFFChunk, err error) {
	chunk, err = decodeRIFFChunk(r)

	return
}

func decodeRIFFChunk(r RIFFReader) (chunk *RIFFChunk, err error) {
	bytes := newBytes(r)

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

	for bytes.Offset < fileSize {
		chunkId = bytes.readBytes(4)
		chunkSize := bytes.readLEUint32()
		offset := bytes.Offset

		if chunkSize%2 == 1 {
			chunkSize += 1
		}

		bytes.Offset += chunkSize

		chunk.Chunks = append(
			chunk.Chunks,
			&Chunk{
				chunkId,
				chunkSize,
				io.NewSectionReader(r, int64(offset), int64(chunkSize))})
	}

	return
}

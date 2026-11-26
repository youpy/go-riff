package riff

import (
	"errors"
	"io"
)

type RIFFReader interface {
	io.Reader
	io.ReaderAt
}

type Reader struct {
	RIFFReader
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

func NewReader(r RIFFReader) *Reader {
	return &Reader{r}
}

func (r *Reader) Read() (chunk *RIFFChunk, err error) {
	chunk, err = readRIFFChunk(r)

	return
}

func readRIFFChunk(r *Reader) (chunk *RIFFChunk, err error) {
	bytes := newByteReader(r)

	chunkId := bytes.readBytes(4)

	if string(chunkId[:]) != "RIFF" {
		err = errors.New("Given bytes is not a RIFF format")
		return
	}

	fileSize := bytes.readLEUint32()
	fileType := bytes.readBytes(4)

	chunk = &RIFFChunk{fileSize, fileType, make([]*Chunk, 0)}

	// fileSize does not include the eight-byte RIFF chunk header.
	for uint64(bytes.offset) < uint64(fileSize)+8 {
		chunkId = bytes.readBytes(4)
		chunkSize := bytes.readLEUint32()
		offset := bytes.offset

		// chunkSize is padded to an even number of bytes by appending a pad byte if necessary.
		bytes.offset += chunkSize + chunkSize%2

		chunk.Chunks = append(
			chunk.Chunks,
			&Chunk{
				chunkId,
				chunkSize,
				io.NewSectionReader(r, int64(offset), int64(chunkSize))})
	}

	return
}

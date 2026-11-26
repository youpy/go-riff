package riff

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"testing"
)

func TestReadRIFFChunkBoundary(t *testing.T) {
	tests := []struct {
		name     string
		payloads []string
	}{
		{"no chunks", nil},
		{"one empty chunk", []string{""}},
		{"two empty chunks", []string{"", ""}},
		{"empty chunk after data", []string{"ab", ""}},
		{"odd chunk before next chunk", []string{"a", "bc"}},
		{"multiple empty chunks after data", []string{"ab", "", ""}},
		{"nonempty chunk", []string{"ab"}},
		{"empty chunk before data", []string{"", "ab"}},
	}

	for _, test := range tests {
		data := []byte("RIFF\x00\x00\x00\x00WAVE")
		for _, payload := range test.payloads {
			header := []byte("data\x00\x00\x00\x00")
			binary.LittleEndian.PutUint32(header[4:], uint32(len(payload)))
			data = append(data, header...)
			data = append(data, []byte(payload)...)
			if len(payload)%2 == 1 {
				data = append(data, 0)
			}
		}
		fileSize := uint32(len(data) - 8)
		binary.LittleEndian.PutUint32(data[4:8], fileSize)
		// This chunk is outside the RIFF container and must not be read.
		data = append(data, []byte("JUNK\x00\x00\x00\x00")...)

		riff, err := NewReader(bytes.NewReader(data)).Read()
		if err != nil {
			t.Fatalf("%s: Read failed: %v", test.name, err)
		}
		if riff.FileSize != fileSize {
			t.Errorf("%s: FileSize = %d, want %d", test.name, riff.FileSize, fileSize)
		}
		if len(riff.Chunks) != len(test.payloads) {
			t.Errorf("%s: got %d chunks, want %d", test.name, len(riff.Chunks), len(test.payloads))
			continue
		}
		for i, chunk := range riff.Chunks {
			if string(chunk.ChunkID) != "data" || chunk.ChunkSize != uint32(len(test.payloads[i])) {
				t.Errorf("%s: chunk %d has ID %q and size %d", test.name, i, chunk.ChunkID, chunk.ChunkSize)
			}
			payload, err := ioutil.ReadAll(chunk)
			if err != nil {
				t.Fatalf("%s: reading chunk %d failed: %v", test.name, i, err)
			}
			if string(payload) != test.payloads[i] {
				t.Errorf("%s: chunk %d data = %q, want %q", test.name, i, payload, test.payloads[i])
			}
		}
	}
}

func TestReadRIFF(t *testing.T) {
	testFiles := []testFile{
		testFile{
			"a.wav",
			3,
			243800,
			"WAVE"},
		testFile{
			"1_webp_a.webp",
			3,
			23396,
			"WEBP"}}

	for _, testFile := range testFiles {
		file, err := fixtureFile(testFile.Name)

		if err != nil {
			t.Fatalf("Failed to open fixture file")
		}

		reader := NewReader(file)
		riff, err := reader.Read()

		if err != nil {
			t.Fatal(err)
		}

		for _, chunk := range riff.Chunks {
			t.Logf("Chunk ID: %s", string(chunk.ChunkID[:]))
		}

		if len(riff.Chunks) != testFile.ChunkSize {
			t.Fatalf("Invalid length of chunks")
		}

		if riff.FileSize != testFile.FileSize {
			t.Fatalf("File size is invalid: %d", riff.FileSize)
		}

		if string(riff.FileType[:]) != testFile.FileType {
			t.Fatalf("File type is invalid: %s", riff.FileType)
		}

		data, err := ioutil.ReadAll(riff.Chunks[0])

		if err != nil {
			t.Fatalf("Can't read data from chunk")
		}

		t.Logf("Length of the first chunk: %d", len(data))

		file.Close()
	}
}

func TestReadNonRIFF(t *testing.T) {
	file, err := fixtureFile("../reader.go")

	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	_, err = reader.Read()

	if err.Error() != "Given bytes is not a RIFF format" {
		t.Fatal("Non-RIFF file should not be read")
	}
}

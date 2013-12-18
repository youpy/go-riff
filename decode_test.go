package riff

import (
	"testing"
)

type TestFile struct {
	Name      string
	ChunkSize int
	FileSize  uint32
	FileType  string
}

func TestDecodeRIFF(t *testing.T) {
	testFiles := []TestFile{
		TestFile{
			"a.wav",
			3,
			243800,
			"WAVE"},
		TestFile{
			"1_webp_a.webp",
			3,
			23396,
			"WEBP"}}

	for _, testFile := range testFiles {
		file, err := fixtureFile(testFile.Name)

		if err != nil {
			t.Fatalf("Failed to open fixture file")
		}

		riff, err := Decode(file)

		if err != nil {
			t.Fatal(err)
		}

		for _, chunk := range riff.Chunks {
			t.Log(string(chunk.ChunkID[:]))
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
	}
}

func TestDecodeNonRIFF(t *testing.T) {
	file, err := fixtureFile("../decode.go")

	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	_, err = Decode(file)

	if err.Error() != "Given bytes is not a RIFF format" {
		t.Fatal("Non-RIFF file should not be decoded")
	}
}

package riff

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncodeRIFF(t *testing.T) {
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

		outfile, err := ioutil.TempFile("/tmp", "outfile")

		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			outfile.Close()
			os.Remove(outfile.Name())
		}()

		encoder := NewEncoder(outfile, riff.FileType, riff.FileSize)

		for _, chunk := range riff.Chunks {
			encoder.WriteChunk(chunk.ChunkID, chunk.ChunkSize, func(w io.Writer) {
				buf, err := ioutil.ReadAll(chunk)

				if err != nil {
					t.Fatal(err)
				}

				w.Write(buf)
			})
		}

		outfile.Close()

		// reopen to check file content
		file, err = os.Open(outfile.Name())

		if err != nil {
			t.Fatal(err)
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

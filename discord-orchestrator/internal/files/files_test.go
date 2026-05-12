package files

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTail_Last20Lines(t *testing.T) {
	path := filepath.Join("testdata", "hundred-lines.txt")
	result, err := Tail(path, 20)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	lines := strings.Split(result, "\n")
	if len(lines) != 20 {
		t.Errorf("expected 20 lines, got %d", len(lines))
	}
	// Last line should be "line 100"
	if lines[19] != "line 100" {
		t.Errorf("last line = %q, want %q", lines[19], "line 100")
	}
	// First returned line should be "line 081"
	if lines[0] != "line 081" {
		t.Errorf("first line = %q, want %q", lines[0], "line 081")
	}
}

func TestReadTail_TruncatedWhenFileLarger(t *testing.T) {
	// Create a file larger than maxBytes
	dir := t.TempDir()
	path := filepath.Join(dir, "big.log")
	content := bytes.Repeat([]byte("x"), 200)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	data, truncated, err := ReadTail(path, 100)
	if err != nil {
		t.Fatalf("ReadTail: %v", err)
	}
	if !truncated {
		t.Error("expected truncated=true for file larger than maxBytes")
	}
	if int64(len(data)) != 100 {
		t.Errorf("expected 100 bytes, got %d", len(data))
	}
}

func TestReadTail_NotTruncatedWhenFileSmaller(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "small.log")
	content := []byte("hello world")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	data, truncated, err := ReadTail(path, 1000)
	if err != nil {
		t.Fatalf("ReadTail: %v", err)
	}
	if truncated {
		t.Error("expected truncated=false for file smaller than maxBytes")
	}
	if string(data) != "hello world" {
		t.Errorf("data = %q, want %q", data, "hello world")
	}
}

func TestChunkFile_ChunkBoundaries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chunked.bin")

	// 2.5 × chunkSize — expect 3 chunks
	chunkSize := int64(100)
	content := bytes.Repeat([]byte("a"), 250)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	chunks, err := ChunkFile(path, chunkSize)
	if err != nil {
		t.Fatalf("ChunkFile: %v", err)
	}
	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks for 2.5× file, got %d", len(chunks))
	}
	if int64(len(chunks[0])) != chunkSize {
		t.Errorf("chunk[0] size = %d, want %d", len(chunks[0]), chunkSize)
	}
	if int64(len(chunks[1])) != chunkSize {
		t.Errorf("chunk[1] size = %d, want %d", len(chunks[1]), chunkSize)
	}
	if int64(len(chunks[2])) != 50 {
		t.Errorf("chunk[2] size = %d, want 50", len(chunks[2]))
	}
}

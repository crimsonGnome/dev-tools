package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	// DefaultTailBytes is the default maximum bytes read by ReadTail (50 KB).
	DefaultTailBytes int64 = 50 * 1024

	// ChunkSize is the maximum bytes per chunk produced by ChunkFile (8 MB).
	ChunkSize int64 = 8 * 1024 * 1024

	// DefaultTailLines is the default line count for Tail.
	DefaultTailLines = 20
)

// Tail returns the last n lines from the file at path, joined by newlines.
func Tail(path string, lines int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("files: open %q: %w", path, err)
	}
	defer f.Close()

	// Read all lines into a ring buffer of size `lines`
	scanner := bufio.NewScanner(f)
	buf := make([]string, 0, lines)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
		if len(buf) > lines {
			buf = buf[1:]
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("files: scan %q: %w", path, err)
	}

	result := ""
	for i, line := range buf {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result, nil
}

// ReadTail reads up to maxBytes from the end of the file at path.
// The bool return value is true if the file was larger than maxBytes
// (i.e. the output is truncated from the beginning).
func ReadTail(path string, maxBytes int64) ([]byte, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false, fmt.Errorf("files: open %q: %w", path, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, false, fmt.Errorf("files: stat %q: %w", path, err)
	}

	size := info.Size()
	truncated := size > maxBytes

	offset := int64(0)
	if truncated {
		offset = size - maxBytes
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, false, fmt.Errorf("files: seek %q: %w", path, err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, false, fmt.Errorf("files: read %q: %w", path, err)
	}
	return data, truncated, nil
}

// ChunkFile reads the entire file at path and splits it into sequential chunks
// of at most chunkSize bytes each.
func ChunkFile(path string, chunkSize int64) ([][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("files: read %q: %w", path, err)
	}

	if len(data) == 0 {
		return [][]byte{{}}, nil
	}

	var chunks [][]byte
	for len(data) > 0 {
		end := int(chunkSize)
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[:end])
		data = data[end:]
	}
	return chunks, nil
}

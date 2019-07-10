package filekeeper

import (
	"io"
)

type ProgressReader struct {
	content  []byte
	seek     int
	callback func(int)
}

func NewProgressReader(content []byte, callback func(int)) *ProgressReader {
	return &ProgressReader{
		content:  content,
		seek:     0,
		callback: callback,
	}
}

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	copy(p, r.content[r.seek:])
	delta := len(p)
	if len(r.content[r.seek:]) < delta {
		delta = len(r.content[r.seek:])
	}
	r.seek += delta
	r.callback(r.seek)
	if len(r.content) == r.seek {
		return delta, io.EOF
	}
	return delta, nil
}

package recorder

import (
	"os"
)

type SerialRecorder struct {
	StdRecorder
}

func (r *SerialRecorder) SetWriter(h *os.File) {
	r.writer = h
}

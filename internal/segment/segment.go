package segment

import (
	"github.com/ISSuh/wal/internal/file"
)

type Segment struct {
	file file.File
}

func NewSegment() (*Segment, error) {
	return nil, nil
}

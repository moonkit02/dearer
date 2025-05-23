package files

import (
	"time"

	"github.com/moonkit02/dearer/pkg/git"
)

type List struct {
	Files     []File
	BaseFiles []File
	Renames   map[string]string
	Chunks    map[string]git.Chunks
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

type LineMapping struct {
	Base,
	Delta int
}

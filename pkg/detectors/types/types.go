package types

import (
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/util/file"
)

type DetectorConstructor func() Detector

type Detector interface {
	AcceptDir(dir *file.Path) (bool, error)
	ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error)
}

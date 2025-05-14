package buildgradle

import (
	grdlparser "github.com/moonkit02/dearer/pkg/detectors/dependencies/buildgradle/parser"
	"github.com/moonkit02/dearer/pkg/detectors/dependencies/depsbase"
	"github.com/moonkit02/dearer/pkg/util/file"
)

// Discover parses build.gradle file and add discovered dependencies to report
func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	return grdlparser.Discover(file)
}

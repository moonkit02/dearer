package operationshelper

import (
	"github.com/moonkit02/dearer/pkg/report/operations"
	"github.com/moonkit02/dearer/pkg/report/source"
)

type Operation struct {
	Source source.Source
	Value  operations.Operation
}

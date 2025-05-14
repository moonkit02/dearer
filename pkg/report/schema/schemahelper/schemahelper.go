package schemahelper

import (
	"github.com/moonkit02/dearer/pkg/report/schema"
	"github.com/moonkit02/dearer/pkg/report/source"
)

type Schema struct {
	Source source.Source
	Value  schema.Schema
}

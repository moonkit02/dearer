package queries

import (
	"github.com/moonkit02/dearer/pkg/parser"
)

type ChildMatch interface {
	Match(*parser.Node) *parser.Node
}

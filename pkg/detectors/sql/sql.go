package sql

import (
	"github.com/moonkit02/dearer/pkg/detectors/types"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/util/file"

	reporttypes "github.com/moonkit02/dearer/pkg/report"
	schemadatatype "github.com/moonkit02/dearer/pkg/report/schema/datatype"
)

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report reporttypes.Report) (bool, error) {
	// general sql
	if file.Language != "SQL" &&
		// postgress
		file.Language != "PLpgSQL" && file.Language != "PLSQL" && file.Language != "SQLPL" &&
		// microsoft sql
		file.Language != "TSQL" {
		return false, nil
	}

	return true, nil
}

func ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*schemadatatype.DataType, error) {
	return nil, nil
}

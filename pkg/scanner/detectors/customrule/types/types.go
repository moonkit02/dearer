package types

import (
	detectortypes "github.com/moonkit02/dearer/pkg/scanner/detectors/types"
	"github.com/moonkit02/dearer/pkg/scanner/variableshape"
)

type Data struct {
	Pattern   string
	Datatypes []*detectortypes.Detection
	Variables variableshape.Values
	Value     string
}

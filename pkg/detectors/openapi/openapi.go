package openapi

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/moonkit02/dearer/pkg/detectors/openapi/v2json"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/yaml" // Need to use this as some features use custom JSON unmarshalling

	"github.com/moonkit02/dearer/pkg/detectors/openapi/v2yaml"
	"github.com/moonkit02/dearer/pkg/detectors/openapi/v3json"
	"github.com/moonkit02/dearer/pkg/detectors/openapi/v3yaml"
	"github.com/moonkit02/dearer/pkg/detectors/types"
	"github.com/moonkit02/dearer/pkg/report/detectors"

	reporttypes "github.com/moonkit02/dearer/pkg/report"
)

type version struct {
	Swagger string `yaml:"swagger" json:"swagger"`
	OpenAPI string `yaml:"openapi" json:"openapi"`
}

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
	var err error

	fileType, err := getFileType(file)
	if err != nil {
		log.Debug().Msgf("error in OpenAPI detector: %s", err.Error())
		return false, nil
	}

	switch fileType {
	case detectors.OpenApi2JSONFile:
		return v2json.ProcessFile(detector.idGenerator, file, report)
	case detectors.OpenApi2YAMLFile:
		return v2yaml.ProcessFile(detector.idGenerator, file, report)
	case detectors.OpenApi3JSONFile:
		return v3json.ProcessFile(detector.idGenerator, file, report)
	case detectors.OpenApi3YAMLFile:
		return v3yaml.ProcessFile(detector.idGenerator, file, report)
	}

	return false, nil
}

func getFileType(file *file.FileInfo) (detectors.OpenAPIFileType, error) {
	ext := file.Extension
	if ext != ".yml" && ext != ".yaml" && ext != ".json" {
		return "", nil
	}

	input, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return "", err
	}

	var version version
	if ext == ".json" {
		if isArray(input) { // fallback to json|yaml detector
			return "", nil
		}

		err := json.Unmarshal(input, &version)
		if err != nil {
			return detectors.OpenAPIFileType(""), err
		}
	} else {
		err := yaml.Unmarshal(input, &version)
		if err != nil {
			return detectors.OpenAPIFileType(""), err
		}
	}

	versionString := version.OpenAPI
	if versionString == "" {
		versionString = version.Swagger
	}

	// if we can't determine openapi version we fallback to json|yaml detector
	if versionString == "" {
		return "", nil
	}

	if semver.Compare("v"+versionString, "v3") >= 0 {
		if ext == ".json" {
			return detectors.OpenApi3JSONFile, nil
		} else {
			return detectors.OpenApi3YAMLFile, nil
		}
	}

	if ext == ".json" {
		return detectors.OpenApi2JSONFile, nil
	} else {
		return detectors.OpenApi2YAMLFile, nil
	}

}

func isArray(input []byte) bool {
	return bytes.HasPrefix(bytes.TrimSpace(input), []byte("["))
}

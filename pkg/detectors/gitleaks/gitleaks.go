package gitleaks

import (
	_ "embed"
	"log"
	"strings"

	"github.com/moonkit02/dearer/pkg/detectors/types"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/secret"
	"github.com/moonkit02/dearer/pkg/report/source"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/pelletier/go-toml"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
)

//go:embed gitlab_config.toml
var RawConfig []byte

type detector struct {
	gitleaksDetector *detect.Detector
	idGenerator      nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	var vc config.ViperConfig
	toml.Unmarshal(RawConfig, &vc) //nolint:all,errcheck
	cfg, err := vc.Translate()
	if err != nil {
		log.Fatal(err)
	}

	gitleaksDetector := detect.NewDetector(cfg)

	return &detector{
		gitleaksDetector: gitleaksDetector,
		idGenerator:      idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	findings, err := detector.gitleaksDetector.DetectFiles(file.Path.AbsolutePath)

	if err != nil {
		return false, err
	}

	for _, finding := range findings {
		text := strings.TrimPrefix(finding.Line, "\n")
		report.AddSecretLeak(secret.Secret{
			Description: finding.Description,
		}, source.Source{
			Filename:          file.Path.RelativePath,
			StartLineNumber:   &finding.StartLine,
			StartColumnNumber: &finding.StartColumn,
			EndLineNumber:     &finding.EndLine,
			EndColumnNumber:   &finding.EndColumn,
			Text:              &text,
		})
	}

	return false, nil
}

package poetry

import (
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/toml"

	"github.com/moonkit02/dearer/pkg/detectors/dependencies/depsbase"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/moonkit02/dearer/pkg/util/stringutil"
)

var language = toml.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(table_array_element
	(pair
		(bare_key) @helper_name
		(#match? @helper_name "^name$")
		(string) @param_dependency
	)
	(pair
		(bare_key) @helper_version
		(#match? @helper_version "^version$")
		(string) @param_version
	)
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "poetry"
	report.Language = "python"
	report.PackageManager = "pypi"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_name"].Content()) != "name" ||
			stringutil.StripQuotes(capture["helper_version"].Content()) != "version" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].StartLineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}

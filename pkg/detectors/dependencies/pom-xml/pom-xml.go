package pomxml

import (
	"github.com/rs/zerolog/log"

	"github.com/moonkit02/dearer/pkg/detectors/dependencies/depsbase"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/sitter/xml"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/moonkit02/dearer/pkg/util/stringutil"
)

var language = xml.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(element
	(start_tag
		(tag_name) @helper_dependency
		(#match? @helper_dependency "^dependency$")
	)
 ) @param_dependency
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "pom-xml"
	report.Language = "java"
	report.PackageManager = "maven"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryConventional(queryDependencies)
	for _, capture := range captures {
		var groupId, artifactId, version string

		dependencyNode := capture["param_dependency"]
		for i := 0; i < dependencyNode.ChildCount(); i++ {
			child := dependencyNode.Child(i)

			if child.Type() != "element" {
				continue
			}

			tag := ""
			tagContent := ""

			for j := 0; j < child.ChildCount(); j++ {
				elementChild := child.Child(j)

				if elementChild.Type() == "start_tag" {
					tag = elementChild.Child(0).Content()
				}

				if elementChild.Type() == "text" {
					tagContent = elementChild.Content()
				}
			}

			switch tag {
			case "groupId":
				groupId = tagContent
			case "artifactId":
				artifactId = tagContent
			case "version":
				version = tagContent
			}
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    artifactId,
			Group:   groupId,
			Version: stringutil.StripQuotes(version),
			Line:    int64(capture["param_dependency"].StartLineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}

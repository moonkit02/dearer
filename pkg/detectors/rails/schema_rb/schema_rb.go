package schema_rb

import (
	"strings"

	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	parserschema "github.com/moonkit02/dearer/pkg/parser/schema"
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/report/schema"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/moonkit02/dearer/pkg/util/pluralize"
)

var (
	language = ruby.GetLanguage()

	rubyDatabaseSchemaQuery = parser.QueryMustCompile(language, `
		(call
			method: (identifier) @table_method
			arguments: (argument_list . (string) @table_name)
			block: (do_block
				parameters: (block_parameters (identifier) @block_param)
				(call
					receiver: (_) @receiver
					method: (_) @type
					arguments: (argument_list . (string) @column_name))
				(#eq @receiver @block_param))
			(#eq @table_method "create_table")) @rule
	`)
)

func ExtractFromDatabaseSchema(
	idGenerator nodeid.Generator,
	file *file.FileInfo,
	report report.Report,
) error {
	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return err
	}
	defer tree.Close()

	uuidHolder := parserschema.NewUUIDHolder()

	err = tree.Query(rubyDatabaseSchemaQuery, func(captures parser.Captures) error {
		tableNode := captures["table_name"]
		tableName := stripQuotes(tableNode.Content())
		columnNode := captures["column_name"]
		columnName := stripQuotes(columnNode.Content())
		columnTypeNode := captures["type"]
		columnType := columnTypeNode.Content()

		if columnType == "index" {
			return nil
		}

		ruleNode := captures["rule"]

		objectUUID := uuidHolder.Assign(tableNode.ID(), idGenerator)
		fieldUUID := uuidHolder.Assign(columnTypeNode.ID(), idGenerator)

		normalizedObjectName := ""
		normalizedFieldName := ""
		normalizedObjectName = pluralize.Singular(strings.ToLower(tableName))
		normalizedFieldName = pluralize.Singular(strings.ToLower(columnName))

		currentSchema := schema.Schema{
			ObjectName:           tableName,
			ObjectUUID:           objectUUID,
			FieldName:            columnName,
			FieldUUID:            fieldUUID,
			FieldType:            columnType,
			SimpleFieldType:      convertToSimpleType(columnType),
			NormalizedObjectName: normalizedObjectName,
			NormalizedFieldName:  normalizedFieldName,
		}

		if report.SchemaGroupShouldClose(tableName) {
			report.SchemaGroupEnd(idGenerator)
		}

		if !report.SchemaGroupIsOpen() {
			source := tableNode.Source(false)
			report.SchemaGroupBegin(
				detectors.DetectorSchemaRb,
				tableNode,
				currentSchema,
				&source,
				ruleNode,
			)
		}
		source := columnNode.Source(false)
		report.SchemaGroupAddItem(
			columnNode,
			currentSchema,
			&source,
		)

		return nil
	})

	report.SchemaGroupEnd(idGenerator)

	return err
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}

func convertToSimpleType(value string) string {
	switch value {

	case "boolean":
		return schema.SimpleTypeBool
	case "date":
		return schema.SimpleTypeDate
	case "timestamp":
		return schema.SimpleTypeDate
	case "datetime":
		return schema.SimpleTypeDate
	case "float":
		return schema.SimpleTypeNumber
	case "integer":
		return schema.SimpleTypeNumber
	case "bigint":
		return schema.SimpleTypeNumber
	case "binary":
		return schema.SimpleTypeString
	case "string":
		return schema.SimpleTypeString
	case "text":
		return schema.SimpleTypeString
	default:
		return schema.SimpleTypeObject
	}
}

package datatypes_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/report/output/dataflow"
	"github.com/moonkit02/dearer/pkg/report/output/dataflow/types"
	"github.com/moonkit02/dearer/pkg/report/output/detectors"
	outputtypes "github.com/moonkit02/dearer/pkg/report/output/types"
	globaltypes "github.com/moonkit02/dearer/pkg/types"
)

func TestDataflowDataType(t *testing.T) {
	config := settings.Config{
		Rules: map[string]*settings.Rule{
			"logger_leak": {
				Stored: true,
			},
		},
	}

	testCases := []struct {
		Name        string
		FileContent string
		Config      settings.Config
		Want        []types.Datatype
	}{
		{
			Name:        "single detection",
			Config:      config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}`,
			Want: []types.Datatype{
				{
					Name: "Username",
					Detectors: []types.DatatypeDetector{
						{
							Name: "ruby",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.rb",
									FullFilename:      "./users.rb",
									StartLineNumber:   25,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
					},
				},
			},
		},
		{
			Name:        "single detection - no classification",
			Config:      config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name"}}`,
			Want:        []types.Datatype{},
		},
		{
			Name:   "single detection - duplicates",
			Config: config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}
{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}`,
			Want: []types.Datatype{
				{
					Name: "Username",
					Detectors: []types.DatatypeDetector{
						{
							Name: "ruby",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.rb",
									FullFilename:      "./users.rb",
									StartLineNumber:   25,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
					},
				},
			},
		},
		{
			Name:   "single detection - with wierd data in report",
			Config: config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}
{"user": true }`,
			Want: []types.Datatype{
				{
					Name: "Username",
					Detectors: []types.DatatypeDetector{
						{
							Name: "ruby",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.rb",
									FullFilename:      "./users.rb",
									StartLineNumber:   25,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
					},
				},
			},
		},
		{
			Name:   "multiple detections - with same object name - deterministic output",
			Config: config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}
{"type": "schema_classified", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12, "start_line_number": 12, "end_line_number": 12, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}`,
			Want: []types.Datatype{
				{
					Name: "Username",
					Detectors: []types.DatatypeDetector{
						{
							Name: "csharp",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.cs",
									FullFilename:      "./users.cs",
									StartLineNumber:   12,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
						{
							Name: "ruby",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.rb",
									FullFilename:      "./users.rb",
									StartLineNumber:   25,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
					},
				},
			},
		},
		{
			Name:   "multiple detections - with different names - deterministic output",
			Config: config,
			FileContent: `{"type": "schema_classified", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25, "start_line_number": 25, "end_line_number": 25, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username"} ,"decision":{"state": "valid"}}}}
{"type": "schema_classified", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12, "start_line_number": 12, "end_line_number": 12, "end_column_number": 30, "start_column_number": 20}, "value": {"field_name": "address", "classification": {"data_type": {"name": "Physical Address"} ,"decision":{"state": "valid"}}}}`,
			Want: []types.Datatype{
				{
					Name: "Physical Address",
					Detectors: []types.DatatypeDetector{
						{
							Name: "csharp",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.cs",
									FullFilename:      "./users.cs",
									StartLineNumber:   12,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "address",
								},
							},
						},
					},
				},
				{
					Name: "Username",
					Detectors: []types.DatatypeDetector{
						{
							Name: "ruby",
							Locations: []types.DatatypeLocation{
								{
									Filename:          "./users.rb",
									FullFilename:      "./users.rb",
									StartLineNumber:   25,
									StartColumnNumber: 20,
									EndColumnNumber:   30,
									FieldName:         "User_name",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			file, err := os.CreateTemp("", "*test.jsonlines")
			if err != nil {
				t.Fatalf("failed to create tmp file for report %s", err)
				return
			}
			defer os.Remove(file.Name())
			_, err = file.Write([]byte(test.FileContent))
			if err != nil {
				t.Fatalf("failed to write to tmp file %s", err)
				return
			}
			file.Close()

			output := &outputtypes.ReportData{}
			if err = detectors.AddReportData(output, globaltypes.Report{
				Path: file.Name(),
			}, test.Config); err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			if err = dataflow.AddReportData(output, test.Config, false, true); err != nil {
				t.Fatalf("failed to get dataflow output %s", err)
				return
			}

			assert.Equal(t, test.Want, output.Dataflow.Datatypes)
		})
	}
}

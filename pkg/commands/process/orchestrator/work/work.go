package work

import (
	"github.com/moonkit02/dearer/pkg/commands/process/filelist/files"
	"github.com/moonkit02/dearer/pkg/scanner/stats"
)

type InitializeResponse struct {
	Error string
}

type ProcessResponse struct {
	FileStats *stats.FileStats
	Error     string
}

type Repository struct {
	Dir               string
	PreviousCommitSHA string
	CommitSHA         string
}

type ProcessRequest struct {
	Repository
	File       files.File
	ReportPath string
}

var RouteInitialize = "/initialize"
var RouteProcess = "/process"
var RouteReduceMemory = "/reduce_memory"

package models

import "strings"

type ModelVersionStage string

func (s ModelVersionStage) String() string {
	return string(s)
}

const (
	ModelVersionStageNone       = "None"
	ModelVersionStageStaging    = "Staging"
	ModelVersionStageProduction = "Production"
	ModelVersionStageArchived   = "Archived"
)

var CanonicalMapping = map[string]ModelVersionStage{
	strings.ToLower(ModelVersionStageNone):       ModelVersionStageNone,
	strings.ToLower(ModelVersionStageStaging):    ModelVersionStageStaging,
	strings.ToLower(ModelVersionStageProduction): ModelVersionStageProduction,
	strings.ToLower(ModelVersionStageArchived):   ModelVersionStageArchived,
}

var DefaultStagesForGetLatestVersions = map[string]ModelVersionStage{
	strings.ToLower(ModelVersionStageStaging):    ModelVersionStageStaging,
	strings.ToLower(ModelVersionStageProduction): ModelVersionStageProduction,
}

func AllModelVersionStages() string {
	pairs := make([]string, 0, len(CanonicalMapping))

	for _, v := range CanonicalMapping {
		pairs = append(pairs, v.String())
	}

	return strings.Join(pairs, ",")
}

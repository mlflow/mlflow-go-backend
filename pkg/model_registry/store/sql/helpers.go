package sql

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/mlflow/mlflow-go-backend/pkg/entities"
)

// GetNextVersion returns the next version number for a given registered model.
func GetNextVersion(sqlRegisteredModel *entities.RegisteredModel) int32 {
	if len(sqlRegisteredModel.Versions) > 0 {
		maxVersion := sqlRegisteredModel.Versions[0].Version
		for _, mv := range sqlRegisteredModel.Versions {
			if mv.Version > maxVersion {
				maxVersion = mv.Version
			}
		}

		return maxVersion + 1
	}

	return 1
}

type ParsedModelURI struct {
	Name    string
	Version string
	Stage   string
	Alias   string
}

//nolint:cyclop,err113,mnd,wrapcheck
func ParseModelURI(uri string) (*ParsedModelURI, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "models" {
		return nil, errors.New("invalid model URI scheme")
	}

	path := strings.TrimLeft(parsed.Path, "/")
	if path == "" {
		return nil, errors.New("invalid model URI path")
	}

	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return nil, errors.New("invalid model URI format")
	}

	if len(parts) == 2 {
		name, suffix := parts[0], parts[1]
		if suffix == "" {
			return nil, errors.New("invalid model URI suffix")
		}

		if _, err := strconv.Atoi(suffix); err == nil {
			// The suffix is a specific version
			return &ParsedModelURI{Name: name, Version: suffix}, nil
		} else if strings.EqualFold(suffix, "latest") {
			// The suffix is "latest"
			return &ParsedModelURI{Name: name}, nil
		}

		// The suffix is a specific stage
		return &ParsedModelURI{Name: name, Stage: suffix}, nil
	}

	// The URI is an alias URI
	aliasParts := strings.SplitN(parts[0], "@", 2)
	if len(aliasParts) != 2 || aliasParts[1] == "" {
		return nil, errors.New("invalid model alias format")
	}

	return &ParsedModelURI{Name: aliasParts[0], Alias: aliasParts[1]}, nil
}

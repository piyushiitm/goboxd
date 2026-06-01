package validator

import (
	"errors"
	"strings"

	"github.com/piyushiitm/goboxd/internal/models"
)

var (
	ErrMissingLanguage         = errors.New("language is required")
	ErrMissingSource           = errors.New("source is required")
	ErrMissingTests            = errors.New("at least one test is required")
	ErrInvalidSourceFilename   = errors.New("source_filename must be a single path component")
	ErrInvalidArtifactFilename = errors.New("artifact_filename must be a single path component")
	ErrUnknownLanguage         = errors.New("language must match a configured language id")
	ErrInvalidBuildFlag        = errors.New("invalid build flag")
	ErrInvalidRunFlag          = errors.New("invalid run flag")
)

func isValidFilename(name string) bool {

	if name == "" {
		return true
	}

	if strings.HasPrefix(name, ".") {
		return false
	}

	if strings.Contains(name, "/") {
		return false
	}

	if strings.Contains(name, "\\") {
		return false
	}

	return true
}

func Validate(req models.RunRequest) error {

	if strings.TrimSpace(req.Language) == "" {
		return ErrMissingLanguage
	}

	if strings.TrimSpace(req.Source) == "" {
		return ErrMissingSource
	}

	if len(req.Tests) == 0 {
		return ErrMissingTests
	}
	if !isValidFilename(req.SourceFilename) {
		return ErrInvalidSourceFilename
	}

	if !isValidFilename(req.ArtifactFilename) {
		return ErrInvalidArtifactFilename
	}
	return nil
}

func ValidateFlags(
	flags []string,
	allowed []string,
) bool {

	allowedSet := make(map[string]struct{})

	for _, flag := range allowed {
		allowedSet[flag] = struct{}{}
	}

	for _, flag := range flags {

		if _, ok := allowedSet[flag]; !ok {
			return false
		}
	}

	return true
}

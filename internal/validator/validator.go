package validator

import (
	"errors"
	"strings"

	"github.com/piyushiitm/goboxd/internal/models"
)

var (
	ErrMissingLanguage = errors.New("language is required")
	ErrMissingSource   = errors.New("source is required")
	ErrMissingTests    = errors.New("at least one test is required")
)

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

	return nil
}

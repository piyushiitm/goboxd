package validator

import (
	"testing"

	"github.com/piyushiitm/goboxd/internal/models"
)

func TestMissingLanguage(t *testing.T) {

	req := models.RunRequest{}

	err := Validate(req)

	if err != ErrMissingLanguage {
		t.Fatal("expected ErrMissingLanguage")
	}
}

func TestMissingSource(t *testing.T) {

	req := models.RunRequest{
		Language: "py3",
	}

	err := Validate(req)

	if err != ErrMissingSource {
		t.Fatal("expected ErrMissingSource")
	}
}

func TestMissingTests(t *testing.T) {

	req := models.RunRequest{
		Language: "py3",
		Source:   "print('hello')",
	}

	err := Validate(req)

	if err != ErrMissingTests {
		t.Fatal("expected ErrMissingTests")
	}
}

func TestValidRequest(t *testing.T) {

	req := models.RunRequest{
		Language: "py3",
		Source:   "print('hello')",
		Tests: []models.TestCase{
			{
				ExpectedStdout: "hello",
			},
		},
	}

	err := Validate(req)

	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidSourceFilename(t *testing.T) {

	req := models.RunRequest{
		Language:       "py3",
		Source:         "print('hello')",
		SourceFilename: "../hack.py",
		Tests: []models.TestCase{
			{
				ExpectedStdout: "hello",
			},
		},
	}

	err := Validate(req)

	if err != ErrInvalidSourceFilename {
		t.Fatal("expected ErrInvalidSourceFilename")
	}
}

func TestInvalidArtifactFilename(t *testing.T) {

	req := models.RunRequest{
		Language:         "py3",
		Source:           "print('hello')",
		ArtifactFilename: "../hack",
		Tests: []models.TestCase{
			{
				ExpectedStdout: "hello",
			},
		},
	}

	err := Validate(req)

	if err != ErrInvalidArtifactFilename {
		t.Fatal("expected ErrInvalidArtifactFilename")
	}
}

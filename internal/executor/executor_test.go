package executor

import "testing"

func TestReplacePlaceholders(t *testing.T) {

	command := []string{
		"python3",
		"{source}",
	}

	result := ReplacePlaceholders(
		command,
		"source.py",
		"program",
	)

	if result[0] != "python3" {
		t.Fail()
	}

	if result[1] != "source.py" {
		t.Fail()
	}
}

func TestUnsupportedLanguage(t *testing.T) {

	_, err := Execute(
		"rust",
		"fn main() {}",
	)

	if err == nil {
		t.Fail()
	}
}

func TestReplaceArtifactPlaceholder(t *testing.T) {

	command := []string{
		"{artifact}",
	}

	result := ReplacePlaceholders(
		command,
		"source.py",
		"program",
	)

	if result[0] != "program" {
		t.Fail()
	}
}

func TestPythonExecution(t *testing.T) {

	result, err := Execute(
		"python",
		"print('hello')",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Stdout != "hello\n" {
		t.Fatalf("unexpected output: %s", result.Stdout)
	}
}

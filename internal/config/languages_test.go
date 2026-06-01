package config

import "testing"

func TestLoadLanguages(t *testing.T) {

	languages, err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	if len(languages) == 0 {
		t.Fatal("no languages loaded")
	}
}

func TestPythonLanguageExists(t *testing.T) {

	languages, err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := languages["py3"]; !ok {
		t.Fatal("py3 language missing")
	}
}

func TestCppLanguageExists(t *testing.T) {

	languages, err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := languages["cpp"]; !ok {
		t.Fatal("cpp language missing")
	}
}

func TestJavaLanguageExists(t *testing.T) {

	languages, err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := languages["java"]; !ok {
		t.Fatal("java language missing")
	}
}

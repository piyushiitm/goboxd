package executor

import (
	"testing"

	"github.com/piyushiitm/goboxd/internal/config"
	"github.com/piyushiitm/goboxd/internal/models"
)

func TestReplaceSourcePlaceholder(t *testing.T) {

	command := []string{
		"python3",
		"{source}",
	}

	result := ReplacePlaceholders(
		command,
		"source.py",
		"program",
		"workspace",
	)

	if result[0] != "python3" {
		t.Fatal("command changed unexpectedly")
	}

	if result[1] != "source.py" {
		t.Fatal("source placeholder not replaced")
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
		"workspace",
	)

	if result[0] != "program" {
		t.Fatal("artifact placeholder not replaced")
	}
}

func TestReplaceWorkspacePlaceholder(t *testing.T) {

	command := []string{
		"{workspace}",
	}

	result := ReplacePlaceholders(
		command,
		"source.py",
		"program",
		"workspace123",
	)

	if result[0] != "workspace123" {
		t.Fatal("workspace placeholder not replaced")
	}
}

func TestReplaceAllPlaceholders(t *testing.T) {

	command := []string{
		"{workspace}",
		"{source}",
		"{artifact}",
	}

	result := ReplacePlaceholders(
		command,
		"main.cpp",
		"app",
		"workspace123",
	)

	if result[0] != "workspace123" {
		t.Fatal("workspace replacement failed")
	}

	if result[1] != "main.cpp" {
		t.Fatal("source replacement failed")
	}

	if result[2] != "app" {
		t.Fatal("artifact replacement failed")
	}
}

func TestResolveLimitsNilOverride(t *testing.T) {

	defaults := config.Limits{
		WallTimeS:    9,
		MemoryKB:     102400,
		MaxProcesses: 100,
	}

	result := ResolveLimits(defaults, nil)

	if result.WallTimeS != 9 {
		t.Fatal("wall time changed unexpectedly")
	}

	if result.MemoryKB != 102400 {
		t.Fatal("memory changed unexpectedly")
	}

	if result.MaxProcesses != 100 {
		t.Fatal("process count changed unexpectedly")
	}
}

func TestResolveLimitsWallTimeOverride(t *testing.T) {

	defaults := config.Limits{
		WallTimeS:    9,
		MemoryKB:     102400,
		MaxProcesses: 100,
	}

	override := &models.Limits{
		WallTimeS: 3,
	}

	result := ResolveLimits(defaults, override)

	if result.WallTimeS != 3 {
		t.Fatal("wall time override not applied")
	}

	if result.MemoryKB != 102400 {
		t.Fatal("memory should remain unchanged")
	}

	if result.MaxProcesses != 100 {
		t.Fatal("process count should remain unchanged")
	}
}

func TestResolveLimitsMemoryOverride(t *testing.T) {

	defaults := config.Limits{
		WallTimeS:    9,
		MemoryKB:     102400,
		MaxProcesses: 100,
	}

	override := &models.Limits{
		MemoryKB: 204800,
	}

	result := ResolveLimits(defaults, override)

	if result.MemoryKB != 204800 {
		t.Fatal("memory override not applied")
	}
}

func TestResolveLimitsProcessOverride(t *testing.T) {

	defaults := config.Limits{
		WallTimeS:    9,
		MemoryKB:     102400,
		MaxProcesses: 100,
	}

	override := &models.Limits{
		MaxProcesses: 50,
	}

	result := ResolveLimits(defaults, override)

	if result.MaxProcesses != 50 {
		t.Fatal("process override not applied")
	}
}

func TestResolveLimitsAllOverrides(t *testing.T) {

	defaults := config.Limits{
		WallTimeS:    9,
		MemoryKB:     102400,
		MaxProcesses: 100,
	}

	override := &models.Limits{
		WallTimeS:    5,
		MemoryKB:     204800,
		MaxProcesses: 50,
	}

	result := ResolveLimits(defaults, override)

	if result.WallTimeS != 5 {
		t.Fatal("wall time override failed")
	}

	if result.MemoryKB != 204800 {
		t.Fatal("memory override failed")
	}

	if result.MaxProcesses != 50 {
		t.Fatal("process override failed")
	}
}

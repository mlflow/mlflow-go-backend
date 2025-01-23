//go:build mage

//nolint:wrapcheck
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/magefile/mage/sh"
)

var errNoVersionInChangelog = errors.New("no version found in changelog")

const writeFilePermission = 0o600

// Update the pyproject.toml version based on the changelog.
func updateVersionFromChangelog(changelogPath, pyprojectPath string) error {
	// Define the regex pattern to match a version line, e.g., ## [0.1.0]
	versionPattern := regexp.MustCompile(`\#\# \[(\d+\.\d+.\d+)\]`)

	// Open the changelog file
	file, err := os.Open(changelogPath)
	if err != nil {
		return fmt.Errorf("failed to open changelog: %w", err)
	}
	defer file.Close()

	// Find the first matching version
	var version string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		matches := versionPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			version = matches[1]

			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading changelog: %w", err)
	}

	if version == "" {
		return errNoVersionInChangelog
	}

	// Read the pyproject.toml file
	pyproject, err := os.ReadFile(pyprojectPath)
	if err != nil {
		return fmt.Errorf("failed to read pyproject.toml: %w", err)
	}

	// Update the version in the pyproject.toml file
	lines := strings.Split(string(pyproject), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "version =") {
			lines[i] = fmt.Sprintf("version = \"%s\"", version)

			break
		}
	}

	// Write the updated pyproject.toml back
	updatedPyproject := strings.Join(lines, "\n")

	err = os.WriteFile(pyprojectPath, []byte(updatedPyproject), writeFilePermission)
	if err != nil {
		return fmt.Errorf("failed to write updated pyproject.toml: %w", err)
	}

	log.Printf("Updated version in pyproject.toml to %s\n", version)

	return nil
}

// Build a Python wheel.
func Build(goos, goarch string) error {
	if err := updateVersionFromChangelog("CHANGELOG.md", "pyproject.toml"); err != nil {
		return err
	}

	if err := sh.RunWithV(map[string]string{
		"TARGET_GOOS":   goos,
		"TARGET_GOARCH": goarch,
	}, "uvx", "--from", "build[uv]", "pyproject-build", "--installer", "uv"); err != nil {
		return err
	}

	return nil
}

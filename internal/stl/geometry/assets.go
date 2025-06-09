package geometry

import (
	"embed"
	"fmt"
	"os"

	"github.com/github/gh-skyline/internal/errors"
)

//go:embed assets/*
var embeddedAssets embed.FS

// writeTempFont writes the embedded font to a temporary file and returns its path.
// The caller is responsible for cleaning up the temporary file.
func writeTempFont(fontName string) (string, func(), error) {
	fontBytes, err := embeddedAssets.ReadFile("assets/" + fontName)
	if err != nil {
		return "", nil, errors.New(errors.IOError, "failed to read embedded font", err)
	}

	// Create temp file with .ttf extension to ensure proper font loading
	tmpFile, err := os.CreateTemp("", "skyline-font-*.ttf")
	if err != nil {
		return "", nil, errors.New(errors.IOError, "failed to create temp font file", err)
	}

	if _, err := tmpFile.Write(fontBytes); err != nil {
		closeErr := tmpFile.Close()
		removeErr := os.Remove(tmpFile.Name())
		return "", nil, errors.New(errors.IOError, "failed to write font to temp file", fmt.Errorf("%w; close error: %v; remove error: %v", err, closeErr, removeErr))
	}
	if err := tmpFile.Close(); err != nil {
		removeErr := os.Remove(tmpFile.Name())
		return "", nil, errors.New(errors.IOError, "failed to close temp font file", fmt.Errorf("%w; remove error: %v", err, removeErr))
	}

	cleanup := func() {
		_ = os.Remove(tmpFile.Name()) // Ignore cleanup errors in defer
	}

	return tmpFile.Name(), cleanup, nil
}

// getEmbeddedImage returns a temporary file path for the embedded image.
// The caller is responsible for cleaning up the temporary file.
func getEmbeddedImage() (string, func(), error) {
	imgBytes, err := embeddedAssets.ReadFile("assets/invertocat.png")
	if err != nil {
		return "", nil, errors.New(errors.IOError, "failed to read embedded image", err)
	}

	tmpFile, err := os.CreateTemp("", "skyline-img-*.png")
	if err != nil {
		return "", nil, errors.New(errors.IOError, "failed to create temp image file", err)
	}

	if _, err := tmpFile.Write(imgBytes); err != nil {
		closeErr := tmpFile.Close()
		removeErr := os.Remove(tmpFile.Name())
		return "", nil, errors.New(errors.IOError, "failed to write image to temp file", fmt.Errorf("%w; close error: %v; remove error: %v", err, closeErr, removeErr))
	}
	if err := tmpFile.Close(); err != nil {
		removeErr := os.Remove(tmpFile.Name())
		return "", nil, errors.New(errors.IOError, "failed to close temp image file", fmt.Errorf("%w; remove error: %v", err, removeErr))
	}

	cleanup := func() {
		_ = os.Remove(tmpFile.Name()) // Ignore cleanup errors in defer
	}

	return tmpFile.Name(), cleanup, nil
}

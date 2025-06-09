package geometry

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWriteTempFont verifies temporary font file creation and cleanup
func TestWriteTempFont(t *testing.T) {
	t.Run("verify valid font extraction", func(t *testing.T) {
		fontPath, cleanup, err := writeTempFont("monasans-medium.ttf")
		if err != nil {
			t.Fatalf("writeTempFont failed: %v", err)
		}
		defer cleanup()

		// Verify file exists and has content
		content, err := os.ReadFile(fontPath)
		if err != nil {
			t.Errorf("Failed to read temp font file: %v", err)
		}
		if len(content) == 0 {
			t.Error("Temp font file is empty")
		}

		// Verify file extension
		if filepath.Ext(fontPath) != ".ttf" {
			t.Errorf("Expected .ttf extension, got %s", filepath.Ext(fontPath))
		}

		// Verify cleanup works
		cleanup()
		if _, err := os.Stat(fontPath); !os.IsNotExist(err) {
			t.Error("Temp font file not cleaned up properly")
		}
	})

	t.Run("verify nonexistent font handling", func(t *testing.T) {
		_, cleanup, err := writeTempFont("nonexistent.ttf")
		if err == nil {
			defer cleanup()
			t.Error("Expected error for nonexistent font")
		}
	})
}

// TestGetEmbeddedImage verifies temporary image file creation and cleanup
func TestGetEmbeddedImage(t *testing.T) {
	t.Run("verify valid image extraction", func(t *testing.T) {
		imagePath, cleanup, err := getEmbeddedImage()
		if err != nil {
			t.Fatalf("getEmbeddedImage failed: %v", err)
		}
		defer cleanup()

		// Verify file exists and has content
		content, err := os.ReadFile(imagePath)
		if err != nil {
			t.Errorf("Failed to read temp image file: %v", err)
		}
		if len(content) == 0 {
			t.Error("Temp image file is empty")
		}

		// Verify file extension
		if filepath.Ext(imagePath) != ".png" {
			t.Errorf("Expected .png extension, got %s", filepath.Ext(imagePath))
		}

		// Verify cleanup works
		cleanup()
		if _, err := os.Stat(imagePath); !os.IsNotExist(err) {
			t.Error("Temp image file not cleaned up properly")
		}
	})

	// Test embedded filesystem access
	t.Run("verify embedded filesystem access", func(t *testing.T) {
		// Try to read the embedded image directly
		_, err := embeddedAssets.ReadFile("assets/invertocat.png")
		if err != nil {
			t.Errorf("Failed to access embedded image: %v", err)
		}
	})
}

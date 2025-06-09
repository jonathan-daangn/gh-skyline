package ascii

import "testing"

func TestBlockConstants(t *testing.T) {
	// Test that block characters are different
	blocks := []rune{
		EmptyBlock,
		FoundationLow,
		FoundationMed,
		FoundationHigh,
		MiddleLow,
		MiddleMed,
		MiddleHigh,
		TopLow,
		TopMed,
		TopHigh,
	}

	for i := 0; i < len(blocks); i++ {
		for j := i + 1; j < len(blocks); j++ {
			if i == j {
				continue
			}
			// Some blocks intentionally use the same character
			if isIntentionallySameChar(blocks[i], blocks[j]) {
				continue
			}
			if blocks[i] == blocks[j] {
				t.Errorf("Block characters at positions %d and %d are the same: %c", i, j, blocks[i])
			}
		}
	}
}

func TestThresholds(t *testing.T) {
	if LowThreshold >= MediumThreshold {
		t.Error("LowThreshold should be less than MediumThreshold")
	}
	if LowThreshold <= 0 || LowThreshold >= 1 {
		t.Error("LowThreshold should be between 0 and 1")
	}
	if MediumThreshold <= 0 || MediumThreshold >= 1 {
		t.Error("MediumThreshold should be between 0 and 1")
	}
}

// Helper function to check if two blocks are intentionally the same character
func isIntentionallySameChar(a, b rune) bool {
	// Some blocks use the same character by design
	intentionallySame := map[rune][]rune{
		'░': {FoundationLow, MiddleLow},
		'▒': {FoundationMed, MiddleMed},
		'▓': {FoundationHigh, MiddleHigh},
	}

	if same, exists := intentionallySame[a]; exists {
		for _, r := range same {
			if r == b {
				return true
			}
		}
	}
	return false
}

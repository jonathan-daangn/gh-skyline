package ascii

// Block character sets for different contribution levels.
// The ASCII art uses different characters depending on the position (foundation, middle, top)
// and intensity (low, medium, high) of contributions.
const (
	// Basic blocks
	EmptyBlock  = ' ' // Represents days with no contributions
	FutureBlock = '.' // Represents future dates

	// Foundation blocks (bottom layer)
	FoundationLow  = '░' // 1-33% intensity
	FoundationMed  = '▒' // 34-66% intensity
	FoundationHigh = '▓' // 67-100% intensity

	// Middle blocks (intermediate layers)
	MiddleLow  = '░'
	MiddleMed  = '▒'
	MiddleHigh = '▓'

	// Top blocks (highest layer, using special characters for visual distinction)
	TopLow  = '╻' // Lower intensity peak
	TopMed  = '┃' // Medium intensity peak
	TopHigh = '╽' // High intensity peak
)

// Contribution level thresholds as percentages of the maximum contribution count
const (
	LowThreshold    = 0.33 // 33% of max contributions
	MediumThreshold = 0.66 // 66% of max contributions
)

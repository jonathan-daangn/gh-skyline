package geometry

import (
	"math"
	"testing"

	"github.com/github/gh-skyline/internal/types"
)

const (
	epsilon = 0.0001 // Tolerance for floating-point comparisons
)

// TestNormalizeContribution verifies contribution normalization logic
func TestNormalizeContribution(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		maxCount int
		want     float64
	}{
		{"zero contributions", 0, 10, 0},
		{"zero max count", 5, 0, MinHeight},
		{"negative max count", 5, -1, MinHeight},
		{"full scale", 100, 100, MaxHeight},
		{"half scale", 25, 100, MinHeight + (MaxHeight-MinHeight)*0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeContribution(tt.count, tt.maxCount)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("NormalizeContribution(%v, %v) = %v, want %v", tt.count, tt.maxCount, got, tt.want)
			}
		})
	}
}

// TestCreateContributionGeometry verifies contribution geometry generation
func TestCreateContributionGeometry(t *testing.T) {
	tests := []struct {
		name        string
		contribs    [][]types.ContributionDay
		yearIndex   int
		maxContrib  int
		wantErr     bool
		triangleLen int
	}{
		{
			name: "empty contributions",
			contribs: [][]types.ContributionDay{
				{{ContributionCount: 0, Date: "2023-01-01"}},
			},
			yearIndex:   0,
			maxContrib:  10,
			wantErr:     false,
			triangleLen: 0,
		},
		{
			name: "single contribution",
			contribs: [][]types.ContributionDay{
				{{ContributionCount: 5, Date: "2023-01-01"}},
			},
			yearIndex:   0,
			maxContrib:  10,
			wantErr:     false,
			triangleLen: 12, // One column = 12 triangles
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triangles, err := CreateContributionGeometry(tt.contribs, tt.yearIndex, tt.maxContrib)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContributionGeometry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(triangles) != tt.triangleLen {
				t.Errorf("CreateContributionGeometry() got %v triangles, want %v", len(triangles), tt.triangleLen)
			}
		})
	}
}

// TestCalculateMultiYearDimensions verifies dimension calculations
func TestCalculateMultiYearDimensions(t *testing.T) {
	tests := []struct {
		name      string
		yearCount int
		wantW     float64
		wantD     float64
	}{
		{"single year", 1, float64(GridSize)*CellSize + 4*CellSize, 7*CellSize + 4*CellSize},
		{"multiple years", 3, float64(GridSize)*CellSize + 4*CellSize, 21*CellSize + 4*CellSize},
		{"zero years", 0, float64(GridSize)*CellSize + 4*CellSize, 4 * CellSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotD := CalculateMultiYearDimensions(tt.yearCount)
			if math.Abs(gotW-tt.wantW) > epsilon {
				t.Errorf("CalculateMultiYearDimensions() width = %v, want %v", gotW, tt.wantW)
			}
			if math.Abs(gotD-tt.wantD) > epsilon {
				t.Errorf("CalculateMultiYearDimensions() depth = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

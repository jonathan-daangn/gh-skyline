package geometry

import (
	"math"
	"testing"

	"github.com/github/gh-skyline/internal/types"
)

// Epsilon is declared in geometry_test.go

// TestVectorOperations verifies vector math operations.
func TestVectorOperations(t *testing.T) {
	t.Run("verify vector normalization", func(t *testing.T) {
		v := types.Point3D{X: 3, Y: 4, Z: 0}
		normalized := normalizeVector(v)
		magnitude := math.Sqrt(float64(
			normalized.X*normalized.X +
				normalized.Y*normalized.Y +
				normalized.Z*normalized.Z))

		if math.Abs(magnitude-1.0) > epsilon {
			t.Errorf("Expected normalized vector magnitude 1.0, got %f", magnitude)
		}
	})
}

func TestValidateVector(t *testing.T) {
	tests := []struct {
		name    string
		vector  types.Point3D
		wantErr bool
	}{
		{
			name:    "valid vector",
			vector:  types.Point3D{X: 1, Y: 2, Z: 3},
			wantErr: false,
		},
		{
			name:    "invalid vector with NaN",
			vector:  types.Point3D{X: math.NaN(), Y: 2, Z: 3},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVector(tt.vector)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateNormal(t *testing.T) {
	tests := []struct {
		name    string
		p1      types.Point3D
		p2      types.Point3D
		p3      types.Point3D
		want    types.Point3D
		wantErr bool
	}{
		{
			name:    "valid triangle",
			p1:      types.Point3D{X: 0, Y: 0, Z: 0},
			p2:      types.Point3D{X: 1, Y: 0, Z: 0},
			p3:      types.Point3D{X: 0, Y: 1, Z: 0},
			want:    types.Point3D{X: 0, Y: 0, Z: 1},
			wantErr: false,
		},
		{
			name:    "degenerate triangle",
			p1:      types.Point3D{X: 0, Y: 0, Z: 0},
			p2:      types.Point3D{X: 0, Y: 0, Z: 0},
			p3:      types.Point3D{X: 0, Y: 0, Z: 0},
			want:    types.Point3D{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateNormal(tt.p1, tt.p2, tt.p3)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateNormal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !vectorsEqual(got, tt.want) {
				t.Errorf("calculateNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZeroVector(t *testing.T) {
	tests := []struct {
		name string
		v    types.Point3D
		want bool
	}{
		{
			name: "zero vector",
			v:    types.Point3D{X: 0, Y: 0, Z: 0},
			want: true,
		},
		{
			name: "non-zero vector",
			v:    types.Point3D{X: 1, Y: 1, Z: 1},
			want: false,
		},
		{
			name: "near-zero vector",
			v:    types.Point3D{X: 1e-11, Y: 1e-11, Z: 1e-11},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZeroVector(tt.v); got != tt.want {
				t.Errorf("isZeroVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorSubtract(t *testing.T) {
	tests := []struct {
		name string
		a    types.Point3D
		b    types.Point3D
		want types.Point3D
	}{
		{
			name: "simple subtraction",
			a:    types.Point3D{X: 1, Y: 2, Z: 3},
			b:    types.Point3D{X: 4, Y: 5, Z: 6},
			want: types.Point3D{X: -3, Y: -3, Z: -3},
		},
		{
			name: "zero vector result",
			a:    types.Point3D{X: 1, Y: 1, Z: 1},
			b:    types.Point3D{X: 1, Y: 1, Z: 1},
			want: types.Point3D{X: 0, Y: 0, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vectorSubtract(tt.a, tt.b); !vectorsEqual(got, tt.want) {
				t.Errorf("vectorSubtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorCross(t *testing.T) {
	tests := []struct {
		name string
		u    types.Point3D
		v    types.Point3D
		want types.Point3D
	}{
		{
			name: "standard cross product",
			u:    types.Point3D{X: 1, Y: 0, Z: 0},
			v:    types.Point3D{X: 0, Y: 1, Z: 0},
			want: types.Point3D{X: 0, Y: 0, Z: 1},
		},
		{
			name: "parallel vectors",
			u:    types.Point3D{X: 1, Y: 0, Z: 0},
			v:    types.Point3D{X: 2, Y: 0, Z: 0},
			want: types.Point3D{X: 0, Y: 0, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vectorCross(tt.u, tt.v); !vectorsEqual(got, tt.want) {
				t.Errorf("vectorCross() = %v, want %v", got, tt.want)
			}
		})
	}
}

// vectorsEqual helps compare two vectors within a small epsilon
func vectorsEqual(a, b types.Point3D) bool {
	return math.Abs(a.X-b.X) < epsilon &&
		math.Abs(a.Y-b.Y) < epsilon &&
		math.Abs(a.Z-b.Z) < epsilon
}

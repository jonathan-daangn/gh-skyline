// Package geometry provides 3D geometry manipulation functions for generating STL models.
package geometry

import (
	"math"

	"github.com/github/gh-skyline/internal/errors"
	"github.com/github/gh-skyline/internal/types"
)

// validateVector checks if a vector's components are valid numbers
func validateVector(v types.Point3D) error {
	if !v.IsValid() {
		return errors.New(errors.ValidationError, "vector contains invalid components", nil)
	}
	return nil
}

// calculateNormal computes the normalized normal vector for a triangle face.
// Returns a unit vector perpendicular to the plane defined by the three input points,
// or an error if the points form a degenerate triangle.
func calculateNormal(p1, p2, p3 types.Point3D) (types.Point3D, error) {
	// Validate input points
	for _, p := range []types.Point3D{p1, p2, p3} {
		if err := validateVector(p); err != nil {
			return types.Point3D{}, err
		}
	}

	u := vectorSubtract(p2, p1)
	v := vectorSubtract(p3, p1)
	normal := vectorCross(u, v)

	// Check for degenerate triangle
	if isZeroVector(normal) {
		return types.Point3D{}, errors.New(errors.ValidationError, "degenerate triangle", nil)
	}

	return normalizeVector(normal), nil
}

// isZeroVector checks if a vector has zero magnitude
func isZeroVector(v types.Point3D) bool {
	const epsilon = 1e-10
	return math.Abs(v.X) < epsilon && math.Abs(v.Y) < epsilon && math.Abs(v.Z) < epsilon
}

// vectorSubtract calculates the vector difference between two 3D points.
// Returns a vector representing the direction from point b to point a.
func vectorSubtract(a, b types.Point3D) types.Point3D {
	return types.Point3D{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// vectorCross computes the cross product of two 3D vectors.
// Returns a vector perpendicular to both input vectors.
func vectorCross(u, v types.Point3D) types.Point3D {
	return types.Point3D{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

// normalizeVector converts a vector to a unit vector (magnitude of 1).
// If the input vector has zero length, returns the original vector unchanged.
func normalizeVector(v types.Point3D) types.Point3D {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if length > 0 {
		return types.Point3D{
			X: v.X / length,
			Y: v.Y / length,
			Z: v.Z / length,
		}
	}
	return v
}

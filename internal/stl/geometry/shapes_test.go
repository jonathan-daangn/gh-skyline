package geometry

import (
	"math"
	"testing"

	"github.com/github/gh-skyline/internal/types"
)

// TestCreateCuboidBase verifies cuboid base generation functionality.
func TestCreateCuboidBase(t *testing.T) {
	t.Run("verify basic cuboid dimensions and normal vectors", func(t *testing.T) {
		// Create a cuboid with equal width and depth
		generatedTriangles, err := CreateCuboidBase(10.0, 10.0)

		if err != nil {
			t.Fatalf("CreateCuboidBase failed: %v", err)
		}

		expectedTriangleCount := 12 // 2 triangles each for: top, bottom, front, back, left, right

		if len(generatedTriangles) != expectedTriangleCount {
			t.Errorf("Expected %d triangles, got %d", expectedTriangleCount, len(generatedTriangles))
		}

		// Ensure all normal vectors are unit vectors (magnitude = 1)
		for triangleIndex, triangle := range generatedTriangles {
			normalLength := math.Sqrt(float64(
				triangle.Normal.X*triangle.Normal.X +
					triangle.Normal.Y*triangle.Normal.Y +
					triangle.Normal.Z*triangle.Normal.Z))
			if math.Abs(normalLength-1.0) > epsilon {
				t.Errorf("Triangle %d has invalid normal vector: magnitude %f", triangleIndex, normalLength)
			}
		}
	})

	// Test edge case where dimensions are zero
	t.Run("ensure non-zero output for zero dimensions", func(t *testing.T) {
		_, err := CreateCuboidBase(0.0, 0.0)
		if err == nil {
			t.Error("Expected error for zero dimensions")
		}
	})
}

// TestCreateColumn verifies column generation functionality.
func TestCreateColumn(t *testing.T) {
	t.Run("verify standard column generation", func(t *testing.T) {
		generatedTriangles, err := CreateColumn(0, 0, 10, 2)
		if err != nil {
			t.Fatalf("CreateColumn failed: %v", err)
		}
		expectedTriangleCount := 12 // 2 triangles each for front, back, left, right, top, bottom

		if len(generatedTriangles) != expectedTriangleCount {
			t.Errorf("Expected %d triangles, got %d", expectedTriangleCount, len(generatedTriangles))
		}
	})

	// Test edge case of zero height
	t.Run("verify column generation with zero height", func(t *testing.T) {
		_, err := CreateColumn(0, 0, 0, 2)
		if err == nil {
			t.Error("Expected error for zero height")
		}
	})

	t.Run("verify column coordinates", func(t *testing.T) {
		x, y, z := 1.0, 2.0, 0.0
		height, size := 5.0, 2.0
		triangles, err := CreateColumn(x, y, height, size)
		if err != nil {
			t.Fatalf("CreateColumn failed: %v", err)
		}

		// Test that coordinates are within expected bounds
		for _, tri := range triangles {
			vertices := []types.Point3D{tri.V1, tri.V2, tri.V3}
			for _, v := range vertices {
				if v.X < x-epsilon || v.X > x+size+epsilon ||
					v.Y < y-epsilon || v.Y > y+size+epsilon ||
					v.Z < z-epsilon || v.Z > z+height+epsilon {
					t.Error("Column vertex coordinates out of expected bounds")
				}
			}
		}
	})

	t.Run("verify negative dimensions", func(t *testing.T) {
		triangles, err := CreateColumn(0, 0, -1, -1)

		if err == nil {
			t.Error("Expected error for negative dimensions")
		}

		if len(triangles) != 0 {
			t.Errorf("Expected zero triangles, got %d", len(triangles))
		}
	})
}

// TestCreateQuad verifies quad creation
func TestCreateQuad(t *testing.T) {
	t.Run("verify valid quad creation", func(t *testing.T) {
		v1 := types.Point3D{X: 0, Y: 0, Z: 0}
		v2 := types.Point3D{X: 1, Y: 0, Z: 0}
		v3 := types.Point3D{X: 1, Y: 1, Z: 0}
		v4 := types.Point3D{X: 0, Y: 1, Z: 0}

		triangles, err := CreateQuad(v1, v2, v3, v4)
		if err != nil {
			t.Fatalf("CreateQuad failed: %v", err)
		}
		if len(triangles) != 2 {
			t.Errorf("Expected 2 triangles, got %d", len(triangles))
		}
	})

	t.Run("verify degenerate quad", func(t *testing.T) {
		v1 := types.Point3D{X: 0, Y: 0, Z: 0}
		v2 := types.Point3D{X: 0, Y: 0, Z: 0}
		v3 := types.Point3D{X: 0, Y: 0, Z: 0}
		v4 := types.Point3D{X: 0, Y: 0, Z: 0}

		_, err := CreateQuad(v1, v2, v3, v4)
		if err == nil {
			t.Error("Expected error for degenerate quad")
		}
	})
}

// TestCreateCube verifies cube creation
func TestCreateCube(t *testing.T) {
	t.Run("verify cube creation", func(t *testing.T) {
		triangles, err := CreateCube(0, 0, 0, 1, 1, 1)
		if err != nil {
			t.Fatalf("CreateCube failed: %v", err)
		}
		expectedTriangles := 12 // 6 faces * 2 triangles per face
		if len(triangles) != expectedTriangles {
			t.Errorf("Expected %d triangles, got %d", expectedTriangles, len(triangles))
		}
	})

	t.Run("verify zero dimensions", func(t *testing.T) {
		_, err := CreateCube(0, 0, 0, 0, 0, 0)
		if err == nil {
			t.Error("Expected error for zero dimensions")
		}
	})

	t.Run("verify normal vectors", func(t *testing.T) {
		triangles, err := CreateCube(0, 0, 0, 1, 1, 1)
		if err != nil {
			t.Fatalf("CreateCube failed: %v", err)
		}
		for i, tri := range triangles {
			normalLen := math.Sqrt(tri.Normal.X*tri.Normal.X +
				tri.Normal.Y*tri.Normal.Y +
				tri.Normal.Z*tri.Normal.Z)
			if math.Abs(normalLen-1.0) > epsilon {
				t.Errorf("Triangle %d has non-unit normal vector: length = %f", i, normalLen)
			}
		}
	})
}

// TestCreateBox verifies internal box creation functionality
func TestCreateBox(t *testing.T) {
	t.Run("verify negative dimensions", func(t *testing.T) {
		_, err := createBox(0, 0, 0, -1, -1, -1)
		if err == nil {
			t.Error("Expected error for negative dimensions")
		}
	})

	t.Run("verify vertex count", func(t *testing.T) {
		triangles, err := createBox(0, 0, 0, 1, 1, 1)
		if err != nil {
			t.Fatalf("createBox failed: %v", err)
		}
		expectedVertices := make(map[types.Point3D]bool)
		for _, tri := range triangles {
			expectedVertices[tri.V1] = true
			expectedVertices[tri.V2] = true
			expectedVertices[tri.V3] = true
		}
		if len(expectedVertices) != 8 {
			t.Errorf("Expected 8 unique vertices, got %d", len(expectedVertices))
		}
	})
}

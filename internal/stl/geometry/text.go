package geometry

import (
	"fmt"
	"image/png"
	"os"

	"github.com/fogleman/gg"
	"github.com/github/gh-skyline/internal/errors"
	"github.com/github/gh-skyline/internal/types"
)

const (
	baseWidthVoxelResolution = 2000 // Number of voxels across the skyline face
	voxelDepth               = 1.0  // Distance to come out of face

	logoScale      = 0.5   // Percent
	logoTopOffset  = -0.18 // Percent (세로 중앙)
	logoLeftOffset = 0.03  // Percent

	usernameFontSize      = 120.0
	usernameJustification = "left" // "left", "center", "right"
	usernameLeftOffset    = 0.1    // Percent

	yearFontSize      = 100.0
	yearJustification = "right" // "left", "center", "right"
	yearLeftOffset    = 0.97    // Percent

	// 추가 텍스트를 위한 상수
	additionalTextFontSize      = 150.0
	additionalTextJustification = "center" // "left", "center", "right"
	additionalTextLeftOffset    = 0.5      // Percent (오른쪽 끝에서 살짝 안쪽)
	additionalTextTopOffset     = 0.5      // Percent (세로 중앙)
)

// Create3DText generates 3D text geometry for the username and year.
func Create3DText(username string, year string, baseWidth float64, baseHeight float64, baseDepth float64, additionalText string) ([]types.Triangle, error) {
	var allTriangles []types.Triangle

	if username != "" {
		usernameTriangles, err := renderText(
			username,
			usernameJustification,
			usernameLeftOffset,
			usernameFontSize,
			baseWidth,
			baseHeight,
		)
		if err != nil {
			return nil, err
		}
		allTriangles = append(allTriangles, usernameTriangles...)
	}

	yearTriangles, err := renderText(
		year,
		yearJustification,
		yearLeftOffset,
		yearFontSize,
		baseWidth,
		baseHeight,
	)
	if err != nil {
		return nil, err
	}
	allTriangles = append(allTriangles, yearTriangles...)

	// 추가 텍스트가 있는 경우 윗면에 생성
	if additionalText != "" {
		topTriangles, err := renderTextOnTop(
			additionalText,
			additionalTextJustification,
			additionalTextLeftOffset,
			additionalTextTopOffset,
			additionalTextFontSize,
			baseWidth,
			baseDepth,
			baseHeight,
		)
		if err != nil {
			return nil, err
		}
		allTriangles = append(allTriangles, topTriangles...)
	}

	return allTriangles, nil
}

// renderText places text on the face of a skyline, offset from the left and vertically-aligned.
// The function takes the text to be displayed, offset from left, and font size.
// It returns an array of types.Triangle.
//
// Parameters:
//
//	text (string): The text to be displayed on the skyline's front face.
//	leftOffsetPercent (float64): The percentage distance from the left to start displaying the text.
//	fontSize (float64): How large to make the text. Note: It scales with the baseWidthVoxelResolution.
//
// Returns:
//
//	([]types.Triangle, error): A slice of triangles representing text.
func renderText(text string, justification string, leftOffsetPercent float64, fontSize float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Create a rendering context for the face of the skyline
	faceWidthRes := baseWidthVoxelResolution
	faceHeightRes := int(float64(faceWidthRes) * baseHeight / baseWidth)

	// Create image representing the skyline face
	dc := gg.NewContext(faceWidthRes, faceHeightRes)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)

	// Load font into context
	fontPath, cleanup, err := writeTempFont(PrimaryFont)
	if err != nil {
		// Try fallback font
		fontPath, cleanup, err = writeTempFont(FallbackFont)
		if err != nil {
			return nil, errors.New(errors.IOError, "failed to load any fonts", err)
		}
	}
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, errors.New(errors.IOError, "failed to load font", err)
	}

	// Draw text on image at desired location
	var triangles []types.Triangle

	// Convert justification to a number
	var justificationPercent float64
	switch justification {
	case "center":
		justificationPercent = 0.5
	case "right":
		justificationPercent = 1.0
	default:
		justificationPercent = 0.0
	}

	dc.DrawStringAnchored(
		text,
		float64(faceWidthRes)*leftOffsetPercent, // Offset from right
		float64(faceHeightRes)*0.5,              // Offset from top
		justificationPercent,                    // Justification (0.0=left, 0.5=center, 1.0=right)
		0.5,                                     // Vertically aligned
	)

	// Convert context image pixels into voxels
	for x := 0; x < faceWidthRes; x++ {
		for y := 0; y < faceHeightRes; y++ {
			if isPixelActive(dc, x, y) {
				voxel, err := createVoxelOnFace(
					float64(x),
					float64(y),
					voxelDepth,
					baseWidth,
					baseHeight,
				)
				if err != nil {
					return nil, errors.New(errors.STLError, "failed to create cube", err)
				}

				triangles = append(triangles, voxel...)
			}
		}
	}

	defer cleanup()

	return triangles, nil
}

// createVoxelOnFace creates a voxel on the face of a skyline by generating a cube at the specified coordinates.
// The function takes in the x, y coordinates and height.
// It returns a slice of types.Triangle representing the cube and an error if the cube creation fails.
//
// Parameters:
//
//	x (float64): The x-coordinate on the skyline face (left to right).
//	y (float64): The y-coordinate on the skyline face (top to bottom).
//	height (float64): Distance coming out of the face.
//
// Returns:
//
//	([]types.Triangle, error): A slice of triangles representing the cube and an error if any.
func createVoxelOnFace(x float64, y float64, height float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Mapping resolution
	xResolution := float64(baseWidthVoxelResolution)
	yResolution := xResolution * baseHeight / baseWidth

	// Pixel size
	voxelSize := 1.0

	// Scale coordinate to face resolution
	x = (x / xResolution) * baseWidth
	y = (y / yResolution) * baseHeight
	voxelSizeX := (voxelSize / xResolution) * baseWidth
	voxelSizeY := (voxelSize / yResolution) * baseHeight

	cube, err := CreateCube(
		// Location (from top left corner of skyline face)
		x,             // x - Left to right
		-height,       // y - Negative comes out of face. Positive goes into face.
		-voxelSizeY-y, // z - Bottom to top

		// Size
		voxelSizeX, // x length - left to right from specified point
		height,     // thickness - distance coming out of face
		voxelSizeY, // y length - bottom to top from specified point
	)

	return cube, err
}

// GenerateImageGeometry creates 3D geometry from the embedded logo image.
func GenerateImageGeometry(baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Get temporary image file
	imgPath, cleanup, err := getEmbeddedImage()
	if err != nil {
		return nil, err
	}

	defer cleanup()

	return renderImage(
		imgPath,
		logoScale,
		voxelDepth,
		logoLeftOffset,
		logoTopOffset,
		baseWidth,
		baseHeight,
	)
}

// renderImage generates 3D geometry for the given image configuration.
func renderImage(filePath string, scale float64, height float64, leftOffsetPercent float64, topOffsetPercent float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {

	// Get voxel resolution of base face
	faceWidthRes := baseWidthVoxelResolution
	faceHeightRes := int(float64(faceWidthRes) * baseHeight / baseWidth)

	// Load image from file
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(errors.IOError, "failed to open image", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			closeErr := errors.New(errors.IOError, "failed to close reader", err)
			// Log the error or handle it appropriately
			fmt.Println(closeErr)
		}
	}()
	img, err := png.Decode(reader)
	if err != nil {
		return nil, errors.New(errors.IOError, "failed to decode PNG", err)
	}

	// Get image size
	bounds := img.Bounds()
	logoWidth := bounds.Max.X
	logoHeight := bounds.Max.Y

	// Transfer image pixels onto face of skyline as voxels
	var triangles []types.Triangle
	for x := 0; x < logoWidth; x++ {
		for y := logoHeight - 1; y >= 0; y-- {
			// Get pixel color and alpha
			r, _, _, a := img.At(x, y).RGBA()

			// If pixel is active (white) and not fully transparent, create a voxel
			if a > 32768 && r > 32768 {

				voxel, err := createVoxelOnFace(
					(leftOffsetPercent*float64(faceWidthRes))+float64(x)*scale,
					(topOffsetPercent*float64(faceHeightRes))+float64(y)*scale,
					height,
					baseWidth,
					baseHeight,
				)

				if err != nil {
					return nil, errors.New(errors.STLError, "failed to create cube", err)
				}

				triangles = append(triangles, voxel...)
			}
		}
	}

	return triangles, nil
}

// isPixelActive checks if a pixel is active (white) in the given context.
func isPixelActive(dc *gg.Context, x, y int) bool {
	r, _, _, _ := dc.Image().At(x, y).RGBA()
	return r > 32768
}

// 윗면(Top Face)에 텍스트를 양각으로 생성하는 함수
func renderTextOnTop(text string, justification string, leftOffsetPercent, topOffsetPercent, fontSize, baseWidth, baseDepth, baseHeight float64) ([]types.Triangle, error) {
	faceWidthRes := baseWidthVoxelResolution
	faceDepthRes := int(float64(faceWidthRes) * baseDepth / baseWidth)

	dc := gg.NewContext(faceWidthRes, faceDepthRes)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)

	fontPath, cleanup, err := writeTempFont(PrimaryFont)
	if err != nil {
		fontPath, cleanup, err = writeTempFont(FallbackFont)
		if err != nil {
			return nil, errors.New(errors.IOError, "failed to load any fonts", err)
		}
	}
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, errors.New(errors.IOError, "failed to load font", err)
	}

	var justificationPercent float64
	switch justification {
	case "center":
		justificationPercent = 0.5
	case "right":
		justificationPercent = 1.0
	default:
		justificationPercent = 0.0
	}
	dc.DrawStringAnchored(
		text,
		float64(faceWidthRes)*leftOffsetPercent,
		float64(faceDepthRes)*topOffsetPercent,
		justificationPercent,
		0.5, // 수직 중앙 정렬
	)

	var triangles []types.Triangle
	for x := 0; x < faceWidthRes; x++ {
		for y := 0; y < faceDepthRes; y++ {
			if isPixelActive(dc, x, y) {
				voxel, err := createVoxelOnTop(
					float64(x),
					float64(y),
					voxelDepth, // 윗면에서 위로 양각
					baseWidth,
					baseDepth,
					baseHeight,
				)
				if err != nil {
					return nil, errors.New(errors.STLError, "failed to create cube", err)
				}
				triangles = append(triangles, voxel...)
			}
		}
	}
	defer cleanup()
	return triangles, nil
}

// 윗면에 voxel을 생성하는 함수
func createVoxelOnTop(x float64, y float64, height float64, baseWidth float64, baseDepth float64, baseHeight float64) ([]types.Triangle, error) {
	xResolution := float64(baseWidthVoxelResolution)
	yResolution := xResolution * baseDepth / baseWidth

	voxelSize := 1.0

	x = (x / xResolution) * baseWidth
	// y축 반전
	y = baseDepth - (y/yResolution)*baseDepth
	voxelSizeX := (voxelSize / xResolution) * baseWidth
	voxelSizeY := (voxelSize / yResolution) * baseDepth

	cube, err := CreateCube(
		x,          // x - Left to right
		y,          // y - Front to back (반전 적용)
		0,          // z - 박스 윗면(z=0)에 딱 붙게
		voxelSizeX, // x 길이
		voxelSizeY, // y 길이
		height,     // z 길이(양각 두께)
	)
	return cube, err
}

// 임의의 경로에서 이미지를 relief로 생성하는 함수
func GenerateImageGeometryWithPath(imgPath string, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// 좌측에 배치: scale, offset은 기존과 동일하게 사용
	return renderImage(
		imgPath,
		logoScale,
		voxelDepth,
		logoLeftOffset,
		logoTopOffset,
		baseWidth,
		baseHeight,
	)
}

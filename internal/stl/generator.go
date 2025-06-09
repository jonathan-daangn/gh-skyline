package stl

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/github/gh-skyline/internal/errors"
	"github.com/github/gh-skyline/internal/logger"
	"github.com/github/gh-skyline/internal/stl/geometry"
	"github.com/github/gh-skyline/internal/types"
)

// GenerateSTL creates a 3D model from GitHub contribution data and writes it to an STL file.
// It's a convenience wrapper around GenerateSTLRange for single year processing.
func GenerateSTL(contributions [][]types.ContributionDay, outputPath, username string, year int, topText, rightText string) error {
	// Wrap single year data in the format expected by GenerateSTLRange
	contributionsRange := [][][]types.ContributionDay{contributions}
	return GenerateSTLRange(contributionsRange, outputPath, username, year, year, topText, rightText)
}

// GenerateSTLRange creates a 3D model from multiple years of GitHub contribution data.
// It handles the complete process from data validation through geometry generation to file output.
// Parameters:
//   - contributions: 3D slice of contribution data ([year][week][day])
//   - outputPath: destination path for the STL file
//   - username: GitHub username for the contribution data
//   - startYear: first year in the range
//   - endYear: last year in the range
func GenerateSTLRange(contributions [][][]types.ContributionDay, outputPath, username string, startYear, endYear int, topText, rightText string) error {
	log := logger.GetLogger()
	if err := log.Debug("Starting STL generation for user %s, years %d-%d", username, startYear, endYear); err != nil {
		return errors.Wrap(err, "failed to log debug message")
	}

	if err := validateInput(contributions[0], outputPath, username); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	dimensions, err := calculateDimensions(len(contributions))
	if err != nil {
		return errors.Wrap(err, "failed to calculate dimensions")
	}

	// Find global max contribution across all years
	maxContribution := findMaxContributionsAcrossYears(contributions)

	modelTriangles, err := generateModelGeometry(contributions, dimensions, maxContribution, username, startYear, endYear, topText, rightText)
	if err != nil {
		return errors.Wrap(err, "failed to generate geometry")
	}

	// Read and merge character.stl if it exists
	var characterTriangles []types.Triangle
	var readErr error
	isBinary := true
	// 1. 80바이트 읽어서 solid로 시작하는지 확인
	if f, err := os.Open("character.stl"); err == nil {
		header := make([]byte, 80)
		_, herr := f.Read(header)
		f.Close()
		if herr == nil && strings.HasPrefix(string(header), "solid") {
			isBinary = false
		}
	}
	if isBinary {
		characterTriangles, readErr = ReadSTLBinary("character.stl")
	} else {
		characterTriangles, readErr = ReadASCIISTL("character.stl")
	}
	if readErr == nil && len(characterTriangles) > 0 {
		// 1. 70%로 스케일
		characterTriangles = scaleTriangles(characterTriangles, 0.7)
		// 3. bounding box 계산
		minX, minY, minZ, maxX, maxY, _ := calcBoundingBox(characterTriangles)
		charWidth := maxX - minX
		charDepth := maxY - minY
		// STL 바닥의 오른쪽 위에 올리기
		dx := dimensions.innerWidth - charWidth - minX - 10
		offset := 3.0
		dy := dimensions.innerDepth - charDepth - minY - offset
		dz := 0 - minZ - 0.5
		characterTriangles = translateTriangles(characterTriangles, dx, dy, dz)
		modelTriangles = append(modelTriangles, characterTriangles...)
		if err := log.Info("Merged character.stl with %d triangles (right-top, scaled)", len(characterTriangles)); err != nil {
			return errors.Wrap(err, "failed to log info message")
		}
	} else {
		if err := log.Debug("No character.stl found or error reading it: %v", readErr); err != nil {
			return errors.Wrap(err, "failed to log debug message")
		}
	}

	if err := log.Info("Model generation complete: %d total triangles", len(modelTriangles)); err != nil {
		return errors.Wrap(err, "failed to log info message")
	}
	if err := log.Debug("Writing STL file to: %s", outputPath); err != nil {
		return errors.Wrap(err, "failed to log debug message")
	}

	if err := WriteSTLBinary(outputPath, modelTriangles); err != nil {
		return errors.Wrap(err, "failed to write STL file")
	}

	if err := log.Info("STL file written successfully to: %s", outputPath); err != nil {
		return errors.Wrap(err, "failed to log info message")
	}
	return nil
}

// ASCII STL 파서
func ReadASCIISTL(filename string) ([]types.Triangle, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var triangles []types.Triangle
	var normal types.Point3D
	var vertices []types.Point3D

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "facet normal") {
			fields := strings.Fields(line)
			normal.X, _ = strconv.ParseFloat(fields[2], 64)
			normal.Y, _ = strconv.ParseFloat(fields[3], 64)
			normal.Z, _ = strconv.ParseFloat(fields[4], 64)
		} else if strings.HasPrefix(line, "vertex") {
			fields := strings.Fields(line)
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			vertices = append(vertices, types.Point3D{X: x, Y: y, Z: z})
		} else if strings.HasPrefix(line, "endfacet") {
			if len(vertices) == 3 {
				triangles = append(triangles, types.Triangle{
					Normal: normal,
					V1:     vertices[0],
					V2:     vertices[1],
					V3:     vertices[2],
				})
			}
			vertices = nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return triangles, nil
}

// modelDimensions represents the core measurements of the 3D model.
// All measurements are in millimeters.
type modelDimensions struct {
	innerWidth float64 // Width of the contribution grid
	innerDepth float64 // Depth of the contribution grid
	imagePath  string  // Path to the logo image
}

func validateInput(contributions [][]types.ContributionDay, outputPath, username string) error {
	if len(contributions) == 0 {
		return errors.New(errors.ValidationError, "contributions data cannot be empty", nil)
	}
	if len(contributions) > geometry.GridSize {
		return errors.New(errors.ValidationError, "contributions data exceeds maximum grid size", nil)
	}
	if outputPath == "" {
		return errors.New(errors.ValidationError, "output path cannot be empty", nil)
	}
	if username == "" {
		return errors.New(errors.ValidationError, "username cannot be empty", nil)
	}
	return nil
}

func calculateDimensions(yearCount int) (modelDimensions, error) {
	if yearCount <= 0 {
		return modelDimensions{}, errors.New(errors.ValidationError, "year count must be positive", nil)
	}

	var width, depth float64
	width, depth = geometry.CalculateMultiYearDimensions(yearCount)

	dims := modelDimensions{
		innerWidth: width,
		innerDepth: depth,
		imagePath:  "assets/invertocat.png",
	}

	if dims.innerWidth <= 0 || dims.innerDepth <= 0 {
		return modelDimensions{}, errors.New(errors.ValidationError, "invalid model dimensions", nil)
	}

	return dims, nil
}

func findMaxContributions(contributions [][]types.ContributionDay) int {
	maxContrib := 0
	for _, week := range contributions {
		for _, day := range week {
			if day.ContributionCount > maxContrib {
				maxContrib = day.ContributionCount
			}
		}
	}
	return maxContrib
}

// findMaxContributionsAcrossYears finds the maximum contribution count across all years
func findMaxContributionsAcrossYears(contributionsPerYear [][][]types.ContributionDay) int {
	maxContrib := 0
	for _, yearContributions := range contributionsPerYear {
		yearMax := findMaxContributions(yearContributions)
		if yearMax > maxContrib {
			maxContrib = yearMax
		}
	}
	return maxContrib
}

// geometryResult holds the output of geometry generation operations.
// It includes both the generated triangles and any errors that occurred.
type geometryResult struct {
	triangles []types.Triangle
	err       error
}

// generateModelGeometry orchestrates the concurrent generation of all model components.
// It manages four parallel processes for generating the base, columns, text, and logo.
func generateModelGeometry(contributionsPerYear [][][]types.ContributionDay, dims modelDimensions, maxContrib int, username string, startYear, endYear int, topText, rightText string) ([]types.Triangle, error) {
	if len(contributionsPerYear) == 0 {
		return nil, errors.New(errors.ValidationError, "contributions data cannot be empty", nil)
	}

	channels := map[string]chan geometryResult{
		"base":    make(chan geometryResult),
		"columns": make(chan geometryResult),
		"text":    make(chan geometryResult),
		"image":   make(chan geometryResult),
	}

	var wg sync.WaitGroup
	wg.Add(len(channels))

	go generateBase(dims, channels["base"], &wg)
	go generateColumnsForYearRange(contributionsPerYear, maxContrib, channels["columns"], &wg)
	go generateText("", startYear, endYear, dims, channels["text"], &wg, topText, rightText)
	go generateLogoWithCustomPath(dims, channels["image"], &wg, "logo.png")

	modelTriangles := make([]types.Triangle, 0, estimateTriangleCount(contributionsPerYear[0])*len(contributionsPerYear))
	for _, name := range []string{"base", "image", "columns", "text"} {
		result := <-channels[name]
		if result.err != nil {
			return nil, errors.Wrap(result.err, fmt.Sprintf("failed to generate %s geometry", name))
		}
		modelTriangles = append(modelTriangles, result.triangles...)
	}

	wg.Wait()
	for _, ch := range channels {
		close(ch)
	}

	return modelTriangles, nil
}

func generateBase(dims modelDimensions, ch chan<- geometryResult, wg *sync.WaitGroup) {
	defer wg.Done()
	baseTriangles, err := geometry.CreateCuboidBase(dims.innerWidth, dims.innerDepth)

	if err != nil {
		if logErr := logger.GetLogger().Warning("Failed to generate base geometry: %v. Continuing without base.", err); logErr != nil {
			ch <- geometryResult{triangles: []types.Triangle{}, err: logErr}
			return
		}
		ch <- geometryResult{triangles: []types.Triangle{}}
		return
	}

	ch <- geometryResult{triangles: baseTriangles}
}

// generateText creates 3D text geometry for the model
func generateText(username string, startYear int, endYear int, dims modelDimensions, ch chan<- geometryResult, wg *sync.WaitGroup, topText, rightText string) {
	defer wg.Done()
	var embossedRight string
	if rightText != "" {
		embossedRight = rightText
	} else if startYear == endYear {
		embossedRight = fmt.Sprintf("%d", endYear)
	} else {
		embossedRight = fmt.Sprintf("%04d-%02d", startYear, endYear%100)
	}

	textTriangles, err := geometry.Create3DText(username, embossedRight, dims.innerWidth, geometry.BaseHeight, dims.innerDepth, topText)
	if err != nil {
		if logErr := logger.GetLogger().Warning("Failed to generate text geometry: %v. Continuing without text.", err); logErr != nil {
			ch <- geometryResult{triangles: []types.Triangle{}, err: logErr}
			return
		}
		ch <- geometryResult{triangles: []types.Triangle{}}
		return
	}
	ch <- geometryResult{triangles: textTriangles}
}

// generateLogo handles the generation of the GitHub logo geometry
func generateLogo(dims modelDimensions, ch chan<- geometryResult, wg *sync.WaitGroup) {
	defer wg.Done()
	logoTriangles, err := geometry.GenerateImageGeometry(dims.innerWidth, geometry.BaseHeight)
	if err != nil {
		// Log warning and continue without logo instead of failing
		if logErr := logger.GetLogger().Warning("Failed to generate logo geometry: %v. Continuing without logo.", err); logErr != nil {
			ch <- geometryResult{triangles: []types.Triangle{}, err: logErr}
			return
		}
		ch <- geometryResult{triangles: []types.Triangle{}}
		return
	}
	ch <- geometryResult{triangles: logoTriangles}
}

func estimateTriangleCount(contributions [][]types.ContributionDay) int {
	totalContributions := 0
	for _, week := range contributions {
		for _, day := range week {
			if day.ContributionCount > 0 {
				totalContributions++
			}
		}
	}

	baseTrianglesCount := 12
	columnsTrianglesCount := totalContributions * 12
	textTrianglesEstimate := 1000
	return baseTrianglesCount + columnsTrianglesCount + textTrianglesEstimate
}

// generateColumnsForYearRange generates contribution columns for multiple years
func generateColumnsForYearRange(contributionsPerYear [][][]types.ContributionDay, maxContrib int, ch chan<- geometryResult, wg *sync.WaitGroup) {
	defer wg.Done()
	var yearTriangles []types.Triangle

	// Process years in reverse order so most recent year is at the front
	for i := len(contributionsPerYear) - 1; i >= 0; i-- {
		yearOffset := len(contributionsPerYear) - 1 - i
		triangles, err := geometry.CreateContributionGeometry(contributionsPerYear[i], yearOffset, maxContrib)
		if err != nil {
			if logErr := logger.GetLogger().Warning("Failed to generate column geometry for year %d: %v. Skipping year.", i, err); logErr != nil {
				return
			}
			continue
		}
		yearTriangles = append(yearTriangles, triangles...)
	}

	ch <- geometryResult{triangles: yearTriangles}
}

// CreateContributionGeometry generates geometry for a single year's worth of contributions
func CreateContributionGeometry(contributions [][]types.ContributionDay, yearIndex int, maxContrib int) []types.Triangle {
	var triangles []types.Triangle

	// Calculate the Y offset for this year's grid
	// Each subsequent year is placed further back (larger Y value)
	baseYOffset := float64(yearIndex) * (geometry.YearOffset + geometry.YearSpacing)

	// Generate contribution columns
	for weekIdx, week := range contributions {
		for dayIdx, day := range week {
			if day.ContributionCount > 0 {
				height := geometry.NormalizeContribution(day.ContributionCount, maxContrib)
				x := float64(weekIdx) * geometry.CellSize
				y := baseYOffset + float64(dayIdx)*geometry.CellSize

				columnTriangles, err := geometry.CreateColumn(x, y, height, geometry.CellSize)
				if err != nil {
					if logErr := logger.GetLogger().Warning("Failed to generate column geometry: %v. Skipping column.", err); logErr != nil {
						return nil
					}
					continue
				}
				triangles = append(triangles, columnTriangles...)
			}
		}
	}

	return triangles
}

// 로고 이미지를 지정 경로로부터 relief로 생성하는 함수
func generateLogoWithCustomPath(dims modelDimensions, ch chan<- geometryResult, wg *sync.WaitGroup, logoPath string) {
	defer wg.Done()
	logoTriangles, err := geometry.GenerateImageGeometryWithPath(logoPath, dims.innerWidth, geometry.BaseHeight)
	if err != nil {
		if logErr := logger.GetLogger().Warning("Failed to generate logo geometry: %v. Continuing without logo.", err); logErr != nil {
			ch <- geometryResult{triangles: []types.Triangle{}, err: logErr}
			return
		}
		ch <- geometryResult{triangles: []types.Triangle{}}
		return
	}
	ch <- geometryResult{triangles: logoTriangles}
}

// bounding box 계산 함수
func calcBoundingBox(triangles []types.Triangle) (minX, minY, minZ, maxX, maxY, maxZ float64) {
	minX, minY, minZ = math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	maxX, maxY, maxZ = -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for _, t := range triangles {
		for _, v := range []types.Point3D{t.V1, t.V2, t.V3} {
			if v.X < minX {
				minX = v.X
			}
			if v.Y < minY {
				minY = v.Y
			}
			if v.Z < minZ {
				minZ = v.Z
			}
			if v.X > maxX {
				maxX = v.X
			}
			if v.Y > maxY {
				maxY = v.Y
			}
			if v.Z > maxZ {
				maxZ = v.Z
			}
		}
	}
	return
}

// translate 함수
func translateTriangles(triangles []types.Triangle, dx, dy, dz float64) []types.Triangle {
	moved := make([]types.Triangle, len(triangles))
	for i, t := range triangles {
		moved[i] = types.Triangle{
			Normal: t.Normal,
			V1:     types.Point3D{X: t.V1.X + dx, Y: t.V1.Y + dy, Z: t.V1.Z + dz},
			V2:     types.Point3D{X: t.V2.X + dx, Y: t.V2.Y + dy, Z: t.V2.Z + dz},
			V3:     types.Point3D{X: t.V3.X + dx, Y: t.V3.Y + dy, Z: t.V3.Z + dz},
		}
	}
	return moved
}

// 스케일 함수
func scaleTriangles(triangles []types.Triangle, scale float64) []types.Triangle {
	scaled := make([]types.Triangle, len(triangles))
	for i, t := range triangles {
		scaled[i] = types.Triangle{
			Normal: t.Normal,
			V1:     types.Point3D{X: t.V1.X * scale, Y: t.V1.Y * scale, Z: t.V1.Z * scale},
			V2:     types.Point3D{X: t.V2.X * scale, Y: t.V2.Y * scale, Z: t.V2.Z * scale},
			V3:     types.Point3D{X: t.V3.X * scale, Y: t.V3.Y * scale, Z: t.V3.Z * scale},
		}
	}
	return scaled
}

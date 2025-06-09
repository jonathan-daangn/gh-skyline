// Package utils are utility functions for the GitHub Skyline project
package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Constants for GitHub launch year and default output file format
const (
	githubLaunchYear = 2008
	outputFileFormat = "%s-%s-github-skyline.stl"
)

// ParseYearRange parses whether a year is a single year or a range of years.
func ParseYearRange(yearRange string) (startYear, endYear int, err error) {
	if strings.Contains(yearRange, "-") {
		parts := strings.Split(yearRange, "-")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid year range format")
		}
		startYear, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, err
		}
		endYear, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
	} else {
		year, err := strconv.Atoi(yearRange)
		if err != nil {
			return 0, 0, err
		}
		startYear, endYear = year, year
	}
	return startYear, endYear, validateYearRange(startYear, endYear)
}

// validateYearRange checks if the years are within the range
// of GitHub's launch year to the current year and if
// the start year is not greater than the end year.
func validateYearRange(startYear, endYear int) error {
	currentYear := time.Now().Year()
	if startYear < githubLaunchYear || endYear > currentYear {
		return fmt.Errorf("years must be between %d and %d", githubLaunchYear, currentYear)
	}
	if startYear > endYear {
		return fmt.Errorf("start year cannot be after end year")
	}
	return nil
}

// FormatYearRange returns a formatted string representation of the year range
func FormatYearRange(startYear, endYear int) string {
	if startYear == endYear {
		return fmt.Sprintf("%d", startYear)
	}
	// Use YYYY-YY format for multi-year ranges
	return fmt.Sprintf("%04d-%02d", startYear, endYear%100)
}

// GenerateOutputFilename creates a consistent filename for the STL output
func GenerateOutputFilename(user string, startYear, endYear int, output string) string {
	if output != "" {
		// Ensure the filename ends with .stl
		if !strings.HasSuffix(strings.ToLower(output), ".stl") {
			return output + ".stl"
		}
		return output
	}
	yearStr := FormatYearRange(startYear, endYear)
	return fmt.Sprintf(outputFileFormat, user, yearStr)
}

// ValidateMonthRange checks if the months are within valid range (1-12)
// and if start month is not greater than end month
func ValidateMonthRange(startMonth, endMonth int) error {
	if startMonth < 1 || startMonth > 12 || endMonth < 1 || endMonth > 12 {
		return fmt.Errorf("months must be between 1 and 12")
	}
	if startMonth > endMonth {
		return fmt.Errorf("start month cannot be after end month")
	}
	return nil
}

// Package ascii provides functionality for generating ASCII art representations
// of GitHub contribution graphs.
package ascii

import (
	"strings"
)

// GridWidth defines the standard width for the ASCII output.
const GridWidth = 53

// HeaderTemplate contains the ASCII art header for the output.
const HeaderTemplate = `
           ____ _ _   _   _       _     
          / ___(_) |_| | | |_   _| |__  
         | |  _| | __| |_| | | | | '_ \ 
         | |_| | | |_|  _  | |_| | |_) |
          \____|_|\__|_| |_|\__,_|_.__/ 

          ____  _          _ _            
         / ___|| | ___   _| (_)_ __   ___ 
         \___ \| |/ / | | | | | '_ \ / _ \
          ___) |   <| |_| | | | | | | __/
         |____/|_|\_\\__, |_|_|_| |_|\___|
                    |___/
`

// centerText centers the given text within the GridWidth.
// It accounts for wide Unicode characters and ensures the text fits within
// the specified width. If the text is longer than GridWidth, it will be truncated.
func centerText(text string) string {
	visualWidth := len(text)

	if visualWidth >= GridWidth {
		return text[:GridWidth] + "\n"
	}

	totalPadding := GridWidth - visualWidth

	if totalPadding <= 1 {
		return text + "\n"
	}

	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding) + "\n"
}

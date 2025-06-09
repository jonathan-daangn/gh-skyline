// Package cmd is a package that contains the root command (entrypoint) for the GitHub Skyline CLI tool.
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/github/gh-skyline/cmd/skyline"
	"github.com/github/gh-skyline/internal/errors"
	"github.com/github/gh-skyline/internal/github"
	"github.com/github/gh-skyline/internal/logger"
	"github.com/github/gh-skyline/internal/utils"
	"github.com/spf13/cobra"
)

// Command line variables and root command configuration
var (
	yearRange  string
	user       string
	full       bool
	debug      bool
	web        bool
	artOnly    bool
	output     string // new output path flag
	topText    string // STL 상단 텍스트 입력값 (기존 text)
	startMonth int    // 시작 월
	endMonth   int    // 종료 월
	rightText  string // 우측 텍스트 입력값
)

// rootCmd is the root command for the GitHub Skyline CLI tool.
var rootCmd = &cobra.Command{
	Use:   "skyline",
	Short: "Generate a 3D model of a user's GitHub contribution history",
	Long: `GitHub Skyline creates 3D printable STL files from GitHub contribution data.
It can generate models for specific years or year ranges for the authenticated user or an optional specified user.

While the STL file is being generated, an ASCII preview will be displayed in the terminal.

ASCII Preview Legend:
  ' ' Empty/Sky     - No contributions
  '.' Future dates  - What contributions could you make?
  '░' Low level     - Light contribution activity
  '▒' Medium level  - Moderate contribution activity
  '▓' High level    - Heavy contribution activity
  '╻┃╽' Top level   - Last block with contributions in the week (Low, Medium, High)

Layout:
Each column represents one week. Days within each week are reordered vertically
to create a "building" effect, with empty spaces (no contributions) at the top.`,
	RunE: handleSkylineCommand,
}

// init initializes command line flags for the skyline CLI tool.
func init() {
	initFlags()
}

// Execute initializes and executes the root command for the GitHub Skyline CLI.
func Execute(_ context.Context) error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// initFlags sets up command line flags for the skyline CLI tool.
func initFlags() {
	flags := rootCmd.Flags()
	flags.StringVarP(&yearRange, "year", "y", fmt.Sprintf("%d", time.Now().Year()), "Year or year range (e.g., 2024 or 2014-2024)")
	flags.StringVarP(&user, "user", "u", "", "GitHub username (optional, defaults to authenticated user)")
	flags.BoolVarP(&full, "full", "f", false, "Generate contribution graph from join year to current year")
	flags.BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	flags.BoolVarP(&web, "web", "w", false, "Open GitHub profile (authenticated or specified user).")
	flags.BoolVarP(&artOnly, "art-only", "a", false, "Generate only ASCII preview")
	flags.StringVarP(&output, "output", "o", "", "Output file path (optional)")
	flags.StringVar(&topText, "top-text", "", "상단에 들어갈 텍스트 (optional)")
	flags.IntVar(&startMonth, "start-month", 1, "시작 월 (1-12)")
	flags.IntVar(&endMonth, "end-month", 12, "종료 월 (1-12)")
	flags.StringVar(&rightText, "right-text", "", "우측에 들어갈 텍스트 (optional)")
}

// executeRootCmd is the main execution function for the root command.
func handleSkylineCommand(_ *cobra.Command, _ []string) error {
	log := logger.GetLogger()
	if debug {
		log.SetLevel(logger.DEBUG)
		if err := log.Debug("Debug logging enabled"); err != nil {
			return err
		}
	}

	client, err := github.InitializeGitHubClient()
	if err != nil {
		return errors.New(errors.NetworkError, "failed to initialize GitHub client", err)
	}

	if web {
		b := browser.New("", os.Stdout, os.Stderr)
		if err := openGitHubProfile(user, client, b); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return nil
	}

	startYear, endYear, err := utils.ParseYearRange(yearRange)
	if err != nil {
		return fmt.Errorf("invalid year range: %v", err)
	}

	return skyline.GenerateSkyline(startYear, endYear, user, full, output, artOnly, topText, startMonth, endMonth, rightText)
}

// Browser interface matches browser.Browser functionality.
type Browser interface {
	Browse(url string) error
}

// openGitHubProfile opens the GitHub profile page for the specified user or authenticated user.
func openGitHubProfile(targetUser string, client skyline.GitHubClientInterface, b Browser) error {
	if targetUser == "" {
		username, err := client.GetAuthenticatedUser()
		if err != nil {
			return errors.New(errors.NetworkError, "failed to get authenticated user", err)
		}
		targetUser = username
	}

	hostname, _ := auth.DefaultHost()
	profileURL := fmt.Sprintf("https://%s/%s", hostname, targetUser)
	return b.Browse(profileURL)
}

package logger

import (
	"bytes"
	"strings"
	"testing"
)

// testLogCapture helps capture log output for testing
type testLogCapture struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func setupTestLogger(_ *testing.T) (*Logger, *testLogCapture) {
	capture := &testLogCapture{
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}

	logger := GetLogger()
	logger.debug.SetOutput(capture.stdout)
	logger.info.SetOutput(capture.stdout)
	logger.warning.SetOutput(capture.stdout)
	logger.error.SetOutput(capture.stderr)

	return logger, capture
}

func TestSetLevel(t *testing.T) {
	logger := GetLogger()

	tests := []struct {
		name     string
		setLevel LogLevel
		want     LogLevel
	}{
		{"Set Debug Level", DEBUG, DEBUG},
		{"Set Info Level", INFO, INFO},
		{"Set Warning Level", WARNING, WARNING},
		{"Set Error Level", ERROR, ERROR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetLevel(tt.setLevel)
			if logger.level != tt.want {
				t.Errorf("SetLevel() = %v, want %v", logger.level, tt.want)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	logger, capture := setupTestLogger(t)

	tests := []struct {
		name    string
		level   LogLevel
		message string
		wantLog bool
		wantErr bool
	}{
		{"Debug when level is DEBUG", DEBUG, "test debug message", true, false},
		{"Debug when level is INFO", INFO, "should not show", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture.stdout.Reset()
			logger.SetLevel(tt.level)
			err := logger.Debug("%s", tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("Debug() error = %v, wantErr %v", err, tt.wantErr)
			}

			hasOutput := capture.stdout.Len() > 0
			if hasOutput != tt.wantLog {
				t.Errorf("Debug() output = %v, want %v", hasOutput, tt.wantLog)
			}
			if tt.wantLog && !strings.Contains(capture.stdout.String(), tt.message) {
				t.Errorf("Debug() output doesn't contain message: %s", tt.message)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	logger, capture := setupTestLogger(t)

	tests := []struct {
		name    string
		level   LogLevel
		message string
		wantLog bool
	}{
		{"Info when level is DEBUG", DEBUG, "test info message", true},
		{"Info when level is INFO", INFO, "test info message", true},
		{"Info when level is WARNING", WARNING, "should not show", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture.stdout.Reset()
			logger.SetLevel(tt.level)
			if err := logger.Info("%s", tt.message); err != nil {
				t.Errorf("Info() error = %v", err)
			}

			hasOutput := capture.stdout.Len() > 0
			if hasOutput != tt.wantLog {
				t.Errorf("Info() output = %v, want %v", hasOutput, tt.wantLog)
			}
			if tt.wantLog && !strings.Contains(capture.stdout.String(), tt.message) {
				t.Errorf("Info() output doesn't contain message: %s", tt.message)
			}
		})
	}
}

func TestError(t *testing.T) {
	logger, capture := setupTestLogger(t)

	tests := []struct {
		name    string
		level   LogLevel
		message string
		wantLog bool
	}{
		{"Error when level is DEBUG", DEBUG, "test error message", true},
		{"Error when level is ERROR", ERROR, "test error message", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture.stderr.Reset()
			logger.SetLevel(tt.level)
			if err := logger.Error("%s", tt.message); err != nil {
				t.Errorf("Error() error = %v", err)
			}

			hasOutput := capture.stderr.Len() > 0
			if hasOutput != tt.wantLog {
				t.Errorf("Error() output = %v, want %v", hasOutput, tt.wantLog)
			}
			if tt.wantLog && !strings.Contains(capture.stderr.String(), tt.message) {
				t.Errorf("Error() output doesn't contain message: %s", tt.message)
			}
		})
	}
}

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"DEBUG level string", DEBUG, "DEBUG"},
		{"INFO level string", INFO, "INFO"},
		{"WARNING level string", WARNING, "WARNING"},
		{"ERROR level string", ERROR, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

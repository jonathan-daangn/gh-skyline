#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Extension name
EXTENSION="gh-skyline"

# Function to print status
print_status() {
	echo -e "${BLUE}[$(date '+%H:%M:%S')] $1${NC}"
}

# Function to check command status
check_status() {
	if [ $? -eq 0 ]; then
		echo -e "${GREEN}✓ Success${NC}"
	else
		echo -e "${RED}✗ Failed${NC}"
		exit 1
	fi
}

# Start time
START_TIME=$(date +%s)

# Check if gh CLI is installed
if ! command -v gh &>/dev/null; then
	echo -e "${RED}Error: GitHub CLI is not installed${NC}"
	exit 1
fi

# Remove existing extension
print_status "Removing existing extension..."
gh extension remove $EXTENSION 2>/dev/null || true

# Build extension
print_status "Building extension..."
go build -o $EXTENSION
check_status

# Install extension
print_status "Installing extension..."
gh extension install .
check_status

# Run extension
print_status "Running skyline..."
gh skyline

# Calculate execution time
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
echo -e "\n${GREEN}Completed in ${DURATION} seconds${NC}"

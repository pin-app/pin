#!/bin/bash

# Test script for development data seeding
# This script uses the Makefile to set up a clean development environment

set -e

echo "Setting up development environment with seeding..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "Docker is not running. Please start Docker first."
    exit 1
fi

# Use the Makefile to start the development environment
echo "Starting development environment with fresh database..."
make dev-reset

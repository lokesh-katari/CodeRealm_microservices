#!/bin/sh

# Check if a code file is provided as an argument
if [ -z "$1" ]; then
    echo "Error: No code file provided."
    exit 1
fi

# Validate file existence and extension
if [ ! -f "$1" ]; then
    echo "Error: File '$1' not found."
    exit 1
fi

if [ "${1##*.}" != "rb" ]; then
    echo "Error: '$1' is not a valid Ruby source file."
    exit 1
fi

# Execute the Ruby code file and capture output
ruby "$1"

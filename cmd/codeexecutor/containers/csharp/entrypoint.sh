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

if [ "${1##*.}" != "cs" ]; then
    echo "Error: '$1' is not a valid C# source file."
    exit 1
fi

# Compile and execute the C# code file
mcs "$1" && mono "${1%.*}.exe"

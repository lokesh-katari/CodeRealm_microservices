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

if [ "${1##*.}" != "rs" ]; then
    echo "Error: '$1' is not a valid Rust source file."
    exit 1
fi

# Compile and execute the Rust code file
rustc "$1" && "${1%.*}"


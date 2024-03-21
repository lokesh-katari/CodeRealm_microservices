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

if [ "${1##*.}" != "java" ]; then
    echo "Error: '$1' is not a valid Java source file."
    exit 1
fi

# Copy the code file provided as an argument into the container
code_file="$1"

javac "$code_file"
if [ $? -ne 0 ]; then
    echo "Error: Compilation failed."
    exit 1
fi

# Execute the code file and capture output
main_class=$(basename "$code_file" .java)

# Execute the code file and capture output
output=$(java "$main_class")

# Output the result
echo "$output"
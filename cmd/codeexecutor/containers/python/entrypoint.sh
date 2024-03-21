#!/bin/sh

# Check if a code file is provided as an argument
if [ -z "$1" ]; then
    echo "Error: No code file provided."
    exit 1
fi

# Copy the code file provided as an argument into the container
code_file="$1"
cp "$code_file" /app/code_file

# Execute the code file and capture output
output=$(python /app/code_file)

# Output the result
echo "$output"

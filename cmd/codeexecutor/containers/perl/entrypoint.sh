#!/bin/sh

# Check if a Perl script file is provided as an argument
if [ -z "$1" ]; then
    echo "Error: No Perl script file provided."
    exit 1
fi

# Validate file existence and extension
if [ ! -f "$1" ]; then
    echo "Error: File '$1' not found."
    exit 1
fi

if [ "${1##*.}" != "pl" ]; then
    echo "Error: '$1' is not a valid Perl script file."
    exit 1
fi

# Execute the Perl script file
perl "$1"

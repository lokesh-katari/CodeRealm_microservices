#!/bin/bash

# Check if a C# file is provided as an argument
if [ -z "$1" ]; then
  echo "Error: No C# file provided."
  exit 1
fi

# Check if the provided file exists
if [ ! -f "$1" ]; then
  echo "Error: File '$1' not found."
  exit 1
fi

# Check if the file extension is .cs
if [ "${1##*.}" != "cs" ]; then
  echo "Error: '$1' is not a valid C# source file."
  exit 1
fi

mcs /app/*.cs 

mono "${1%.*}.exe"
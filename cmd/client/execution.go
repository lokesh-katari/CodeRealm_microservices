package main

import (
	"encoding/json"
	"fmt"
	"lokesh-katari/code-realm/cmd/client/models"
	"regexp"
	"strconv"
)

func GenerateCode(language string, userCode string, template models.Templates, reqType string) (string, error) {
	var finalCode string

	if reqType == "submit" {
		switch language {
		case "Python":
			finalCode = userCode + template.Python.HiddenTestCode
		case "JavaScript":
			finalCode = userCode + template.JavaScript.HiddenTestCode
		case "Golang":
			finalCode = userCode + template.Golang.HiddenTestCode
		case "Java":
			finalCode = userCode + template.Java.HiddenTestCode
		case "C":
			finalCode = userCode + template.C.HiddenTestCode
		case "Cpp":
			finalCode = userCode + template.Cpp.HiddenTestCode
		default:
			return "", fmt.Errorf("Invalid language: %s", language)
		}
	} else {
		switch language {
		case "Python":
			finalCode = userCode + template.Python.RunTestCode
		case "JavaScript":
			finalCode = userCode + template.JavaScript.RunTestCode
		case "Golang":
			finalCode = userCode + template.Golang.RunTestCode
		case "Java":
			finalCode = userCode + template.Java.RunTestCode
		case "C":
			finalCode = userCode + template.C.RunTestCode
		case "Cpp":
			finalCode = userCode + template.Cpp.RunTestCode
		default:
			return "", fmt.Errorf("Invalid language: %s", language)
		}
	}
	return finalCode, nil
}

func Seperateoutput(output string, language string) (string, error) {
	runtimeRegex := regexp.MustCompile(`Runtime: (\d+\.\d+) seconds`)
	statusRegex := regexp.MustCompile(`All Test Cases Passed: (True|False|true|false)`)

	var runtime float64
	var status bool
	var passedTestCasesInt []int

	if match := runtimeRegex.FindStringSubmatch(output); len(match) > 1 {
		runtime, _ = strconv.ParseFloat(match[1], 64)
	} else {
		runtime = 0.0
	}

	if match := statusRegex.FindStringSubmatch(output); len(match) > 1 {
		status = (match[1] == "True" || match[1] == "true")
	}

	// Extract the test cases that passed
	// testCaseRegex := regexp.MustCompile(`Test Case (\d+): .*Status: Accepted`)
	// matches := testCaseRegex.FindAllStringSubmatch(output, -1)
	// for _, match := range matches {
	// 	testCaseNumber, _ := strconv.Atoi(match[1])
	// 	passedTestCases = append(passedTestCases, testCaseNumber)
	// }
	pattern := `Passed Test Cases:\s*\[([^]]+)\]`

	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Find the matches in the output string
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		// Extract the matched string containing the passed test cases array
		passedTestCasesStr := matches[1]

		// Split the string into individual test case values
		passedTestCases := regexp.MustCompile(`\s*,\s*`).Split(passedTestCasesStr, -1)
		for _, testCase := range passedTestCases {
			tc, err := strconv.Atoi(testCase)
			if err != nil {
				fmt.Println("Error converting test case value to integer:", err)
				return "", err
			}
			passedTestCasesInt = append(passedTestCasesInt, tc)
		}

		fmt.Println("Passed Test Cases:", passedTestCasesInt)
	} else {
		fmt.Println("Passed Test Cases not found in the output string.")
	}

	// Create a map to hold the JSON response
	jsonResponse := map[string]interface{}{
		"output":          output,
		"runtime":         runtime,
		"language":        language,
		"status":          status,
		"passedTestCases": passedTestCasesInt,
	}

	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		panic(err)
	}
	return string(jsonResponseBytes), nil
}

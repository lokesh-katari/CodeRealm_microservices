package main

import (
	"fmt"
	"lokesh-katari/code-realm/cmd/client/models"
)

func GenerateCode(language string, userCode string, problem models.CodeQue) (string, error) {
	var finalCode string

	switch language {
	case "Python":
		finalCode = userCode + problem.Templates.Python.HiddenTestCode
	case "JavaScript":
		finalCode = userCode + problem.Templates.JavaScript.HiddenTestCode
	case "Golang":
		finalCode = userCode + problem.Templates.Golang.HiddenTestCode
	case "Java":
		finalCode = userCode + problem.Templates.Java.HiddenTestCode
	case "C":
		finalCode = userCode + problem.Templates.C.HiddenTestCode
	case "Cpp":
		finalCode = userCode + problem.Templates.Cpp.HiddenTestCode
	default:
		return "", fmt.Errorf("Invalid language: %s", language)
	}
	return finalCode, nil
}

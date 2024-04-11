package main

func GenerateCode(lang string, userCode string,hiddentestcasesSnippet string ){
	var finalCode string;
	switch lang {
	case "python":
		finalCode = userCode+hiddentestcasesSnippet

		
	
	case "java":
		finalCode= `import java.util.*; `+userCode+hiddentestcasesSnippet

	case "c":
		finalCode= `#include <stdio.h> `+userCode+hiddentestcasesSnippet

	case "cpp":
		finalCode= `#include <iostream> `+userCode+hiddentestcasesSnippet
		
	
}
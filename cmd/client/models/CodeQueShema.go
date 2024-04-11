package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// LanguageTemplate represents the user code and hidden test code for a specific language.
type LanguageTemplate struct {
	UserCode       string `json:"userCode"`
	HiddenTestCode string `json:"hiddenTestCode"`
}
type Languages struct {
	Python     LanguageTemplate `json:"python"`
	JavaScript LanguageTemplate `json:"javascript"`
	Golang     LanguageTemplate `json:"golang"`
	Java       LanguageTemplate `json:"java"`
	C          LanguageTemplate `json:"c"`
	Cpp        LanguageTemplate `json:"cpp"`
}

// Submission represents the number of correct and wrong submissions.
type Submission struct {
	Correct int `json:"correct"`
	Wrong   int `json:"wrong"`
}

// TestCase represents the input and output test cases.
type TestCase struct {
	Input  []interface{} `json:"input"`
	Output []interface{} `json:"output"`
}

// Problem represents a problem in the MongoDB document.
type Problem struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProblemID   int                `json:"problemId" bson:"problemId"`
	Difficulty  string             `json:"difficulty" bson:"difficulty"`
	Description string             `json:"description" bson:"description"`
	Submissions Submission         `json:"submissions" bson:"submissions"`
	Templates   Languages          `json:"templates" bson:"templates"`
	TestCases   []TestCase         `json:"testCases" bson:"testCases"`
}

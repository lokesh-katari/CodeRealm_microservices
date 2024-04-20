package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// LanguageTemplate represents the user code and hidden test code for a specific language.
type LanguageTemplate struct {
	UserCode       string `json:"userCode"`
	HiddenTestCode string `json:"hiddenTestCode"`
	RunTestCode    string `json:"runTestCode"`
}
type Templates struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Python     LanguageTemplate   `json:"python"`
	JavaScript LanguageTemplate   `json:"javascript"`
	Golang     LanguageTemplate   `json:"golang"`
	Java       LanguageTemplate   `json:"java"`
	C          LanguageTemplate   `json:"c"`
	Cpp        LanguageTemplate   `json:"cpp"`
}

// Submission represents the number of correct and wrong submissions.
type Submission struct {
	Correct int `json:"correct"`
	Wrong   int `json:"wrong"`
}
type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// TestCase represents the input and output test cases.
type CodeQue struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Difficulty  string             `json:"difficulty" bson:"difficulty"`
	Category    string             `json:"category" bson:"category"`
	Description string             `json:"description" bson:"description"`
	Submissions Submission         `json:"submissions" bson:"submissions"`
	TemplateID  primitive.ObjectID `json:"templateId" bson:"templateId"`
	TestCases   []TestCase         `json:"testCases" bson:"testCases"`
}

// model CodeQue {
// 	id              String            @id @default(auto()) @map("_id") @db.ObjectId
// 	title           String
// 	difficulty      String
// 	category        String
// 	description     String
// 	submissions     Submission
// 	templates       Languages         @relation(fields: [templateId], references: [id])
// 	templateId      String            @unique @db.ObjectId
// 	testCases       TestCase[]
// 	Codesubmissions Codesubmissions[]

// 	@@map("CodeQues")
//   }

package codeservice

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CodeService is the service that provides the code execution functionalityf

type CodeService interface {
	// codeRepo is the repository that provides the code execution functionality
	ExecuteCode(ctx context.Context, code string) (string, error)
}

type CodeSub struct {
	PID       string
	Email     string
	Code      string
	TimeStamp string
}

type CodeQuestion struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Template        []LanguageTemplate `bson:"template"`
	Category        string             `bson:"category"`
	Submissions     SubmissionStats    `bson:"submissions"`
	ProblemTitle    string             `bson:"problemTitle"`
	ProblemDesc     string             `bson:"problemDesc"`
	CreatedAt       time.Time          `bson:"createdAt"`
	DifficultyLevel string             `bson:"difficultyLevel"`
	Output          []string           `bson:"output"`
	HiddenOutputs   []string           `bson:"hiddenOutputs"`
	HiddenTestCases []HiddenTestCase   `bson:"hiddenTestCases"`
}

type LanguageTemplate struct {
	Language     string `bson:"language"`
	CodeTemplate string `bson:"codeTemplate"`
}

type SubmissionStats struct {
	Wrong   int `bson:"wrong"`
	Correct int `bson:"correct"`
}

type HiddenTestCase struct {
	Name   []string `bson:"name"`
	Values []string `bson:"values"`
}

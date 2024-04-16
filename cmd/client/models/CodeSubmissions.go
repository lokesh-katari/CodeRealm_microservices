package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CodeSubmission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PID         string             `json:"pid"`
	QueID       primitive.ObjectID `bson:"queId,omitempty"`
	Email       string             `json:"email"`
	Language    string             `json:"language"`
	Code        string             `json:"code"`
	Output      string             `json:"output"`
	SubmittedAT time.Time          `json:"submittedAt"`
	Runtime     string             `json:"runtime"`
	Memory      string             `json:"memory"`
}

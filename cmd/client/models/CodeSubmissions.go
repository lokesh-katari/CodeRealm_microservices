package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CodeSubmission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PID         string             `json:"pid"`
	QueID       primitive.ObjectID `bson:"queId,omitempty"`
	Title       string             `json:"title"`
	Accepted    bool               `json:"accepted"`
	Email       string             `json:"email"`
	Code        string             `json:"code"`
	Language    string             `json:"language"`
	Testcases   []int              `json:"testcases"`
	Runtime     string             `json:"runtime"`
	SubmittedAT time.Time          `json:"submittedAt"`
	Output      string             `json:"output"`
}

// model Codesubmissions {
// 	id          String   @id @default(auto()) @map("_id") @db.ObjectId
// 	pid         String   @unique
// 	queId       String   @db.ObjectId
// 	title       String
// 	Accepted    Boolean
// 	email       String?  @unique
// 	code        String
// 	language    String
// 	testcases   Int[]
// 	runtime     String
// 	memory      String
// 	submittedAt DateTime @default(now())
// 	output      String
// 	question    CodeQue? @relation(fields: [queId], references: [id])

// 	@@map("submissions")
//   }

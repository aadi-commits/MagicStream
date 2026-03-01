package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIReviewLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	MovieID string	`bson:"movie_id"`
	AdminReview string	`bson:"admin_review"`
	AIRawResponse string	`bson:"ai_raw_response"`
	ParsedRank int	`bson:"parsed_rank"`
	Status string	`bson:"status"`
	ErrorMessage string	`bson:"error_message,omitempty"`
	CreatedAt time.Time	`bson:"created_at"`
}
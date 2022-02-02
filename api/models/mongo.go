package models

import (
	"time"

	"github.com/s3f4/ginterview/api/library"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Record holds collection data
type Record struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Key        string             `json:"key" bson:"key"`
	Value      string             `json:"value" bson:"value"`
	Counts     []int              `json:"counts" bson:"counts"`
	TotalCount int                `json:"totalCount"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

// Requests

// MongoRequests holds the data coming from user
type MongoRequest struct {
	MinCount  int                 `json:"minCount" validate:"required"`
	MaxCount  int                 `json:"maxCount" validate:"required"`
	StartDate *library.CustomDate `json:"startDate" validate:"required,datetime"`
	EndDate   *library.CustomDate `json:"endDate" validate:"required,datetime"`
}

// Responses
// ResponseRecord is formatted record
type ResponseRecord struct {
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}

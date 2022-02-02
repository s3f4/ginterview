package repository

import (
	"context"
	"log"

	"github.com/s3f4/ginterview/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository ..
type MongoRepository interface {
	List(context.Context, *models.MongoRequest) ([]*models.Record, error)
}

var (
	// DB ...
	DB = "getir-case-study" //os.Getenv("MONGO_DATABASE")
	// recordsTable
	recordsTable = "records" //os.Getenv("MONGO_IMAGE_TABLE")
)

type mongoRepository struct {
	client *mongo.Client
}

// NewMongoRepository returns an MongoRepository object
func NewMongoRepository(client *mongo.Client) MongoRepository {
	return &mongoRepository{
		client,
	}
}

func (r *mongoRepository) getCollection() *mongo.Collection {
	return r.client.Database(DB).Collection(recordsTable)
}

// List gets values from mongodb and returns *models.Record array
func (r *mongoRepository) List(ctx context.Context, request *models.MongoRequest) ([]*models.Record, error) {
	var records []*models.Record
	collection := r.getCollection()
	filter := bson.A{
		// create a totalCount property on document
		bson.M{
			"$set": bson.M{
				"totalCount": bson.M{
					"$sum": "$counts",
				},
			},
		},
		// filter by request's startDate and endDate
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(request.StartDate.Time),
					"$lte": primitive.NewDateTimeFromTime(request.EndDate.Time),
				},
			},
		},
		// filter by request's minCount and maxCount
		bson.M{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gte": request.MinCount,
					"$lte": request.MaxCount,
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, filter)
	if err != nil {
		log.Print("mongoRepository.List ", err)
		return nil, err
	}

	// Fill recoreds
	for cursor.Next(ctx) {
		record := models.Record{}
		err := cursor.Decode(&record)
		// var x interface{}
		// fmt.Println(x)
		if err != nil {
			log.Print("mongoRepository.List cursor.Decode ", err)
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

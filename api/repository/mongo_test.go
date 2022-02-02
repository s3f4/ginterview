package repository

import (
	"context"
	"testing"
	"time"

	"github.com/s3f4/ginterview/api/library"
	"github.com/s3f4/ginterview/api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_MongoRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	createdAt1 := primitive.NewDateTimeFromTime(time.Now())
	createdAt2 := primitive.NewDateTimeFromTime(time.Now())

	mt.Run("success", func(mt *mtest.T) {
		mongoRepository := NewMongoRepository(mt.Client)
		first := mtest.CreateCursorResponse(1, "mongo.record", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"key", "dcJUSDLR"},
			{"value", "lVZeOUSDkQjx"},
			{"createdAt", createdAt1},
			{"counts", bson.A{50, 50, 50}},
		})
		second := mtest.CreateCursorResponse(1, "mongo.record", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"key", "NOdGNUDn"},
			{"value", "lVZeOUSDkQjx"},
			{"createdAt", createdAt2},
			{"counts", bson.A{50, 50, 50}},
		})
		killCursors := mtest.CreateCursorResponse(0, "mongo.record", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		records, err := mongoRepository.List(context.Background(), &models.MongoRequest{
			StartDate: &library.CustomDate{Time: createdAt1.Time().UTC()},
			EndDate:   &library.CustomDate{Time: createdAt2.Time().UTC()},
		})

		resultRecords := []*models.Record{
			{ID: id1, Key: "dcJUSDLR", Value: "lVZeOUSDkQjx", Counts: []int{50, 50, 50}, CreatedAt: createdAt1.Time().UTC()},
			{ID: id2, Key: "NOdGNUDn", Value: "lVZeOUSDkQjx", Counts: []int{50, 50, 50}, CreatedAt: createdAt2.Time().UTC()},
		}

		assert.Nil(t, err)
		assert.Equal(t, records, resultRecords)
	})

	mt.Run("command_fail", func(mt *mtest.T) {
		mongoRepository := NewMongoRepository(mt.Client)
		mt.AddMockResponses(bson.D{{"ok", 0}})
		_, err := mongoRepository.List(context.Background(), &models.MongoRequest{
			StartDate: &library.CustomDate{Time: createdAt1.Time().UTC()},
			EndDate:   &library.CustomDate{Time: createdAt2.Time().UTC()},
		})

		assert.NotNil(t, err)
		assert.IsType(t, err, mongo.CommandError{})
	})

	mt.Run("cursor_fail", func(mt *mtest.T) {
		mongoRepository := NewMongoRepository(mt.Client)

		first := mtest.CreateCursorResponse(0, "mongo.record", mtest.FirstBatch, bson.D{
			{"_id", 1},
		})
		mt.AddMockResponses(first)

		records, err := mongoRepository.List(context.Background(), &models.MongoRequest{
			StartDate: &library.CustomDate{Time: createdAt1.Time().UTC()},
			EndDate:   &library.CustomDate{Time: createdAt2.Time().UTC()},
		})

		assert.Nil(t, records)
		assert.NotNil(t, err)
	})
}

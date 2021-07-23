package result

import (
	"context"

	"github.com/kubeshop/kubetest/pkg/api/kubetest"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "results"

func NewMongoRespository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		Coll: db.Collection(CollectionName),
	}
}

type MongoRepository struct {
	Coll *mongo.Collection
}

func (r *MongoRepository) Get(ctx context.Context, id string) (result kubetest.ScriptExecution, err error) {
	err = r.Coll.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	return
}

func (r *MongoRepository) Insert(ctx context.Context, result kubetest.ScriptExecution) (err error) {
	_, err = r.Coll.InsertOne(ctx, result)
	return
}

func (r *MongoRepository) Update(ctx context.Context, result kubetest.ScriptExecution) (err error) {
	_, err = r.Coll.ReplaceOne(ctx, bson.M{"id": result.Id}, result)
	return
}

func (r *MongoRepository) QueuePull(ctx context.Context) (result kubetest.ScriptExecution, err error) {
	err = r.Coll.FindOneAndUpdate(ctx, bson.M{"status": kubetest.ExecutionStatusQueued}, bson.M{"$set": bson.M{"status": kubetest.ExecutionStatusPending}}).Decode(&result)
	return
}
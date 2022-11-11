package repository

import (
	"context"
	"errors"
	"os"

	pkg_mongodb "go-rengan/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-rengan/todo/models"
	"go-rengan/utils"
)

// MongoTodoRepository represent the todo repository contract
type MongoTodoRepository interface {
	FindAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, error)
	CountFindAll(ctx context.Context, keyword string) (int, error)
	FindById(ctx context.Context, id string) (*models.Todo, error)
	CountFindByID(ctx context.Context, id string) (int, error)
	Store(ctx context.Context, value *models.Todo) (*models.Todo, error)
	Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, id string) error
}

type MongoTodoRepositoryImpl struct {
	mongoDB pkg_mongodb.MongoDB
}

// NewMongoTodoRepository will create an object that represent the TodoRepository interface
func NewMongoTodoRepository(mongoDB pkg_mongodb.MongoDB) MongoTodoRepository {
	return &MongoTodoRepositoryImpl{
		mongoDB: mongoDB,
	}
}

// FindAll - find all todo
func (m *MongoTodoRepositoryImpl) FindAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, error) {
	var results []*models.Todo

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")
	cur, err := collection.Find(ctx, bson.M{"title": bson.M{"$regex": keyword, "$options": "i"}}, findOptions)
	if err != nil {
		return []*models.Todo{}, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Todo
		err := cur.Decode(&elem)
		if err != nil {
			return []*models.Todo{}, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return []*models.Todo{}, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

// CountFindAll - count find all todo
func (m *MongoTodoRepositoryImpl) CountFindAll(ctx context.Context, keyword string) (int, error) {
	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")

	total, err := collection.CountDocuments(ctx, bson.M{"title": bson.M{"$regex": keyword, "$options": "i"}})
	if err != nil {
		return int(total), err
	}

	return int(total), nil
}

// FindById - find todo by id
func (m *MongoTodoRepositoryImpl) FindById(ctx context.Context, id string) (*models.Todo, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("not found")
	}

	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")

	result := &models.Todo{}
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return result, errors.New("not found")
		}

		return result, err
	}

	return result, nil
}

// CountFindByID - find count todo by id
func (m *MongoTodoRepositoryImpl) CountFindByID(ctx context.Context, id string) (int, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("not found")
	}

	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")
	total, err := collection.CountDocuments(ctx, bson.M{"_id": docID})
	if err != nil {
		return 0, err
	}

	if total <= 0 {
		return 0, errors.New("not found")
	}

	return int(total), nil
}

// Store - store todo
func (m *MongoTodoRepositoryImpl) Store(ctx context.Context, value *models.Todo) (*models.Todo, error) {
	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")

	timeNow := utils.GetTimeNow()
	res, err := collection.InsertOne(ctx, bson.M{
		"title":       value.Title,
		"description": value.Description,
		"createdAt":   timeNow,
		"updatedAt":   timeNow,
	})
	if err != nil {
		return &models.Todo{}, err
	}

	result := &models.Todo{
		ID:          res.InsertedID.(primitive.ObjectID),
		Title:       value.Title,
		Description: value.Description,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	return result, nil
}

// Update - update todo by id
func (m *MongoTodoRepositoryImpl) Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("not found")
	}

	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")

	timeNow := utils.GetTimeNow()
	bsonValue := bson.D{
		{Key: "title", Value: value.Title},
		{Key: "description", Value: value.Description},
		{Key: "updatedAt", Value: timeNow},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.D{{Key: "$set", Value: bsonValue}})
	if err != nil {
		return nil, err
	}

	result := &models.Todo{
		ID: docID,
	}

	return result, nil
}

// Delete - delete todo by id
func (m *MongoTodoRepositoryImpl) Delete(ctx context.Context, id string) error {
	client := m.mongoDB.Get()
	collection := client.Database(os.Getenv("DB_NAME")).Collection("todo")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("not found")
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		return err
	}

	if result.DeletedCount <= 0 {
		return errors.New("not found")
	}

	return nil
}

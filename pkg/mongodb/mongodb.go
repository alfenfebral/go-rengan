package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type MongoDB interface {
	Get() *mongo.Client
}

type MongoDBImpl struct {
	client *mongo.Client
}

func NewMongoDB() (MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DB_URL")
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(uri)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Checking the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	logrus.Println("Mongo Database connected")

	return &MongoDBImpl{
		client: client,
	}, err
}

func (m *MongoDBImpl) Get() *mongo.Client {
	return m.client
}

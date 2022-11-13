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
	Disconnect() error
}

type MongoDBImpl struct {
	ctx    context.Context
	client *mongo.Client
}

func NewMongoDB() (MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DB_URL")
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor() // add mongo opentelemetry tracing
	opts.ApplyURI(uri)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		logrus.Error("cannot connect")
		return nil, err
	}

	// Checking the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	logrus.Println("Mongo Database connected")

	return &MongoDBImpl{
		ctx:    ctx,
		client: client,
	}, err
}

func (m *MongoDBImpl) Get() *mongo.Client {
	return m.client
}

func (m *MongoDBImpl) Disconnect() error {
	if err := m.client.Disconnect(m.ctx); err != nil {
		return err
	}

	return nil
}

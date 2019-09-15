package mongodb

import (
	"context"
	"fmt"
	"github.com/metajar/coreSNMP/internal/controller"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	resourceCollection string = "resources"
)

type MongoBackend struct {
	Host     string
	Username string
	Password string
	Database string
	Client   *mongo.Client
}

func (m *MongoBackend) Init(ctx context.Context) error {
	var uri string
	if m.Username != "" {
		uri = fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
			m.Username,
			m.Password,
			m.Host,
			m.Database,
		)
	} else {
		uri = fmt.Sprintf(`mongodb://%s/%s`,
			m.Host,
			m.Database, )
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	m.Client = client
	return nil
}

func (m *MongoBackend) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (m *MongoBackend) Put(ctx context.Context) error {
	return nil
}

func (m *MongoBackend) Get(ctx context.Context) error {
	return nil
}

func (m *MongoBackend) Update(ctx context.Context) error {
	return nil
}

func (m *MongoBackend) TestWrite(ctx context.Context, c controller.CoreSNMPResource) error {
	collection := m.Client.Database(m.Database).Collection(resourceCollection)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	fmt.Println(res.InsertedID, "Inserted.")
	return nil

}

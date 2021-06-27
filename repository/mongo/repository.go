package mongo

import (
	"context"
	"time"

	errs "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/igorcavalcanti/go_shortener/shortener"
)

type repository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repo := &repository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client

	return repo, nil
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (this *repository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	ret := &shortener.Redirect{}
	collection := this.client.Database(this.database).Collection("redirects")
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = shortener.ErrRedirectNotFound
		}
		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}
	return ret, nil
}

func (this *repository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	collection := this.client.Database(this.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"create_at": redirect.CreateAt,
	})
	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

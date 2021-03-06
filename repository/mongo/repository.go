package mongo

import (
	"context"
	"log"
	"time"

	"github.com/Falcer/go-login/model"
	"github.com/Falcer/go-login/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	mongoRepository struct {
		db *mongo.Client
	}
)

const databaseName = "users"

// NewMongoRepository function
func NewMongoRepository(url string) (repository.UserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return &mongoRepository{client}, nil
}

func (r *mongoRepository) GetAllUser() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.db.Database(databaseName).Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []model.User
	for cur.Next(ctx) {
		var result model.User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return users, nil
}
func (r *mongoRepository) AddUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.Database(databaseName).Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.Database(databaseName).Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

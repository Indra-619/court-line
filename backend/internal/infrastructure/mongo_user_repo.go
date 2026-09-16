package infrastructure

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/internal/models"
)

// MongoUserRepository implements repository.UserRepository on MongoDB.
type MongoUserRepository struct {
	coll *mongo.Collection
}

// NewMongoUserRepository builds a user repository bound to the given
// client and database name.
func NewMongoUserRepository(client *mongo.Client, dbName string) *MongoUserRepository {
	return &MongoUserRepository{coll: client.Database(dbName).Collection("users")}
}

func (r *MongoUserRepository) findBy(ctx context.Context, filter bson.M) (*entity.User, error) {
	var user models.User
	if err := r.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return userToEntity(&user), nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	objID, err := parseID(id)
	if err != nil {
		return nil, err
	}
	return r.findBy(ctx, bson.M{"_id": objID})
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.findBy(ctx, bson.M{"email": email})
}

func (r *MongoUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	return r.findBy(ctx, bson.M{"googleId": googleID})
}

func (r *MongoUserRepository) UpsertGoogleUser(ctx context.Context, user *entity.User) error {
	now := time.Now()
	filter := bson.M{"googleId": user.GoogleID}
	update := bson.M{
		"$set": bson.M{
			"email":     user.Email,
			"name":      user.Name,
			"picture":   user.Picture,
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"googleId":  user.GoogleID,
			"role":      user.Role,
			"createdAt": now,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoUserRepository) Create(ctx context.Context, user *entity.User) error {
	objID := primitive.NewObjectID()
	if user.ID != "" {
		parsed, err := parseID(user.ID)
		if err != nil {
			return err
		}
		objID = parsed
	}
	user.ID = objID.Hex()

	doc := models.User{
		ID:        objID,
		GoogleID:  user.GoogleID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	_, err := r.coll.InsertOne(ctx, doc)
	return err
}

func (r *MongoUserRepository) Update(ctx context.Context, user *entity.User) error {
	objID, err := parseID(user.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"googleId":  user.GoogleID,
			"email":     user.Email,
			"name":      user.Name,
			"picture":   user.Picture,
			"role":      user.Role,
			"updatedAt": user.UpdatedAt,
		},
	}

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

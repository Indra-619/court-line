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

// timeNow is a variable so tests can control clock-dependent filters.
var timeNow = time.Now

// MongoRefreshTokenRepository implements repository.RefreshTokenRepository
// on MongoDB. FindByHash filters expired and revoked records in the
// query so an unusable token is indistinguishable from an unknown one.
type MongoRefreshTokenRepository struct {
	coll *mongo.Collection
}

// NewMongoRefreshTokenRepository builds a refresh-token repository
// bound to the given client and database name.
func NewMongoRefreshTokenRepository(client *mongo.Client, dbName string) *MongoRefreshTokenRepository {
	return &MongoRefreshTokenRepository{coll: client.Database(dbName).Collection("refresh_tokens")}
}

func (r *MongoRefreshTokenRepository) Create(ctx context.Context, record *entity.RefreshToken) error {
	objID := primitive.NewObjectID()
	if record.ID != "" {
		parsed, err := parseID(record.ID)
		if err != nil {
			return err
		}
		objID = parsed
	}
	record.ID = objID.Hex()

	doc := models.RefreshToken{
		ID:        objID,
		UserID:    record.UserID,
		TokenHash: record.TokenHash,
		ExpiresAt: record.ExpiresAt,
		CreatedAt: record.CreatedAt,
		Revoked:   record.Revoked,
	}
	_, err := r.coll.InsertOne(ctx, doc)
	return err
}

func (r *MongoRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	filter := bson.M{
		"tokenHash": hash,
		"revoked":   false,
		"expiresAt": bson.M{"$gt": timeNow()},
	}

	var doc models.RefreshToken
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &entity.RefreshToken{
		ID:        doc.ID.Hex(),
		UserID:    doc.UserID,
		TokenHash: doc.TokenHash,
		ExpiresAt: doc.ExpiresAt,
		CreatedAt: doc.CreatedAt,
		Revoked:   doc.Revoked,
	}, nil
}

func (r *MongoRefreshTokenRepository) DeleteByHash(ctx context.Context, hash string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"tokenHash": hash})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *MongoRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}

// MongoRevokedTokenRepository implements repository.RevokedTokenRepository
// on MongoDB. Entries are dropped automatically by the TTL index on
// expiresAt once the blacklisted JWT has expired.
type MongoRevokedTokenRepository struct {
	coll *mongo.Collection
}

// NewMongoRevokedTokenRepository builds a revoked-token (jti blacklist)
// repository bound to the given client and database name.
func NewMongoRevokedTokenRepository(client *mongo.Client, dbName string) *MongoRevokedTokenRepository {
	return &MongoRevokedTokenRepository{coll: client.Database(dbName).Collection("revoked_tokens")}
}

func (r *MongoRevokedTokenRepository) Create(ctx context.Context, record *entity.RevokedToken) error {
	// Upsert on jti so repeated logouts of the same token stay cheap;
	// the TTL index on expiresAt removes entries once the JWT expires.
	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"jti": record.JTI},
		bson.M{"$set": bson.M{"jti": record.JTI, "expiresAt": record.ExpiresAt}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *MongoRevokedTokenRepository) Exists(ctx context.Context, jti string) (bool, error) {
	filter := bson.M{"jti": jti, "expiresAt": bson.M{"$gt": timeNow()}}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

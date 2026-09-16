package infrastructure

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/internal/models"
)

// MongoCourtRepository implements repository.CourtRepository on MongoDB.
type MongoCourtRepository struct {
	coll *mongo.Collection
}

// NewMongoCourtRepository builds a court repository bound to the given
// client and database name.
func NewMongoCourtRepository(client *mongo.Client, dbName string) *MongoCourtRepository {
	return &MongoCourtRepository{coll: client.Database(dbName).Collection("courts")}
}

func parseID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, repository.ErrInvalidID
	}
	return objID, nil
}

func (r *MongoCourtRepository) FindAll(ctx context.Context) ([]*entity.Court, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courts []models.Court
	if err := cursor.All(ctx, &courts); err != nil {
		return nil, err
	}

	out := make([]*entity.Court, 0, len(courts))
	for i := range courts {
		out = append(out, courtToEntity(&courts[i]))
	}
	return out, nil
}

func (r *MongoCourtRepository) FindByID(ctx context.Context, id string) (*entity.Court, error) {
	objID, err := parseID(id)
	if err != nil {
		return nil, err
	}

	var court models.Court
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&court); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return courtToEntity(&court), nil
}

func (r *MongoCourtRepository) Create(ctx context.Context, court *entity.Court) error {
	objID := primitive.NewObjectID()
	if court.ID != "" {
		parsed, err := parseID(court.ID)
		if err != nil {
			return err
		}
		objID = parsed
	}
	court.ID = objID.Hex()

	doc := models.Court{
		ID:           objID,
		Name:         court.Name,
		Type:         court.Type,
		Location:     court.Location,
		Description:  court.Description,
		PricePerHour: court.PricePerHour,
		ImageURL:     court.ImageURL,
		Facilities:   court.Facilities,
		IsAvailable:  court.IsAvailable,
	}
	_, err := r.coll.InsertOne(ctx, doc)
	return err
}

func (r *MongoCourtRepository) Update(ctx context.Context, court *entity.Court) error {
	objID, err := parseID(court.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":         court.Name,
			"type":         court.Type,
			"location":     court.Location,
			"description":  court.Description,
			"pricePerHour": court.PricePerHour,
			"imageUrl":     court.ImageURL,
			"facilities":   court.Facilities,
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

func (r *MongoCourtRepository) Delete(ctx context.Context, id string) error {
	objID, err := parseID(id)
	if err != nil {
		return err
	}

	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

package repository

import (
	"context"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
	domainRepo "github.com/your-username/book-lapangan/backend/internal/domain/repository"
	"github.com/your-username/book-lapangan/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type courtMongoRepository struct {
	collection *mongo.Collection
}

func NewCourtMongoRepository(db *mongo.Database) domainRepo.CourtRepository {
	return &courtMongoRepository{
		collection: db.Collection("courts"),
	}
}

func (r *courtMongoRepository) FindAll(ctx context.Context) ([]*entity.Court, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mCourts []models.Court
	if err := cursor.All(ctx, &mCourts); err != nil {
		return nil, err
	}

	var entities []*entity.Court
	for _, m := range mCourts {
		entities = append(entities, toCourtEntity(&m))
	}
	return entities, nil
}

func (r *courtMongoRepository) FindByID(ctx context.Context, id string) (*entity.Court, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var m models.Court
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&m)
	if err != nil {
		return nil, err
	}
	return toCourtEntity(&m), nil
}

func (r *courtMongoRepository) Create(ctx context.Context, court *entity.Court) error {
	m := fromCourtEntity(court)
	m.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	court.ID = m.ID.Hex()
	return nil
}

func (r *courtMongoRepository) Update(ctx context.Context, court *entity.Court) error {
	objID, err := primitive.ObjectIDFromHex(court.ID)
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
			"isAvailable":  court.IsAvailable,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *courtMongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func toCourtEntity(m *models.Court) *entity.Court {
	return &entity.Court{
		ID:           m.ID.Hex(),
		Name:         m.Name,
		Type:         m.Type,
		Location:     m.Location,
		Description:  m.Description,
		PricePerHour: m.PricePerHour,
		ImageURL:     m.ImageURL,
		Facilities:   m.Facilities,
		IsAvailable:  m.IsAvailable,
	}
}

func fromCourtEntity(e *entity.Court) *models.Court {
	var id primitive.ObjectID
	if e.ID != "" {
		id, _ = primitive.ObjectIDFromHex(e.ID)
	}
	return &models.Court{
		ID:           id,
		Name:         e.Name,
		Type:         e.Type,
		Location:     e.Location,
		Description:  e.Description,
		PricePerHour: e.PricePerHour,
		ImageURL:     e.ImageURL,
		Facilities:   e.Facilities,
		IsAvailable:  e.IsAvailable,
	}
}

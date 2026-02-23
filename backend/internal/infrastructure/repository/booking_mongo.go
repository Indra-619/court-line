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

type bookingMongoRepository struct {
	collection *mongo.Collection
}

func NewBookingMongoRepository(db *mongo.Database) domainRepo.BookingRepository {
	return &bookingMongoRepository{
		collection: db.Collection("bookings"),
	}
}

func (r *bookingMongoRepository) FindAll(ctx context.Context) ([]*entity.Booking, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mBookings []models.Booking
	if err := cursor.All(ctx, &mBookings); err != nil {
		return nil, err
	}

	var entities []*entity.Booking
	for _, m := range mBookings {
		entities = append(entities, toBookingEntity(&m))
	}
	return entities, nil
}

func (r *bookingMongoRepository) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var m models.Booking
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&m)
	if err != nil {
		return nil, err
	}
	return toBookingEntity(&m), nil
}

func (r *bookingMongoRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mBookings []models.Booking
	if err := cursor.All(ctx, &mBookings); err != nil {
		return nil, err
	}

	var entities []*entity.Booking
	for _, m := range mBookings {
		entities = append(entities, toBookingEntity(&m))
	}
	return entities, nil
}

func (r *bookingMongoRepository) Create(ctx context.Context, booking *entity.Booking) error {
	m := fromBookingEntity(booking)
	m.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	booking.ID = m.ID.Hex()
	return nil
}

func (r *bookingMongoRepository) Update(ctx context.Context, booking *entity.Booking) error {
	objID, err := primitive.ObjectIDFromHex(booking.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":    booking.Status,
			"updatedAt": booking.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func toBookingEntity(m *models.Booking) *entity.Booking {
	return &entity.Booking{
		ID:            m.ID.Hex(),
		CourtID:       m.CourtID.Hex(),
		UserID:        m.UserID.Hex(),
		CustomerName:  m.CustomerName,
		CustomerPhone: m.CustomerPhone,
		Date:          m.Date,
		StartTime:     m.StartTime,
		EndTime:       m.EndTime,
		TotalPrice:    m.TotalPrice,
		Status:        entity.BookingStatus(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func fromBookingEntity(e *entity.Booking) *models.Booking {
	var id, courtID, userID primitive.ObjectID
	if e.ID != "" {
		id, _ = primitive.ObjectIDFromHex(e.ID)
	}
	if e.CourtID != "" {
		courtID, _ = primitive.ObjectIDFromHex(e.CourtID)
	}
	if e.UserID != "" {
		userID, _ = primitive.ObjectIDFromHex(e.UserID)
	}
	return &models.Booking{
		ID:            id,
		CourtID:       courtID,
		UserID:        userID,
		CustomerName:  e.CustomerName,
		CustomerPhone: e.CustomerPhone,
		Date:          e.Date,
		StartTime:     e.StartTime,
		EndTime:       e.EndTime,
		TotalPrice:    e.TotalPrice,
		Status:        models.BookingStatus(e.Status),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

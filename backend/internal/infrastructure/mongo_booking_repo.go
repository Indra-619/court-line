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

// MongoBookingRepository implements repository.BookingRepository on MongoDB.
type MongoBookingRepository struct {
	coll *mongo.Collection
}

// NewMongoBookingRepository builds a booking repository bound to the
// given client and database name.
func NewMongoBookingRepository(client *mongo.Client, dbName string) *MongoBookingRepository {
	return &MongoBookingRepository{coll: client.Database(dbName).Collection("bookings")}
}

func (r *MongoBookingRepository) FindAll(ctx context.Context) ([]*entity.Booking, error) {
	return r.find(ctx, bson.M{})
}

func (r *MongoBookingRepository) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	objID, err := parseID(id)
	if err != nil {
		return nil, err
	}

	var booking models.Booking
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return bookingToEntity(&booking), nil
}

func (r *MongoBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	objID, err := parseID(userID)
	if err != nil {
		return nil, err
	}
	return r.find(ctx, bson.M{"userId": objID})
}

func (r *MongoBookingRepository) FindByCourtID(ctx context.Context, courtID string) ([]*entity.Booking, error) {
	objID, err := parseID(courtID)
	if err != nil {
		return nil, err
	}
	return r.find(ctx, bson.M{"courtId": objID})
}

func (r *MongoBookingRepository) FindActiveByCourtAndDate(ctx context.Context, courtID, date string) ([]*entity.Booking, error) {
	objID, err := parseID(courtID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"courtId": objID,
		"date":    date,
		"status": bson.M{"$in": []string{
			string(models.BookingStatusPending),
			string(models.BookingStatusConfirmed),
		}},
	}
	return r.find(ctx, filter)
}

func (r *MongoBookingRepository) find(ctx context.Context, filter bson.M) ([]*entity.Booking, error) {
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	out := make([]*entity.Booking, 0, len(bookings))
	for i := range bookings {
		out = append(out, bookingToEntity(&bookings[i]))
	}
	return out, nil
}

func (r *MongoBookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	objID := primitive.NewObjectID()
	if booking.ID != "" {
		parsed, err := parseID(booking.ID)
		if err != nil {
			return err
		}
		objID = parsed
	}
	courtObjID, err := parseID(booking.CourtID)
	if err != nil {
		return err
	}
	userObjID, err := parseID(booking.UserID)
	if err != nil {
		return err
	}
	booking.ID = objID.Hex()

	doc := models.Booking{
		ID:            objID,
		CourtID:       courtObjID,
		UserID:        userObjID,
		CustomerName:  booking.CustomerName,
		CustomerPhone: booking.CustomerPhone,
		Date:          booking.Date,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		TotalPrice:    booking.TotalPrice,
		Status:        models.BookingStatus(booking.Status),
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
	_, err = r.coll.InsertOne(ctx, doc)
	return err
}

func (r *MongoBookingRepository) Update(ctx context.Context, booking *entity.Booking) error {
	objID, err := parseID(booking.ID)
	if err != nil {
		return err
	}
	courtObjID, err := parseID(booking.CourtID)
	if err != nil {
		return err
	}
	userObjID, err := parseID(booking.UserID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"courtId":       courtObjID,
			"userId":        userObjID,
			"customerName":  booking.CustomerName,
			"customerPhone": booking.CustomerPhone,
			"date":          booking.Date,
			"startTime":     booking.StartTime,
			"endTime":       booking.EndTime,
			"totalPrice":    booking.TotalPrice,
			"status":        booking.Status,
			"updatedAt":     booking.UpdatedAt,
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

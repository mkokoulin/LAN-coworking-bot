package services

import (
	"context"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoworkingRegistrationMongo struct {
	col *mongo.Collection
}

func NewCoworkingRegistrationMongo(db *mongo.Database) (*CoworkingRegistrationMongo, error) {
	col := db.Collection("coworking_registrations")

	_, err := col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "chat_id", Value: 1}}},
		{Keys: bson.D{{Key: "telegram_user_id", Value: 1}}},
		{Keys: bson.D{{Key: "phone", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return nil, err
	}

	return &CoworkingRegistrationMongo{col: col}, nil
}

func (r *CoworkingRegistrationMongo) Create(ctx context.Context, reg types.CoworkingRegistration) error {
	now := time.Now().UTC()
	reg.CreatedAt = now
	reg.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, reg)
	return err
}

func (r *CoworkingRegistrationMongo) GetPendingByChatID(ctx context.Context, chatID int64) (*types.CoworkingRegistration, error) {
	var reg types.CoworkingRegistration
	err := r.col.FindOne(ctx, bson.M{
		"chat_id": chatID,
		"status":  types.RegistrationPending,
	}, options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})).Decode(&reg)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *CoworkingRegistrationMongo) GetLatestApprovedByPhone(ctx context.Context, phone string) (*types.CoworkingRegistration, error) {
	var reg types.CoworkingRegistration
	err := r.col.FindOne(ctx, bson.M{
		"phone":  phone,
		"status": types.RegistrationApproved,
	}, options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})).Decode(&reg)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *CoworkingRegistrationMongo) GetLatestApprovedByChatID(ctx context.Context, chatID int64) (*types.CoworkingRegistration, error) {
	var reg types.CoworkingRegistration

	err := r.col.FindOne(
		ctx,
		bson.M{
			"chat_id": chatID,
			"status":  types.RegistrationApproved,
		},
		options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	).Decode(&reg)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &reg, nil
}

func (r *CoworkingRegistrationMongo) UpdateStatusByChatID(
	ctx context.Context,
	chatID int64,
	status types.RegistrationStatus,
	adminID int64,
	comment string,
) error {
	set := bson.M{
		"status":        status,
		"admin_comment": comment,
		"updated_at":    time.Now().UTC(),
		"approved_by":   adminID,
	}
	if status == types.RegistrationApproved {
		now := time.Now().UTC()
		set["approved_at"] = now
	}

	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"chat_id": chatID, "status": types.RegistrationPending},
		bson.M{"$set": set},
	)
	return err
}
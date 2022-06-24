package account

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	CreateAccount(ctx context.Context, acc Account) (string, error)
	UpdateAccount(ctx context.Context, id string, acc Account) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context) ([]Account, error)
	//	ActivateAccount(ctx context.Context, id string) error
	//	DeactivateAccount(ctx context.Context, id string) error
}

type Account struct {
}

var (
	ErrAccountNotFound = errors.New("Account not found")
	ErrNotImplemented  = errors.New("Not implemented")
)

type accountService struct {
	collection mongo.Collection
}

func NewAccountService(c mongo.Collection) Service {
	return &accountService{
		collection: c,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, acc Account) (string, error) {
	res, err := s.collection.InsertOne(ctx, acc)
	if err != nil {
		return "", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	} else {
		return "", errors.New("Unable to mashall InsertID object to string")
	}
}

// TODO
func (s *accountService) UpdateAccount(ctx context.Context, id string, acc Account) (*Account, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := s.collection.UpdateByID(ctx, bson.M{"_id": oid}, acc)
	if err != nil {
		return nil, err
	}

	if res.ModifiedCount < 1 {
		return nil, ErrAccountNotFound
	}

	return s.GetAccount(ctx, id)
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.collection.FindOne(ctx, bson.M{"_id": oid})

	var acc Account
	if err = res.Decode(&acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

func (s *accountService) ListAccounts(ctx context.Context) ([]Account, error) {
	cur, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var result []Account

	for cur.Next(ctx) {
		var item Account
		err = cur.Decode(&item)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

package store

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/brunovale91/payment-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentStoreConfig struct {
	URL        string
	Database   string
	Collection string
}

type PaymentStore interface {

	// Create payment in data store and return the created payment
	CreatePayment(*types.Payment) (*types.Payment, error)

	// Update payment attributes in data store and return the update payment
	UpdatePayment(string, *types.PaymentAttributes) (*types.Payment, error)

	// Delete payment in data store
	DeletePayment(string) (bool, error)

	// Get payment from data store
	GetPayment(string) (*types.Payment, error)

	// Get slice of payments from data store
	GetPayments() ([]*types.Payment, error)
}

type PaymentStoreImpl struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewPaymentStore(config *PaymentStoreConfig) (PaymentStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URL))
	if err != nil {
		return nil, err
	}
	collection := client.Database(config.Database).Collection(config.Collection)
	return PaymentStoreImpl{
		client:     client,
		collection: collection,
	}, nil
}

func (s PaymentStoreImpl) CreatePayment(payment *types.Payment) (*types.Payment, error) {
	_, err := s.collection.InsertOne(context.Background(), paymentToDoc(payment))
	if err != nil {
		log.Printf("Error creating payment with id %s: %s", payment.Id, err.Error())
		return nil, err
	}
	return payment, nil
}

func (s PaymentStoreImpl) UpdatePayment(id string, attributes *types.PaymentAttributes) (*types.Payment, error) {
	updateDoc := bson.M{
		"$inc": bson.M{
			"Version": 1,
		},
		"$set": bson.M{
			"Attributes": attributesToDoc(attributes),
		},
	}

	elem := &bson.D{}
	err := s.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": id}, updateDoc).Decode(elem)
	if err != nil {
		if isNoDocuments(err.Error()) {
			return nil, nil
		}
		log.Printf("Error updating payment with id %s: %s", id, err.Error())
		return nil, err
	}

	payment := docToPayment(*elem)
	payment.Attributes = attributes
	payment.Version += 1
	return payment, nil
}

func (s PaymentStoreImpl) DeletePayment(id string) (bool, error) {
	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error deleting payment with id %s: %s", id, err.Error())
		return false, err
	}
	return result.DeletedCount == 1, nil
}

func (s PaymentStoreImpl) GetPayment(id string) (*types.Payment, error) {
	elem := &bson.D{}
	err := s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(elem)
	if err != nil {
		if isNoDocuments(err.Error()) {
			return nil, nil
		}
		log.Printf("Error fetching payment with id %s: %s", id, err.Error())
		return nil, err
	}
	return docToPayment(*elem), nil
}

func (s PaymentStoreImpl) GetPayments() ([]*types.Payment, error) {
	cursor, err := s.collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error fetching payments: %s", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	payments, err := decodePayments(cursor)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func decodePayments(cursor *mongo.Cursor) ([]*types.Payment, error) {
	payments := make([]*types.Payment, 0)
	for cursor.Next(context.Background()) {
		payment, err := decodePayment(cursor)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Error fetching payments: %s", err)
		return nil, err
	}
	return payments, nil
}

func decodePayment(cursor *mongo.Cursor) (*types.Payment, error) {
	elem := &bson.D{}
	err := cursor.Decode(elem)
	if err != nil {
		log.Printf("Error parsing payment: %s", err)
		return nil, err
	}
	return docToPayment(*elem), nil
}

func isNoDocuments(message string) bool {
	return message == "mongo: no documents in result"
}

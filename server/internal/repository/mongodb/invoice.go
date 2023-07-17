package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/avalonprod/invoicepro/server/internal/apperrors"
	"github.com/avalonprod/invoicepro/server/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepository struct {
	db *mongo.Collection
}

func NewInvoiceRepository(db *mongo.Database) *InvoiceRepository {
	return &InvoiceRepository{
		db: db.Collection(invoiceCollection),
	}
}

func (r *InvoiceRepository) Create(ctx context.Context, input model.Invoice) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	doc, err := r.db.InsertOne(nCtx, input)
	if err != nil {
		return "", err
	}

	id := doc.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *InvoiceRepository) GetById(ctx context.Context, userID string, id string) (model.Invoice, error) {
	var invoice model.Invoice
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Invoice{}, err
	}
	if err := r.db.FindOne(nCtx, bson.M{"_id": ObjectID, "userID": userID}).Decode(&invoice); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Invoice{}, apperrors.ErrDocumentNotFound
		}
		return model.Invoice{}, err
	}
	return invoice, err
}

func (r *InvoiceRepository) SetMarkedById(ctx context.Context, userID string, id string, value bool) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"userID": userID, "id": id}, bson.M{"$set": bson.M{"isMarked": value}})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}

	return nil
}

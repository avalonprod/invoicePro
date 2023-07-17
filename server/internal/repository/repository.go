package repository

import (
	"github.com/avalonprod/invoicepro/server/internal/domain/repository"
	"github.com/avalonprod/invoicepro/server/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	User    repository.UserRepository
	Invoice repository.InvoiceRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User:    mongodb.NewUserRepository(db),
		Invoice: mongodb.NewInvoiceRepository(db),
	}
}

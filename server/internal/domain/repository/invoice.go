package repository

import (
	"context"

	"github.com/avalonprod/invoicepro/server/internal/domain/model"

)

type InvoiceRepository interface {
	Create(ctx context.Context, input model.Invoice) (string, error)
	// GetAllInvoce(ctx context.Context, userID string) ([]model.Invoice, error)
	GetById(ctx context.Context, userID string, id string) (model.Invoice, error)
	SetMarkedById(ctx context.Context, userID string, id string, value bool) error
	// GetLastCreatedInvoce()
	// DeleteById(ctx context.Context, userID string, id string) error
}

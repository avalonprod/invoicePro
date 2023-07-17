package service

import (
	"context"
	"errors"

	"github.com/avalonprod/invoicepro/server/internal/apperrors"
	"github.com/avalonprod/invoicepro/server/internal/domain/model"
	"github.com/avalonprod/invoicepro/server/internal/domain/repository"
	"github.com/avalonprod/invoicepro/server/internal/domain/types"
)

type InvoiceService struct {
	repository repository.InvoiceRepository
}

func NewInvoiceService(repository repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		repository: repository,
	}
}

func (s *InvoiceService) Create(ctx context.Context, input types.InvoiceDTO) (string, error) {
	var invList []model.InvItem
	for _, item := range input.InvList {

		invList = append(invList, model.InvItem{
			Title:       item.Title,
			Description: item.Description,
			Rate:        item.Rate,
			Qty:         item.Qty,
			Amount:      s.GetAmountForInvoiceItem(item.Rate, item.Qty).Amount,
		})
	}

	invoice := model.Invoice{
		UserID:      input.UserID,
		InvTitle:    input.InvTitle,
		InvNum:      input.InvNum,
		CreatedDate: input.CreatedDate,
		Balance:     input.Balance,
		Notes:       input.Notes,
		Dispatch:    input.Dispatch,
		Discount:    input.Discount,
		ColorLine:   input.ColorLine,
		Currency:    input.Currency,
		From: model.From{
			Name:      input.From.Name,
			EmailFrom: input.From.EmailFrom,
			Address: model.Address{
				Street:    input.From.Address.Street,
				CityState: input.From.Address.CityState,
				ZipCode:   input.From.Address.ZipCode,
			},
			Phone:         input.From.Phone,
			BusinessPhone: input.From.BusinessPhone,
		},
		To: model.To{
			Name:    input.To.Name,
			EmailTo: input.To.EmailTo,
			Address: model.Address{
				Street:    input.To.Address.Street,
				CityState: input.To.Address.CityState,
				ZipCode:   input.To.Address.ZipCode,
			},
			Phone: input.To.Phone,
		},
		IsMarked: false,
		InvList:  invList,
	}
	id, err := s.repository.Create(ctx, invoice)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *InvoiceService) GetAmountForInvoiceItem(rate float64, qty float64) *types.InvItemAmountDTO {
	res := rate * qty

	return &types.InvItemAmountDTO{
		Amount: res,
	}
}

func (s *InvoiceService) GetById(ctx context.Context, userID string, id string) (types.InvoiceDTO, error) {
	invoice, err := s.repository.GetById(ctx, userID, id)
	if err != nil {

	}

	var invList []types.InvItem
	for _, item := range invoice.InvList {

		invList = append(invList, types.InvItem{
			Title:       item.Title,
			Description: item.Description,
			Rate:        item.Rate,
			Qty:         item.Qty,
			Amount:      item.Amount,
		})
	}

	res := types.InvoiceDTO{
		ID:          invoice.ID,
		InvTitle:    invoice.InvTitle,
		InvNum:      invoice.InvNum,
		CreatedDate: invoice.CreatedDate,
		Balance:     invoice.Balance,
		Notes:       invoice.Notes,
		Dispatch:    invoice.Dispatch,
		Discount:    invoice.Discount,
		ColorLine:   invoice.ColorLine,
		Currency:    invoice.Currency,
		From: types.From{
			Name:      invoice.From.Name,
			EmailFrom: invoice.From.EmailFrom,
			Address: types.Address{
				Street:    invoice.From.Address.Street,
				CityState: invoice.From.Address.CityState,
				ZipCode:   invoice.From.Address.ZipCode,
			},
			Phone:         invoice.From.Phone,
			BusinessPhone: invoice.From.BusinessPhone,
		},
		To: types.To{
			Name:    invoice.To.Name,
			EmailTo: invoice.To.EmailTo,
			Address: types.Address{
				Street:    invoice.To.Address.Street,
				CityState: invoice.To.Address.CityState,
				ZipCode:   invoice.To.Address.ZipCode,
			},
			Phone: invoice.To.Phone,
		},
		InvList: invList,
	}
	return res, nil
}

func (s *InvoiceService) SetMarkedById(ctx context.Context, userID string, id string, value bool) error {
	err := s.repository.SetMarkedById(ctx, userID, id, value)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return err
		}
		return err
	}
	return nil
}

// func (s *InvoiceService) GetAll(ctx context.Context, userID string) ([]types.InvoiceDTO, error)

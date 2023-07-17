package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avalonprod/invoicepro/server/internal/apperrors"
	"github.com/avalonprod/invoicepro/server/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initInvoiceRoutes(api *gin.RouterGroup) {
	invoice := api.Group("/invoice")
	{
		authenticated := invoice.Group("/", h.userIdentity)
		{
			authenticated.POST("/create", h.Create)
			authenticated.POST("/amount", h.getAmountForInvoiceItem)
			authenticated.GET("/:id", h.getById)
		}

	}
}

type InvoiceCreateInput struct {
	InvTitle    string    `json:"invTitle"`
	InvNum      int       `json:"invNum"`
	CreatedDate time.Time `json:"createdDate"`
	Balance     float64   `json:"balance"`
	Notes       string    `json:"notes"`
	Dispatch    bool      `json:"dispatch"`
	Discount    bool      `json:"discount"`
	ColorLine   string    `json:"colorLine"`
	Currency    string    `json:"currency"`
	From        From      `json:"from"`
	To          To        `json:"to"`
	InvList     []InvItem `json:"invList"`
}

type From struct {
	Name          string  `json:"name"`
	EmailFrom     string  `json:"emailFrom"`
	Address       Address `json:"address"`
	Phone         string  `json:"phone"`
	BusinessPhone string  `json:"businessPhone"`
}

type To struct {
	Name    string  `json:"name"`
	EmailTo string  `json:"emailTo"`
	Address Address `json:"address"`
	Phone   string  `json:"phone"`
}

type Address struct {
	Street    string `json:"street"`
	CityState string `json:"cityState"`
	ZipCode   string `json:"zipCode"`
}

type InvItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Qty         float64 `json:"qty"`
}

type InvItemAmountResponse struct {
	Amount float64 `json:"amount"`
}

type InvItemAmountInput struct {
	Rate float64 `json:"rate"`
	Qty  float64 `json:"qty"`
}

type InvoiceResponse struct {
	ID          string            `json:"id"`
	InvTitle    string            `json:"invTitle"`
	InvNum      int               `json:"invNum"`
	CreatedDate time.Time         `json:"createdDate"`
	Balance     float64           `json:"balance"`
	Notes       string            `json:"notes"`
	Dispatch    bool              `json:"dispatch"`
	Discount    bool              `json:"discount"`
	ColorLine   string            `json:"colorLine"`
	Currency    string            `json:"currency"`
	From        From              `json:"from"`
	To          To                `json:"to"`
	IsMarked    bool              `json:"isMarked"`
	InvList     []InvItemResponse `json:"invList"`
}

type InvItemResponse struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Qty         float64 `json:"qty"`
	Amount      float64 `json:"amount"`
}

type SetMarker struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

func (h *HandlerV1) Create(c *gin.Context) {
	var input InvoiceCreateInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	var invList []types.InvItem
	for _, item := range input.InvList {
		invList = append(invList, types.InvItem{
			Description: item.Description,
			Title:       item.Title,
			Rate:        item.Rate,
			Qty:         item.Qty,
		})
	}
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	invoice := types.InvoiceDTO{
		UserID:      userID,
		InvTitle:    input.InvTitle,
		InvNum:      input.InvNum,
		CreatedDate: input.CreatedDate,
		Balance:     input.Balance,
		Notes:       input.Notes,
		Dispatch:    input.Dispatch,
		Discount:    input.Discount,
		ColorLine:   input.ColorLine,
		Currency:    input.Currency,
		From: types.From{
			Name:      input.From.Name,
			EmailFrom: input.From.EmailFrom,
			Address: types.Address{
				Street:    input.From.Address.Street,
				CityState: input.From.Address.CityState,
				ZipCode:   input.From.Address.ZipCode,
			},
			Phone:         input.From.Phone,
			BusinessPhone: input.From.BusinessPhone,
		},
		To: types.To{
			Name:    input.To.Name,
			EmailTo: input.To.EmailTo,
			Address: types.Address{
				Street:    input.To.Address.Street,
				CityState: input.To.Address.CityState,
				ZipCode:   input.To.Address.ZipCode,
			},
			Phone: input.To.Phone,
		},
		InvList: invList,
	}

	id, err := h.service.InvoiceService.Create(c.Request.Context(), invoice)

	if err != nil {

		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (h *HandlerV1) getAmountForInvoiceItem(c *gin.Context) {
	var input InvItemAmountInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	res := h.service.InvoiceService.GetAmountForInvoiceItem(input.Rate, input.Qty)

	c.JSON(http.StatusOK, InvItemAmountResponse{
		Amount: res.Amount,
	})
}

func (h *HandlerV1) getById(c *gin.Context) {
	id := c.Param("id")
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	res, err := h.service.InvoiceService.GetById(c.Request.Context(), userID, id)

	var invList []InvItemResponse
	for _, item := range res.InvList {
		invList = append(invList, InvItemResponse{
			Description: item.Description,
			Title:       item.Title,
			Rate:        item.Rate,
			Qty:         item.Qty,
			Amount:      item.Amount,
		})
	}
	c.JSON(http.StatusOK, InvoiceResponse{
		ID:          res.ID,
		InvTitle:    res.InvTitle,
		InvNum:      res.InvNum,
		CreatedDate: res.CreatedDate,
		Balance:     res.Balance,
		Notes:       res.Notes,
		Dispatch:    res.Dispatch,
		Discount:    res.Discount,
		ColorLine:   res.ColorLine,
		Currency:    res.Currency,
		From: From{
			Name:      res.From.Name,
			EmailFrom: res.From.EmailFrom,
			Address: Address{
				Street:    res.From.Address.Street,
				CityState: res.From.Address.CityState,
				ZipCode:   res.From.Address.ZipCode,
			},
			Phone:         res.From.Phone,
			BusinessPhone: res.From.BusinessPhone,
		},
		To: To{
			Name:    res.To.Name,
			EmailTo: res.To.EmailTo,
			Address: Address{
				Street:    res.To.Address.Street,
				CityState: res.To.Address.CityState,
				ZipCode:   res.To.Address.ZipCode,
			},
			Phone: res.To.Phone,
		},
		InvList: invList,
	})
}

func (h *HandlerV1) SetMarkedById(c *gin.Context) {
	var input SetMarker

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	if err := h.service.InvoiceService.SetMarkedById(c.Request.Context(), userID, input.Id, input.Value); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}
	c.Status(http.StatusOK)
}

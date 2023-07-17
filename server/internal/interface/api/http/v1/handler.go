package v1

import (
	"time"

	"github.com/avalonprod/invoicepro/server/internal/domain/service"
	"github.com/avalonprod/invoicepro/server/pkg/auth/manager"
	"github.com/gin-gonic/gin"
)

type HandlerV1 struct {
	service         *service.Service
	JWTManager      *manager.JWTManager
	refreshTokenTTL time.Duration
}

func NewHandler(service *service.Service, JWTManager *manager.JWTManager, refreshTokenTTL time.Duration) *HandlerV1 {
	return &HandlerV1{
		service:         service,
		JWTManager:      JWTManager,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (h *HandlerV1) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initInvoiceRoutes(v1)
	}
}

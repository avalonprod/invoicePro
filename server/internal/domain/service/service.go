package service

import (
	"context"
	"time"

	"github.com/avalonprod/invoicepro/server/internal/config"
	"github.com/avalonprod/invoicepro/server/internal/domain/repository"
	"github.com/avalonprod/invoicepro/server/internal/domain/types"
	"github.com/avalonprod/invoicepro/server/pkg/auth/manager"
	"github.com/avalonprod/invoicepro/server/pkg/codegenerator"
	"github.com/avalonprod/invoicepro/server/pkg/email"
	"github.com/avalonprod/invoicepro/server/pkg/hash"
)

type UserServiceI interface {
	SignUp(ctx context.Context, input types.UserSignUpDTO) (types.VerificationCodeDTO, error)
	SignIn(ctx context.Context, input types.UserSignInDTO) (types.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (types.Tokens, error)
	Verify(ctx context.Context, email string, verificationCode string) error
}

type InvoiceServiceI interface {
	Create(ctx context.Context, input types.InvoiceDTO) (string, error)
	GetAmountForInvoiceItem(rate float64, qty float64) *types.InvItemAmountDTO
	GetById(ctx context.Context, userID string, id string) (types.InvoiceDTO, error)
	SetMarkedById(ctx context.Context, userID string, id string, value bool) error
}

type Service struct {
	UserService    UserServiceI
	InvoiceService InvoiceServiceI
}

type Deps struct {
	Hasher              *hash.Hasher
	UserRepository      repository.UserRepository
	InvoiceRepository   repository.InvoiceRepository
	JWTManager          *manager.JWTManager
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	VerificationCodeTTL time.Duration
	Sender              email.Sender
	EmailConfig         config.EmailConfig
	CodeGenerator       *codegenerator.CodeGenerator
}

func NewService(deps *Deps) *Service {
	emailService := NewEmailService(deps.Sender, deps.EmailConfig)
	return &Service{
		UserService:    NewUserService(deps.Hasher, deps.UserRepository, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService),
		InvoiceService: NewInvoiceService(deps.InvoiceRepository),
	}
}

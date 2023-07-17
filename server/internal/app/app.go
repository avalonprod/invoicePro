package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/avalonprod/invoicepro/server/internal/config"
	"github.com/avalonprod/invoicepro/server/internal/domain/service"
	httpRoutes "github.com/avalonprod/invoicepro/server/internal/interface/api/http"
	"github.com/avalonprod/invoicepro/server/internal/repository"
	"github.com/avalonprod/invoicepro/server/pkg/auth/manager"
	"github.com/avalonprod/invoicepro/server/pkg/codegenerator"
	"github.com/avalonprod/invoicepro/server/pkg/db/mongodb"
	"github.com/avalonprod/invoicepro/server/pkg/email/smtp"
	"github.com/avalonprod/invoicepro/server/pkg/hash"
	"github.com/avalonprod/invoicepro/server/pkg/logger"
)

const configDir = "configs"

func Run() {

	cfg, err := config.Init(configDir)

	if err != nil {
		logger.Errorf("error parse config. err: %v", err)
	}

	// -----
	hasher := hash.NewHasher(cfg.Auth.PasswordSalt)
	mongoClient, err := mongodb.NewConnection(cfg.MongoDB.Url, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		logger.Errorf("failed to create new mongo client. err: %v", err)
	}

	JWTManager, err := manager.NewJWTManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)
		return
	}
	codeGenerator := codegenerator.NewCodeGenerator()
	// -----

	mongodb := mongoClient.Database(cfg.MongoDB.DBName)
	repository := repository.NewRepository(mongodb)
	service := service.NewService(&service.Deps{
		UserRepository:      repository.User,
		InvoiceRepository:   repository.Invoice,
		Hasher:              hasher,
		JWTManager:          JWTManager,
		AccessTokenTTL:      cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:     cfg.Auth.JWT.RefreshTokenTTL,
		VerificationCodeTTL: cfg.Auth.VerificationCodeTTL,
		Sender:              emailSender,
		EmailConfig:         cfg.Email,
		CodeGenerator:       codeGenerator,
	})
	handler := httpRoutes.NewHandler(service, JWTManager, cfg.Auth.JWT.RefreshTokenTTL)

	srv := NewServer(cfg, handler.InitRoutes())
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	if err := srv.Shotdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %x", err)
	}
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Errorf("error disconnect to mongoClient. err: %v", err)
	}
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes << 20,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	port := strings.Replace(s.httpServer.Addr, ":", "", 1)

	logger.Infof("Server has ben started on port: %s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shotdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

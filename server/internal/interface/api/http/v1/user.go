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

func (h *HandlerV1) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.SignUp)
		users.POST("/sign-in", h.SignIn)
		users.POST("/verify", h.userVerify)
		users.GET("/auth/refresh", h.userRefresh)
		// authenticated := users.Group("/", h.userIdentity)
		// {

		// }
	}

}

type UserSignUpInput struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	CompanyName string `json:"companyName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type UserSignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserVerifyInput struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verificationCode"`
}

// type refreshTokenInput struct {
// 	Token string `json:"token" binding:"required"`
// }

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type verifyResponse struct {
	Email                       string        `json:"email"`
	VerificationCodeExpiresTime time.Duration `json:"verificationCodeExpiresTime"`
}

func (h *HandlerV1) SignUp(c *gin.Context) {
	var input UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	res, err := h.service.UserService.SignUp(c.Request.Context(), types.UserSignUpDTO{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		CompanyName: input.CompanyName,
		Email:       input.Email,
		Password:    input.Password,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectEmailFormat) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectPasswordFormat) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectUserData) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, verifyResponse{
		Email:                       res.Email,
		VerificationCodeExpiresTime: res.VerificationCodeExpiresTime,
	})
}

func (h *HandlerV1) SignIn(c *gin.Context) {
	var input UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	tokens, err := h.service.UserService.SignIn(c.Request.Context(), types.UserSignInDTO{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, int(h.refreshTokenTTL.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *HandlerV1) userRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.service.UserService.RefreshTokens(c.Request.Context(), refreshToken)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	c.SetCookie("refresh_token", res.RefreshToken, int(h.refreshTokenTTL.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *HandlerV1) userVerify(c *gin.Context) {
	var input UserVerifyInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	err := h.service.UserService.Verify(c.Request.Context(), input.Email, input.VerificationCode)

	if err != nil {
		if errors.Is(err, apperrors.ErrIncorrectVerificationCode) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrUserAlreadyVerifyed) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrVerificationCodeExpired) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.String(http.StatusOK, "succes")
}

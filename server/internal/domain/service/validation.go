package service

import (
	"regexp"

	"github.com/avalonprod/invoicepro/server/internal/apperrors"
)

func validateCredentials(email, password string) error {
	if email == "" {
		return apperrors.ErrIncorrectEmailFormat
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return apperrors.ErrIncorrectEmailFormat
	}

	if password == "" {
		return apperrors.ErrIncorrectPasswordFormat
	}

	if len(password) < 8 || !regexp.MustCompile(`[A-Z]+`).MatchString(password) || !regexp.MustCompile(`\d+`).MatchString(password) {
		return apperrors.ErrIncorrectPasswordFormat
	}
	return nil
}

func validateUserData(firstName string, lastName string, companyName string) error {
	if len(firstName) < 2 || len(lastName) < 2 || len(companyName) < 2 {
		return apperrors.ErrIncorrectUserData
	}
	return nil
}

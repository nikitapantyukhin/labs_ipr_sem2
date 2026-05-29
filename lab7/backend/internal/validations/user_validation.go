package validations

import (
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func ValidateUserData(ctx *gin.Context, wrapper *service_wrapper.Wrapper, email string, phoneNumber string, socialNetworkLink string) (bool, error) {
	emailNotExists, emailValidationError := validateEmail(ctx, wrapper, email)
	if emailValidationError != nil {
		return false, emailValidationError
	}

	phoneNumberNotExists, phoneNumberValidationError := validatePhoneNumber(ctx, wrapper, phoneNumber)
	if phoneNumberValidationError != nil {
		return false, phoneNumberValidationError
	}

	socialNetworkNotExists, socialNetworkValidationError := validateSocialNetwork(ctx, wrapper, socialNetworkLink)
	if socialNetworkValidationError != nil {
		return false, socialNetworkValidationError
	}

	return emailNotExists && phoneNumberNotExists && socialNetworkNotExists, nil
}

func validateEmail(ctx *gin.Context, wrapper *service_wrapper.Wrapper, email string) (bool, error) {
	emailExists, err := wrapper.Db.Queries.CheckIfEmailIsRegistered(ctx, email)
	if err != nil {
		return false, err
	}

	return !emailExists, nil
}

func validatePhoneNumber(ctx *gin.Context, wrapper *service_wrapper.Wrapper, phoneNumber string) (bool, error) {
	phoneNumberExists, err := wrapper.Db.Queries.CheckIfPhoneIsRegistered(ctx, phoneNumber)
	if err != nil {
		return false, err
	}

	return !phoneNumberExists, nil
}

func validateSocialNetwork(ctx *gin.Context, wrapper *service_wrapper.Wrapper, socialNetworkLink string) (bool, error) {
	socialNetworkExists, err := wrapper.Db.Queries.CheckIfSocialNetworkIsRegistered(ctx, socialNetworkLink)
	if err != nil {
		return false, err
	}

	return !socialNetworkExists, nil
}

package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/create_user"
	"sport_platform/application/models/shared"
	"sport_platform/internal/mapper"
	"sport_platform/internal/observability"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"
	"sport_platform/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUserHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request create_user.CreateUserRequest
	if err := ctx.ShouldBind(&request); err != nil {
		fmt.Printf("CreateUserHandler bind error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse request: %s", err),
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		fmt.Printf("CreateUserHandler validation error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid registration data",
		})
		return
	}

	isDataValid, validationError := validations.ValidateUserData(ctx, wrapper, request.Email, request.PhoneNumber, request.SocialNetworkLink)

	if validationError != nil {
		fmt.Println(validationError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	if !isDataValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User with entered data already exists",
		})
		return
	}

	var createParams db_queries.CreateUserParams

	paramsMappingError := mapper.Mapper{}.Map(
		&createParams,
		request,
		struct {
			Password []byte
			Role     string
		}{
			Password: wrapper.PasswordHandler.HashPassword(request.Password),
			Role:     shared.Student,
		},
	)

	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	user, err := wrapper.Db.Queries.CreateUser(ctx, createParams)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var userClaims claims.UserClaims

	claimsMappingError := mapper.Mapper{}.Map(&userClaims, user)
	if claimsMappingError != nil {
		fmt.Println(claimsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	accessToken, refreshToken, tokenGenerationError := wrapper.JwtHandler.GenerateJwtPair(userClaims, fmt.Sprintf("%d", user.ID))
	if tokenGenerationError != nil {
		fmt.Println(tokenGenerationError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response create_user.CreateUserResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			AccessToken  string
			RefreshToken string
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)

	if responseMappingError != nil {
		fmt.Println(responseMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	ctx.JSON(
		http.StatusCreated,
		response,
	)
	observability.RecordBusinessEvent("user_registered")
}

package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/create_club"
	"sport_platform/application/models/shared"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/minio_config"
	"sport_platform/internal/observability"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CreateClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	userClaims := claimsRaw.(claims.UserClaims)

	if userClaims.Role != shared.Teacher && userClaims.Role != shared.Admin {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	var request create_club.CreateClubRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	if err := validate.Struct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Validation failed: %s", err.Error()),
			"details": err.Error(),
		})
		return
	}

	if userClaims.Role == shared.Teacher {
		request.TeacherID = userClaims.ID
	}

	if request.TeacherID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "teacher_id is required",
		})
		return
	}

	existsSport, err := wrapper.Db.Queries.CheckSportTypeExists(ctx, request.SportTypeID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error checking sport type"})
		return
	}
	if !existsSport {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid SportTypeID provided"})
		return
	}

	existsEducation, err := wrapper.Db.Queries.CheckEducationLevelExistsByName(ctx, request.EducationLevelName)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error checking education level"})
		return
	}
	if !existsEducation {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid EducationLevelName provided"})
		return
	}

	var createParams db_queries.CreateClubParams

	paramsMappingError := mapper.Mapper{}.Map(&createParams, request)
	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	club, err := wrapper.Db.Queries.CreateClub(ctx, createParams)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse form"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse form"})
		return
	}

	files := form.File["attachments"]
	var uploadedFilesUrls []string

	for _, fileHeader := range files {
		var uploadParams db_queries.UploadAttachmentParams
		minioID, err := minio_config.UploadFile(ctx, wrapper.Minio, fileHeader, "clubs")
		if err != nil {
			fmt.Printf("Failed to upload file %s: %v", fileHeader.Filename, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		minioUrl := fmt.Sprintf("http://minio:9000/clubs/%s", minioID)
		paramsMappingError := mapper.Mapper{}.Map(
			&uploadParams,
			struct {
				ClubID        int64
				AttachmentUrl string
			}{
				ClubID:        club.ID,
				AttachmentUrl: minioUrl,
			},
		)

		if paramsMappingError != nil {
			fmt.Println(paramsMappingError)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unknown error",
			})
			return
		}
		if dbError := wrapper.Db.Queries.UploadAttachment(ctx, uploadParams); dbError != nil {
			fmt.Printf("Error happened uploading attachment: %s\n", dbError)

			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Something unusual happened",
				},
			)
			return
		}

		uploadedFilesUrls = append(uploadedFilesUrls, minioUrl)
	}

	var response create_club.CreateClubResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		club,
		struct {
			Attachments []string
		}{
			Attachments: uploadedFilesUrls,
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
	observability.RecordBusinessEvent("club_created")
}

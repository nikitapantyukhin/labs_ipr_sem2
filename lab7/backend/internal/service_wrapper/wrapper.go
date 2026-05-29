package service_wrapper

import (
	"sport_platform/application/models/claims"
	"sport_platform/internal/jwt"
	"sport_platform/internal/password"
	"sport_platform/internal/sqlc/db"

	"github.com/minio/minio-go/v7"
)

type Wrapper struct {
	Db              *db.DbClient
	PasswordHandler password.IPasswordHandler
	JwtHandler      jwt.IJwtHandler[claims.UserClaims]
	Minio           *minio.Client
}

func (w Wrapper) Close() error {
	return w.Db.Close()
}

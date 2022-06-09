package services

import (
	"database/sql"

	"github.com/maslow123/users/pkg/client"
	"github.com/maslow123/users/pkg/utils"
)

type Server struct {
	Jwt            utils.JwtWrapper
	DB             *sql.DB
	BalanceService client.BalanceServiceClient
}

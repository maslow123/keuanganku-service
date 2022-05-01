package services

import (
	"database/sql"

	"github.com/maslow123/transactions/pkg/client"
)

type Server struct {
	DB             *sql.DB
	PosService     client.PosServiceClient
	BalanceService client.BalanceServiceClient
}

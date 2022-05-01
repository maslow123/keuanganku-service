package services

import (
	"context"
	"log"
	"net/http"

	"github.com/maslow123/balance/pkg/pb"
)

func (s *Server) UpsertBalance(ctx context.Context, req *pb.UpsertBalanceRequest) (*pb.UpsertBalanceResponse, error) {
	if req.UserId == 0 {
		return genericUpsertBalanceResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Type > 1 && req.Type < 0 {
		return genericUpsertBalanceResponse(http.StatusBadRequest, "invalid-type")
	}

	q := `
		INSERT INTO balance (user_id, type, total)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, type)
		DO UPDATE SET total = EXCLUDED.total RETURNING id, total
	`

	row := s.DB.QueryRowContext(ctx, q,
		&req.UserId,
		&req.Type,
		&req.Total,
	)

	var lastInsertedId, currentBalance int32
	err := row.Scan(&lastInsertedId, &currentBalance)

	if err != nil {
		log.Println(err)
		return genericUpsertBalanceResponse(http.StatusInternalServerError, err.Error())
	}

	resp := &pb.UpsertBalanceResponse{
		Status:         http.StatusCreated,
		Error:          "",
		Id:             lastInsertedId,
		CurrentBalance: currentBalance,
	}
	return resp, nil
}

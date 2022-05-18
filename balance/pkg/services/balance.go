package services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/maslow123/balance/pkg/pb"
)

func (s *Server) UpsertBalance(ctx context.Context, req *pb.UpsertBalanceRequest) (*pb.UpsertBalanceResponse, error) {
	if req.UserId == 0 {
		return genericUpsertBalanceResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Type != 0 && req.Type != 1 {
		return genericUpsertBalanceResponse(http.StatusBadRequest, "invalid-type")
	}
	if req.Action != 0 && req.Action != 1 {
		return genericUpsertBalanceResponse(http.StatusBadRequest, "invalid-action")
	}

	q := `
		INSERT INTO balance (user_id, type, total)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, type)
		DO UPDATE SET 
	`
	if req.Action == 1 {
		q = fmt.Sprintf("%s total = balance.total - EXCLUDED.total RETURNING id, total", q)
	} else {
		q = fmt.Sprintf("%s total = balance.total + EXCLUDED.total RETURNING id, total", q)
	}

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

func (s *Server) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.GetUserBalanceResponse, error) {
	if req.UserId == 0 {
		return genericGetUserBalanceResponse(http.StatusBadRequest, "invalid-user-id")
	}

	q := `
		SELECT type, total FROM balance WHERE user_id = $1
	`

	rows, err := s.DB.QueryContext(ctx, q, req.UserId)
	if err != nil {
		log.Println(err)
		return genericGetUserBalanceResponse(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var balances []*pb.UserBalance

	for rows.Next() {
		var balance pb.UserBalance
		if err := rows.Scan(
			&balance.Type,
			&balance.Total,
		); err != nil {
			log.Println(err)
			return genericGetUserBalanceResponse(http.StatusInternalServerError, err.Error())
		}

		balances = append(balances, &balance)
	}

	if err := rows.Close(); err != nil {
		log.Println(err)
		return genericGetUserBalanceResponse(http.StatusInternalServerError, err.Error())
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return genericGetUserBalanceResponse(http.StatusInternalServerError, err.Error())
	}

	if len(balances) == 0 {
		return genericGetUserBalanceResponse(http.StatusNotFound, "user-balance-not-found")
	}

	resp := &pb.GetUserBalanceResponse{
		Status:   http.StatusOK,
		Error:    "",
		Balances: balances,
	}

	return resp, nil
}

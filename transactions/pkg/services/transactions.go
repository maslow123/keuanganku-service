package services

import (
	"context"
	"log"
	"net/http"

	"github.com/maslow123/transactions/pkg/pb"
)

func (s *Server) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	if req.UserId == 0 {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.PosId == 0 {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-pos-id")
	}
	if req.Total == 0 {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-total")
	}
	if req.Details == "" {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-details")
	}
	if req.ActionType != 0 && req.ActionType != 1 {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-action-type")
	}
	if req.Type != 0 && req.Type != 1 {
		return genericCreateTransactionResponse(http.StatusBadRequest, "invalid-type")
	}
	// check existing pos
	pos, err := s.PosService.PosDetail(req.PosId)
	if err != nil || pos.Status != int32(http.StatusOK) {
		log.Println(err)
		return genericCreateTransactionResponse(int(pos.Status), pos.Error)
	}

	// insert tx
	q := `
		INSERT INTO transactions
		(user_id, pos_id, total, details, type)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id
	`

	row := s.DB.QueryRowContext(ctx, q,
		&req.UserId,
		&req.PosId,
		&req.Total,
		&req.Details,
		&req.Type,
	)
	var lastInsertedId int
	err = row.Scan(&lastInsertedId)

	if err != nil {
		log.Println(err)
		return genericCreateTransactionResponse(http.StatusInternalServerError, err.Error())
	}

	// Update pos total
	updatePos, err := s.PosService.UpdateTotalPosByUser(req.PosId, req.Total)
	if err != nil || updatePos.Status != int32(http.StatusOK) {
		log.Println(err)
		return genericCreateTransactionResponse(int(pos.Status), pos.Error)
	}
	log.Printf("===== Pos %s currently has Rp.%d =====", pos.Pos.Name, updatePos.Total)

	// Update balance
	log.Println(req.ActionType, pos.Pos.Type)
	updateBalance, err := s.BalanceService.UpsertBalance(req.UserId, pos.Pos.Type, req.ActionType, req.Total)
	if err != nil || updateBalance.Status != int32(http.StatusCreated) {
		log.Println("error: ", err)
		return genericCreateTransactionResponse(int(updateBalance.Status), updateBalance.Error)
	}
	log.Printf("===== Balance %d currently has Rp.%d =====", updateBalance.Id, updateBalance.CurrentBalance)

	resp := &pb.CreateTransactionResponse{
		Status: http.StatusCreated,
		Error:  "",
		Id:     int32(lastInsertedId),
	}
	return resp, nil
}

func (s *Server) GetTransactionByUser(ctx context.Context, req *pb.GetTransactionListRequest) (*pb.GetTransactionListResponse, error) {
	if req.UserId == 0 {
		return genericGetTransactionListByUserResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Page == 0 {
		return genericGetTransactionListByUserResponse(http.StatusBadRequest, "invalid-page")
	}
	if req.Limit == 0 {
		return genericGetTransactionListByUserResponse(http.StatusBadRequest, "invalid-limit")
	}

	q := `
		SELECT 
			t.id, t.total, t.details, t.type,
			p."name" pos_name, p.type pos_type, p.total pos_total, p.color pos_color
		FROM transactions t
		LEFT JOIN pos p ON p.id = t.pos_id
		LEFT JOIN users u ON u.id = t.user_id
		WHERE u.id = $1 ORDER BY t.created_at
		LIMIT $2 OFFSET $3
	`

	offset := (req.Page - 1) * req.Limit
	rows, err := s.DB.QueryContext(ctx, q, req.UserId, req.Limit, offset)
	if err != nil {
		log.Println(err)
		return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var transactions []*pb.Transaction

	for rows.Next() {
		var transaction pb.Transaction
		var pos pb.Pos
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Total,
			&transaction.Details,
			&transaction.Type,

			&pos.Name,
			&pos.Type,
			&pos.Total,
			&pos.Color,
		); err != nil {
			log.Println(err)
			return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
		}

		transaction.Pos = &pos
		transactions = append(transactions, &transaction)
	}
	if err := rows.Close(); err != nil {
		return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
	}

	if err := rows.Err(); err != nil {
		return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
	}

	if len(transactions) == 0 {
		return genericGetTransactionListByUserResponse(http.StatusNotFound, "transaction-not-found")
	}

	resp := &pb.GetTransactionListResponse{
		Status:      http.StatusOK,
		Error:       "",
		Limit:       req.Limit,
		Page:        req.Page,
		Transaction: transactions,
	}

	return resp, nil
}

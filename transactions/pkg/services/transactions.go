package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

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
		(user_id, pos_id, total, details, type, action, created_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	dt := time.Unix(int64(req.Date), 0)
	currentDate := time.Now().Format("2006-01-02")
	selectedDate := dt.Format("2006-01-02")

	dt = time.Unix(int64(req.Date), 0)
	if currentDate == selectedDate {
		dt = time.Unix(int64(req.Date), 0).Local().UTC()
	}
	// Start transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return genericCreateTransactionResponse(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, q,
		&req.UserId,
		&req.PosId,
		&req.Total,
		&req.Details,
		&req.Type,
		&req.ActionType,
		&dt,
	)
	var lastInsertedId int
	err = row.Scan(&lastInsertedId)

	if err != nil {
		log.Println(err)
		return genericCreateTransactionResponse(http.StatusInternalServerError, err.Error())
	}

	// Update pos total
	updatePos, err := s.PosService.UpdateTotalPosByUser(req.PosId, req.Total, 0)
	if err != nil || updatePos.Status != int32(http.StatusOK) {
		log.Println(err)
		return genericCreateTransactionResponse(int(pos.Status), pos.Error)
	}
	log.Printf("===== Pos %s currently has Rp.%d =====", pos.Pos.Name, updatePos.Total)

	// Update balance
	updateBalance, err := s.BalanceService.UpsertBalance(req.UserId, req.Type, req.ActionType, req.Total)
	if err != nil || updateBalance.Status != int32(http.StatusCreated) {
		log.Println(err)
		return genericCreateTransactionResponse(int(updateBalance.Status), updateBalance.Error)
	}
	log.Printf("===== Balance %d currently has Rp.%d =====", updateBalance.Id, updateBalance.CurrentBalance)

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return genericCreateTransactionResponse(http.StatusInternalServerError, err.Error())
	}

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
	if req.Action != 0 && req.Action != 1 && req.Action != 2 {
		return genericGetTransactionListByUserResponse(http.StatusBadRequest, "invalid-type")
	}

	params := 1
	args := make([]interface{}, 0)
	args = append(args, req.UserId)
	q := `
		SELECT 
			t.id, t.total, t.details, t.type, t.created_at,
			p."name" pos_name, p.type pos_type, p.total pos_total, p.color pos_color
		FROM transactions t
		LEFT JOIN pos p ON p.id = t.pos_id
		LEFT JOIN users u ON u.id = t.user_id
		WHERE u.id = $1		
	`
	if req.Action != 2 {
		params++
		q = fmt.Sprintf("%s AND t.action = $%d", q, params)
		args = append(args, req.Action)
	}

	var startDate, endDate string
	if req.StartDate != 0 && req.EndDate != 0 {
		startDate = time.Unix(int64(req.StartDate), 0).Format("2006-01-02")
		endDate = time.Unix(int64(req.EndDate), 0).Format("2006-01-02")
	} else {
		startDate = time.Now().Format("2006-01-02")
		endDate = time.Now().Format("2006-01-02")
	}

	q = fmt.Sprintf("%s AND t.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", q, startDate, endDate)

	q = fmt.Sprintf(
		"%s ORDER BY t.created_at DESC, t.details ASC LIMIT $%d OFFSET $%d", q, params+1, params+2,
	)

	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		log.Println(err)
		return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var transactions []*pb.Transaction
	var createdAt time.Time

	for rows.Next() {
		var transaction pb.Transaction
		var pos pb.Pos
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Total,
			&transaction.Details,
			&transaction.Type,
			&createdAt,

			&pos.Name,
			&pos.Type,
			&pos.Total,
			&pos.Color,
		); err != nil {
			log.Println(err)
			return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
		}

		transaction.CreatedAt = int32(createdAt.Unix())
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

	// Get user total transaction by date
	q = fmt.Sprintf(`
		SELECT COALESCE(SUM(total), 0) as total_transaction FROM transactions 
		WHERE user_id = $1 AND action = $2  AND created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'
	`, startDate, endDate)
	row := s.DB.QueryRowContext(ctx, q, req.UserId, req.Action)
	var totalTransaction int32
	errTotalTx := row.Scan(&totalTransaction)
	if errTotalTx != nil {
		log.Println(errTotalTx)
		return genericGetTransactionListByUserResponse(http.StatusInternalServerError, err.Error())
	}

	resp := &pb.GetTransactionListResponse{
		Status:           http.StatusOK,
		Error:            "",
		Limit:            req.Limit,
		Page:             req.Page,
		Transaction:      transactions,
		TotalTransaction: totalTransaction,
	}

	return resp, nil
}

func (s *Server) DetailTransaction(ctx context.Context, req *pb.DetailTransactionRequest) (*pb.DetailTransactionResponse, error) {
	if req.UserId == 0 {
		return genericDetailTransactionResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Id == 0 {
		return genericDetailTransactionResponse(http.StatusBadRequest, "invalid-transaction-id")
	}

	q := `
		SELECT 
			t.id, t.total, t.details, t.type, t.created_at,
			p."name" pos_name, p.type pos_type, p.total pos_total, p.color pos_color
		FROM transactions t
		LEFT JOIN pos p ON p.id = t.pos_id
		LEFT JOIN users u ON u.id = t.user_id
		WHERE u.id = $1 AND t.id = $2		
	`
	var transaction pb.Transaction
	var pos pb.Pos
	var createdAt time.Time

	row := s.DB.QueryRowContext(ctx, q, req.UserId, req.Id)
	err := row.Scan(
		&transaction.Id,
		&transaction.Total,
		&transaction.Details,
		&transaction.Type,
		&createdAt,
		&pos.Name,
		&pos.Type,
		&pos.Total,
		&pos.Color,
	)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericDetailTransactionResponse(http.StatusNotFound, "transaction-not-found")
		}
		return genericDetailTransactionResponse(http.StatusInternalServerError, err.Error())
	}

	transaction.CreatedAt = int32(createdAt.Unix())
	transaction.Pos = &pos

	resp := &pb.DetailTransactionResponse{
		Status:      http.StatusOK,
		Error:       "",
		Transaction: &transaction,
	}

	log.Println(resp)

	return resp, nil
}

func (s *Server) DeleteTransactionByUser(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.DeleteTransactionResponse, error) {
	if req.Id == 0 {
		return genericDeleteTransactionResponse(http.StatusBadRequest, "invalid-transaction-id")
	}
	if req.UserId == 0 {
		return genericDeleteTransactionResponse(http.StatusBadRequest, "invalid-user-id")
	}

	q := `
		DELETE FROM transactions 
		WHERE id = $1 AND user_id = $2
		RETURNING pos_id, total, user_id, type
	`

	row := s.DB.QueryRowContext(ctx, q, req.Id, req.UserId)
	var posId, posTotal, userId, paymentType int32
	err := row.Scan(&posId, &posTotal, &userId, &paymentType)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericDeleteTransactionResponse(http.StatusNotFound, "transaction-not-found")
		}
		return genericDeleteTransactionResponse(http.StatusInternalServerError, err.Error())
	}

	// update pos total
	updatePos, err := s.PosService.UpdateTotalPosByUser(posId, posTotal, 1)
	if err != nil || updatePos.Status != int32(http.StatusOK) {
		log.Println(err)
		return genericDeleteTransactionResponse(int(updatePos.Status), updatePos.Error)
	}

	// Update balance
	updateBalance, err := s.BalanceService.UpsertBalance(userId, paymentType, 1, posTotal)
	if err != nil || updateBalance.Status != int32(http.StatusCreated) {
		log.Println(err)
		return genericDeleteTransactionResponse(int(updateBalance.Status), updateBalance.Error)
	}

	resp := &pb.DeleteTransactionResponse{
		Status: http.StatusOK,
		Error:  "",
	}
	return resp, nil
}

func (s *Server) GetPercentageExpenditure(ctx context.Context, req *pb.GetPercentageExpenditureRequest) (*pb.GetPercentageExpenditureResponse, error) {
	d := "2006-01-02"
	if req.UserId == 0 {
		return genericGetPercentageExpenditureResponse(http.StatusBadRequest, "invalid-user-id")
	}

	_, err := time.Parse(d, req.StartDate)
	if req.StartDate == "" || err != nil {
		return genericGetPercentageExpenditureResponse(http.StatusBadRequest, "invalid-start-date")
	}

	_, err = time.Parse(d, req.EndDate)
	if req.EndDate == "" || err != nil {
		return genericGetPercentageExpenditureResponse(http.StatusBadRequest, "invalid-end-date")
	}

	// Formula: ((yesterday_expenses - today_expenses) / yesterday_expenses) * 100
	q := `
		SELECT 
			today_expenditure, other_expenditure, 
			( CAST( (other_expenditure - today_expenditure) AS DECIMAL) / other_expenditure ) * 100 percentage
		FROM 
			(
				(
					SELECT SUM(total) AS today_expenditure, action
					FROM transactions
					WHERE action = 1 AND user_id = $1 AND created_at::date = $2
					GROUP BY action
				) te 
				JOIN (
						SELECT SUM(total) AS other_expenditure, action
						FROM transactions			
						WHERE action = 1 and user_id = $1 AND created_at::date = $3
						GROUP BY action
				) oe ON te.action = oe.action
			)
		GROUP BY today_expenditure, other_expenditure
	`

	row := s.DB.QueryRowContext(ctx, q, req.UserId, req.StartDate, req.EndDate)
	var todayExpenses, otherDayExpenses, percentage float32

	err = row.Scan(&todayExpenses, &otherDayExpenses, &percentage)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericGetPercentageExpenditureResponse(http.StatusNotFound, "transaction-not-found")
		}
		return genericGetPercentageExpenditureResponse(http.StatusInternalServerError, err.Error())
	}

	resp := &pb.GetPercentageExpenditureResponse{
		Status:           http.StatusOK,
		Error:            "",
		TodayExpenses:    todayExpenses,
		OtherDayExpenses: otherDayExpenses,
		Percentage:       percentage,
	}

	return resp, nil
}

package services

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/maslow123/pos/pkg/pb"
)

func (s *Server) CreatePos(ctx context.Context, req *pb.CreatePosRequest) (*pb.CreatePosResponse, error) {
	if req.UserId == 0 {
		return genericCreatePosResponse(http.StatusBadRequest, "invalid-user")
	}
	if req.Name == "" {
		return genericCreatePosResponse(http.StatusBadRequest, "invalid-name")
	}
	if req.Type != 0 && req.Type != 1 {
		return genericCreatePosResponse(http.StatusBadRequest, "invalid-type")
	}
	if req.Color == "" {
		return genericCreatePosResponse(http.StatusBadRequest, "invalid-color")
	}

	q := `
		INSERT INTO pos (user_id, name, type, color)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := s.DB.QueryRowContext(ctx, q,
		&req.UserId,
		&req.Name,
		&req.Type,
		&req.Color,
	)

	var lastInsertedId int64

	err := row.Scan(&lastInsertedId)
	if err != nil {
		log.Println(err)
		return genericCreatePosResponse(http.StatusInternalServerError, err.Error())
	}

	resp := &pb.CreatePosResponse{
		Status: http.StatusCreated,
		Id:     lastInsertedId,
		Error:  "",
	}

	return resp, nil
}

func (s *Server) GetPosByUser(ctx context.Context, req *pb.GetPosListRequest) (*pb.GetPosListResponse, error) {
	if req.UserId == 0 {
		return genericListPosByUserResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Limit == 0 {
		return genericListPosByUserResponse(http.StatusBadRequest, "invalid-limit")
	}
	if req.Page == 0 {
		return genericListPosByUserResponse(http.StatusBadRequest, "invalid-page")
	}

	q := `
		SELECT id, name, type, total, color
		FROM pos
		WHERE user_id = $1
		LIMIT $2
		OFFSET $3
	`

	offset := (req.Page - 1) * req.Limit
	rows, err := s.DB.QueryContext(ctx, q, req.UserId, req.Limit, offset)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericListPosByUserResponse(http.StatusNotFound, "pos-not-found")
		}
		return genericListPosByUserResponse(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var pos []*pb.Pos

	for rows.Next() {
		var p pb.Pos
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Type,
			&p.Total,
			&p.Color,
		); err != nil {
			log.Println(err)
			return genericListPosByUserResponse(http.StatusInternalServerError, err.Error())
		}

		pos = append(pos, &p)
	}

	if err := rows.Close(); err != nil {
		return genericListPosByUserResponse(http.StatusInternalServerError, err.Error())
	}

	if err := rows.Err(); err != nil {
		return genericListPosByUserResponse(http.StatusInternalServerError, err.Error())
	}

	if len(pos) == 0 {
		return genericListPosByUserResponse(http.StatusNotFound, "pos-not-found")
	}

	resp := &pb.GetPosListResponse{
		Status: http.StatusOK,
		Error:  "",
		Limit:  req.Limit,
		Page:   req.Page,
		Pos:    pos,
	}
	return resp, nil
}

func (s *Server) PosDetail(ctx context.Context, req *pb.PosDetailRequest) (*pb.PosDetailResponse, error) {
	if req.Id == 0 {
		return genericPosDetailResponse(http.StatusBadRequest, "invalid-id")
	}
	q := `
		SELECT id, name, type, total, color, created_at, updated_at
		FROM pos
		WHERE id = $1
	`
	var pos pb.Pos
	var createdAt, updatedAt time.Time

	row := s.DB.QueryRowContext(ctx, q, req.Id)
	err := row.Scan(
		&pos.Id,
		&pos.Name,
		&pos.Type,
		&pos.Total,
		&pos.Color,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericPosDetailResponse(http.StatusNotFound, "pos-not-found")
		}
		return genericPosDetailResponse(http.StatusInternalServerError, err.Error())
	}

	pos.CreatedAt = createdAt.Unix()
	pos.UpdatedAt = updatedAt.Unix()

	resp := &pb.PosDetailResponse{
		Status: http.StatusOK,
		Pos:    &pos,
		Error:  "",
	}
	return resp, nil
}

func (s *Server) UpdatePosByUser(ctx context.Context, req *pb.UpdatePosRequest) (*pb.UpdatePosResponse, error) {
	if req.Id == 0 {
		return genericUpdatePosByUserResponse(http.StatusBadRequest, "invalid-id")
	}
	if req.Name == "" {
		return genericUpdatePosByUserResponse(http.StatusBadRequest, "invalid-name")
	}
	if req.Color == "" {
		return genericUpdatePosByUserResponse(http.StatusBadRequest, "invalid-color")
	}

	q := `
		UPDATE pos
		SET name = $2, color = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, type, total, color, created_at, updated_at	
	`

	row := s.DB.QueryRowContext(ctx, q,
		&req.Id,
		&req.Name,
		&req.Color,
	)
	var p pb.Pos
	var createdAt, updatedAt time.Time
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.Type,
		&p.Total,
		&p.Color,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericUpdatePosByUserResponse(http.StatusNotFound, "pos-not-found")
		}
		return genericUpdatePosByUserResponse(http.StatusInternalServerError, err.Error())
	}

	p.CreatedAt = createdAt.Unix()
	p.UpdatedAt = updatedAt.Unix()

	resp := &pb.UpdatePosResponse{
		Status: http.StatusOK,
		Error:  "",
		Pos:    &p,
	}

	return resp, nil
}

func (s *Server) DeletePosByUser(ctx context.Context, req *pb.DeletePosRequest) (*pb.DeletePosResponse, error) {
	if req.Id == 0 {
		return genericDeletePosByUserResponse(http.StatusBadRequest, "invalid-id")
	}

	q := `DELETE FROM pos WHERE id = $1`

	res, err := s.DB.ExecContext(ctx, q, req.Id)
	if err != nil {
		log.Println(err)
		return genericDeletePosByUserResponse(http.StatusInternalServerError, err.Error())
	}
	count, err := res.RowsAffected()
	if err == nil && count == 0 {
		return genericDeletePosByUserResponse(http.StatusNotFound, "pos-not-found")
	}

	return genericDeletePosByUserResponse(http.StatusOK, "")
}
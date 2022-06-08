package services

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/maslow123/users/pkg/pb"
	"github.com/maslow123/users/pkg/utils"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Name == "" {
		return genericRegisterResponse(http.StatusBadRequest, "invalid-name")
	}
	if req.Email == "" {
		return genericRegisterResponse(http.StatusBadRequest, "invalid-email")
	}
	if req.Password == "" {
		return genericRegisterResponse(http.StatusBadRequest, "invalid-password")
	}
	if req.ConfirmPassword == "" {
		return genericRegisterResponse(http.StatusBadRequest, "invalid-confirm-password")
	}
	if req.Password != req.ConfirmPassword {
		return genericRegisterResponse(http.StatusBadRequest, "password-not-match")
	}

	hashedPassword := utils.HashPassword(req.Password)

	// Check email already exists
	q := `SELECT email FROM users where email = $1`

	row := s.DB.QueryRowContext(ctx, q, req.Email)
	var email string

	_ = row.Scan(&email)

	if email != "" {
		return genericRegisterResponse(http.StatusBadRequest, "email-already-exists")
	}

	q = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	row = s.DB.QueryRowContext(ctx, q,
		&req.Name,
		&req.Email,
		&hashedPassword,
	)

	var lastInsertedId int32

	err := row.Scan(&lastInsertedId)
	if err != nil {
		log.Println(err)
		return genericRegisterResponse(http.StatusInternalServerError, err.Error())
	}

	return genericRegisterResponse(http.StatusCreated, "")
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" {
		return genericLoginResponse(http.StatusBadRequest, "invalid-email")
	}
	if req.Password == "" {
		return genericLoginResponse(http.StatusBadRequest, "invalid-password")
	}

	var user pb.User
	var userPass string
	q := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	row := s.DB.QueryRowContext(ctx, q, req.Email)

	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&userPass,
	)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return genericLoginResponse(http.StatusNotFound, "user-not-found")
		}
		return genericLoginResponse(http.StatusInternalServerError, err.Error())
	}

	match := utils.CheckPasswordHash(req.Password, userPass)
	if !match {
		return genericLoginResponse(http.StatusUnauthorized, "password-not-match")
	}

	token, _ := s.Jwt.GenerateToken(user.Id)
	resp := &pb.LoginResponse{
		Status: http.StatusOK,
		Error:  "",
		User:   &user,
		Token:  token,
	}

	return resp, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: claims.UserId,
	}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.Id == 0 {
		return genericUpdateProfileResponse(http.StatusBadRequest, "invalid-user-id")
	}
	if req.Name == "" {
		return genericUpdateProfileResponse(http.StatusBadRequest, "invalid-name")
	}
	if req.Email == "" {
		return genericUpdateProfileResponse(http.StatusBadRequest, "invalid-email")
	}

	q := `
		UPDATE users SET email = $2, name = $3
		WHERE id = $1
	`

	res, err := s.DB.ExecContext(ctx, q, req.Id, req.Email, req.Name)
	if err != nil {
		log.Println(err)
		return genericUpdateProfileResponse(http.StatusInternalServerError, err.Error())
	}
	count, err := res.RowsAffected()
	if err == nil && count == 0 {
		return genericUpdateProfileResponse(http.StatusNotFound, "user-not-found")
	}

	return genericUpdateProfileResponse(http.StatusOK, "")
}

package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type CreateTransactionRequest struct {
	PosId      int32  `json:"pos_id"`
	Total      int32  `json:"total"`
	Details    string `json:"details"`
	ActionType int32  `json:"action_type"`
	Type       int32  `json:"type"`
}

func CreateTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	req := CreateTransactionRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userID := ctx.Value("user_id").(int32)
	request := &pb.CreateTransactionRequest{
		UserId:     int32(userID),
		PosId:      req.PosId,
		Total:      req.Total,
		Details:    req.Details,
		ActionType: req.ActionType,
		Type:       req.Type,
	}
	log.Println(request)
	res, err := c.CreateTransaction(context.Background(), request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int32(http.StatusCreated) {
		ctx.JSON(int(res.Status), res.Error)
		return
	}
	utils.SendProtoMessage(ctx, res, http.StatusCreated)
}

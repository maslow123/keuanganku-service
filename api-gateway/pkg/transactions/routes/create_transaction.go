package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type CreateTransactionRequest struct {
	UserId  int64  `json:"user_id"`
	PosId   int64  `json:"pos_id"`
	Total   int64  `json:"total"`
	Details string `json:"details"`
}

func CreateTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	req := CreateTransactionRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := c.CreateTransaction(context.Background(), &pb.CreateTransactionRequest{
		UserId:  req.UserId,
		PosId:   req.PosId,
		Total:   req.Total,
		Details: req.Details,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int64(http.StatusCreated) {
		ctx.JSON(int(res.Status), res.Error)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

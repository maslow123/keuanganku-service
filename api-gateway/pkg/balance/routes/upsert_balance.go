package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type UpsertBalanceRequest struct {
	Type   int32 `json:"type"`
	Total  int32 `json:"total"`
	Action int32 `json:"action"`
}

func UpsertBalance(ctx *gin.Context, c pb.BalanceServiceClient) {
	req := UpsertBalanceRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userID := ctx.Value("user_id").(int32)
	res, err := c.UpsertBalance(context.Background(), &pb.UpsertBalanceRequest{
		UserId: userID,
		Type:   req.Type,
		Total:  req.Total,
		Action: pb.UpsertBalanceRequest_ActionType(req.Action),
	})

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int32(http.StatusCreated) {
		ctx.JSON(int(res.Status), res)
		return
	}

	ctx.JSON(int(res.Status), &res)
}

package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type CreatePosRequest struct {
	Name  string `json:"name"`
	Type  int32  `json:"type"`
	Color string `json:"color"`
}

func CreatePos(ctx *gin.Context, c pb.PosServiceClient) {
	req := CreatePosRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userID := ctx.Value("user_id").(int32)

	res, err := c.CreatePos(context.Background(), &pb.CreatePosRequest{
		UserId: userID,
		Name:   req.Name,
		Type:   req.Type,
		Color:  req.Color,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int32(http.StatusCreated) {
		ctx.JSON(int(res.Status), res)
		return
	}

	ctx.JSON(int(res.Status), &res)
}

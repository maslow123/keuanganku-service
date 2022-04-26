package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type UpdatePosRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func UpdatePosByUser(ctx *gin.Context, c pb.PosServiceClient) {
	var req UpdatePosRequest

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := c.UpdatePosByUser(context.Background(), &pb.UpdatePosRequest{
		Id:    id,
		Name:  req.Name,
		Color: req.Color,
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	if res.Status != int64(http.StatusOK) {
		ctx.JSON(int(res.Status), res.Error)
		return
	}
	ctx.JSON(int(res.Status), &res)
}

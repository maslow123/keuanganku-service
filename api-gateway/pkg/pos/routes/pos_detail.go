package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/pos/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func PosDetail(ctx *gin.Context, c pb.PosServiceClient) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := c.PosDetail(context.Background(), &pb.PosDetailRequest{
		Id: int64(id),
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	if res.Status != int64(http.StatusCreated) {
		ctx.JSON(int(res.Status), res.Error)
		return
	}

	ctx.JSON(int(res.Status), &res)
}

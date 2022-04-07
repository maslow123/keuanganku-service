package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/pos/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func DeletePosByUser(ctx *gin.Context, c pb.PosServiceClient) {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	res, err := c.DeletePosByUser(context.Background(), &pb.DeletePosRequest{
		Id: id,
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

package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func GetPosList(ctx *gin.Context, c pb.PosServiceClient) {

	limitString := ctx.Query("limit")
	pageString := ctx.Query("page")
	typeString := ctx.Query("type")
	userID := ctx.Value("user_id").(int32)

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	parsingType, err := strconv.Atoi(typeString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	res, err := c.GetPosByUser(context.Background(), &pb.GetPosListRequest{
		UserId: userID,
		Limit:  int32(limit),
		Page:   int32(page),
		Type:   int32(parsingType),
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	if res.Status != int32(http.StatusOK) {
		ctx.JSON(int(res.Status), res)
		return
	}
	utils.SendProtoMessage(ctx, res)
}

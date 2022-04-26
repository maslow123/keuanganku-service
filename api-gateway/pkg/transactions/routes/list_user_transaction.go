package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func GetUserTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	userIDString := ctx.Query("user_id")
	limitString := ctx.Query("limit")
	pageString := ctx.Query("page")

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

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

	res, err := c.GetTransactionByUser(context.Background(), &pb.GetTransactionListRequest{
		UserId: int64(userID),
		Limit:  int64(limit),
		Page:   int64(page),
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
package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func DetailUserTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	transactionId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}
	userID := ctx.Value("user_id").(int32)

	res, err := c.DetailTransaction(context.Background(), &pb.DetailTransactionRequest{
		Id:     int32(transactionId),
		UserId: userID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int32(http.StatusOK) {
		ctx.JSON(int(res.Status), res)
		return
	}

	log.Println(res)
	utils.SendProtoMessage(ctx, res, http.StatusOK)
}

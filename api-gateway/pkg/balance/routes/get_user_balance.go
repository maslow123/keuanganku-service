package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func GetUserBalance(ctx *gin.Context, c pb.BalanceServiceClient) {
	userID := ctx.Value("user_id").(int32)
	res, err := c.GetUserBalance(context.Background(), &pb.GetUserBalanceRequest{
		UserId: userID,
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

	utils.SendProtoMessage(ctx, res, http.StatusOK)
}

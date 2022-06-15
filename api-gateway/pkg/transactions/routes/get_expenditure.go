package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func GetPercentageExpenditure(ctx *gin.Context, c pb.TransactionServiceClient) {
	startDateString := ctx.Query("start_date")
	endDateString := ctx.Query("end_date")

	userID := ctx.Value("user_id").(int32)

	res, err := c.GetPercentageExpenditure(context.Background(), &pb.GetPercentageExpenditureRequest{
		UserId:    userID,
		StartDate: startDateString,
		EndDate:   endDateString,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if res.Status != int32(http.StatusOK) {
		utils.SendProtoMessage(ctx, res, int(res.Status))
		return
	}

	log.Println(res)
	utils.SendProtoMessage(ctx, res, int(res.Status))
}

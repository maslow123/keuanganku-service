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

func GetUserTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	limitString := ctx.Query("limit")
	pageString := ctx.Query("page")
	actionString := ctx.Query("action")
	startDateString := ctx.Query("start_date")
	endDateString := ctx.Query("end_date")

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

	action, err := strconv.Atoi(actionString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	startDate, err := strconv.Atoi(startDateString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	endDate, err := strconv.Atoi(endDateString)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	userID := ctx.Value("user_id").(int32)

	res, err := c.GetTransactionByUser(context.Background(), &pb.GetTransactionListRequest{
		UserId:    userID,
		Limit:     int32(limit),
		Page:      int32(page),
		Action:    int32(action),
		StartDate: int32(startDate),
		EndDate:   int32(endDate),
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

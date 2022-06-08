package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/users/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

type UpdateProfileBody struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func UpdateProfile(ctx *gin.Context, c pb.UserServiceClient) {
	userID := ctx.Value("user_id").(int32)
	req := UpdateProfileBody{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := c.UpdateProfile(context.Background(), &pb.UpdateProfileRequest{
		Id:    userID,
		Email: req.Email,
		Name:  req.Name,
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
		return
	}

	utils.SendProtoMessage(ctx, res, int(res.Status))
}

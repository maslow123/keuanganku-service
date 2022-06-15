package routes

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/users/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
)

func UploadImage(ctx *gin.Context, c pb.UserServiceClient) {
	userID := ctx.Value("user_id").(int32)
	stream, err := c.UploadImage(context.Background())

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
	}

	imageType := filepath.Ext(header.Filename)
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				UserId:    userID,
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
			break
		}
		size += n
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
			break
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err))
	}

	utils.SendProtoMessage(ctx, res, http.StatusOK)
}

package pos

import (
	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/pos/routes"
	"github.com/maslow123/api-gateway/pkg/users"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, userSvc *users.ServiceClient) *ServiceClient {
	a := users.InitAuthMiddleware(userSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
		Router: r,
	}

	routes := r.Group("/pos")
	routes.Use(a.AuthRequired)
	routes.POST("/create", svc.CreatePos)
	routes.GET("/list", svc.GetPosList)
	routes.GET("/:id", svc.PosDetail)
	routes.PUT("/:id", svc.UpdatePosByUser)
	routes.DELETE("/:id", svc.DeletePosByUser)

	return svc
}

func (svc *ServiceClient) CreatePos(ctx *gin.Context) {
	routes.CreatePos(ctx, svc.Client)
}

func (svc *ServiceClient) GetPosList(ctx *gin.Context) {
	routes.GetPosList(ctx, svc.Client)
}

func (svc *ServiceClient) PosDetail(ctx *gin.Context) {
	routes.PosDetail(ctx, svc.Client)
}

func (svc *ServiceClient) UpdatePosByUser(ctx *gin.Context) {
	routes.UpdatePosByUser(ctx, svc.Client)
}

func (svc *ServiceClient) DeletePosByUser(ctx *gin.Context) {
	routes.DeletePosByUser(ctx, svc.Client)
}

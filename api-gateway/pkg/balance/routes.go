package pos

import (
	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/balance/routes"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/users"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, userSvc *users.ServiceClient) *ServiceClient {
	a := users.InitAuthMiddleware(userSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
		Router: r,
	}
	r.Use(a.CORSMiddleware)
	routes := r.Group("/balance")
	routes.Use(a.AuthRequired)
	routes.POST("/upsert", svc.UpsertBalance)

	return svc
}

func (svc *ServiceClient) UpsertBalance(ctx *gin.Context) {
	routes.UpsertBalance(ctx, svc.Client)
}

package users

import (
	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/users/routes"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {

	svc := &ServiceClient{
		Client: InitServiceClient(c),
		Router: r,
	}
	a := InitAuthMiddleware(svc)

	routes := r.Group("/users")
	routes.Use(a.CORSMiddleware)
	routes.POST("/register", svc.Register)
	routes.POST("/login", svc.Login)

	routes.Use(a.AuthRequired)
	routes.PUT("/update", svc.UpdateProfile)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	routes.Login(ctx, svc.Client)
}

func (svc *ServiceClient) UpdateProfile(ctx *gin.Context) {
	routes.UpdateProfile(ctx, svc.Client)
}

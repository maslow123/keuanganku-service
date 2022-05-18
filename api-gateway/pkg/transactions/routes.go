package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/transactions/routes"
	"github.com/maslow123/api-gateway/pkg/users"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, userSvc *users.ServiceClient) *ServiceClient {
	a := users.InitAuthMiddleware(userSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
		Router: r,
	}

	r.Use(a.CORSMiddleware)
	routes := r.Group("/transactions")
	routes.Use(a.AuthRequired)
	routes.POST("/create", svc.CreateTransaction)
	routes.GET("/list", svc.GetUserTransaction)
	routes.DELETE("/:id", svc.DeleteTransactionByUser)

	return svc
}

func (svc *ServiceClient) CreateTransaction(ctx *gin.Context) {
	routes.CreateTransaction(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserTransaction(ctx *gin.Context) {
	routes.GetUserTransaction(ctx, svc.Client)
}

func (svc *ServiceClient) DeleteTransactionByUser(ctx *gin.Context) {
	routes.DeleteTransactionByUser(ctx, svc.Client)
}

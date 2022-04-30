package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/pos"
	"github.com/maslow123/api-gateway/pkg/transactions"
	"github.com/maslow123/api-gateway/pkg/users"
)

func main() {
	c, err := config.LoadConfig("./pkg/config/envs", "dev")

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	r := gin.Default()

	userService := *users.RegisterRoutes(r, &c)
	_ = pos.RegisterRoutes(r, &c, &userService)
	_ = transactions.RegisterRoutes(r, &c, &userService)

	r.Run(c.Port)
}

package subscription

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Logger  *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, logger *zap.Logger) {
	h := &Handler{Service: service, Logger: logger}
	fmt.Println(h)
	// rg.POST("/diary/link_doc/", h.LinkDoc)
}
